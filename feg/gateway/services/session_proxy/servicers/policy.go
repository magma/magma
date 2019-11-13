/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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

func (srv *CentralSessionController) sendInitialGxRequest(imsi string, pReq *protos.CreateSessionRequest) (*gx.CreditControlAnswer, error) {

	var qos *gx.QosRequestInfo
	if pReq.GetQosInfo() != nil {
		qos = &gx.QosRequestInfo{
			ApnAggMaxBitRateDL: pReq.GetQosInfo().GetApnAmbrDl(),
			ApnAggMaxBitRateUL: pReq.GetQosInfo().GetApnAmbrUl(),
			QosClassIdentifier: pReq.GetQosInfo().GetQosClassId(),
			PriLevel:           pReq.GetQosInfo().GetPriorityLevel(),
			PreCapability:      pReq.GetQosInfo().GetPreemptionCapability(),
			PreVulnerability:   pReq.GetQosInfo().GetPreemptionVulnerability(),
		}
	}

	request := &gx.CreditControlRequest{
		SessionID:     pReq.SessionId,
		Type:          credit_control.CRTInit,
		IMSI:          imsi,
		RequestNumber: 1,
		IPAddr:        pReq.UeIpv4,
		SpgwIPV4:      pReq.SpgwIpv4,
		Apn:           pReq.Apn,
		Msisdn:        pReq.Msisdn,
		Imei:          pReq.Imei,
		PlmnID:        pReq.PlmnId,
		UserLocation:  pReq.UserLocation,
		GcID:          pReq.GcId,
		Qos:           qos,
		HardwareAddr:  pReq.HardwareAddr,
		RATType:       getRATType(pReq.RatType),
		IPCANType:     getIPCANType(pReq.RatType),
	}

	return getGxAnswerOrError(request, srv.policyClient, srv.cfg.PCRFConfig, srv.cfg.RequestTimeout)
}

func getRATType(pRATType protos.RATType) credit_control.RATType {
	switch pRATType {
	case protos.RATType_TGPP_LTE:
		return credit_control.RAT_EUTRAN
	case protos.RATType_TGPP_WLAN:
		return credit_control.RAT_WLAN
	default:
		return credit_control.RAT_EUTRAN
	}
}

// Since we don't specify the IP CAN type at session initialization, and we
// only support WLAN and EUTRAN, we will infer the IP CAN type from RAT type.
func getIPCANType(pRATType protos.RATType) credit_control.IPCANType {
	switch pRATType {
	case protos.RATType_TGPP_LTE:
		return credit_control.IPCAN_3GPP
	case protos.RATType_TGPP_WLAN:
		return credit_control.IPCAN_Non3GPP
	default:
		return credit_control.IPCAN_Non3GPP
	}
}

func (srv *CentralSessionController) sendTerminationGxRequest(pRequest *protos.SessionTerminateRequest) (*gx.CreditControlAnswer, error) {
	reports := make([]*gx.UsageReport, 0, len(pRequest.MonitorUsages))
	for _, update := range pRequest.MonitorUsages {
		reports = append(reports, getGxUsageReportFromUsageUpdate(update))
	}
	request := &gx.CreditControlRequest{
		SessionID:     pRequest.SessionId,
		Type:          credit_control.CRTTerminate,
		IMSI:          removeSidPrefix(pRequest.Sid),
		RequestNumber: pRequest.RequestNumber,
		IPAddr:        pRequest.UeIpv4,
		UsageReports:  reports,
		RATType:       getRATType(pRequest.RatType),
		IPCANType:     getIPCANType(pRequest.RatType),
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

func getPolicyRulesFromDefinitions(ruleDefs []*gx.RuleDefinition) []*protos.PolicyRule {
	policyRules := make([]*protos.PolicyRule, 0, len(ruleDefs))
	for _, def := range ruleDefs {
		policyRules = append(policyRules, def.ToProto())
	}
	return policyRules
}

func getUsageMonitorsFromCCA(imsi string, sessionID string, gxCCA *gx.CreditControlAnswer) []*protos.UsageMonitoringUpdateResponse {
	monitors := make([]*protos.UsageMonitoringUpdateResponse, 0, len(gxCCA.UsageMonitors))
	for _, monitor := range gxCCA.UsageMonitors {
		monitors = append(monitors, &protos.UsageMonitoringUpdateResponse{
			Credit:    gx.GetUsageMonitorCreditFromAVP(monitor),
			SessionId: sessionID,
			Sid:       addSidPrefix(imsi),
			Success:   true,
		})
	}
	return monitors
}

// getGxUpdateRequestsFromUsage returns a slice of CCRs from usage update protos
func getGxUpdateRequestsFromUsage(updates []*protos.UsageMonitoringUpdateRequest) []*gx.CreditControlRequest {
	requests := []*gx.CreditControlRequest{}
	for _, update := range updates {
		requests = append(requests, &gx.CreditControlRequest{
			SessionID:     update.SessionId,
			RequestNumber: update.RequestNumber,
			Type:          credit_control.CRTUpdate,
			IMSI:          removeSidPrefix(update.Sid),
			IPAddr:        update.UeIpv4,
			HardwareAddr:  update.HardwareAddr,
			UsageReports:  []*gx.UsageReport{getGxUsageReportFromUsageUpdate(update.Update)},
			RATType:       getRATType(update.RatType),
			IPCANType:     getIPCANType(update.RatType),
		})
	}
	return requests
}

func getGxUsageReportFromUsageUpdate(update *protos.UsageMonitorUpdate) *gx.UsageReport {
	return &gx.UsageReport{
		MonitoringKey: update.MonitoringKey,
		Level:         gx.MonitoringLevel(update.Level),
		InputOctets:   update.BytesTx,
		OutputOctets:  update.BytesRx, // receive == output
		TotalOctets:   update.BytesTx + update.BytesRx,
	}
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
			Sid:       addSidPrefix(ccr.IMSI),
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
func (srv *CentralSessionController) getSingleUsageMonitorResponseFromCCA(answer *gx.CreditControlAnswer, request *gx.CreditControlRequest) *protos.UsageMonitoringUpdateResponse {
	res := &protos.UsageMonitoringUpdateResponse{
		Success:       answer.ResultCode == diameter.SuccessCode || answer.ResultCode == 0,
		SessionId:     request.SessionID,
		Sid:           addSidPrefix(request.IMSI),
		ResultCode:    answer.ResultCode,
		RulesToRemove: srv.getRulesToRemoveFromAVP(answer.RuleRemoveAVP),
	}
	if len(answer.UsageMonitors) == 0 {
		glog.Infof("No usage monitor response in CCA for subscriber %s", request.IMSI)
		res.Credit =
			&protos.UsageMonitoringCredit{
				Action:        protos.UsageMonitoringCredit_DISABLE,
				MonitoringKey: request.UsageReports[0].MonitoringKey,
				Level:         protos.MonitoringLevel(request.UsageReports[0].Level)}
	} else {
		res.Credit = gx.GetUsageMonitorCreditFromAVP(answer.UsageMonitors[0])
	}

	res.EventTriggers, res.RevalidationTime = gx.GetEventTriggersRelatedInfo(
		answer.EventTriggers,
		answer.RevalidationTime,
	)
	return res
}

func (srv *CentralSessionController) getRulesToRemoveFromAVP(rulesToRemoveAVP []*gx.RuleRemoveAVP) []string {
	var ruleNames []string
	for _, rule := range rulesToRemoveAVP {
		ruleNames = append(ruleNames, rule.RuleNames...)
		if len(rule.RuleBaseNames) > 0 {
			ruleNames = append(ruleNames, srv.dbClient.GetRuleIDsForBaseNames(rule.RuleBaseNames)...)
		}
	}
	return ruleNames
}
