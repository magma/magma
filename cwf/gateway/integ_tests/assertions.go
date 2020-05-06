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
	"reflect"

	"fbc/lib/go/radius/rfc2869"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"

	"github.com/stretchr/testify/assert"
)

// Trigger a UE Authentication with the IMSI. Assert that the authentication
// succeeded.
func (tr *TestRunner) AuthenticateAndAssertSuccess(imsi string) {
	radiusP, err := tr.Authenticate(imsi, defaultCalledStationID)
	assert.NoError(tr.t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(tr.t, eapMessage, fmt.Sprintf("EAP Message from authentication is nil"))
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), fmt.Sprintf("UE Authentication did not return success"))
}

// Trigger a UE Authentication with the IMSI and called station ID.
// Assert that the authentication succeeded.
func (tr *TestRunner) AuthenticateWithCalledIDAndAssertSuccess(imsi, calledStationID string) {
	radiusP, err := tr.Authenticate(imsi, calledStationID)
	assert.NoError(tr.t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(tr.t, eapMessage, fmt.Sprintf("EAP Message from authentication is nil"))
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), fmt.Sprintf("UE Authentication did not return success"))
}

// Trigger a UE Authentication with the IMSI. Assert that the authentication
// failed.
func (tr *TestRunner) AuthenticateAndAssertFail(imsi string) {
	radiusP, err := tr.Authenticate(imsi, defaultCalledStationID)
	assert.NoError(tr.t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(tr.t, eapMessage)
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.FailureCode))
}

// Trigger a UE Disconnect and assert it succeeds.
func (tr *TestRunner) DisconnectAndAssertSuccess(imsi string) {
	_, err := tr.Disconnect(imsi, defaultCalledStationID)
	assert.NoError(tr.t, err)
}

// Query policy usage records from pipelined, and assert that there is no entry
// for the IMSI. This means no policy flows are installed for the IMSI.
func (tr *TestRunner) AssertPolicyEnforcementRecordIsNil(imsi string) {
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(tr.t, err)
	assert.Empty(tr.t, recordsBySubID[prependIMSIPrefix(imsi)])
}

// Query assertion result from MockPCRF and assert all expectations were met.
// Only applicable when MockDriver is used.
func (tr *TestRunner) AssertAllGxExpectationsMetNoError() {
	fmt.Println("Asserting all Gx expectations were met...")
	resultByIndex, errByIndex, err := getPCRFAssertExpectationsResult()
	tr.assertAllExpectationsMetNoError(resultByIndex, errByIndex, err)
	fmt.Println("Passed!")
}

// Query assertion result from MockOCS and assert all expectations were met.
// Only applicable when MockDriver is used.
func (tr *TestRunner) AssertAllGyExpectationsMetNoError() {
	fmt.Println("Asserting all Gy expectations were met...")
	resultByIndex, errByIndex, err := getOCSAssertExpectationsResult()
	tr.assertAllExpectationsMetNoError(resultByIndex, errByIndex, err)
	fmt.Println("Passed!")
}

func (tr *TestRunner) assertAllExpectationsMetNoError(resByIdx []*protos.ExpectationResult, errByIdx []*protos.ErrorByIndex, err error) {
	expectedResults := makeDefaultExpectationResults(len(resByIdx))
	assert.NoError(tr.t, err)
	assert.Empty(tr.t, errByIdx)
	assert.ElementsMatch(tr.t, expectedResults, resByIdx)
}

func makeDefaultExpectationResults(n int) []*protos.ExpectationResult {
	expectedResults := make([]*protos.ExpectationResult, n)
	for i := 0; i < n; i++ {
		expectedResults[i] = &protos.ExpectationResult{
			ExpectationIndex: int32(i),
			ExpectationMet:   true,
		}
	}
	return expectedResults
}
