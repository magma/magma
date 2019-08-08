/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"fmt"

	"magma/feg/cloud/go/protos"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/golang/glog"
)

const (
	// MaxDiamRTRetries - number of retries for responding to RTR
	MaxDiamRTRetries = 1
)

type fegRelayClient struct{}

func (r *fegRelayClient) RelayFromFeg() (protos.ErrorCode, error) {
	return protos.ErrorCode_COMMAND_UNSUPORTED, fmt.Errorf("Relay for RTR is unimplemented")
}

func handleRTR(s *swxProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("handling RTR %v\n", m)
		var rtr RTR
		err := m.Unmarshal(&rtr)
		if err != nil {
			glog.Errorf("RTR Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		code, err := s.Relay.RelayFromFeg()
		if err != nil {
			glog.Error(err)
		}
		err = s.sendRTA(c, m, code, &rtr, MaxDiamRTRetries)
		if err != nil {
			glog.Errorf("Failed to send RTA: %s", err.Error())
		}
	}
}

func (s *swxProxy) sendRTA(c diam.Conn, m *diam.Message, code protos.ErrorCode, rtr *RTR, retries uint) error {
	ans := m.Answer(uint32(code))
	// SessionID is required to be the AVP in position 1
	ans.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(rtr.SessionID)))
	ans.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(rtr.AuthSessionState))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Host))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Realm))
	if s.originStateID != 0 {
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(s.originStateID))
	}
	_, err := ans.WriteToWithRetry(c, retries)
	return err
}
