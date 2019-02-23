/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package servicers implements Swx GRPC proxy service which sends MAR/SAR messages over
// diameter connection, waits (blocks) for diameter's MAA/SAAs and returns their RPC representation
package servicers

import (
	"fmt"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterImpl sends SAR (code 301) over diameter
// waits (blocks) for SAA and returns its RPC representation
func (s *swxProxy) RegisterImpl(req *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	res := &protos.RegistrationAnswer{}
	err := validateRegistrationRequest(req)
	if err != nil {
		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}

	sid := s.genSID()
	ch := make(chan interface{})
	s.requestTracker.RegisterRequest(sid, ch)
	// if request hasn't been removed by end of transaction, remove it
	defer s.requestTracker.DeregisterRequest(sid)

	sarMsg := s.createSAR(sid, req.GetUserName(), ServerAssignmentType_REGISTRATION)
	err = s.sendDiameterMsg(sarMsg, MAX_DIAM_RETRIES)
	if err != nil {
		glog.Errorf("Error while sending SAR with SID %s: %s", sid, err)
		return res, err
	}
	select {
	case resp, open := <-ch:
		if !open {
			err = status.Errorf(codes.Aborted, "SAA for Session ID: %s is cancelled", sid)
			glog.Error(err)
			return res, err
		}
		saa, ok := resp.(*SAA)
		if !ok {
			err = status.Errorf(codes.Internal, "Invalid Response Type: %T, SAA expected.", resp)
			glog.Error(err)
			return res, err
		}
		err = diameter.TranslateDiamResultCode(saa.ResultCode)
		// If there is no base diameter error, check that there is no experimental error either
		if err == nil {
			err = diameter.TranslateDiamResultCode(saa.ExperimentalResult.ExperimentalResultCode)
		}
	case <-time.After(time.Second * TIMEOUT_SECONDS):
		err = status.Errorf(codes.DeadlineExceeded, "SAA Timed Out for Session ID: %s", sid)
		glog.Error(err)
	}
	return res, err
}

// createSAR creates a Server Assignment Request with provided SessionID (sid),
// UserName, and ServerAssignmentType (saType) to be sent over diameter to HSS
func (s *swxProxy) createSAR(sid, userName string, saType uint32) *diam.Message {
	msg := diameter.NewProxiableRequest(diam.ServerAssignment, diam.TGPP_SWX_APP_ID, dict.Default)
	msg.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.clientCfg.Host))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.clientCfg.Realm))
	msg.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(userName))
	msg.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	msg.NewAVP(avp.ServerAssignmentType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(saType))
	return msg
}

func handleSAA(s *swxProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var saa SAA
		err := m.Unmarshal(&saa)
		if err != nil {
			glog.Errorf("SAA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(saa.SessionID)
		if ch != nil {
			ch <- &saa
		} else {
			glog.Errorf("SAA SessionID %s not found. Message: %s, Remote: %s", saa.SessionID, m, c.RemoteAddr())
		}
	}
}

func validateRegistrationRequest(req *protos.RegistrationRequest) error {
	if req == nil {
		return fmt.Errorf("Nil registration request provided")
	}
	if len(req.GetUserName()) == 0 {
		return fmt.Errorf("Empty user-name provided in registration request")
	}
	// imsi cannot be greater than 15 digits according to 3GPP Spec 23.003
	if len(req.GetUserName()) > 15 {
		return fmt.Errorf("Provided username %s is greater than 15 digits", req.GetUserName())
	}
	return nil
}
