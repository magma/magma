package calculations_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/lte/analytics/calculations"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestUserCalculations(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(
		"n0",
		configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g0", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw0"},
		serdes.Entity,
	)

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

	userMetricsCalculation := calculations.UserMetricsCalculation{}
	results, err := userMetricsCalculation.Calculate(nil)
	assert.NoError(t, err)
	assert.Equal(t, len(results), 3)
	resultMetricMap := make(map[string]float64)
	for _, result := range results {
		resultMetricMap[result.GetMetricName()] = result.GetValue()
	}
	assert.Equal(t, resultMetricMap[calculations.ConfiguredSubscribersMetric], float64(0))
	assert.Equal(t, resultMetricMap[calculations.ActiveSessionAPNMetric], float64(1))
	assert.Equal(t, resultMetricMap[calculations.ActualSubscribersMetric], float64(2))
}

func TestGeneralCalculations(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	configurator.CreateNetwork(configurator.Network{ID: "n0_1", Type: "LTE"}, serdes.Network)
	configurator.CreateNetwork(configurator.Network{ID: "n0_2", Type: "LTE"}, serdes.Network)
	configurator.CreateNetwork(configurator.Network{ID: "n1", Type: "FEG_LTE"}, serdes.Network)
	configurator.CreateNetwork(configurator.Network{ID: "n2_0", Type: "FEG"}, serdes.Network)
	configurator.CreateNetwork(configurator.Network{ID: "n2_1", Type: "FEG"}, serdes.Network)
	configurator.CreateNetwork(configurator.Network{ID: "n2_2", Type: "FEG"}, serdes.Network)
	generalCalculation := calculations.GeneralMetricsCalculation{}
	results, err := generalCalculation.Calculate(nil)
	resultMetricMap := make(map[string]float64)
	for _, result := range results {
		resultMetricMap[result.GetLabels()["networkType"]] = result.GetValue()
	}
	assert.NoError(t, err)
	assert.Equal(t, resultMetricMap["LTE"], float64(2))
	assert.Equal(t, resultMetricMap["FEG_LTE"], float64(1))
	assert.Equal(t, resultMetricMap["FEG"], float64(3))
}

func TestSiteCalculations(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(
		"n0",
		configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g0", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw0"},
		serdes.Entity,
	)

	ctx := test_utils.GetContextWithCertificate(t, "hw0")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw0"))
	siteMetricsCalculation := calculations.SiteMetricsCalculation{}
	results, err := siteMetricsCalculation.Calculate(nil)
	assert.NoError(t, err)
	resultMetricMap := make(map[string]float64)
	for _, result := range results {
		resultMetricMap[result.GetMetricName()] = result.GetValue()
	}
	assert.Equal(t, resultMetricMap[calculations.GatewayMagmaVersionMetric], float64(1))
}
