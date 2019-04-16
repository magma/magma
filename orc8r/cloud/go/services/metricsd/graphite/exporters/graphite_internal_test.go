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
	"magma/orc8r/cloud/go/services/metricsd/test_common"

	dto "github.com/prometheus/client_model/go"

	"github.com/stretchr/testify/assert"
)

const (
	testFamilyName       = "testFamily"
	testNetwork          = "test_network"
	testGateway          = "test_gateway"
	testMetricLabelName  = "metricLabelName"
	testMetricLabelValue = "metricLabelValue"

	gIDName = "gatewayID"
	nIDName = "networkID"
)

var (
	testLabelPair = &dto.LabelPair{
		Name:  test_common.MakeStringPointer(testMetricLabelName),
		Value: test_common.MakeStringPointer(testMetricLabelValue),
	}
)

func TestMakeGraphiteName(t *testing.T) {
	expectedName := fmt.Sprintf("%s;%s=%s;%s=%s", testFamilyName, gIDName, "defaultGateway", nIDName, "defaultNetwork")
	testMakeGraphiteNameHelper(t, "", "", testFamilyName, []*dto.LabelPair{}, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s;%s=%s", testFamilyName, gIDName, testGateway, testMetricLabelName, testMetricLabelValue, nIDName, testNetwork)
	testMakeGraphiteNameHelper(t, testNetwork, testGateway, testFamilyName, []*dto.LabelPair{testLabelPair}, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s;%s=%s", testFamilyName, gIDName, "defaultGateway", testMetricLabelName, testMetricLabelValue, nIDName, testNetwork)
	testMakeGraphiteNameHelper(t, testNetwork, "", testFamilyName, []*dto.LabelPair{testLabelPair}, expectedName)

	expectedName = fmt.Sprintf("%s;%s=%s;%s=%s;%s=%s", testFamilyName, gIDName, "defaultGateway", testMetricLabelName, testMetricLabelValue, nIDName, "defaultNetwork")
	testMakeGraphiteNameHelper(t, "", "", testFamilyName, []*dto.LabelPair{testLabelPair}, expectedName)
}

func testMakeGraphiteNameHelper(t *testing.T, networkID, gatewayID, familyName string, labels []*dto.LabelPair, expectedName string) {
	metric := test_common.MakePromoGauge(100)
	metric.Label = labels
	family := test_common.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{})
	family.Name = &familyName
	ctx := exporters.MetricsContext{
		MetricName:  protos.GetDecodedName(family),
		NetworkID:   networkID,
		GatewayID:   gatewayID,
		DecodedName: protos.GetDecodedName(family),
	}

	graphiteName := makeGraphiteName(&metric, ctx)
	assert.Equal(t, expectedName, graphiteName)
}
