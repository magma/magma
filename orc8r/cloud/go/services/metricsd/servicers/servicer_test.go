/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers_test

import (
	"errors"
	"flag"
	"strconv"
	"testing"
	"time"

	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/servicers"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/cloud/go/services/metricsd/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type testMetricExporter struct {
	queue []exporters.Sample

	// error to return
	retErr error
}

const (
	MetricName = protos.MetricName_process_virtual_memory_bytes
	LabelName  = protos.MetricLabelName_result
	LabelValue = "success"
)

// Set verbosity so we can capture exporter error logging
var _ = flag.Set("vmodule", "*=2")

func (e *testMetricExporter) Submit(metrics []exporters.MetricAndContext) error {
	for _, metricAndContext := range metrics {
		family := metricAndContext.Family
		for _, metric := range family.GetMetric() {
			convertedMetricAndContext := exporters.MakeProtoMetric(metricAndContext)
			e.queue = append(
				e.queue,
				exporters.GetSamplesForMetrics(convertedMetricAndContext, metric)...,
			)
		}
	}

	return e.retErr
}

func (e *testMetricExporter) Start() {}

func TestCollect(t *testing.T) {
	device_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	e := &testMetricExporter{}
	test_init.StartNewTestExporter(t, e)
	srv := servicers.NewMetricsControllerServer()

	// Create test network
	networkID := "metricsd_servicer_test_network"
	test_utils.RegisterNetwork(t, networkID, "Test Network Name")

	// Register a fake gateway
	gatewayID := "2876171d-bf38-4254-b4da-71a713952904"
	id := protos.NewGatewayIdentity(gatewayID, "testNwId", "testLogicalId")
	ctx := id.NewContextWithIdentity(context.Background())
	test_utils.RegisterGateway(t, networkID, gatewayID, &models.GatewayDevice{HardwareID: gatewayID})

	name := strconv.Itoa(int(MetricName))
	key := strconv.Itoa(int(LabelName))
	value := LabelValue
	float := 1.0
	int_val := uint64(1)
	counter_type := dto.MetricType_COUNTER
	counters := protos.MetricsContainer{
		GatewayId: gatewayID,
		Family: []*dto.MetricFamily{{
			Type: &counter_type,
			Name: &name,
			Metric: []*dto.Metric{
				{
					Label:   []*dto.LabelPair{{Name: &key, Value: &value}},
					Counter: &dto.Counter{Value: &float}}}}}}

	gauge_type := dto.MetricType_GAUGE
	gauges := protos.MetricsContainer{
		GatewayId: gatewayID,
		Family: []*dto.MetricFamily{{
			Type: &gauge_type,
			Name: &name,
			Metric: []*dto.Metric{
				{
					Label: []*dto.LabelPair{{Name: &key, Value: &value}},
					Gauge: &dto.Gauge{Value: &float}}}}}}

	summary_type := dto.MetricType_SUMMARY
	summaries := protos.MetricsContainer{
		GatewayId: gatewayID,
		Family: []*dto.MetricFamily{{
			Type: &summary_type,
			Name: &name,
			Metric: []*dto.Metric{
				{
					Label:   []*dto.LabelPair{{Name: &key, Value: &value}},
					Summary: &dto.Summary{SampleSum: &float, SampleCount: &int_val}}}}}}

	histogram_type := dto.MetricType_HISTOGRAM
	histograms := protos.MetricsContainer{
		GatewayId: gatewayID,
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
	_, err := srv.Collect(ctx, &counters)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.queue))
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.queue[0].Value())
	// check that label protos are converted
	assert.True(t, tests.HasLabelName(e.queue[0].Labels(), protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name)))
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	// clear queue
	e.queue = e.queue[:0]

	// Collect gauges
	_, err = srv.Collect(ctx, &gauges)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.queue))
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.queue[0].Value())
	assert.True(t, tests.HasLabelName(e.queue[0].Labels(), protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name)))
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	//assert.Equal(t, protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name), e.queue[0].Labels()[0].GetName())
	e.queue = e.queue[:0]

	// Collect summaries
	_, err = srv.Collect(ctx, &summaries)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(e.queue))
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.queue[0].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.queue[0].Value())
	assert.True(t, tests.HasLabelName(e.queue[0].Labels(), protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name)))
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	//assert.Equal(t, protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name), e.queue[0].Labels()[0].GetName())
	e.queue = e.queue[:0]

	// Collect histograms
	_, err = srv.Collect(ctx, &histograms)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(e.queue))
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.queue[0].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'E', -1, 64), e.queue[1].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'E', -1, 64), e.queue[2].Value())
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.queue[3].Value())
	assert.True(t, tests.HasLabelName(e.queue[0].Labels(), protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name)))
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	//assert.Equal(t, protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name), e.queue[0].Labels()[0].GetName())
	e.queue = e.queue[:0]

	// Test Collect with empty collection
	_, err = srv.Collect(ctx, &protos.MetricsContainer{})
	assert.NoError(t, err)

	// Exporter error should not result in error returned from Collect
	// But with verbosity set to 2 at the top of the test, we should log
	prevInfoLines := glog.Stats.Info.Lines()
	e.retErr = errors.New("mock exporter error")
	_, err = srv.Collect(ctx, &gauges)
	assert.NoError(t, err)
	assert.Equal(t, prevInfoLines+1, glog.Stats.Info.Lines())
}

func TestCollectMismatchedGateway(t *testing.T) {
	device_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	e := &testMetricExporter{}
	test_init.StartNewTestExporter(t, e)
	srv := servicers.NewMetricsControllerServer()

	// Create test network
	networkID := "metricsd_servicer_test_network"
	test_utils.RegisterNetwork(t, networkID, "Test Network Name")

	// Register a fake gateway
	gatewayID := "2876171d-bf38-4254-b4da-71a713952904"
	id := protos.NewGatewayIdentity(gatewayID, "testNwId", "testLogicalId")
	ctx := id.NewContextWithIdentity(context.Background())
	test_utils.RegisterGateway(t, networkID, gatewayID, &models.GatewayDevice{HardwareID: gatewayID})

	// Mismatched gateway ID
	mismatchgatewayID := "2876171d-bf38-4254-b4da-71a713954029"

	name := strconv.Itoa(int(MetricName))
	key := strconv.Itoa(int(LabelName))
	value := LabelValue
	float := 1.0
	counter_type := dto.MetricType_COUNTER
	counters := protos.MetricsContainer{
		GatewayId: mismatchgatewayID,
		Family: []*dto.MetricFamily{
			{
				Type: &counter_type,
				Name: &name,
				Metric: []*dto.Metric{
					{
						Label:   []*dto.LabelPair{{Name: &key, Value: &value}},
						Counter: &dto.Counter{Value: &float},
					},
				},
			},
		},
	}

	// Collect counters
	_, err := srv.Collect(ctx, &counters)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.queue))
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.queue[0].Value())
	// check that label protos are converted
	assert.True(t, tests.HasLabelName(e.queue[0].Labels(), protos.GetEnumNameIfPossible(key, protos.MetricLabelName_name)))
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.GatewayLabelName, gatewayID))

}

func TestConsume(t *testing.T) {
	metricsChan := make(chan *dto.MetricFamily)
	e := &testMetricExporter{}

	test_init.StartNewTestExporter(t, e)
	srv := servicers.NewMetricsControllerServer()

	go srv.ConsumeCloudMetrics(metricsChan, "Host_name_place_holder")
	fam1 := "test1"
	fam2 := "test2"
	go func() {
		metricsChan <- &dto.MetricFamily{Name: &fam1, Metric: []*dto.Metric{{}}}
		metricsChan <- &dto.MetricFamily{Name: &fam2, Metric: []*dto.Metric{{}}}
	}()
	time.Sleep(time.Second)
	assert.Equal(t, 2, len(e.queue))
}

func TestPush(t *testing.T) {
	device_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	e := &testMetricExporter{}
	test_init.StartNewTestExporter(t, e)
	srv := servicers.NewMetricsControllerServer()

	// Create test network
	networkID := "metricsd_servicer_test_network"
	test_utils.RegisterNetwork(t, networkID, "Test Network Name")

	metricName := "test_metric"
	value := 8.2
	testLabel := &protos.LabelPair{Name: "labelName", Value: "labelValue"}
	timestamp := int64(123456)

	protoMet := protos.PushedMetric{
		MetricName:  metricName,
		Value:       value,
		TimestampMS: timestamp,
		Labels:      []*protos.LabelPair{testLabel},
	}
	pushedMetrics := protos.PushedMetricsContainer{
		NetworkId: networkID,
		Metrics:   []*protos.PushedMetric{&protoMet},
	}

	_, err := srv.Push(context.Background(), &pushedMetrics)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.queue))
	assert.Equal(t, metricName, e.queue[0].Name())
	assert.Equal(t, 2, len(e.queue[0].Labels()))
	assert.Equal(t, testLabel.Name, *e.queue[0].Labels()[0].Name)
	assert.Equal(t, testLabel.Value, *e.queue[0].Labels()[0].Value)
	assert.Equal(t, metrics.NetworkLabelName, *e.queue[0].Labels()[1].Name)
	assert.Equal(t, networkID, *e.queue[0].Labels()[1].Value)
	assert.Equal(t, timestamp, e.queue[0].TimestampMs())
	assert.Equal(t, strconv.FormatFloat(value, 'f', -1, 64), e.queue[0].Value())
}
