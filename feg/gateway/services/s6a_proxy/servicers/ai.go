/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package service implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package servicers

import (
	"log"
	"strconv"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/s6a_proxy/metrics"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"google.golang.org/grpc/codes"
)

// sendAIR - sends AIR with given Session ID (sid)
func (s *s6aProxy) sendAIR(sid string, req *protos.AuthenticationInformationRequest, retryCount uint) error {
	c, err := s.connMan.GetConnection(s.smClient, s.serverCfg)
	if err != nil {
		return err
	}
	var irp uint32
	if req.ImmediateResponsePreferred {
		irp = 1
	}
	m := diameter.NewProxiableRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	s.addDiamOriginAVPs(m)

	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.UserName))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.OctetString(req.VisitedPlmn))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(req.NumRequestedEutranVectors)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(irp)),
		},
	}
	if len(req.ResyncInfo) > 0 {
		resyncInfo := diam.NewAVP(avp.ResynchronizationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP,
			datatype.OctetString(req.ResyncInfo))
		authInfo.AddAVP(resyncInfo)
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)

	err = c.SendRequest(m, retryCount)
	if err != nil {
		err = Error(codes.DataLoss, err)
	}
	return err
}

// S6a AIA
func handleAIA(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var aia AIA
		err := m.Unmarshal(&aia)
		if err != nil {
			log.Printf("AIA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(aia.SessionID)
		if ch != nil {
			ch <- &aia
		} else {
			log.Printf("AIA SessionID %s not found. Message: %s, Remote: %s", aia.SessionID, m, c.RemoteAddr())
		}
	}
}

// AuthenticationInformationImpl sends AIR over diameter connection,
// waits (blocks) for AIA & returns its RPC representation
func (s *s6aProxy) AuthenticationInformationImpl(
	req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error) {

	res := &protos.AuthenticationInformationAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil AI Request")
	}

	sid := s.genSID()
	ch := make(chan interface{})
	s.requestTracker.RegisterRequest(sid, ch)
	// if request hasn't been removed by end of transaction, remove it
	defer s.requestTracker.DeregisterRequest(sid)

	var (
		err     error
		retries uint = MAX_DIAM_RETRIES
	)

	err = s.sendAIR(sid, req, retries)

	if err != nil {
		metrics.AIRSendFailures.Inc()
		log.Printf("Error sending AIR with SID %s: %v", sid, err)
	}

	if err == nil {
		metrics.AIRRequests.Inc()
		select {
		case resp, open := <-ch:
			if open {
				aia, ok := resp.(*AIA)
				if ok {
					metrics.S6aResultCodes.WithLabelValues(strconv.FormatUint(uint64(aia.ResultCode), 10)).Inc()
					err = diameter.TranslateDiamResultCode(aia.ResultCode)
					res.ErrorCode = protos.ErrorCode(aia.ExperimentalResult.ExperimentalResultCode)
					for _, ai := range aia.AIs {
						for _, ev := range ai.EUtranVectors {
							res.EutranVectors = append(
								res.EutranVectors,
								&protos.AuthenticationInformationAnswer_EUTRANVector{
									Rand:  ev.RAND.Serialize(),
									Xres:  ev.XRES.Serialize(),
									Autn:  ev.AUTN.Serialize(),
									Kasme: ev.KASME.Serialize()})
						}
					}
					return res, err // the only successful "exit" is here
				} else {
					err = Errorf(codes.Internal, "Invalid Response Type: %T, AIA expected.", resp)
					metrics.S6aUnparseableMsg.Inc()
				}
			} else {
				err = Errorf(codes.Aborted, "AIR for Session ID: %s is canceled", sid)
			}
		case <-time.After(time.Second * TIMEOUT_SECONDS):
			err = Errorf(codes.DeadlineExceeded, "AIR Timed Out for Session ID: %s", sid)
			metrics.S6aTimeouts.Inc()
		}
	}

	return res, err
}
