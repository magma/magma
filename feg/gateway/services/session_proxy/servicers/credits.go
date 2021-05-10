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

// makeCCRInit creates the first update request to send to the
// OCS when a session is established.
func makeCCRInit(
	imsi string,
	pReq *protos.CreateSessionRequest,
	keys []policydb.ChargingKey,
) *gy.CreditControlRequest {
	common := pReq.GetCommonContext()
	request := &gy.CreditControlRequest{
		SessionID:     pReq.GetSessionId(),
		Type:          credit_control.CRTInit,
		IMSI:          imsi,
		RequestNumber: 0,
		UeIPV4:        common.GetUeIpv4(),
		Apn:           common.GetApn(),
		Msisdn:        common.GetMsisdn(),
		RatType:       gy.GetRATType(common.GetRatType()),
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
				request.Qos = &gy.QosRequestInfo{
					ApnAggMaxBitRateDL: lteContext.GetQosInfo().GetApnAmbrDl(),
					ApnAggMaxBitRateUL: lteContext.GetQosInfo().GetApnAmbrUl(),
				}
			}
		}
	} else {
		glog.Warning("No RatSpecificContext is specified")
	}
	request.Credits = makeUsedCreditsForCCRInit(pReq.RequestedUnits, keys)
	return request
}

func makeUsedCreditsForCCRInit(
	requestedUnits *protos.RequestedUnits,
	keys []policydb.ChargingKey) []*gy.UsedCredits {
	usedCredits := make([]*gy.UsedCredits, 0, len(keys))
	for _, key := range keys {
		uc := &gy.UsedCredits{
			RatingGroup:    key.RatingGroup,
			RequestedUnits: getRequestedUnitsOrDefault(requestedUnits),
		}
		if key.ServiceIdTracking {
			sid := key.ServiceIdentifier
			uc.ServiceIdentifier = &sid
		}
		usedCredits = append(usedCredits, uc)
	}
	return usedCredits
}

// TODO: function for backwards compatibility. Delete once older AGW are updated
func getRequestedUnitsOrDefault(requestedUnits *protos.RequestedUnits) *protos.RequestedUnits {
	if requestedUnits == nil {
		return &protos.RequestedUnits{Total: 100000, Tx: 100000, Rx: 100000}
	}
	return requestedUnits
}

// makeCCRInitWithoutChargingKeys creates a CreditControlRequest for an INIT
// message, defaulting the request number to 0 and not including credit usage
func makeCCRInitWithoutChargingKeys(imsi string, pReq *protos.CreateSessionRequest) *gy.CreditControlRequest {
	return makeCCRInit(imsi, pReq, nil)
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
			Success:    false,
			Sid:        imsi,
			SessionId:  request.SessionID,
			ResultCode: answer.ResultCode,
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
		SessionId:   request.SessionID,
		Sid:         imsi,
		ChargingKey: receivedCredit.RatingGroup,
		Credit:      getSingleChargingCreditFromCCA(receivedCredit),
		TgppCtx:     tgppCtx,
		ResultCode:  receivedCredit.ResultCode, //answer.ResultCode is returned in case of general failure
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
	chargingCredit := &protos.ChargingCredit{
		GrantedUnits: credits.GrantedUnits.ToProto(),
		Type:         protos.ChargingCredit_BYTES,
		ValidityTime: credits.ValidityTime,
	}
	if credits.FinalUnitIndication != nil {
		chargingCredit.IsFinal = true
		chargingCredit.FinalAction = protos.ChargingCredit_FinalAction(credits.FinalUnitIndication.FinalAction)
		chargingCredit.RedirectServer = credits.FinalUnitIndication.RedirectServer.ToProto()
		chargingCredit.RestrictRules = credits.FinalUnitIndication.RestrictRules
	}
	return chargingCredit
}

// getUpdateRequestsFromUsage returns a slice of CCRs from usage update protos
func getGyUpdateRequestsFromUsage(updates []*protos.CreditUsageUpdate) []*gy.CreditControlRequest {
	requests := []*gy.CreditControlRequest{}
	for _, update := range updates {
		requests = append(requests, (&gy.CreditControlRequest{}).FromCreditUsageUpdate(update))
	}
	return requests
}

// getTerminateRequestFromUsage returns a slice of CCRs from usage update protos
func getTerminateRequestFromUsage(termination *protos.SessionTerminateRequest) *gy.CreditControlRequest {
	usedCredits := make([]*gy.UsedCredits, 0, len(termination.CreditUsages))
	for _, usage := range termination.CreditUsages {
		usedCredits = append(usedCredits, (&gy.UsedCredits{}).FromCreditUsage(usage))
	}
	common := termination.GetCommonContext()
	return &gy.CreditControlRequest{
		SessionID:               termination.SessionId,
		IMSI:                    credit_control.RemoveIMSIPrefix(common.GetSid().GetId()),
		Apn:                     common.GetApn(),
		RequestNumber:           termination.RequestNumber,
		Credits:                 usedCredits,
		UeIPV4:                  common.GetUeIpv4(),
		Msisdn:                  common.GetMsisdn(),
		SpgwIPV4:                termination.SpgwIpv4,
		Imei:                    termination.Imei,
		PlmnID:                  termination.PlmnId,
		UserLocation:            termination.UserLocation,
		Type:                    credit_control.CRTTerminate,
		RatType:                 gy.GetRATType(common.GetRatType()),
		TgppCtx:                 termination.GetTgppCtx(),
		ChargingCharacteristics: termination.ChargingCharacteristics,
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
