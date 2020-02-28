/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package expectations

import (
	"magma/feg/cloud/go/protos"
	merrors "magma/orc8r/lib/go/errors"
)

type ExpectationManagement struct {
	upcomingExpectations  []*protos.Expectation
	fulfilledExpectations []*protos.ExpectationResult
	unexpectedRequests    []interface{}
	failureBehavior       protos.SetExpectationsRequest_UnexpectedRequestBehavior
	defaultAnswer         interface{}
}

type Matcher func(request interface{}, expectation *protos.Expectation) (bool, interface{})

func (e *ExpectationManagement) InitializeExpectations(expectations []*protos.Expectation,
	failureBehavior protos.SetExpectationsRequest_UnexpectedRequestBehavior, defaultAnswer interface{}) {
	e.upcomingExpectations = expectations
	e.fulfilledExpectations = []*protos.ExpectationResult{}
	e.unexpectedRequests = []interface{}{}
	e.failureBehavior = failureBehavior
	e.defaultAnswer = defaultAnswer
}

func (e *ExpectationManagement) GetNextAnswer(matcher Matcher, message interface{}) (interface{}, error) {
	if len(e.upcomingExpectations) == 0 {
		return e.getDefaultAnswerOrNil(), merrors.ErrNotFound
	}
	nextExpectation := e.upcomingExpectations[0]
	matched, answerMessage := matcher(message, nextExpectation)
	if matched {
		e.fulfilledExpectations = append(e.fulfilledExpectations, nextExpectation.ToExpectationResult(true))
		e.upcomingExpectations = e.upcomingExpectations[1:]
		return answerMessage, nil
	} else {
		return e.getDefaultAnswerOrNil(), merrors.ErrNotFound
	}
}

func (e *ExpectationManagement) MarkExpectationAsFulfilled(expectation *protos.Expectation) {
	e.fulfilledExpectations = append(e.fulfilledExpectations, &protos.ExpectationResult{Expectation: expectation, ExpectationMet: true})
	// remove the fulfilled expectation from upcoming
	if len(e.upcomingExpectations) != 0 {
		e.upcomingExpectations = e.upcomingExpectations[1:]
	}
}

func (e *ExpectationManagement) AddToUnexpectedRequests(request interface{}) {
	e.unexpectedRequests = append(e.unexpectedRequests, request)
}

func (e *ExpectationManagement) getDefaultAnswerOrNil() interface{} {
	if e.failureBehavior == protos.SetExpectationsRequest_CONTINUE_WITH_DEFAULT_ANSWER {
		return e.defaultAnswer
	}
	return nil
}

func (e *ExpectationManagement) GetResultsAndClear() ([]*protos.ExpectationResult, []interface{}) {
	result := []*protos.ExpectationResult{}
	result = append(result, e.fulfilledExpectations...)
	// ignore request ordering for now
	for _, expectation := range e.upcomingExpectations {
		result = append(result, &protos.ExpectationResult{Expectation: expectation, ExpectationMet: false})
	}
	unexpectedRequests := []interface{}{}
	unexpectedRequests = append(unexpectedRequests, e.unexpectedRequests...)

	e.unexpectedRequests = nil
	e.fulfilledExpectations = nil
	e.fulfilledExpectations = nil
	return result, unexpectedRequests
}
