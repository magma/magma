/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package integration

import (
	"fmt"
	"reflect"
	"time"

	"fbc/lib/go/radius"
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
	assert.NotNil(tr.t, eapMessage, "EAP Message from authentication is nil")
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), "UE Authentication did not return success")
}

// Trigger a UE Authentication with the IMSI and called station ID.
// Assert that the authentication succeeded.
func (tr *TestRunner) AuthenticateWithCalledIDAndAssertSuccess(imsi, calledStationID string) {
	radiusP, err := tr.Authenticate(imsi, calledStationID)
	assert.NoError(tr.t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(tr.t, eapMessage, "EAP Message from authentication is nil")
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), "UE Authentication did not return success")
}

// AuthenticateAndAssertSuccessWithRetries triggers a UE Authentication with the IMSI. Assert that the authentication
// succeeded with retrials. Use this function only for those tests that deal with service restart
// Otherwise, use the client without retries. Retries shuldn't happen on a healthy system
func (tr *TestRunner) AuthenticateAndAssertSuccessWithRetries(imsi string, maxRetries int) {
	if maxRetries < 0 {
		panic("Authentication maxRetries must be positive!")
	}
	var (
		err           error
		radiusP       *radius.Packet
		totalAttempts = maxRetries + 1
		eapMessage    radius.Attribute
	)
	for i := 0; i < totalAttempts; i++ {
		radiusP, err = tr.Authenticate(imsi, defaultCalledStationID)
		eapMessage = radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
		// do not print the info for the last attempt
		if i < totalAttempts-1 {
			if err != nil {
				fmt.Printf("...Authentication failed with radius message nul. Retrying...!\n")
				time.Sleep(1 * time.Second)
				continue
			}
			if eapMessage == nil || !reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode) {
				fmt.Printf("...Authentication failed with eap message either nul or not succelful: %+v. Retrying...!\n", eapMessage)
				time.Sleep(1 * time.Second)
				continue
			}
		}
		break
	}
	assert.NoError(tr.t, err)
	assert.NotNil(tr.t, eapMessage, "EAP Message from authentication is nil")
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), "UE Authentication did not return success")

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
	// This string is taken from pipelined.yml enforcement information
	delete(recordsBySubID[prependIMSIPrefix(imsi)], "internal_default_drop_flow_rule")
	assert.Empty(tr.t, recordsBySubID[prependIMSIPrefix(imsi)])
}

func (tr *TestRunner) AssertEventuallyAllRulesRemovedAfterDisconnect(imsi string) {
	checkFn := func() bool {
		fmt.Printf("Waiting until all rules are removed in enforcement stats for %s...\n", imsi)
		records, err := tr.GetPolicyUsage()
		if err != nil {
			return false
		}
		if len(records[prependIMSIPrefix(imsi)]) == 0 {
			return true
		}
		return false
	}
	assert.Eventually(tr.t, checkFn, 10*time.Second, 2*time.Second)
	fmt.Println("All enforcement stats are gone!")
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
	matches := assert.ElementsMatch(tr.t, expectedResults, resByIdx)
	if !matches {
		tr.t.Log(errByIdx)
	}
}

func (tr *TestRunner) AssertPolicyUsage(imsi, rule string, minBytes, maxBytes uint64) uint64 {
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(tr.t, err)
	assert.NotNil(tr.t, recordsBySubID[prependIMSIPrefix(imsi)], fmt.Sprintf("Policy usage record for %s not found", imsi))
	record := recordsBySubID[prependIMSIPrefix(imsi)][rule]
	assert.NotNil(tr.t, record, fmt.Sprintf("Policy usage record for %s not found for %s", rule, imsi))
	if record != nil {
		assert.GreaterOrEqual(tr.t, record.BytesTx, minBytes, fmt.Sprintf("%s actual=%d < expected=%d", record.RuleId, record.BytesTx, minBytes))
		assert.LessOrEqual(tr.t, record.BytesTx, maxBytes, fmt.Sprintf("%s actual=%d > expected=%d", record.RuleId, record.BytesTx, maxBytes))
		return record.BytesTx
	}
	return 0
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
