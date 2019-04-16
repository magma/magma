/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testCloudName   = "gateway_checkin_status_gatewayId_idfaceb00cfaceb00cface6031973e8372_networkId_mesh_tobias_dogfooding"
	testRegularName = "regular_metric_name"
	testNoGatewayID = "no_gwID_networkId=test_network"
)

func TestParseCloudMetricName(t *testing.T) {
	networkID, gatewayID := unpackCloudMetricName(testCloudName)
	assert.Equal(t, "mesh_tobias_dogfooding", networkID)
	assert.Equal(t, "idfaceb00cfaceb00cface6031973e8372", gatewayID)

	networkID, gatewayID = unpackCloudMetricName(testRegularName)
	assert.Equal(t, "", networkID)
	assert.Equal(t, "", gatewayID)

	networkID, gatewayID = unpackCloudMetricName(testNoGatewayID)
	assert.Equal(t, "test_network", networkID)
	assert.Equal(t, "", gatewayID)
}

func TestRemoveCloudMetricLabels(t *testing.T) {
	expectedName := "gateway_checkin_status"
	strippedName := removeCloudMetricLabels(testCloudName)
	assert.Equal(t, expectedName, strippedName)

	strippedName = removeCloudMetricLabels(testRegularName)
	assert.Equal(t, testRegularName, strippedName)

	expectedName = "no_gwID"
	strippedName = removeCloudMetricLabels(testNoGatewayID)
	assert.Equal(t, expectedName, strippedName)
}
