package storage

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/pyroscope-io/pyroscope/pkg/storage/types"
)

type putInputWithCtx struct {
	pi  *types.PutInput
	ctx context.Context
}

func (s *Storage) Enqueue(ctx context.Context, input *types.PutInput) {
	select {
	case s.queue <- &putInputWithCtx{input, ctx}:
		return
	case <-s.stop:
	default:
	}
	s.metrics.discardedTotal.Inc()
}

func (s *Storage) startQueueWorkers() {
	s.queueWorkersWG.Add(s.queueWorkers)
	for i := 0; i < s.queueWorkers; i++ {
		go s.runQueueWorker()
	}
}

func (s *Storage) runQueueWorker() {
	defer s.queueWorkersWG.Done()
	for {
		select {
		case input, ok := <-s.queue:
			if ok {
				if err := s.safePut(input.ctx, input.pi); err != nil {
					s.logger.WithField("key", input.pi.Key.Normalized()).WithError(err).Error("error happened while ingesting data")
				}
			}
		case <-s.stop:
			return
		}
	}
}

func (s *Storage) safePut(ctx context.Context, input *types.PutInput) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v; %v", r, string(debug.Stack()))
		}
	}()
	// TODO(petethepig): not sure if retaining context is a good practice
	return s.Put(ctx, input)
}
