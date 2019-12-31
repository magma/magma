// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package service

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
	"github.com/fiorix/go-diameter/v4/examples/s6a_proxy/protos"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sendAIR - sends AIR with given Session ID (sid)
func (s *s6aProxy) sendAIR(sid string, req *protos.AuthenticationInformationRequest) error {
	glog.V(4).Infof("Got into sendAIR request from gateway: %v\n", req)
	c := s.conn
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return Errorf(codes.Internal, "peer metadata unavailable for AIR")
	}
	var irp uint32
	if req.ImmediateResponsePreferred {
		irp = 1
	}
	// randomize stream of first message and append it to SID to test SCTP multi-streaming support
	stream := uint(rand.Int31n(diam.MaxOutboundSCTPStreams - 2))
	sid = fmt.Sprintf("%s;stream:%d", sid, stream)
	for i := stream; i > 0; i-- {
		sid += " " // variable len
	}
	m := diam.NewRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.cfg.Host))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.cfg.Realm))
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.UserName))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, VENDOR_3GPP, datatype.OctetString(req.VisitedPlmn))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				VENDOR_3GPP,
				datatype.Unsigned32(req.NumRequestedEutranVectors)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(irp)),
		},
	}
	if len(req.ResyncInfo) > 0 {
		resyncInfo := diam.NewAVP(avp.ResynchronizationInfo, avp.Vbit|avp.Mbit, VENDOR_3GPP,
			datatype.OctetString(req.ResyncInfo))
		authInfo.AddAVP(resyncInfo)
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, VENDOR_3GPP, authInfo)

	writer, isMulti := c.Connection().(diam.MultistreamWriter)
	if !isMulti {
		panic("Must be *diam.SCTPConn...")
	}
	// randomize stream of first message to test SCTP multistreaming support
	// and send the AIR in two pieces to simulate message fragmentation
	mBytes, err := m.Serialize()
	if err != err {
		return Error(codes.DataLoss, err)
	}
	l2 := len(mBytes) / 2
	// Don't allow simultaneous sends on the same stream, take stream scoped lock
	s.airSendLocks[stream].Lock()
	defer s.airSendLocks[stream].Unlock()

	_, err = writer.WriteStream(mBytes[:l2], stream)
	if err != nil {
		return Error(codes.DataLoss, err)
	}
	time.Sleep(time.Millisecond * 3)
	_, err = writer.WriteStream(mBytes[l2:], stream)
	if err != nil {
		return Error(codes.DataLoss, err)
	}
	return nil
}

// S6a AIA
func handleAIA(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(4).Infof("Got into handleAIA, diam msg: %v\n", m)
		var aia AIA
		err := m.Unmarshal(&aia)
		if err != nil {
			log.Printf("AIA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		s.sessionsMu.Lock()

		var msgStream uint
		idx := strings.LastIndex(aia.SessionID, ";stream:")
		sid := aia.SessionID[0:idx]
		if _, err = fmt.Sscanf(aia.SessionID[idx:], ";stream:%d", &msgStream); err != nil {
			panic(err)
		}
		ch, ok := s.sessions[sid]
		if ok {
			msc := c.Connection().(diam.MultistreamConn)
			stream := msc.CurrentStream()
			// Current diam.Server implementation reads & handles messages from a single thread/routine,
			// so - the receiver stream should not change until the message is handled. The test server should also
			// send responses on the the requests' streams
			if stream != msgStream {
				panic(
					fmt.Sprintf("STREAM %d From SessionID %q != conn stream %d; Conn: %+v\n",
						msgStream, aia.SessionID, stream, msc))
			}
			delete(s.sessions, sid)
			s.sessionsMu.Unlock()
			ch <- &aia
		} else {
			s.sessionsMu.Unlock()
			log.Printf("AIA SessionID %q (%q) not found. Message: %q, Remote: %q", aia.SessionID, sid, m, c.RemoteAddr())
		}
	}
}

// AuthenticationInformation sends AIR over diameter connection,
// waits (blocks) for AIA & returns its RPC representation
func (s *s6aProxy) AuthenticationInformationImpl(
	req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error) {
	glog.V(4).Infof("Got AI request from gateway\n")
	res := &protos.AuthenticationInformationAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil AI Request")
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
		err = s.sendAIR(sid, req)
		s.releaseConnection() // we can unlock reader after send
		if err != nil {
			log.Printf("Error sending AIR with SID %s: %v", sid, err)
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
				aia, ok := resp.(*AIA)
				if ok {
					err = TranslateBaseDiamResultCode(aia.ResultCode)
					res.ErrorCode = protos.ErrorCode(aia.ExperimentalResult.ExperimentalResultCode)
					for _, ai := range aia.AIs {
						res.EutranVectors = append(
							res.EutranVectors,
							&protos.AuthenticationInformationAnswer_EUTRANVector{
								Rand:  ai.EUtranVector.RAND.Serialize(),
								Xres:  ai.EUtranVector.XRES.Serialize(),
								Autn:  ai.EUtranVector.AUTN.Serialize(),
								Kasme: ai.EUtranVector.KASME.Serialize()})
					}
					return res, err // the only successful "exit" is here
				} else {
					err = Errorf(codes.Internal, "Invalid Response Type: %T, AIA expected.", resp)
				}
			} else {
				err = Errorf(codes.Aborted, "AIR for Session ID: %s is canceled", sid)
			}
		case <-time.After(time.Second * TIMEOUT_SECONDS):
			err = Errorf(codes.DeadlineExceeded, "AIR Timed Out for Session ID: %s", sid)
		}
	}
	s.cleanupSession(sid)
	return res, err
}
