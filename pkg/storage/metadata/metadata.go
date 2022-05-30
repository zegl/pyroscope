package metadata

type Units string

const (
	SamplesUnits         Units = "samples"
	ObjectsUnits               = "objects"
	BytesUnits                 = "bytes"
	LockNanosecondsUnits       = "lock_nanoseconds"
	LockSamplesUnits           = "lock_samples"
)

type AggregationType string

const (
	AverageAggregationType AggregationType = "average"
	SumAggregationType     AggregationType = "sum"
)

func (a AggregationType) String() string {
	return a.String()
}

func (u Units) String() string {
	return u.String()
}
