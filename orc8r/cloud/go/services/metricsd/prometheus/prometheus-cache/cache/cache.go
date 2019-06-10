/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package cache

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// MetricCache serves as a replacement for the prometheus pushgateway. Accepts
// timestamps with metrics, and stores them in a queue to allow multiple
// datapoints per metric series to be scraped
type MetricCache struct {
	familyMap     map[string]*familyAndMetrics
	queueCapacity int
	sync.Mutex
}

func NewMetricCache(queueCapacity int) *MetricCache {
	return &MetricCache{
		familyMap:     make(map[string]*familyAndMetrics),
		queueCapacity: queueCapacity,
	}
}

// Receive is a handler function to receive metric pushes
func (c *MetricCache) Receive(ctx echo.Context) error {
	var parser expfmt.TextParser
	var err error

	parsedFamilies, err := parser.TextToMetricFamilies(ctx.Request().Body)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("error parsing metrics: %v", err))
	}

	c.Lock()
	defer c.Unlock()
	for _, fam := range parsedFamilies {
		if fAndM, ok := c.familyMap[fam.GetName()]; ok {
			fAndM.addMetrics(fam.Metric)
		} else {
			c.familyMap[fam.GetName()] = newFamilyAndMetrics(fam, c.queueCapacity)
		}
	}
	return ctx.NoContent(http.StatusOK)
}

// Scrape is a handler function for prometheus scrape requests. Formats the
// metrics for scraping and sends the oldest datapoint per metric series if
// there are multiple in the queue since prometheus cannot handle multiple
// datapoints of the same metric in a scrape.
func (c *MetricCache) Scrape(ctx echo.Context) error {
	c.Lock()
	defer c.Unlock()
	respStr := strings.Builder{}
	for _, fam := range c.familyMap {
		pullFamily := fam.popOldestDatapoint()
		familyStr, err := familyToString(pullFamily)
		if err != nil {
			glog.Errorf("metric %s dropped. error converting metric to string: %v", *pullFamily.Name, err)
		} else {
			respStr.WriteString(familyStr)
		}
	}
	c.clearEmptyMetrics()
	return ctx.String(http.StatusOK, respStr.String())
}

// clearEmptyMetrics deletes families from the cache if they have no more
// datapoints
func (c *MetricCache) clearEmptyMetrics() {
	for familyName, family := range c.familyMap {
		family.prune()
		if len(family.metrics) == 0 {
			delete(c.familyMap, familyName)
		}
	}
}

// familyAndMetrics stores the metrics in a MetricFamily in a MetricQueue
type familyAndMetrics struct {
	family        *dto.MetricFamily
	metrics       map[string]*MetricQueue
	queueCapacity int
}

func newFamilyAndMetrics(family *dto.MetricFamily, queueCapacity int) *familyAndMetrics {
	metricMap := make(map[string]*MetricQueue)
	for _, metric := range family.Metric {
		name := makeLabeledName(metric, family.GetName())
		if metricQueue, ok := metricMap[name]; ok {
			metricQueue.Push(metric)
		} else {
			queue := NewMetricQueue(queueCapacity)
			queue.Push(metric)
			metricMap[name] = queue
		}
	}
	// clear metrics in family because we are keeping them in the queues
	family.Metric = []*dto.Metric{}

	return &familyAndMetrics{
		family:        family,
		metrics:       metricMap,
		queueCapacity: queueCapacity,
	}
}

func (f *familyAndMetrics) addMetrics(newMetrics []*dto.Metric) {
	for _, metric := range newMetrics {
		metricName := makeLabeledName(metric, f.family.GetName())
		if queue, ok := f.metrics[metricName]; ok {
			queue.Push(metric)
		} else {
			f.metrics[metricName] = NewMetricQueue(f.queueCapacity)
			f.metrics[metricName].Push(metric)
		}
	}
}

// Returns a prometheus MetricFamily populated with the oldest datapoint for
// each of the unique metric series stored in the struct. Pops the oldest point
// off of each queue.
func (f *familyAndMetrics) popOldestDatapoint() *dto.MetricFamily {
	pullFamily := f.copyFamily()
	for _, queue := range f.metrics {
		poppedMetric := queue.Pop()
		if poppedMetric == nil {
			continue
		}
		pullFamily.Metric = append(pullFamily.Metric, poppedMetric)
	}
	return &pullFamily
}

// return a copy of the MetricFamily that can be modified safely
func (f *familyAndMetrics) copyFamily() dto.MetricFamily {
	return *f.family
}

// prune deletes metricQueues which no longer have any datapoints
func (f *familyAndMetrics) prune() {
	for metricName, metricQueue := range f.metrics {
		if metricQueue.size == 0 {
			delete(f.metrics, metricName)
		}
	}
}

// makeLabeledName builds a unique name from a metric LabelPairs
func makeLabeledName(metric *dto.Metric, metricName string) string {
	labels := metric.GetLabel()
	sort.Sort(ByName(labels))

	labeledName := strings.Builder{}
	labeledName.WriteString(metricName)
	for _, labelPair := range labels {
		labeledName.WriteString(fmt.Sprintf("_%s_%s", labelPair.GetName(), labelPair.GetValue()))
	}
	return labeledName.String()
}

func familyToString(family *dto.MetricFamily) (string, error) {
	var buf bytes.Buffer
	_, err := expfmt.MetricFamilyToText(&buf, family)
	if err != nil {
		return "", fmt.Errorf("error writing family string: %v", err)
	}
	return buf.String(), nil
}

// MetricQueue implements a fixed-size moving window queue so that in the case
// of more frequent pushes than scrapes, memory usage does not increase constantly.
type MetricQueue struct {
	capacity int
	size     int
	writeIdx int
	queue    []*dto.Metric
}

func NewMetricQueue(capacity int) *MetricQueue {
	return &MetricQueue{
		capacity: capacity,
		size:     0,
		writeIdx: 0,
		queue:    make([]*dto.Metric, capacity),
	}
}

// Push adds a metric to the queue. The oldest metric is overwritten when
// called while the buffer is full.
func (q *MetricQueue) Push(metric *dto.Metric) {
	q.queue[q.writeIdx] = metric
	q.writeIdx = (q.writeIdx + 1) % q.capacity
	if q.size < q.capacity {
		q.size++
	}
}

// Pop returns the oldest written metric. Overwrites the oldest metric when
func (q *MetricQueue) Pop() *dto.Metric {
	if q.size == 0 {
		return nil
	}
	readIdx := int(math.Abs(float64((q.writeIdx - q.size) % q.capacity)))
	metric := q.queue[readIdx]
	q.size--
	return metric
}

// ByName is an interface for sorting LabelPairs by name
type ByName []*dto.LabelPair

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].GetName() < a[j].GetName() }
