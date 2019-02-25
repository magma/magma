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

// NewSAA outputs a server assignment answer (SAA) to reply to a server
// assignment request (SAR) message.
func NewSAA(srv *HomeSubscriberServer, msg *diam.Message) (*diam.Message, error) {
	err := ValidateSAR(msg)
	if err != nil {
		return msg.Answer(diam.MissingAVP), err
	}

	var sar servicers.SAR
	if err := msg.Unmarshal(&sar); err != nil {
		return msg.Answer(diam.UnableToComply), fmt.Errorf("SAR Unmarshal failed for message: %v failed: %v", msg, err)
	}

	// TODO(vikg): Actually handle the SAR here
	return ConstructSuccessAnswer(msg, sar.SessionID, srv.Config.Server), nil
}

// ValidateSAR returns an error if the message is missing any mandatory AVPs.
// Mandatory AVPs are specified in 3GPP TS 29.273 Table 8.1.2.2.2.1/1.
func ValidateSAR(msg *diam.Message) error {
	if msg == nil {
		return errors.New("Message is nil")
	}
	_, err := msg.FindAVP(avp.UserName, dict.UndefinedVendorID)
	if err != nil {
		return errors.New("Missing IMSI in message")
	}
	_, err = msg.FindAVP(avp.ServerAssignmentType, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing server assignment type in message")
	}
	return nil
}
