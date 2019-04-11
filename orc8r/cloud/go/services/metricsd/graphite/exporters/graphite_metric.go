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
	"strconv"
	"strings"
	"time"

	"github.com/marpaia/graphite-golang"
	dto "github.com/prometheus/client_model/go"
)

type GraphiteMetric interface {
	Update(metric *dto.Metric)
	Register(metric *dto.Metric, name string, exporter *GraphiteExporter)
	Export(exporter *GraphiteExporter) error
}

type graphiteBase struct {
	name       string
	updateTime int64
}

func (g *graphiteBase) setTimeNow() {
	g.updateTime = time.Now().Unix()
}

type GraphiteCounter struct {
	graphiteBase
	value float64
}

func NewGraphiteCounter() GraphiteMetric {
	return &GraphiteCounter{}
}

func (c *GraphiteCounter) Update(metric *dto.Metric) {
	c.value = metric.GetCounter().GetValue()
	c.setTimeNow()
}

func (c *GraphiteCounter) Register(metric *dto.Metric, name string, exporter *GraphiteExporter) {
	c.name = name
	exporter.registeredMetrics[c.name] = c
	c.value = metric.GetCounter().GetValue()
	c.setTimeNow()
}

func (c *GraphiteCounter) Export(exporter *GraphiteExporter) error {
	err := exporter.graphite.SendMetric(graphite.Metric{
		Name:      c.name,
		Value:     floatToString(c.value),
		Timestamp: c.updateTime,
	})
	if err != nil {
		return fmt.Errorf("could not send metric %v to graphite: %v", c.name, err)
	}
	return nil
}

type GraphiteGauge struct {
	graphiteBase
	value float64
}

func NewGraphiteGauge() GraphiteMetric {
	return &GraphiteGauge{}
}

func (g *GraphiteGauge) Update(metric *dto.Metric) {
	g.value = metric.GetGauge().GetValue()
	g.setTimeNow()
}

func (g *GraphiteGauge) Register(metric *dto.Metric, name string, exporter *GraphiteExporter) {
	g.name = name
	g.value = metric.GetGauge().GetValue()
	g.setTimeNow()
	exporter.registeredMetrics[g.name] = g
}

func (g *GraphiteGauge) Export(exporter *GraphiteExporter) error {
	err := exporter.graphite.SendMetric(graphite.Metric{
		Name:      g.name,
		Value:     floatToString(g.value),
		Timestamp: g.updateTime,
	})
	if err != nil {
		return fmt.Errorf("could not send metric %v to graphite: %v", g.name, err)
	}
	return nil
}

type GraphiteSummary struct {
	graphiteBase
	sumValue   float64
	countValue float64
}

func NewGraphiteSummary() GraphiteMetric {
	return &GraphiteSummary{}
}

func (s *GraphiteSummary) Update(metric *dto.Metric) {
	s.sumValue = metric.GetSummary().GetSampleSum()
	s.countValue = float64(metric.GetSummary().GetSampleCount())
	s.setTimeNow()
}

func (s *GraphiteSummary) Register(metric *dto.Metric, name string, exporter *GraphiteExporter) {
	s.name = name
	exporter.registeredMetrics[s.name] = s

	s.sumValue = metric.GetSummary().GetSampleSum()
	s.countValue = float64(metric.GetSummary().GetSampleCount())
	s.setTimeNow()
}

func (s *GraphiteSummary) Export(exporter *GraphiteExporter) error {
	sumMetric := graphite.Metric{
		Name:      makeGraphiteSumName(s.name),
		Value:     floatToString(s.sumValue),
		Timestamp: s.updateTime,
	}

	countMetric := graphite.Metric{
		Name:      makeGraphiteCountName(s.name),
		Value:     floatToString(s.countValue),
		Timestamp: s.updateTime,
	}
	err := exporter.graphite.SendMetrics([]graphite.Metric{sumMetric, countMetric})
	if err != nil {
		return fmt.Errorf("could not send metric %v to graphite: %v", s.name, err)
	}
	return nil
}

type GraphiteHistogram struct {
	graphiteBase
	sumValue     float64
	countValue   float64
	bucketCounts map[string]float64
	bucketLEs    map[string]float64
}

func NewGraphiteHistogram() GraphiteMetric {
	return &GraphiteHistogram{
		bucketCounts: make(map[string]float64),
		bucketLEs:    make(map[string]float64),
	}
}

func (h *GraphiteHistogram) Register(metric *dto.Metric, name string, exporter *GraphiteExporter) {
	h.name = name
	exporter.registeredMetrics[h.name] = h

	h.sumValue = metric.GetHistogram().GetSampleSum()
	h.countValue = float64(metric.GetHistogram().GetSampleCount())

	for i, bucket := range metric.GetHistogram().GetBucket() {
		bucketLEName := makeGraphiteBucketLEName(h.name, i)
		bucketCountName := makeGraphiteBucketCountName(h.name, i)

		h.bucketLEs[bucketLEName] = bucket.GetUpperBound()
		h.bucketCounts[bucketCountName] = float64(bucket.GetCumulativeCount())
	}
}

func (h *GraphiteHistogram) Update(metric *dto.Metric) {
	h.setTimeNow()
	h.sumValue = metric.GetHistogram().GetSampleSum()
	h.countValue = float64(metric.GetHistogram().GetSampleCount())

	for i, bucket := range metric.GetHistogram().GetBucket() {
		bucketLEName := makeGraphiteBucketLEName(h.name, i)
		if bucketLE, ok := h.bucketLEs[bucketLEName]; ok {
			bucketLE = bucket.GetUpperBound()
		} else {
			bucketLE = bucket.GetUpperBound()
			h.bucketLEs[bucketLEName] = bucketLE
		}

		bucketCountName := makeGraphiteBucketCountName(h.name, i)
		if bucketCount, ok := h.bucketCounts[bucketCountName]; ok {
			bucketCount = float64(bucket.GetCumulativeCount())
		} else {
			bucketCount = float64(bucket.GetCumulativeCount())
			h.bucketCounts[bucketCountName] = bucketCount
		}
	}
}

func (h *GraphiteHistogram) Export(exporter *GraphiteExporter) error {
	var metricsToSend []graphite.Metric
	sumMetric := graphite.Metric{
		Name:      makeGraphiteSumName(h.name),
		Value:     floatToString(h.sumValue),
		Timestamp: h.updateTime,
	}
	metricsToSend = append(metricsToSend, sumMetric)

	countMetric := graphite.Metric{
		Name:      makeGraphiteCountName(h.name),
		Value:     floatToString(h.countValue),
		Timestamp: h.updateTime,
	}
	metricsToSend = append(metricsToSend, countMetric)

	for bucketLEName, bucketLE := range h.bucketLEs {
		bucketLEMetric := graphite.Metric{
			Name:      bucketLEName,
			Value:     floatToString(bucketLE),
			Timestamp: h.updateTime,
		}
		metricsToSend = append(metricsToSend, bucketLEMetric)
	}

	for bucketCountName, bucketCount := range h.bucketCounts {
		bucketCountMetric := graphite.Metric{
			Name:      bucketCountName,
			Value:     floatToString(bucketCount),
			Timestamp: h.updateTime,
		}
		metricsToSend = append(metricsToSend, bucketCountMetric)
	}
	err := exporter.graphite.SendMetrics(metricsToSend)
	if err != nil {
		return fmt.Errorf("error sending histogram metrics to graphite: %v", err)
	}
	return nil
}

func floatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func makeGraphiteSumName(name string) string {
	firstTagIndex := strings.Index(name, ";")
	if firstTagIndex == -1 {
		return fmt.Sprintf("%s_sum", name)
	}
	return fmt.Sprintf("%s_sum%s", name[:firstTagIndex], name[firstTagIndex:])
}

func makeGraphiteCountName(name string) string {
	firstTagIndex := strings.Index(name, ";")
	if firstTagIndex == -1 {
		return fmt.Sprintf("%s_count", name)
	}
	return fmt.Sprintf("%s_count%s", name[:firstTagIndex], name[firstTagIndex:])
}

func makeGraphiteBucketLEName(name string, bucketNum int) string {
	firstTagIndex := strings.Index(name, ";")
	if firstTagIndex == -1 {
		return fmt.Sprintf("%s_bucket_%d_le", name, bucketNum)
	}
	return fmt.Sprintf("%s_bucket_%d_le%s", name[:firstTagIndex], bucketNum, name[firstTagIndex:])
}

func makeGraphiteBucketCountName(name string, bucketNum int) string {
	firstTagIndex := strings.Index(name, ";")
	if firstTagIndex == -1 {
		return fmt.Sprintf("%s_bucket_%d_count", name, bucketNum)
	}
	return fmt.Sprintf("%s_bucket_%d_count%s", name[:firstTagIndex], bucketNum, name[firstTagIndex:])
}
