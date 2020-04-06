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
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateMultipleUEs(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticate...")
	tr := NewTestRunner()
	ues, err := tr.ConfigUEs(3)
	assert.NoError(t, err)

	for _, ue := range ues {
		tr.AuthenticateAndAssertSuccess(t, ue.GetImsi())
		tr.Disconnect(ue.GetImsi())
	}
	time.Sleep(1 * time.Second)
	// Clear hss, ocs, and pcrf
	assert.NoError(t, tr.CleanUp())
}

func TestAuthenticateFail(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateFail...")
	tr := NewTestRunner()
	assert.NoError(t, usePCRFMockDriver())

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

	radiusP, err := tr.Authenticate(imsiFail)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.FailureCode))

	resultByIndex, errByIndex, err := getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{{ExpectationIndex: 0, ExpectationMet: true}}
	assert.ElementsMatch(t, expectedResult, resultByIndex)
	// Since CCR/A-I failed, there should be no rules installed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.Empty(t, recordsBySubID["IMSI"+imsiFail])

	// Clear hss, ocs, and pcrf
	assert.NoError(t, clearPCRFMockDriver())
	assert.NoError(t, tr.CleanUp())
}
