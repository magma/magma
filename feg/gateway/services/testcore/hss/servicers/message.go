/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"time"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
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
