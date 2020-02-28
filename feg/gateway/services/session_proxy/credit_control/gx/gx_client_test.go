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
	imsi1Rules     = []string{"rule1", "rule2"}
	imsi1BaseRules = []string{"rule1", "rule2"}
	imsi2Rules     = []string{"rule1", "rule3"}
	imsi3Rules     = []string{"rule1", "rule2"}
	imsi3BaseRules = []string{"rule1", "rule2"}
)

// TestGxClient tests CCR init and terminate messages using a fake PCRF
func TestGxClient(t *testing.T) {
	serverConfig := diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:     "127.0.0.1:3898",
		Protocol: "tcp"},
	}
	serverConfig1 := serverConfig
	serverConfig2 := serverConfig
	clientConfig := getClientConfig()
	startServer(clientConfig, &serverConfig1)

	gxClient := gx.NewGxClient(
		clientConfig,
		&serverConfig2,
		getMockReAuthHandler(),
		nil,
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

func TestGxClientUsageMonitoring(t *testing.T) {
	serverConfig := diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:     "127.0.0.1:3899",
		Protocol: "tcp"},
	}
	serverConfig1 := serverConfig
	serverConfig2 := serverConfig
	clientConfig := getClientConfig()
	startServer(clientConfig, &serverConfig1)
	gxClient := gx.NewGxClient(
		clientConfig,
		&serverConfig2,
		getMockReAuthHandler(),
		nil,
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

func getClientConfig() *diameter.DiameterClientConfig {
	return &diameter.DiameterClientConfig{
		Host:        "test.test.com",
		Realm:       "test.com",
		ProductName: "gx_test",
		AppID:       diam.GX_CHARGING_CONTROL_APP_ID,
	}
}

func startServer(
	client *diameter.DiameterClientConfig,
	server *diameter.DiameterServerConfig,
) {
	serverStarted := make(chan struct{})
	go func() {
		log.Printf("Starting server")
		pcrf := mock_pcrf.NewPCRFDiamServer(
			client,
			&mock_pcrf.PCRFConfig{ServerConfig: server},
		)
		ctx := context.Background()
		pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI1})
		pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI2})
		pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI3})
		pcrf.CreateAccount(ctx, &protos.SubscriberID{Id: testIMSI4})
		monitoringKey := "mkey"
		monitoringKey3 := "mkey3"
		var rg1 uint32 = 1
		var rg3 uint32 = 3
		var rg5 uint32 = 5

		redirect := &protos.RedirectInformation{
			Support:       protos.RedirectInformation_ENABLED,
			AddressType:   protos.RedirectInformation_URL,
			ServerAddress: "http://www.example.com/",
		}

		maxReqBWUL := uint32(128000)
		maxReqBWDL := uint32(128000)
		gbrDL := uint32(64000)
		gbrUL := uint32(64000)
		qci := protos.FlowQos_Qci(8)

		qos := &protos.FlowQos{
			MaxReqBwUl: maxReqBWUL,
			MaxReqBwDl: maxReqBWDL,
			GbrDl:      gbrDL,
			GbrUl:      gbrUL,
			Qci:        qci,
		}

		pcrf.SetRules(
			ctx,
			&fegprotos.AccountRules{
				Imsi:          testIMSI1,
				RuleNames:     imsi1Rules,
				RuleBaseNames: imsi1BaseRules,
				RuleDefinitions: []*fegprotos.RuleDefinition{
					{
						RuleName:       "dynrule1",
						RatingGroup:    rg1,
						Precedence:     100,
						MonitoringKey:  monitoringKey,
						QosInformation: qos,
					},
				},
			},
		)
		pcrf.SetRules(
			ctx,
			&fegprotos.AccountRules{
				Imsi:      testIMSI2,
				RuleNames: imsi2Rules,
			},
		)
		pcrf.SetRules(
			ctx,
			&fegprotos.AccountRules{
				Imsi:          testIMSI3,
				RuleNames:     imsi3Rules,
				RuleBaseNames: imsi3BaseRules,
				RuleDefinitions: []*fegprotos.RuleDefinition{
					{
						RuleName:            "dynrule3",
						RatingGroup:         rg3,
						Precedence:          300,
						MonitoringKey:       monitoringKey3,
						RedirectInformation: redirect,
					},
				},
			},
		)
		pcrf.SetRules(
			ctx,
			&fegprotos.AccountRules{
				Imsi: testIMSI4,
				RuleDefinitions: []*fegprotos.RuleDefinition{
					{
						RuleName:      "dynrule4",
						Precedence:    300,
						MonitoringKey: monitoringKey,
					},
					{
						RuleName:    "dynrule5",
						RatingGroup: rg5,
						Precedence:  100,
					},
				},
			},
		)
		pcrf.SetUsageMonitors(
			ctx,
			&fegprotos.SetUsageMonitorRequest{
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
			},
		)
		serverStarted <- struct{}{}
		lis, err := pcrf.StartListener()
		if err != nil {
			log.Fatalf("Could not start listener for PCRF, %s", err.Error())
		}
		server.Addr = lis.Addr().String()
		err = pcrf.Start(lis)
		if err != nil {
			log.Fatalf("Could not start test PCRF server, %s", err.Error())
			return
		}
	}()
	<-serverStarted
	time.Sleep(time.Millisecond)
}

func getMockReAuthHandler() gx.ReAuthHandler {
	return func(request *gx.ReAuthRequest) *gx.ReAuthAnswer {
		return &gx.ReAuthAnswer{
			SessionID:  request.SessionID,
			ResultCode: diam.Success,
		}
	}
}
