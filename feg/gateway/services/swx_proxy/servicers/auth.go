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

// AuthenticateImpl sends MAR over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s *swxProxy) AuthenticateImpl(req *protos.AuthenticationRequest) (*protos.AuthenticationAnswer, error) {
	res := &protos.AuthenticationAnswer{}
	err := validateAuthRequest(req)
	if err != nil {
		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}
	res.UserName = req.GetUserName()

	sid := s.genSID()
	ch := make(chan interface{})
	s.requestTracker.RegisterRequest(sid, ch)
	// if request hasn't been removed by end of transaction, remove it
	defer s.requestTracker.DeregisterRequest(sid)

	marMsg := s.createMAR(sid, req)
	err = s.sendDiameterMsg(marMsg, MAX_DIAM_RETRIES)
	if err != nil {
		glog.Errorf("Error while sending MAR with SID %s: %s", sid, err)
		return res, err
	}
	select {
	case resp, open := <-ch:
		if !open {
			err = status.Errorf(codes.Aborted, "MAA for Session ID: %s is canceled", sid)
			glog.Error(err)
			return res, err
		}
		maa, ok := resp.(*MAA)
		if !ok {
			err = status.Errorf(codes.Internal, "Invalid Response Type: %T, MAA expected.", resp)
			glog.Error(err)
			return res, status.Errorf(codes.Internal, "Invalid Response Type: %T, MAA expected.", resp)
		}
		err = diameter.TranslateDiamResultCode(maa.ResultCode)
		// If there is no base diameter error, check that there is no experimental error either
		if err == nil {
			err = diameter.TranslateDiamResultCode(maa.ExperimentalResult.ExperimentalResultCode)
		}
		// According to spec 29.273, SIP-Auth-Data-Item(s) only present on SUCCESS
		if err != nil {
			return res, err
		}
		res.SipAuthVectors = getSIPAuthenticationVectors(maa.SIPAuthDataItems)

	case <-time.After(time.Second * TIMEOUT_SECONDS):
		err = status.Errorf(codes.DeadlineExceeded, "MAA Timed Out for Session ID: %s", sid)
		glog.Error(err)
	}
	return res, err
}

// createMAR creates a Multimedia Authentication Request diameter msg with provided SessionID (sid)
// to be sent to HSS
func (s *swxProxy) createMAR(sid string, req *protos.AuthenticationRequest) *diam.Message {
	msg := diameter.NewProxiableRequest(diam.MultimediaAuthentication, diam.TGPP_SWX_APP_ID, dict.Default)
	msg.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.clientCfg.Host))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.clientCfg.Realm))
	msg.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.GetUserName()))
	msg.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	msg.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, uint32(diameter.Vendor3GPP), datatype.Unsigned32(req.GetSipNumAuthVectors()))
	msg.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(RadioAccessTechnologyType_WLAN))
	authDataAvp := []*diam.AVP{
		diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, uint32(diameter.Vendor3GPP), datatype.UTF8String(req.GetAuthenticationScheme())),
	}
	if len(req.GetResyncInfo()) > 0 {
		authDataAvp = append(
			authDataAvp,
			diam.NewAVP(avp.SIPAuthorization, avp.Mbit|avp.Vbit, uint32(diameter.Vendor3GPP), datatype.OctetString(req.GetResyncInfo())),
		)
	}
	msg.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{AVP: authDataAvp})
	return msg
}

func handleMAA(s *swxProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var maa MAA
		err := m.Unmarshal(&maa)
		if err != nil {
			glog.Errorf("MAA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(maa.SessionID)
		if ch != nil {
			ch <- &maa
		} else {
			glog.Errorf("MAA SessionID %s not found. Message: %s, Remote: %s", maa.SessionID, m, c.RemoteAddr())
		}
	}
}

func getSIPAuthenticationVectors(items []SIPAuthDataItem) []*protos.AuthenticationAnswer_SIPAuthVector {
	var authVectors []*protos.AuthenticationAnswer_SIPAuthVector
	for _, item := range items {
		// If the auth scheme is unrecognized, don't include the vector
		authScheme, err := getAuthenticationScheme(item.AuthScheme)
		if err != nil {
			glog.Error(err)
			continue
		}
		authVectors = append(
			authVectors,
			&protos.AuthenticationAnswer_SIPAuthVector{
				AuthenticationScheme: protos.AuthenticationScheme(authScheme),
				RandAutn:             item.Authenticate.Serialize(),
				Xres:                 item.Authorization.Serialize(),
				ConfidentialityKey:   item.ConfidentialityKey.Serialize(),
				IntegrityKey:         item.IntegrityKey.Serialize()})
	}
	return authVectors
}

func validateAuthRequest(req *protos.AuthenticationRequest) error {
	if req == nil {
		return fmt.Errorf("Nil authentication request provided")
	}
	if len(req.GetUserName()) == 0 {
		return fmt.Errorf("Empty user-name provided in authentication request")
	}
	if req.SipNumAuthVectors == 0 {
		return fmt.Errorf("SIPNumAuthVectors in authentication request must be greater than 0")
	}
	// imsi cannot be greater than 15 digits according to 3GPP Spec 23.003
	if len(req.GetUserName()) > 15 {
		return fmt.Errorf("Provided username %s is greater than 15 digits", req.GetUserName())
	}
	return nil
}

func getAuthenticationScheme(maaScheme string) (protos.AuthenticationScheme, error) {
	switch maaScheme {
	case "EAP-AKA":
		return protos.AuthenticationScheme_EAP_AKA, nil
	case "EAP-AKA'":
		return protos.AuthenticationScheme_EAP_AKA_PRIME, nil
	default:
		return protos.AuthenticationScheme_EAP_AKA, fmt.Errorf("Unrecognized Authentication Scheme returned: %s", maaScheme)
	}
}
