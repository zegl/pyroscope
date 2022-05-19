package storage

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/pyroscope-io/pyroscope/pkg/storage/exemplars"
)

const (
	exemplarBatches       = 5
	exemplarsPerBatch     = 10 << 10 // 10K
	exemplarBatchDuration = time.Second * 5
)

func (s *Storage) initExemplarsStorage(db BadgerDBWithCache, reg prometheus.Registerer) {
	e := exemplars.New(&exemplars.ExemplarsConfig{
		Db:                    *db.DBInstance(),
		Dicts:                 s.dicts.CacheInstance(),
		Logger:                s.logger,
		MaxNodesSerialization: s.config.maxNodesSerialization,
		Reg:                   reg,
	})
	s.exemplars = e
	s.tasksWG.Add(1)

	go func() {
		retentionTicker := time.NewTicker(s.retentionTaskInterval)
		batchFlushTicker := time.NewTicker(exemplarBatchDuration)
		defer func() {
			batchFlushTicker.Stop()
			retentionTicker.Stop()
			s.tasksWG.Done()
		}()
		for {
			select {
			default:
			case batch, ok := <-e.Batches():
				if ok {
					e.Flush(batch)
				}
			}

			select {
			case <-s.stop:
				s.logger.Debug("flushing batches queue")
				e.FlushBatchQueue()
				return

			case <-batchFlushTicker.C:
				s.logger.Debug("flushing current batch")
				e.FlushCurrentBatch()

			case batch, ok := <-e.Batches():
				if ok {
					e.Flush(batch)
				}

			case <-retentionTicker.C:
				s.exemplarsRetentionTask()
			}
		}
	}()
}
