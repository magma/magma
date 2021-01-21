package calculations_test

import (
	lte_calculations "magma/lte/cloud/go/services/lte/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/query_api/mocks"
	"magma/orc8r/lib/go/metrics"
	"testing"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserThroughput(t *testing.T) {
	metric1 := model.Metric{}
	metric1["networkID"] = "testNetwork1"
	var vec model.Vector
	vec = []*model.Sample{
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 1, Timestamp: 123},
		{Metric: metric1, Value: 10, Timestamp: 123},
	}
	metric2 := model.Metric{}
	metric2["networkID"] = "testNetwork2"
	vec = append(vec, []*model.Sample{
		{Metric: metric2, Value: 10, Timestamp: 123},
		{Metric: metric2, Value: 10, Timestamp: 123},
		{Metric: metric2, Value: 12, Timestamp: 123},
		{Metric: metric2, Value: 13, Timestamp: 123},
		{Metric: metric2, Value: 15, Timestamp: 123},
	}...)
	successClient := &mocks.PrometheusAPI{}
	successClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(vec, nil, nil)
	analyticsConfig := &calculations.AnalyticsConfig{
		Metrics: map[string]calculations.MetricConfig{
			"user_throughput": {
				Export:   true,
				Register: true,
			},
		},
	}
	calc := lte_calculations.UserThroughputCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{
				Name:            "user_throughput",
				AnalyticsConfig: analyticsConfig,
			},
		},
		Direction: "up",
	}
	results, err := calc.Calculate(successClient)
	assert.NoError(t, err)

	// Verify that for testNetwork1 0.5 = 1 and 0.95  = 10
	// and for testNetwork2 0.5 = 12 and 0.95 = 15
	expTestNetwork1P50 := 1.0
	expTestNetwork1P95 := 10.0
	expTestNetwork2P50 := 12.0
	expTestNetwork2P95 := 15.0

	for _, result := range results {
		labels := result.GetLabels()
		if labels[metrics.NetworkLabelName] == "testNetwork1" && labels[metrics.QuantileLabel] == "0.5" {
			assert.InEpsilon(t, expTestNetwork1P50, result.GetValue(), 1e-8)
		}
		if labels[metrics.NetworkLabelName] == "testNetwork1" && labels[metrics.QuantileLabel] == "0.95" {
			assert.InEpsilon(t, expTestNetwork1P95, result.GetValue(), 1e-8)
		}
		if labels[metrics.NetworkLabelName] == "testNetwork2" && labels[metrics.QuantileLabel] == "0.5" {
			assert.InEpsilon(t, expTestNetwork2P50, result.GetValue(), 1e-8)
		}
		if labels[metrics.NetworkLabelName] == "testNetwork2" && labels[metrics.QuantileLabel] == "0.95" {
			assert.InEpsilon(t, expTestNetwork2P95, result.GetValue(), 1e-8)
		}
	}
}
