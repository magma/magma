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

package southbound

import (
	"context"
	"errors"
	"flag"
	"strconv"
	"testing"

	"github.com/golang/glog"
	dto "github.com/prometheus/client_model/go"
	prometheusProto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/test_common"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/cloud/go/services/metricsd/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"
)

var (
	testLabels = []*prometheusProto.LabelPair{{Name: tests.MakeStrPtr("labelName"), Value: tests.MakeStrPtr("labelValue")}}
)

type testMetricExporter struct {
	queue []exporters.Sample

	// error to return
	retErr error
}

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

// Set verbosity so we can capture exporter error logging
var _ = flag.Set("vmodule", "*=2")

func TestMetricsContainerToMetricAndContexts(t *testing.T) {
	testFamily := tests.MakeTestMetricFamily(prometheusProto.MetricType_GAUGE, 1, testLabels)
	container := protos.MetricsContainer{
		GatewayId: "gw1",
		Family:    []*prometheusProto.MetricFamily{testFamily},
	}

	metricAndContext := metricsContainerToMetricAndContexts(&container, "testNetwork", "gw1")

	assert.Equal(t, 1, len(metricAndContext))
	ctx := metricAndContext[0].Context
	family := metricAndContext[0].Family
	assert.NotNil(t, ctx.AdditionalContext)
	assert.Equal(t, "gw1", ctx.AdditionalContext.(*exporters.GatewayMetricContext).GatewayID)
	assert.Equal(t, "testNetwork", ctx.AdditionalContext.(*exporters.GatewayMetricContext).NetworkID)

	labels := family.GetMetric()[0].Label
	assert.Equal(t, 3, len(labels))
	assert.Equal(t, testFamily.GetName(), ctx.MetricName)
	assert.True(t, tests.HasLabel(labels, metrics.NetworkLabelName, "testNetwork"))
	assert.True(t, tests.HasLabel(labels, metrics.GatewayLabelName, "gw1"))
	assert.True(t, tests.HasLabel(labels, testLabels[0].GetName(), testLabels[0].GetValue()))
}
func TestCollect(t *testing.T) {
	device_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	e := &testMetricExporter{}
	test_init.StartNewTestExporter(t, e)
	srv := NewMetricsControllerServer()

	// Create test network
	networkID := "metricsd_servicer_test_network"
	test_utils.RegisterNetwork(t, networkID, "Test Network Name")

	// Register a fake gateway
	gatewayID := "2876171d-bf38-4254-b4da-71a713952904"
	id := protos.NewGatewayIdentity(gatewayID, "testNwId", "testLogicalId")
	ctx := id.NewContextWithIdentity(context.Background())
	test_utils.RegisterGateway(t, networkID, gatewayID, &models.GatewayDevice{HardwareID: gatewayID})

	name := strconv.Itoa(int(test_common.MetricName))
	key := strconv.Itoa(int(test_common.LabelName))
	value := test_common.LabelValue
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
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	// clear queue
	e.queue = e.queue[:0]

	// Collect gauges
	_, err = srv.Collect(ctx, &gauges)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.queue))
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.queue[0].Value())
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	e.queue = e.queue[:0]

	// Collect summaries
	_, err = srv.Collect(ctx, &summaries)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(e.queue))
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.queue[0].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'f', -1, 64), e.queue[0].Value())
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	e.queue = e.queue[:0]

	// Collect histograms
	_, err = srv.Collect(ctx, &histograms)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(e.queue))
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.queue[0].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'E', -1, 64), e.queue[1].Value())
	assert.Equal(t, strconv.FormatFloat(float, 'E', -1, 64), e.queue[2].Value())
	assert.Equal(t, strconv.FormatUint(int_val, 10), e.queue[3].Value())
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
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
	srv := NewMetricsControllerServer()

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

	name := strconv.Itoa(int(test_common.MetricName))
	key := strconv.Itoa(int(test_common.LabelName))
	value := test_common.LabelValue
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
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.NetworkLabelName, networkID))
	assert.True(t, tests.HasLabel(e.queue[0].Labels(), metrics.GatewayLabelName, gatewayID))

}
