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
	"fmt"

	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

// ReAuth initiates a credit reauth on the gateway
func (srv *FegToGwRelayServer) ChargingReAuth(
	ctx context.Context,
	req *protos.ChargingReAuthRequest,
) (*protos.ChargingReAuthAnswer, error) {
	if err := validateFegContext(ctx); err != nil {
		return &protos.ChargingReAuthAnswer{Result: protos.ReAuthResult_OTHER_FAILURE}, err
	}
	hwID, err := getHwIDFromIMSI(ctx, req.Sid)
	if err != nil {
		return &protos.ChargingReAuthAnswer{Result: protos.ReAuthResult_SESSION_NOT_FOUND},
			fmt.Errorf("unable to get HwID from IMSI %v. err: %v", req.Sid, err)
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwSessiondService, hwID)
	if err != nil {
		return &protos.ChargingReAuthAnswer{Result: protos.ReAuthResult_OTHER_FAILURE},
			fmt.Errorf("unable to get connection to the gateway ID: %s", hwID)
	}
	return protos.NewSessionProxyResponderClient(conn).ChargingReAuth(ctx, req)
}
