package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/segmentio/parquet-go"

	"github.com/grafana/pyroscope/pkg/iter"
	"github.com/grafana/pyroscope/pkg/objstore"
	parquet2 "github.com/grafana/pyroscope/pkg/objstore/parquet"
	"github.com/grafana/pyroscope/pkg/phlaredb/block"
	"github.com/grafana/pyroscope/pkg/phlaredb/schemas/v1"
)

type ParquetReader[M v1.Models, P v1.PersisterName] struct {
	persister P
	file      *parquet.File
	reader    objstore.ReaderAtCloser
	size      int64
	metrics   *Metrics
}

const parquetReadBufferSize = 2 * 1024 * 1024 // 2MB

func (r *ParquetReader[M, P]) Open(ctx context.Context, bucketReader objstore.BucketReader) error {
	r.metrics = GetMetricsFromContext(ctx)
	filePath := r.persister.Name() + block.ParquetSuffix

	if r.size == 0 {
		attrs, err := bucketReader.Attributes(ctx, filePath)
		if err != nil {
			return errors.Wrapf(err, "getting attributes for '%s'", filePath)
		}
		r.size = attrs.Size
	}
	// the same reader is used to serve all requests, so we pass context.Background() here
	ra, err := bucketReader.ReaderAt(context.Background(), filePath)
	if err != nil {
		return errors.Wrapf(err, "create reader '%s'", filePath)
	}
	ra = parquet2.NewOptimizedReader(ra)
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
		parquet.ReadBufferSize(parquetReadBufferSize),
	}
	// now open it for real
	r.file, err = parquet.OpenFile(ra, r.size, opts...)
	if err != nil {
		return errors.Wrapf(err, "opening parquet file '%s'", filePath)
	}

	return nil
}

func (r *ParquetReader[M, P]) Close() error {
	if r.reader != nil {
		return r.reader.Close()
	}
	r.reader = nil
	r.file = nil
	return nil
}

func (r *ParquetReader[M, P]) File() *parquet.File { return r.file }

func (r *ParquetReader[M, P]) SetSize(s int64) { r.size = s }

func (r *ParquetReader[M, P]) RelPath() string {
	return r.persister.Name() + block.ParquetSuffix
}

func (r *ParquetReader[M, P]) ColumnIter(ctx context.Context, columnName string, predicate Predicate, alias string) Iterator {
	index, _ := GetColumnIndexByPath(r.file, columnName)
	if index == -1 {
		return NewErrIterator(fmt.Errorf("column '%s' not found in parquet file '%s'", columnName, r.RelPath()))
	}
	ctx = AddMetricsToContext(ctx, r.metrics)
	return NewSyncIterator(ctx, r.file.RowGroups(), index, columnName, 1000, predicate, alias)
}

type Source interface {
	Schema() *parquet.Schema
	RowGroups() []parquet.RowGroup
}

func RepeatedColumnIter[T any](ctx context.Context, source Source, columnName string, rows iter.Iterator[T]) iter.Iterator[*RepeatedRow[T]] {
	column, found := source.Schema().Lookup(strings.Split(columnName, ".")...)
	if !found {
		return iter.NewErrIterator[*RepeatedRow[T]](fmt.Errorf("column '%s' not found in parquet file", columnName))
	}

	opentracing.SpanFromContext(ctx).SetTag("columnName", columnName)
	return NewRepeatedPageIterator(ctx, rows, source.RowGroups(), column.ColumnIndex, 1e4)
}
