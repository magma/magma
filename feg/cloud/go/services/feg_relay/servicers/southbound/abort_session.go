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
	"log"

	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

// AbortSession relays the AbortSessionRequest sent from PCRF/OCS->FeG->Access Gateway
func (srv *FegToGwRelayServer) AbortSession(
	ctx context.Context, req *protos.AbortSessionRequest) (*protos.AbortSessionResult, error) {

	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.AbortSessionUnverified(ctx, req)
}

// AbortSessionUnverified relays the AbortSessionRequest sent from PCRF/OCS->FeG->Access Gateway
// without FeG Identity verification
func (srv *FegToGwRelayServer) AbortSessionUnverified(
	ctx context.Context, req *protos.AbortSessionRequest) (*protos.AbortSessionResult, error) {

	hwId, err := getHwIDFromIMSI(ctx, req.UserName)
	if err != nil {
		msg := fmt.Sprintf("unable to get HwID from IMSI %v. err: %v", req.GetUserName(), err)
		log.Print(msg)
		return &protos.AbortSessionResult{
			Code:         protos.AbortSessionResult_GATEWAY_NOT_FOUND,
			ErrorMessage: msg}, nil
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwAbortSessionService, hwId)
	if err != nil {
		msg := fmt.Sprintf("unable to connect to GW %s to abbort session for IMSI: %s. err: %v",
			hwId, req.GetUserName(), err)
		log.Print(msg)
		return &protos.AbortSessionResult{
			Code:         protos.AbortSessionResult_GATEWAY_NOT_FOUND,
			ErrorMessage: msg}, nil
	}
	client := protos.NewAbortSessionResponderClient(conn)
	return client.AbortSession(ctx, req)
}
