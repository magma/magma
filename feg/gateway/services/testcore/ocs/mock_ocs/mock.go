/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_ocs

import (
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/testcore/mock_driver"

	"github.com/fiorix/go-diameter/v4/diam"
)

// Here we wrap the protobuf definitions to easily define instance methods
type GyExpectation struct {
	*protos.GyCreditControlExpectation
}

type GyAnswer struct {
	*protos.GyCreditControlAnswer
}

func (e GyExpectation) GetAnswer() interface{} {
	return GyAnswer{e.Answer}
}

func (e GyExpectation) DoesMatch(message interface{}) error {
	expected := e.ExpectedRequest
	ccr := message.(ccrMessage)
	expectedPK := mock_driver.NewCCRequestPK(expected.Imsi, expected.RequestType)
	actualImsi := ccr.GetIMSI()
	actualPK := mock_driver.NewCCRequestPK(actualImsi, protos.CCRequestType(ccr.RequestType))
	// For better readability of errors, we will check for the IMSI and the request type first.
	if expectedPK != actualPK {
		return fmt.Errorf("Expected: %v, Received: %v", expectedPK, actualPK)
	}
	if !compareMsccAgainstExpected(ccr.MSCC, expected.GetMscc(), expected.GetUsageReportDelta()) {
		return fmt.Errorf("For Request=%v, Expected: %v, Received: %v", actualPK, expected.GetMscc(), ccr.MSCC)
	}
	return nil
}

func (answer GyAnswer) toAVPs() ([]*diam.AVP, uint32) {
	avps := make([]*diam.AVP, 0, len(answer.GetQuotaGrants()))
	for _, grant := range answer.QuotaGrants {
		avps = append(avps, toGrantedUnitsAVP(grant.GetResultCode(), grant.GetValidityTime(), grant.GetGrantedServiceUnit(), grant.GetIsFinalCredit(), grant.GetFinalUnitAction(), grant.GetRatingGroup()))
	}
	return avps, answer.GetResultCode()
}

func compareMsccAgainstExpected(actualMscc []*ccrCredit, expectedMscc []*protos.MultipleServicesCreditControl, delta uint64) bool {
	if expectedMscc == nil {
		return true
	}
	expectedCreditByKey := toExpectedCreditByKey(expectedMscc)
	actualCreditByKey := toActualCreditByKey(actualMscc)
	for rg, expectedCredit := range expectedCreditByKey {
		actualCredit, exists := actualCreditByKey[rg]
		if !exists {
			return false
		}
		actualTotal := actualCredit.UsedServiceUnit.TotalOctets
		expectedTotal := expectedCredit.UsedServiceUnit.TotalOctets
		if !mock_driver.EqualWithinDelta(actualTotal, expectedTotal, delta) {
			return false
		}
	}
	return true
}

func toExpectedCreditByKey(mscc []*protos.MultipleServicesCreditControl) map[uint32]*protos.MultipleServicesCreditControl {
	msccByRG := map[uint32]*protos.MultipleServicesCreditControl{}
	for _, credit := range mscc {
		msccByRG[credit.RatingGroup] = credit
	}
	return msccByRG
}

func toActualCreditByKey(mscc []*ccrCredit) map[uint32]*ccrCredit {
	msccByRG := map[uint32]*ccrCredit{}
	for _, credit := range mscc {
		msccByRG[credit.RatingGroup] = credit
	}
	return msccByRG
}
