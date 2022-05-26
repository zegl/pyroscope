package remote

import (
	"context"

	"github.com/pyroscope-io/pyroscope/pkg/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Putter interface {
	Put(ctx context.Context, pi *storage.PutInput) error
}

type RemoteWriter struct {
	log *logrus.Logger
}

func NewRemoteWriter(logger *logrus.Logger) *RemoteWriter {
	return &RemoteWriter{
		log: logger,
	}
}

func (r *RemoteWriter) Put(_ context.Context, pi *storage.PutInput) error {
	logrus.Debugf("Putting in remote storage")
	return nil
}

// TODO(eh-am): rename to something else more clear
type StorageOrchestrator struct {
	log     *logrus.Logger
	putters []Putter
}

// TODO(eh-am): move this somewhere else, probably on its own package
// TODO(eh-am): should we use an array of putters?
// maybe this orchestrator should create the putters?
func NewStorageOrchestrator(log *logrus.Logger, putters ...Putter) *StorageOrchestrator {
	// TODO(eh-am): wrap the logger with a svc?
	return &StorageOrchestrator{
		log:     log,
		putters: putters,
	}
}

// Put puts data into the underlying Putters
// By default it uses a parallel strategy
func (so *StorageOrchestrator) Put(ctx context.Context, pi *storage.PutInput) error {
	logrus.Debugf("Putting in storage orchestrator")

	// TODO(eh-am): maybe we should have different strategies here?
	// like writing in parallel, or writing sequentially
	g, ctx := errgroup.WithContext(ctx)
	for _, p := range so.putters {
		// https://golang.org/doc/faq#closures_and_goroutines
		p := p

		g.Go(func() error {
			return p.Put(ctx, pi)
		})
	}

	return g.Wait()
}
