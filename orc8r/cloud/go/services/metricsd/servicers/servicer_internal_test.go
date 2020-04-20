package servicers

import (
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"

	tests "magma/orc8r/cloud/go/services/metricsd/test_common"

	prometheusProto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	testLabels = []*prometheusProto.LabelPair{{Name: tests.MakeStrPtr("labelName"), Value: tests.MakeStrPtr("labelValue")}}
)

func TestPreprocessCloudMetrics(t *testing.T) {
	testFamily := tests.MakeTestMetricFamily(prometheusProto.MetricType_GAUGE, 1, testLabels)
	metricAndContext := preprocessCloudMetrics(testFamily, "hostA")

	assert.NotNil(t, metricAndContext.Context.AdditionalContext)
	assert.Equal(t, "hostA", metricAndContext.Context.AdditionalContext.(*exporters.CloudMetricContext).CloudHost)

	labels := metricAndContext.Family.GetMetric()[0].Label
	assert.Equal(t, 2, len(labels))
	assert.True(t, tests.HasLabel(labels, "cloudHost", "hostA"))
	assert.True(t, tests.HasLabel(labels, testLabels[0].GetName(), testLabels[0].GetValue()))
}

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
	assert.True(t, tests.HasLabel(labels, metrics.NetworkLabelName, "testNetwork"))
	assert.True(t, tests.HasLabel(labels, metrics.GatewayLabelName, "gw1"))
	assert.True(t, tests.HasLabel(labels, testLabels[0].GetName(), testLabels[0].GetValue()))
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
