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
	familyMap map[string]*familyAndMetrics
	sync.Mutex
}

func NewMetricCache() *MetricCache {
	return &MetricCache{
		familyMap: make(map[string]*familyAndMetrics),
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
	c.cacheMetrics(parsedFamilies)
	return ctx.NoContent(http.StatusOK)
}

func (c *MetricCache) cacheMetrics(families map[string]*dto.MetricFamily) {
	c.Lock()
	defer c.Unlock()
	for _, fam := range families {
		if fAndM, ok := c.familyMap[fam.GetName()]; ok {
			fAndM.addMetrics(fam.Metric)
		} else {
			c.familyMap[fam.GetName()] = newFamilyAndMetrics(fam)
		}
	}
}

// Scrape is a handler function for prometheus scrape requests. Formats the
// metrics for scraping.
func (c *MetricCache) Scrape(ctx echo.Context) error {
	c.Lock()
	scrapeMetrics := c.familyMap
	c.clearMetrics()
	c.Unlock()

	return ctx.String(http.StatusOK, c.exposeMetrics(scrapeMetrics))
}

func (c *MetricCache) clearMetrics() {
	c.familyMap = make(map[string]*familyAndMetrics)
}

func (c *MetricCache) exposeMetrics(familyMap map[string]*familyAndMetrics) string {
	respStr := strings.Builder{}
	for _, fam := range familyMap {
		pullFamily := fam.popSortedDatapoints()
		familyStr, err := familyToString(pullFamily)
		if err != nil {
			glog.Errorf("metric %s dropped. error converting metric to string: %v", *pullFamily.Name, err)
		} else {
			respStr.WriteString(familyStr)
		}
	}
	return respStr.String()
}

type familyAndMetrics struct {
	family  *dto.MetricFamily
	metrics map[string][]*dto.Metric
}

func newFamilyAndMetrics(family *dto.MetricFamily) *familyAndMetrics {
	metricMap := make(map[string][]*dto.Metric)
	for _, metric := range family.Metric {
		name := makeLabeledName(metric, family.GetName())
		if metricQueue, ok := metricMap[name]; ok {
			metricMap[name] = append(metricQueue, metric)
		} else {
			metricMap[name] = []*dto.Metric{metric}
		}
	}
	// clear metrics in family because we are keeping them in the queues
	family.Metric = []*dto.Metric{}

	return &familyAndMetrics{
		family:  family,
		metrics: metricMap,
	}
}

func (f *familyAndMetrics) addMetrics(newMetrics []*dto.Metric) {
	for _, metric := range newMetrics {
		metricName := makeLabeledName(metric, f.family.GetName())
		if queue, ok := f.metrics[metricName]; ok {
			f.metrics[metricName] = append(queue, metric)
		} else {
			f.metrics[metricName] = []*dto.Metric{metric}
		}
	}
}

// Returns a prometheus MetricFamily populated with all datapoints, sorted so
// that the earliest datapoint appears first
func (f *familyAndMetrics) popSortedDatapoints() *dto.MetricFamily {
	pullFamily := f.copyFamily()
	for _, queue := range f.metrics {
		if len(queue) == 0 {
			continue
		}
		// Sort metrics by timestamp
		sort.Slice(queue, func(i, j int) bool {
			return *queue[i].TimestampMs < *queue[j].TimestampMs
		})
		pullFamily.Metric = append(pullFamily.Metric, queue...)
	}
	return &pullFamily
}

// return a copy of the MetricFamily that can be modified safely
func (f *familyAndMetrics) copyFamily() dto.MetricFamily {
	return *f.family
}

// makeLabeledName builds a unique name from a metric LabelPairs
func makeLabeledName(metric *dto.Metric, metricName string) string {
	labels := metric.GetLabel()
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].GetName() < labels[j].GetName()
	})

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
