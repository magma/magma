package exporters

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCloudMetricName(t *testing.T) {
	metricName := "gateway_checkin_status_gatewayId_idfaceb00cfaceb00cface6031973e8372_networkId_mesh_tobias_dogfooding"
	networkID, gatewayID, err := UnpackCloudMetricName(metricName)
	assert.NoError(t, err)
	assert.Equal(t, "mesh_tobias_dogfooding", networkID)
	assert.Equal(t, "idfaceb00cfaceb00cface6031973e8372", gatewayID)

	metricName = "regular_metric_name"
	networkID, gatewayID, err = UnpackCloudMetricName(metricName)
	expectedErr := fmt.Sprintf("no gateway or network label in cloud metric: %s", metricName)
	assert.EqualError(t, err, expectedErr)
}
