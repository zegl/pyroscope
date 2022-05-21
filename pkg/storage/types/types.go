package types

//revive:disable:max-public-structs TODO: we will refactor this later

import (
	"context"

	"time"

	"github.com/pyroscope-io/pyroscope/pkg/flameql"
	"github.com/pyroscope-io/pyroscope/pkg/storage/metadata"
	"github.com/pyroscope-io/pyroscope/pkg/storage/segment"
	"github.com/pyroscope-io/pyroscope/pkg/storage/tree"
)

type PutInput struct {
	StartTime       time.Time
	EndTime         time.Time
	Key             *segment.Key
	Val             *tree.Tree
	SpyName         string
	SampleRate      uint32
	Units           metadata.Units
	AggregationType metadata.AggregationType
}

type Putter interface {
	Put(ctx context.Context, pi *PutInput) error
}

type GetInput struct {
	StartTime time.Time
	EndTime   time.Time
	Key       *segment.Key
	Query     *flameql.Query
}

type GetOutput struct {
	Tree            *tree.Tree
	Timeline        *segment.Timeline
	SpyName         string
	SampleRate      uint32
	Count           uint64
	Units           metadata.Units
	AggregationType metadata.AggregationType
}

type Getter interface {
	Get(ctx context.Context, gi *GetInput) (*GetOutput, error)
}

type Enqueuer interface {
	Enqueue(ctx context.Context, input *PutInput)
}

type MergeProfilesInput struct {
	AppName  string
	Profiles []string
}

type MergeProfilesOutput struct {
	Tree *tree.Tree
}

type Merger interface {
	MergeProfiles(ctx context.Context, mi MergeProfilesInput) (o MergeProfilesOutput, err error)
}

type GetLabelKeysByQueryInput struct {
	StartTime time.Time
	EndTime   time.Time
	Query     string
}

type GetLabelKeysByQueryOutput struct {
	Keys []string
}

type LabelsGetter interface {
	GetKeys(ctx context.Context, cb func(string) bool)
	GetKeysByQuery(ctx context.Context, in GetLabelKeysByQueryInput) (GetLabelKeysByQueryOutput, error)
}

type GetLabelValuesByQueryInput struct {
	StartTime time.Time
	EndTime   time.Time
	Label     string
	Query     string
}

type GetLabelValuesByQueryOutput struct {
	Values []string
}

type LabelValuesGetter interface {
	GetValues(ctx context.Context, key string, cb func(v string) bool)
	GetValuesByQuery(ctx context.Context, in GetLabelValuesByQueryInput) (GetLabelValuesByQueryOutput, error)
}

type AppNameGetter interface {
	GetAppNames(ctx context.Context) []string
}

type DeleteInput struct {
	// Key must match exactly one segment.
	Key *segment.Key
}
