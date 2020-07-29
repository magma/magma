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

package mock_pcrf_test

import (
	"context"
	"log"
	"testing"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/eap/test"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/testcore/pcrf/mock_pcrf"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

// TestGxClient tests CCR init and terminate messages using a fake PCRF
func TestPCRFExpectations(t *testing.T) {
	serverConfig := diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:     "127.0.0.1:0",
		Protocol: "tcp"},
	}
	clientConfig := getClientConfig()

	// E/A that should be met
	expectedInitReq := fegprotos.NewGxCCRequest(test.IMSI1, fegprotos.CCRequestType_INITIAL)
	usageMonitoringQuotaGrant := &fegprotos.UsageMonitoringInformation{
		MonitoringLevel: fegprotos.MonitoringLevel_RuleLevel,
		MonitoringKey:   []byte("mkey1"),
		Octets:          &fegprotos.Octets{TotalOctets: 1024},
	}
	dynamicRuleToInstall := &fegprotos.RuleDefinition{

		RuleName:         "rule1",
		RatingGroup:      9,
		Precedence:       10,
		MonitoringKey:    "m1",
		FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
		RedirectInformation: &lteprotos.RedirectInformation{
			Support:     lteprotos.RedirectInformation_ENABLED,
			AddressType: lteprotos.RedirectInformation_IPv4,
		},
		QosInformation: &lteprotos.FlowQos{
			MaxReqBwDl: 15,
			MaxReqBwUl: 30,
		},
	}
	activationTime := time.Now().Round(1 * time.Second)
	pActivationTime, err := ptypes.TimestampProto(activationTime)
	assert.NoError(t, err)
	deactivationTime := time.Now().Round(1 * time.Second).Add(5 * time.Second)
	pDeactivationTime, err := ptypes.TimestampProto(deactivationTime)
	assert.NoError(t, err)
	expectedInitAns := fegprotos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"rule1", "rule2"}, []string{"base1", "base2"}).
		SetDynamicRuleInstall(dynamicRuleToInstall).
		SetRuleActivationTime(pActivationTime).
		SetRuleDeactivationTime(pDeactivationTime).
		SetUsageMonitorInfo(usageMonitoringQuotaGrant)
	expectedInit := fegprotos.NewGxCreditControlExpectation().Expect(expectedInitReq).Return(expectedInitAns)

	// Update Request
	expectedUpdateReq := fegprotos.NewGxCCRequest(test.IMSI1, fegprotos.CCRequestType_UPDATE).
		SetUsageMonitorReport(usageMonitoringQuotaGrant).
		SetUsageReportDelta(100)
	expectedUpdateAns := fegprotos.NewGxCCAnswer(diam.Success).
		SetUsageMonitorInfo(usageMonitoringQuotaGrant)
	expectedUpdate := fegprotos.NewGxCreditControlExpectation().Expect(expectedUpdateReq).Return(expectedUpdateAns)

	// E/A that will not be met
	expectedReqNotMet := fegprotos.NewGxCCRequest(test.IMSI1, fegprotos.CCRequestType_UPDATE)
	answerNotMet := fegprotos.NewGxCCAnswer(diam.UnableToComply)
	expectationNotMet := fegprotos.NewGxCreditControlExpectation().Expect(expectedReqNotMet).Return(answerNotMet)

	defaultCCA := &fegprotos.GxCreditControlAnswer{
		ResultCode: 2001,
	}
	pcrf := startServerWithExpectations(
		clientConfig, &serverConfig,
		[]*fegprotos.GxCreditControlExpectation{expectedInit, expectedUpdate, expectationNotMet},
		fegprotos.UnexpectedRequestBehavior_CONTINUE_WITH_DEFAULT_ANSWER,
		defaultCCA)
	pcrf.CreateAccount(context.Background(), &lteprotos.SubscriberID{Id: test.IMSI1, Type: lteprotos.SubscriberID_IMSI})
	pcrf.CreateAccount(context.Background(), &lteprotos.SubscriberID{Id: test.IMSI2, Type: lteprotos.SubscriberID_IMSI})
	gxClient := gx.NewGxClient(clientConfig, &serverConfig, getMockReAuthHandler(), nil, nil)
	// send init
	ccrInit := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTInit,
		IMSI:          test.IMSI1,
		RequestNumber: 1,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
	}
	done := make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrInit))
	actualAnswer := gx.GetAnswer(done)
	assertCCAIsEqualToExpectedAnswer(t, actualAnswer, expectedInitAns)

	ccrUpdate := &gx.CreditControlRequest{
		SessionID:     "1",
		Type:          credit_control.CRTUpdate,
		IMSI:          test.IMSI1,
		RequestNumber: 2,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
		UsageReports: []*gx.UsageReport{
			{
				MonitoringKey: []byte("mkey1"),
				Level:         gx.RuleLevel,
				TotalOctets:   950,
			},
		},
	}
	done = make(chan interface{}, 1000)
	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrUpdate))
	actualAnswer = gx.GetAnswer(done)
	assertCCAIsEqualToExpectedAnswer(t, actualAnswer, expectedUpdateAns)

	// send an unexpected request
	ccrUpdateUnexpected := &gx.CreditControlRequest{
		SessionID:     "2",
		Type:          credit_control.CRTTerminate,
		IMSI:          test.IMSI2,
		RequestNumber: 3,
		IPAddr:        "192.168.1.1",
		SpgwIPV4:      "10.10.10.10",
	}
	done = make(chan interface{}, 1000)

	assert.NoError(t, gxClient.SendCreditControlRequest(&serverConfig, done, ccrUpdateUnexpected))
	unexpectedAnswer := gx.GetAnswer(done)
	assertCCAIsEqualToExpectedAnswer(t, unexpectedAnswer, defaultCCA)

	// should complain about an unexpected request
	res, err := pcrf.AssertExpectations(context.Background(), &orcprotos.Void{})
	assert.Nil(t, err)

	expectedResult := []*fegprotos.ExpectationResult{
		{ExpectationMet: true, ExpectationIndex: 0},
		{ExpectationMet: true, ExpectationIndex: 1},
		{ExpectationMet: false, ExpectationIndex: 2},
	}
	assert.ElementsMatch(t, expectedResult, res.Results)
	expectedErrors := []*fegprotos.ErrorByIndex{
		{
			Index: 2,
			Error: "Expected: Imsi: 001010000000055, Type: UPDATE, " +
				"Received: Imsi: 001010000000043, Type: TERMINATION",
		},
	}
	assert.ElementsMatch(t, expectedErrors, res.Errors)
}

func startServerWithExpectations(
	client *diameter.DiameterClientConfig,
	server *diameter.DiameterServerConfig,
	expectations []*fegprotos.GxCreditControlExpectation,
	failureBehavior fegprotos.UnexpectedRequestBehavior,
	defaultCCA *fegprotos.GxCreditControlAnswer,
) *mock_pcrf.PCRFServer {
	serverStarted := make(chan struct{})
	pcrf := mock_pcrf.NewPCRFServer(client, server)
	go func() {
		log.Printf("Starting server")
		ctx := context.Background()
		pcrf.SetPCRFConfigs(ctx, &fegprotos.PCRFConfigs{UseMockDriver: true})
		pcrf.SetExpectations(ctx, &fegprotos.GxCreditControlExpectations{
			Expectations:              expectations,
			UnexpectedRequestBehavior: failureBehavior,
			GxDefaultCca:              defaultCCA,
		})

		lis, err := pcrf.StartListener()
		if err != nil {
			log.Fatalf("Could not start listener for PCRF, %s", err.Error())
		}
		server.Addr = lis.Addr().String()
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

func assertCCAIsEqualToExpectedAnswer(t *testing.T, actual *gx.CreditControlAnswer, expectation *fegprotos.GxCreditControlAnswer) {
	ruleNames, ruleBaseNames, ruleDefinitions := getRuleInstallsFromCCA(actual)
	assert.ElementsMatch(t, expectation.GetRuleInstalls().GetRuleNames(), ruleNames)
	assert.ElementsMatch(t, expectation.GetRuleInstalls().GetRuleBaseNames(), ruleBaseNames)
	assert.ElementsMatch(t, expectation.GetRuleInstalls().GetRuleDefinitions(), ruleDefinitions)
	assertRuleInstallTimeStampsMatch(t, expectation.GetRuleInstalls(), actual.RuleInstallAVP)
	usageMonitors := getUsageMonitorsFromCCA(actual)
	assert.ElementsMatch(t, expectation.GetUsageMonitoringInfos(), usageMonitors)
}

func getRuleInstallsFromCCA(cca *gx.CreditControlAnswer) ([]string, []string, []*fegprotos.RuleDefinition) {
	var ruleNames []string
	var ruleBaseNames []string
	var ruleDefinitions []*fegprotos.RuleDefinition
	for _, installRule := range cca.RuleInstallAVP {
		ruleNames = append(ruleNames, installRule.RuleNames...)
		ruleBaseNames = append(ruleBaseNames, installRule.RuleBaseNames...)
		ruleDefinitions = append(ruleDefinitions, toProtosRuleDefinitions(installRule.RuleDefinitions)...)
	}
	return ruleNames, ruleBaseNames, ruleDefinitions
}

func assertRuleInstallTimeStampsMatch(t *testing.T, expected *fegprotos.RuleInstalls, actual []*gx.RuleInstallAVP) {
	expectedActivationTime, _ := ptypes.Timestamp(expected.GetActivationTime())
	expectedDeactivationTime, _ := ptypes.Timestamp(expected.GetDeactivationTime())

	for _, ruleInstall := range actual {
		if expected.GetActivationTime() != nil {
			assert.True(t, expectedActivationTime.Equal(*ruleInstall.RuleActivationTime))
		}
		if expected.GetDeactivationTime() != nil {
			assert.True(t, expectedDeactivationTime.Equal(*ruleInstall.RuleDeactivationTime))
		}
	}
}

func toProtosRuleDefinitions(gxRuleDfs []*gx.RuleDefinition) []*fegprotos.RuleDefinition {
	ruleDefs := []*fegprotos.RuleDefinition{}
	for _, ruleDef := range gxRuleDfs {
		ruleDefs = append(ruleDefs, &fegprotos.RuleDefinition{
			RuleName:            ruleDef.RuleName,
			RatingGroup:         swag.Uint32Value(ruleDef.RatingGroup),
			Precedence:          ruleDef.Precedence,
			MonitoringKey:       string(ruleDef.MonitoringKey),
			FlowDescriptions:    ruleDef.FlowDescriptions,
			RedirectInformation: ruleDef.RedirectInformation.ToProto(),
			QosInformation:      ruleDef.Qos.ToProto(),
		})
	}
	return ruleDefs
}

func getUsageMonitorsFromCCA(cca *gx.CreditControlAnswer) []*fegprotos.UsageMonitoringInformation {
	monitors := []*fegprotos.UsageMonitoringInformation{}
	for _, usageMonitor := range cca.UsageMonitors {
		monitors = append(monitors, &fegprotos.UsageMonitoringInformation{
			MonitoringKey:   usageMonitor.MonitoringKey,
			MonitoringLevel: fegprotos.MonitoringLevel(usageMonitor.Level),
			Octets:          grantedServiceUnitToOctet(usageMonitor.GrantedServiceUnit),
		})
	}
	return monitors
}

func grantedServiceUnitToOctet(gsu *credit_control.GrantedServiceUnit) *fegprotos.Octets {
	return &fegprotos.Octets{
		TotalOctets:  swag.Uint64Value(gsu.TotalOctets),
		InputOctets:  swag.Uint64Value(gsu.InputOctets),
		OutputOctets: swag.Uint64Value(gsu.OutputOctets),
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

func getMockReAuthHandler() gx.PolicyReAuthHandler {
	return func(request *gx.PolicyReAuthRequest) *gx.PolicyReAuthAnswer {
		return &gx.PolicyReAuthAnswer{
			SessionID:  request.SessionID,
			ResultCode: diam.Success,
		}
	}
}
