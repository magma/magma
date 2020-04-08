/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_ocs_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/eap/test"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/testcore/ocs/mock_ocs"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestOCSExpectations(t *testing.T) {
	serverConfig := diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:     "127.0.0.1:0",
		Protocol: "tcp"},
	}

	initRequest := fegprotos.NewGyCCRequest(test.IMSI1, fegprotos.CCRequestType_INITIAL, 1)
	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 100,
		},
		IsFinalCredit:   true,
		FinalUnitAction: fegprotos.FinalUnitAction_Terminate,
		ResultCode:      2001,
	}
	initAnswer := fegprotos.NewGyCCAnswer(diameter.SuccessCode).SetQuotaGrant(quotaGrant)
	initExpectation := fegprotos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	updateReq := fegprotos.NewGyCCRequest(test.IMSI1, fegprotos.CCRequestType_UPDATE, 3).
		SetMSCC(&fegprotos.MultipleServicesCreditControl{RatingGroup: 1, UsedServiceUnit: &fegprotos.Octets{TotalOctets: 100}})
	updateAnswer := fegprotos.NewGyCCAnswer(diam.Success)
	updateExpectation := fegprotos.NewGyCreditControlExpectation().Expect(updateReq).Return(updateAnswer)

	terminateReq := fegprotos.NewGyCCRequest(test.IMSI1, fegprotos.CCRequestType_TERMINATION, 4)
	terminateAnswer := fegprotos.NewGyCCAnswer(diam.Success)
	terminateExpectation := fegprotos.NewGyCreditControlExpectation().Expect(terminateReq).Return(terminateAnswer)

	expectations := []*fegprotos.GyCreditControlExpectation{initExpectation, updateExpectation, terminateExpectation}
	failureBehavior := fegprotos.UnexpectedRequestBehavior_CONTINUE_WITH_DEFAULT_ANSWER
	defaultCCA := &fegprotos.GyCreditControlAnswer{}

	clientConfig := getClientConfig()
	ocs := startServerWithExpectations(clientConfig, &serverConfig, gy.PerSessionInit, expectations, failureBehavior, defaultCCA)
	gyGlobalConfig := getGyGlobalConfig("")
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)
	ocs.CreateAccount(context.Background(), &lteprotos.SubscriberID{Type: lteprotos.SubscriberID_IMSI, Id: test.IMSI1})
	ocs.CreateAccount(context.Background(), &lteprotos.SubscriberID{Type: lteprotos.SubscriberID_IMSI, Id: test.IMSI2})

	// send Init
	ccrInit := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          test.IMSI1,
		RequestNumber: 1,
	}
	done := make(chan interface{}, 1000)
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	actualAnswer := gy.GetAnswer(done)
	assertCCAIsEqualToExpectedAnswer(t, actualAnswer, initAnswer)

	ccrUpdate := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTUpdate,
		IMSI:          test.IMSI1,
		RequestNumber: 2,
		Credits:       []*gy.UsedCredits{{TotalOctets: 100, RatingGroup: 1}},
	}
	done = make(chan interface{}, 1000)
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrUpdate))
	actualAnswer = gy.GetAnswer(done)
	assertCCAIsEqualToExpectedAnswer(t, actualAnswer, updateAnswer)

	res, err := ocs.AssertExpectations(context.Background(), &orcprotos.Void{})
	assert.Nil(t, err)
	expectedResult := []*fegprotos.ExpectationResult{
		{ExpectationMet: true, ExpectationIndex: 0},
		{ExpectationMet: true, ExpectationIndex: 1},
		{ExpectationMet: false, ExpectationIndex: 2}, //no terminate
	}
	assert.ElementsMatch(t, expectedResult, res.Results)
}

func assertCCAIsEqualToExpectedAnswer(t *testing.T, actual *gy.CreditControlAnswer, expectation *fegprotos.GyCreditControlAnswer) {
	assert.Equal(t, actual.ResultCode, expectation.ResultCode)
	actualCreditsByKey := getCreditByKey(actual.Credits)
	expectedCreditsByKey := getExpectedCreditByKey(expectation.QuotaGrants)
	for rg, credit := range expectedCreditsByKey {
		actualCredit, found := actualCreditsByKey[rg]
		assert.True(t, found, fmt.Sprintf("Expected %v in answer but it doesn't exist", credit))
		assert.Equal(t, int(credit.GetResultCode()), int(actualCredit.ResultCode))
		assert.Equal(t, credit.GetValidityTime(), actualCredit.ValidityTime)
		assert.Equal(t, credit.GetIsFinalCredit(), actualCredit.IsFinal)
		if credit.IsFinalCredit {
			assert.Equal(t, int(credit.GetFinalUnitAction()), int(actualCredit.FinalAction))
			assert.Equal(t, credit.GetRedirectServer().GetRedirectServerAddress(), actualCredit.RedirectServer.RedirectServerAddress)
		}
		expectedOctet := credit.GetGrantedServiceUnit()
		actualOctet := actualCredit.GrantedUnits
		assert.Equal(t, expectedOctet.GetTotalOctets(), swag.Uint64Value(actualOctet.TotalOctets))
		assert.Equal(t, expectedOctet.GetOutputOctets(), swag.Uint64Value(actualOctet.OutputOctets))
		assert.Equal(t, expectedOctet.GetInputOctets(), swag.Uint64Value(actualOctet.InputOctets))
	}
}

func startServerWithExpectations(
	client *diameter.DiameterClientConfig,
	server *diameter.DiameterServerConfig,
	initMethod gy.InitMethod,
	expectations []*fegprotos.GyCreditControlExpectation,
	failureBehavior fegprotos.UnexpectedRequestBehavior,
	defaultCCA *fegprotos.GyCreditControlAnswer,
) *mock_ocs.OCSDiamServer {
	serverStarted := make(chan struct{})
	ocs := mock_ocs.NewOCSDiamServer(
		client,
		&mock_ocs.OCSConfig{
			ServerConfig: server,
			GyInitMethod: initMethod,
		},
	)
	go func() {
		log.Printf("Starting server")
		ctx := context.Background()
		ocs.SetOCSSettings(ctx, &fegprotos.OCSConfig{UseMockDriver: true})
		ocs.SetExpectations(ctx, &fegprotos.GyCreditControlExpectations{
			Expectations:              expectations,
			UnexpectedRequestBehavior: failureBehavior,
			GyDefaultCca:              defaultCCA,
		})

		lis, err := ocs.StartListener()
		if err != nil {
			log.Fatalf("Could not start listener for PCRF, %s", err.Error())
		}
		server.Addr = lis.Addr().String()
		serverStarted <- struct{}{}
		err = ocs.Start(lis)
		if err != nil {
			log.Fatalf("Could not start test PCRF server, %s", err.Error())
			return
		}
	}()
	<-serverStarted
	return ocs
}

func getCreditByKey(credits []*gy.ReceivedCredits) map[uint32]*gy.ReceivedCredits {
	creditsByKey := make(map[uint32]*gy.ReceivedCredits, len(credits))
	for _, credit := range credits {
		creditsByKey[credit.RatingGroup] = credit
	}
	return creditsByKey
}

func getExpectedCreditByKey(credits []*fegprotos.QuotaGrant) map[uint32]*fegprotos.QuotaGrant {
	creditsByKey := make(map[uint32]*fegprotos.QuotaGrant, len(credits))
	for _, credit := range credits {
		creditsByKey[credit.RatingGroup] = credit
	}
	return creditsByKey
}

func grantedServiceUnitToOctet(gsu *credit_control.GrantedServiceUnit) *fegprotos.Octets {
	return &fegprotos.Octets{
		TotalOctets:  swag.Uint64Value(gsu.TotalOctets),
		InputOctets:  swag.Uint64Value(gsu.InputOctets),
		OutputOctets: swag.Uint64Value(gsu.OutputOctets),
	}
}

func getGyGlobalConfig(ocsOverwriteApn string) *gy.GyGlobalConfig {
	return &gy.GyGlobalConfig{
		OCSOverwriteApn: ocsOverwriteApn,
	}
}

func getClientConfig() *diameter.DiameterClientConfig {
	return &diameter.DiameterClientConfig{
		Host:        "test.test.com",
		Realm:       "test.com",
		ProductName: "gy_test",
		AppID:       diam.CHARGING_CONTROL_APP_ID,
	}
}

func getReAuthHandler() gy.ChargingReAuthHandler {
	return func(request *gy.ChargingReAuthRequest) *gy.ChargingReAuthAnswer {
		return &gy.ChargingReAuthAnswer{
			SessionID:  request.SessionID,
			ResultCode: diam.Success,
		}
	}
}
