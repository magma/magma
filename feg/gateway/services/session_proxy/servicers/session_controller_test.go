/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"fmt"
	"testing"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/session_proxy/servicers"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/gateway/mconfig"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thoas/go-funk"
	"golang.org/x/net/context"
)

const (
	IMSI1 = "IMSI00101"
	IMSI2 = "IMSI00102"
)

type MockPolicyClient struct {
	mock.Mock
}

func (p *MockPolicyClient) SendCreditControlRequest(
	server *diameter.DiameterServerConfig,
	done chan interface{},
	request *gx.CreditControlRequest,
) error {
	args := p.Called(server, done, request)
	return args.Error(0)
}

func (p *MockPolicyClient) IgnoreAnswer(request *gx.CreditControlRequest) {
	return
}

func (p *MockPolicyClient) EnableConnections() error {
	p.Called()
	return nil
}

func (p *MockPolicyClient) DisableConnections(period time.Duration) {
	p.Called(period)
	return
}

type MockPolicyDBClient struct {
	mock.Mock
}

func (client *MockPolicyDBClient) GetChargingKeysForRules(ruleIDs []string, ruleDefs []*protos.PolicyRule) []policydb.ChargingKey {

	args := client.Called(ruleIDs)
	return args.Get(0).([]policydb.ChargingKey)
}

func (client *MockPolicyDBClient) GetRuleIDsForBaseNames(baseNames []string) []string {
	args := client.Called(baseNames)
	return args.Get(0).([]string)
}

func (client *MockPolicyDBClient) GetPolicyRuleByID(id string) (*protos.PolicyRule, error) {
	return nil, nil
}

func (client *MockPolicyDBClient) GetOmnipresentRules() ([]string, []string) {
	args := client.Called()
	return args.Get(0).([]string), args.Get(1).([]string)
}

type MockCreditClient struct {
	mock.Mock
}

func (cc *MockCreditClient) SendCreditControlRequest(
	server *diameter.DiameterServerConfig,
	done chan interface{},
	request *gy.CreditControlRequest,
) error {
	args := cc.Called(server, done, request)
	return args.Error(0)
}

func (cc *MockCreditClient) IgnoreAnswer(request *gy.CreditControlRequest) {
	return
}

func (cc *MockCreditClient) EnableConnections() error {
	cc.Called()
	return nil
}

func (cc *MockCreditClient) DisableConnections(period time.Duration) {
	cc.Called(period)
	return
}

type sessionMocks struct {
	gx       *MockPolicyClient
	gy       *MockCreditClient
	policydb *MockPolicyDBClient
}

func TestSessionControllerPerSessionInit(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerSessionInit),
	)
	standardUsageTest(t, srv, mocks, gy.PerSessionInit)
}

func TestSessionControllerPerKeyInit(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}

	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerKeyInit),
	)

	standardUsageTest(t, srv, mocks, gy.PerKeyInit)
}

func TestStartSessionGxFail(t *testing.T) {
	// Set up mocks
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}

	// Send back DIAMETER_RATING_FAILED (5031) from gx
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.DiameterRatingFailed),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()
	// If gx fails gy should not be used at all

	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerKeyInit),
	)
	ctx := context.Background()
	_, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	mocks.gx.AssertExpectations(t)
	assert.Error(t, err)
}

func TestStartSessionGyFail(t *testing.T) {
	// Set up mocks
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		activationTime := time.Unix(1, 0)
		deactivationTime := time.Unix(2, 0)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:            []string{"static_rule_1"},
				RuleActivationTime:   &activationTime,
				RuleDeactivationTime: &deactivationTime,
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()

	mocks.policydb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{{RatingGroup: 1}}, nil).Once()
	// no omnipresent rules
	mocks.policydb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mocks.policydb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Send back DIAMETER_RATING_FAILED (5031) from gy
	mocks.gy.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gy.CreditControlRequest)
		done <- &gy.CreditControlAnswer{
			ResultCode:    uint32(diameter.DiameterRatingFailed),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()

	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerKeyInit),
	)
	ctx := context.Background()
	_, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	mocks.gx.AssertExpectations(t)
	assert.Error(t, err)
}

func standardUsageTest(
	t *testing.T,
	srv *servicers.CentralSessionController,
	mocks *sessionMocks,
	initMethod gy.InitMethod) {
	ctx := context.Background()

	maxReqBWUL := uint32(128000)
	maxReqBWDL := uint32(128000)
	key1 := []byte("key1")

	// send static rules back
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		qos := gx.QosInformation{MaxReqBwUL: &maxReqBWUL, MaxReqBwDL: &maxReqBWDL}
		redirect := gx.RedirectInformation{
			RedirectSupport:       1,
			RedirectAddressType:   2,
			RedirectServerAddress: "http://www.example.com/",
		}

		var (
			rg20 uint32 = 20
			si20 uint32 = 201
			rg21 uint32 = 21
		)
		activationTime := time.Unix(1, 0)
		deactivationTime := time.Unix(2, 0)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:     []string{"static_rule_1", "static_rule_2"},
				RuleBaseNames: []string{"base_10"},
				RuleDefinitions: []*gx.RuleDefinition{
					&gx.RuleDefinition{
						RuleName:            "dyn_rule_20",
						RatingGroup:         &rg20,
						ServiceIdentifier:   &si20,
						Precedence:          100,
						MonitoringKey:       key1,
						RedirectInformation: &redirect,
						Qos:                 &qos,
						FlowDescriptions: []string{
							"permit out ip from any to any",
							"permit in ip from any to 0.0.0.1",
						},
					},
					&gx.RuleDefinition{
						RuleName:    "dyn_rule_21",
						RatingGroup: &rg21,
						Precedence:  200,
					},
				},
				RuleActivationTime:   &activationTime,
				RuleDeactivationTime: &deactivationTime,
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()

	// send rating groups back
	mocks.policydb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})
	mocks.policydb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{
			policydb.ChargingKey{RatingGroup: 1},
			policydb.ChargingKey{RatingGroup: 2},
			policydb.ChargingKey{RatingGroup: 10},
			policydb.ChargingKey{RatingGroup: 11},
			policydb.ChargingKey{RatingGroup: 11},
			policydb.ChargingKey{RatingGroup: 20, ServiceIdTracking: true, ServiceIdentifier: 201},
			policydb.ChargingKey{RatingGroup: 21}}, nil).Once()
	// no omnipresent rules
	mocks.policydb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mocks.policydb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()
	multiReqType := credit_control.CRTInit // type of CCR sent to get credits
	if initMethod == gy.PerSessionInit {
		mocks.gy.On(
			"SendCreditControlRequest",
			mock.Anything,
			mock.Anything,
			mock.MatchedBy(getGyCCRMatcher(credit_control.CRTInit)),
		).Return(nil).Run(returnDefaultGyResponse).Once()
		multiReqType = credit_control.CRTUpdate // on per session init, credits are received through CCR-Updates
	}
	// return default responses for gy CCR's, depending on init method
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(multiReqType)),
	).Return(nil).Run(returnDefaultGyResponse).Once()
	createResponse, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	mocks.gx.AssertExpectations(t)
	mocks.gy.AssertExpectations(t)
	mocks.policydb.AssertExpectations(t)
	assert.Equal(t, 6, len(createResponse.Credits)) // 2 static, 2 dynamic, 2 base
	assert.Equal(t, 2, len(createResponse.DynamicRules))

	allRuleIDs := []string{}
	for _, staticRule := range createResponse.StaticRules {
		allRuleIDs = append(allRuleIDs, staticRule.RuleId)
		assert.Equal(t, &timestamp.Timestamp{Seconds: 1}, staticRule.ActivationTime)
		assert.Equal(t, &timestamp.Timestamp{Seconds: 2}, staticRule.DeactivationTime)
	}
	assert.ElementsMatch(t, allRuleIDs, []string{"static_rule_1", "static_rule_2", "base_rule_1", "base_rule_2"})

	for _, rule := range createResponse.DynamicRules {
		if rule.PolicyRule.Id == "dyn_rule_20" {
			assert.Equal(t, protos.RedirectInformation_ENABLED, rule.PolicyRule.Redirect.Support)
			assert.Equal(t, protos.RedirectInformation_URL, rule.PolicyRule.Redirect.AddressType)
			assert.Equal(t, "http://www.example.com/", rule.PolicyRule.Redirect.ServerAddress)
			assert.Equal(t, maxReqBWUL, rule.PolicyRule.Qos.MaxReqBwUl)
			assert.Equal(t, maxReqBWDL, rule.PolicyRule.Qos.MaxReqBwDl)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 1}, rule.ActivationTime)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 2}, rule.DeactivationTime)
		} else if rule.PolicyRule.Id == "dyn_rule_21" {
			assert.Nil(t, rule.PolicyRule.Redirect)
			assert.Nil(t, rule.PolicyRule.Qos)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 1}, rule.ActivationTime)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 2}, rule.DeactivationTime)
		} else {
			assert.Fail(t, "Unknown rule id returned")
		}
	}
	ratingGroups := []uint32{}
	for _, update := range createResponse.Credits {
		assert.True(t, update.Success)
		assert.Equal(t, IMSI1, update.Sid)
		ratingGroups = append(ratingGroups, update.ChargingKey)
		if update.ChargingKey == 20 {
			assert.NotNil(t, update.ServiceIdentifier)
			assert.Equal(t, uint32(201), update.GetServiceIdentifier().GetValue())
		} else {
			assert.Nil(t, update.ServiceIdentifier)
		}
		assert.Equal(t, uint64(2048), update.Credit.GrantedUnits.Total.Volume)
		assert.True(t, update.Credit.GrantedUnits.Total.IsValid)
		assert.False(t, update.Credit.GrantedUnits.Rx.IsValid)
		assert.False(t, update.Credit.GrantedUnits.Tx.IsValid)
		assert.Equal(t, uint32(3600), update.Credit.ValidityTime)
		assert.Equal(t, protos.CreditUpdateResponse_UPDATE, update.Type)
	}
	assert.ElementsMatch(t, ratingGroups, []uint32{1, 2, 10, 11, 20, 21})

	// updates
	mocks.gy.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(returnDefaultGyResponse).Times(2)
	updateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
		},
	})
	mocks.gy.AssertExpectations(t)
	assert.Equal(t, 2, len(updateResponse.Responses))
	for _, update := range updateResponse.Responses {
		assert.True(t, update.Success)
		assert.Equal(t, IMSI1, update.Sid)
		assert.True(t, update.ChargingKey == 1 || update.ChargingKey == 2)
	}

	// Connection Manager tests
	mocks.gx.On("DisableConnections", mock.Anything).Return()
	mocks.gy.On("DisableConnections", mock.Anything).Return()
	void, err := srv.Disable(ctx, &fegprotos.DisableMessage{
		DisablePeriodSecs: 10,
	})
	mocks.gx.AssertExpectations(t)
	mocks.gy.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, void)

	mocks.gx.On("EnableConnections").Return()
	mocks.gy.On("EnableConnections").Return()
	void, err = srv.Enable(ctx, &orcprotos.Void{})

	mocks.gx.AssertExpectations(t)
	mocks.gy.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, void)
}

func TestSessionCreateWithOmnipresentRules(t *testing.T) {
	// Set up mocks
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}

	// send static rules back
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:     []string{"static_rule_1", "static_rule_2"},
				RuleBaseNames: []string{"base_10"},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mocks.policydb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})
	mocks.policydb.On("GetRuleIDsForBaseNames", []string{"omnipresent_base_1"}).Return([]string{"omnipresent_rule_2"})
	mocks.policydb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return([]policydb.ChargingKey{}, nil).Once()
	mocks.policydb.On("GetOmnipresentRules").Return([]string{"omnipresent_rule_1"}, []string{"omnipresent_base_1"})
	ctx := context.Background()
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerKeyInit),
	)
	response, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	assert.NoError(t, err)

	mocks.gx.AssertExpectations(t)
	mocks.policydb.AssertExpectations(t)

	assert.Equal(t, 6, len(response.StaticRules))
	expectedRuleIDs := []string{"static_rule_1", "static_rule_2", "base_rule_1", "base_rule_2", "omnipresent_rule_1", "omnipresent_rule_2"}
	actualRuleIDs := funk.Map(response.StaticRules, func(ruleInstall *protos.StaticRuleInstall) string { return ruleInstall.RuleId }).([]string)
	assert.ElementsMatch(t, expectedRuleIDs, actualRuleIDs)
}

func TestSessionControllerTimeouts(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerSessionInit),
	)

	ctx := context.Background()

	// depending on request number, "lose" request
	var units uint64 = 2048
	mocks.gy.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gy.CreditControlRequest)
		if request.RequestNumber%2 == 0 {
			return
		} else {
			done <- &gy.CreditControlAnswer{
				ResultCode:    uint32(diameter.SuccessCode),
				SessionID:     request.SessionID,
				RequestNumber: request.RequestNumber,
				Credits: []*gy.ReceivedCredits{&gy.ReceivedCredits{
					RatingGroup:  request.Credits[0].RatingGroup,
					GrantedUnits: &credit_control.GrantedServiceUnit{TotalOctets: &units},
					ValidityTime: 3600,
				}},
			}
		}
	}).Return(nil).Times(3)
	updateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI2, 2, 2, protos.CreditUsage_TERMINATED),
			createUsageUpdate(IMSI1, 1, 2, protos.CreditUsage_TERMINATED),
		},
	})
	mocks.gy.AssertExpectations(t)
	assert.Equal(t, 3, len(updateResponse.Responses))
	// Every other request will fail
	countFailed := 0
	for _, update := range updateResponse.Responses {
		if !update.Success {
			countFailed++
		}
	}
	assert.Equal(t, 2, countFailed)
}

func TestSessionTermination(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerSessionInit),
	)
	ctx := context.Background()

	// Return success for Gx termination
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTTerminate)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()
	// Return success for Gy terminations
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTTerminate)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gy.CreditControlRequest)
		done <- &gy.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()

	termResponse, err := srv.TerminateSession(ctx, &protos.SessionTerminateRequest{
		Sid:       IMSI2,
		SessionId: fmt.Sprintf("%s-1234", IMSI2),
		CreditUsages: []*protos.CreditUsage{
			createUsage(2, protos.CreditUsage_TERMINATED),
			createUsage(1, protos.CreditUsage_TERMINATED),
		},
	})
	mocks.gy.AssertExpectations(t)
	mocks.gx.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, IMSI2, termResponse.Sid)
	assert.Equal(t, fmt.Sprintf("%s-1234", IMSI2), termResponse.SessionId)
}

func TestGxUsageMonitoring(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerSessionInit),
	)
	ctx := context.Background()

	// Return success for Gx Update
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGyResponse).Times(2)
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGxUpdateResponse).Times(2)

	updateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
		},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI1, "mkey2", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
		},
	})
	mocks.gy.AssertExpectations(t)
	mocks.gx.AssertExpectations(t)
	assert.Equal(t, 2, len(updateResponse.Responses))
	assert.Equal(t, 2, len(updateResponse.UsageMonitorResponses))
	for _, update := range updateResponse.Responses {
		assert.True(t, update.Success)
		assert.Equal(t, IMSI1, update.Sid)
		assert.True(t, update.ChargingKey == 1 || update.ChargingKey == 2)
	}
	for _, update := range updateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.Equal(t, IMSI1, update.Sid)
		assert.Equal(t, protos.UsageMonitoringCredit_CONTINUE, update.Credit.Action)
		assert.Equal(t, uint64(2048), update.Credit.GrantedUnits.Total.Volume)
		if string(update.Credit.MonitoringKey) == "mkey" {
			assert.Equal(t, protos.MonitoringLevel_SESSION_LEVEL, update.Credit.Level)
		} else if string(update.Credit.MonitoringKey) == "mkey2" {
			assert.Equal(t, protos.MonitoringLevel_PCC_RULE_LEVEL, update.Credit.Level)
		} else {
			assert.True(t, false)
		}
	}

	// test usage monitoring disabling
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(returnEmptyGxUpdateResponse).Times(1)

	emptyUpdateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocks.gx.AssertExpectations(t)
	assert.Equal(t, 1, len(emptyUpdateResponse.UsageMonitorResponses))
	update := emptyUpdateResponse.UsageMonitorResponses[0]
	assert.True(t, update.Success)
	assert.Equal(t, IMSI1, update.Sid)
	assert.Equal(t, protos.UsageMonitoringCredit_DISABLE, update.Credit.Action)
	assert.Nil(t, update.Credit.GrantedUnits)
	assert.Equal(t, protos.MonitoringLevel_SESSION_LEVEL, update.Credit.Level)

	// Test that static rule install avp in CCA-Update by rule names gets propagated properly
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleInstallGxUpdateResponse([]string{"static1", "static2"}, []string{})).Times(1)

	ruleInstallUpdateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocks.gx.AssertExpectations(t)
	assert.Equal(t, 1, len(ruleInstallUpdateResponse.UsageMonitorResponses))
	update = ruleInstallUpdateResponse.UsageMonitorResponses[0]
	assert.True(t, update.Success)
	assert.Equal(t, IMSI1, update.Sid)
	assert.Nil(t, update.Credit.GrantedUnits)
	assert.Equal(t, "static1", update.StaticRulesToInstall[0].RuleId)
	assert.Equal(t, "static2", update.StaticRulesToInstall[1].RuleId)

	// Test that static rule install avp in CCA-Update by rule base names gets propagated properly
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleInstallGxUpdateResponse([]string{}, []string{"base_10"})).Times(1)
	mocks.policydb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})

	ruleInstallUpdateResponse, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocks.gx.AssertExpectations(t)
	assert.Equal(t, 1, len(ruleInstallUpdateResponse.UsageMonitorResponses))
	update = ruleInstallUpdateResponse.UsageMonitorResponses[0]
	assert.True(t, update.Success)
	assert.Equal(t, IMSI1, update.Sid)
	assert.Nil(t, update.Credit.GrantedUnits)
	assert.Equal(t, "base_rule_1", update.StaticRulesToInstall[0].RuleId)
	assert.Equal(t, "base_rule_2", update.StaticRulesToInstall[1].RuleId)

	// Test that dynamic rule install avp in CCA-Update gets propagated properly
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(returnDynamicRuleInstallGxUpdateResponse).Times(1)

	ruleInstallUpdateResponse, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocks.gx.AssertExpectations(t)
	assert.Equal(t, 1, len(ruleInstallUpdateResponse.UsageMonitorResponses))
	update = ruleInstallUpdateResponse.UsageMonitorResponses[0]
	assert.True(t, update.Success)
	assert.Equal(t, IMSI1, update.Sid)
	assert.Nil(t, update.Credit.GrantedUnits)
	assert.Equal(t, "dyn_rule_20", update.DynamicRulesToInstall[0].PolicyRule.Id)

	// Test that rule remove avp in CCA-Update by rule names gets propagated properly
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleDisableGxUpdateResponse([]string{"rule1", "rule2"}, []string{})).Times(1)

	ruleDisableUpdateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocks.gx.AssertExpectations(t)
	assert.Equal(t, 1, len(ruleDisableUpdateResponse.UsageMonitorResponses))
	update = ruleDisableUpdateResponse.UsageMonitorResponses[0]
	assert.True(t, update.Success)
	assert.Equal(t, IMSI1, update.Sid)
	assert.Nil(t, update.Credit.GrantedUnits)
	assert.Equal(t, []string{"rule1", "rule2"}, update.RulesToRemove)

	// Test that rule remove avp in CCA-Update by base names gets propagated properly
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleDisableGxUpdateResponse([]string{}, []string{"base_10"})).Times(1)
	mocks.policydb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})

	ruleDisableUpdateResponse, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocks.gx.AssertExpectations(t)
	assert.Equal(t, 1, len(ruleDisableUpdateResponse.UsageMonitorResponses))
	update = ruleDisableUpdateResponse.UsageMonitorResponses[0]
	assert.True(t, update.Success)
	assert.Equal(t, IMSI1, update.Sid)
	assert.Nil(t, update.Credit.GrantedUnits)
	assert.Equal(t, []string{"base_rule_1", "base_rule_2"}, update.RulesToRemove)
}

func TestGetHealthStatus(t *testing.T) {
	err := initMconfig()
	assert.NoError(t, err)

	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		getTestConfig(gy.PerSessionInit),
	)
	ctx := context.Background()

	// Return success for Gx/Gy CCR-Update
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGyResponse).Times(2)
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGxUpdateResponse).Times(2)

	_, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
		},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI1, "mkey2", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
		},
	})
	mocks.gy.AssertExpectations(t)
	mocks.gx.AssertExpectations(t)

	status, err := srv.GetHealthStatus(ctx, &orcprotos.Void{})
	assert.NoError(t, err)
	assert.Equal(t, fegprotos.HealthStatus_HEALTHY, status.Health)

	// Return error for Gx/Gy CCR-Updatee
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTUpdate)),
	).Return(fmt.Errorf("Failed to establish new diameter connection; will retry upon first request.")).Run(returnDefaultGyResponse).Times(2)
	mocks.gx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(credit_control.CRTUpdate)),
	).Return(fmt.Errorf("Failed to establish new diameter connection; will retry upon first request.")).Run(returnDefaultGxUpdateResponse).Times(2)

	_, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
		},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI1, "mkey2", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
		},
	})
	mocks.gy.AssertExpectations(t)
	mocks.gx.AssertExpectations(t)

	status, err = srv.GetHealthStatus(ctx, &orcprotos.Void{})
	assert.NoError(t, err)
	assert.Equal(t, fegprotos.HealthStatus_UNHEALTHY, status.Health)
}

func getTestConfig(initMethod gy.InitMethod) *servicers.SessionControllerConfig {
	return &servicers.SessionControllerConfig{
		OCSConfig: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:     "127.0.0.1:3869",
			Protocol: "tcp"},
		},
		PCRFConfig: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:     "127.0.0.1:3870",
			Protocol: "tcp"},
		},
		RequestTimeout: time.Millisecond,
		InitMethod:     initMethod,
	}
}

func createUsageUpdate(
	sid string,
	chargingKey uint32,
	requestNumber uint32,
	requestType protos.CreditUsage_UpdateType,
) *protos.CreditUsageUpdate {
	return &protos.CreditUsageUpdate{
		Usage:         createUsage(chargingKey, requestType),
		SessionId:     fmt.Sprintf("%s-1234", sid),
		RequestNumber: requestNumber,
		Sid:           sid,
	}
}

func createUsageMonitoringRequest(
	sid string,
	monitoringKey string,
	requestNumber uint32,
	monitoringLevel protos.MonitoringLevel,
) *protos.UsageMonitoringUpdateRequest {
	return &protos.UsageMonitoringUpdateRequest{
		Update: &protos.UsageMonitorUpdate{
			BytesTx:       1024,
			BytesRx:       2048,
			MonitoringKey: []byte(monitoringKey),
			Level:         monitoringLevel,
		},
		SessionId:     fmt.Sprintf("%s-1234", sid),
		RequestNumber: requestNumber,
		Sid:           sid,
	}
}

func createUsage(
	chargingKey uint32,
	requestType protos.CreditUsage_UpdateType,
) *protos.CreditUsage {
	return &protos.CreditUsage{
		BytesTx:     1024,
		BytesRx:     2048,
		ChargingKey: chargingKey,
		Type:        requestType,
	}
}

func returnDefaultGyResponse(args mock.Arguments) {
	var units uint64 = 2048
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := make([]*gy.ReceivedCredits, 0, len(request.Credits))

	for _, credit := range request.Credits {
		credits = append(credits, &gy.ReceivedCredits{
			RatingGroup:       credit.RatingGroup,
			ServiceIdentifier: credit.ServiceIdentifier,
			GrantedUnits:      &credit_control.GrantedServiceUnit{TotalOctets: &units},
			ValidityTime:      3600,
			ResultCode:        uint32(diameter.SuccessCode),
		})
	}

	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}

func returnDefaultGxUpdateResponse(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gx.CreditControlRequest)
	monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
	for _, report := range request.UsageReports {
		totalOctets := uint64(2048)
		monitors = append(monitors, &gx.UsageMonitoringInfo{
			MonitoringKey: report.MonitoringKey,
			GrantedServiceUnit: &credit_control.GrantedServiceUnit{
				TotalOctets: &totalOctets,
			},
			Level: report.Level,
		})
	}
	done <- &gx.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		UsageMonitors: monitors,
	}
}

func initMconfig() error {
	fegConfig := `{
		"configsByKey": {
			"session_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.SessionProxyConfig",
				"logLevel": "INFO",
				"gx": {
					"server": {
						 "protocol": "tcp",
						 "address": "",
						 "retransmits": 3,
						 "watchdogInterval": 1,
						 "retryCount": 5,
						 "productName": "magma",
		 				"realm": "magma.com",
		 				"host": "magma-fedgw.magma.com"
					}
				},
				"gy": {
					"server": {
						 "protocol": "tcp",
						 "address": "",
						 "retransmits": 3,
						 "watchdogInterval": 1,
						 "retryCount": 5,
						 "productName": "magma",
		 				 "realm": "magma.com",
		 				 "host": "magma-fedgw.magma.com"
					},
					"initMethod": "PER_KEY"
				},
				"requestFailureThreshold": 0.5,
   				"minimumRequestThreshold": 1
			}
		}
	}`

	err := mconfig.CreateLoadTempConfig(fegConfig)
	if err != nil {
		return err
	}
	return nil
}

func returnEmptyGxUpdateResponse(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gx.CreditControlRequest)
	monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
	for _, report := range request.UsageReports {
		monitors = append(monitors, &gx.UsageMonitoringInfo{
			MonitoringKey:      report.MonitoringKey,
			GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
			Level:              report.Level,
		})
	}
	done <- &gx.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		UsageMonitors: monitors,
	}
}

func getRuleInstallGxUpdateResponse(ruleNames, baseNames []string) func(mock.Arguments) {
	return func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
		for _, report := range request.UsageReports {
			monitors = append(monitors, &gx.UsageMonitoringInfo{
				MonitoringKey:      report.MonitoringKey,
				GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
				Level:              report.Level,
			})
		}
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
			UsageMonitors: monitors,
			RuleInstallAVP: []*gx.RuleInstallAVP{
				{
					RuleNames:     ruleNames,
					RuleBaseNames: baseNames,
				},
			},
		}
	}
}

func returnDynamicRuleInstallGxUpdateResponse(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gx.CreditControlRequest)
	monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
	for _, report := range request.UsageReports {
		monitors = append(monitors, &gx.UsageMonitoringInfo{
			MonitoringKey:      report.MonitoringKey,
			GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
			Level:              report.Level,
		})
	}
	activationTime := time.Unix(1, 0)
	deactivationTime := time.Unix(2, 0)
	done <- &gx.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		UsageMonitors: monitors,
		RuleInstallAVP: []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleDefinitions: []*gx.RuleDefinition{
					&gx.RuleDefinition{
						RuleName: "dyn_rule_20",
						//RatingGroup: swag.Uint32(20),
					},
				},
				RuleActivationTime:   &activationTime,
				RuleDeactivationTime: &deactivationTime,
			},
		},
	}
}

func getRuleDisableGxUpdateResponse(ruleNames []string, ruleBaseNames []string) func(mock.Arguments) {
	return func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
		for _, report := range request.UsageReports {
			monitors = append(monitors, &gx.UsageMonitoringInfo{
				MonitoringKey:      report.MonitoringKey,
				GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
				Level:              report.Level,
			})
		}
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
			UsageMonitors: monitors,
			RuleRemoveAVP: []*gx.RuleRemoveAVP{
				{
					RuleNames:     ruleNames,
					RuleBaseNames: ruleBaseNames,
				},
			},
		}
	}
}

func getGyCCRMatcher(ccrType credit_control.CreditRequestType) interface{} {
	return func(request *gy.CreditControlRequest) bool {
		return request.Type == ccrType
	}
}

func getGxCCRMatcher(ccrType credit_control.CreditRequestType) interface{} {
	return func(request *gx.CreditControlRequest) bool {
		return request.Type == ccrType
	}
}

/***** UseGyForAuthOnlySuccess Test Cases *****/
func TestSessionControllerUseGyForAuthOnlySuccess(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}
	activationTime := time.Unix(1, 0)
	deactivationTime := time.Unix(2, 0)
	// send static rules back
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:            []string{"static_rule_1"},
				RuleActivationTime:   &activationTime,
				RuleDeactivationTime: &deactivationTime,
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()

	mocks.policydb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{{RatingGroup: 3}}, nil).Once()
	mocks.policydb.On("GetOmnipresentRules").Return([]string{"omnipresent_1"}, []string{}).Once()
	mocks.policydb.On("GetRuleIDsForBaseNames", []string{}).Return([]string{}).Once()

	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessNoRatingGroup).Once()

	cfg := getTestConfig(gy.PerKeyInit)
	cfg.UseGyForAuthOnly = true
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		cfg,
	)
	ctx := context.Background()
	res, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	mocks.gx.AssertExpectations(t)
	assert.NoError(t, err)
	expectedStaticRule1 := &protos.StaticRuleInstall{
		RuleId:           "static_rule_1",
		ActivationTime:   gx.ConvertToProtoTimestamp(&activationTime),
		DeactivationTime: gx.ConvertToProtoTimestamp(&deactivationTime),
	}
	assert.ElementsMatch(t, []*protos.StaticRuleInstall{{RuleId: "omnipresent_1"}, expectedStaticRule1}, res.StaticRules)
}

func TestSessionControllerUseGyForAuthOnlyNoRatingGroup(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:       []string{"static_rule_1"},
				RuleDefinitions: []*gx.RuleDefinition{},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mocks.policydb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{}, nil).Once()
	// no omnipresent rule
	mocks.policydb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mocks.policydb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Even if there are no rating groups, gy CCR-I will be called.
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessNoRatingGroup).Once()

	cfg := getTestConfig(gy.PerKeyInit)
	cfg.UseGyForAuthOnly = true
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		cfg,
	)
	ctx := context.Background()
	_, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	mocks.gx.AssertExpectations(t)
	assert.NoError(t, err)
}

func returnGySuccessNoRatingGroup(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := make([]*gy.ReceivedCredits, 0, len(request.Credits))
	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}

func TestSessionControllerUseGyForAuthOnlyCreditLimitReached(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:       []string{"static_rule_1"},
				RuleDefinitions: []*gx.RuleDefinition{},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mocks.policydb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{}, nil).Once()
	// no omnipresent rule
	mocks.policydb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mocks.policydb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Even if there are no rating groups, gy CCR-I will be called.
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessCreditLimitReached).Once()

	cfg := getTestConfig(gy.PerKeyInit)
	cfg.UseGyForAuthOnly = true
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		cfg,
	)
	ctx := context.Background()
	_, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	mocks.gx.AssertExpectations(t)
	assert.NoError(t, err)
}

func returnGySuccessCreditLimitReached(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := []*gy.ReceivedCredits{
		&gy.ReceivedCredits{
			ResultCode: diameter.DiameterCreditLimitReached,
		},
	}

	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}

func TestSessionControllerUseGyForAuthOnlySubscriberBarred(t *testing.T) {
	mocks := &sessionMocks{
		gy:       &MockCreditClient{},
		gx:       &MockPolicyClient{},
		policydb: &MockPolicyDBClient{},
	}

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocks.gx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:       []string{"static_rule_1"},
				RuleDefinitions: []*gx.RuleDefinition{},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mocks.policydb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{}, nil).Once()
	// no omnipresent rule
	mocks.policydb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mocks.policydb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Even if there are no rating groups, gy CCR-I will be called.
	mocks.gy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessSubscriberBarred).Once()

	cfg := getTestConfig(gy.PerKeyInit)
	cfg.UseGyForAuthOnly = true
	srv := servicers.NewCentralSessionController(
		mocks.gy,
		mocks.gx,
		mocks.policydb,
		cfg,
	)
	ctx := context.Background()
	_, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: "00101-1234",
	})
	mocks.gx.AssertExpectations(t)
	assert.Error(t, err)
}

func returnGySuccessSubscriberBarred(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := []*gy.ReceivedCredits{
		&gy.ReceivedCredits{
			ResultCode: diameter.DiameterRatingFailed,
		},
	}

	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}
