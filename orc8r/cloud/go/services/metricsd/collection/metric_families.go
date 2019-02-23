/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package collection

import (
	"github.com/prometheus/client_model/go"
)

type MetricLabel struct {
	Name, Value string
}

// MakeSingleGaugeFamily returns a MetricFamily with a single gauge value
// as specified by the function arguments.
// label is nil-able - a nil input will return a gauge metric without a label.
func MakeSingleGaugeFamily(
	name string, help string,
	label *MetricLabel,
	value float64,
) *io_prometheus_client.MetricFamily {
	mtype := io_prometheus_client.MetricType_GAUGE
	return &io_prometheus_client.MetricFamily{
		Name:   &name,
		Help:   &help,
		Type:   &mtype,
		Metric: []*io_prometheus_client.Metric{MakeSingleGaugeMetric(label, value)},
	}
}

// MakeMultiGaugeFamily returns a single MetricsFamily with multiple labeled
// gauge metrics.
func MakeMultiGaugeFamily(
	name string, help string,
	gaugesByLabel map[MetricLabel]float64,
) *io_prometheus_client.MetricFamily {
	metrics := make([]*io_prometheus_client.Metric, 0, len(gaugesByLabel))

	for label, val := range gaugesByLabel {
		clonedLabel := MetricLabel{Name: label.Name, Value: label.Value}

		metrics = append(
			metrics,
			&io_prometheus_client.Metric{
				Label: []*io_prometheus_client.LabelPair{
					{Name: &clonedLabel.Name, Value: &clonedLabel.Value},
				},
				Gauge: &io_prometheus_client.Gauge{Value: &val},
			},
		)
	}

	mtype := io_prometheus_client.MetricType_GAUGE
	return &io_prometheus_client.MetricFamily{
		Name:   &name,
		Help:   &help,
		Type:   &mtype,
		Metric: metrics,
	}
}

// MakeSingleGaugeMetric returns a Metric with a single gauge value as
// specified by the function argument.
// label is nil-able - a nil input will return a gauge metric without a label.
func MakeSingleGaugeMetric(
	label *MetricLabel,
	value float64,
) *io_prometheus_client.Metric {
	if label == nil {
		return &io_prometheus_client.Metric{
			Gauge: &io_prometheus_client.Gauge{Value: &value},
		}
	} else {
		return &io_prometheus_client.Metric{
			Label: []*io_prometheus_client.LabelPair{
				{Name: &label.Name, Value: &label.Value},
			},
			Gauge: &io_prometheus_client.Gauge{Value: &value},
		}
	}
}
