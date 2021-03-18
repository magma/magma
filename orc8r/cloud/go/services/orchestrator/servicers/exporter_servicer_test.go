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
	"regexp"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/cloud/go/services/metricsd/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/servicers"
	"magma/orc8r/lib/go/metrics"

	prometheus_models "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

const (
	sampleNetworkID  = "sampleNetwork"
	sampleGatewayID  = "sampleGateway"
	sampleMetricName = "metric_A"

	// Below copied from gaugeconverter.go
	bucketPostfix            = "_bucket"
	histogramBucketLabelName = "le"
	summaryQuantileLabelName = "quantile"
)

var (
	sampleLabels = []*prometheus_models.LabelPair{
		{Name: tests.MakeStrPtr(metrics.NetworkLabelName), Value: tests.MakeStrPtr(sampleNetworkID)},
		{Name: tests.MakeStrPtr("testLabel"), Value: tests.MakeStrPtr("testValue")},
	}

	sampleGatewayContext = exporters.MetricContext{
		MetricName: sampleMetricName,
		AdditionalContext: &exporters.GatewayMetricContext{
			NetworkID: sampleNetworkID,
			GatewayID: sampleGatewayID,
		},
	}
)

func TestEnsureHTTP(t *testing.T) {
	addrs := []string{"http://prometheus-cache:9091", "prometheus-cache:9091", "https://prometheus-cache:9091"}
	srv := servicers.NewPushExporterServicer(addrs).(*servicers.PushExporterServicer)
	protocolMatch := regexp.MustCompile("(http|https)://")
	for _, addr := range srv.PushAddresses {
		assert.True(t, protocolMatch.MatchString(addr))
	}
}

func TestPushExporterServicer_Submit_Gauge(t *testing.T) {
	_, exp := makeTestCustomPushExporter(t)
	err := submitNewMetric(exp, prometheus_models.MetricType_GAUGE, sampleGatewayContext)
	assert.NoError(t, err)

	err = submitNewMetric(exp, prometheus_models.MetricType_GAUGE, sampleGatewayContext)
	assert.NoError(t, err)
}

func TestPushExporterServicer_Submit_Counter(t *testing.T) {
	_, exp := makeTestCustomPushExporter(t)
	err := submitNewMetric(exp, prometheus_models.MetricType_COUNTER, sampleGatewayContext)
	assert.NoError(t, err)

	err = submitNewMetric(exp, prometheus_models.MetricType_COUNTER, sampleGatewayContext)
	assert.NoError(t, err)
}

func TestPushExporterServicer_Submit_Histogram(t *testing.T) {
	_, exp := makeTestCustomPushExporter(t)
	err := submitNewMetric(exp, prometheus_models.MetricType_HISTOGRAM, sampleGatewayContext)
	assert.NoError(t, err)

	err = submitNewMetric(exp, prometheus_models.MetricType_HISTOGRAM, sampleGatewayContext)
	assert.NoError(t, err)
}

func TestPushExporterServicer_Submit_Summary(t *testing.T) {
	_, exp := makeTestCustomPushExporter(t)
	err := submitNewMetric(exp, prometheus_models.MetricType_SUMMARY, sampleGatewayContext)
	assert.NoError(t, err)

	err = submitNewMetric(exp, prometheus_models.MetricType_SUMMARY, sampleGatewayContext)
	assert.NoError(t, err)
}

func TestPushExporterServicer_Submit_Untyped(t *testing.T) {
	_, exp := makeTestCustomPushExporter(t)
	err := submitNewMetric(exp, prometheus_models.MetricType_UNTYPED, sampleGatewayContext)
	assert.NoError(t, err)

	err = submitNewMetric(exp, prometheus_models.MetricType_UNTYPED, sampleGatewayContext)
	assert.NoError(t, err)
}

func TestPushExporterServicer_Submit_InvalidMetrics(t *testing.T) {
	// Submitting a metric family with 0 metrics should not register the family
	_, exp := makeTestCustomPushExporter(t)
	noMetricFamily := tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 0, sampleLabels)
	mc := exporters.MetricAndContext{
		Family:  noMetricFamily,
		Context: sampleGatewayContext,
	}

	err := exp.Submit([]exporters.MetricAndContext{mc})
	assert.NoError(t, err)
}

func TestPushExporterServicer_Submit_InvalidName(t *testing.T) {
	// Submitting a metric with an invalid name should submit a renamed metric
	testInvalidName(t, "invalid metric name", "invalid_metric_name")
	testInvalidName(t, "0starts_with_number", "_0starts_with_number")
	testInvalidName(t, "bad?-/$chars", "bad____chars")
}

func testInvalidName(t *testing.T, inputName, expectedName string) {
	_, exp := makeTestCustomPushExporter(t)
	mf := tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 1, sampleLabels)

	mc := exporters.MetricAndContext{
		Family: mf,
		Context: exporters.MetricContext{
			MetricName: inputName,
		},
	}

	err := exp.Submit([]exporters.MetricAndContext{mc})
	assert.NoError(t, err)
}

func TestPushExporterServicer_Submit_InvalidLabel(t *testing.T) {
	// Submitting a metric with invalid labelnames should not include that metric
	_, exp := makeTestCustomPushExporter(t)
	mf := tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 5, sampleLabels)
	extraMetric := tests.MakePromoGauge(10)
	mf.Metric[2] = &extraMetric
	mf.Metric[2].Label = append(mf.Metric[2].Label, &prometheus_models.LabelPair{Name: tests.MakeStrPtr("1"), Value: tests.MakeStrPtr("badLabelName")})

	mc := exporters.MetricAndContext{
		Family:  mf,
		Context: sampleGatewayContext,
	}

	err := exp.Submit([]exporters.MetricAndContext{mc})
	assert.NoError(t, err)

	// If all metrics are invalid, the family should not be submitted
	_, exp = makeTestCustomPushExporter(t)
	mf = tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 1, sampleLabels)
	badMetric := tests.MakePromoGauge(10)
	mf.Metric[0] = &badMetric
	mf.Metric[0].Label = append(mf.Metric[0].Label, &prometheus_models.LabelPair{Name: tests.MakeStrPtr("1"), Value: tests.MakeStrPtr("badLabelName")})

	mc = exporters.MetricAndContext{
		Family:  mf,
		Context: sampleGatewayContext,
	}

	err = exp.Submit([]exporters.MetricAndContext{mc})
	assert.NoError(t, err)
}

func submitNewMetric(exp exporters.Exporter, mtype prometheus_models.MetricType, ctx exporters.MetricContext) error {
	mc := exporters.MetricAndContext{
		Family:  tests.MakeTestMetricFamily(mtype, 1, sampleLabels),
		Context: ctx,
	}
	return exp.Submit([]exporters.MetricAndContext{mc})
}

func makeTestCustomPushExporter(t *testing.T) (*servicers.PushExporterServicer, exporters.Exporter) {
	srv := servicers.NewPushExporterServicer([]string{})
	test_init.StartTestServiceInternal(t, srv)

	exporterSrv := srv.(*servicers.PushExporterServicer)
	exporter := exporters.NewRemoteExporter(metricsd.ServiceName)
	return exporterSrv, exporter
}
