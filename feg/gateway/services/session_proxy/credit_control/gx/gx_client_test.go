/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package gx_test tests gx protocol messages
package gx_test

import (
	"log"
	"net"
	"testing"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/testcore/pcrf/mock_pcrf"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	testIMSI1 = "1234"
	testIMSI2 = "4321"
	testIMSI3 = "4499"
	testIMSI4 = "5000"
)

var (
	imsi1Rules               = []string{"rule1", "rule2"}
	imsi1BaseRules           = []string{"rule1", "rule2"}
	imsi2Rules               = []string{"rule1", "rule3"}
	imsi3Rules               = []string{"rule1", "rule2"}
	imsi3BaseRules           = []string{"rule1", "rule2"}
	defaultLocalServerConfig = diameter.DiameterServerConfig{
		DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:     "127.0.0.1:0",
			Protocol: "tcp"},
	}
)

// TestGxClient tests CCR init and terminate messages using a fake PCRF
func TestGxClient(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	globalConfig := getGxGlobalConfig("")
	pcrf := startServer(clientConfig, &serverConfig)
	seedAccountConfigurations(pcrf)

	gxClient := gx.NewGxClient(
		clientConfig,
		&serverConfig,
		getMockReAuthHandler(),
		nil,
		globalConfig,
	)

	// send init
	ccrInit := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
	}
	done := make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assert.Equal(t, 5, len(answer.RuleInstallAVP))
	calledStationID, err := mock_pcrf.GetAVP(pcrf.LastMessageReceived, "Called-Station-Id")
	assert.NoError(t, err)
	assert.Equal(t, "", calledStationID)

	var ruleNames []string
	var ruleBaseNames []string
	var ruleDefinitions []*gx.RuleDefinition
	for _, installRule := range answer.RuleInstallAVP {
		ruleNames = append(ruleNames, installRule.RuleNames...)
		ruleBaseNames = append(ruleBaseNames, installRule.RuleBaseNames...)
		ruleDefinitions = append(ruleDefinitions, installRule.RuleDefinitions...)
	}
	assert.ElementsMatch(t, imsi1Rules, ruleNames)
	assert.ElementsMatch(t, imsi1BaseRules, ruleBaseNames)
	assert.Equal(t, 1, len(ruleDefinitions))
	assert.Equal(t, "dynrule1", ruleDefinitions[0].RuleName)
	assert.Equal(t, "mkey", string(ruleDefinitions[0].MonitoringKey))
	assert.Equal(t, uint32(128000), *ruleDefinitions[0].Qos.MaxReqBwUL)
	assert.Equal(t, uint32(128000), *ruleDefinitions[0].Qos.MaxReqBwDL)
	if ruleDefinitions[0].Qos.GbrUL != nil {
		assert.Equal(t, uint32(64000), *ruleDefinitions[0].Qos.GbrUL)
	}
	if ruleDefinitions[0].Qos.GbrDL != nil {
		assert.Equal(t, uint32(64000), *ruleDefinitions[0].Qos.GbrDL)
	}
	if ruleDefinitions[0].Qos.Qci != nil {
		assert.Equal(t, int32(8), int32(*ruleDefinitions[0].Qos.Qci))
	}

	// send terminate
	ccrTerminate := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTTerminate,
		IMSI:          testIMSI1,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
	}
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrTerminate))
	terminate := gx.GetAnswer(done)
	assert.Equal(t, ccrTerminate.SessionID, terminate.SessionID)
	assert.Equal(t, ccrTerminate.RequestNumber, terminate.RequestNumber)
	assert.Empty(t, terminate.RuleInstallAVP)

	// send init
	ccrInit = &gx.CreditControlRequest{
		SessionID:     "2",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI3,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
	}
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer = gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assert.Equal(t, 5, len(answer.RuleInstallAVP))

	ruleNames = []string{}
	ruleBaseNames = []string{}
	ruleDefinitions = []*gx.RuleDefinition{}
	for _, installRule := range answer.RuleInstallAVP {
		ruleNames = append(ruleNames, installRule.RuleNames...)
		ruleBaseNames = append(ruleBaseNames, installRule.RuleBaseNames...)
		ruleDefinitions = append(ruleDefinitions, installRule.RuleDefinitions...)
	}
	assert.ElementsMatch(t, imsi3Rules, ruleNames)
	assert.ElementsMatch(t, imsi3BaseRules, ruleBaseNames)
	assert.Equal(t, 1, len(ruleDefinitions))
	assert.Equal(t, "dynrule3", ruleDefinitions[0].RuleName)
	assert.Equal(t, "mkey3", string(ruleDefinitions[0].MonitoringKey))
	assert.Equal(t, uint32(1), ruleDefinitions[0].RedirectInformation.RedirectSupport)
	assert.Equal(t, uint32(2), ruleDefinitions[0].RedirectInformation.RedirectAddressType)
	assert.Equal(t, "http://www.example.com/", ruleDefinitions[0].RedirectInformation.RedirectServerAddress)

	// send terminate
	ccrTerminate = &gx.CreditControlRequest{
		SessionID:     "2",
		Type:          credit_control.CRTTerminate,
		IMSI:          testIMSI3,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
	}
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrTerminate))
	terminate = gx.GetAnswer(done)
	assert.Equal(t, ccrTerminate.SessionID, terminate.SessionID)
	assert.Equal(t, ccrTerminate.RequestNumber, terminate.RequestNumber)
	assert.Empty(t, terminate.RuleInstallAVP)

	// Connection Disabling should cause CCR to fail
	gxClient.DisableConnections(10 * time.Second)
	assert.Error(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))

	// CCR Success after Enabling
	gxClient.EnableConnections()
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))

	hwaddr, err := net.ParseMAC("00:00:5e:00:53:01")
	assert.NoError(t, err)
	ipv6addr := gx.Ipv6PrefixFromMAC(hwaddr)
	assert.Equal(t, ipv6addr[:6], []byte{0, 0x80, 0xfd, 0xfa, 0xce, 0xb0})
	assert.NotEqual(t, ipv6addr[6:10], []byte{0x0c, 0xab, 0xcd, 0xef})
	assert.Equal(t, ipv6addr[10:], []byte{0x2, 0x0, 0x5e, 0xff, 0xfe, 0x0, 0x53, 0x1})

}

// TestGxClient tests CCR init and terminate messages using a fake PCRF and a specific GxGlobalConfig
func TestGxClientWithGyGlobalConf(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	overWriteApn := "gx.Apn.magma.com"
	globalConfig := getGxGlobalConfig(overWriteApn)
	pcrf := startServer(clientConfig, &serverConfig)
	seedAccountConfigurations(pcrf)

	gxClient := gx.NewGxClient(
		clientConfig,
		&serverConfig,
		getMockReAuthHandler(),
		nil,
		globalConfig,
	)

	// send init
	ccrInit := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI1,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
	}
	done := make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assert.Equal(t, 5, len(answer.RuleInstallAVP))
	calledStationID, err := mock_pcrf.GetAVP(pcrf.LastMessageReceived, "Called-Station-Id")
	assert.NoError(t, err)
	assert.Equal(t, overWriteApn, calledStationID)
}

func TestGxClientUsageMonitoring(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	globalConfig := getGxGlobalConfig("")
	pcrf := startServer(clientConfig, &serverConfig)
	seedAccountConfigurations(pcrf)

	gxClient := gx.NewGxClient(
		clientConfig,
		&serverConfig,
		getMockReAuthHandler(),
		nil,
		globalConfig,
	)
	done := make(chan interface{}, 1000)

	// Usage Monitoring
	init := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI4,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
	}

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, init))
	initAnswer := gx.GetAnswer(done)
	assert.Equal(t, init.SessionID, initAnswer.SessionID)
	assert.Equal(t, init.RequestNumber, initAnswer.RequestNumber)
	assert.Equal(t, 2, len(initAnswer.UsageMonitors))
	for _, monitor := range initAnswer.UsageMonitors {
		if string(monitor.MonitoringKey) == "mkey" {
			assert.Equal(t, *(*monitor.GrantedServiceUnit).TotalOctets, uint64(1024))
			assert.Equal(t, monitor.Level, gx.RuleLevel)
		} else if string(monitor.MonitoringKey) == "mkey3" {
			assert.Equal(t, *(*monitor.GrantedServiceUnit).TotalOctets, uint64(2048))
			assert.Equal(t, monitor.Level, gx.SessionLevel)
		} else {
			assert.True(t, false)
		}
	}

	update := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTUpdate,
		IMSI:          testIMSI4,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
		UsageReports: []*gx.UsageReport{
			{
				MonitoringKey: []byte("mkey"),
				Level:         gx.RuleLevel,
				InputOctets:   24,
				OutputOctets:  0,
				TotalOctets:   24,
			},
			{
				MonitoringKey: []byte("mkey3"),
				Level:         gx.SessionLevel,
				InputOctets:   4000,
				OutputOctets:  0,
				TotalOctets:   4000,
			},
		},
	}
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, update))
	updateAnswer := gx.GetAnswer(done)
	assert.Equal(t, update.SessionID, updateAnswer.SessionID)
	assert.Equal(t, update.RequestNumber, updateAnswer.RequestNumber)
	assert.Equal(t, 2, len(updateAnswer.UsageMonitors))
	for _, monitor := range updateAnswer.UsageMonitors {
		if string(monitor.MonitoringKey) == "mkey" {
			assert.Equal(t, *(*monitor.GrantedServiceUnit).TotalOctets, uint64(1024))
			assert.Equal(t, monitor.Level, gx.RuleLevel)
		} else if string(monitor.MonitoringKey) == "mkey3" {
			assert.Equal(t, *(*monitor.GrantedServiceUnit).TotalOctets, uint64(96))
			assert.Equal(t, monitor.Level, gx.SessionLevel)
		} else {
			assert.Fail(t, "Unknown monitoring key")
		}
	}
}

func TestGxReAuthRemoveRules(t *testing.T) {
	log.Printf("Start TestGxReAuth")
	if t != nil {
		return
	}
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	globalConfig := getGxGlobalConfig("")
	pcrf := startServer(clientConfig, &serverConfig)
	seedAccountConfigurations(pcrf)

	gxClient := gx.NewGxClient(
		clientConfig,
		&serverConfig,
		getMockReAuthHandler(),
		nil,
		globalConfig,
	)

	// send one init to set user context in OCS
	ccrInit := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          testIMSI4,
		RequestNumber: 0,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
	}
	done := make(chan interface{}, 1000)
	log.Printf("Sending CCR-Init")
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)

	// send a RAR request to request a rule removal
	rulesToRemove := &fegprotos.RuleRemovals{RuleNames: []string{"dynrule4"}, RuleBaseNames: []string{""}}
	raa, err := pcrf.ReAuth(
		context.Background(),
		&fegprotos.PolicyReAuthTarget{Imsi: testIMSI4, RulesToRemove: rulesToRemove},
	)
	// success reauth
	assert.NoError(t, err)
	assert.Equal(t, "1", raa.SessionId)
	assert.Equal(t, uint32(diam.Success), raa.ResultCode)
}

func getClientConfig() *diameter.DiameterClientConfig {
	return &diameter.DiameterClientConfig{
		Host:        "test.test.com",
		Realm:       "test.com",
		ProductName: "gx_test",
		AppID:       diam.GX_CHARGING_CONTROL_APP_ID,
	}
}

func getGxGlobalConfig(pcrfOverwriteApn string) *gx.GxGlobalConfig {
	return &gx.GxGlobalConfig{
		PCFROverwriteApn: pcrfOverwriteApn,
	}
}

func startServer(
	client *diameter.DiameterClientConfig,
	server *diameter.DiameterServerConfig,
) *mock_pcrf.PCRFDiamServer {
	serverStarted := make(chan struct{})
	var pcrf *mock_pcrf.PCRFDiamServer
	go func() {
		log.Printf("Starting server")
		pcrf = mock_pcrf.NewPCRFDiamServer(
			client,
			&mock_pcrf.PCRFConfig{ServerConfig: server},
		)

		lis, err := pcrf.StartListener()
		if err != nil {
			log.Fatalf("Could not start listener for PCRF, %s", err.Error())
		}
		// Overwrite config addr with the allocated port
		server.Addr = lis.Addr().String()
		log.Printf("Server Addr: %v", server.Addr)
		serverStarted <- struct{}{}

		err = pcrf.Start(lis)
		if err != nil {
			log.Fatalf("Could not start test PCRF server, %s", err.Error())
			return
		}
	}()
	<-serverStarted
	return pcrf
}

func getMockReAuthHandler() gx.PolicyReAuthHandler {
	return func(request *gx.PolicyReAuthRequest) *gx.PolicyReAuthAnswer {
		return &gx.PolicyReAuthAnswer{
			SessionID:  request.SessionID,
			ResultCode: diam.Success,
		}
	}
}

func seedAccountConfigurations(pcrf *mock_pcrf.PCRFDiamServer) {
	monitoringKey := "mkey"
	monitoringKey3 := "mkey3"
	ruleImsi1 := &fegprotos.AccountRules{
		Imsi:                testIMSI1,
		StaticRuleNames:     imsi1Rules,
		StaticRuleBaseNames: imsi1BaseRules,
		DynamicRuleDefinitions: []*fegprotos.RuleDefinition{
			{
				RuleName:      "dynrule1",
				RatingGroup:   1,
				Precedence:    100,
				MonitoringKey: monitoringKey,
				QosInformation: &protos.FlowQos{
					MaxReqBwUl: 128000,
					MaxReqBwDl: 128000,
					GbrDl:      64000,
					GbrUl:      64000,
					Qci:        8,
				},
			},
		},
	}
	ruleImsi2 := &fegprotos.AccountRules{
		Imsi:            testIMSI2,
		StaticRuleNames: imsi2Rules,
	}
	ruleImsi3 := &fegprotos.AccountRules{
		Imsi:                testIMSI3,
		StaticRuleNames:     imsi3Rules,
		StaticRuleBaseNames: imsi3BaseRules,
		DynamicRuleDefinitions: []*fegprotos.RuleDefinition{
			{
				RuleName:      "dynrule3",
				RatingGroup:   3,
				Precedence:    300,
				MonitoringKey: monitoringKey3,
				RedirectInformation: &protos.RedirectInformation{
					Support:       protos.RedirectInformation_ENABLED,
					AddressType:   protos.RedirectInformation_URL,
					ServerAddress: "http://www.example.com/",
				},
			},
		},
	}
	ruleImsi4 := &fegprotos.AccountRules{
		Imsi: testIMSI4,
		DynamicRuleDefinitions: []*fegprotos.RuleDefinition{
			{
				RuleName:      "dynrule4",
				Precedence:    300,
				MonitoringKey: monitoringKey,
			},
			{
				RuleName:    "dynrule5",
				RatingGroup: 5,
				Precedence:  100,
			},
		},
	}
	usageMonitorImsi4 := &fegprotos.UsageMonitorConfiguration{
		Imsi: testIMSI4,
		UsageMonitorCredits: []*fegprotos.UsageMonitor{
			{
				MonitorInfoPerRequest: &fegprotos.UsageMonitoringInformation{
					MonitoringKey:   []byte(monitoringKey),
					MonitoringLevel: fegprotos.MonitoringLevel_RuleLevel,
					Octets:          &fegprotos.Octets{TotalOctets: 1024},
				},
				TotalQuota: &fegprotos.Octets{TotalOctets: 4096},
			},
			{
				MonitorInfoPerRequest: &fegprotos.UsageMonitoringInformation{
					MonitoringKey:   []byte(monitoringKey3),
					MonitoringLevel: fegprotos.MonitoringLevel_SessionLevel,
					Octets:          &fegprotos.Octets{TotalOctets: 2048},
				},
				TotalQuota: &fegprotos.Octets{TotalOctets: 4096},
			},
		},
	}
	ctx := context.Background()
	// IMSI 1
	pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI1})
	pcrf.SetRules(ctx, ruleImsi1)
	// IMSI 2
	pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI2})
	pcrf.SetRules(ctx, ruleImsi2)
	// IMSI 3
	pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI3})
	pcrf.SetRules(ctx, ruleImsi3)
	// IMSI 4
	pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI4})
	pcrf.SetRules(ctx, ruleImsi4)
	pcrf.SetUsageMonitors(ctx, usageMonitorImsi4)
}
