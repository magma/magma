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

package servicers

import (
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/gateway/services/aaa/session_manager"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/lib/go/protos"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	IMSI_PREFIX   = "IMSI"
	DEFAULT_UE_IP = "192.168.88.88"
)

// UESimServerHssLess tracks all the UEs being simulated.
type UESimServerHssLess struct {
	store blobstore.BlobStorageFactory
	cfg   *UESimConfig
}

// NewUESimServerHssLess initializes a UESimServer with an empty store map.
// Output: a new UESimServerHssLess
func NewUESimServerHssLess(factory blobstore.BlobStorageFactory) (*UESimServerHssLess, error) {
	config, err := GetUESimConfig()
	if err != nil {
		return nil, err
	}
	return &UESimServerHssLess{
		store: factory,
		cfg:   config,
	}, nil

}

func (srv *UESimServerHssLess) AddUE(ctx context.Context, ue *cwfprotos.UEConfig) (*protos.Void, error) {
	ret := &protos.Void{}

	err := validateUEData(ue)
	if err != nil {
		err = ConvertStorageErrorToGrpcStatus(err)
		return ret, err
	}

	err = validateUEDataForHssLess(ue)
	if err != nil {
		err = ConvertStorageErrorToGrpcStatus(err)
		return ret, err
	}
	addUeToStore(srv.store, ue)

	return ret, nil
}
func (srv *UESimServerHssLess) Disconnect(ctx context.Context, id *cwfprotos.DisconnectRequest) (*cwfprotos.DisconnectResponse, error) {

	ue, err := getUE(srv.store, id.GetImsi())
	if err != nil {
		return &cwfprotos.DisconnectResponse{}, err
	}

	req := &lte_protos.LocalEndSessionRequest{
		Sid: makeSubscriberId(ue.GetImsi()),
		Apn: ue.GetHsslessCfg().GetApn(),
	}
	_, err = session_manager.EndSession(req)
	if err != nil {
		return &cwfprotos.DisconnectResponse{}, err
	}
	return &cwfprotos.DisconnectResponse{}, nil

}
func (srv *UESimServerHssLess) Authenticate(ctx context.Context, id *cwfprotos.AuthenticateRequest) (*cwfprotos.AuthenticateResponse, error) {

	ue, err := getUE(srv.store, id.GetImsi())
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	sid := makeSubscriberId(ue.GetImsi())
	req := &lte_protos.LocalCreateSessionRequest{
		CommonContext: &lte_protos.CommonSessionContext{
			Sid:     sid,
			UeIpv4:  DEFAULT_UE_IP,
			Apn:     ue.GetHsslessCfg().GetApn(),
			Msisdn:  ([]byte)(ue.GetHsslessCfg().GetMsisdn()),
			RatType: (lte_protos.RATType)(ue.GetHsslessCfg().GetRat()),
		},
		RatSpecificContext: &lte_protos.RatSpecificContext{
			Context: &lte_protos.RatSpecificContext_WlanContext{
				WlanContext: &lte_protos.WLANSessionContext{
					MacAddrBinary:   ([]byte)(srv.cfg.brMac),
					MacAddr:         srv.cfg.brMac,
					RadiusSessionId: "sessiond1",
				},
			},
		},
	}
	resp, err := session_manager.CreateSession(req)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{SessionId: ""}, err
	}
	activateReq := &lte_protos.UpdateTunnelIdsRequest{
		Sid:      sid,
		BearerId: 0,
		EnbTeid:  0,
		AgwTeid:  0,
	}

	// activate is needed to install all hte flows
	_, err = session_manager.UpdateTunnelIds(activateReq)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{SessionId: ""}, err
	}

	return &cwfprotos.AuthenticateResponse{SessionId: resp.GetSessionId()}, nil
}
func (srv *UESimServerHssLess) GenTraffic(ctx context.Context, req *cwfprotos.GenTrafficRequest) (*cwfprotos.GenTrafficResponse, error) {
	glog.V(2).Infof("Reached GenTraffic...")
	return nil, nil
}

func makeSubscriberId(imsi string) *lte_protos.SubscriberID {
	if !strings.HasPrefix(imsi, IMSI_PREFIX) {
		imsi = IMSI_PREFIX + imsi
	}
	return &lte_protos.SubscriberID{Id: imsi, Type: lte_protos.SubscriberID_IMSI}
}
