package model

import (
	"sort"
	"testing"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"

	"github.com/grafana/pyroscope/pkg/phlaredb"
)

func TestLabelsForProfiles(t *testing.T) {
	for _, tt := range []struct {
		name     string
		in       Labels
		expected Labels
	}{
		{
			"default",
			Labels{{Name: model.MetricNameLabel, Value: "cpu"}},
			Labels{
				{Name: model.MetricNameLabel, Value: "cpu"},
				{Name: LabelNameUnit, Value: "unit"},
				{Name: LabelNameProfileType, Value: "cpu:type:unit:type:unit"},
				{Name: LabelNameType, Value: "type"},
				{Name: LabelNamePeriodType, Value: "type"},
				{Name: LabelNamePeriodUnit, Value: "unit"},
			},
		},
		{
			"with service_name",
			Labels{
				{Name: model.MetricNameLabel, Value: "cpu"},
				{Name: LabelNameServiceName, Value: "service_name"},
			},
			Labels{
				{Name: model.MetricNameLabel, Value: "cpu"},
				{Name: LabelNameUnit, Value: "unit"},
				{Name: LabelNameProfileType, Value: "cpu:type:unit:type:unit"},
				{Name: LabelNameType, Value: "type"},
				{Name: LabelNamePeriodType, Value: "type"},
				{Name: LabelNamePeriodUnit, Value: "unit"},
				{Name: LabelNameServiceName, Value: "service_name"},
				{Name: LabelNameServiceName, Value: "service_name"},
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.expected)
			result, fps := LabelsForProfile(phlaredb.newProfileFoo(), tt.in...)
			require.Equal(t, tt.expected, result[0])
			require.Equal(t, model.Fingerprint(tt.expected.Hash()), fps[0])
		})
	}
}
