/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"fmt"

	dto "github.com/prometheus/client_model/go"
)

const (
	bucketPostfix = "_bucket"
	countPostfix  = "_count"
	sumPostfix    = "_sum"

	histogramBucketLabelName = "le"
	summaryQuantileLabelName = "quantile"
)

var (
	gaugeType = dto.MetricType_GAUGE
)

func convertFamilyToGauges(baseFamily *dto.MetricFamily) []*dto.MetricFamily {
	gaugeFamilies := make([]*dto.MetricFamily, 0)
	switch *baseFamily.Type {
	case dto.MetricType_GAUGE:
		gaugeFamilies = append(gaugeFamilies, baseFamily)
	case dto.MetricType_COUNTER:
		gaugeFamilies = append(gaugeFamilies, counterToGauge(baseFamily))
	case dto.MetricType_HISTOGRAM:
		gaugeFamilies = append(gaugeFamilies, histogramToGauges(baseFamily)...)
	case dto.MetricType_SUMMARY:
		gaugeFamilies = append(gaugeFamilies, summaryToGauges(baseFamily)...)
	case dto.MetricType_UNTYPED:
		gaugeFamilies = append(gaugeFamilies, untypedToGauge(baseFamily))
	}
	return gaugeFamilies
}

// counterToGauge takes a counter and converts it to a gauge with the
// same value
func counterToGauge(family *dto.MetricFamily) *dto.MetricFamily {
	counterFamily := dto.MetricFamily{
		Name: makeStringPointer(family.GetName()),
		Type: &gaugeType,
	}
	for _, metric := range family.Metric {
		if metric.Counter == nil {
			continue
		}
		counterValue := float64(*metric.Counter.Value)
		counterMetric := dto.Metric{
			Label: metric.Label,
			Gauge: &dto.Gauge{
				Value: &counterValue,
			},
		}
		counterFamily.Metric = append(counterFamily.Metric, &counterMetric)
	}
	return &counterFamily
}

// histogramToGauges converts a histogram into 3 families of gauges, one for the
// buckets, sum, and count each.
func histogramToGauges(family *dto.MetricFamily) []*dto.MetricFamily {
	baseName := family.GetName()
	bucketFamily := dto.MetricFamily{
		Name: makeStringPointer(baseName + bucketPostfix),
		Type: &gaugeType,
	}
	sumFamily := dto.MetricFamily{
		Name: makeStringPointer(baseName + sumPostfix),
		Type: &gaugeType,
	}
	countFamily := dto.MetricFamily{
		Name: makeStringPointer(baseName + countPostfix),
		Type: &gaugeType,
	}

	for _, metric := range family.Metric {
		if metric.Histogram == nil {
			continue
		}
		sumValue := float64(*metric.Histogram.SampleSum)
		sumMetric := dto.Metric{
			Label: metric.Label,
			Gauge: &dto.Gauge{
				Value: &sumValue,
			},
		}
		sumFamily.Metric = append(sumFamily.Metric, &sumMetric)

		countValue := float64(*metric.Histogram.SampleCount)
		countMetric := dto.Metric{
			Label: metric.Label,
			Gauge: &dto.Gauge{
				Value: &countValue,
			},
		}
		countFamily.Metric = append(countFamily.Metric, &countMetric)

		for _, bucket := range metric.Histogram.Bucket {
			bucketValue := float64(*bucket.CumulativeCount)
			bucketMetric := dto.Metric{
				Label: append(metric.Label, &dto.LabelPair{
					Name:  makeStringPointer(histogramBucketLabelName),
					Value: makeStringPointer(fmt.Sprintf("%g", bucket.GetUpperBound())),
				}),
				Gauge: &dto.Gauge{
					Value: &bucketValue,
				},
			}
			bucketFamily.Metric = append(bucketFamily.Metric, &bucketMetric)
		}
	}
	return []*dto.MetricFamily{&bucketFamily, &sumFamily, &countFamily}
}

// summaryToGauges converts a summary to 3 gauge families, one for the quantiles,
// sum, and count each
func summaryToGauges(family *dto.MetricFamily) []*dto.MetricFamily {
	baseName := family.GetName()
	quantFamily := dto.MetricFamily{
		Name: makeStringPointer(baseName),
		Type: &gaugeType,
	}
	sumFamily := dto.MetricFamily{
		Name: makeStringPointer(baseName + sumPostfix),
		Type: &gaugeType,
	}
	countFamily := dto.MetricFamily{
		Name: makeStringPointer(baseName + countPostfix),
		Type: &gaugeType,
	}

	for _, metric := range family.Metric {
		if metric.Summary == nil {
			continue
		}
		sumValue := float64(*metric.Summary.SampleSum)
		sumMetric := dto.Metric{
			Label: metric.Label,
			Gauge: &dto.Gauge{
				Value: &sumValue,
			},
		}
		sumFamily.Metric = append(sumFamily.Metric, &sumMetric)

		countValue := float64(*metric.Summary.SampleCount)
		countMetric := dto.Metric{
			Label: metric.Label,
			Gauge: &dto.Gauge{
				Value: &countValue,
			},
		}
		countFamily.Metric = append(countFamily.Metric, &countMetric)

		for _, quant := range metric.Summary.Quantile {
			quantValue := *quant.Value
			quantMetric := dto.Metric{
				Label: append(metric.Label, &dto.LabelPair{
					Name:  makeStringPointer(summaryQuantileLabelName),
					Value: makeStringPointer(fmt.Sprintf("%g", *quant.Quantile)),
				}),
				Gauge: &dto.Gauge{
					Value: &quantValue,
				},
			}
			quantFamily.Metric = append(quantFamily.Metric, &quantMetric)
		}
	}
	return []*dto.MetricFamily{&quantFamily, &sumFamily, &countFamily}
}

// untypedToGauge takes an untyped metric and converts it to a gauge with the
// same value
func untypedToGauge(family *dto.MetricFamily) *dto.MetricFamily {
	untypedFamily := dto.MetricFamily{
		Name: makeStringPointer(family.GetName()),
		Type: &gaugeType,
	}
	for _, metric := range family.Metric {
		if metric.Untyped == nil {
			continue
		}
		untypedValue := float64(*metric.Untyped.Value)
		untypedMetric := dto.Metric{
			Label: metric.Label,
			Gauge: &dto.Gauge{
				Value: &untypedValue,
			},
		}
		untypedFamily.Metric = append(untypedFamily.Metric, &untypedMetric)
	}
	return &untypedFamily
}
