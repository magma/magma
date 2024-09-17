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

package gy

import (
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/lte/cloud/go/protos"
)

func (redirectServer *RedirectServer) ToProto() *protos.RedirectServer {
	if redirectServer == nil {
		return &protos.RedirectServer{}
	}
	return &protos.RedirectServer{
		RedirectAddressType:   protos.RedirectServer_RedirectAddressType(redirectServer.RedirectAddressType),
		RedirectServerAddress: redirectServer.RedirectServerAddress,
	}
}

func (credits *UsedCredits) FromCreditUsage(usage *protos.CreditUsage) *UsedCredits {
	credits.RatingGroup = usage.ChargingKey
	credits.InputOctets = usage.BytesTx  // transmit == input
	credits.OutputOctets = usage.BytesRx // receive == output
	credits.TotalOctets = usage.BytesTx + usage.BytesRx
	credits.Type = UsedCreditsType(usage.Type)
	credits.ServiceIdentifier = fromServiceIdentifier(usage.ServiceIdentifier)
	return credits
}

// FromCreditUsageUpdates returns a slice of CCRs from usage update protos
// It merges updates from same session into one single request
func FromCreditUsageUpdates(updates []*protos.CreditUsageUpdate) []*CreditControlRequest {
	updatesPerSession := make(map[string][]*protos.CreditUsageUpdate)

	// sort updates per session
	for _, update := range updates {
		updatesPerSession[update.SessionId] = append(updatesPerSession[update.SessionId], update)
	}

	// merge updates for the same sessions
	requests := []*CreditControlRequest{}
	for _, listUpdates := range updatesPerSession {
		firstUpdate := listUpdates[0]
		request := &CreditControlRequest{}
		common := firstUpdate.GetCommonContext()
		request.SessionID = firstUpdate.SessionId
		request.RequestNumber = firstUpdate.RequestNumber
		request.IMSI = credit_control.RemoveIMSIPrefix(common.GetSid().GetId())
		request.Msisdn = common.GetMsisdn()
		request.UeIPV4 = common.GetUeIpv4()
		request.SpgwIPV4 = firstUpdate.SpgwIpv4
		request.Apn = common.GetApn()
		request.Imei = firstUpdate.Imei
		request.PlmnID = firstUpdate.PlmnId
		request.UserLocation = firstUpdate.UserLocation
		request.ChargingCharacteristics = firstUpdate.ChargingCharacteristics
		request.Type = credit_control.CRTUpdate
		request.RatType = GetRATType(common.GetRatType())
		request.TgppCtx = firstUpdate.GetTgppCtx()

		request.Credits = []*UsedCredits{}
		for _, updateN := range listUpdates {
			request.Credits = append(request.Credits, &UsedCredits{
				RatingGroup:       updateN.Usage.ChargingKey,
				ServiceIdentifier: fromServiceIdentifier(updateN.Usage.ServiceIdentifier),
				InputOctets:       updateN.Usage.BytesTx, // transmit == input
				OutputOctets:      updateN.Usage.BytesRx, // receive == output
				TotalOctets:       updateN.Usage.BytesTx + updateN.Usage.BytesRx,
				Type:              UsedCreditsType(updateN.Usage.Type),
				RequestedUnits:    updateN.Usage.GetRequestedUnits(),
			})
		}
		requests = append(requests, request)
	}
	return requests
}

func fromServiceIdentifier(si *protos.ServiceIdentifier) *uint32 {
	if si == nil {
		return nil
	}
	return &si.Value
}

func GetRATType(prt protos.RATType) string {
	switch prt {
	case protos.RATType_TGPP_WLAN:
		return RAT_TYPE_WLAN
	default: // including protos.RATType_TGPP_LTE
		return RAT_TYPE_EUTRAN
	}
}
