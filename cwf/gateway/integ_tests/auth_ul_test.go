/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"fbc/lib/go/radius/rfc2869"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUplinkTraffic(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateUplinkTraffic...")
	tr := NewTestRunner()
	assert.NoError(t, usePCRFMockDriver())

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	imsi := ues[0].GetImsi()
	usageMonitorInfo := []*protos.UsageMonitoringInformation{
		{
			MonitoringLevel: protos.MonitoringLevel_RuleLevel,
			MonitoringKey:   []byte("mkey1"),
			Octets:          &protos.Octets{TotalOctets: 250 * KiloBytes},
		},
	}
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL, 1)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetDynamicRuleInstalls([]*protos.RuleDefinition{getPassAllRuleDefinition("dynamic-pass-all", "mkey1", 100)}).
		SetUsageMonitorInfos(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// return success with credit on unexpected requests
	defaultAnswer := protos.NewGxCCAnswer(2001).SetUsageMonitorInfos(usageMonitorInfo)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "100K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	time.Sleep(3 * time.Second)

	resultByIndex, errByIndex, err := getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)
	assert.NoError(t, clearPCRFMockDriver())

	// Clear hss, ocs, and pcrf
	assert.NoError(t, clearPCRFMockDriver())
	assert.NoError(t, tr.CleanUp())
}
