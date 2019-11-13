/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"
	"sync"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/session_proxy/metrics"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
)

// CentralSessionController acts as the gRPC server for accepting calls from
// gateways to start new UE sessions and retrieve traffic policy and credits.
type CentralSessionController struct {
	creditClient  gy.CreditClient
	policyClient  gx.PolicyClient
	dbClient      policydb.PolicyDBClient
	cfg           *SessionControllerConfig
	healthTracker *metrics.SessionHealthTracker
}

// SessionControllerConfig stores all the needed configuration for running
// gx and gy clients
type SessionControllerConfig struct {
	OCSConfig      *diameter.DiameterServerConfig
	PCRFConfig     *diameter.DiameterServerConfig
	RequestTimeout time.Duration
	InitMethod     gy.InitMethod
}

type ruleTiming struct {
	activationTime   *timestamp.Timestamp
	deactivationTime *timestamp.Timestamp
}

// NewCentralSessionController constructs a CentralSessionController
// and registers external handlers
func NewCentralSessionController(
	creditClient gy.CreditClient,
	policyClient gx.PolicyClient,
	dbClient policydb.PolicyDBClient,
	cfg *SessionControllerConfig,
) *CentralSessionController {
	return &CentralSessionController{
		creditClient:  creditClient,
		policyClient:  policyClient,
		dbClient:      dbClient,
		cfg:           cfg,
		healthTracker: metrics.NewSessionHealthTracker(),
	}
}

// CreateSession begins a UE session by requesting rules from PCEF
// and credit from OCS (if RatingGroup is present) and returning them.
func (srv *CentralSessionController) CreateSession(
	ctx context.Context,
	request *protos.CreateSessionRequest,
) (*protos.CreateSessionResponse, error) {
	glog.V(2).Info("Trying to create session")
	imsi := removeSidPrefix(request.Subscriber.Id)
	sessionID := request.SessionId
	gxCCAInit, err := srv.sendInitialGxRequest(imsi, request)
	metrics.UpdateGxRecentRequestMetrics(err)
	if err != nil {
		metrics.PcrfCcrInitSendFailures.Inc()
		glog.Errorf("Failed to send initial Gx request: %s", err)
		return nil, err
	}
	metrics.PcrfCcrInitRequests.Inc()

	var ruleNames []string
	var ruleDefs []*gx.RuleDefinition
	for _, rule := range gxCCAInit.RuleInstallAVP {
		ruleNames = append(ruleNames, rule.RuleNames...)
		if len(rule.RuleBaseNames) > 0 {
			ruleNames = append(ruleNames, srv.dbClient.GetRuleIDsForBaseNames(rule.RuleBaseNames)...)
		}
		ruleDefs = append(ruleDefs, rule.RuleDefinitions...)
	}

	policyRules := getPolicyRulesFromDefinitions(ruleDefs)
	keys, err := srv.dbClient.GetChargingKeysForRules(ruleNames, policyRules)
	if err != nil {
		glog.Errorf("Failed to get charging keys for rules: %s", err)
		return nil, err
	}
	keys = removeDuplicateChargingKeys(keys)
	credits := []*protos.CreditUpdateResponse{}

	if len(keys) > 0 {
		if srv.cfg.InitMethod == gy.PerSessionInit {
			_, err = srv.sendSingleCreditRequest(getCCRInitRequest(imsi, request))
			metrics.UpdateGyRecentRequestMetrics(err)
			if err != nil {
				metrics.OcsCcrInitSendFailures.Inc()
				glog.Errorf("Failed to send first single credit request: %s", err)
				return nil, err
			}
			metrics.OcsCcrInitRequests.Inc()
		}

		updateRequest := getCCRInitialCreditRequest(imsi, request, keys, srv.cfg.InitMethod)
		ans, err := srv.sendSingleCreditRequest(updateRequest)
		metrics.UpdateGyRecentRequestMetrics(err)
		if err != nil {
			metrics.OcsCcrInitSendFailures.Inc()
			glog.Errorf("Failed to send second single credit request: %s", err)
			return nil, err
		}
		metrics.OcsCcrInitRequests.Inc()
		credits = getInitialCreditResponsesFromCCA(ans, updateRequest)
	}

	staticRules, dynamicRules := gx.ParseRuleInstallAVPs(
		srv.dbClient,
		gxCCAInit.RuleInstallAVP,
	)

	return &protos.CreateSessionResponse{
		Credits:       credits,
		StaticRules:   staticRules,
		DynamicRules:  dynamicRules,
		UsageMonitors: getUsageMonitorsFromCCA(imsi, sessionID, gxCCAInit),
	}, nil
}

func removeDuplicateChargingKeys(keysIn []uint32) []uint32 {
	keysOut := []uint32{}
	keyMap := make(map[uint32]bool)
	for _, k := range keysIn {
		if _, ok := keyMap[k]; !ok {
			keysOut = append(keysOut, k)
			keyMap[k] = true
		}
	}
	return keysOut
}

// UpdateSession handles periodic updates from gateways that include quota
// exhaustion and terminations.
func (srv *CentralSessionController) UpdateSession(
	ctx context.Context,
	request *protos.UpdateSessionRequest,
) (*protos.UpdateSessionResponse, error) {
	// Then send out updates
	var wg sync.WaitGroup
	wg.Add(2)

	var gxUpdateResponses []*protos.UsageMonitoringUpdateResponse
	go func() {
		defer wg.Done()
		requests := getGxUpdateRequestsFromUsage(request.UsageMonitors)
		gxUpdateResponses = srv.sendMultipleGxRequestsWithTimeout(requests, srv.cfg.RequestTimeout)
	}()
	var gyUpdateResponses []*protos.CreditUpdateResponse
	go func() {
		defer wg.Done()
		requests := getGyUpdateRequestsFromUsage(request.Updates)
		gyUpdateResponses = srv.sendMultipleGyRequestsWithTimeout(requests, srv.cfg.RequestTimeout)
	}()
	wg.Wait()

	return &protos.UpdateSessionResponse{
		Responses:             gyUpdateResponses,
		UsageMonitorResponses: gxUpdateResponses,
	}, nil
}

// TerminateSession handles a session termination by sending single CCR-T on Gx
// sending CCR-T per rating group on Gy
func (srv *CentralSessionController) TerminateSession(
	ctx context.Context,
	request *protos.SessionTerminateRequest,
) (*protos.SessionTerminateResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := srv.sendTerminationGxRequest(request)
		metrics.UpdateGxRecentRequestMetrics(err)
		if err != nil {
			metrics.PcrfCcrTerminateSendFailures.Inc()
			glog.Errorf("Error sending gx termination: %s", err)
		} else {
			metrics.PcrfCcrTerminateRequests.Inc()
		}
	}()
	go func() {
		defer wg.Done()
		_, err := srv.sendSingleCreditRequest(getTerminateRequestFromUsage(request))
		metrics.UpdateGyRecentRequestMetrics(err)
		if err != nil {
			metrics.OcsCcrTerminateSendFailures.Inc()
			glog.Errorf("Error sending gy termination: %s", err)
		} else {
			metrics.OcsCcrTerminateRequests.Inc()
		}
	}()
	wg.Wait()
	// in the event of any errors on Gx or Gy, the session should regardless be
	// terminated, so there are no errors sent back
	return &protos.SessionTerminateResponse{
		Sid:       request.Sid,
		SessionId: request.SessionId,
	}, nil
}

// Disable closes all existing diameter connections and disables
// connection creation for the time specified in the request
func (srv *CentralSessionController) Disable(
	ctx context.Context,
	req *fegprotos.DisableMessage,
) (*orcprotos.Void, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil Disable Request")
	}
	srv.policyClient.DisableConnections(time.Duration(req.DisablePeriodSecs) * time.Second)
	srv.creditClient.DisableConnections(time.Duration(req.DisablePeriodSecs) * time.Second)
	return &orcprotos.Void{}, nil
}

// Enable enables diameter connection creation
// If creation is already enabled, Enable has no effect
func (srv *CentralSessionController) Enable(
	ctx context.Context,
	void *orcprotos.Void,
) (*orcprotos.Void, error) {
	srv.policyClient.EnableConnections()
	srv.creditClient.EnableConnections()
	return &orcprotos.Void{}, nil
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
func (srv *CentralSessionController) GetHealthStatus(
	ctx context.Context,
	void *orcprotos.Void,
) (*fegprotos.HealthStatus, error) {
	currentMetrics, err := metrics.GetCurrentHealthMetrics()
	if err != nil {
		return &fegprotos.HealthStatus{
			Health:        fegprotos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("Error occured while retrieving health metrics: %s", err),
		}, err
	}
	deltaMetrics, err := srv.healthTracker.Metrics.GetDelta(currentMetrics)
	if err != nil {
		return &fegprotos.HealthStatus{
			Health:        fegprotos.HealthStatus_UNHEALTHY,
			HealthMessage: err.Error(),
		}, err
	}
	gxReqTotal := deltaMetrics.PcrfInitTotal + deltaMetrics.PcrfInitSendFailures +
		deltaMetrics.PcrfUpdateTotal + deltaMetrics.PcrfUpdateSendFailures +
		deltaMetrics.PcrfTerminateTotal + deltaMetrics.PcrfTerminateSendFailures
	gxFailureTotal := deltaMetrics.PcrfInitSendFailures + deltaMetrics.PcrfUpdateSendFailures +
		deltaMetrics.PcrfTerminateSendFailures + deltaMetrics.GxTimeouts + deltaMetrics.GxUnparseableMsg

	gxStatus := srv.getHealthStatusForGxRequests(gxFailureTotal, gxReqTotal)
	if gxStatus.Health == fegprotos.HealthStatus_UNHEALTHY {
		return gxStatus, nil
	}

	gyReqTotal := deltaMetrics.OcsInitTotal + deltaMetrics.OcsInitSendFailures +
		deltaMetrics.OcsUpdateTotal + deltaMetrics.OcsUpdateSendFailures +
		deltaMetrics.OcsTerminateTotal + deltaMetrics.OcsTerminateSendFailures
	gyFailureTotal := deltaMetrics.OcsInitSendFailures + deltaMetrics.OcsUpdateSendFailures +
		deltaMetrics.OcsTerminateSendFailures + deltaMetrics.GyTimeouts + deltaMetrics.GyUnparseableMsg

	gyStatus := srv.getHealthStatusForGyRequests(gyFailureTotal, gyReqTotal)
	if gyStatus.Health == fegprotos.HealthStatus_UNHEALTHY {
		return gyStatus, nil
	}
	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "All metrics appear healthy",
	}, nil
}

func (srv *CentralSessionController) getHealthStatusForGxRequests(failures, total int64) *fegprotos.HealthStatus {
	gxExceedsThreshold := total >= int64(srv.healthTracker.MinimumRequestThreshold) &&
		float64(failures)/float64(total) >= float64(srv.healthTracker.RequestFailureThreshold)
	if gxExceedsThreshold {
		unhealthyMsg := fmt.Sprintf("Metric Gx Request Failure Ratio >= threshold %f; %d / %d",
			srv.healthTracker.RequestFailureThreshold,
			failures,
			total,
		)
		return &fegprotos.HealthStatus{
			Health:        fegprotos.HealthStatus_UNHEALTHY,
			HealthMessage: unhealthyMsg,
		}
	}
	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "Gx metrics appear healthy",
	}
}

func (srv *CentralSessionController) getHealthStatusForGyRequests(failures, total int64) *fegprotos.HealthStatus {
	gyExceedsThreshold := total >= int64(srv.healthTracker.MinimumRequestThreshold) &&
		float64(failures)/float64(total) >= float64(srv.healthTracker.RequestFailureThreshold)
	if gyExceedsThreshold {
		unhealthyMsg := fmt.Sprintf("Metric Gy Request Failure Ratio >= threshold %f; %d / %d",
			srv.healthTracker.RequestFailureThreshold,
			failures,
			total,
		)
		return &fegprotos.HealthStatus{
			Health:        fegprotos.HealthStatus_UNHEALTHY,
			HealthMessage: unhealthyMsg,
		}
	}
	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "Gy metrics appear healthy",
	}
}
