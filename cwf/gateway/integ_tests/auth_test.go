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

	"magma/feg/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"
)

// - Initialize 3 UEs and initiate Authentication. Assert that it is successful.
// - Disconnect all UEs.
func TestAuthenticateMultipleUEs(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticate...")
	tr := NewTestRunner(t)
	ues, err := tr.ConfigUEs(3)
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
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	// Test Authentication Fail
	imsiFail := ues[0].GetImsi()
	initRequest := protos.NewGxCCRequest(imsiFail, protos.CCRequestType_INITIAL, 1)
	initAnswer := protos.NewGxCCAnswer(diam.AuthenticationRejected).
		SetDynamicRuleInstalls([]*protos.RuleDefinition{getPassAllRuleDefinition("dynamic-pass-all", "mkey1", 100)})
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	defaultAnswer := protos.NewGxCCAnswer(diam.AuthenticationRejected)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	tr.AuthenticateAndAssertFail(imsiFail)

	resultByIndex, errByIndex, err := getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{{ExpectationIndex: 0, ExpectationMet: true}}
	assert.ElementsMatch(t, expectedResult, resultByIndex)
	// Since CCR/A-I failed, there should be no rules installed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.Empty(t, recordsBySubID["IMSI"+imsiFail])
}
