package model

import (
	"github.com/google/pprof/profile"
	"github.com/samber/lo"

	"github.com/grafana/pyroscope/api/gen/proto/go/ingester/v1"
)

type ProfileSampleByPartition map[uint64]ProfileSampleMap

func (m ProfileSampleByPartition) Add(mapping uint64, key uint32, value int64) {
	if _, ok := m[mapping]; !ok {
		m[mapping] = make(ProfileSampleMap)
	}
	m[mapping].Add(key, value)
}

func (m ProfileSampleByPartition) ForEach(f func(mapping uint64, samples ProfileSampleMap) error) error {
	for mapping, samples := range m {
		if err := f(mapping, samples); err != nil {
			return err
		}
	}
	return nil
}

func (m ProfileSampleByPartition) StacktraceSamples() []*profile.Sample {
	var result []*profile.Sample
	for _, samples := range m {
		result = append(result, lo.Values(samples)...)
	}
	return result
}

type ProfileSampleMap map[uint32]*profile.Sample

func (m ProfileSampleMap) Add(key uint32, value int64) {
	if _, ok := m[key]; ok {
		m[key].Value[0] += value
		return
	}
	m[key] = &profile.Sample{
		Value: []int64{value},
	}
}

func (m ProfileSampleMap) Ids() []uint32 {
	return lo.Keys(m)
}

type StacktracesByPartition map[uint64]StacktraceSampleMap

func (m StacktracesByPartition) Add(partition uint64, key uint32, value int64) {
	if _, ok := m[partition]; !ok {
		m[partition] = make(StacktraceSampleMap)
	}
	m[partition].Add(key, value)
}

func (m StacktracesByPartition) ForEach(f func(mapping uint64, samples StacktraceSampleMap) error) error {
	for mapping, samples := range m {
		if err := f(mapping, samples); err != nil {
			return err
		}
	}
	return nil
}

func (m StacktracesByPartition) StacktraceSamples() []*ingesterv1.StacktraceSample {
	var result []*ingesterv1.StacktraceSample
	for _, samples := range m {
		result = append(result, lo.Values(samples)...)
	}
	return result
}

type StacktraceSampleMap map[uint32]*ingesterv1.StacktraceSample

func (m StacktraceSampleMap) Add(key uint32, value int64) {
	if _, ok := m[key]; ok {
		m[key].Value += value
		return
	}
	m[key] = &ingesterv1.StacktraceSample{
		Value: value,
	}
}

func (m StacktraceSampleMap) Ids() []uint32 {
	return lo.Keys(m)
}
