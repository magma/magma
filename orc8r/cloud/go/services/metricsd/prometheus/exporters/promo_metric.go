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
	"sort"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const (
	NetworkLabelNetwork = "networkID"
	NetworkLabelGateway = "gatewayID"
	NetworkLabelService = "service"
	NetworkLabelHost    = "host"

	MetricPostfixCount  = "__count"
	MetricPostfixBucket = "__bucket"
	MetricPostfixLE     = "__le"
	MetricPostfixSum    = "__sum"
	MetricInfixBucket   = "_bucket_"
)

// PrometheusSubmittable provides Register and Update functions to facilitate
// exporting metrics to Prometheus and updating them
type PrometheusMetric interface {
	Register(metric *dto.Metric, name string, exporter *PrometheusExporter, networkLabels prometheus.Labels) error
	Update(metric *dto.Metric, networkLabels prometheus.Labels) error
}

// PrometheusGauge wraps a prometheus Gauge
type PrometheusGauge struct {
	gaugeVec *prometheus.GaugeVec
}

// NewPrometheusGauge returns new PrometheusGauge
func NewPrometheusGauge() PrometheusMetric {
	return &PrometheusGauge{}
}

// Register registers gauge with given exporter
func (g *PrometheusGauge) Register(metric *dto.Metric,
	name string,
	exporter *PrometheusExporter,
	networkLabels prometheus.Labels,
) error {
	exporter.registeredMetrics[name] = g

	networkLabelNames := extractNetworkLabelKeys(networkLabels)
	g.gaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: name}, networkLabelNames)
	g.gaugeVec.With(networkLabels).Set(metric.GetGauge().GetValue())

	err := exporter.Registry.Register(g.gaugeVec)
	if err != nil {
		return fmt.Errorf("could not register Gauge: %v", err)
	}
	return nil
}

// Update updates the gauge
func (g *PrometheusGauge) Update(metric *dto.Metric, networkLabels prometheus.Labels) error {
	g.gaugeVec.With(networkLabels).Set(metric.GetGauge().GetValue())
	return nil
}

// PrometheusSummary wraps a prometheus Summary
type PrometheusSummary struct {
	sumGaugeVec   *prometheus.GaugeVec
	countGaugeVec *prometheus.GaugeVec
}

// NewPrometheusSummary returns a new PrometheusSummary
func NewPrometheusSummary() PrometheusMetric {
	return &PrometheusSummary{}
}

// Register registers summary with given exporter
func (s *PrometheusSummary) Register(metric *dto.Metric,
	name string,
	exporter *PrometheusExporter,
	networkLabels prometheus.Labels,
) error {
	exporter.registeredMetrics[name] = s

	sumName := makeSumName(name)
	networkLabelNames := extractNetworkLabelKeys(networkLabels)
	s.sumGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: sumName, Help: sumName}, networkLabelNames)
	s.sumGaugeVec.With(networkLabels).Set(metric.GetSummary().GetSampleSum())
	err := exporter.Registry.Register(s.sumGaugeVec)
	if err != nil {
		return fmt.Errorf("could not register Summary: %v", err)
	}

	countName := makeCountName(name)
	s.countGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: countName, Help: countName}, networkLabelNames)
	s.countGaugeVec.With(networkLabels).Set(float64(metric.GetSummary().GetSampleCount()))
	err = exporter.Registry.Register(s.countGaugeVec)
	if err != nil {
		return fmt.Errorf("could not register Summary: %v", err)
	}
	return nil
}

// Update updates the summary
func (s *PrometheusSummary) Update(metric *dto.Metric, networkLabels prometheus.Labels) error {
	s.sumGaugeVec.With(networkLabels).Set(metric.GetSummary().GetSampleSum())
	s.countGaugeVec.With(networkLabels).Set(float64(metric.GetSummary().GetSampleCount()))
	return nil
}

// PrometheusCounter wraps a prometheus counter. Includes
type PrometheusCounter struct {
	counterVec        *prometheus.CounterVec
	counterValues     map[string]float64
	name              string
	exporterReference *PrometheusExporter // Need to keep exporter in case of re-registering
}

// NewPrometheusCounter returns a new PrometheusCounter
func NewPrometheusCounter(exporter *PrometheusExporter) PrometheusMetric {
	return &PrometheusCounter{
		exporterReference: exporter,
		counterValues:     make(map[string]float64),
	}
}

// Register registers counter with a given exporter
func (c *PrometheusCounter) Register(metric *dto.Metric,
	name string,
	exporter *PrometheusExporter,
	networkLabels prometheus.Labels,
) error {
	exporter.registeredMetrics[name] = c

	networkLabelNames := extractNetworkLabelKeys(networkLabels)
	c.counterVec = prometheus.NewCounterVec(prometheus.CounterOpts{Name: name, Help: name}, networkLabelNames)
	c.counterVec.With(networkLabels).Add(metric.GetCounter().GetValue())
	c.counterValues[makeNetworkLabelString(networkLabelNames, networkLabels)] = metric.GetCounter().GetValue()
	c.exporterReference = exporter
	c.name = name

	err := exporter.Registry.Register(c.counterVec)
	if err != nil {
		return fmt.Errorf("could not register Counter: %v", err)
	}
	return nil
}

// Update updates the counter and handles a decreased counter by deleting
// and then re-adding the counter with the new value
func (c *PrometheusCounter) Update(metric *dto.Metric, networkLabels prometheus.Labels) error {
	newVal := metric.Counter.GetValue()
	networkLabelNames := extractNetworkLabelKeys(networkLabels)
	if oldVal, ok := c.counterValues[makeNetworkLabelString(networkLabelNames, networkLabels)]; ok {
		difference := newVal - oldVal
		if difference < 0 {
			return c.handleDecreasedCounter(newVal, networkLabels)
		}
		c.counterVec.With(networkLabels).Add(difference)
	} else {
		c.counterVec.With(networkLabels).Add(newVal)
	}
	c.counterValues[makeNetworkLabelString(networkLabelNames, networkLabels)] = newVal
	return nil
}

func (c *PrometheusCounter) handleDecreasedCounter(newVal float64, networkLabels prometheus.Labels) error {
	if ok := c.counterVec.Delete(networkLabels); !ok {
		return fmt.Errorf("could not update counter: failed to delete decreased counter %v", c.name)
	}
	c.counterVec.With(networkLabels).Add(newVal)
	networkLabelNames := extractNetworkLabelKeys(networkLabels)
	c.counterValues[makeNetworkLabelString(networkLabelNames, networkLabels)] = newVal
	return nil
}

// PrometheusHistogram wraps prometheus Histogram
type PrometheusHistogram struct {
	baseName      string
	sumGaugeVec   *prometheus.GaugeVec
	countGaugeVec *prometheus.GaugeVec
	bucketCounts  map[string]*prometheus.GaugeVec
	bucketLEs     map[string]*prometheus.GaugeVec
}

// NewPrometheusHistogram returns a new PrometheusHistogram
func NewPrometheusHistogram(baseName string) PrometheusMetric {
	return &PrometheusHistogram{
		baseName:     baseName,
		bucketLEs:    make(map[string]*prometheus.GaugeVec),
		bucketCounts: make(map[string]*prometheus.GaugeVec),
	}
}

// Register registers histogram with a given exporter
func (h *PrometheusHistogram) Register(metric *dto.Metric,
	name string,
	exporter *PrometheusExporter,
	networkLabels prometheus.Labels,
) error {
	exporter.registeredMetrics[name] = h

	networkLabelNames := extractNetworkLabelKeys(networkLabels)
	sumName := makeSumName(name)
	h.sumGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: sumName, Help: sumName}, networkLabelNames)
	h.sumGaugeVec.With(networkLabels).Set(metric.GetHistogram().GetSampleSum())
	err := exporter.Registry.Register(h.sumGaugeVec)
	if err != nil {
		return fmt.Errorf("could not register Histogram sum: %v", err)
	}

	countName := makeCountName(name)
	h.countGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: countName, Help: countName}, networkLabelNames)
	h.countGaugeVec.With(networkLabels).Set(float64(metric.GetHistogram().GetSampleCount()))
	err = exporter.Registry.Register(h.countGaugeVec)
	if err != nil {
		return fmt.Errorf("could not register Histogram count: %v", err)
	}

	for i, bucket := range metric.GetHistogram().GetBucket() {
		bucketLEName := makeBucketLEName(h.baseName, networkLabelNames, networkLabels, i)
		bucketLE := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: bucketLEName, Help: bucketLEName}, networkLabelNames)
		bucketLE.With(networkLabels).Set(bucket.GetUpperBound())
		h.bucketLEs[bucketLEName] = bucketLE
		err = exporter.Registry.Register(bucketLE)
		if err != nil {
			return fmt.Errorf("could not register Histogram bucketLE: %v", err)
		}

		bucketCountName := makeBucketCountName(h.baseName, networkLabelNames, networkLabels, i)
		bucketCount := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: bucketCountName, Help: bucketCountName}, networkLabelNames)
		bucketCount.With(networkLabels).Set(float64(bucket.GetCumulativeCount()))
		h.bucketCounts[bucketCountName] = bucketCount
		err = exporter.Registry.Register(bucketCount)
		if err != nil {
			return fmt.Errorf("could not register Histogram bucketCount: %v", err)
		}
	}
	return nil
}

// Update updates the histogram
func (h *PrometheusHistogram) Update(metric *dto.Metric, networkLabels prometheus.Labels) error {
	h.sumGaugeVec.With(networkLabels).Set(metric.GetHistogram().GetSampleSum())
	h.countGaugeVec.With(networkLabels).Set(float64(metric.GetHistogram().GetSampleCount()))
	networkLabelNames := extractNetworkLabelKeys(networkLabels)

	for i, bucket := range metric.GetHistogram().GetBucket() {
		bucketLEName := makeBucketLEName(h.baseName, networkLabelNames, networkLabels, i)
		bucketCountName := makeBucketCountName(h.baseName, networkLabelNames, networkLabels, i)
		if bucketLE, ok := h.bucketLEs[bucketLEName]; ok {
			bucketLE.With(networkLabels).Set(bucket.GetUpperBound())
		} else {
			bucketLE = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: bucketLEName, Help: bucketLEName}, networkLabelNames)
			bucketLE.With(networkLabels).Set(bucket.GetUpperBound())
			h.bucketLEs[bucketLEName] = bucketLE
		}
		if bucketCount, ok := h.bucketCounts[bucketCountName]; ok {
			bucketCount.With(networkLabels).Set(float64(bucket.GetCumulativeCount()))
		} else {
			bucketCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: bucketCountName, Help: bucketCountName}, networkLabelNames)
			bucketCount.With(networkLabels).Set(float64(bucket.GetCumulativeCount()))
			h.bucketCounts[bucketCountName] = bucketCount
		}
	}
	return nil
}

func makeNetworkLabelString(labelNames []string, networkLabels prometheus.Labels) string {
	str := ""
	for idx, labelName := range labelNames {
		if idx == len(labelNames)-1 {
			str += labelName + "." + networkLabels[labelName]
		} else {
			str += "." + labelName + "." + networkLabels[labelName]
		}
	}
	return str
}

func makeCountName(name string) string {
	return name + MetricPostfixCount
}

func makeSumName(name string) string {
	return name + MetricPostfixSum
}

func makeBucketLEName(baseName string, labelNames []string, networkLabels prometheus.Labels, bucketNum int) string {
	bucketLEName := baseName
	for _, labelName := range labelNames {
		bucketLEName += "_" + networkLabels[labelName]
	}
	bucketLEName += MetricInfixBucket + strconv.Itoa(bucketNum) + MetricPostfixLE
	return bucketLEName
}

func makeBucketCountName(baseName string, labelNames []string, networkLabels prometheus.Labels, bucketNum int) string {
	bucketCountName := baseName
	for _, labelName := range labelNames {
		bucketCountName += "_" + networkLabels[labelName]
	}
	bucketCountName += MetricInfixBucket + strconv.Itoa(bucketNum) + MetricPostfixCount
	return bucketCountName
}

func extractNetworkLabelKeys(networkLabels prometheus.Labels) []string {
	var keys []string
	for k := range networkLabels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
