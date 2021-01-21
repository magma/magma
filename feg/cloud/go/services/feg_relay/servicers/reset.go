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
	"context"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay/utils"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

// Reset (Code 322) over diameter connection,
// Not implemented
func (srv *FegToGwRelayServer) Reset(ctx context.Context, in *protos.ResetRequest) (*protos.ResetAnswer, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.ResetUnverified(ctx, in)
}

// ResetUnverified called directly in test server for unit test.
// Skip identity check
func (srv *FegToGwRelayServer) ResetUnverified(
	ctx context.Context,
	req *protos.ResetRequest,
) (*protos.ResetAnswer, error) {

	res := &protos.ResetAnswer{ErrorCode: protos.ErrorCode_SUCCESS}
	if req == nil {
		return res, nil
	}
	if len(req.UserId) == 0 { // empty IDs list - reset all IMSIs
		hwIds, err := utils.GetAllGatewayIDs(ctx)
		// Reset ALL is a "best effort" request, we reset what we can & just log error for what we cannot
		if err != nil {
			glog.Error(err)
		}
		go func() {
			for _, hwId := range hwIds {
				sendReset(hwId, req)
			}
		}()
	} else { // group IDs by their respective Gateways
		uidMap := map[string]struct{}{}
		hwIdMap := map[string][]string{}
		for _, uid := range req.UserId {
			if _, ok := uidMap[uid]; !ok {
				uidMap[uid] = struct{}{}
				hwId, err := getHwIDFromIMSI(ctx, uid)
				if err != nil {
					glog.Errorf("Reset: unable to get Hw ID for UID %s: %v.", uid, err)
					continue
				}
				if uidList, ok := hwIdMap[hwId]; ok {
					hwIdMap[hwId] = append(uidList, uid)
				} else {
					hwIdMap[hwId] = []string{uid}
				}
			}
		}
		go func() {
			for hwId, uidList := range hwIdMap {
				req.UserId = uidList
				sendReset(hwId, req)
			}
		}()
	}
	// Always ACK success, even if some GWs were unreachable. Our Reset support is "best effort"
	return &protos.ResetAnswer{ErrorCode: protos.ErrorCode_SUCCESS}, nil
}

// sendReset - sends reset request to a GW with given hwId, logs errors if any
func sendReset(hwId string, req *protos.ResetRequest) {
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwS6aService, hwId)
	if err != nil {
		glog.Errorf("Reset: unable to get connection to the gateway Hw ID: %s.", hwId)
		return
	}
	client := protos.NewS6AGatewayServiceClient(conn)
	ans, err := client.Reset(ctx, req)
	if err != nil {
		glog.Errorf("Reset error %v for gateway Hw ID: %s", err, hwId)
	} else if ans.ErrorCode > protos.ErrorCode_LIMITED_SUCCESS {
		diamErrName := protos.ErrorCode_name[int32(ans.ErrorCode)]
		glog.Errorf("Reset Diameter error %d (%s) for gateway Hw ID: %s.", ans.ErrorCode, diamErrName, hwId)
	}
}
