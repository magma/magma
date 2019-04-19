/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tc "magma/orc8r/cloud/go/services/metricsd/test_common"

	dto "github.com/prometheus/client_model/go"

	"github.com/stretchr/testify/assert"
)

const (
	testFamilyName = "testFamily"
	testNetwork    = "test_network"
	testGateway    = "test_gateway"
)

func TestMakeGraphiteName(t *testing.T) {
	expectedName := fmt.Sprintf("%s;%s=%s;%s=%s", testFamilyName, GatewayTagName, defaultGateway, NetworkTagName, defaultNetwork)
	testMakeGraphiteNameHelper(t, "", "", testFamilyName, []*dto.LabelPair{}, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s;%s=%s", testFamilyName, GatewayTagName, testGateway, tc.SimpleLabelName, tc.SimpleLabelValue, NetworkTagName, testNetwork)
	testMakeGraphiteNameHelper(t, testNetwork, testGateway, testFamilyName, tc.SimpleLabels, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s;%s=%s", testFamilyName, GatewayTagName, defaultGateway, tc.SimpleLabelName, tc.SimpleLabelValue, NetworkTagName, testNetwork)
	testMakeGraphiteNameHelper(t, testNetwork, "", testFamilyName, tc.SimpleLabels, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s;%s=%s", testFamilyName, GatewayTagName, defaultGateway, tc.SimpleLabelName, tc.SimpleLabelValue, NetworkTagName, defaultNetwork)
	testMakeGraphiteNameHelper(t, "", "", testFamilyName, tc.SimpleLabels, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s", testFamilyName, GatewayTagName, tc.TestGateway, NetworkTagName, tc.TestNetwork)
	testMakeGraphiteNameHelper(t, tc.TestNetwork, tc.TestGateway, testFamilyName, tc.NetworkLabels, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s", testFamilyName, GatewayTagName, defaultGateway, NetworkTagName, defaultNetwork)
	testMakeGraphiteNameHelper(t, "", "", testFamilyName, tc.GatewayLabels, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s", testFamilyName, GatewayTagName, tc.TestGateway, NetworkTagName, tc.TestNetwork)
	testMakeGraphiteNameHelper(t, tc.TestNetwork, tc.TestGateway, testFamilyName, tc.NetworkAndGatewayLabels, expectedName)
}

func testMakeGraphiteNameHelper(t *testing.T, networkID, gatewayID, familyName string, labels []*dto.LabelPair, expectedName string) {
	metric := tc.MakePromoGauge(100)
	metric.Label = labels
	family := tc.MakeTestMetricFamily(dto.MetricType_GAUGE, familyName, 1, []*dto.LabelPair{})
	ctx := exporters.MetricsContext{
		MetricName:  protos.GetDecodedName(family),
		NetworkID:   networkID,
		GatewayID:   gatewayID,
		DecodedName: protos.GetDecodedName(family),
	}

	graphiteName := makeGraphiteName(&metric, ctx)
	assert.Equal(t, expectedName, graphiteName)
}
