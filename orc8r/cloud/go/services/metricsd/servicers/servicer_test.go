/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"strconv"
	"testing"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/servicers"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type TestMetricExporter struct {
	Queue []exporters.Sample
}

const MetricName = protos.MetricName_process_virtual_memory_bytes
const LabelName = protos.MetricLabelName_result
const LabelValue = "success"

func (e *TestMetricExporter) Submit(family *dto.MetricFamily, context exporters.MetricsContext) error {
	for _, metric := range family.GetMetric() {
		for _, s := range exporters.GetSamplesForMetrics(context.DecodedName, family.GetType(), metric, context.OriginatingEntity) {
			e.Queue = append(e.Queue, s)
		}
	}
	return nil
}

func (e *TestMetricExporter) Start() {}

func NewTestMetricExporter() *TestMetricExporter {
	e := new(TestMetricExporter)
	return e
}

func TestCollect(t *testing.T) {
	magmad_test_init.StartTestService(t)

	e := NewTestMetricExporter()
	ctx := context.Background()
	srv := servicers.NewMetricsControllerServer()
	srv.RegisterExporter(e)

	// Create test network
	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"},
		"metricsd_servicer_test_network")
	if err != nil {
		t.Fatalf("Magmad Register Network '%s' Error: %s", testNetworkId, err)
	}

	// Register a fake gateway
	gatewayId := "2876171d-bf38-4254-b4da-71a713952904"
	hwId := protos.AccessGatewayID{Id: gatewayId}
	logicalId, err := magmad.RegisterGateway(testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &hwId, Name: "bla"})
	if err != nil || logicalId == "" {
		t.Fatalf("Magmad Register Error: %s, logical ID: %#v", err, logicalId)
	}

	name := strconv.Itoa(int(MetricName))
	key := strconv.Itoa(int(LabelName))
	value := LabelValue
	float := 1.0
	int_val := uint64(1)
	counter_type := dto.MetricType_COUNTER
	counters := protos.MetricsContainer{
		GatewayId: gatewayId,
		Family: []*dto.MetricFamily{{
			Type: &counter_type,
			Name: &name,
			Metric: []*dto.Metric{
				{
					Label:   []*dto.LabelPair{{Name: &key, Value: &value}},
					Counter: &dto.Counter{Value: &float}}}}}}

	gauge_type := dto.MetricType_GAUGE
	gauges := protos.MetricsContainer{
		GatewayId: gatewayId,
		Family: []*dto.MetricFamily{{
			Type: &gauge_type,
			Name: &name,
			Metric: []*dto.Metric{
				{
					Label: []*dto.LabelPair{{Name: &key, Value: &value}},
					Gauge: &dto.Gauge{Value: &float}}}}}}

	summary_type := dto.MetricType_SUMMARY
	summaries := protos.MetricsContainer{
		GatewayId: gatewayId,
		Family: []*dto.MetricFamily{{
			Type: &summary_type,
			Name: &name,
			Metric: []*dto.Metric{
				{
					Label:   []*dto.LabelPair{{Name: &key, Value: &value}},
					Summary: &dto.Summary{SampleSum: &float, SampleCount: &int_val}}}}}}

	histogram_type := dto.MetricType_HISTOGRAM
	histograms := protos.MetricsContainer{
		GatewayId: gatewayId,
		Family: []*dto.MetricFamily{{
			Type: &histogram_type,
			Name: &name,
			Metric: []*dto.Metric{
				{
					Label: []*dto.LabelPair{{Name: &key, Value: &value}},
					Histogram: &dto.Histogram{SampleSum: &float, SampleCount: &int_val,
						Bucket: []*dto.Bucket{
							{CumulativeCount: &int_val,
								UpperBound: &float}}}}}}}}

	// Collect counters
	_, err = srv.Collect(ctx, &counters)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.Queue))
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.Queue[0].Value())
	// clear queue
	e.Queue = e.Queue[:0]

	// Collect gauges
	_, err = srv.Collect(ctx, &gauges)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.Queue))
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.Queue[0].Value())
	e.Queue = e.Queue[:0]

	// Collect summaries
	_, err = srv.Collect(ctx, &summaries)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(e.Queue))
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.Queue[0].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.Queue[0].Value())
	e.Queue = e.Queue[:0]

	// Collect histograms
	_, err = srv.Collect(ctx, &histograms)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(e.Queue))
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.Queue[0].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'E', -1, 64), e.Queue[1].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'E', -1, 64), e.Queue[2].Value())
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.Queue[3].Value())
	e.Queue = e.Queue[:0]

	// Test Collect can raise Write errors
	_, err = srv.Collect(ctx, &protos.MetricsContainer{})
	assert.Error(t, err)
}

func TestConsume(t *testing.T) {
	metricsChan := make(chan *dto.MetricFamily)
	e := NewTestMetricExporter()

	srv := servicers.NewMetricsControllerServer()
	srv.RegisterExporter(e)

	go srv.ConsumeCloudMetrics(metricsChan, "Host_name_place_holder")
	fam1 := "test1"
	fam2 := "test2"
	go func() {
		metricsChan <- &dto.MetricFamily{Name: &fam1, Metric: []*dto.Metric{{}}}
		metricsChan <- &dto.MetricFamily{Name: &fam2, Metric: []*dto.Metric{{}}}
	}()
	time.Sleep(time.Second)
	assert.Equal(t, 2, len(e.Queue))
}
