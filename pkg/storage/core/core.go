package core

import (
	"sync"
	"time"

	"github.com/pyroscope-io/pyroscope/pkg/storage/tree"
	"github.com/sirupsen/logrus"
)

type ObjectCache interface {
	Put(key string, val interface{})
	Delete(key string) error
	Discard(key string)
	DiscardPrefix(prefix string) error
	GetOrCreate(key string) (interface{}, error)
	Lookup(key string) (interface{}, bool)
}

type LabelsStore interface {
	Put(k, v string) error
	GetKeys(cb func(k string) bool) error
	Delete(key, value string) error
	GetValues(key string, cb func(v string) bool) error
}

type ExemplarsStore interface {
	Insert(appName, profileID string, v *tree.Tree, timestamp time.Time) error
}

type Core struct {
	// TODO(petethepig): get rid of this lock, it's too broad
	putMutex sync.Mutex

	segments   ObjectCache
	dimensions ObjectCache
	dicts      ObjectCache
	trees      ObjectCache
	main       ObjectCache
	labels     LabelsStore
	exemplars  ExemplarsStore

	logger logrus.FieldLogger
}

func New() *Core {
	return &Core{}
}
