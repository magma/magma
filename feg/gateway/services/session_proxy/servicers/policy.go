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
	"fmt"
	"strconv"
	"time"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/metrics"
	"magma/lte/cloud/go/protos"

	"github.com/golang/glog"
)

// sendInitialGxRequestOrGenerateEmptyResponse generates an empty response in case Gx is disabled.
// otherwise it sends the initial request to PCRF
func (srv *CentralSessionController) sendInitialGxRequestOrGenerateEmptyResponse(imsi string, pReq *protos.CreateSessionRequest) (*gx.CreditControlAnswer, error) {
	if srv.cfg.DisableGx {
		return generateGxLessCCAInit()
	}
	return srv.sendInitialGxRequest(imsi, pReq)
}

// sendInitialGxRequest sends the initial request to PCRF. Returns a response
func (srv *CentralSessionController) sendInitialGxRequest(imsi string, pReq *protos.CreateSessionRequest) (*gx.CreditControlAnswer, error) {
	common := pReq.GetCommonContext()
	ratType := common.GetRatType()
	request := &gx.CreditControlRequest{
		SessionID:      pReq.GetSessionId(),
		Type:           credit_control.CRTInit,
		IMSI:           imsi,
		RequestNumber:  0,
		IPAddr:         common.GetUeIpv4(),
		IPv6Addr:       common.GetUeIpv6(),
		Apn:            common.GetApn(),
		Msisdn:         common.GetMsisdn(),
		RATType:        gx.GetRATType(ratType),
		IPCANType:      gx.GetIPCANType(ratType),
		AccessTimezone: pReq.GetAccessTimezone(),
	}

	if pReq.RatSpecificContext != nil {
		ratSpecific := pReq.GetRatSpecificContext().GetContext()
		switch context := ratSpecific.(type) {
		case *protos.RatSpecificContext_LteContext:
			lteContext := context.LteContext
			request.SpgwIPV4 = lteContext.GetSpgwIpv4()
			request.Imei = lteContext.GetImei()
			request.PlmnID = lteContext.GetPlmnId()
			request.UserLocation = lteContext.GetUserLocation()
			request.ChargingCharacteristics = lteContext.GetChargingCharacteristics()
			if lteContext.GetQosInfo() != nil {
				request.Qos = (&gx.QosRequestInfo{}).FromProtos(lteContext.GetQosInfo())
			}
		case *protos.RatSpecificContext_WlanContext:
			request.HardwareAddr = context.WlanContext.GetMacAddrBinary()
		}
	} else {
		glog.Warning("No RatSpecificContext is specified")
	}

	return getGxAnswerOrError(request, srv.policyClient, srv.cfg.PCRFConfig, srv.cfg.RequestTimeout)
}

// generateGxLessCCAInit generates a creditControlAnswer for the case Gx is disabled.
func generateGxLessCCAInit() (*gx.CreditControlAnswer, error) {
	return &gx.CreditControlAnswer{
		ResultCode:             uint32(diameter.SuccessCode),
		ExperimentalResultCode: 0,
		SessionID:              "",
		OriginHost:             "gx-less.magma.com",
		RequestNumber:          0,
		RuleInstallAVP:         []*gx.RuleInstallAVP{},
		RuleRemoveAVP:          []*gx.RuleRemoveAVP{},
		UsageMonitors:          []*gx.UsageMonitoringInfo{},
		EventTriggers:          []gx.EventTrigger{},
		RevalidationTime:       nil,
	}, nil
}

func (srv *CentralSessionController) sendTerminationGxRequest(pRequest *protos.SessionTerminateRequest) (*gx.CreditControlAnswer, error) {
	reports := make([]*gx.UsageReport, 0, len(pRequest.MonitorUsages))
	for _, update := range pRequest.MonitorUsages {
		reports = append(reports, (&gx.UsageReport{}).FromUsageMonitorUpdate(update))
	}
	common := pRequest.GetCommonContext()
	request := &gx.CreditControlRequest{
		SessionID:               pRequest.SessionId,
		Type:                    credit_control.CRTTerminate,
		IMSI:                    credit_control.RemoveIMSIPrefix(common.GetSid().GetId()),
		RequestNumber:           pRequest.RequestNumber,
		IPAddr:                  common.GetUeIpv4(),
		UsageReports:            reports,
		RATType:                 gx.GetRATType(common.GetRatType()),
		IPCANType:               gx.GetIPCANType(common.GetRatType()),
		TgppCtx:                 pRequest.GetTgppCtx(),
		ChargingCharacteristics: pRequest.ChargingCharacteristics,
	}
	return getGxAnswerOrError(request, srv.policyClient, srv.cfg.PCRFConfig, srv.cfg.RequestTimeout)
}

func getGxAnswerOrError(
	request *gx.CreditControlRequest,
	policyClient gx.PolicyClient,
	pcrfConfig *diameter.DiameterServerConfig,
	requestTimeout time.Duration,
) (*gx.CreditControlAnswer, error) {
	done := make(chan interface{}, 1)
	err := policyClient.SendCreditControlRequest(pcrfConfig, done, request)
	if err != nil {
		return nil, err
	}
	select {
	case resp := <-done:
		answer := resp.(*gx.CreditControlAnswer)
		metrics.GxResultCodes.WithLabelValues(strconv.FormatUint(uint64(answer.ResultCode), 10)).Inc()
		if answer.ResultCode != diameter.SuccessCode {
			return nil, fmt.Errorf(
				"Received unsuccessful result code from PCRF, ResultCode: %d, ExperimentalResultCode: %d",
				answer.ResultCode, answer.ExperimentalResultCode)
		}
		return answer, nil
	case <-time.After(requestTimeout):
		metrics.GxTimeouts.Inc()
		policyClient.IgnoreAnswer(request)
		return nil, fmt.Errorf("CCA wait timeout for session: %s after %s", request.SessionID, requestTimeout.String())
	}
}

func getUsageMonitorsFromCCA_I(
	imsi, sessionID, gyOriginHost string, gxCCAInit *gx.CreditControlAnswer) []*protos.UsageMonitoringUpdateResponse {
	monitors := make([]*protos.UsageMonitoringUpdateResponse, 0, len(gxCCAInit.UsageMonitors))
	// If there is a message wide revalidation time, apply it to every Usage Monitor
	triggers, revalidationTime := gx.GetEventTriggersRelatedInfo(gxCCAInit.EventTriggers, gxCCAInit.RevalidationTime)
	tgppCtx := &protos.TgppContext{GxDestHost: gxCCAInit.OriginHost, GyDestHost: gyOriginHost}

	for _, monitor := range gxCCAInit.UsageMonitors {
		monitors = append(monitors, &protos.UsageMonitoringUpdateResponse{
			Credit:           monitor.ToUsageMonitoringCredit(),
			SessionId:        sessionID,
			TgppCtx:          tgppCtx,
			Sid:              credit_control.AddIMSIPrefix(imsi),
			Success:          true,
			EventTriggers:    triggers,
			RevalidationTime: revalidationTime,
		})
	}
	return monitors
}

// getGxUpdateRequestsFromUsage returns a slice of CCRs from usage update protos
func getGxUpdateRequestsFromUsage(updates []*protos.UsageMonitoringUpdateRequest) []*gx.CreditControlRequest {
	requests := []*gx.CreditControlRequest{}
	for _, update := range updates {
		requests = append(requests, (&gx.CreditControlRequest{}).FromUsageMonitorUpdate(update))
	}
	return requests
}

// sendMultipleGxRequestsWithTimeout sends a batch of update requests to the PCRF
// and returns a response for every request, even during timeouts.
func (srv *CentralSessionController) sendMultipleGxRequestsWithTimeout(
	requests []*gx.CreditControlRequest,
	timeoutDuration time.Duration,
) []*protos.UsageMonitoringUpdateResponse {
	done := make(chan interface{}, len(requests))
	srv.sendGxUpdateRequestsToConnections(requests, done)
	return srv.waitForGxResponses(requests, done, timeoutDuration)
}

// sendGxUpdateRequestsToConnections sends batches of requests to PCRF's
func (srv *CentralSessionController) sendGxUpdateRequestsToConnections(
	requests []*gx.CreditControlRequest,
	done chan interface{},
) {
	sendErrors := []error{}
	for _, request := range requests {
		err := srv.policyClient.SendCreditControlRequest(srv.cfg.PCRFConfig, done, request)
		if err != nil {
			sendErrors = append(sendErrors, err)
			metrics.PcrfCcrUpdateSendFailures.Inc()
		} else {
			metrics.PcrfCcrUpdateRequests.Inc()
		}
	}
	if len(sendErrors) > 0 {
		go func() {
			for _, err := range sendErrors {
				done <- err
			}
		}()
	}
}

// waitForGxResponses waits for CreditControlAnswers on the done channel. It stops
// no matter how many responses it has gotten within the given timeout. If any
// responses are not received, it manually adds them and returns. It is ensured
// that the number of requests matches the number of responses
func (srv *CentralSessionController) waitForGxResponses(
	requests []*gx.CreditControlRequest,
	done chan interface{},
	timeoutDuration time.Duration,
) []*protos.UsageMonitoringUpdateResponse {
	requestMap := createGxRequestKeyMap(requests)
	responses := []*protos.UsageMonitoringUpdateResponse{}
	timeout := time.After(timeoutDuration)
	numResponses := len(requests)
loop:
	for i := 0; i < numResponses; i++ {
		select {
		case resp := <-done:
			switch ans := resp.(type) {
			case error:
				glog.Errorf("Error encountered in request: %s", ans.Error())
			case *gx.CreditControlAnswer:
				metrics.GxResultCodes.WithLabelValues(strconv.FormatUint(uint64(ans.ResultCode), 10)).Inc()
				metrics.UpdateGxRecentRequestMetrics(nil)
				key := credit_control.GetRequestKey(credit_control.Gx, ans.SessionID, ans.RequestNumber)
				newResponse := srv.getSingleUsageMonitorResponseFromCCA(ans, requestMap[key])
				responses = append(responses, newResponse)
				// satisfied request, remove
				delete(requestMap, key)
			default:
				glog.Errorf("Unknown type sent to CCA done channel")
			}
		case <-timeout:
			glog.Errorf("Timed out receiving responses from PCRF\n")
			// tell client to ignore answers to timed out responses
			srv.ignoreGxTimedOutRequests(requestMap)
			// add missing responses
			break loop
		}
	}
	responses = addMissingGxResponses(responses, requestMap)
	return responses
}

// createRequestKeyMap takes a list of requests and returns a map of request key
// (SESSIONID-REQUESTNUM) to request. This is used to identify responses as they
// come through and ensure every request is responded to
func createGxRequestKeyMap(requests []*gx.CreditControlRequest) map[credit_control.RequestKey]*gx.CreditControlRequest {
	requestMap := make(map[credit_control.RequestKey]*gx.CreditControlRequest)
	for _, request := range requests {
		requestKey := credit_control.GetRequestKey(credit_control.Gx, request.SessionID, request.RequestNumber)
		requestMap[requestKey] = request
	}
	return requestMap
}

// ignoreGxTimedOutRequests tells the gx client to ignore any requests that have
// timed out. This ensures the gx client does not leak request trackings
func (srv *CentralSessionController) ignoreGxTimedOutRequests(
	leftoverRequests map[credit_control.RequestKey]*gx.CreditControlRequest,
) {
	for _, ccr := range leftoverRequests {
		metrics.GxTimeouts.Inc()
		srv.policyClient.IgnoreAnswer(ccr)
	}
}

// addMissingGxResponses looks through leftoverRequests to see if there are any
// requests that did not receive responses, and manually adds an errored
// response to responses
func addMissingGxResponses(
	responses []*protos.UsageMonitoringUpdateResponse,
	leftoverRequests map[credit_control.RequestKey]*gx.CreditControlRequest,
) []*protos.UsageMonitoringUpdateResponse {
	for _, ccr := range leftoverRequests {
		responses = append(responses, &protos.UsageMonitoringUpdateResponse{
			Success:   false,
			SessionId: ccr.SessionID,
			Sid:       credit_control.AddIMSIPrefix(ccr.IMSI),
			Credit: &protos.UsageMonitoringCredit{
				MonitoringKey: ccr.UsageReports[0].MonitoringKey,
				Level:         protos.MonitoringLevel(ccr.UsageReports[0].Level),
			},
		})
		metrics.UpdateGxRecentRequestMetrics(fmt.Errorf("Gx update failure"))
	}
	return responses
}

// getSingleUsageMonitorResponseFromCCA creates a UsageMonitoringUpdateResponse proto from a CCA
func (srv *CentralSessionController) getSingleUsageMonitorResponseFromCCA(
	answer *gx.CreditControlAnswer, request *gx.CreditControlRequest) *protos.UsageMonitoringUpdateResponse {

	staticRules, dynamicRules := gx.ParseRuleInstallAVPs(
		srv.dbClient,
		answer.RuleInstallAVP,
	)
	rulesToRemove := gx.ParseRuleRemoveAVPs(
		srv.dbClient,
		answer.RuleRemoveAVP,
	)
	tgppCtx := request.TgppCtx
	if len(answer.OriginHost) > 0 {
		if tgppCtx == nil {
			tgppCtx = new(protos.TgppContext)
		}
		tgppCtx.GxDestHost = answer.OriginHost
	}
	res := &protos.UsageMonitoringUpdateResponse{
		Success:               answer.ResultCode == diameter.SuccessCode || answer.ResultCode == 0,
		SessionId:             request.SessionID,
		Sid:                   credit_control.AddIMSIPrefix(request.IMSI),
		ResultCode:            answer.ResultCode,
		RulesToRemove:         rulesToRemove,
		StaticRulesToInstall:  staticRules,
		DynamicRulesToInstall: dynamicRules,
		TgppCtx:               tgppCtx,
	}
	if len(answer.UsageMonitors) != 0 {
		res.Credit = answer.UsageMonitors[0].ToUsageMonitoringCredit()
	} else if len(request.UsageReports) != 0 {
		glog.Infof("No usage monitor response in CCA for subscriber %s", request.IMSI)
		res.Credit = &protos.UsageMonitoringCredit{
			Action:        protos.UsageMonitoringCredit_DISABLE,
			MonitoringKey: request.UsageReports[0].MonitoringKey,
			Level:         protos.MonitoringLevel(request.UsageReports[0].Level)}

	}

	res.EventTriggers, res.RevalidationTime = gx.GetEventTriggersRelatedInfo(
		answer.EventTriggers,
		answer.RevalidationTime,
	)
	return res
}
