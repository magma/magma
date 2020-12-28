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

// Package service implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package servicers

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/services/s6a_proxy/metrics"
)

// sendAIR - sends AIR with given Session ID (sid)
func (s *s6aProxy) sendAIR(sid string, req *protos.AuthenticationInformationRequest, retryCount uint) error {
	c, err := s.connMan.GetConnection(s.smClient, s.config.ServerCfg)
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
	if req.NumRequestedEutranVectors > 0 || req.NumRequestedUtranGeranVectors == 0 {
		m.NewAVP(
			avp.RequestedEUTRANAuthenticationInfo,
			avp.Vbit|avp.Mbit,
			diameter.Vendor3GPP,
			genAuthInfoAvp(req.NumRequestedEutranVectors, irp, req.ResyncInfo))
	}
	if req.NumRequestedUtranGeranVectors > 0 {
		m.NewAVP(
			avp.RequestedUTRANGERANAuthenticationInfo,
			avp.Vbit|avp.Mbit,
			diameter.Vendor3GPP,
			genAuthInfoAvp(req.NumRequestedUtranGeranVectors, irp, req.UtranGeranResyncInfo))
	}
	glog.V(2).Infof("Sending S6a AIR message\n%s\n", m)
	err = c.SendRequest(m, retryCount)
	if err != nil {
		err = Error(codes.DataLoss, err)
	}
	return err
}

func genAuthInfoAvp(requestedVectorNum, irp uint32, resyncInfo []byte) *diam.GroupedAVP {
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(requestedVectorNum)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(irp)),
		},
	}
	if len(resyncInfo) > 0 {
		resyncInfo := diam.NewAVP(avp.ResynchronizationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP,
			datatype.OctetString(resyncInfo))
		authInfo.AddAVP(resyncInfo)
	}
	return authInfo
}

// S6a AIA
func handleAIA(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received S6a AIA message:\n%s\n", m)
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

	// check if PLMN list exists, and if IMSI belongs to any of that PLMN
	if !plmn_filter.CheckImsiOnPlmnIdListIfAny(req.UserName, s.config.PlmnIds) {
		return &protos.AuthenticationInformationAnswer{
			ErrorCode:     protos.ErrorCode_AUTHENTICATION_REJECTED,
			EutranVectors: nil,
		}, nil
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
					if len(aia.AI.EUtranVectors) > 0 {
						sort.Slice(aia.AI.EUtranVectors, func(i, j int) bool {
							return aia.AI.EUtranVectors[i].ItemNumber < aia.AI.EUtranVectors[j].ItemNumber
						})
						for _, ev := range aia.AI.EUtranVectors {
							res.EutranVectors = append(
								res.EutranVectors,
								&protos.AuthenticationInformationAnswer_EUTRANVector{
									Rand:  ev.RAND.Serialize(),
									Xres:  ev.XRES.Serialize(),
									Autn:  ev.AUTN.Serialize(),
									Kasme: ev.KASME.Serialize()})
						}
					}
					if len(aia.AI.UtranVectors) > 0 {
						sort.Slice(aia.AI.UtranVectors, func(i, j int) bool {
							return aia.AI.UtranVectors[i].ItemNumber < aia.AI.UtranVectors[j].ItemNumber
						})
						for _, uv := range aia.AI.UtranVectors {
							res.UtranVectors = append(
								res.UtranVectors,
								&protos.AuthenticationInformationAnswer_UTRANVector{
									Rand:               uv.RAND.Serialize(),
									Xres:               uv.XRES.Serialize(),
									Autn:               uv.AUTN.Serialize(),
									ConfidentialityKey: uv.CK.Serialize(),
									IntegrityKey:       uv.IK.Serialize()})
						}
					}
					if len(aia.AI.GeranVectors) > 0 {
						sort.Slice(aia.AI.GeranVectors, func(i, j int) bool {
							return aia.AI.GeranVectors[i].ItemNumber < aia.AI.GeranVectors[j].ItemNumber
						})
						for _, gv := range aia.AI.GeranVectors {
							res.GeranVectors = append(
								res.GeranVectors,
								&protos.AuthenticationInformationAnswer_GERANVector{
									Rand: gv.RAND.Serialize(),
									Sres: gv.SRES.Serialize(),
									Kc:   gv.Kc.Serialize()})
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
