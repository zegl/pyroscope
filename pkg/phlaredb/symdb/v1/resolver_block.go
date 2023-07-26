package v1

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/google/pprof/profile"
	"github.com/grafana/dskit/runutil"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/segmentio/parquet-go"

	googlev1 "github.com/grafana/pyroscope/api/gen/proto/go/google/v1"
	ingestv1 "github.com/grafana/pyroscope/api/gen/proto/go/ingester/v1"
	"github.com/grafana/pyroscope/pkg/iter"
	phlaremodel "github.com/grafana/pyroscope/pkg/model"
	phlareobj "github.com/grafana/pyroscope/pkg/objstore"
	parquetobj "github.com/grafana/pyroscope/pkg/objstore/parquet"
	"github.com/grafana/pyroscope/pkg/phlaredb/block"
	"github.com/grafana/pyroscope/pkg/phlaredb/query"
	schemav1 "github.com/grafana/pyroscope/pkg/phlaredb/schemas/v1"
	"github.com/grafana/pyroscope/pkg/phlaredb/symdb"
	"github.com/grafana/pyroscope/pkg/util"
)

type SymbolsResolver struct {
	stacktraces stacktraces
	locations   inMemoryparquetReader[*schemav1.InMemoryLocation, *schemav1.LocationPersister]
	mappings    inMemoryparquetReader[*schemav1.InMemoryMapping, *schemav1.MappingPersister]
	functions   inMemoryparquetReader[*schemav1.InMemoryFunction, *schemav1.FunctionPersister]
	strings     inMemoryparquetReader[string, *schemav1.StringPersister]
}

type stacktraces interface {
	Open(ctx context.Context) error
	Close() error

	// Load the database into memory entirely.
	// This method is used at compaction.
	Load(context.Context) error
	WriteStats(partition uint64, s *symdb.Stats)

	Resolve(ctx context.Context, partition uint64, locs symdb.StacktraceInserter, stacktraceIDs []uint32) error
}

// TODO: Copy from block_querier.go

func NewSymbolsResolver() *SymbolsResolver {
	return new(SymbolsResolver)
}

func (r *SymbolsResolver) Open(ctx context.Context, reader phlareobj.BucketReader) error {
	return nil
}

func (r *SymbolsResolver) Close() error {
	return nil
}

func (r *SymbolsResolver) resolveLocations(ctx context.Context, mapping uint64, locs locationsIdsByStacktraceID, stacktraceIDs []uint32) error {
	sort.Slice(stacktraceIDs, func(i, j int) bool {
		return stacktraceIDs[i] < stacktraceIDs[j]
	})
	return r.stacktraces.Resolve(ctx, mapping, locs, stacktraceIDs)
}

func (r *SymbolsResolver) resolvePprofSymbols(ctx context.Context, profileSampleByMapping phlaremodel.ProfileSampleByPartition) (*profile.Profile, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "ResolvePprofSymbols - Block")
	defer sp.Finish()

	locationsIdsByStacktraceID := newLocationsIdsByStacktraceID(len(profileSampleByMapping) * 1024)

	// gather stacktraces
	if err := profileSampleByMapping.ForEach(func(mapping uint64, samples phlaremodel.ProfileSampleMap) error {
		stacktraceIDs := samples.Ids()
		sp.LogFields(
			otlog.Int("stacktraces", len(stacktraceIDs)),
			otlog.Uint64("mapping", mapping),
		)
		return r.resolveLocations(ctx, mapping, locationsIdsByStacktraceID, stacktraceIDs)
	}); err != nil {
		return nil, err
	}

	// gather locations
	var (
		functionIDs         = newUniqueIDs[struct{}]()
		mappingIDs          = newUniqueIDs[lo.Tuple2[*profile.Mapping, *googlev1.Mapping]]()
		locations           = r.locations.retrieveRows(ctx, locationsIdsByStacktraceID.locationIds().iterator())
		locationModelsByIds = map[uint64]*profile.Location{}
		functionModelsByIds = map[uint32]*profile.Function{}
	)
	for locations.Next() {
		s := locations.At()
		m, ok := mappingIDs[int64(s.Result.MappingId)]
		if !ok {
			m = lo.T2(&profile.Mapping{
				ID: uint64(s.Result.MappingId),
			}, &googlev1.Mapping{
				Id: uint64(s.Result.MappingId),
			})
			mappingIDs[int64(s.Result.MappingId)] = m
		}
		loc := &profile.Location{
			ID:       s.Result.Id,
			Address:  s.Result.Address,
			IsFolded: s.Result.IsFolded,
			Mapping:  m.A,
		}
		for _, line := range s.Result.Line {
			functionIDs[int64(line.FunctionId)] = struct{}{}
			fn, ok := functionModelsByIds[line.FunctionId]
			if !ok {
				fn = &profile.Function{
					ID: uint64(line.FunctionId),
				}
				functionModelsByIds[line.FunctionId] = fn
			}

			loc.Line = append(loc.Line, profile.Line{
				Line:     int64(line.Line),
				Function: fn,
			})
		}
		locationModelsByIds[uint64(s.RowNum)] = loc
	}
	if err := locations.Err(); err != nil {
		return nil, err
	}

	// gather functions
	var (
		stringsIds    = newUniqueIDs[int64]()
		functions     = r.functions.retrieveRows(ctx, functionIDs.iterator())
		functionsById = map[int64]*googlev1.Function{}
	)
	for functions.Next() {
		s := functions.At()
		functionsById[int64(s.Result.Id)] = &googlev1.Function{
			Id:         s.Result.Id,
			Name:       int64(s.Result.Name),
			SystemName: int64(s.Result.SystemName),
			Filename:   int64(s.Result.Filename),
			StartLine:  int64(s.Result.StartLine),
		}
		stringsIds[int64(s.Result.Name)] = 0
		stringsIds[int64(s.Result.Filename)] = 0
		stringsIds[int64(s.Result.SystemName)] = 0
	}
	if err := functions.Err(); err != nil {
		return nil, err
	}
	// gather mapping
	mapping := r.mappings.retrieveRows(ctx, mappingIDs.iterator())
	for mapping.Next() {
		cur := mapping.At()
		m := mappingIDs[int64(cur.Result.Id)]
		m.B.Filename = int64(cur.Result.Filename)
		m.B.BuildId = int64(cur.Result.BuildId)
		m.A.Start = cur.Result.MemoryStart
		m.A.Limit = cur.Result.MemoryLimit
		m.A.Offset = cur.Result.FileOffset
		m.A.HasFunctions = cur.Result.HasFunctions
		m.A.HasFilenames = cur.Result.HasFilenames
		m.A.HasLineNumbers = cur.Result.HasLineNumbers
		m.A.HasInlineFrames = cur.Result.HasInlineFrames

		stringsIds[int64(cur.Result.Filename)] = 0
		stringsIds[int64(cur.Result.BuildId)] = 0
	}
	// gather strings
	var (
		names   = make([]string, len(stringsIds))
		strings = r.strings.retrieveRows(ctx, stringsIds.iterator())
		idx     = int64(0)
	)
	for strings.Next() {
		s := strings.At()
		names[idx] = s.Result
		stringsIds[s.RowNum] = idx
		idx++
	}
	if err := strings.Err(); err != nil {
		return nil, err
	}

	for _, model := range mappingIDs {
		model.A.File = names[stringsIds[model.B.Filename]]
		model.A.BuildID = names[stringsIds[model.B.BuildId]]
	}

	mappingResult := make([]*profile.Mapping, 0, len(mappingIDs))
	for _, model := range mappingIDs {
		mappingResult = append(mappingResult, model.A)
	}

	_ = profileSampleByMapping.ForEach(func(_ uint64, samples phlaremodel.ProfileSampleMap) error {
		for id, model := range samples {
			locsId := locationsIdsByStacktraceID.byStacktraceID[int64(id)]
			model.Location = make([]*profile.Location, len(locsId))
			for i, locId := range locsId {
				model.Location[i] = locationModelsByIds[uint64(locId)]
			}
			// todo labels.
		}
		return nil
	})

	for id, model := range functionModelsByIds {
		fn := functionsById[int64(id)]
		model.Name = names[stringsIds[fn.Name]]
		model.Filename = names[stringsIds[fn.Filename]]
		model.SystemName = names[stringsIds[fn.SystemName]]
		model.StartLine = fn.StartLine
	}
	result := &profile.Profile{
		Sample:   profileSampleByMapping.StacktraceSamples(),
		Location: lo.Values(locationModelsByIds),
		Function: lo.Values(functionModelsByIds),
		Mapping:  mappingResult,
	}
	normalizeProfileIds(result)

	return result, nil
}

func (r *SymbolsResolver) resolveSymbols(ctx context.Context, stacktracesByMapping phlaremodel.StacktracesByPartition) (*ingestv1.MergeProfilesStacktracesResult, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "ResolveSymbols - Block")
	defer sp.Finish()
	locationsIdsByStacktraceID := newLocationsIdsByStacktraceID(len(stacktracesByMapping) * 1024)

	// gather stacktraces
	if err := stacktracesByMapping.ForEach(func(mapping uint64, samples phlaremodel.StacktraceSampleMap) error {
		stacktraceIDs := samples.Ids()
		sp.LogFields(
			otlog.Int("stacktraces", len(stacktraceIDs)),
			otlog.Uint64("mapping", mapping),
		)
		return r.resolveLocations(ctx, mapping, locationsIdsByStacktraceID, stacktraceIDs)
	}); err != nil {
		return nil, err
	}

	sp.LogFields(otlog.Int("locationIDs", len(locationsIdsByStacktraceID.locationIds())))

	// gather locations
	sp.LogFields(otlog.String("msg", "gather locations"))
	var (
		locationIDsByFunctionID = newUniqueIDs[[]int64]()
		locations               = r.locations.retrieveRows(ctx, locationsIdsByStacktraceID.locationIds().iterator())
	)
	for locations.Next() {
		s := locations.At()

		for _, line := range s.Result.Line {
			locationIDsByFunctionID[int64(line.FunctionId)] = append(locationIDsByFunctionID[int64(line.FunctionId)], s.RowNum)
		}
	}
	if err := locations.Err(); err != nil {
		return nil, err
	}
	sp.LogFields(otlog.Int("functions", len(locationIDsByFunctionID)))

	// gather functions
	sp.LogFields(otlog.String("msg", "gather functions"))
	var (
		functionIDsByStringID = newUniqueIDs[[]int64]()
		functions             = r.functions.retrieveRows(ctx, locationIDsByFunctionID.iterator())
	)
	for functions.Next() {
		s := functions.At()

		functionIDsByStringID[int64(s.Result.Name)] = append(functionIDsByStringID[int64(s.Result.Name)], s.RowNum)
	}
	if err := functions.Err(); err != nil {
		return nil, err
	}

	// gather strings
	sp.LogFields(otlog.String("msg", "gather strings"))
	var (
		names   = make([]string, len(functionIDsByStringID))
		idSlice = make([][]int64, len(functionIDsByStringID))
		strings = r.strings.retrieveRows(ctx, functionIDsByStringID.iterator())
		idx     = 0
	)
	for strings.Next() {
		s := strings.At()
		names[idx] = s.Result
		idSlice[idx] = []int64{s.RowNum}
		idx++
	}
	if err := strings.Err(); err != nil {
		return nil, err
	}

	sp.LogFields(otlog.String("msg", "build MergeProfilesStacktracesResult"))
	// idSlice contains stringIDs and gets rewritten into functionIDs
	for nameID := range idSlice {
		var functionIDs []int64
		for _, stringID := range idSlice[nameID] {
			functionIDs = append(functionIDs, functionIDsByStringID[stringID]...)
		}
		idSlice[nameID] = functionIDs
	}

	// idSlice contains functionIDs and gets rewritten into locationIDs
	for nameID := range idSlice {
		var locationIDs []int64
		for _, functionID := range idSlice[nameID] {
			locationIDs = append(locationIDs, locationIDsByFunctionID[functionID]...)
		}
		idSlice[nameID] = locationIDs
	}

	// write a map locationID two nameID
	nameIDbyLocationID := make(map[int64]int32)
	for nameID := range idSlice {
		for _, locationID := range idSlice[nameID] {
			nameIDbyLocationID[locationID] = int32(nameID)
		}
	}
	_ = stacktracesByMapping.ForEach(func(_ uint64, stacktraceSamples phlaremodel.StacktraceSampleMap) error {
		// write correct string ID into each sample
		for stacktraceID, samples := range stacktraceSamples {
			locationIDs := locationsIdsByStacktraceID.byStacktraceID[int64(stacktraceID)]

			functionIDs := make([]int32, len(locationIDs))
			for idx := range functionIDs {
				functionIDs[idx] = nameIDbyLocationID[int64(locationIDs[idx])]
			}
			samples.FunctionIds = functionIDs
		}
		return nil
	})

	return &ingestv1.MergeProfilesStacktracesResult{
		Stacktraces:   stacktracesByMapping.StacktraceSamples(),
		FunctionNames: names,
	}, nil
}

type stacktraceResolverV1 struct {
	stacktraces  query.ParquetReader[*schemav1.Stacktrace, *schemav1.StacktracePersister]
	bucketReader phlareobj.Bucket
}

func (r *stacktraceResolverV1) Open(ctx context.Context) error {
	return r.stacktraces.Open(ctx, r.bucketReader)
}

func (r *stacktraceResolverV1) Close() error {
	return r.stacktraces.Close()
}

func (r *stacktraceResolverV1) Resolve(ctx context.Context, _ uint64, locs symdb.StacktraceInserter, stacktraceIDs []uint32) error {
	stacktraces := query.RepeatedColumnIter(ctx, r.stacktraces.File(), "LocationIDs.list.element", iter.NewSliceIterator(stacktraceIDs))
	defer stacktraces.Close()
	t := make([]int32, 0, 64)
	for stacktraces.Next() {
		s := stacktraces.At()
		t = grow(t, len(s.Values))
		for i, v := range s.Values {
			t[i] = v.Int32()
		}
		locs.InsertStacktrace(s.Row, t)
	}
	return stacktraces.Err()
}

func (r *stacktraceResolverV1) WriteStats(_ uint64, s *symdb.Stats) {
	s.StacktracesTotal = int(r.stacktraces.File().NumRows())
	s.MaxStacktraceID = s.StacktracesTotal
}

func (r *stacktraceResolverV1) Load(context.Context) error {
	// FIXME(kolesnikovae): Loading all stacktraces from parquet file
	//  into memory is likely a bad choice. Instead we could convert
	//  it to symdb first.
	return nil
}

type stacktraceResolverV2 struct {
	reader       *symdb.Reader
	bucketReader phlareobj.Bucket
}

func (r *stacktraceResolverV2) Open(ctx context.Context) error {
	if r.reader != nil {
		return nil
	}
	var err error
	r.reader, err = symdb.Open(ctx, r.bucketReader)
	return err
}

func (r *stacktraceResolverV2) Close() error {
	return nil
}

func (r *stacktraceResolverV2) Resolve(ctx context.Context, partition uint64, locs symdb.StacktraceInserter, stacktraceIDs []uint32) error {
	mr, ok := r.reader.SymbolsResolver(partition)
	if !ok {
		return nil
	}
	resolver := mr.StacktraceResolver()
	defer resolver.Release()
	return resolver.ResolveStacktraces(ctx, locs, stacktraceIDs)
}

func (r *stacktraceResolverV2) Load(ctx context.Context) error {
	return r.reader.Load(ctx)
}

func (r *stacktraceResolverV2) WriteStats(partition uint64, s *symdb.Stats) {
	mr, ok := r.reader.SymbolsResolver(partition)
	if ok {
		mr.WriteStats(s)
	}
}

func newStacktraceResolverV1(bucketReader phlareobj.Bucket, meta *block.Meta) stacktraces {
	q := &stacktraceResolverV1{
		bucketReader: bucketReader,
	}
	for _, f := range meta.Files {
		switch f.RelPath {
		case q.stacktraces.RelPath():
			q.stacktraces.SetSize(int64(f.SizeBytes))
		}
	}
	return q
}

func newStacktraceResolverV2(bucketReader phlareobj.Bucket) stacktraces {
	return &stacktraceResolverV2{
		bucketReader: bucketReader,
	}
}

type locationsIdsByStacktraceID struct {
	byStacktraceID map[int64][]int32
	ids            uniqueIDs[struct{}]
}

func newLocationsIdsByStacktraceID(size int) locationsIdsByStacktraceID {
	return locationsIdsByStacktraceID{
		byStacktraceID: make(map[int64][]int32, size),
		ids:            newUniqueIDs[struct{}](),
	}
}

func (l locationsIdsByStacktraceID) InsertStacktrace(stacktraceID uint32, locs []int32) {
	s := make([]int32, len(locs))
	l.byStacktraceID[int64(stacktraceID)] = s
	for i, locationID := range locs {
		l.ids[int64(locationID)] = struct{}{}
		s[i] = locationID
	}
}

func (l locationsIdsByStacktraceID) locationIds() uniqueIDs[struct{}] {
	return l.ids
}

type uniqueIDs[T any] map[int64]T

func newUniqueIDs[T any]() uniqueIDs[T] {
	return uniqueIDs[T](make(map[int64]T))
}

func (m uniqueIDs[T]) iterator() iter.Iterator[int64] {
	ids := lo.Keys(m)
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	return iter.NewSliceIterator(ids)
}

type ResultWithRowNum[M any] struct {
	Result M
	RowNum int64
}

type inMemoryparquetReader[M schemav1.Models, P schemav1.Persister[M]] struct {
	persister P
	file      *parquet.File
	size      int64
	reader    phlareobj.ReaderAtCloser
	cache     []M
}

func (r *inMemoryparquetReader[M, P]) open(ctx context.Context, bucketReader phlareobj.BucketReader) error {
	filePath := r.persister.Name() + block.ParquetSuffix

	if r.size == 0 {
		attrs, err := bucketReader.Attributes(ctx, filePath)
		if err != nil {
			return errors.Wrapf(err, "getting attributes for '%s'", filePath)
		}
		r.size = attrs.Size
	}
	ra, err := bucketReader.ReaderAt(ctx, filePath)
	if err != nil {
		return errors.Wrapf(err, "create reader '%s'", filePath)
	}
	ra = parquetobj.NewOptimizedReader(ra)

	r.reader = ra

	// first try to open file, this is required otherwise OpenFile panics
	parquetFile, err := parquet.OpenFile(ra, r.size, parquet.SkipPageIndex(true), parquet.SkipBloomFilters(true))
	if err != nil {
		return errors.Wrapf(err, "opening parquet file '%s'", filePath)
	}
	if parquetFile.NumRows() == 0 {
		return fmt.Errorf("error parquet file '%s' contains no rows", filePath)
	}
	opts := []parquet.FileOption{
		parquet.SkipBloomFilters(true), // we don't use bloom filters
		parquet.FileReadMode(parquet.ReadModeAsync),
		parquet.ReadBufferSize(2 << 20),
	}
	// now open it for real
	r.file, err = parquet.OpenFile(ra, r.size, opts...)
	if err != nil {
		return errors.Wrapf(err, "opening parquet file '%s'", filePath)
	}

	// read all rows into memory
	r.cache = make([]M, r.file.NumRows())
	var offset int64
	for _, rg := range r.file.RowGroups() {
		rows := rg.NumRows()
		dst := r.cache[offset : offset+rows]
		offset += rows
		if err = r.readRG(dst, rg); err != nil {
			return errors.Wrapf(err, "reading row group from parquet file '%s'", filePath)
		}
	}
	err = r.reader.Close()
	r.reader = nil
	r.file = nil
	return err
}

// parquet.CopyRows uses hardcoded buffer size:
// defaultRowBufferSize = 42
const inMemoryReaderRowsBufSize = 1 << 10

func (r *inMemoryparquetReader[M, P]) readRG(dst []M, rg parquet.RowGroup) (err error) {
	rr := parquet.NewRowGroupReader(rg)
	defer runutil.CloseWithLogOnErr(util.Logger, rr, "closing parquet row group reader")
	buf := make([]parquet.Row, inMemoryReaderRowsBufSize)
	for i := 0; i < len(dst); {
		n, err := rr.ReadRows(buf)
		if n > 0 {
			for _, row := range buf[:n] {
				_, v, err := r.persister.Reconstruct(row)
				if err != nil {
					return err
				}
				dst[i] = v
				i++
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
	return nil
}

func (r *inMemoryparquetReader[M, P]) Close() error {
	if r.reader != nil {
		return r.reader.Close()
	}
	r.reader = nil
	r.file = nil
	r.cache = nil
	return nil
}

func (r *inMemoryparquetReader[M, P]) relPath() string {
	return r.persister.Name() + block.ParquetSuffix
}

func (r *inMemoryparquetReader[M, P]) retrieveRows(_ context.Context, rowNumIterator iter.Iterator[int64]) iter.Iterator[ResultWithRowNum[M]] {
	return &cacheIterator[M]{
		cache:          r.cache,
		rowNumIterator: rowNumIterator,
	}
}

type cacheIterator[M any] struct {
	cache          []M
	rowNumIterator iter.Iterator[int64]
}

func (c *cacheIterator[M]) Next() bool {
	if !c.rowNumIterator.Next() {
		return false
	}
	if c.rowNumIterator.At() >= int64(len(c.cache)) {
		return false
	}
	return true
}

func (c *cacheIterator[M]) At() ResultWithRowNum[M] {
	return ResultWithRowNum[M]{
		Result: c.cache[c.rowNumIterator.At()],
		RowNum: c.rowNumIterator.At(),
	}
}

func (c *cacheIterator[M]) Err() error {
	return nil
}

func (c *cacheIterator[M]) Close() error {
	return nil
}
