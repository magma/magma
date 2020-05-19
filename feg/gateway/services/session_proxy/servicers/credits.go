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
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/session_proxy/metrics"
	"magma/lte/cloud/go/protos"

	"github.com/golang/glog"
)

// sendSingleCreditRequest sends a CCR message through the gy client
// and waits for a response based on the grpc server's set timeout
func (srv *CentralSessionController) sendSingleCreditRequest(request *gy.CreditControlRequest) (*gy.CreditControlAnswer, error) {
	done := make(chan interface{}, 1)
	err := srv.creditClient.SendCreditControlRequest(srv.cfg.OCSConfig, done, request)
	if err != nil {
		return nil, err
	}

	select {
	case resp := <-done:
		answer := resp.(*gy.CreditControlAnswer)
		metrics.GyResultCodes.WithLabelValues(strconv.FormatUint(uint64(answer.ResultCode), 10)).Inc()
		if answer.ResultCode != diameter.SuccessCode {
			return nil, fmt.Errorf("Received unsuccessful result code from OCS: %d for session: %s, IMSI: %s",
				answer.ResultCode, request.SessionID, request.IMSI)
		}
		return answer, nil
	case <-time.After(srv.cfg.RequestTimeout):
		metrics.GyTimeouts.Inc()
		srv.creditClient.IgnoreAnswer(request)
		return nil, fmt.Errorf("Did not receive Gy CCA for session: %s, IMSI: %s", request.SessionID, request.IMSI)
	}
}

// getCCRInitRequest creates a CreditControlRequest for an INIT message,
// defaulting the request number to 0 and not including credit usage
func getCCRInitRequest(
	imsi string,
	pReq *protos.CreateSessionRequest,
) *gy.CreditControlRequest {

	var qos *gy.QosRequestInfo

	if pReq.GetQosInfo() != nil {
		qos = &gy.QosRequestInfo{
			ApnAggMaxBitRateDL: pReq.GetQosInfo().GetApnAmbrDl(),
			ApnAggMaxBitRateUL: pReq.GetQosInfo().GetApnAmbrUl(),
		}
	}

	return &gy.CreditControlRequest{
		SessionID:     pReq.SessionId,
		RequestNumber: 0,
		IMSI:          imsi,
		UeIPV4:        pReq.UeIpv4,
		SpgwIPV4:      pReq.SpgwIpv4,
		Apn:           pReq.Apn,
		Msisdn:        pReq.Msisdn,
		Imei:          pReq.Imei,
		PlmnID:        pReq.PlmnId,
		UserLocation:  pReq.UserLocation,
		GcID:          pReq.GcId,
		Qos:           qos,
		Type:          credit_control.CRTInit,
		RatType:       gy.GetRATType(pReq.GetRatType()),
	}
}

// getCCRInitialUpdateRequest creates the first update request to send to the
// OCS when a session is established.
func getCCRInitialCreditRequest(
	imsi string,
	pReq *protos.CreateSessionRequest,
	keys []policydb.ChargingKey,
) *gy.CreditControlRequest {
	var msgType credit_control.CreditRequestType
	var qos *gy.QosRequestInfo

	if pReq.GetQosInfo() != nil {
		qos = &gy.QosRequestInfo{
			ApnAggMaxBitRateDL: pReq.GetQosInfo().GetApnAmbrDl(),
			ApnAggMaxBitRateUL: pReq.GetQosInfo().GetApnAmbrUl(),
		}
	}

	msgType = credit_control.CRTInit
	usedCredits := make([]*gy.UsedCredits, 0, len(keys))
	for _, key := range keys {
		uc := &gy.UsedCredits{RatingGroup: key.RatingGroup}
		if key.ServiceIdTracking {
			sid := key.ServiceIdentifier
			uc.ServiceIdentifier = &sid
		}
		usedCredits = append(usedCredits, uc)
	}
	return &gy.CreditControlRequest{
		SessionID:     pReq.SessionId,
		RequestNumber: 0,
		IMSI:          imsi,
		UeIPV4:        pReq.UeIpv4,
		SpgwIPV4:      pReq.SpgwIpv4,
		Apn:           pReq.Apn,
		Msisdn:        pReq.Msisdn,
		Imei:          pReq.Imei,
		PlmnID:        pReq.PlmnId,
		UserLocation:  pReq.UserLocation,
		GcID:          pReq.GcId,
		Qos:           qos,
		Credits:       usedCredits,
		Type:          msgType,
		RatType:       gy.GetRATType(pReq.GetRatType()),
	}
}

// sendMultipleRequestsWithTimeout sends a batch of update requests to the OCS
// and returns a response for every request, even during timeouts.
// This method takes in a slice of requests, groups them by shared OCS
// connection, and sends them through.
func (srv *CentralSessionController) sendMultipleGyRequestsWithTimeout(
	requests []*gy.CreditControlRequest,
	timeoutDuration time.Duration,
) []*protos.CreditUpdateResponse {
	done := make(chan interface{}, len(requests))
	srv.sendGyUpdateRequestsToConnections(requests, done)
	return srv.waitForGyResponses(requests, done, timeoutDuration)
}

// sendGyUpdateRequestsToConnections sends batches of requests to OCS's
func (srv *CentralSessionController) sendGyUpdateRequestsToConnections(
	requests []*gy.CreditControlRequest,
	done chan interface{},
) {
	sendErrors := []error{}
	for _, request := range requests {
		err := srv.creditClient.SendCreditControlRequest(srv.cfg.OCSConfig, done, request)
		if err != nil {
			sendErrors = append(sendErrors, err)
			metrics.OcsCcrUpdateSendFailures.Inc()
		} else {
			metrics.OcsCcrUpdateRequests.Inc()
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

// waitForGyResponses waits for CreditControlAnswers on the done channel. It stops
// no matter how many responses it has gotten within the given timeout. If any
// responses are not received, it manually adds them and returns. It is ensured
// that the number of requests matches the number of responses
func (srv *CentralSessionController) waitForGyResponses(
	requests []*gy.CreditControlRequest,
	done chan interface{},
	timeoutDuration time.Duration,
) []*protos.CreditUpdateResponse {
	requestMap := createRequestKeyMap(requests)
	responses := []*protos.CreditUpdateResponse{}
	timeout := time.After(timeoutDuration) // TODO constant
	numResponses := len(requests)
loop:
	for i := 0; i < numResponses; i++ {
		select {
		case resp := <-done:
			switch ans := resp.(type) {
			case error:
				glog.Errorf("Error encountered in request: %s", ans.Error())
			case *gy.CreditControlAnswer:
				metrics.GyResultCodes.WithLabelValues(strconv.FormatUint(uint64(ans.ResultCode), 10)).Inc()
				metrics.UpdateGyRecentRequestMetrics(nil)
				key := credit_control.GetRequestKey(credit_control.Gy, ans.SessionID, ans.RequestNumber)
				newResponse := getSingleCreditResponseFromCCA(ans, requestMap[key])
				responses = append(responses, newResponse)
				// satisfied request, remove
				delete(requestMap, key)
			default:
				glog.Errorf("Unknown type sent to CCA done channel")
			}
		case <-timeout:
			glog.Errorf("Timed out receiving responses from OCS\n")
			// tell client to ignore answers to timed out responses
			srv.ignoreTimedOutRequests(requestMap)
			// add missing responses
			break loop
		}
	}
	responses = addMissingResponses(responses, requestMap)
	return responses
}

// createRequestKeyMap takes a list of requests and returns a map of request key
// (SESSIONID-REQUESTNUM) to request. This is used to identify responses as they
// come through and ensure every request is responded to
func createRequestKeyMap(requests []*gy.CreditControlRequest) map[credit_control.RequestKey]*gy.CreditControlRequest {
	requestMap := make(map[credit_control.RequestKey]*gy.CreditControlRequest)
	for _, request := range requests {
		requestKey := credit_control.GetRequestKey(credit_control.Gy, request.SessionID, request.RequestNumber)
		requestMap[requestKey] = request
	}
	return requestMap
}

// ignoreTimedOutRequests tells the gy client to ignore any requests that have
// timed out. This ensures the gy client does not leak request trackings
func (srv *CentralSessionController) ignoreTimedOutRequests(
	leftoverRequests map[credit_control.RequestKey]*gy.CreditControlRequest,
) {
	for _, ccr := range leftoverRequests {
		metrics.GyTimeouts.Inc()
		srv.creditClient.IgnoreAnswer(ccr)
	}
}

// addMissingResponses looks through leftoverRequests to see if there are any
// requests that did not receive responses, and manually adds an errored
// response to responses
func addMissingResponses(
	responses []*protos.CreditUpdateResponse,
	leftoverRequests map[credit_control.RequestKey]*gy.CreditControlRequest,
) []*protos.CreditUpdateResponse {
	for _, ccr := range leftoverRequests {
		resp := &protos.CreditUpdateResponse{
			Success:     false,
			Sid:         credit_control.AddIMSIPrefix(ccr.IMSI),
			ChargingKey: ccr.Credits[0].RatingGroup,
		}
		if ccr.Credits[0].ServiceIdentifier != nil {
			resp.ServiceIdentifier = &protos.ServiceIdentifier{Value: *ccr.Credits[0].ServiceIdentifier}
		}
		responses = append(responses, resp)
		metrics.UpdateGyRecentRequestMetrics(fmt.Errorf("Gy update failure"))
	}
	return responses
}

// getSingleCreditResponseFromCCA creates a CreditUpdateResponse proto from a CCA
func getSingleCreditResponseFromCCA(
	answer *gy.CreditControlAnswer,
	request *gy.CreditControlRequest,
) *protos.CreditUpdateResponse {
	success := answer.ResultCode == diameter.SuccessCode
	imsi := credit_control.AddIMSIPrefix(request.IMSI)
	if len(answer.Credits) == 0 {
		return &protos.CreditUpdateResponse{
			Success: false,
			Sid:     imsi,
		}
	}
	receivedCredit := answer.Credits[0]
	msccSuccess := receivedCredit.ResultCode == diameter.SuccessCode || receivedCredit.ResultCode == 0 // 0: not set
	tgppCtx := request.TgppCtx
	if len(answer.OriginHost) > 0 {
		if tgppCtx == nil {
			tgppCtx = new(protos.TgppContext)
		}
		tgppCtx.GyDestHost = answer.OriginHost
	}
	res := &protos.CreditUpdateResponse{
		Success:     success && msccSuccess,
		Sid:         imsi,
		ChargingKey: receivedCredit.RatingGroup,
		Credit:      getSingleChargingCreditFromCCA(receivedCredit),
		TgppCtx:     tgppCtx,
		ResultCode:  answer.ResultCode,
	}

	if receivedCredit.ServiceIdentifier != nil {
		res.ServiceIdentifier = &protos.ServiceIdentifier{Value: *receivedCredit.ServiceIdentifier}
	}
	return res
}

func getInitialCreditResponsesFromCCA(request *gy.CreditControlRequest, answer *gy.CreditControlAnswer) []*protos.CreditUpdateResponse {
	responses := make([]*protos.CreditUpdateResponse, 0, len(answer.Credits))
	tgppCtx := request.TgppCtx
	if len(answer.OriginHost) > 0 {
		if tgppCtx == nil {
			tgppCtx = new(protos.TgppContext)
		}
		tgppCtx.GyDestHost = answer.OriginHost
	}
	for _, credit := range answer.Credits {
		success := credit.ResultCode == diameter.SuccessCode || credit.ResultCode == 0
		response := &protos.CreditUpdateResponse{
			Success:     success,
			Sid:         credit_control.AddIMSIPrefix(request.IMSI),
			ChargingKey: credit.RatingGroup,
			Credit:      getSingleChargingCreditFromCCA(credit),
			ResultCode:  credit.ResultCode,
			TgppCtx:     tgppCtx,
		}
		if credit.ServiceIdentifier != nil {
			response.ServiceIdentifier = &protos.ServiceIdentifier{Value: *credit.ServiceIdentifier}
		}
		responses = append(responses, response)
	}
	return responses
}

// getSingleChargingCreditFromCCA returns a ChargingCredit proto from received
// credits over gy
func getSingleChargingCreditFromCCA(
	credits *gy.ReceivedCredits,
) *protos.ChargingCredit {
	return &protos.ChargingCredit{
		GrantedUnits:   credits.GrantedUnits.ToProto(),
		Type:           protos.ChargingCredit_BYTES,
		ValidityTime:   credits.ValidityTime,
		IsFinal:        credits.IsFinal,
		FinalAction:    protos.ChargingCredit_FinalAction(credits.FinalAction),
		RedirectServer: credits.RedirectServer.ToProto(),
	}
}

// getUpdateRequestsFromUsage returns a slice of CCRs from usage update protos
func getGyUpdateRequestsFromUsage(updates []*protos.CreditUsageUpdate) []*gy.CreditControlRequest {
	requests := []*gy.CreditControlRequest{}
	for _, update := range updates {
		requests = append(requests, &gy.CreditControlRequest{
			SessionID:     update.SessionId,
			RequestNumber: update.RequestNumber,
			IMSI:          credit_control.RemoveIMSIPrefix(update.Sid),
			Msisdn:        update.Msisdn,
			UeIPV4:        update.UeIpv4,
			SpgwIPV4:      update.SpgwIpv4,
			Apn:           update.Apn,
			Imei:          update.Imei,
			PlmnID:        update.PlmnId,
			UserLocation:  update.UserLocation,
			Type:          credit_control.CRTUpdate,
			Credits: []*gy.UsedCredits{&gy.UsedCredits{
				RatingGroup:  update.Usage.ChargingKey,
				InputOctets:  update.Usage.BytesTx, // transmit == input
				OutputOctets: update.Usage.BytesRx, // receive == output
				TotalOctets:  update.Usage.BytesTx + update.Usage.BytesRx,
				Type:         gy.UsedCreditsType(update.Usage.Type),
			}},
			RatType: gy.GetRATType(update.GetRatType()),
			TgppCtx: update.GetTgppCtx(),
		})
	}
	return requests
}

// getTerminateRequestFromUsage returns a slice of CCRs from usage update protos
func getTerminateRequestFromUsage(termination *protos.SessionTerminateRequest) *gy.CreditControlRequest {
	usedCredits := make([]*gy.UsedCredits, 0, len(termination.CreditUsages))
	for _, usage := range termination.CreditUsages {
		usedCredits = append(usedCredits, (&gy.UsedCredits{}).FromCreditUsage(usage))
	}
	return &gy.CreditControlRequest{
		SessionID:     termination.SessionId,
		IMSI:          credit_control.RemoveIMSIPrefix(termination.Sid),
		Apn:           termination.Apn,
		RequestNumber: termination.RequestNumber,
		Credits:       usedCredits,
		UeIPV4:        termination.UeIpv4,
		Msisdn:        termination.Msisdn,
		SpgwIPV4:      termination.SpgwIpv4,
		Imei:          termination.Imei,
		PlmnID:        termination.PlmnId,
		UserLocation:  termination.UserLocation,
		Type:          credit_control.CRTTerminate,
		RatType:       gy.GetRATType(termination.GetRatType()),
		TgppCtx:       termination.GetTgppCtx(),
	}
}

func validateGyCCAIMSCC(gyCCAInit *gy.CreditControlAnswer) error {
	/* Here we need to go through the result codes within MSCC received */

	for _, credit := range gyCCAInit.Credits {
		switch credit.ResultCode {
		case diameter.SuccessCode:
			{
				glog.V(2).Infof("MSCC Avp Result code %v for Rating group %v",
					credit.ResultCode, credit.RatingGroup)
			}
		case diameter.DiameterCreditLimitReached:
			{
				glog.V(2).Infof("MSCC Avp Result code %v for Rating group %v. Subscriber out of credit on OCS",
					credit.ResultCode, credit.RatingGroup)
			}
		default:
			{
				return fmt.Errorf(
					"Received unsuccessful result code from OCS, ResultCode: %d, Rating Group: %d",
					credit.ResultCode, credit.RatingGroup)
			}
		}
	}
	return nil
}
