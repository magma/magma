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
// It also handles CLR & RSR sends sync rpc request to gateway, then returns a CLA/RSA over diameter connection.
package servicers

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s6a_proxy"
)

const (
	// MaxDiamRsRetries - number of retries for forwarding RSR to a gateway
	MaxDiamRsRetries = 1
)

// S6a CLR
func handleRSR(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("handling RSR\n")
		var code protos.ErrorCode //result-code
		var rsr RSR
		err := m.Unmarshal(&rsr)
		if err != nil {
			glog.Errorf("RSR Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		var retries = MaxSyncRPCRetries
		for ; retries >= 0; retries-- {
			code, err = forwardRSRToGateway(&rsr)
			if err != nil {
				glog.Errorf("Failed to forward RSR to gateway. err: %v. Retries left: %v\n", err, retries)
			} else {
				break
			}
		}
		err = s.sendRSA(c, m, code, &rsr, MaxDiamRsRetries)
		if err != nil {
			glog.Errorf("Failed to send RSA: %s", err.Error())
		}
	}
}

func forwardRSRToGateway(rsr *RSR) (protos.ErrorCode, error) {
	if rsr == nil {
		return diam.MissingAVP, nil
	}
	in := new(protos.ResetRequest)

	res, err := s6a_proxy.GWS6AProxyReset(in)
	if err != nil {
		return protos.ErrorCode_UNABLE_TO_DELIVER, err
	}
	return res.ErrorCode, nil
}

func (s *s6aProxy) sendRSA(c diam.Conn, m *diam.Message, code protos.ErrorCode, rsr *RSR, retries uint) error {
	ans := m.Answer(uint32(code))
	// SessionID is required to be the AVP in position 1
	ans.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(rsr.SessionID)))
	ans.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(rsr.AuthSessionState))
	s.addDiamOriginAVPs(ans)

	_, err := ans.WriteToWithRetry(c, retries)
	return err

}
