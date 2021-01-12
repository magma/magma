// Copyright 2020 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build with_builtin_radius

// package dae implements Radius Dynamic Authorization Extensions API (https://tools.ietf.org/html/rfc5176)
package dae

import (
	"context"

	"github.com/golang/glog"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa/protos"
	aaa_radius "magma/feg/gateway/services/aaa/radius"
)

type daeServerCfg struct {
	*mconfig.RadiusConfig
}

// NewDAEServicer returns servicer for built in DAE service
func NewDAEServicer(cfg *mconfig.RadiusConfig) DAE {
	return daeServerCfg{aaa_radius.ValidateConfigs(cfg)}
}

// Disconnect is DAE's Disconnect Messages equivalent
func (s daeServerCfg) Disconnect(aaaCtx *protos.Context) error {
	if s.RadiusConfig == nil {
		glog.V(1).Info("nil DAE configuration")
		return nil
	}
	if len(s.RadiusConfig.GetDAEAddr()) == 0 {
		glog.V(1).Info("empty DAE server address")
		return nil
	}
	rr := radius.Request{
		Packet: &radius.Packet{
			Code:       radius.CodeDisconnectRequest,
			Identifier: 0,
			Secret:     s.RadiusConfig.GetSecret(),
		},
	}
	rr.Add(rfc2866.AcctSessionID_Type, radius.Attribute(aaaCtx.GetSessionId()))
	rr.Add(rfc2865.CallingStationID_Type, radius.Attribute(aaaCtx.GetMacAddr()))
	resp, err := radius.Exchange(context.Background(), rr.Packet, s.RadiusConfig.GetDAEAddr())
	if err != nil {
		glog.Errorf("failed radius DAE request to %s: %v", s.RadiusConfig.GetDAEAddr(), err)
	} else {
		glog.V(2).Infof("DAE Disconnect Response Code %s", resp.Code)
	}
	return err
}
