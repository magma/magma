/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	s6a "magma/feg/gateway/services/s6a_proxy/servicers"
	swx "magma/feg/gateway/services/swx_proxy/servicers"
	"magma/lte/cloud/go/services/eps_authentication/servicers"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/golang/glog"
)

// ConstructFailureAnswer creates an answer for the message with an embedded
// Experimental-Result AVP. This answer informs the peer that the request has failed.
// See 3GPP TS 29.272 section 7.4.3 (permanent errors) and section 7.4.4 (transient errors).
func ConstructFailureAnswer(msg *diam.Message, sessionID datatype.UTF8String, serverCfg *mconfig.DiamServerConfig, resultCode uint32) *diam.Message {
	newMsg := diam.NewMessage(
		msg.Header.CommandCode,
		msg.Header.CommandFlags&^diam.RequestFlag, // Reset the Request bit.
		msg.Header.ApplicationID,
		msg.Header.HopByHopID,
		msg.Header.EndToEndID,
		msg.Dictionary(),
	)
	AddStandardAnswerAVPS(newMsg, sessionID, serverCfg, resultCode)
	return newMsg
}

// ConvertAuthErrorToFailureMessage creates a corresponding diameter failure message for an auth error.
func ConvertAuthErrorToFailureMessage(err error, msg *diam.Message, sessionID datatype.UTF8String, serverCfg *mconfig.DiamServerConfig) *diam.Message {
	switch err.(type) {
	case servicers.AuthRejectedError:
		return ConstructFailureAnswer(msg, sessionID, serverCfg, uint32(protos.ErrorCode_AUTHORIZATION_REJECTED))
	case servicers.AuthDataUnavailableError:
		return ConstructFailureAnswer(msg, sessionID, serverCfg, uint32(protos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE))
	default:
		return ConstructFailureAnswer(msg, sessionID, serverCfg, uint32(diam.UnableToComply))
	}
}

// ConstructSuccessAnswer returns a message response with a success result code
// and with the server config AVPs already added.
func ConstructSuccessAnswer(msg *diam.Message, sessionID datatype.UTF8String, serverCfg *mconfig.DiamServerConfig) *diam.Message {
	answer := msg.Answer(diam.Success)
	AddStandardAnswerAVPS(answer, sessionID, serverCfg, diam.Success)
	return answer
}

// AddStandardAnswerAVPS adds the SessionID, ExperimentalResult, OriginHost, OriginRealm, and OriginStateID AVPs to a message.
func AddStandardAnswerAVPS(answer *diam.Message, sessionID datatype.UTF8String, serverCfg *mconfig.DiamServerConfig, resultCode uint32) {
	// SessionID is required to be the AVP in position 1
	answer.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, sessionID))
	answer.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
			diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(resultCode)),
		},
	})

	answer.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(serverCfg.DestHost))
	answer.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(serverCfg.DestRealm))
	answer.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(time.Now().Unix()))
}

// replyFunc creates a response message (or an error) to reply to a request message.
type replyFunc func(*HomeSubscriberServer, *diam.Message) (*diam.Message, error)

// handleMessage processes incoming request messages and sends an answer message
// back which is constructed using a replyFunc.
func (srv *HomeSubscriberServer) handleMessage(reply replyFunc) diam.HandlerFunc {
	return func(conn diam.Conn, msg *diam.Message) {
		if msg == nil {
			glog.Error("Received nil message")
			return
		}

		glog.V(2).Infof("Message received in hss service: %s", msg.String())

		answer, err := reply(srv, msg)
		if err != nil {
			glog.Error(err)
		}

		_, err = answer.WriteTo(conn)
		if err != nil {
			glog.Errorf("Failed to send response: %s", err.Error())
		}
	}
}

// getRedirectMessage returns a response message which can be used to redirect
// the user to a different 3GPP AAA server.
func getRedirectMessage(msg *diam.Message, sessionID datatype.UTF8String, serverCfg *mconfig.DiamServerConfig, aaaServer datatype.DiameterIdentity) *diam.Message {
	answer := msg.Answer(diam.RedirectIndication)
	AddStandardAnswerAVPS(answer, sessionID, serverCfg, uint32(protos.SwxErrorCode_IDENTITY_ALREADY_REGISTERED))
	answer.NewAVP(avp.TGPPAAAServerName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, aaaServer)
	return answer
}

func isRATTypeAllowed(ratType uint32) bool {
	return ratType == swx.RadioAccessTechnologyType_WLAN || ratType == s6a.RadioAccessTechnologyType_EUTRAN
}
