package storage

//revive:disable:max-public-structs TODO: we will refactor this later

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/pyroscope-io/pyroscope/pkg/storage/cache"
	"github.com/pyroscope-io/pyroscope/pkg/storage/types"
	"github.com/pyroscope-io/pyroscope/pkg/util/bytesize"
)

type PutInput types.PutInput
type Putter types.Putter
type GetInput types.GetInput
type GetOutput types.GetOutput
type Getter types.Getter
type Enqueuer types.Enqueuer
type MergeProfilesInput types.MergeProfilesInput
type MergeProfilesOutput types.MergeProfilesOutput
type Merger types.Merger
type GetLabelKeysByQueryInput types.GetLabelKeysByQueryInput
type GetLabelKeysByQueryOutput types.GetLabelKeysByQueryOutput
type LabelsGetter types.LabelsGetter
type GetLabelValuesByQueryInput types.GetLabelValuesByQueryInput
type GetLabelValuesByQueryOutput types.GetLabelValuesByQueryOutput
type LabelValuesGetter types.LabelValuesGetter
type AppNameGetter types.AppNameGetter

// Other functions from storage.Storage:
// type Backend interface {
// 	Put(ctx context.Context, pi *PutInput) error
// 	Get(ctx context.Context, gi *GetInput) (*GetOutput, error)

// 	Enqueue(ctx context.Context, input *PutInput)

// 	GetAppNames(ctx context.Context, ) []string
// 	GetKeys(ctx context.Context, cb func(string) bool)
// 	GetKeysByQuery(ctx context.Context, query string, cb func(_k string) bool) error
// 	GetValues(ctx context.Context, key string, cb func(v string) bool)
// 	GetValuesByQuery(ctx context.Context, label string, query string, cb func(v string) bool) error
// 	DebugExport(ctx context.Context, w http.ResponseWriter, r *http.Request)

// 	Delete(ctx context.Context, di *DeleteInput) error
// 	DeleteApp(ctx context.Context, appname string) error
// }

type BadgerDB interface {
	Update(func(txn *badger.Txn) error) error
	View(func(txn *badger.Txn) error) error
	NewWriteBatch() *badger.WriteBatch
	MaxBatchCount() int64
}

type CacheLayer interface {
	Put(key string, val interface{})
	Evict(percent float64)
	WriteBack()
	Delete(key string) error
	Discard(key string)
	DiscardPrefix(prefix string) error
	GetOrCreate(key string) (interface{}, error)
	Lookup(key string) (interface{}, bool)
}

type BadgerDBWithCache interface {
	BadgerDB
	CacheLayer

	Close()
	Size() bytesize.ByteSize
	CacheSize() uint64

	DBInstance() *badger.DB
	CacheInstance() *cache.Cache
	Name() string
}
