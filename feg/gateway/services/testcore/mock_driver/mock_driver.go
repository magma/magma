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

package mock_driver

import (
	"sync"

	"magma/feg/cloud/go/protos"

	"github.com/golang/glog"
)

type Expectation interface {
	DoesMatch(interface{}) error
	GetAnswer() interface{}
}

type MockDriver struct {
	expectationsSet bool
	// Should not modified to maintain the expectation order
	expectations              []Expectation
	expectationIndex          int
	resultByIndex             map[int]*protos.ExpectationResult
	errorMessages             []*protos.ErrorByIndex
	unexpectedRequestBehavior protos.UnexpectedRequestBehavior
	defaultAnswer             interface{}
	sync.Mutex
}

func NewMockDriver(expectations []Expectation, behavior protos.UnexpectedRequestBehavior, defaultAnswer interface{}) *MockDriver {
	glog.Infof("MockDriver: Initializing a new set of Expectations: %s", expectations)
	resultByIndex := make(map[int]*protos.ExpectationResult, len(expectations))
	for i := range expectations {
		resultByIndex[i] = &protos.ExpectationResult{ExpectationIndex: int32(i), ExpectationMet: false}
	}
	return &MockDriver{
		expectationsSet:           true,
		expectations:              expectations,
		expectationIndex:          0,
		resultByIndex:             resultByIndex,
		unexpectedRequestBehavior: behavior,
		defaultAnswer:             defaultAnswer,
		errorMessages:             []*protos.ErrorByIndex{},
	}
}

// GetAnswerFromExpectations will use the message passed in to determine if
// the message matches the next upcoming expectation. If it does, it will
// return the answer specified in the expectation.
// On failure, if the unexpected request behavior is set to
// CONTINUE_WITH_DEFAULT_ANSWER, it will return the default answer.
// Otherwise, it will return nil.
func (e *MockDriver) GetAnswerFromExpectations(message interface{}) interface{} {
	if !e.expectationsSet {
		return nil
	}
	if len(e.expectations) == 0 || e.expectationIndex >= len(e.expectations) {
		return e.getAnswerForUnexpectedMessage()
	}
	expectation := e.expectations[e.expectationIndex]
	err := expectation.DoesMatch(message)
	if err != nil {
		glog.Errorf("MockDriver: Expectation was not met: %s, err: %s", expectation, err)
		errByIndex := &protos.ErrorByIndex{Index: int32(e.expectationIndex), Error: err.Error()}
		e.errorMessages = append(e.errorMessages, errByIndex)
		return e.getAnswerForUnexpectedMessage()
	}

	e.resultByIndex[e.expectationIndex].ExpectationMet = true
	e.expectationIndex++
	return expectation.GetAnswer()
}

// AggregateResults will aggregate resultByIndex and errorsByIndex.
func (e *MockDriver) AggregateResults() ([]*protos.ExpectationResult, []*protos.ErrorByIndex) {
	results := make([]*protos.ExpectationResult, len(e.expectations))
	for i := range e.expectations {
		results[i] = e.resultByIndex[i]
	}
	e.expectationsSet = false
	return results, e.errorMessages
}

func (e *MockDriver) getAnswerForUnexpectedMessage() interface{} {
	switch e.unexpectedRequestBehavior {
	case protos.UnexpectedRequestBehavior_CONTINUE_WITH_DEFAULT_ANSWER:
		glog.Infof("MockDriver: Returning default answer for an unexpected request")
		return e.defaultAnswer
	default:
		return nil
	}
}
