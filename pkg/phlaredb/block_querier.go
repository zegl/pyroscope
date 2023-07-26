package phlaredb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gogo/status"
	"github.com/google/pprof/profile"
	"github.com/grafana/dskit/multierror"
	"github.com/grafana/dskit/runutil"
	"github.com/oklog/ulid"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/samber/lo"
	"github.com/segmentio/parquet-go"
	"github.com/thanos-io/objstore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"

	ingestv1 "github.com/grafana/pyroscope/api/gen/proto/go/ingester/v1"
	typesv1 "github.com/grafana/pyroscope/api/gen/proto/go/types/v1"
	"github.com/grafana/pyroscope/pkg/iter"
	phlaremodel "github.com/grafana/pyroscope/pkg/model"
	phlareobj "github.com/grafana/pyroscope/pkg/objstore"
	phlarecontext "github.com/grafana/pyroscope/pkg/phlare/context"
	"github.com/grafana/pyroscope/pkg/phlaredb/block"
	"github.com/grafana/pyroscope/pkg/phlaredb/query"
	schemav1 "github.com/grafana/pyroscope/pkg/phlaredb/schemas/v1"
	"github.com/grafana/pyroscope/pkg/phlaredb/tsdb/index"
	"github.com/grafana/pyroscope/pkg/util"
)

const defaultBatchSize = 4096

type BlockQuerier struct {
	phlarectx context.Context
	logger    log.Logger

	bkt phlareobj.Bucket

	queriers     []*singleBlockQuerier
	queriersLock sync.RWMutex
}

func NewBlockQuerier(phlarectx context.Context, bucketReader phlareobj.Bucket) *BlockQuerier {
	return &BlockQuerier{
		phlarectx: contextWithBlockMetrics(phlarectx,
			newBlocksMetrics(
				phlarecontext.Registry(phlarectx),
			),
		),
		logger: phlarecontext.Logger(phlarectx),
		bkt:    bucketReader,
	}
}

func (b *BlockQuerier) Queriers() Queriers {
	b.queriersLock.RLock()
	defer b.queriersLock.RUnlock()

	res := make([]Querier, 0, len(b.queriers))
	for _, q := range b.queriers {
		res = append(res, q)
	}
	return res
}

func (b *BlockQuerier) BlockMetas(ctx context.Context) (metas []*block.Meta, _ error) {
	var names []ulid.ULID
	if err := b.bkt.Iter(ctx, "", func(n string) error {
		ulid, ok := block.IsBlockDir(n)
		if !ok {
			return nil
		}
		names = append(names, ulid)
		return nil
	}); err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(16)
	metas = make([]*block.Meta, len(names))
	for pos := range names {
		func(pos int) {
			g.Go(util.RecoverPanic(func() error {
				path := filepath.Join(names[pos].String(), block.MetaFilename)
				metaReader, err := b.bkt.Get(ctx, path)
				if err != nil {
					level.Error(b.logger).Log("msg", "error reading block meta", "block", path, "err", err)
					return nil
				}

				metas[pos], err = block.Read(metaReader)
				if err != nil {
					level.Error(b.logger).Log("msg", "error parsing block meta", "block", path, "err", err)
					return nil
				}
				return nil
			}))
		}(pos)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// sort slice and make sure nils are last
	sort.Slice(metas, func(i, j int) bool {
		if metas[i] == nil {
			return false
		}
		if metas[j] == nil {
			return true
		}
		return metas[i].MinTime < metas[j].MinTime
	})

	// iterate from the end and cut of till the first non-nil
	var pos int
	for pos = len(metas) - 1; pos >= 0; pos-- {
		if metas[pos] != nil {
			break
		}
	}

	return metas[0 : pos+1], nil
}

// Sync gradually scans the available blocks. If there are any changes to the
// last run it will Open/Close new/no longer existing ones.
func (b *BlockQuerier) Sync(ctx context.Context) error {
	observedMetas, err := b.BlockMetas(ctx)
	if err != nil {
		return err
	}

	// hold write lock to queriers
	b.queriersLock.Lock()

	// build lookup map

	querierByULID := make(map[ulid.ULID]*singleBlockQuerier)

	for pos := range b.queriers {
		querierByULID[b.queriers[pos].meta.ULID] = b.queriers[pos]
	}

	// ensure queries has the right length
	lenQueriers := len(observedMetas)
	if cap(b.queriers) < lenQueriers {
		b.queriers = make([]*singleBlockQuerier, lenQueriers)
	} else {
		b.queriers = b.queriers[:lenQueriers]
	}

	for pos, m := range observedMetas {

		q, ok := querierByULID[m.ULID]
		if ok {
			b.queriers[pos] = q
			delete(querierByULID, m.ULID)
			continue
		}

		b.queriers[pos] = NewSingleBlockQuerierFromMeta(b.phlarectx, b.bkt, m)
	}
	// ensure queriers are in ascending order.
	sort.Slice(b.queriers, func(i, j int) bool {
		return b.queriers[i].meta.MinTime < b.queriers[j].meta.MinTime
	})
	b.queriersLock.Unlock()

	// now close no longer available queries
	for _, q := range querierByULID {
		if err := q.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (b *BlockQuerier) AddBlockQuerierByMeta(m *block.Meta) {
	q := NewSingleBlockQuerierFromMeta(b.phlarectx, b.bkt, m)
	b.queriersLock.Lock()
	defer b.queriersLock.Unlock()
	i := sort.Search(len(b.queriers), func(i int) bool {
		return b.queriers[i].meta.MinTime >= m.MinTime
	})
	if i < len(b.queriers) && b.queriers[i].meta.ULID == m.ULID {
		// Block with this meta is already present, skipping.
		return
	}
	b.queriers = append(b.queriers, q) // Ensure we have enough capacity.
	copy(b.queriers[i+1:], b.queriers[i:])
	b.queriers[i] = q
}

// evict removes the block with the given ULID from the querier.
func (b *BlockQuerier) evict(blockID ulid.ULID) (bool, error) {
	b.queriersLock.Lock()
	// N.B: queriers are sorted by meta.MinTime.
	j := -1
	for i, q := range b.queriers {
		if q.meta.ULID.Compare(blockID) == 0 {
			j = i
			break
		}
	}
	if j < 0 {
		b.queriersLock.Unlock()
		return false, nil
	}
	blockQuerier := b.queriers[j]
	// Delete the querier from the slice and make it eligible for GC.
	copy(b.queriers[j:], b.queriers[j+1:])
	b.queriers[len(b.queriers)-1] = nil
	b.queriers = b.queriers[:len(b.queriers)-1]
	b.queriersLock.Unlock()
	return true, blockQuerier.Close()
}

func (b *BlockQuerier) Close() error {
	b.queriersLock.Lock()
	defer b.queriersLock.Unlock()

	errs := multierror.New()
	for pos := range b.queriers {
		if err := b.queriers[pos].Close(); err != nil {
			errs.Add(err)
		}
	}
	return errs.Err()
}

type singleBlockQuerier struct {
	logger  log.Logger
	metrics *blocksMetrics

	bkt  phlareobj.Bucket
	meta *block.Meta

	openLock sync.Mutex
	opened   bool
	index    *index.Reader
	profiles query.ParquetReader[*schemav1.Profile, *schemav1.ProfilePersister]
	symdb    symbolsReader
}

type symbolsReader interface {
	Open(context.Context, phlareobj.BucketReader) error
	io.Closer

	ResolveSymbols(context.Context, phlaremodel.StacktracesByPartition) (*ingestv1.MergeProfilesStacktracesResult, error)
	ResolvePprofSymbols(context.Context, phlaremodel.ProfileSampleByPartition) (*profile.Profile, error)
}

func NewSingleBlockQuerierFromMeta(phlarectx context.Context, bucketReader phlareobj.Bucket, meta *block.Meta) *singleBlockQuerier {
	q := &singleBlockQuerier{
		logger:  phlarecontext.Logger(phlarectx),
		metrics: contextBlockMetrics(phlarectx),

		bkt:  phlareobj.NewPrefixedBucket(bucketReader, meta.ULID.String()),
		meta: meta,
	}

	profilesFilePath := q.profiles.RelPath()
	for _, f := range meta.Files {
		if f.RelPath == profilesFilePath {
			q.profiles.SetSize(int64(f.SizeBytes))
		}
	}

	switch meta.Version {
	case block.MetaVersion1:
		q.symdb = newStacktraceResolverV1(q.bkt)
	case block.MetaVersion2:
		// TODO: Wrap
		// br := phlareobj.NewPrefixedBucket(q.bkt, symdb.DefaultDirName)
		// q.symdb = symdb.Open(phlarectx, br)
	default:
		panic(fmt.Errorf("unsupported block version %d", meta.Version))
	}

	return q
}

func (b *singleBlockQuerier) Profiles() []parquet.RowGroup {
	return b.profiles.File().RowGroups()
}

func (b *singleBlockQuerier) Index() IndexReader {
	return b.index
}

func (b *singleBlockQuerier) Symbols() SymbolsReader {
	return &inMemorySymbolsReader{
		partitions: make(map[uint64]*inMemorySymbolsResolver),
	}
}

func (b *singleBlockQuerier) Meta() block.Meta {
	if b.meta == nil {
		return block.Meta{}
	}
	return *b.meta
}

func (b *singleBlockQuerier) Close() error {
	b.openLock.Lock()
	defer func() {
		b.openLock.Unlock()
		b.metrics.blockOpened.Dec()
	}()
	errs := multierror.New()
	if b.index != nil {
		err := b.index.Close()
		b.index = nil
		if err != nil {
			errs.Add(err)
		}
	}
	if err := b.profiles.Close(); err != nil {
		errs.Add(err)
	}
	if b.symdb != nil {
		if err := b.symdb.Close(); err != nil {
			errs.Add(err)
		}
	}
	b.opened = false
	return errs.Err()
}

func (b *singleBlockQuerier) Bounds() (model.Time, model.Time) {
	return b.meta.MinTime, b.meta.MaxTime
}

type labelsInfo struct {
	fp  model.Fingerprint
	lbs phlaremodel.Labels
}

type Profile interface {
	StacktracePartition() uint64
	Timestamp() model.Time
	Fingerprint() model.Fingerprint
	Labels() phlaremodel.Labels
}

type Querier interface {
	Bounds() (model.Time, model.Time)
	SelectMatchingProfiles(ctx context.Context, params *ingestv1.SelectProfilesRequest) (iter.Iterator[Profile], error)
	MergeByStacktraces(ctx context.Context, rows iter.Iterator[Profile]) (*ingestv1.MergeProfilesStacktracesResult, error)
	MergeByLabels(ctx context.Context, rows iter.Iterator[Profile], by ...string) ([]*typesv1.Series, error)
	MergePprof(ctx context.Context, rows iter.Iterator[Profile]) (*profile.Profile, error)
	Open(ctx context.Context) error
	// Sorts profiles for retrieval.
	Sort([]Profile) []Profile
}

func InRange(q Querier, start, end model.Time) bool {
	min, max := q.Bounds()
	if start > max {
		return false
	}
	if end < min {
		return false
	}
	return true
}

type Queriers []Querier

func (queriers Queriers) Open(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(128)
	for _, q := range queriers {
		q := q
		g.Go(func() error {
			if err := q.Open(ctx); err != nil {
				return err
			}
			return nil
		})
	}
	return g.Wait()
}

func (queriers Queriers) SelectMatchingProfiles(ctx context.Context, params *ingestv1.SelectProfilesRequest) (iter.Iterator[Profile], error) {
	iters, err := SelectMatchingProfiles(ctx, params, queriers)
	if err != nil {
		return nil, err
	}
	return iter.NewMergeIterator(maxBlockProfile, true, iters...), nil
}

func (queriers Queriers) ForTimeRange(_ context.Context, start, end model.Time) (Queriers, error) {
	result := make(Queriers, 0, len(queriers))
	for _, q := range queriers {
		if InRange(q, start, end) {
			result = append(result, q)
		}
	}
	return result, nil
}

type BlockGetter func(ctx context.Context, start, end model.Time) (Queriers, error)

// SelectMatchingProfiles returns a list iterator of profiles matching the given request.
func SelectMatchingProfiles(ctx context.Context, request *ingestv1.SelectProfilesRequest, queriers Queriers) ([]iter.Iterator[Profile], error) {
	g, ctx := errgroup.WithContext(ctx)
	iters := make([]iter.Iterator[Profile], len(queriers))

	for i, querier := range queriers {
		i := i
		querier := querier
		g.Go(util.RecoverPanic(func() error {
			profiles, err := querier.SelectMatchingProfiles(ctx, request)
			if err != nil {
				return err
			}
			iters[i] = iter.NewBufferedIterator(profiles, 1024)
			return nil
		}))
	}

	if err := g.Wait(); err != nil {
		for _, it := range iters {
			if it != nil {
				runutil.CloseWithLogOnErr(util.Logger, it, "closing buffered iterator")
			}
		}
		return nil, err
	}
	return iters, nil
}

func MergeProfilesStacktraces(ctx context.Context, stream *connect.BidiStream[ingestv1.MergeProfilesStacktracesRequest, ingestv1.MergeProfilesStacktracesResponse], blockGetter BlockGetter) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "MergeProfilesStacktraces")
	defer sp.Finish()

	r, err := stream.Receive()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return connect.NewError(connect.CodeCanceled, errors.New("client closed stream"))
		}
		return err
	}

	if r.Request == nil {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("missing initial select request"))
	}
	request := r.Request
	sp.LogFields(
		otlog.String("start", model.Time(request.Start).Time().String()),
		otlog.String("end", model.Time(request.End).Time().String()),
		otlog.String("selector", request.LabelSelector),
		otlog.String("profile_id", request.Type.ID),
	)

	queriers, err := blockGetter(ctx, model.Time(request.Start), model.Time(request.End))
	if err != nil {
		return err
	}

	iters, err := SelectMatchingProfiles(ctx, request, queriers)
	if err != nil {
		return err
	}

	// send batches of profiles to client and filter via bidi stream.
	selectedProfiles, err := filterProfiles[
		BidiServerMerge[*ingestv1.MergeProfilesStacktracesResponse, *ingestv1.MergeProfilesStacktracesRequest],
		*ingestv1.MergeProfilesStacktracesResponse,
		*ingestv1.MergeProfilesStacktracesRequest](ctx, iters, defaultBatchSize, stream)
	if err != nil {
		return err
	}

	m := phlaremodel.NewStackTraceMerger()
	g, ctx := errgroup.WithContext(ctx)

	for i, querier := range queriers {
		querier := querier
		i := i
		if len(selectedProfiles[i]) == 0 {
			continue
		}
		// Sort profiles for better read locality.
		// Merge async the result so we can continue streaming profiles.
		g.Go(util.RecoverPanic(func() error {
			merge, err := querier.MergeByStacktraces(ctx, iter.NewSliceIterator(querier.Sort(selectedProfiles[i])))
			if err != nil {
				return err
			}
			m.MergeStackTraces(merge.Stacktraces, merge.FunctionNames)
			return nil
		}))
	}

	// Signals the end of the profile streaming by sending an empty response.
	// This allows the client to not block other streaming ingesters.
	sp.LogFields(otlog.String("msg", "signaling the end of the profile streaming"))
	if err = stream.Send(&ingestv1.MergeProfilesStacktracesResponse{}); err != nil {
		return err
	}

	if err = g.Wait(); err != nil {
		return err
	}

	// sends the final result to the client.
	sp.LogFields(otlog.String("msg", "sending the final result to the client"))
	err = stream.Send(&ingestv1.MergeProfilesStacktracesResponse{
		Result: &ingestv1.MergeProfilesStacktracesResult{
			Format:    ingestv1.StacktracesMergeFormat_MERGE_FORMAT_TREE,
			TreeBytes: m.TreeBytes(r.GetMaxNodes()),
		},
	})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return connect.NewError(connect.CodeCanceled, errors.New("client closed stream"))
		}
		return err
	}

	return nil
}

func MergeProfilesLabels(ctx context.Context, stream *connect.BidiStream[ingestv1.MergeProfilesLabelsRequest, ingestv1.MergeProfilesLabelsResponse], blockGetter BlockGetter) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "MergeProfilesLabels")
	defer sp.Finish()

	r, err := stream.Receive()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return connect.NewError(connect.CodeCanceled, errors.New("client closed stream"))
		}
		return err
	}

	if r.Request == nil {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("missing initial select request"))
	}
	request := r.Request
	by := r.By
	sort.Strings(by)
	sp.LogFields(
		otlog.String("start", model.Time(request.Start).Time().String()),
		otlog.String("end", model.Time(request.End).Time().String()),
		otlog.String("selector", request.LabelSelector),
		otlog.String("profile_id", request.Type.ID),
		otlog.String("by", strings.Join(by, ",")),
	)

	queriers, err := blockGetter(ctx, model.Time(request.Start), model.Time(request.End))
	if err != nil {
		return err
	}

	iters, err := SelectMatchingProfiles(ctx, request, queriers)
	if err != nil {
		return err
	}
	// send batches of profiles to client and filter via bidi stream.
	selectedProfiles, err := filterProfiles[
		BidiServerMerge[*ingestv1.MergeProfilesLabelsResponse, *ingestv1.MergeProfilesLabelsRequest],
		*ingestv1.MergeProfilesLabelsResponse,
		*ingestv1.MergeProfilesLabelsRequest](ctx, iters, defaultBatchSize, stream)
	if err != nil {
		return err
	}

	// Signals the end of the profile streaming by sending an empty request.
	// This allows the client to not block other streaming ingesters.
	if err := stream.Send(&ingestv1.MergeProfilesLabelsResponse{}); err != nil {
		return err
	}

	result := make([][]*typesv1.Series, 0, len(queriers))
	g, ctx := errgroup.WithContext(ctx)
	sync := lo.Synchronize()
	for i, querier := range queriers {
		i := i
		querier := querier
		if len(selectedProfiles[i]) == 0 {
			continue
		}
		// Sort profiles for better read locality.
		// And merge async the result for each queriers.
		g.Go(util.RecoverPanic(func() error {
			merge, err := querier.MergeByLabels(ctx,
				iter.NewSliceIterator(querier.Sort(selectedProfiles[i])),
				by...)
			if err != nil {
				return err
			}
			sync.Do(func() {
				result = append(result, merge)
			})

			return nil
		}))
	}

	if err := g.Wait(); err != nil {
		return err
	}

	// sends the final result to the client.
	err = stream.Send(&ingestv1.MergeProfilesLabelsResponse{
		Series: phlaremodel.SumSeries(result...),
	})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return connect.NewError(connect.CodeCanceled, errors.New("client closed stream"))
		}
		return err
	}

	return nil
}

func MergeProfilesPprof(ctx context.Context, stream *connect.BidiStream[ingestv1.MergeProfilesPprofRequest, ingestv1.MergeProfilesPprofResponse], blockGetter BlockGetter) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "MergeProfilesPprof")
	defer sp.Finish()

	r, err := stream.Receive()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return connect.NewError(connect.CodeCanceled, errors.New("client closed stream"))
		}
		return err
	}

	if r.Request == nil {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("missing initial select request"))
	}
	request := r.Request
	sp.LogFields(
		otlog.String("start", model.Time(request.Start).Time().String()),
		otlog.String("end", model.Time(request.End).Time().String()),
		otlog.String("selector", request.LabelSelector),
		otlog.String("profile_id", request.Type.ID),
	)

	queriers, err := blockGetter(ctx, model.Time(request.Start), model.Time(request.End))
	if err != nil {
		return err
	}

	iters, err := SelectMatchingProfiles(ctx, request, queriers)
	if err != nil {
		return err
	}

	// send batches of profiles to client and filter via bidi stream.
	selectedProfiles, err := filterProfiles[
		BidiServerMerge[*ingestv1.MergeProfilesPprofResponse, *ingestv1.MergeProfilesPprofRequest],
		*ingestv1.MergeProfilesPprofResponse,
		*ingestv1.MergeProfilesPprofRequest](ctx, iters, defaultBatchSize, stream)
	if err != nil {
		return err
	}

	result := make([]*profile.Profile, 0, len(queriers))
	var lock sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	for i, querier := range queriers {
		i := i
		querier := querier
		if len(selectedProfiles[i]) == 0 {
			continue
		}
		// Sort profiles for better read locality.
		// Merge async the result so we can continue streaming profiles.
		g.Go(util.RecoverPanic(func() error {
			merge, err := querier.MergePprof(ctx, iter.NewSliceIterator(querier.Sort(selectedProfiles[i])))
			if err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			result = append(result, merge)
			return nil
		}))
	}

	// Signals the end of the profile streaming by sending an empty response.
	// This allows the client to not block other streaming ingesters.
	if err := stream.Send(&ingestv1.MergeProfilesPprofResponse{}); err != nil {
		return err
	}

	if err := g.Wait(); err != nil {
		return err
	}
	if len(result) == 0 {
		result = append(result, &profile.Profile{})
	}
	for _, p := range result {
		phlaremodel.SetProfileMetadata(p, request.Type)
		p.TimeNanos = model.Time(r.Request.End).UnixNano()
	}
	p, err := profile.Merge(result)
	if err != nil {
		return err
	}

	// connect go already handles compression.
	var buf bytes.Buffer
	if err := p.WriteUncompressed(&buf); err != nil {
		return err
	}
	// sends the final result to the client.
	err = stream.Send(&ingestv1.MergeProfilesPprofResponse{
		Result: buf.Bytes(),
	})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return connect.NewError(connect.CodeCanceled, errors.New("client closed stream"))
		}
		return err
	}

	return nil
}

var maxBlockProfile Profile = BlockProfile{
	ts: model.Time(math.MaxInt64),
}

type BlockProfile struct {
	labels              phlaremodel.Labels
	fp                  model.Fingerprint
	ts                  model.Time
	stacktracePartition uint64
	RowNum              int64
}

func (p BlockProfile) StacktracePartition() uint64 {
	return p.stacktracePartition
}

func (p BlockProfile) RowNumber() int64 {
	return p.RowNum
}

func (p BlockProfile) Labels() phlaremodel.Labels {
	return p.labels
}

func (p BlockProfile) Timestamp() model.Time {
	return p.ts
}

func (p BlockProfile) Fingerprint() model.Fingerprint {
	return p.fp
}

func retrieveStacktracePartition(buf [][]parquet.Value, pos int) uint64 {
	if len(buf) > pos && len(buf[pos]) == 1 {
		return buf[pos][0].Uint64()
	}

	// return 0 stacktrace partition
	return uint64(0)
}

func (b *singleBlockQuerier) SelectMatchingProfiles(ctx context.Context, params *ingestv1.SelectProfilesRequest) (iter.Iterator[Profile], error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "SelectMatchingProfiles - Block")
	defer sp.Finish()
	if err := b.Open(ctx); err != nil {
		return nil, err
	}
	matchers, err := parser.ParseMetricSelector(params.LabelSelector)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse label selectors: "+err.Error())
	}
	if params.Type == nil {
		return nil, errors.New("no profileType given")
	}
	matchers = append(matchers, phlaremodel.SelectorFromProfileType(params.Type))

	postings, err := PostingsForMatchers(b.index, nil, matchers...)
	if err != nil {
		return nil, err
	}

	var (
		lbls       = make(phlaremodel.Labels, 0, 6)
		chks       = make([]index.ChunkMeta, 1)
		lblsPerRef = make(map[int64]labelsInfo)
	)

	// get all relevant labels/fingerprints
	for postings.Next() {
		fp, err := b.index.Series(postings.At(), &lbls, &chks)
		if err != nil {
			return nil, err
		}
		if lblsExisting, exists := lblsPerRef[int64(chks[0].SeriesIndex)]; exists {
			// Compare to check if there is a clash
			if phlaremodel.CompareLabelPairs(lbls, lblsExisting.lbs) != 0 {
				panic("label hash conflict")
			}
		} else {
			lblsPerRef[int64(chks[0].SeriesIndex)] = labelsInfo{
				fp:  model.Fingerprint(fp),
				lbs: lbls,
			}
			lbls = make(phlaremodel.Labels, 0, 6)
		}
	}

	var buf [][]parquet.Value

	pIt := query.NewBinaryJoinIterator(
		0,
		b.profiles.ColumnIter(ctx, "SeriesIndex", query.NewMapPredicate(lblsPerRef), "SeriesIndex"),
		b.profiles.ColumnIter(ctx, "TimeNanos", query.NewIntBetweenPredicate(model.Time(params.Start).UnixNano(), model.Time(params.End).UnixNano()), "TimeNanos"),
	)

	if b.meta.Version >= 2 {
		pIt = query.NewBinaryJoinIterator(
			0,
			pIt,
			b.profiles.ColumnIter(ctx, "StacktracePartition", nil, "StacktracePartition"),
		)
		buf = make([][]parquet.Value, 3)
	} else {
		buf = make([][]parquet.Value, 2)
	}

	iters := make([]iter.Iterator[Profile], 0, len(lblsPerRef))
	defer pIt.Close()

	currSeriesIndex := int64(-1)
	var currentSeriesSlice []Profile
	for pIt.Next() {
		res := pIt.At()
		buf = res.Columns(buf, "SeriesIndex", "TimeNanos", "StacktracePartition")
		seriesIndex := buf[0][0].Int64()
		if seriesIndex != currSeriesIndex {
			currSeriesIndex = seriesIndex
			if len(currentSeriesSlice) > 0 {
				iters = append(iters, iter.NewSliceIterator(currentSeriesSlice))
			}
			currentSeriesSlice = make([]Profile, 0, 100)
		}

		currentSeriesSlice = append(currentSeriesSlice, BlockProfile{
			labels:              lblsPerRef[seriesIndex].lbs,
			fp:                  lblsPerRef[seriesIndex].fp,
			ts:                  model.TimeFromUnixNano(buf[1][0].Int64()),
			stacktracePartition: retrieveStacktracePartition(buf, 2),
			RowNum:              res.RowNumber[0],
		})
	}
	if len(currentSeriesSlice) > 0 {
		iters = append(iters, iter.NewSliceIterator(currentSeriesSlice))
	}

	return iter.NewMergeIterator(maxBlockProfile, false, iters...), nil
}

func (b *singleBlockQuerier) Sort(in []Profile) []Profile {
	// Sort by RowNumber to avoid seeking back and forth in the file.
	sort.Slice(in, func(i, j int) bool {
		return in[i].(BlockProfile).RowNum < in[j].(BlockProfile).RowNum
	})
	return in
}

func newByteSliceFromBucketReader(ctx context.Context, bucketReader objstore.BucketReader, path string) (index.RealByteSlice, error) {
	f, err := bucketReader.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (q *singleBlockQuerier) Open(ctx context.Context) error {
	q.openLock.Lock()
	defer q.openLock.Unlock()

	// already open
	if q.opened {
		return nil
	}
	if err := q.openFiles(ctx); err != nil {
		return err
	}
	q.metrics.blockOpened.Inc()
	q.opened = true
	return nil
}

// openFiles opens the parquet and tsdb files so they are ready for usage.
func (q *singleBlockQuerier) openFiles(ctx context.Context) error {
	start := time.Now()
	sp, ctx := opentracing.StartSpanFromContext(ctx, "BlockQuerier - open")
	defer func() {
		q.metrics.blockOpeningLatency.Observe(time.Since(start).Seconds())
		sp.LogFields(
			otlog.String("block_ulid", q.meta.ULID.String()),
		)
		sp.Finish()
	}()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(util.RecoverPanic(func() error {
		// open tsdb index
		indexBytes, err := newByteSliceFromBucketReader(ctx, q.bkt, block.IndexFilename)
		if err != nil {
			return errors.Wrap(err, "error reading tsdb index")
		}

		q.index, err = index.NewReader(indexBytes)
		if err != nil {
			return errors.Wrap(err, "opening tsdb index")
		}
		return nil
	}))

	// open parquet files
	g.Go(util.RecoverPanic(func() error {
		return q.profiles.Open(query.AddMetricsToContext(ctx, q.metrics.query), q.bkt)
	}))
	g.Go(util.RecoverPanic(func() error {
		return q.symdb.Open(query.AddMetricsToContext(ctx, q.metrics.query), q.bkt)
	}))

	return g.Wait()
}
