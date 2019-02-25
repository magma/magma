/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"errors"
	"fmt"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/servicers"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/dict"
)

// NewMAA outputs a multimedia authentication answer (MAA) to reply to a multimedia
// authentication request (MAR) message.
func NewMAA(srv *HomeSubscriberServer, msg *diam.Message) (*diam.Message, error) {
	err := ValidateMAR(msg)
	if err != nil {
		return msg.Answer(diam.MissingAVP), err
	}

	var mar servicers.MAR
	if err := msg.Unmarshal(&mar); err != nil {
		return msg.Answer(diam.UnableToComply), fmt.Errorf("MAR Unmarshal failed for message: %v failed: %v", msg, err)
	}

	// TODO(vikg): Actually handle the MAR here
	return ConstructSuccessAnswer(msg, mar.SessionID, srv.Config.Server), nil
}

// ValidateMAR returns an error if the message is missing any mandatory AVPs.
// Mandatory AVPs are specified in 3GPP TS 29.273 Table 8.1.2.1.1/1.
func ValidateMAR(msg *diam.Message) error {
	if msg == nil {
		return errors.New("Message is nil")
	}
	_, err := msg.FindAVP(avp.UserName, dict.UndefinedVendorID)
	if err != nil {
		return errors.New("Missing IMSI in message")
	}
	_, err = msg.FindAVP(avp.SIPNumberAuthItems, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing SIP-Number-Auth-Items in message")
	}
	_, err = msg.FindAVP(avp.SIPAuthDataItem, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing SIP-Auth-Data-Item in message")
	}
	_, err = msg.FindAVP(avp.SIPAuthenticationScheme, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing SIP-Authentication-Scheme in message")
	}
	_, err = msg.FindAVP(avp.RATType, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing RAT type in message")
	}
	return nil
}
