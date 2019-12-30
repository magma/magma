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

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"google.golang.org/grpc/codes"
)

// sendULR - sends ULR with given Session ID (sid)
func (s *s6aProxy) sendULR(sid string, req *protos.UpdateLocationRequest, retryCount uint) error {
	c, err := s.connMan.GetConnection(s.smClient, s.serverCfg)
	if err != nil {
		return err
	}
	m := diameter.NewProxiableRequest(diam.UpdateLocation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	s.addDiamOriginAVPs(m)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.UserName))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	m.NewAVP(avp.RATType, avp.Mbit, uint32(diameter.Vendor3GPP), datatype.Enumerated(ULR_RAT_TYPE))
	m.NewAVP(avp.ULRFlags, avp.Vbit|avp.Mbit, uint32(diameter.Vendor3GPP), datatype.Unsigned32(ULR_FLAGS))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.OctetString(req.VisitedPlmn))

	err = c.SendRequest(m, retryCount)
	if err != nil {
		err = Error(codes.DataLoss, err)
	}
	return err
}

// S6a ULA
func handleULA(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var ula ULA
		err := m.Unmarshal(&ula)
		if err != nil {
			log.Printf("ULA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(ula.SessionID)
		if ch != nil {
			ch <- &ula
		} else {
			log.Printf("ULA SessionID %s not found. Message: %s, Remote: %s", ula.SessionID, m, c.RemoteAddr())
		}
	}
}

// UpdateLocationImpl sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULA & returns its RPC representation
func (s *s6aProxy) UpdateLocationImpl(req *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error,
) {
	res := &protos.UpdateLocationAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil UL Request")
	}
	sid := s.genSID()
	ch := make(chan interface{})
	s.requestTracker.RegisterRequest(sid, ch)
	defer s.requestTracker.DeregisterRequest(sid)

	var (
		err     error
		retries uint = MAX_DIAM_RETRIES
	)

	err = s.sendULR(sid, req, retries)

	if err != nil {
		metrics.ULRSendFailures.Inc()
		log.Printf("Error sending ULR with SID %s: %v", sid, err)
	}

	if err == nil {
		metrics.ULRRequests.Inc()
		select {
		case resp, open := <-ch:
			if open {
				ula, ok := resp.(*ULA)
				if ok {
					metrics.S6aResultCodes.WithLabelValues(strconv.FormatUint(uint64(ula.ResultCode), 10)).Inc()
					err = diameter.TranslateDiamResultCode(ula.ResultCode)
					res.ErrorCode = protos.ErrorCode(ula.ExperimentalResult.ExperimentalResultCode)
					res.Msisdn = ula.SubscriptionData.MSISDN.Serialize()
					res.DefaultContextId = ula.SubscriptionData.APNConfigurationProfile.ContextIdentifier
					res.TotalAmbr = &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
						MaxBandwidthUl: ula.SubscriptionData.AMBR.MaxRequestedBandwidthUL,
						MaxBandwidthDl: ula.SubscriptionData.AMBR.MaxRequestedBandwidthDL,
					}
					res.AllApnsIncluded =
						ula.SubscriptionData.APNConfigurationProfile.AllAPNConfigurationsIncludedIndicator == 0
					res.NetworkAccessMode = protos.UpdateLocationAnswer_NetworkAccessMode(ula.SubscriptionData.NetworkAccessMode)

					for _, apnCfg := range ula.SubscriptionData.APNConfigurationProfile.APNConfigs {
						res.Apn = append(
							res.Apn,
							&protos.UpdateLocationAnswer_APNConfiguration{
								ContextId:        apnCfg.ContextIdentifier,
								Pdn:              protos.UpdateLocationAnswer_APNConfiguration_PDNType(apnCfg.PDNType),
								ServiceSelection: apnCfg.ServiceSelection,
								QosProfile: &protos.UpdateLocationAnswer_APNConfiguration_QoSProfile{
									ClassId:                 apnCfg.EPSSubscribedQoSProfile.QoSClassIdentifier,
									PriorityLevel:           apnCfg.EPSSubscribedQoSProfile.AllocationRetentionPriority.PriorityLevel,
									PreemptionCapability:    apnCfg.EPSSubscribedQoSProfile.AllocationRetentionPriority.PreemptionCapability == 0,
									PreemptionVulnerability: apnCfg.EPSSubscribedQoSProfile.AllocationRetentionPriority.PreemptionVulnerability == 0,
								},
								Ambr: &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
									MaxBandwidthUl: apnCfg.AMBR.MaxRequestedBandwidthUL,
									MaxBandwidthDl: apnCfg.AMBR.MaxRequestedBandwidthDL,
								},
							})
					}
					return res, err
				} else {
					err = Errorf(codes.Internal, "Invalid Response Type: %T, ULA expected.", resp)
					metrics.S6aUnparseableMsg.Inc()
				}
			} else {
				err = Errorf(codes.Aborted, "ULR for Session ID: %s is canceled", sid)
			}
		case <-time.After(time.Second * TIMEOUT_SECONDS):
			err = Errorf(codes.DeadlineExceeded, "ULR Timed Out for Session ID: %s", sid)
			metrics.S6aTimeouts.Inc()
		}
	}
	return res, err
}
