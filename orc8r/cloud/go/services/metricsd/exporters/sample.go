/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters

import (
	"fmt"
	"sort"
	"strconv"

	"magma/orc8r/lib/go/metrics"

	dto "github.com/prometheus/client_model/go"
)

type Sample struct {
	name        string
	value       string
	timestampMs int64
	labels      []*dto.LabelPair
	// entity identifies the entity that created this sample
	// for samples to be exported to ODS, entity has to be in the form of <ID1>.<ID2>
	// e.g. networkId.logicalId, or cloud.host_name
	entity string
}

func NewSample(
	name string,
	value string,
	timestampMs int64,
	labels []*dto.LabelPair,
	entity string,
) Sample {
	return Sample{
		name:        name,
		value:       value,
		timestampMs: timestampMs,
		labels:      labels,
		entity:      entity,
	}
}

func (s *Sample) Name() string {
	return s.name
}

func (s *Sample) Labels() []*dto.LabelPair {
	return s.labels
}

func (s *Sample) Value() string {
	return s.value
}

func (s *Sample) Entity() string {
	return s.entity
}

func (s *Sample) TimestampMs() int64 {
	return s.timestampMs
}

// GetSamplesForMetrics takes a Metric protobuf and extracts Samples from them
// since there may be multiple Samples per a single Metric instance
func GetSamplesForMetrics(metricAndContext MetricAndContext, metric *dto.Metric) []Sample {
	context := metricAndContext.Context
	family := metricAndContext.Family
	var entity string
	labels := metric.Label

	switch additionalCtx := metricAndContext.Context.AdditionalContext.(type) {
	case *CloudMetricContext:
		entity = fmt.Sprintf("cloud.%s", additionalCtx.CloudHost)
	case *GatewayMetricContext:
		entity = fmt.Sprintf("%s.%s", additionalCtx.NetworkID, additionalCtx.GatewayID)
	case *PushedMetricContext:
		gatewayID := popLabel(&labels, metrics.GatewayLabelName)
		if gatewayID == "" {
			entity = additionalCtx.NetworkID
		} else {
			entity = fmt.Sprintf("%s.%s", additionalCtx.NetworkID, gatewayID)
		}
	}

	name := context.MetricName
	timestampMs := metric.GetTimestampMs()
	sort.Sort(ByName(labels))

	switch family.GetType() {
	case dto.MetricType_COUNTER:
		return getCounterSamples(name, labels, timestampMs, metric.GetCounter(), entity)
	case dto.MetricType_GAUGE:
		return getGaugeSamples(name, labels, timestampMs, metric.GetGauge(), entity)
	case dto.MetricType_SUMMARY:
		return getSummarySamples(name, labels, timestampMs, metric.GetSummary(), entity)
	case dto.MetricType_HISTOGRAM:
		return getHistogramSamples(name, labels, timestampMs, metric.GetHistogram(), entity)
	}
	// I don't know what this is, return empty
	return []Sample{}
}

// getCounterSamples will extract a single counter sample from a Counter
func getCounterSamples(name string,
	labels []*dto.LabelPair,
	timestampMs int64,
	c *dto.Counter,
	entity string,
) []Sample {
	samples := make([]Sample, 1)
	samples[0] = Sample{
		name:        name,
		labels:      labels,
		timestampMs: timestampMs,
		value:       strconv.FormatFloat(c.GetValue(), 'f', -1, 64),
		entity:      entity,
	}
	return samples
}

// GetGaugeSamples will extract a single gauge sample from a Gauge
func getGaugeSamples(name string,
	labels []*dto.LabelPair,
	timestampMs int64,
	g *dto.Gauge,
	entity string,
) []Sample {
	samples := make([]Sample, 1)
	samples[0] = Sample{
		name:        name,
		labels:      labels,
		timestampMs: timestampMs,
		value:       strconv.FormatFloat(g.GetValue(), 'f', -1, 64),
		entity:      entity,
	}
	return samples
}

// GetSummarySamples will extract a two samples from a Summary
// one for count and another for sum
func getSummarySamples(name string,
	labels []*dto.LabelPair,
	timestampMs int64,
	s *dto.Summary,
	entity string,
) []Sample {
	samples := make([]Sample, 2)
	samples[0] = Sample{
		name:        name + "_count",
		labels:      labels,
		timestampMs: timestampMs,
		value:       strconv.FormatUint(s.GetSampleCount(), 10),
		entity:      entity,
	}
	samples[1] = Sample{
		name:        name + "_sum",
		labels:      labels,
		timestampMs: timestampMs,
		value:       strconv.FormatFloat(s.GetSampleSum(), 'f', -1, 64),
		entity:      entity}
	return samples
}

// GetHistogramSamples will extract 2 + 2*dim(Buckets) Samples
// for each Histogram instance
func getHistogramSamples(name string,
	labels []*dto.LabelPair,
	timestampMs int64,
	h *dto.Histogram,
	entity string,
) []Sample {
	samples := make([]Sample, len(h.GetBucket())*2+2)
	samples[0] = Sample{
		name:        name + "_count",
		labels:      labels,
		timestampMs: timestampMs,
		value:       strconv.FormatUint(h.GetSampleCount(), 10),
		entity:      entity,
	}
	samples[1] = Sample{
		name:        name + "_sum",
		labels:      labels,
		timestampMs: timestampMs,
		value:       strconv.FormatFloat(h.GetSampleSum(), 'E', -1, 64),
		entity:      entity,
	}
	for i, b := range h.GetBucket() {
		samples[i+2] = Sample{
			name:        fmt.Sprintf("%s_bucket_%d_le", name, i),
			labels:      labels,
			timestampMs: timestampMs,
			value:       strconv.FormatFloat(b.GetUpperBound(), 'E', -1, 64),
			entity:      entity,
		}
		samples[i+3] = Sample{
			name:        fmt.Sprintf("%s_bucket_%d_count", name, i),
			labels:      labels,
			timestampMs: timestampMs,
			value:       strconv.FormatUint(b.GetCumulativeCount(), 10),
			entity:      entity,
		}
	}
	return samples
}

func popLabel(labels *[]*dto.LabelPair, labelToRemove string) string {
	var ret string
	for i, label := range *labels {
		if label.GetName() == labelToRemove {
			ret = label.GetValue()
			*labels = removeLabel(*labels, i)
			break
		}
	}
	return ret
}

func removeLabel(labels []*dto.LabelPair, i int) []*dto.LabelPair {
	labels[len(labels)-1], labels[i] = labels[i], labels[len(labels)-1]
	return labels[:len(labels)-1]
}

// ByName is an interface for sorting LabelPairs by name
type ByName []*dto.LabelPair

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].GetName() < a[j].GetName() }
