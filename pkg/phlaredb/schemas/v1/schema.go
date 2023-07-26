package v1

import profilev1 "github.com/grafana/pyroscope/api/gen/proto/go/google/v1"

type Models interface {
	*Profile | *InMemoryProfile |
		*profilev1.Location | *InMemoryLocation |
		*profilev1.Function | *InMemoryFunction |
		*profilev1.Mapping | *InMemoryMapping |
		*Stacktrace |
		string
}
