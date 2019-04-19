/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tc "magma/orc8r/cloud/go/services/metricsd/test_common"

	dto "github.com/prometheus/client_model/go"

	"github.com/stretchr/testify/assert"
)

const (
	testCloudNetwork = "mesh_tobias_dogfooding"
	testCloudGateway = "idfaceb00cfaceb00cface6031973e8372"

	testRegularName = "regular_metric_name"
	testNoGatewayID = "no_gwID_networkId=test_network"
)

var (
	testCloudName = fmt.Sprintf("gateway_checkin_status_gatewayId_%s_networkId_%s", testCloudGateway, testCloudNetwork)
)

func TestParseCloudMetricName(t *testing.T) {
	networkID, gatewayID := unpackCloudMetricName(testCloudName)
	assert.Equal(t, testCloudNetwork, networkID)
	assert.Equal(t, testCloudGateway, gatewayID)

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

type networkAndGatewayTestCase struct {
	metricLabels    []*dto.LabelPair
	familyName      string
	expectedNetwork string
	expectedGateway string
}

func (c networkAndGatewayTestCase) runTest(t *testing.T) {
	family := tc.MakeTestMetricFamily(dto.MetricType_GAUGE, c.familyName, 1, c.metricLabels)
	networkID, gatewayID := determineCloudNetworkAndGatewayID(family, tc.DefaultGateway)
	assert.Equal(t, c.expectedNetwork, networkID)
	assert.Equal(t, c.expectedGateway, gatewayID)
}

func makeNetworkAndGatewayTestCase(labels []*dto.LabelPair, name, expectedNetwork, expectedGateway string) networkAndGatewayTestCase {
	return networkAndGatewayTestCase{
		metricLabels:    labels,
		familyName:      name,
		expectedNetwork: expectedNetwork,
		expectedGateway: expectedGateway,
	}
}

func TestDetermineCloudNetworkAndGatewayID(t *testing.T) {
	testCases := []networkAndGatewayTestCase{
		// Use labels when none in metric name
		makeNetworkAndGatewayTestCase(tc.SimpleLabels, testRegularName, exporters.CloudMetricID, tc.DefaultGateway),
		makeNetworkAndGatewayTestCase(tc.GatewayLabels, testRegularName, exporters.CloudMetricID, tc.TestGateway),
		makeNetworkAndGatewayTestCase(tc.NetworkLabels, testRegularName, tc.TestNetwork, tc.DefaultGateway),
		makeNetworkAndGatewayTestCase(tc.NetworkAndGatewayLabels, testRegularName, tc.TestNetwork, tc.TestGateway),

		// Use metric name instead of labels
		makeNetworkAndGatewayTestCase(tc.SimpleLabels, testCloudName, testCloudNetwork, testCloudGateway),
		makeNetworkAndGatewayTestCase(tc.GatewayLabels, testCloudName, testCloudNetwork, testCloudGateway),
		makeNetworkAndGatewayTestCase(tc.NetworkLabels, testCloudName, testCloudNetwork, testCloudGateway),
		makeNetworkAndGatewayTestCase(tc.NetworkAndGatewayLabels, testCloudName, testCloudNetwork, testCloudGateway),
	}
	for _, c := range testCases {
		c.runTest(t)
	}
}
