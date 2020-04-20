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

	fegprotos "magma/feg/cloud/go/protos"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/errors"
)

// CancelLocation relays the CancelLocationRequest to a corresponding
// dispatcher service instance, who will in turn relay the request to the
// corresponding gateway
func (srv *FegToGwRelayServer) CancelLocation(
	ctx context.Context,
	req *fegprotos.CancelLocationRequest,
) (*fegprotos.CancelLocationAnswer, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.CancelLocationUnverified(ctx, req)
}

// CancelLocationUnverified called directly in test server for unit test.
// Skip identity check
func (srv *FegToGwRelayServer) CancelLocationUnverified(
	ctx context.Context,
	req *fegprotos.CancelLocationRequest,
) (*fegprotos.CancelLocationAnswer, error) {
	hwId, err := getHwIDFromIMSI(ctx, req.UserName)
	if err != nil {
		fmt.Printf("unable to get HwID from IMSI %v. err: %v", req.UserName, err)
		if _, ok := err.(errors.ClientInitError); ok {
			return &fegprotos.CancelLocationAnswer{ErrorCode: fegprotos.ErrorCode_UNABLE_TO_DELIVER}, nil
		}
		return &fegprotos.CancelLocationAnswer{ErrorCode: fegprotos.ErrorCode_USER_UNKNOWN}, nil
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(
		gateway_registry.GwS6aService, hwId)
	if err != nil {
		fmt.Printf("unable to get connection to the gateway ID: %s", hwId)
		return &fegprotos.CancelLocationAnswer{ErrorCode: fegprotos.ErrorCode_UNABLE_TO_DELIVER}, nil
	}
	client := fegprotos.NewS6AGatewayServiceClient(conn)
	return client.CancelLocation(ctx, req)
}
