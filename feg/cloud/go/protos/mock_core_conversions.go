/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package protos

func (m *Expectation) ToExpectationResult(met bool) *ExpectationResult {
	return &ExpectationResult{Expectation: m, ExpectationMet: met}
}

func NewGxCreditControlExpectation() *GxCreditControlExpectation {
	return &GxCreditControlExpectation{}
}

func (m *GxCreditControlExpectation) Expect(ccr *GxCreditControlRequest) *GxCreditControlExpectation {
	m.ExpectedRequest = ccr
	return m
}

func (m GxCreditControlExpectation) Return(cca *GxCreditControlAnswer) *Expectation {
	m.Answer = cca
	return &Expectation{Expectation: &Expectation_GxCcExpectation{GxCcExpectation: &m}}
}
