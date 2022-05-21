package storage

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/pyroscope-io/pyroscope/pkg/storage/segment"
)

const defaultBatchSize = 1 << 10 // 1K items

func (s *Storage) enforceRetentionPolicy(ctx context.Context, rp *segment.RetentionPolicy) {
	observer := prometheus.ObserverFunc(s.metrics.retentionTaskDuration.Observe)
	timer := prometheus.NewTimer(observer)
	defer timer.ObserveDuration()

	s.retention.EnforceRetentionPolicy(ctx, rp)
}
