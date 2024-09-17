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

package protected

import (
	"context"
	"flag"
	"strconv"
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
	prometheusProto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/cloud/go/services/metricsd/test_init"
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

func TestPreprocessCloudMetrics(t *testing.T) {
	testFamily := tests.MakeTestMetricFamily(prometheusProto.MetricType_GAUGE, 1, testLabels)
	metricAndContext := preprocessCloudMetrics(testFamily, "hostA")

	assert.NotNil(t, metricAndContext.Context.AdditionalContext)
	assert.Equal(t, "hostA", metricAndContext.Context.AdditionalContext.(*exporters.CloudMetricContext).CloudHost)

	labels := metricAndContext.Family.GetMetric()[0].Label
	assert.Equal(t, 2, len(labels))
	assert.True(t, tests.HasLabel(labels, "cloudHost", "hostA"))
	assert.True(t, tests.HasLabel(labels, testLabels[0].GetName(), testLabels[0].GetValue()))
	assert.Equal(t, testFamily.GetName(), metricAndContext.Context.MetricName)
}

func TestPushedMetricsToMetricsAndContext(t *testing.T) {
	container := protos.PushedMetricsContainer{
		NetworkId: "testNetwork",
		Metrics: []*protos.PushedMetric{{
			MetricName:  "metricA",
			Value:       10,
			TimestampMS: 1234,
			Labels:      []*protos.LabelPair{{Name: "labelName", Value: "labelValue"}},
		},
		},
	}

	metricAndContext := pushedMetricsToMetricsAndContext(&container)

	assert.Equal(t, 1, len(metricAndContext))
	ctx := metricAndContext[0].Context
	family := metricAndContext[0].Family
	assert.NotNil(t, ctx.AdditionalContext)
	assert.Equal(t, "testNetwork", ctx.AdditionalContext.(*exporters.PushedMetricContext).NetworkID)

	labels := family.GetMetric()[0].Label
	assert.Equal(t, 2, len(labels))
	assert.True(t, tests.HasLabel(labels, metrics.NetworkLabelName, "testNetwork"))
	assert.True(t, tests.HasLabel(labels, "labelName", "labelValue"))
}

func TestConsume(t *testing.T) {
	metricsChan := make(chan *dto.MetricFamily)
	e := &testMetricExporter{}

	test_init.StartNewTestExporter(t, e)
	srv := NewCloudMetricsControllerServer()

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
	srv := NewCloudMetricsControllerServer()

	// Create test network
	networkID := "metricsd_servicer_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network Name")

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

func TestPushRaw(t *testing.T) {
	device_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	e := &testMetricExporter{}
	test_init.StartNewTestExporter(t, e)
	srv := NewCloudMetricsControllerServer()

	fam := collection.MakeSingleGaugeFamily("name1", "help1", nil, 12.34)

	c := &protos.RawMetricsContainer{
		HostName: "someHostName",
		Families: []*prometheusProto.MetricFamily{fam},
	}

	_, err := srv.PushRaw(context.Background(), c)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.queue))
	assert.Equal(t, "name1", e.queue[0].Name())
	assert.Equal(t, 2, len(e.queue[0].Labels()))
	assert.Equal(t, "cloudHost", *e.queue[0].Labels()[0].Name)
	assert.Equal(t, "someHostName", *e.queue[0].Labels()[0].Value)
}
