package storage

import (
	"context"

	"github.com/pyroscope-io/pyroscope/pkg/storage/tree"
)

func (s *Storage) MergeProfiles(ctx context.Context, mi MergeProfilesInput) (o MergeProfilesOutput, err error) {
	o.Tree = tree.New()
	return o, s.exemplars.fetch(ctx, mi.AppName, mi.Profiles, func(t *tree.Tree) error {
		o.Tree.Merge(t)
		return nil
	})
}
