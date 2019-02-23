/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicQuery(t *testing.T) {
	testQueryHelper(t, "up", []string{"up"})
}

func TestQueryWithFunction(t *testing.T) {
	testQueryHelper(t, "sum(up)", []string{"up"})
}

func TestQueryWithLabels(t *testing.T) {
	testQueryHelper(t, "up{existingLabelName=\"existingLabelValue\"}", []string{"up"})
}

func TestQueryWithMultipleMetrics(t *testing.T) {
	testQueryHelper(t, "metric1 or metric2", []string{"metric1", "metric2"})
}

func TestQueryWithMultipleMetricsAndLabels(t *testing.T) {
	testQueryHelper(t, "metric1 or metric2{existingLabelName=\"existingLabelValue\"}", []string{"metric1", "metric2"})
}

func TestQueryWithMatrixSelector(t *testing.T) {
	testQueryHelper(t, "up[5m]", []string{"up"})
}

func TestQueryWithMatrixAndFunctions(t *testing.T) {
	testQueryHelper(t, "sum_over_time(metric1[5m]) or sum_over_time(metric2[5m])", []string{"metric1", "metric2"})
}

func testQueryHelper(t *testing.T, query string, metricsInQuery []string) {
	singleLabel := map[string]string{"name1": "value1"}
	restrictedBasicQuery, err := createRestrictedQuery(query, singleLabel)
	assert.NoError(t, err)
	checkMetricsHaveLabels(t, metricsInQuery, restrictedBasicQuery, singleLabel)

	multipleLabels := map[string]string{"name1": "value1", "name2": "value2", "name3": "value3"}
	restrictedBasicQuery, err = createRestrictedQuery(query, multipleLabels)
	assert.NoError(t, err)
	checkMetricsHaveLabels(t, metricsInQuery, restrictedBasicQuery, multipleLabels)
}

// Asserts that each metric in a query is restricted with some labels
func checkMetricsHaveLabels(t *testing.T, metrics []string, query string, labels map[string]string) {
	labelsRegexString := "{(.*=\".*\")+(,.*=\".*\")*}"
	for _, metric := range metrics {
		metricRegex := regexp.MustCompile(metric + labelsRegexString)
		assert.Regexp(t, metricRegex, query)

		metricStrings := metricRegex.FindAllString(query, -1)
		for _, metricString := range metricStrings {
			for name, val := range labels {
				labelPair := fmt.Sprintf("%s=\"%s\"", name, val)
				assert.True(t, strings.Contains(metricString, labelPair))
			}
		}
	}
}

func createRestrictedQuery(query string, labels map[string]string) (string, error) {
	restrictor := NewQueryRestrictor(labels)
	return restrictor.RestrictQuery(query)
}
