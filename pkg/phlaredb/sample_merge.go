package phlaredb

import (
	"context"
	"sort"

	"github.com/google/pprof/profile"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/model"
	"github.com/samber/lo"

	ingestv1 "github.com/grafana/pyroscope/api/gen/proto/go/ingester/v1"
	typesv1 "github.com/grafana/pyroscope/api/gen/proto/go/types/v1"
	"github.com/grafana/pyroscope/pkg/iter"
	phlaremodel "github.com/grafana/pyroscope/pkg/model"
	"github.com/grafana/pyroscope/pkg/phlaredb/query"
)

func (b *singleBlockQuerier) MergeByStacktraces(ctx context.Context, rows iter.Iterator[Profile]) (*ingestv1.MergeProfilesStacktracesResult, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "MergeByStacktraces - Block")
	defer sp.Finish()
	stacktraceAggrValues := make(phlaremodel.StacktracesByPartition)
	if err := mergeByStacktraces(ctx, b.profiles.File(), rows, stacktraceAggrValues); err != nil {
		return nil, err
	}
	return b.symdb.ResolveSymbols(ctx, stacktraceAggrValues)
}

func (b *singleBlockQuerier) MergePprof(ctx context.Context, rows iter.Iterator[Profile]) (*profile.Profile, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "MergeByStacktraces - Block")
	defer sp.Finish()
	stacktraceAggrValues := make(phlaremodel.ProfileSampleByPartition)
	if err := mergeByStacktraces(ctx, b.profiles.File(), rows, stacktraceAggrValues); err != nil {
		return nil, err
	}
	return b.symdb.ResolvePprofSymbols(ctx, stacktraceAggrValues)
}

type mapAdder interface {
	Add(partition uint64, key uint32, value int64)
}

func mergeByStacktraces(ctx context.Context, profileSource query.Source, rows iter.Iterator[Profile], m mapAdder) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "mergeByStacktraces")
	defer sp.Finish()
	// clone the rows to be able to iterate over them twice
	multiRows, err := iter.CloneN(rows, 2)
	if err != nil {
		return err
	}
	it := query.NewMultiRepeatedPageIterator(
		query.RepeatedColumnIter(ctx, profileSource, "Samples.list.element.StacktraceID", multiRows[0]),
		query.RepeatedColumnIter(ctx, profileSource, "Samples.list.element.Value", multiRows[1]),
	)
	defer it.Close()

	for it.Next() {
		values := it.At().Values
		for i := 0; i < len(values[0]); i++ {
			m.Add(it.At().Row.StacktracePartition(), uint32(values[0][i].Int64()), values[1][i].Int64())
		}
	}
	return nil
}

func (b *singleBlockQuerier) MergeByLabels(ctx context.Context, rows iter.Iterator[Profile], by ...string) ([]*typesv1.Series, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "MergeByLabels - Block")
	defer sp.Finish()

	m := make(seriesByLabels)
	columnName := "TotalValue"
	if b.meta.Version == 1 {
		columnName = "Samples.list.element.Value"
	}
	if err := mergeByLabels(ctx, b.profiles.File(), columnName, rows, m, by...); err != nil {
		return nil, err
	}
	return m.normalize(), nil
}

type seriesByLabels map[string]*typesv1.Series

func (m seriesByLabels) normalize() []*typesv1.Series {
	result := lo.Values(m)
	sort.Slice(result, func(i, j int) bool {
		return phlaremodel.CompareLabelPairs(result[i].Labels, result[j].Labels) < 0
	})
	// we have to sort the points in each series because labels reduction may have changed the order
	for _, s := range result {
		sort.Slice(s.Points, func(i, j int) bool {
			return s.Points[i].Timestamp < s.Points[j].Timestamp
		})
	}
	return result
}

func mergeByLabels(ctx context.Context, profileSource query.Source, columnName string, rows iter.Iterator[Profile], m seriesByLabels, by ...string) error {
	it := query.RepeatedColumnIter(ctx, profileSource, columnName, rows)

	defer it.Close()

	labelsByFingerprint := map[model.Fingerprint]string{}
	labelBuf := make([]byte, 0, 1024)

	for it.Next() {
		values := it.At()
		p := values.Row
		var total int64
		for _, e := range values.Values {
			total += e.Int64()
		}
		labelsByString, ok := labelsByFingerprint[p.Fingerprint()]
		if !ok {
			labelBuf = p.Labels().BytesWithLabels(labelBuf, by...)
			labelsByString = string(labelBuf)
			labelsByFingerprint[p.Fingerprint()] = labelsByString
			if _, ok := m[labelsByString]; !ok {
				m[labelsByString] = &typesv1.Series{
					Labels: p.Labels().WithLabels(by...),
					Points: []*typesv1.Point{
						{
							Timestamp: int64(p.Timestamp()),
							Value:     float64(total),
						},
					},
				}
				continue
			}
		}
		series := m[labelsByString]
		series.Points = append(series.Points, &typesv1.Point{
			Timestamp: int64(p.Timestamp()),
			Value:     float64(total),
		})
	}
	return it.Err()
}
