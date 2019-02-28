/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package eapauth (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
//
//   The usage sequence should be as follows:
//
//     import (
//         "feg/gateway/services/eap_auth"
//         "feg/gateway/services/eap_auth/protos"
//     )
//     ...
//     providerTypes := eap.SupportedTypes()
//     // either select the desired type or iterate through all of them
//     ...
//     // on new connection:
//
//     var tp eap_auth.Provider
//     for _, tp := providerTypes {
//       idMsg, err := eap_auth.IdentityRequest(tp)
//       if err != nil { panic("should never happen") }
//       /* Send idMsg to peer;
//        * get peer response => idResp;
//        * if success: create protos.EapMessage from idResp -> eapMsg;
//        * else -> continue
//        */
//        break;
//      }
//      go func () {
//        r, err := eap_auth.HandleIdentityResponse(tp, eapMsg)
//        // check err, etc...
//        // forward r to the AP
//        for {
//            // get AP resp, convert to msg *protos.EapRequest, copy Ctx from r to msg
//            r, err := eap_auth.Handle(msg)
//            if err != nil { // handle error & exit }
//            p := r.GetPayload()
// 			  switch p.(type) {
// 				case *protos.EapResult_Response:
// 					// forward payload p to the AP
// 					// wait for AP resp, convert to msg & continue
// 					continue
// 				case *protos.EapResult_Success:
// 					// we are done, forward success to AP, save state and return
// 					return
// 				default:
// 					panic("shouldn't happen")
// 			  }
//       }()
//       ...
//go:generate protoc --go_out=plugins=grpc,paths=source_relative:. protos/eap_auth.proto
//
package eapauth

import (
	"errors"
	"fmt"

	"magma/feg/gateway/services/eapauth/protos"
)

const (
	// EAP Related Consts
	EapMethodIdentity = uint8(protos.EapType_Identity)
	EapCodeResponse   = uint8(protos.EapCode_Response)
)

// Handle handles passed EAP payload & returns corresponding EAP result
// NOTE: Identity Request is handled by APs & does not involve EAP Authenticator's support
func HandleIdentityResponse(providerType uint8, msg *protos.EapMessage) (*protos.EapResult, error) {
	if msg == nil {
		return nil, errors.New("Nil EAP Request")
	}
	err := verifyEapPayload(msg.Payload)
	if err != nil {
		return nil, err
	}
	if msg.Payload[EapMsgMethodType] != EapMethodIdentity {
		return nil, fmt.Errorf(
			"Invalid EAP Method Type for Identity Response: %d. Expecting EAP Identity (%d)",
			msg.Payload[EapMsgMethodType], EapMethodIdentity)
	}
	p := getProvider(providerType)
	if p == nil {
		return nil, unsupportedProviderError(providerType)
	}
	return p.Handle(&protos.EapRequest{Payload: msg.Payload})
}

// Handle handles passed EAP payload & returns corresponding EAP result
func Handle(msg *protos.EapRequest) (*protos.EapResult, error) {
	if msg == nil {
		return nil, errors.New("Nil EAP Message")
	}
	err := verifyEapPayload(msg.Payload)
	if err != nil {
		return nil, err
	}
	p := getProvider(msg.Payload[EapMsgMethodType])
	if p == nil {
		return nil, unsupportedProviderError(msg.Payload[EapMsgMethodType])
	}
	return p.Handle(msg)
}

// verifyEapPayload checks validity of EAP message & it's length
func verifyEapPayload(payload []byte) error {
	el := len(payload)
	if el < EapMsgData {
		return fmt.Errorf("EAP Message is too short: %d bytes", el)
	}
	mLen := uint16(payload[EapMsgLenHigh])<<8 + uint16(payload[EapMsgLenLow])
	if el < int(mLen) {
		return fmt.Errorf("Invalid EAP Message: bytes received %d are below specified length %d", el, mLen)
	}
	if payload[EapMsgCode] != EapCodeResponse {
		return fmt.Errorf(
			"Unsupported EAP Code: %d. Expecting EAP-Response (%d)",
			payload[EapMsgCode], EapCodeResponse)
	}
	return nil
}

func unsupportedProviderError(methodType uint8) error {
	return fmt.Errorf("Unsupported EAP Provider for Method Type: %d", methodType)
}
