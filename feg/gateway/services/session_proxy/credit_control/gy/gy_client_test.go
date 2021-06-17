/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package diameter_test tests diameter calls within the magma setting
package gy_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/testcore/ocs/mock_ocs"
	"magma/lte/cloud/go/protos"
)

const (
	testIMSI1      = "000000000000001"
	testIMSI2      = "4321"
	returnedOctets = 1024
	validityTime   = 3600
	restrictRule   = "restrict-rule-1"
)

var (
	defaultLocalServerConfig = diameter.DiameterServerConfig{
		DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:     "127.0.0.1:0",
			Protocol: "tcp",
		},
	}
	defaultRSU = &protos.RequestedUnits{Total: 10000, Tx: 10000, Rx: 10000}
)

var defaultfinalUnitConfig = mock_ocs.FinalUnitIndication{FinalUnitAction: fegprotos.FinalUnitAction(gy.Terminate)}

// TestGyClient tests CCR init, update, and terminate messages using a fake
// server
func TestGyClient(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	ocs := startServer(clientConfig, &serverConfig, gy.PerSessionInit, defaultfinalUnitConfig)
	seedAccountConfigurations(ocs)
	gyGlobalConfig := getGyGlobalConfig("", "")
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)
	si11 := uint32(11)
	// send init
	ccrInit := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 0,
		UeIPV4:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
		Apn:           "gy.Apn.magma.com",
		Credits: []*gy.UsedCredits{
			{
				RatingGroup:    1,
				RequestedUnits: defaultRSU,
			},
			{
				RatingGroup:    2,
				RequestedUnits: defaultRSU,
			},
			{
				RatingGroup:       3,
				ServiceIdentifier: &si11,
				RequestedUnits:    defaultRSU,
			},
		},
	}
	done := make(chan interface{}, 1000)

	log.Printf("Sending CCR-Init")
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gy.GetAnswer(done)
	log.Printf("Received CCA-Init")

	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assert.Equal(t, 3, len(answer.Credits))
	assert.Equal(t, uint32(diam.Success), answer.ResultCode)
	assertReceivedAPNonOCS(t, ocs, ccrInit.Apn)

	// send multiple updates
	ccrUpdates := []*gy.CreditControlRequest{
		{
			SessionID:     "1",
			Type:          credit_control.CRTUpdate,
			IMSI:          testIMSI1,
			RequestNumber: 1,
			Credits: []*gy.UsedCredits{{
				RatingGroup:    1,
				InputOctets:    1024,
				OutputOctets:   2048,
				TotalOctets:    3072,
				RequestedUnits: defaultRSU,
			},
			}},
		{
			SessionID:     "2",
			Type:          credit_control.CRTUpdate,
			IMSI:          testIMSI2,
			RequestNumber: 1,
			Credits: []*gy.UsedCredits{{
				RatingGroup:    1,
				InputOctets:    1024,
				OutputOctets:   2048,
				TotalOctets:    3072,
				RequestedUnits: defaultRSU,
			},
			}},
	}

	for _, update := range ccrUpdates {
		assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, update))
	}

	for i := 0; i < 2; i++ {
		update := gy.GetAnswer(done)
		assert.Equal(t, uint64(returnedOctets), *update.Credits[0].GrantedUnits.TotalOctets)
		assert.Equal(t, uint32(validityTime), update.Credits[0].ValidityTime)
		assert.Equal(t, ccrUpdates[i].SessionID, update.SessionID)
		assert.Equal(t, ccrUpdates[i].RequestNumber, update.RequestNumber)
	}

	// send terminates
	ccrTerminate := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTTerminate,
		IMSI:          testIMSI1,
		RequestNumber: 2,
		Credits: []*gy.UsedCredits{{
			RatingGroup:  1,
			InputOctets:  1024,
			OutputOctets: 2048,
			TotalOctets:  3072,
		}},
	}
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrTerminate))
	terminate := gy.GetAnswer(done)
	assert.Equal(t, uint32(diam.Success), terminate.ResultCode)
	assert.Equal(t, 0, len(terminate.Credits))
	assert.Equal(t, ccrTerminate.SessionID, terminate.SessionID)
	assert.Equal(t, ccrTerminate.RequestNumber, terminate.RequestNumber)

	// Connection disabling should cause CCR to fail
	gyClient.DisableConnections(10 * time.Second)
	assert.Error(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))

	// CCR Success after enabling connections
	gyClient.EnableConnections()
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
}

// TestGyClient test different options on global configuration
func TestGyClientWithGyGlobalConf(t *testing.T) {
	serverConfig := defaultLocalServerConfig

	clientConfig := getClientConfig()
	ocs := startServer(clientConfig, &serverConfig, gy.PerSessionInit, defaultfinalUnitConfig)
	seedAccountConfigurations(ocs)
	matchApn := ".*\\.magma.*"
	overwriteApn := "gy.Apn.magma.com"
	gyGlobalConfig := getGyGlobalConfig(matchApn, overwriteApn)
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)

	// send init
	ccrInit := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 0,
		UeIPV4:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
		Apn:           "gy.Apn.magma.com",
		Credits: []*gy.UsedCredits{
			{
				RatingGroup:    1,
				RequestedUnits: defaultRSU,
			},
		},
	}
	done := make(chan interface{}, 1000)

	log.Printf("Sending CCR-Init with custom global parameters")
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gy.GetAnswer(done)
	log.Printf("Received CCA-Init")
	assert.Equal(t, uint32(diam.Success), answer.ResultCode)
	assertReceivedAPNonOCS(t, ocs, overwriteApn)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assert.Equal(t, 1, len(answer.Credits))
}

func TestGyClientOutOfCredit(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	ocs := startServer(clientConfig, &serverConfig, gy.PerSessionInit, defaultfinalUnitConfig)
	seedAccountConfigurations(ocs)
	gyGlobalConfig := getGyGlobalConfig("", "")
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)

	// send init
	ccrInit := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 0,
		UeIPV4:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
		Credits: []*gy.UsedCredits{
			{
				RatingGroup:    1,
				RequestedUnits: defaultRSU,
			},
		},
	}
	done := make(chan interface{}, 1000)
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gy.GetAnswer(done)
	assert.Equal(t, uint32(diam.Success), answer.ResultCode)

	// send request with (total credits - used credits) < max usage (final units)
	ccrUpdate := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTUpdate,
		IMSI:          testIMSI1,
		RequestNumber: 1,
		Credits: []*gy.UsedCredits{{
			RatingGroup:    1,
			InputOctets:    999990,
			OutputOctets:   0,
			TotalOctets:    999990,
			RequestedUnits: defaultRSU,
		}},
	}

	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrUpdate))
	update := gy.GetAnswer(done)
	assert.Equal(t, uint32(diam.Success), update.ResultCode)
	assert.Equal(t, uint64(10), *update.Credits[0].GrantedUnits.TotalOctets)
	assert.Equal(t, gy.Terminate, update.Credits[0].FinalUnitIndication.FinalAction)
}

func TestGyClientPerKeyInit(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	ocs := startServer(clientConfig, &serverConfig, gy.PerKeyInit, defaultfinalUnitConfig)
	seedAccountConfigurations(ocs)
	gyGlobalConfig := getGyGlobalConfig("", "")
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)

	// send inits
	ccrInits := []*gy.CreditControlRequest{
		{
			SessionID:     "1",
			Type:          credit_control.CRTInit,
			IMSI:          testIMSI1,
			RequestNumber: 1,
			UeIPV4:        "192.168.1.1",
			SpgwIPV4:      "10.10.10.10",
			Credits: []*gy.UsedCredits{{
				RatingGroup:    1,
				RequestedUnits: defaultRSU,
			},
			}},
		{
			SessionID:     "1",
			Type:          credit_control.CRTInit,
			IMSI:          testIMSI1,
			RequestNumber: 2,
			UeIPV4:        "192.168.1.1",
			SpgwIPV4:      "10.10.10.10",
			Credits: []*gy.UsedCredits{{
				RatingGroup:    2,
				RequestedUnits: defaultRSU,
			}},
		},
	}
	done := make(chan interface{}, 1000)

	log.Printf("Sending CCR-Updates")
	for _, init := range ccrInits {
		assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, init))
	}

	for i := 0; i < 2; i++ {
		update := gy.GetAnswer(done)
		assert.Equal(t, uint64(returnedOctets), *update.Credits[0].GrantedUnits.TotalOctets)
		assert.Equal(t, uint32(validityTime), update.Credits[0].ValidityTime)
		assert.Equal(t, ccrInits[i].SessionID, update.SessionID)
		assert.Equal(t, ccrInits[i].RequestNumber, update.RequestNumber)
	}
}

func TestGyClientMultipleCredits(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	ocs := startServer(clientConfig, &serverConfig, gy.PerKeyInit, defaultfinalUnitConfig)
	seedAccountConfigurations(ocs)
	gyGlobalConfig := getGyGlobalConfig("", "")
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)

	// send inits
	ccrInit := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 1,
		UeIPV4:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
		Credits: []*gy.UsedCredits{
			{
				RatingGroup:    1,
				RequestedUnits: defaultRSU,
			},
			{
				RatingGroup:    2,
				RequestedUnits: defaultRSU,
			},
			{
				RatingGroup:    3,
				RequestedUnits: defaultRSU,
			},
		},
	}
	done := make(chan interface{}, 1000)

	log.Printf("Sending CCR-Init")
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))

	ans := gy.GetAnswer(done)
	assert.Equal(t, uint32(diam.Success), ans.ResultCode)
	assert.Equal(t, ans.SessionID, ccrInit.SessionID)
	assert.Equal(t, ans.RequestNumber, ccrInit.RequestNumber)
	assert.Equal(t, 3, len(ans.Credits))
	for _, credit := range ans.Credits {
		assert.Contains(t, []uint32{1, 2, 3}, credit.RatingGroup)
		assert.Equal(t, uint64(returnedOctets), *credit.GrantedUnits.TotalOctets)
		assert.Equal(t, uint32(validityTime), credit.ValidityTime)
	}
}

func TestGyReAuth(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	ocs := startServer(clientConfig, &serverConfig, gy.PerKeyInit, defaultfinalUnitConfig)
	seedAccountConfigurations(ocs)
	gyGlobalConfig := getGyGlobalConfig("", "")
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)

	// send one init to set user context in OCS
	sessionID := fmt.Sprintf("IMSI%s-%d", testIMSI1, 1234)
	ccrInit := &gy.CreditControlRequest{
		SessionID:     sessionID,
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 1,
		UeIPV4:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
		Credits: []*gy.UsedCredits{
			{
				RatingGroup:    1,
				RequestedUnits: defaultRSU,
			},
		},
	}
	done := make(chan interface{}, 1000)

	log.Printf("Sending CCR-Init")
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gy.GetAnswer(done)
	assert.Equal(t, uint32(diam.Success), answer.ResultCode)

	// success reauth
	var rg uint32 = 1
	raa, err := ocs.ReAuth(
		context.Background(),
		&fegprotos.ChargingReAuthTarget{Imsi: testIMSI1, RatingGroup: rg},
	)
	assert.NoError(t, err)
	assert.Equal(t, sessionID, raa.SessionId)
	assert.Equal(t, uint32(diam.Success), raa.ResultCode)
}

func TestGyClientOutOfCreditRestrict(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	finalUnitConfig := mock_ocs.FinalUnitIndication{
		RestrictRules:   []string{restrictRule},
		FinalUnitAction: fegprotos.FinalUnitAction(gy.RestrictAccess),
	}
	ocs := startServer(clientConfig, &serverConfig, gy.PerSessionInit, finalUnitConfig)
	seedAccountConfigurations(ocs)
	gyGlobalConfig := getGyGlobalConfig("", "")
	gyClient := gy.NewGyClient(
		clientConfig,
		&serverConfig,
		getReAuthHandler(), nil, gyGlobalConfig,
	)

	// send init
	ccrInit := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 0,
		UeIPV4:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
		Credits: []*gy.UsedCredits{
			{
				RatingGroup:    1,
				RequestedUnits: defaultRSU,
			},
		},
	}
	done := make(chan interface{}, 1000)
	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gy.GetAnswer(done)
	assert.Equal(t, uint32(diam.Success), answer.ResultCode)

	// send request with (total credits - used credits) < max usage (final units)
	ccrUpdate := &gy.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTUpdate,
		IMSI:          testIMSI1,
		RequestNumber: 1,
		Credits: []*gy.UsedCredits{{
			RatingGroup:    1,
			InputOctets:    999990,
			OutputOctets:   0,
			TotalOctets:    999990,
			RequestedUnits: defaultRSU,
		}},
	}

	assert.NoError(t, gyClient.SendCreditControlRequest(&serverConfig, done, ccrUpdate))
	update := gy.GetAnswer(done)
	assert.Equal(t, uint64(10), *update.Credits[0].GrantedUnits.TotalOctets)
	assert.Equal(t, gy.RestrictAccess, update.Credits[0].FinalUnitIndication.FinalAction)
	assert.Equal(t, restrictRule, update.Credits[0].FinalUnitIndication.RestrictRules[0])
}

func getClientConfig() *diameter.DiameterClientConfig {
	return &diameter.DiameterClientConfig{
		Host:        "test.test.com",
		Realm:       "test.com",
		ProductName: "gy_test",
		AppID:       diam.CHARGING_CONTROL_APP_ID,
	}
}

func getGyGlobalConfig(apnFilter, apnOverwrite string) *gy.GyGlobalConfig {
	rule := &credit_control.VirtualApnRule{}
	err := rule.FromMconfig(&mconfig.VirtualApnRule{ApnFilter: apnFilter, ApnOverwrite: apnOverwrite})
	if err != nil {
		return &gy.GyGlobalConfig{}
	}
	return &gy.GyGlobalConfig{
		VirtualApnRules: []*credit_control.VirtualApnRule{rule},
	}
}

func startServer(client *diameter.DiameterClientConfig, server *diameter.DiameterServerConfig, initMethod gy.InitMethod, finalUnitIndication mock_ocs.FinalUnitIndication) *mock_ocs.OCSDiamServer {
	serverStarted := make(chan struct{})
	var ocs *mock_ocs.OCSDiamServer
	go func() {
		log.Printf("Starting server")
		ocs = mock_ocs.NewOCSDiamServer(
			client,
			&mock_ocs.OCSConfig{
				MaxUsageOctets:      &fegprotos.Octets{TotalOctets: returnedOctets},
				MaxUsageTime:        1000,
				ValidityTime:        validityTime,
				ServerConfig:        server,
				GyInitMethod:        initMethod,
				FinalUnitIndication: finalUnitIndication,
			},
		)
		lis, err := ocs.StartListener()
		if err != nil {
			log.Fatalf("Could not start listener, %s", err.Error())
			return
		}
		server.Addr = lis.Addr().String()
		log.Printf("Server Addr: %v", server.Addr)
		serverStarted <- struct{}{}
		err = ocs.Start(lis)
		if err != nil {
			log.Fatalf("Could not start server, %s", err.Error())
			return
		}
	}()
	<-serverStarted
	time.Sleep(time.Millisecond)
	return ocs
}

func getReAuthHandler() gy.ChargingReAuthHandler {
	return func(request *gy.ChargingReAuthRequest) *gy.ChargingReAuthAnswer {
		return &gy.ChargingReAuthAnswer{
			SessionID:  request.SessionID,
			ResultCode: diam.Success,
		}
	}
}

func seedAccountConfigurations(ocs *mock_ocs.OCSDiamServer) {
	ctx := context.Background()
	ocs.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI1})
	ocs.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI2})
	ocs.SetCredit(
		ctx,
		&fegprotos.CreditInfo{
			Imsi:        testIMSI1,
			ChargingKey: 1,
			Volume:      &fegprotos.Octets{TotalOctets: 1000000},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)
	ocs.SetCredit(
		ctx,
		&fegprotos.CreditInfo{
			Imsi:        testIMSI1,
			ChargingKey: 2,
			Volume:      &fegprotos.Octets{TotalOctets: 1000000},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)
	ocs.SetCredit(
		ctx,
		&fegprotos.CreditInfo{
			Imsi:        testIMSI1,
			ChargingKey: 3,
			Volume:      &fegprotos.Octets{TotalOctets: 1000000},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)
	ocs.SetCredit(
		ctx,
		&fegprotos.CreditInfo{
			Imsi:        testIMSI2,
			ChargingKey: 1,
			Volume:      &fegprotos.Octets{TotalOctets: 1000000},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)
}

// assertReceivedAPNonOCS checks if the last received AVP contains the expected APN
func assertReceivedAPNonOCS(t *testing.T, ocs *mock_ocs.OCSDiamServer, expectedAPN string) {
	avpReceived, err := ocs.GetLastAVPreceived()
	assert.NoError(t, err)
	receivedAPN, err := avpReceived.FindAVP("Called-Station-Id", 0)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("UTF8String{%s},Padding:0", expectedAPN), receivedAPN.Data.String())
}
