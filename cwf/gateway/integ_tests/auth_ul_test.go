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
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install for a pass-all dynamic rule and 250KB of
//   quota.
//   Trigger a authentication and assert the CCR-I is received.
// - Generate traffic to put traffic through the newly installed rule.
func TestAuthenticateUplinkTraffic(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateUplinkTraffic...")
	tr := NewTestRunner(t)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, tr.CleanUp())
	}()

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

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "100K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	resultByIndex, errByIndex, err := getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)

	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
}
