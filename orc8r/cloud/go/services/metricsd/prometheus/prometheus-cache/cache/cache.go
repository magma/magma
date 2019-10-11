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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

const (
	internalMetricCacheSize  = "cache_size"
	internalMetricCacheLimit = "cache_limit"
	scrapeWorkerPoolSize     = 100
)

// MetricCache serves as a replacement for the prometheus pushgateway. Accepts
// timestamps with metrics, and stores them in a queue to allow multiple
// datapoints per metric series to be scraped
type MetricCache struct {
	metricFamiliesByName map[string]*familyAndMetrics
	internalMetrics      map[string]prometheus.Gauge
	limit                int
	stats                cacheStats
	sync.Mutex
	scrapeTimeout int
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

func NewMetricCache(limit int, scrapeTimeout int) *MetricCache {
	if limit > 0 {
		glog.Infof("Prometheus-Cache created with a limit of %d\n", limit)
	} else {
		glog.Info("Prometheus-Cache created with no limit\n")
	}

	cacheLimit := prometheus.NewGauge(prometheus.GaugeOpts{Name: internalMetricCacheLimit, Help: "Maximum number of datapoints in cache", ConstLabels: prometheus.Labels{"networkID": "internal"}})
	cacheSize := prometheus.NewGauge(prometheus.GaugeOpts{Name: internalMetricCacheSize, Help: "Number of datapoints in cache", ConstLabels: prometheus.Labels{"networkID": "internal"}})
	internalMetrics := map[string]prometheus.Gauge{internalMetricCacheLimit: cacheLimit, internalMetricCacheSize: cacheSize}

	cacheLimit.Set(float64(limit))

	return &MetricCache{
		metricFamiliesByName: make(map[string]*familyAndMetrics),
		internalMetrics:      internalMetrics,
		limit:                limit,
		scrapeTimeout:        scrapeTimeout,
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

	newDatapoints := 0
	for _, fam := range parsedFamilies {
		newDatapoints += len(fam.Metric)
	}

	// Check if new datapoints will exceed the specified limit
	if c.limit > 0 {
		if c.stats.currentCountDatapoints+newDatapoints > c.limit {
			errString := fmt.Sprintf("Not accepting push of size %d. Would overfill cache limit of %d. Current cache size: %d\n", newDatapoints, c.limit, c.stats.currentCountDatapoints)
			glog.Error(errString)
			return ctx.String(http.StatusNotAcceptable, errString)
		}
	}

	c.cacheMetrics(parsedFamilies)

	c.stats.lastReceiveTime = time.Now().Unix()
	c.stats.lastReceiveSize = ctx.Request().ContentLength
	c.stats.lastReceiveNumFamilies = len(parsedFamilies)
	c.stats.currentCountDatapoints += newDatapoints
	c.internalMetrics[internalMetricCacheSize].Set(float64(c.stats.currentCountDatapoints))

	return ctx.NoContent(http.StatusOK)
}

func (c *MetricCache) cacheMetrics(families map[string]*dto.MetricFamily) {
	c.Lock()
	defer c.Unlock()
	for _, fam := range families {
		if families, ok := c.metricFamiliesByName[fam.GetName()]; ok {
			families.addMetrics(fam.Metric)
		} else {
			c.metricFamiliesByName[fam.GetName()] = newFamilyAndMetrics(fam)
		}
	}
}

// Scrape is a handler function for prometheus scrape requests. Formats the
// metrics for scraping.
func (c *MetricCache) Scrape(ctx echo.Context) error {
	c.Lock()
	scrapeMetrics := c.metricFamiliesByName
	c.clearMetrics()
	c.Unlock()

	expositionString := c.exposeMetrics(scrapeMetrics, scrapeWorkerPoolSize)
	expositionString += c.exposeInternalMetrics()

	c.stats.lastScrapeTime = time.Now().Unix()
	c.stats.lastScrapeSize = int64(len(expositionString))
	c.stats.lastScrapeNumFamilies = len(scrapeMetrics)
	c.stats.currentCountDatapoints = 0
	c.internalMetrics[internalMetricCacheSize].Set(0)

	return ctx.String(http.StatusOK, expositionString)
}

func (c *MetricCache) clearMetrics() {
	c.metricFamiliesByName = make(map[string]*familyAndMetrics)
}

func (c *MetricCache) exposeMetrics(metricFamiliesByName map[string]*familyAndMetrics, workers int) string {
	fams := make(chan *familyAndMetrics, workers)
	results := make(chan string, workers)
	respStrChannel := make(chan string, 1)

	waitGroup := &sync.WaitGroup{}

	for i := 0; i < workers; i++ {
		waitGroup.Add(1)
		go processFamilyWorker(fams, results, waitGroup)
	}

	go processFamilyStringsWorker(results, respStrChannel)

	for _, fam := range metricFamiliesByName {
		fams <- fam
	}

	close(fams)
	waitGroup.Wait()
	close(results)

	select {
	case respStr := <-respStrChannel:
		return respStr
	case <-time.After(time.Duration(c.scrapeTimeout) * time.Second):
		glog.Errorf("Timeout reached for building metrics string. Returning empty string.")
		return ""
	}
}

func processFamilyWorker(fams <-chan *familyAndMetrics, results chan<- string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for fam := range fams {
		pullFamily := fam.popSortedDatapoints()
		familyStr, err := familyToString(pullFamily)
		if err != nil {
			glog.Errorf("metric %s dropped. error converting metric to string: %v", *pullFamily.Name, err)
		} else {
			results <- familyStr
		}
	}
}

func processFamilyStringsWorker(results <-chan string, respStrChannel chan<- string) {
	respStr := strings.Builder{}

	for result := range results {
		respStr.WriteString(result)
	}
	respStrChannel <- respStr.String()
}

func (c *MetricCache) exposeInternalMetrics() string {
	strBuilder := strings.Builder{}
	for name, metric := range c.internalMetrics {
		str, err := writeInternalMetric(metric, name, dto.MetricType_GAUGE)
		if err != nil {
			continue
		}
		strBuilder.WriteString(str)
	}
	return strBuilder.String()
}

// Debug is a handler function to show the current state of the cache without
// consuming any datapoints
func (c *MetricCache) Debug(ctx echo.Context) error {
	verbose := ctx.QueryParam("verbose")

	c.updateCountStats()
	hostname, _ := os.Hostname()
	var limitValue, utilizationValue string
	if c.limit <= 0 {
		limitValue = "None"
		utilizationValue = "0"
	} else {
		limitValue = strconv.Itoa(c.limit)
		utilizationValue = strconv.FormatFloat(float64(c.stats.currentCountDatapoints)*100/float64(c.limit), 'f', 2, 64)
	}

	debugString := fmt.Sprintf(`Prometheus Cache running on %s
Cache Limit:       %s
Cache Utilization: %s%%

Last Scrape: %d
	Scrape Size: %d
	Number of Familes: %d

Last Receive: %d
	Receive Size: %d
	Number of Families: %d

Current Count Families:   %d
Current Count Series:     %d
Current Count Datapoints: %d `, hostname, limitValue, utilizationValue,
		c.stats.lastScrapeTime, c.stats.lastScrapeSize, c.stats.lastScrapeNumFamilies,
		c.stats.lastReceiveTime, c.stats.lastReceiveSize, c.stats.lastReceiveNumFamilies,
		c.stats.currentCountFamilies, c.stats.currentCountSeries, c.stats.currentCountDatapoints)

	if verbose != "" {
		debugString += fmt.Sprintf("\n\nCurrent Exposition Text:\n%s\n%s", c.exposeMetrics(c.metricFamiliesByName, scrapeWorkerPoolSize), c.exposeInternalMetrics())
	}

	return ctx.String(http.StatusOK, debugString)
}

func (c *MetricCache) updateCountStats() {
	numFamilies := len(c.metricFamiliesByName)
	numSeries := 0
	numDatapoints := 0
	for _, family := range c.metricFamiliesByName {
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

func writeInternalMetric(metric prometheus.Metric, name string, familyType dto.MetricType) (string, error) {
	var dtoMetric dto.Metric
	err := metric.Write(&dtoMetric)
	if err != nil {
		return "", err
	}
	fam := dto.MetricFamily{
		Name:   &name,
		Type:   &familyType,
		Metric: []*dto.Metric{&dtoMetric},
	}
	cacheSizeStr, err := familyToString(&fam)
	if err != nil {
		return "", err
	}
	return cacheSizeStr, nil
}
