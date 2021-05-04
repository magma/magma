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

package servicers

import (
	"context"
	"fmt"
	"sync"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/session_proxy/metrics"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/errors"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/thoas/go-funk"
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
	// This flag enables a specific type of behavior.
	// 1. Ensures a Gy CCR-I is called in CreateSession when Gx CCR-I succeeds,
	// even if there is no rating group returned by Gx CCR-A.
	// 2. Ensures all Multi Service Credit Control entities have 2001 result
	// code for CreateSession to succeed.
	UseGyForAuthOnly bool
	DisableGx        bool
	DisableGy        bool
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
	_ context.Context,
	request *protos.CreateSessionRequest,
) (*protos.CreateSessionResponse, error) {
	glog.V(2).Info("Trying to create session")
	if err := checkCreateSessionRequest(request); err != nil {
		return nil, err
	}
	imsi := credit_control.RemoveIMSIPrefix(request.GetCommonContext().GetSid().Id)

	// GX
	// In case srv.cfg.DisableGx this function will return a default response
	gxCCAInit, err := srv.sendInitialGxRequestOrGenerateEmptyResponse(imsi, request)
	metrics.ReportCreateGxSession(err)
	if err != nil {
		glog.Errorf("Failed to send initial Gx request: %s", err)
		return nil, err
	}
	err = injectOmnipresentRules(srv, gxCCAInit)
	if err != nil {
		glog.Errorf("Could not inject omnipresent Rules, skipping. Error: %+v", err)
	}

	staticRuleInstalls, dynamicRuleInstalls := gx.ParseRuleInstallAVPs(srv.dbClient, gxCCAInit.RuleInstallAVP)
	chargingKeys := srv.getChargingKeysFromRuleInstalls(staticRuleInstalls, dynamicRuleInstalls)
	eventTriggers, revalidationTime := gx.GetEventTriggersRelatedInfo(gxCCAInit.EventTriggers, gxCCAInit.RevalidationTime)
	gxOriginHost, gyOriginHost := gxCCAInit.OriginHost, ""

	credits := []*protos.CreditUpdateResponse{}
	// Gy
	if !srv.cfg.DisableGy {
		if srv.cfg.UseGyForAuthOnly {
			return srv.handleUseGyForAuthOnly(
				imsi, request, staticRuleInstalls, dynamicRuleInstalls, gxCCAInit)
		}
		if len(chargingKeys) > 0 {
			gyCCRInit := makeCCRInit(imsi, request, chargingKeys)
			gyCCAInit, err := srv.sendSingleCreditRequest(gyCCRInit)
			metrics.ReportCreateGySession(err)
			if err != nil {
				glog.Errorf("Failed to send initial credit request: %s", err)
				return nil, err
			}
			gyOriginHost = gyCCAInit.OriginHost
			credits = getInitialCreditResponsesFromCCA(gyCCRInit, gyCCAInit)
		}
	}
	usageMonitors := getUsageMonitorsFromCCA_I(imsi, request.SessionId, gyOriginHost, gxCCAInit)

	return &protos.CreateSessionResponse{
		Credits:          credits,
		StaticRules:      staticRuleInstalls,
		DynamicRules:     dynamicRuleInstalls,
		UsageMonitors:    usageMonitors,
		TgppCtx:          &protos.TgppContext{GxDestHost: gxOriginHost, GyDestHost: gyOriginHost},
		SessionId:        request.SessionId,
		EventTriggers:    eventTriggers,
		RevalidationTime: revalidationTime,
	}, nil
}

func (srv *CentralSessionController) handleUseGyForAuthOnly(
	imsi string,
	pReq *protos.CreateSessionRequest,
	staticRuleInstalls []*protos.StaticRuleInstall,
	dynamicRuleInstalls []*protos.DynamicRuleInstall,
	gxCCAInit *gx.CreditControlAnswer) (*protos.CreateSessionResponse, error) {

	gyCCRInit := makeCCRInitWithoutChargingKeys(imsi, pReq)
	gyCCAInit, err := srv.sendSingleCreditRequest(gyCCRInit)
	metrics.ReportCreateGySession(err)
	if err != nil {
		glog.Errorf("Failed to send second single credit request: %s", err)
		return nil, err
	}

	err = validateGyCCAIMSCC(gyCCAInit)
	if err != nil {
		glog.Errorf("MSCC Avp Failure: %s", err)
		return nil, err
	}
	gxOriginHost, gyOriginHost := gxCCAInit.OriginHost, gyCCAInit.OriginHost
	return &protos.CreateSessionResponse{
		StaticRules:   staticRuleInstalls,
		DynamicRules:  dynamicRuleInstalls,
		UsageMonitors: getUsageMonitorsFromCCA_I(imsi, gyOriginHost, gyCCAInit.SessionID, gxCCAInit),
		TgppCtx:       &protos.TgppContext{GxDestHost: gxOriginHost, GyDestHost: gyCCAInit.OriginHost},
		SessionId:     pReq.SessionId,
	}, nil
}

func (srv *CentralSessionController) getChargingKeysFromRuleInstalls(
	staticRuleInstalls []*protos.StaticRuleInstall,
	dynamicRuleInstalls []*protos.DynamicRuleInstall,
) []policydb.ChargingKey {
	staticRuleIDs := funk.Map(staticRuleInstalls, func(s *protos.StaticRuleInstall) string { return s.RuleId }).([]string)
	dynamicRuleDef := funk.Map(dynamicRuleInstalls, func(d *protos.DynamicRuleInstall) *protos.PolicyRule { return d.PolicyRule }).([]*protos.PolicyRule)
	keys := srv.dbClient.GetChargingKeysForRules(staticRuleIDs, dynamicRuleDef)
	return removeDuplicateChargingKeys(keys)
}

func removeDuplicateChargingKeys(keysIn []policydb.ChargingKey) []policydb.ChargingKey {
	keysOut := []policydb.ChargingKey{}
	keyMap := make(map[policydb.ChargingKey]struct{})
	for _, k := range keysIn {
		if _, ok := keyMap[k]; !ok {
			keysOut = append(keysOut, k)
			keyMap[k] = struct{}{}
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

	var (
		gxUpdateResponses []*protos.UsageMonitoringUpdateResponse
		gyUpdateResponses []*protos.CreditUpdateResponse
		gxTgppCtx         = make([]struct {
			sid string
			ctx *protos.TgppContext
		}, 0, len(request.UsageMonitors))
		gyTgppCtx = map[string]*protos.TgppContext{}
	)
	go func() {
		defer wg.Done()
		if srv.cfg.DisableGx {
			return
		}
		requests := getGxUpdateRequestsFromUsage(request.UsageMonitors)
		gxUpdateResponses = srv.sendMultipleGxRequestsWithTimeout(requests, srv.cfg.RequestTimeout)
		for _, mur := range gxUpdateResponses {
			if mur != nil {
				if mur.TgppCtx != nil {
					gxTgppCtx = append(gxTgppCtx, struct {
						sid string
						ctx *protos.TgppContext
					}{mur.Sid, mur.TgppCtx})
				}
			}
		}
	}()
	go func() {
		defer wg.Done()
		if srv.cfg.DisableGy {
			return
		}
		requests := getGyUpdateRequestsFromUsage(request.Updates)
		gyUpdateResponses = srv.sendMultipleGyRequestsWithTimeout(requests, srv.cfg.RequestTimeout)
		for _, cur := range gyUpdateResponses {
			if cur != nil {
				if cur.TgppCtx != nil {
					gyTgppCtx[cur.GetSid()] = cur.TgppCtx
				}
			}
		}
	}()
	wg.Wait()

	// Update destination hosts in all results with common SIDs
	if len(gxTgppCtx) > 0 && len(gyTgppCtx) > 0 {
		for _, pair := range gxTgppCtx {
			if gyCtx, ok := gyTgppCtx[pair.sid]; ok {
				if len(pair.ctx.GxDestHost) > 0 {
					gyCtx.GxDestHost = pair.ctx.GxDestHost
				}
				if len(gyCtx.GyDestHost) > 0 {
					pair.ctx.GyDestHost = gyCtx.GyDestHost
				}
			}
		}
	}

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
		if srv.cfg.DisableGx {
			return
		}
		_, err := srv.sendTerminationGxRequest(request)
		metrics.ReportTerminateGxSession(err)
		if err != nil {
			glog.Errorf("Error sending gx termination: %s", err)
			return
		}

	}()
	go func() {
		defer wg.Done()
		if srv.cfg.DisableGy {
			return
		}
		_, err := srv.sendSingleCreditRequest(getTerminateRequestFromUsage(request))
		metrics.ReportTerminateGySession(err)
		if err != nil {
			glog.Errorf("Error sending gy termination: %s", err)
			return
		}
	}()
	wg.Wait()
	// in the event of any errors on Gx or Gy, the session should regardless be
	// terminated, so there are no errors sent back
	return &protos.SessionTerminateResponse{
		Sid:       request.GetCommonContext().GetSid().GetId(),
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
	if !srv.cfg.DisableGx {
		srv.policyClient.DisableConnections(time.Duration(req.DisablePeriodSecs) * time.Second)
	}
	if !srv.cfg.DisableGy {
		srv.creditClient.DisableConnections(time.Duration(req.DisablePeriodSecs) * time.Second)
	}
	return &orcprotos.Void{}, nil
}

// Enable enables diameter connection creation and gets a connection to the
// diameter server(s). If creation is already enabled and a connection already
// exists, Enable has no effect
func (srv *CentralSessionController) Enable(
	ctx context.Context,
	void *orcprotos.Void,
) (*orcprotos.Void, error) {
	multiError := errors.NewMulti()
	if !srv.cfg.DisableGx {
		err := srv.policyClient.EnableConnections()
		if err != nil {
			multiError.Add(fmt.Errorf("An error occurred while enabling connections; policyClient err: %s", err))
		}
	}
	if !srv.cfg.DisableGy {
		err := srv.creditClient.EnableConnections()
		if err != nil {
			multiError.Add(fmt.Errorf("An error occurred while enabling connections; creditClient err: %s", err))
		}
	}
	return &orcprotos.Void{}, multiError.AsError()
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
			HealthMessage: fmt.Sprintf("Error occurred while retrieving health metrics: %s", err),
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
	if !srv.cfg.DisableGx {
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
	}
	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "Gx metrics appear healthy",
	}
}

func (srv *CentralSessionController) getHealthStatusForGyRequests(failures, total int64) *fegprotos.HealthStatus {
	if !srv.cfg.DisableGy {
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
	}
	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "Gy metrics appear healthy",
	}
}

func checkCreateSessionRequest(req *protos.CreateSessionRequest) error {
	subscriber := req.GetCommonContext().GetSid()
	if subscriber == nil || subscriber.GetId() == "" {
		return fmt.Errorf("Missing Subscriber information on CreateSessionRequest %+v", req)
	}
	if req.SessionId == "" {
		return fmt.Errorf("Missing magma SessionId information on CreateSessionRequest %+v", req)
	}
	return nil
}

// injectOmnipresentRules adds Omnipresent rules and basenames to the gxCCAInit in form of RuleInstallAVP
func injectOmnipresentRules(srv *CentralSessionController, gxCCAInit *gx.CreditControlAnswer) error {
	if gxCCAInit == nil {
		return fmt.Errorf("gxCCAInit is nul")
	}
	omnipresentRuleIDs, omnipresentBaseNames := srv.dbClient.GetOmnipresentRules()
	if len(omnipresentRuleIDs) > 0 || len(omnipresentBaseNames) > 0 {
		glog.V(2).Infof("Adding omnipresent rules %v and omnipresent rule base %v", omnipresentRuleIDs, omnipresentBaseNames)
		omniRuleInstallAVPRuleInstallAVP := &gx.RuleInstallAVP{
			RuleNames:            omnipresentRuleIDs,
			RuleBaseNames:        omnipresentBaseNames,
			RuleDefinitions:      []*gx.RuleDefinition{},
			RuleActivationTime:   nil,
			RuleDeactivationTime: nil,
		}
		gxCCAInit.RuleInstallAVP = append(gxCCAInit.RuleInstallAVP, omniRuleInstallAVPRuleInstallAVP)
	}
	return nil
}
