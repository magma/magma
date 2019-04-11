/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"math"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/test_common"

	"github.com/stretchr/testify/assert"
)

var (
	exporter = NewGraphiteExporter("", 0).(*GraphiteExporter)
)

func TestGraphiteGauge_Register(t *testing.T) {
	g := NewGraphiteGauge()

	testGaugeValue := 123.0
	gauge := test_common.MakePromoGauge(testGaugeValue)
	g.Register(&gauge, "testGauge", exporter)

	compareFloats(t, 123.0, g.(*GraphiteGauge).value)
}

func TestGraphiteGauge_Update(t *testing.T) {
	g := NewGraphiteGauge()
	gauge := test_common.MakePromoGauge(0.0)

	g.Register(&gauge, "testGauge", exporter)

	updatedGaugeValue := 123.0
	gauge.GetGauge().Value = &updatedGaugeValue
	g.Update(&gauge)

	compareFloats(t, 123.0, g.(*GraphiteGauge).value)
}

func TestGraphiteCounter_Register(t *testing.T) {
	c := NewGraphiteCounter()

	testValue := 123.0
	counter := test_common.MakePromoCounter(testValue)
	c.Register(&counter, "testCounter", exporter)

	compareFloats(t, 123.0, c.(*GraphiteCounter).value)
}

func TestGraphiteCounter_Update(t *testing.T) {
	c := NewGraphiteCounter()

	counter := test_common.MakePromoCounter(0.0)
	c.Register(&counter, "testCounter", exporter)

	updatedValue := 123.0
	counter.Counter.Value = &updatedValue
	c.Update(&counter)

	compareFloats(t, 123.0, c.(*GraphiteCounter).value)
}

func TestGraphiteSummary_Register(t *testing.T) {
	s := NewGraphiteSummary()

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	observations := []float64{0.5, 0.6, 0.7, 0.8}
	summary := test_common.MakePromoSummary(objectives, observations)

	metricName := "testSummary"
	s.Register(&summary, metricName, exporter)

	expectedCount := float64(len(observations))
	expectedSum := 0.5 + 0.6 + 0.7 + 0.8
	assert.Equal(t, expectedCount, s.(*GraphiteSummary).countValue)
	compareFloats(t, expectedSum, s.(*GraphiteSummary).sumValue)
}

func TestGraphiteSummary_Update(t *testing.T) {
	s := NewGraphiteSummary()

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	observations := []float64{0.5, 0.6, 0.7}
	summary := test_common.MakePromoSummary(objectives, observations)

	metricName := "testSummary"
	s.Register(&summary, metricName, exporter)

	newObservations := []float64{0.8, 0.9}
	for _, obs := range newObservations {
		observations = append(observations, obs)
	}
	updatedSummary := test_common.MakePromoSummary(objectives, observations)
	s.Update(&updatedSummary)

	expectedCount := float64(len(observations))
	expectedSum := 0.5 + 0.6 + 0.7 + 0.8 + 0.9
	assert.Equal(t, expectedCount, s.(*GraphiteSummary).countValue)
	compareFloats(t, expectedSum, s.(*GraphiteSummary).sumValue)
}

func TestGraphiteHistogram_Register(t *testing.T) {
	metricBaseName := "testBaseName"
	h := NewGraphiteHistogram()

	buckets := []float64{1.0, 5.0, 10.0}
	observations := []float64{0.5, 0.8, 2.0, 7.2, 9.2}
	histogram := test_common.MakePromoHistogram(buckets, observations)

	h.Register(&histogram, metricBaseName, exporter)
	expectedCount := float64(len(observations))
	expectedSum := 0.5 + 0.8 + 2.0 + 7.2 + 9.2
	assert.Equal(t, expectedCount, h.(*GraphiteHistogram).countValue)
	compareFloats(t, expectedSum, h.(*GraphiteHistogram).sumValue)
}

func TestGraphiteHistogram_Update(t *testing.T) {
	metricBaseName := "testBaseName"
	h := NewGraphiteHistogram()

	buckets := []float64{1.0, 5.0, 10.0}
	observations := []float64{0.5, 0.8, 2.0, 7.2, 9.2}
	histogram := test_common.MakePromoHistogram(buckets, observations)

	h.Register(&histogram, metricBaseName, exporter)

	newObservations := []float64{0.4, 2.5, 8.0}
	for _, obs := range newObservations {
		observations = append(observations, obs)
	}
	updatedHistogram := test_common.MakePromoHistogram(buckets, observations)

	h.Update(&updatedHistogram)

	expectedCount := float64(len(observations))
	expectedSum := float64(0.5 + 0.8 + 2.0 + 7.2 + 9.2 + 0.4 + 2.5 + 8.0)
	assert.Equal(t, expectedCount, h.(*GraphiteHistogram).countValue)
	compareFloats(t, expectedSum, h.(*GraphiteHistogram).sumValue)
}

func compareFloats(t *testing.T, x, y float64) {
	if math.Abs(x-y) > 0.0001 {
		t.Fail()
	}
}

const (
	testNameWithTags = "metricName;tag1=x;tag2=y"
	testNameNoTags   = "metricName"
)

func TestMakeGraphiteSumName(t *testing.T) {
	sumWithTags := makeGraphiteSumName(testNameWithTags)
	assert.Equal(t, "metricName_sum;tag1=x;tag2=y", sumWithTags)
	sumNoTags := makeGraphiteSumName(testNameNoTags)
	assert.Equal(t, "metricName_sum", sumNoTags)
}

func TestMakeGraphiteCountName(t *testing.T) {
	countWithTags := makeGraphiteCountName(testNameWithTags)
	assert.Equal(t, "metricName_count;tag1=x;tag2=y", countWithTags)
	countNoTags := makeGraphiteCountName(testNameNoTags)
	assert.Equal(t, "metricName_count", countNoTags)
}

func TestMakeGraphiteBucketLEName(t *testing.T) {
	leWithTags := makeGraphiteBucketLEName(testNameWithTags, 0)
	assert.Equal(t, "metricName_bucket_0_le;tag1=x;tag2=y", leWithTags)
	leNoTags := makeGraphiteBucketLEName(testNameNoTags, 0)
	assert.Equal(t, "metricName_bucket_0_le", leNoTags)
}

func TestMakeGraphiteBucketCountName(t *testing.T) {
	countWithTags := makeGraphiteBucketCountName(testNameWithTags, 0)
	assert.Equal(t, "metricName_bucket_0_count;tag1=x;tag2=y", countWithTags)
	countNoTags := makeGraphiteBucketCountName(testNameNoTags, 0)
	assert.Equal(t, "metricName_bucket_0_count", countNoTags)
}
