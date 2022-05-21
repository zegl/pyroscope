package core

import (
	"context"

	"github.com/pyroscope-io/pyroscope/pkg/storage/tree"
	"github.com/pyroscope-io/pyroscope/pkg/storage/types"
)

func (s *Core) MergeProfiles(ctx context.Context, mi types.MergeProfilesInput) (o types.MergeProfilesOutput, err error) {
	o.Tree = tree.New()
	return o, s.exemplars.Fetch(ctx, mi.AppName, mi.Profiles, func(t *tree.Tree) error {
		o.Tree.Merge(t)
		return nil
	})
}
