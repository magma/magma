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

package mock_driver_test

import (
	"fmt"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/testcore/mock_driver"

	"github.com/stretchr/testify/assert"
)

type TestExpectation struct {
	request string
	answer  string
}

func (t TestExpectation) DoesMatch(iReq interface{}) error {
	request, ok := iReq.(string)
	if !ok {
		return fmt.Errorf("request is not of type string")
	}
	if request != t.request {
		return fmt.Errorf("Expected: %v, Received: %v", t.request, request)
	}
	return nil
}

func (t TestExpectation) GetAnswer() interface{} {
	return t.answer
}

func TestAllExpectationsMet(t *testing.T) {
	// define some expectations
	expectation1 := TestExpectation{request: "req1", answer: "ans1"}
	expectation2 := TestExpectation{request: "req2", answer: "ans2"}
	defaultAnswer := "default"

	expectations := []mock_driver.Expectation{expectation1, expectation2}
	// initialize expectations
	md := mock_driver.NewMockDriver(expectations, protos.UnexpectedRequestBehavior_CONTINUE_WITH_DEFAULT_ANSWER, defaultAnswer)

	// meet both expectations
	request1 := "req1"
	request2 := "req2"

	answer := md.GetAnswerFromExpectations(request1)
	assert.EqualValues(t, "ans1", answer)

	answer = md.GetAnswerFromExpectations(request2)
	assert.EqualValues(t, "ans2", answer)

	result, errs := md.AggregateResults()
	// assert there was no unexpected requests
	assert.Empty(t, errs)
	expectedResult := []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
		{ExpectationIndex: 1, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, result)
}

func TestSomeExpectationsNotMet(t *testing.T) {
	// base case : no expectations
	md := mock_driver.NewMockDriver(nil, protos.UnexpectedRequestBehavior_CONTINUE_WITH_DEFAULT_ANSWER, "default")

	request1 := "req1"
	answer := md.GetAnswerFromExpectations(request1)
	assert.Equal(t, "default", answer)

	result, errs := md.AggregateResults()
	// assert there was no unexpected requests
	assert.Empty(t, result)
	assert.Empty(t, errs)

	// Two expectations are set, we will meet only the first one
	expectation1 := TestExpectation{request: "req1", answer: "ans1"}
	expectation2 := TestExpectation{request: "req2", answer: "ans2"}
	expectation4 := TestExpectation{request: "req4", answer: "ans4"}

	expectations := []mock_driver.Expectation{expectation1, expectation2, expectation4}
	md = mock_driver.NewMockDriver(expectations, protos.UnexpectedRequestBehavior_CONTINUE_WITH_ERROR, nil)

	answer = md.GetAnswerFromExpectations(request1)
	assert.EqualValues(t, "ans1", answer)
	request3 := "bad-req2-1"
	answer = md.GetAnswerFromExpectations(request3)
	assert.EqualValues(t, nil, answer)

	request4 := "bad-req2-2"
	answer = md.GetAnswerFromExpectations(request4)
	assert.EqualValues(t, nil, answer)

	result, errs = md.AggregateResults()
	expectedResult := []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
		{ExpectationIndex: 1, ExpectationMet: false},
		{ExpectationIndex: 2, ExpectationMet: false},
	}
	assert.ElementsMatch(t, expectedResult, result)
	expectedErrors := []*protos.ErrorByIndex{
		{Index: 1, Error: "Expected: req2, Received: bad-req2-1"},
		{Index: 1, Error: "Expected: req2, Received: bad-req2-2"},
	}
	assert.ElementsMatch(t, expectedErrors, errs)
}
