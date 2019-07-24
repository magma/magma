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
	"os"
	"sort"
	"strings"
	"sync"
	"time"

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
	limit     int
	stats     cacheStats
	sync.Mutex
}

type cacheStats struct {
	lastScrapeTime        int64
	lastScrapeSize        int64
	lastScrapeNumFamilies int

	lastReceiveTime        int64
	lastReceiveSize        int64
	lastReceiveNumFamilies int

	currentCountFamilies   int
	currentCountSeries     int
	currentCountDatapoints int
}

func NewMetricCache(limit int) *MetricCache {
	if limit > 0 {
		glog.Infof("Prometheus-Cache created with a limit of %d\n", limit)
	} else {
		glog.Info("Prometheus-Cache created with no limit\n")
	}

	return &MetricCache{
		familyMap: make(map[string]*familyAndMetrics),
		limit:     limit,
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

	// Check if new datapoints will exceed the specified limit
	if c.limit > 0 {
		newDatapoints := 0
		for _, fam := range parsedFamilies {
			newDatapoints += len(fam.Metric)
		}
		if c.stats.currentCountDatapoints+newDatapoints > c.limit {
			fmt.Println("Not accepting push. Would overfill cache limit")
			return ctx.NoContent(http.StatusNotAcceptable)
		}
	}

	c.cacheMetrics(parsedFamilies)

	c.stats.lastReceiveTime = time.Now().Unix()
	c.stats.lastReceiveSize = ctx.Request().ContentLength
	c.stats.lastReceiveNumFamilies = len(parsedFamilies)
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

	expositionString := c.exposeMetrics(scrapeMetrics)
	c.stats.lastScrapeTime = time.Now().Unix()
	c.stats.lastScrapeSize = int64(len(expositionString))
	c.stats.lastScrapeNumFamilies = len(scrapeMetrics)

	return ctx.String(http.StatusOK, expositionString)
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

// Debug is a handler function to show the current state of the cache without
// consuming any datapoints
func (c *MetricCache) Debug(ctx echo.Context) error {
	c.updateCountStats()
	hostname, _ := os.Hostname()
	debugString := fmt.Sprintf(`Prometheus Cache running on %s
Last Scrape: %d
	Scrape Size: %d
	Number of Familes: %d

Last Receive: %d
	Receive Size: %d
	Number of Families: %d

Current Count Families:   %d
Current Count Series:     %d
Current Count Datapoints: %d

Current Exposition Text:

%s`, hostname, c.stats.lastScrapeTime, c.stats.lastScrapeSize, c.stats.lastScrapeNumFamilies,
		c.stats.lastReceiveTime, c.stats.lastReceiveSize, c.stats.lastReceiveNumFamilies,
		c.stats.currentCountFamilies, c.stats.currentCountSeries, c.stats.currentCountDatapoints, c.exposeMetrics(c.familyMap))

	return ctx.String(http.StatusOK, debugString)
}

func (c *MetricCache) updateCountStats() {
	numFamilies := len(c.familyMap)
	numSeries := 0
	numDatapoints := 0
	for _, family := range c.familyMap {
		numSeries += len(family.metrics)
		for _, series := range family.metrics {
			numDatapoints += len(series)
		}
	}
	c.stats.currentCountFamilies = numFamilies
	c.stats.currentCountSeries = numSeries
	c.stats.currentCountDatapoints = numDatapoints
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
