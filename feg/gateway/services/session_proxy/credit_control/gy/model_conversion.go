/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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
	request.SessionID = update.SessionId
	request.RequestNumber = update.RequestNumber
	request.IMSI = credit_control.RemoveIMSIPrefix(update.Sid)
	request.Msisdn = update.Msisdn
	request.UeIPV4 = update.UeIpv4
	request.SpgwIPV4 = update.SpgwIpv4
	request.Apn = update.Apn
	request.Imei = update.Imei
	request.PlmnID = update.PlmnId
	request.UserLocation = update.UserLocation
	request.Type = credit_control.CRTUpdate
	request.Credits = []*UsedCredits{(&UsedCredits{}).FromCreditUsage(update.Usage)}
	request.RatType = GetRATType(update.GetRatType())
	request.TgppCtx = update.GetTgppCtx()
	return request
}

func GetRATType(prt protos.RATType) string {
	switch prt {
	case protos.RATType_TGPP_WLAN:
		return RAT_TYPE_WLAN
	default: // including protos.RATType_TGPP_LTE
		return RAT_TYPE_EUTRAN
	}
}
