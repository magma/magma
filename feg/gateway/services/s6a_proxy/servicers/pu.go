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

// Package servicers implements S6a GRPC proxy service which sends AIR, ULR, PUR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs, PURs & returns their RPC representation
package servicers

import (
	"log"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"google.golang.org/grpc/codes"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
)

// sendPUR - sends PUR with given Session ID (sid)
func (s *s6aProxy) sendPUR(sid string, req *protos.PurgeUERequest, retryCount uint) error {
	c, err := s.connMan.GetConnection(s.smClient, s.config.ServerCfg)
	if err != nil {
		return err
	}
	m := diameter.NewProxiableRequest(diam.PurgeUE, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	s.addDiamOriginAVPs(m)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.UserName))

	err = c.SendRequest(m, retryCount)
	if err != nil {
		err = Error(codes.DataLoss, err)
	}
	return err
}

// S6a PUA
func handlePUA(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var pua PUA
		err := m.Unmarshal(&pua)
		if err != nil {
			log.Printf("PUA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(pua.SessionID)
		if ch != nil {
			ch <- &pua
		} else {
			log.Printf("PUA SessionID %s not found. Message: %s, Remote: %s", pua.SessionID, m, c.RemoteAddr())
		}
	}
}

// PurgeUEImpl sends PUR over diameter connection,
// waits (blocks) for PUA & returns its RPC representation
func (s *s6aProxy) PurgeUEImpl(req *protos.PurgeUERequest) (*protos.PurgeUEAnswer, error) {
	res := &protos.PurgeUEAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil PU Request")
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

	err = s.sendPUR(sid, req, retries)

	if err != nil {
		log.Printf("Error sending PUR with SID %s: %v", sid, err)
	}
	if err == nil {
		select {
		case resp, open := <-ch:
			if open {
				pua, ok := resp.(*PUA)
				if ok {
					err = diameter.TranslateDiamResultCode(pua.ResultCode)
					res.ErrorCode = protos.ErrorCode(pua.ResultCode)
					return res, err // the only successful "exit" is here
				}
				err = Errorf(codes.Internal, "Invalid Response Type: %T, PUA expected.", resp)
			} else {
				err = Errorf(codes.Aborted, "PUR for Session ID: %s is canceled", sid)
			}
		case <-time.After(time.Second * TIMEOUT_SECONDS):
			err = Errorf(codes.DeadlineExceeded, "PUR Timed Out for Session ID: %s", sid)
		}
	}
	return res, err
}
