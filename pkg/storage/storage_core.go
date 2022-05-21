package storage

import (
	"context"

	"github.com/pyroscope-io/pyroscope/pkg/storage/types"
)

func (s *Storage) Get(ctx context.Context, gi *types.GetInput) (*types.GetOutput, error) {
	s.metrics.getTotal.Inc()
	return s.core.Get(ctx, gi)
}

func (s *Storage) Delete(ctx context.Context, di *types.DeleteInput) error {
	return s.core.Delete(ctx, di)
}

func (s *Storage) DeleteApp(ctx context.Context, appname string) error {
	return s.core.DeleteApp(ctx, appname)
}

func (s *Storage) Put(ctx context.Context, pi *types.PutInput) error {
	s.metrics.putTotal.Inc()

	if s.hc.IsOutOfDiskSpace() {
		return errOutOfSpace
	}
	if pi.StartTime.Before(s.retentionPolicy().LowerTimeBoundary()) {
		return errRetention
	}

	return s.core.Put(ctx, pi)
}

//revive:disable-next-line:get-return callback is used
func (s *Storage) GetKeys(ctx context.Context, cb func(string) bool) {
	s.core.GetKeys(ctx, cb)
}

//revive:disable-next-line:get-return callback is used
func (s *Storage) GetValues(ctx context.Context, key string, cb func(v string) bool) {
	s.core.GetValues(ctx, key, cb)
}

func (s *Storage) GetKeysByQuery(ctx context.Context, in types.GetLabelKeysByQueryInput) (types.GetLabelKeysByQueryOutput, error) {
	return s.core.GetKeysByQuery(ctx, in)
}

// GetAppNames returns the list of all app's names
func (s *Storage) GetAppNames(ctx context.Context) []string {
	return s.core.GetAppNames(ctx)
}

func (s *Storage) MergeProfiles(ctx context.Context, mi types.MergeProfilesInput) (o types.MergeProfilesOutput, err error) {
	return s.core.MergeProfiles(ctx, mi)
}
