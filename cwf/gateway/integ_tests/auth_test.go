/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integration

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

// - Initialize 3 UEs and initiate Authentication. Assert that it is successful.
// - Disconnect all UEs.
func TestAuthenticateMultipleUEs(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticate...")
	tr := NewTestRunner(t)
	ues, err := tr.ConfigUEs(3, MockPCRFRemote, MockOCSRemote)
	assert.NoError(t, err)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, tr.CleanUp())
	}()

	for _, ue := range ues {
		tr.AuthenticateAndAssertSuccess(ue.GetImsi())
		tr.Disconnect(ue.GetImsi())
	}
	time.Sleep(1 * time.Second)
}

// - Expect a CCR-I to come into PCRF, and return with Authentication Reject.
// - Configure a UE and trigger a authentication. Assert that the expectation was
//   met, and the authentication failed.
func TestAuthenticateFail(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateFail...")
	tr := NewTestRunner(t)
	assert.NoError(t, usePCRFMockDriver(MockPCRFRemote))
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver(MockPCRFRemote))
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1, MockPCRFRemote, MockOCSRemote)
	assert.NoError(t, err)

	// Test Authentication Fail
	imsiFail := ues[0].GetImsi()
	initRequest := protos.NewGxCCRequest(imsiFail, protos.CCRequestType_INITIAL, 1)
	initAnswer := protos.NewGxCCAnswer(diam.AuthenticationRejected).
		SetDynamicRuleInstalls([]*protos.RuleDefinition{getPassAllRuleDefinition("dynamic-pass-all", "mkey1", 100)})
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	defaultAnswer := protos.NewGxCCAnswer(diam.AuthenticationRejected)
	assert.NoError(t, setPCRFExpectations(MockPCRFRemote, []*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	tr.AuthenticateAndAssertFail(imsiFail)

	resultByIndex, errByIndex, err := getAssertExpectationsResult(MockPCRFRemote)
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{{ExpectationIndex: 0, ExpectationMet: true}}
	assert.ElementsMatch(t, expectedResult, resultByIndex)
	// Since CCR/A-I failed, there should be no rules installed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.Empty(t, recordsBySubID["IMSI"+imsiFail])
}

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install for a pass-all dynamic rule and 250KB of
//   quota.
//   Trigger a authentication and assert the CCR-I is received.
// - Generate traffic to put traffic through the newly installed rule.
func TestAuthenticateUplinkTraffic(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateUplinkTraffic...")
	tr := NewTestRunner(t)
	assert.NoError(t, usePCRFMockDriver(MockPCRFRemote))
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver(MockPCRFRemote))
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1, MockPCRFRemote, MockOCSRemote)
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
	assert.NoError(t, setPCRFExpectations(MockPCRFRemote, []*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "100K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	resultByIndex, errByIndex, err := getAssertExpectationsResult(MockPCRFRemote)
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
