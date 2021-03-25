package calculations_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_calculations "magma/lte/cloud/go/services/lte/analytics/calculations"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/lib/go/metrics"

	"github.com/stretchr/testify/assert"
)

func TestUserCalculations(t *testing.T) {
	configurator_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(
		"n0",
		configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g0", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw0"},
		serdes.Entity,
	)
	assert.NoError(t, err)

	ctx := test_utils.GetContextWithCertificate(t, "hw0")

	subState0 := state.ArbitraryJSON{
		"oai.ipv4": []map[string]interface{}{
			{
				"apn":             "oai.ipv4",
				"lifecycle_state": "SESSION_ACTIVE",
			},
		},
	}
	subState1 := state.ArbitraryJSON{
		"oai.ipv4": []map[string]interface{}{
			{
				"apn":             "oai.ipv4",
				"lifecycle_state": "SESSION_TERMINATED",
			},
		},
	}
	test_utils.ReportState(t, ctx, lte.SubscriberStateType, "IMSI1234567890", &subState0, serdes.State)
	test_utils.ReportState(t, ctx, lte.SubscriberStateType, "IMSI0987654321", &subState1, serdes.State)

	analyticsConfig := &calculations.AnalyticsConfig{
		Metrics: map[string]calculations.MetricConfig{
			metrics.ConfiguredSubscribersMetric: {
				Export:   true,
				Register: true,
			},
			metrics.ActiveSessionAPNMetric: {
				Export:   true,
				Register: true,
			},
			metrics.ActualSubscribersMetric: {
				Export:   true,
				Register: true,
			},
		},
	}
	userMetricsCalculation := lte_calculations.UserMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{
				AnalyticsConfig: analyticsConfig,
			},
		},
	}
	results, err := userMetricsCalculation.Calculate(nil)
	assert.NoError(t, err)
	assert.Equal(t, len(results), 3)
	resultMetricMap := make(map[string]float64)
	for _, result := range results {
		resultMetricMap[result.GetMetricName()] = result.GetValue()
	}
	assert.Equal(t, resultMetricMap[metrics.ConfiguredSubscribersMetric], float64(0))
	assert.Equal(t, resultMetricMap[metrics.ActiveSessionAPNMetric], float64(1))
	assert.Equal(t, resultMetricMap[metrics.ActualSubscribersMetric], float64(2))
}

func TestSiteCalculations(t *testing.T) {
	configurator_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(
		"n0",
		configurator.NetworkEntity{
			Type:       lte.CellularGatewayEntityType,
			Key:        "g0",
			Config:     &models.MagmadGatewayConfigs{},
			PhysicalID: "hw0"},
		serdes.Entity,
	)
	assert.NoError(t, err)

	ctx := test_utils.GetContextWithCertificate(t, "hw0")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw0"))
	analyticsConfig := &calculations.AnalyticsConfig{
		Metrics: map[string]calculations.MetricConfig{
			metrics.EnodebConnectedMetric: {
				Export:   true,
				Register: true,
			},
			metrics.GatewayMagmaVersionMetric: {
				Export:   true,
				Register: true,
			},
		},
	}
	siteMetricsCalculation := lte_calculations.SiteMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{
				AnalyticsConfig: analyticsConfig,
			},
		},
	}
	results, err := siteMetricsCalculation.Calculate(nil)
	assert.NoError(t, err)
	t.Log(results)
	resultMetricMap := make(map[string]string)
	for _, result := range results {
		resultMetricMap[result.GetLabels()[metrics.GatewayLabelName]] = result.GetLabels()[metrics.GatewayMagmaVersionLabel]
	}
	assert.Equal(t, resultMetricMap["hw0"], "0.0.0.0")
}
