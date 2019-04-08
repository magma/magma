/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (client *MockClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	args := client.Called(url, data)
	return args.Get(0).(*http.Response), args.Error(1)
}

func makeStringPointer(s string) *string {
	return &s
}

func makeMetricTypePointer(t dto.MetricType) *dto.MetricType {
	return &t
}

func TestODSSubmit(t *testing.T) {
	exporter := exporters.NewODSExporter(
		"",
		"",
		"",
		"magma",
		2,
		time.Second*10,
	)

	singleMetricTestFamily := makeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{})
	entity := "testId1.testId2"
	context := exporters.MetricsContext{DecodedName: "test", OriginatingEntity: entity}
	err := exporter.Submit([]exporters.MetricAndContext{{Family: singleMetricTestFamily, Context: context}})
	assert.NoError(t, err)

	// Submitting to a full queue should drop metrics
	multiMetricTestFamily := makeTestMetricFamily(dto.MetricType_GAUGE, 100, []*dto.LabelPair{})
	err = exporter.Submit([]exporters.MetricAndContext{
		{Family: multiMetricTestFamily, Context: context},
		{Family: singleMetricTestFamily, Context: context},
	})
	assert.EqualError(t, err, "ODS queue full, dropping 100 samples")

	err = exporter.Submit([]exporters.MetricAndContext{{Family: singleMetricTestFamily, Context: context}})
	assert.EqualError(t, err, "ODS queue full, dropping 1 samples")
}

func TestExport(t *testing.T) {
	exporter := exporters.NewODSExporter(
		"",
		"",
		"",
		"magma",
		2,
		time.Second*10,
	)
	entity := "testId1.testId2"
	nameStr := "test"
	tagLabelPair := dto.LabelPair{
		Name:  makeStringPointer("tags"),
		Value: makeStringPointer("Tag1,Tag2"),
	}
	sample := exporters.NewSample(nameStr, "0", int64(0), []*dto.LabelPair{&tagLabelPair}, entity)

	client := new(MockClient)
	resp := &http.Response{StatusCode: 200}
	datapoints := []exporters.ODSDatapoint{}
	datapoints = append(datapoints, exporters.ODSDatapoint{
		Entity: fmt.Sprintf("magma.%s.%s", "testId1", "testId2"),
		Key:    exporter.FormatKey(sample),
		Tags:   exporter.GetTags(sample),
		Value:  sample.Value()})
	datapointsJson, err := json.Marshal(datapoints)
	assert.NoError(t, err)
	client.On("PostForm", mock.AnythingOfType("string"), url.Values{"datapoints": {string(datapointsJson)}}).Return(resp, nil)

	// Export called on empty queue
	err = exporter.Export(client)
	assert.NoError(t, err)
	client.AssertNotCalled(t, "PostForm", mock.AnythingOfType("string"), mock.AnythingOfType("url.Values"))

	testFamily := makeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{&tagLabelPair})
	context := exporters.MetricsContext{DecodedName: "test", OriginatingEntity: entity, NetworkID: "nID", GatewayID: "gID"}
	err = exporter.Submit([]exporters.MetricAndContext{{Family: testFamily, Context: context}})
	assert.NoError(t, err)

	err = exporter.Export(client)
	assert.NoError(t, err)
	client.AssertExpectations(t)

	// Fill queue (drop some samples), assert we didn't exceed queue length cap
	multiTestFamily := makeTestMetricFamily(dto.MetricType_GAUGE, 100, []*dto.LabelPair{&tagLabelPair})
	err = exporter.Submit([]exporters.MetricAndContext{{Family: multiTestFamily, Context: context}})
	assert.EqualError(t, err, "ODS queue full, dropping 98 samples")

	// We expect 2 samples to be exported
	datapoints = append(datapoints, datapoints...)
	datapointsJson, err = json.Marshal(datapoints)
	assert.NoError(t, err)
	client.On("PostForm", mock.Anything, url.Values{"datapoints": {string(datapointsJson)}}).Return(resp, nil)

	err = exporter.Export(client)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func makeTestMetricFamily(metricType dto.MetricType, count int, labels []*dto.LabelPair) *dto.MetricFamily {
	var testMetric dto.Metric
	switch metricType {
	case dto.MetricType_COUNTER:
		testMetric = makePromoCounter(0)
	case dto.MetricType_SUMMARY:
		testMetric = makePromoSummary(map[float64]float64{0.1: 0.01}, []float64{})
	case dto.MetricType_HISTOGRAM:
		testMetric = makePromoHistogram([]float64{1, 5, 10}, []float64{})
	default:
		testMetric = makePromoGauge(0)
	}

	testMetric.Label = labels
	metrics := make([]*dto.Metric, 0, count)
	for i := 0; i < count; i++ {
		metrics = append(metrics, &testMetric)
	}
	return &dto.MetricFamily{
		Name:   makeStringPointer("testFamily"),
		Help:   makeStringPointer("testFamilyHelp"),
		Type:   makeMetricTypePointer(metricType),
		Metric: metrics,
	}
}

func makePromoGauge(value float64) dto.Metric {
	var metric dto.Metric
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: "testGauge", Help: "testGaugeHelp"})
	gauge.Set(value)
	gauge.Write(&metric)
	return metric
}

func makePromoCounter(value float64) dto.Metric {
	var metric dto.Metric
	counter := prometheus.NewCounter(prometheus.CounterOpts{Name: "testCounter", Help: "testCounterHelp"})
	counter.Add(value)
	counter.Write(&metric)
	return metric
}

func makePromoSummary(objectives map[float64]float64, observations []float64) dto.Metric {
	var metric dto.Metric
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "testSummary",
			Help:       "testSummaryHelp",
			Objectives: objectives,
		},
	)
	for _, obs := range observations {
		summary.Observe(obs)
	}
	summary.Write(&metric)
	return metric
}

func makePromoHistogram(buckets []float64, observations []float64) dto.Metric {
	var metric dto.Metric
	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "testHistogram",
			Help:    "testHistogramHelp",
			Buckets: buckets,
		},
	)
	for _, obs := range observations {
		histogram.Observe(obs)
	}
	histogram.Write(&metric)
	return metric
}

func TestFormatKey(t *testing.T) {
	magmad_test_init.StartTestService(t)

	// Init exporter
	exporter := exporters.NewODSExporter(
		"",
		"",
		"",
		"magma",
		2,
		time.Second*10,
	)
	entity := "testId1.testId2"
	// Test where key should be prepended with service name
	testSampleWithService := exporters.NewSample(
		"test_sample_with_service",
		"val",
		int64(0),
		[]*dto.LabelPair{
			{
				Name:  makeStringPointer("service"),
				Value: makeStringPointer("mme"),
			},
			{
				Name:  makeStringPointer("result"),
				Value: makeStringPointer("success"),
			},
			{
				Name:  makeStringPointer("cause"),
				Value: makeStringPointer("foo"),
			},
		},
		entity,
	)

	key := exporter.FormatKey(testSampleWithService)
	assert.Equal(t, key, "mme.test_sample_with_service.result-success.cause-foo")

	// Test where no service name provided
	testSampleNoService := exporters.NewSample(
		"test_sample_no_service",
		"val",
		int64(0),
		[]*dto.LabelPair{
			{
				Name:  makeStringPointer("result"),
				Value: makeStringPointer("success"),
			},
			{
				Name:  makeStringPointer("cause"),
				Value: makeStringPointer("foo"),
			},
		},
		entity,
	)

	key = exporter.FormatKey(testSampleNoService)
	assert.Equal(t, key, "test_sample_no_service.result-success.cause-foo")

	// Test where tags are provided and no service name is provided.
	testSampleWithTags := exporters.NewSample(
		"test_sample_with_tags",
		"val",
		int64(0),
		[]*dto.LabelPair{
			{
				Name:  makeStringPointer("results"),
				Value: makeStringPointer("success"),
			},
			{
				Name:  makeStringPointer("cause"),
				Value: makeStringPointer("foo"),
			},
			{
				Name:  makeStringPointer("tags"),
				Value: makeStringPointer("Magma"),
			},
		},
		entity,
	)

	key = exporter.FormatKey(testSampleWithTags)
	assert.Equal(t, key, "test_sample_with_tags.results-success.cause-foo")

	// Test where both tags and service name is provided.
	testSampleWithTagsAndService := exporters.NewSample(
		"test_sample_with_tags_and_service",
		"val",
		int64(0),
		[]*dto.LabelPair{
			{
				Name:  makeStringPointer("service"),
				Value: makeStringPointer("mme"),
			},
			{
				Name:  makeStringPointer("results"),
				Value: makeStringPointer("success"),
			},
			{
				Name:  makeStringPointer("cause"),
				Value: makeStringPointer("foo"),
			},
			{
				Name:  makeStringPointer("tags"),
				Value: makeStringPointer("Magma"),
			},
		},
		entity,
	)

	key = exporter.FormatKey(testSampleWithTagsAndService)
	assert.Equal(t, key, "mme.test_sample_with_tags_and_service.results-success.cause-foo")
}

func TestFormatTags(t *testing.T) {
	magmad_test_init.StartTestService(t)

	// Init exporter
	exporter := exporters.NewODSExporter(
		"",
		"",
		"",
		"magma",
		2,
		time.Second*10,
	)
	entity := "testId1.testId2"
	// Test where key should be prepended with service name
	testSampleWithService := exporters.NewSample(
		"test_sample_with_service",
		"val",
		int64(0),
		[]*dto.LabelPair{
			{
				Name:  makeStringPointer("service"),
				Value: makeStringPointer("mme"),
			},
			{
				Name:  makeStringPointer("result"),
				Value: makeStringPointer("success"),
			},
			{
				Name:  makeStringPointer("cause"),
				Value: makeStringPointer("foo"),
			},
			{
				Name:  makeStringPointer("tags"),
				Value: makeStringPointer("Magma,Bootcamp"),
			},
		},
		entity,
	)

	tag := exporter.GetTags(testSampleWithService)
	assert.Equal(t, tag, []string{"Magma", "Bootcamp"})
}
