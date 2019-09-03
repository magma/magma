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
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/prometheus/common/expfmt"
)

// Amount of metric datapoints in the cache when
// testing Receive and Scrape
var (
	powersOfTenToTest = []int{0, 1, 2}
	numBucketsToTest  = [...]int{500, 1000, 2000, 20000, 100000}
)

const (
	templateMetric  = "http_requests_total_%d{method=\"post\",code=\"400\"}    3 %d\n"
	maxSizeOfString = 10000000
)

func BenchmarkReceiveMetrics(b *testing.B) {
	familiesMap := prepareNewFamiliesMap(powersOfTenToTest)

	cache := NewMetricCache(0, 10)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(generateRandomMetricsString(0)))
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	for _, k := range powersOfTenToTest {
		num := int(math.Pow(10, float64(k)))
		for _, buckets := range numBucketsToTest {
			currBuckets := len(familiesMap[k])
			insertNRecordsIntoCacheBucketRange(cache, num, currBuckets, buckets)

			b.Run(fmt.Sprintf("%d-Datapoints-%d-Families", num, buckets), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					cache.metricFamiliesByName = familiesMap[k]
					cache.Receive(c)
				}
			})
		}
	}
}

func BenchmarkScrapeMetrics(b *testing.B) {
	familiesMap := prepareNewFamiliesMap(powersOfTenToTest)

	cache := NewMetricCache(0, 10)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	for _, k := range powersOfTenToTest {
		num := int(math.Pow(10, float64(k)))
		for _, buckets := range numBucketsToTest {
			currBuckets := len(familiesMap[k])
			cache.metricFamiliesByName = familiesMap[k]
			insertNRecordsIntoCacheBucketRange(cache, num, currBuckets, buckets)

			b.Run(fmt.Sprintf("%d-Datapoints-%d-Families", num, buckets), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					cache.metricFamiliesByName = familiesMap[k]
					cache.Scrape(c)
				}
			})
		}
	}
}

func generateRandomMetricsString(b int) string {
	timestamp := rand.Intn(10000000)
	return fmt.Sprintf(templateMetric, b, timestamp)
}

func generateNRandomMetricsStrings(n int, bucketStart int, bucketEnd int) string {
	var buf bytes.Buffer
	for b := bucketStart; b < bucketEnd; b++ {
		for i := 0; i < n; i++ {
			buf.WriteString(generateRandomMetricsString(b))
		}
	}

	return buf.String()
}

func prepareNewFamiliesMap(powersOfTen []int) map[int]map[string]*familyAndMetrics {
	familiesMap := make(map[int]map[string]*familyAndMetrics)

	for _, n := range powersOfTen {
		cache := NewMetricCache(0, 10)
		total := int(math.Pow(10, float64(n)))
		insertNRecordsIntoCacheBucketRange(cache, total, 0, numBucketsToTest[0])
		familiesMap[int(n)] = cache.metricFamiliesByName
	}
	return familiesMap
}

func insertNRecordsIntoCacheBucketRange(cache *MetricCache, total int, bucketStart int, bucketEnd int) {
	var parser expfmt.TextParser

	done := 0

	for done < total {
		records := int(math.Min(maxSizeOfString, float64(total-done)))
		metrics := generateNRandomMetricsStrings(records, bucketStart, bucketEnd)
		done += records
		parsedFamilies, err := parser.TextToMetricFamilies(strings.NewReader(metrics))

		if err != nil {
			fmt.Printf("Bad parsing: %s", err)
		}

		cache.cacheMetrics(parsedFamilies)
	}
}
