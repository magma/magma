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

// package service implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
// It also handles CLR, sends sync rpc request to gateway, then returns a CLA over diameter connection.
package servicers

import (
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s6a_proxy"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"
)

const (
	MaxSyncRPCRetries = 3
	MaxDiamClRetries  = 3
)

// S6a CLR
func handleCLR(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received S6a CLR message:\n%s\n", m)
		var code uint32 //result-code
		var clr CLR
		err := m.Unmarshal(&clr)
		if err != nil {
			glog.Errorf("CLR Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		var retries = MaxSyncRPCRetries
		for ; retries >= 0; retries-- {
			code, err = forwardCLRToGateway(&clr)
			if err != nil {
				glog.Errorf("Failed to forward CLR to gateway. err: %v. Retries left: %v\n", err, retries)
			} else {
				break
			}
		}
		err = s.sendCLA(c, m, code, &clr, MaxDiamClRetries)
		if err != nil {
			glog.Errorf("Failed to send CLA: %s", err.Error())
		} else {
			glog.V(2).Infof("Successfully sent CLA\n")
		}
	}
}

func mapProtoToDiamResult(protoErr protos.ErrorCode) int {
	switch protoErr {
	case protos.ErrorCode_SUCCESS, protos.ErrorCode_USER_UNKNOWN, protos.ErrorCode_UNKNOWN_SESSION_ID:
		return diam.Success
	default:
		if protoErr >= diam.MultiRoundAuth && protoErr <= diam.NoCommonSecurity {
			return int(protoErr)
		}
		return diam.UnableToDeliver
	}
}

func forwardCLRToGateway(clr *CLR) (uint32, error) {
	cancelLocationType := protos.CancelLocationRequest_CancellationType(clr.CancellationType)
	in := &protos.CancelLocationRequest{UserName: clr.UserName, CancellationType: cancelLocationType}
	res, err := s6a_proxy.GWS6AProxyCancelLocation(in)
	if err != nil {
		if res != nil {
			return uint32(mapProtoToDiamResult(res.ErrorCode)), err
		}
		return diam.UnableToDeliver, err
	}
	return uint32(mapProtoToDiamResult(res.ErrorCode)), nil
}

func (s *s6aProxy) sendCLA(c diam.Conn, m *diam.Message, code uint32, clr *CLR, retries uint) error {
	ans := m.Answer(code)
	// SessionID is required to be the AVP in position 1
	ans.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(clr.SessionID)))
	ans.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(clr.AuthSessionState))
	s.addDiamOriginAVPs(m)
	glog.V(2).Infof("Sending S6a CLA message\n%s\n", m)
	_, err := ans.WriteToWithRetry(c, retries)
	return err

}
