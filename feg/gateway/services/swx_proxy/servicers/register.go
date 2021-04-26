/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package servicers implements Swx GRPC proxy service which sends MAR/SAR messages over
// diameter connection, waits (blocks) for diameter's MAA/SAAs and returns their RPC representation
package servicers

import (
	"fmt"
	"strconv"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/metrics"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterImpl sends SAR (code 301) over diameter
// waits (blocks) for SAA and returns its RPC representation
func (s *swxProxy) RegisterImpl(req *protos.RegistrationRequest, serverAssignmentType uint32) (*protos.RegistrationAnswer, error) {
	sid := req.GetSessionId()
	if len(sid) == 0 {
		sid = s.genSID(req.GetUserName())
	}
	res := &protos.RegistrationAnswer{SessionId: sid}
	err := validateRegistrationRequest(req)
	if err != nil {
		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}
	_, err = s.sendSAR(req.GetUserName(), serverAssignmentType, sid)
	return res, err
}

func (s *swxProxy) sendSAR(userName string, serverAssignmentType uint32, sid string) (*SAA, error) {
	return s.sendSARExt(userName, serverAssignmentType, s.config.ClientCfg.Host, s.config.ClientCfg.Realm, sid)
}

func (s *swxProxy) sendSARExt(
	userName string, serverAssignmentType uint32, originHost, originRealm, sid string) (*SAA, error) {
	if len(sid) == 0 {
		sid = s.genSID(userName)
	}
	ch := make(chan interface{})
	s.requestTracker.RegisterRequest(sid, ch)
	// if request hasn't been removed by end of transaction, remove it
	defer s.requestTracker.DeregisterRequest(sid)

	sarMsg := s.createSAR(sid, userName, serverAssignmentType, originHost, originRealm)

	sarStartTime := time.Now()
	err := s.sendDiameterMsg(sarMsg, MAX_DIAM_RETRIES)
	if err != nil {
		metrics.SARSendFailures.Inc()
		glog.Errorf("Error while sending SAR with SID %s: %s", sid, err)
		return nil, err
	}
	metrics.SARRequests.Inc()
	select {
	case resp, open := <-ch:
		metrics.SARLatency.Observe(time.Since(sarStartTime).Seconds())
		if !open {
			metrics.SwxInvalidSessions.Inc()
			err = status.Errorf(codes.Aborted, "SAA for Session ID: %s is canceled", sid)
			glog.Error(err)
			return nil, err
		}
		saa, ok := resp.(*SAA)
		if !ok {
			metrics.SwxUnparseableMsg.Inc()
			err = status.Errorf(codes.Internal, "Invalid Response Type: %T, SAA expected.", resp)
			glog.Error(err)
			return nil, err
		}
		err = diameter.TranslateDiamResultCode(saa.ResultCode)
		metrics.SwxResultCodes.WithLabelValues(strconv.FormatUint(uint64(saa.ResultCode), 10)).Inc()
		// If there is no base diameter error, check that there is no experimental error either
		if err == nil {
			err = diameter.TranslateDiamResultCode(saa.ExperimentalResult.ExperimentalResultCode)
			metrics.SwxExperimentalResultCodes.WithLabelValues(strconv.FormatUint(uint64(saa.ExperimentalResult.ExperimentalResultCode), 10)).Inc()
		}
		return saa, err

	case <-time.After(time.Second * TIMEOUT_SECONDS):
		metrics.SARLatency.Observe(time.Since(sarStartTime).Seconds())
		metrics.SwxTimeouts.Inc()
		err = status.Errorf(codes.DeadlineExceeded, "SAA Timed Out for Session ID: %s", sid)
		glog.Error(err)
		return nil, err
	}
}

// createSAR creates a Server Assignment Request with provided SessionID (sid),
// UserName, ServerAssignmentType (saType), originHost and originRealm to be sent over diameter to HSS
func (s *swxProxy) createSAR(sid, userName string, saType uint32, originHost, originRealm string) *diam.Message {
	msg := diameter.NewProxiableRequest(diam.ServerAssignment, diam.TGPP_SWX_APP_ID, dict.Default)
	msg.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	msg.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_SWX_APP_ID)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
		},
	})
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(originHost))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(originRealm))

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
			metrics.SwxUnparseableMsg.Inc()
			glog.Errorf("SAA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(saa.SessionID)
		if ch != nil {
			ch <- &saa
		} else {
			metrics.SwxInvalidSessions.Inc()
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
