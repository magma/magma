/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_driver

import (
	"fmt"

	"magma/feg/cloud/go/protos"
)

type Expectation interface {
	DoesMatch(interface{}) bool
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
}

func NewMockDriver(expectations []Expectation, behavior protos.UnexpectedRequestBehavior, defaultAnswer interface{}) *MockDriver {
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
// Otherwise, if the unexpected request behavior is set to CONTINUE_WITH_DEFAULT_ANSWER,
// it will return the default answer.
// Otherwise, it will return nil.
func (e *MockDriver) GetAnswerFromExpectations(message interface{}) interface{} {
	if !e.expectationsSet {
		return nil
	}
	if len(e.expectations) == 0 {
		return e.getAnswerForUnexpectedMessage()
	}
	expectation := e.expectations[e.expectationIndex]
	doesMatch := expectation.DoesMatch(message)
	if !doesMatch {
		err := &protos.ErrorByIndex{Index: int32(e.expectationIndex), Error: fmt.Sprintf("Expected: %v, Received: %v", expectation, message)}
		e.errorMessages = append(e.errorMessages, err)
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
		return e.defaultAnswer
	default:
		return nil
	}
}
