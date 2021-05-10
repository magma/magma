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

// Package gx_test tests gx protocol messages
package gx_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/testcore/pcrf/mock_pcrf"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/stretchr/testify/assert"
)

const (
	testIMSI1   = "1234"
	testIMSI2   = "4321"
	testIMSI3   = "4499"
	testIMSI4   = "5000"
	HOUR_IN_MIN = 60
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
	globalConfig := getGxGlobalConfig("", "", "")
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
		IPv6Addr:      "2001:0db8:0a0b:12f0:0000:0000:0000:FFFF",
		SpgwIPV4:      "10.10.10.10",
		Apn:           "gx.Apn.magma.com",
	}
	done := make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assert.Equal(t, 5, len(answer.RuleInstallAVP))
	assertReceivedAPNonPCRF(t, pcrf, ccrInit.Apn)
	assertReceivedIPv4onPCRF(t, pcrf, ccrInit.IPAddr)
	assertReceivedIPv6onPCRF(t, pcrf, ccrInit.IPv6Addr)

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
	matchApn := ".*\\.magma.*"
	matchCC := "12"
	overwriteApn := "gx.overwritten.Apn.magma.com"
	globalConfig := getGxGlobalConfig(matchApn, matchCC, overwriteApn)
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
		SessionID:               "1",
		Type:                    credit_control.CRTInit,
		IMSI:                    testIMSI1,
		RequestNumber:           0,
		IPAddr:                  "192.168.1.1",
		SpgwIPV4:                "10.10.10.10",
		Apn:                     "gx.Apn.magma.com",
		ChargingCharacteristics: "12",
	}
	done := make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assert.Equal(t, 5, len(answer.RuleInstallAVP))
	assertReceivedAPNonPCRF(t, pcrf, overwriteApn)
}

// Test VirtualAPN configuration when one of the APN/ChargingCharacteristics
// RegEx is not satisfied. (We expect the APN to be not modified in this case)
func TestGxClientVirtualAPNNoMatch(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	matchApn := ".*\\.magma.*"
	matchCC := "13"
	overwriteApn := "gx.overwritten.Apn.magma.com"
	globalConfig := getGxGlobalConfig(matchApn, matchCC, overwriteApn)
	pcrf := startServer(clientConfig, &serverConfig)
	seedAccountConfigurations(pcrf)

	gxClient := gx.NewGxClient(
		clientConfig,
		&serverConfig,
		getMockReAuthHandler(),
		nil,
		globalConfig,
	)

	// 1. First fail the charging characteristics regex
	originalAPN := "gx.Apn.magma.com"
	ccrInit := &gx.CreditControlRequest{
		SessionID:               "1",
		Type:                    credit_control.CRTInit,
		IMSI:                    testIMSI1,
		RequestNumber:           0,
		IPAddr:                  "192.168.1.1",
		SpgwIPV4:                "10.10.10.10",
		Apn:                     originalAPN,
		ChargingCharacteristics: "12",
	}
	done := make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer := gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assertReceivedAPNonPCRF(t, pcrf, originalAPN)

	// 2. Now fail the APN regex
	originalAPN = "gx.Apn.m-a-g-m-a.com"
	ccrInit.Apn = originalAPN
	ccrInit.ChargingCharacteristics = "13"
	done = make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	answer = gx.GetAnswer(done)
	assert.Equal(t, ccrInit.SessionID, answer.SessionID)
	assert.Equal(t, ccrInit.RequestNumber, answer.RequestNumber)
	assertReceivedAPNonPCRF(t, pcrf, originalAPN)
}

func TestGxClientUsageMonitoring(t *testing.T) {
	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	globalConfig := getGxGlobalConfig("", "", "")
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
	globalConfig := getGxGlobalConfig("", "", "")
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

// Test cases that don't support FramedIpv6 AVPs
// When these env variables are set, the FrameIP Address should be overwritten
// with the value, and no FramedIPv6Prefix AVP should be sent
func TestDefaultFramedIpv4Addr(t *testing.T) {
	var defaultIpv4 = "10.10.10.11"
	log.Printf("Start TestDefaultFramedIpv4Addr with default=%v", defaultIpv4)
	os.Setenv(gx.FramedIPv4AddrRequiredEnv, "1")
	os.Setenv(gx.DefaultFramedIPv4AddrEnv, defaultIpv4)
	defer func() {
		os.Unsetenv(gx.FramedIPv4AddrRequiredEnv)
		os.Unsetenv(gx.DefaultFramedIPv4AddrEnv)
	}()

	serverConfig := defaultLocalServerConfig
	clientConfig := getClientConfig()
	globalConfig := getGxGlobalConfig("", "", "")
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
		IPAddr:        "2001:db8:0:1:1:1:1:1",
		SpgwIPV4:      "10.10.10.10",
	}
	done := make(chan interface{}, 1000)
	log.Printf("Sending CCR-Init")
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	gx.GetAnswer(done)

	lastMsg, err := pcrf.GetLastAVPreceived()
	assert.NoError(t, err)
	avpValue, err := lastMsg.FindAVP(avp.FramedIPAddress, 0)
	assert.NoError(t, err)
	actualIPv4, err := datatype.DecodeIPv4(avpValue.Data.Serialize())
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("IPv4{%v}", defaultIpv4), actualIPv4.String())

	_, err = lastMsg.FindAVP(avp.FramedIPv6Prefix, 0)
	assert.EqualError(t, err, "AVP not found")
}

func TestTimezoneConversion(t *testing.T) {
	// Test behind UTC (negative offset)
	pTimezone := &protos.Timezone{OffsetMinutes: -6 * HOUR_IN_MIN}
	convertedTimezone := gx.GetTimezoneByte(pTimezone)
	assert.Equal(t, byte(0x4a), convertedTimezone)

	pTimezone = &protos.Timezone{OffsetMinutes: -8 * HOUR_IN_MIN}
	convertedTimezone = gx.GetTimezoneByte(pTimezone)
	assert.Equal(t, byte(0x2b), convertedTimezone)

	pTimezone = &protos.Timezone{OffsetMinutes: -7 * HOUR_IN_MIN}
	convertedTimezone = gx.GetTimezoneByte(pTimezone)
	assert.Equal(t, byte(0x8a), convertedTimezone)

	// Test ahead UTC (positive offset)
	pTimezone = &protos.Timezone{OffsetMinutes: 1 * HOUR_IN_MIN}
	convertedTimezone = gx.GetTimezoneByte(pTimezone)
	assert.Equal(t, byte(0x40), convertedTimezone)
}

func getClientConfig() *diameter.DiameterClientConfig {
	return &diameter.DiameterClientConfig{
		Host:        "test.test.com",
		Realm:       "test.com",
		ProductName: "gx_test",
		AppID:       diam.GX_CHARGING_CONTROL_APP_ID,
	}
}

func getGxGlobalConfig(apnFilter, chargingCharacteristicsFilter, apnOverwrite string) *gx.GxGlobalConfig {
	rule := &credit_control.VirtualApnRule{}
	mconfigRule := &mconfig.VirtualApnRule{
		ApnFilter:                     apnFilter,
		ChargingCharacteristicsFilter: chargingCharacteristicsFilter,
		ApnOverwrite:                  apnOverwrite,
	}
	err := rule.FromMconfig(mconfigRule)
	if err != nil {
		return &gx.GxGlobalConfig{}
	}
	return &gx.GxGlobalConfig{
		VirtualApnRules: []*credit_control.VirtualApnRule{rule},
	}
}

func startServer(
	client *diameter.DiameterClientConfig,
	server *diameter.DiameterServerConfig,
) *mock_pcrf.PCRFServer {
	serverStarted := make(chan struct{})
	var pcrf *mock_pcrf.PCRFServer
	go func() {
		log.Printf("Starting server")
		pcrf = mock_pcrf.NewPCRFServer(client, server)

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

func seedAccountConfigurations(pcrf *mock_pcrf.PCRFServer) {
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

// assertReceivedAPNonPCRF checks if the last received AVP contains the expected APN
func assertReceivedAPNonPCRF(t *testing.T, pcrf *mock_pcrf.PCRFServer, expectedAPN string) {
	avpReceived, err := pcrf.GetLastAVPreceived()
	assert.NoError(t, err)
	receivedAPN, err := avpReceived.FindAVP("Called-Station-Id", 0)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("UTF8String{%s},Padding:0", expectedAPN), receivedAPN.Data.String())
}

func assertReceivedIPv4onPCRF(t *testing.T, pcrf *mock_pcrf.PCRFServer, expectedIPv4 string) {
	// convert ip string into ip AVP (octetstring)
	ipv4 := net.ParseIP(expectedIPv4).To4()
	ipv4OctetString := datatype.OctetString([]byte(ipv4))
	expectedIPv4avp := diam.NewAVP(avp.FramedIPAddress, avp.Mbit, 0, ipv4OctetString)

	avpReceived, err := pcrf.GetLastAVPreceived()
	assert.NoError(t, err)
	receivedFramedIPAddress, err := avpReceived.FindAVP(avp.FramedIPAddress, 0)
	assert.NoError(t, err)

	assert.Equal(t, expectedIPv4avp.Data, receivedFramedIPAddress.Data)
}

func assertReceivedIPv6onPCRF(t *testing.T, pcrf *mock_pcrf.PCRFServer, expectedIPv6 string) {
	// convert ip string into ip AVP (octetstring)
	ipv6 := net.ParseIP(expectedIPv6).To16()
	ipv6OctetString := datatype.OctetString([]byte(ipv6))
	expectedIPv6avp := diam.NewAVP(avp.FramedIPv6Prefix, avp.Mbit, 0, ipv6OctetString[0:8])

	avpReceived, err := pcrf.GetLastAVPreceived()
	assert.NoError(t, err)
	receivedFramedIPv6Prefix, err := avpReceived.FindAVP(avp.FramedIPv6Prefix, 0)
	assert.NoError(t, err)

	assert.Equal(t, expectedIPv6avp.Data, receivedFramedIPv6Prefix.Data)
}
