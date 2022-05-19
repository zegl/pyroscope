package core

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type ObjectCache interface {
	Put(key string, val interface{})
	Evict(percent float64)
	WriteBack()
	Delete(key string) error
	Discard(key string)
	DiscardPrefix(prefix string) error
	GetOrCreate(key string) (interface{}, error)
	Lookup(key string) (interface{}, bool)
}
type LabelsStore interface {
}
type ExemplarsStore interface {
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
