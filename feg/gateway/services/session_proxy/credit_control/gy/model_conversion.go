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
	return credits
}

func (request *CreditControlRequest) FromCreditUsageUpdate(update *protos.CreditUsageUpdate) *CreditControlRequest {
	common := update.GetCommonContext()
	request.SessionID = update.SessionId
	request.RequestNumber = update.RequestNumber
	request.IMSI = credit_control.RemoveIMSIPrefix(common.GetSid().GetId())
	request.Msisdn = common.GetMsisdn()
	request.UeIPV4 = common.GetUeIpv4()
	request.SpgwIPV4 = update.SpgwIpv4
	request.Apn = common.GetApn()
	request.Imei = update.Imei
	request.PlmnID = update.PlmnId
	request.UserLocation = update.UserLocation
	request.ChargingCharacteristics = update.ChargingCharacteristics
	request.Type = credit_control.CRTUpdate

	request.Credits = []*UsedCredits{{
		RatingGroup:       update.Usage.ChargingKey,
		ServiceIdentifier: fromServiceIdentifier(update.Usage.ServiceIdentifier),
		InputOctets:       update.Usage.BytesTx, // transmit == input
		OutputOctets:      update.Usage.BytesRx, // receive == output
		TotalOctets:       update.Usage.BytesTx + update.Usage.BytesRx,
		Type:              UsedCreditsType(update.Usage.Type),
		RequestedUnits:    update.Usage.GetRequestedUnits(),
	}}
	request.RatType = GetRATType(common.GetRatType())
	request.TgppCtx = update.GetTgppCtx()
	return request
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
