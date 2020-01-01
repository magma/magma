// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package service

import (
	"log"
	"math/rand"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
	"github.com/fiorix/go-diameter/v4/examples/s6a_proxy/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sendULR - sends ULR with given Session ID (sid)
func (s *s6aProxy) sendULR(sid string, req *protos.UpdateLocationRequest) error {
	c := s.conn

	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return Errorf(codes.Internal, "peer metadata unavailable for ULR")
	}
	m := diam.NewRequest(diam.UpdateLocation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.cfg.Host))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.cfg.Realm))
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.UserName))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	m.NewAVP(avp.RATType, avp.Mbit, VENDOR_3GPP, datatype.Enumerated(ULR_RAT_TYPE))
	m.NewAVP(avp.ULRFlags, avp.Vbit|avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(ULR_FLAGS))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, VENDOR_3GPP, datatype.OctetString(req.VisitedPlmn))

	// randomize stream of first message and append it to SID to test SCTP multi-streaming support
	stream := uint(rand.Int31n(diam.MaxOutboundSCTPStreams - 2))
	s.airSendLocks[stream].Lock()
	defer s.airSendLocks[stream].Unlock()
	_, err := m.WriteToStream(c, stream)
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
		s.sessionsMu.Lock()
		ch, ok := s.sessions[ula.SessionID]
		if ok {
			delete(s.sessions, ula.SessionID)
			s.sessionsMu.Unlock()
			ch <- &ula
		} else {
			s.sessionsMu.Unlock()
			log.Printf("ULA SessionID %s not found. Message: %s, Remote: %s", ula.SessionID, m, c.RemoteAddr())
		}
	}
}

// UpdateLocation sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULAA & returns its RPC representation
func (s *s6aProxy) UpdateLocationImpl(req *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error,
) {
	res := &protos.UpdateLocationAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil UL Request")
	}
	sid := genSID()
	ch := make(chan interface{})
	s.updateSession(sid, ch)

	var (
		err     error
		retries int = MAX_DIAM_RETRIES
		c       diam.Conn
	)
	for ; retries >= 0; retries-- {
		c, err = s.acquireConnection()

		if err != nil {
			s.releaseConnection()
			s.cleanupSession(sid)
			log.Printf("Cannot connect to %s://%s; %v", s.cfg.Protocol, s.cfg.HssAddr, err)
			return res, Error(codes.Unavailable, err)
		}
		err = s.sendULR(sid, req)

		s.releaseConnection() // we can unlock reader after send

		if err != nil {
			log.Printf("Error sending ULR with SID %s: %v", sid, err)
			if status, ok := status.FromError(err); ok && status != nil && status.Code() == codes.DataLoss {
				s.cleanupConn(c)
				continue
			}
		}
		break
	}

	if err == nil {
		select {
		case resp, open := <-ch:
			if open {
				ula, ok := resp.(*ULA)
				if ok {
					err = TranslateBaseDiamResultCode(ula.ResultCode)
					res.ErrorCode = protos.ErrorCode(ula.ExperimentalResult.ExperimentalResultCode)
					res.DefaultContextId = ula.SubscriptionData.APNConfigurationProfile.ContextIdentifier
					res.TotalAmbr = &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
						MaxBandwidthUl: ula.SubscriptionData.AMBR.MaxRequestedBandwidthUL,
						MaxBandwidthDl: ula.SubscriptionData.AMBR.MaxRequestedBandwidthDL,
					}
					res.AllApnsIncluded =
						ula.SubscriptionData.APNConfigurationProfile.AllAPNConfigurationsIncludedIndicator == 0

					for _, apnCfg := range ula.SubscriptionData.APNConfigurationProfile.APNConfigs {
						res.Apn = append(
							res.Apn,
							&protos.UpdateLocationAnswer_APNConfiguration{
								ContextId:        apnCfg.ContextIdentifier,
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
				}
			} else {
				err = Errorf(codes.Aborted, "ULR for Session ID: %s is canceled", sid)
			}
		case <-time.After(time.Second * TIMEOUT_SECONDS):
			err = Errorf(codes.DeadlineExceeded, "ULR Timed Out for Session ID: %s", sid)
		}
	}
	s.cleanupSession(sid)
	return res, err
}
