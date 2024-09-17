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

	"magma/feg/cloud/go/protos"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

// TerminateRegistration relays the RegistrationTerminationRequest sent from HSS->FeG->Access Gateway
func (srv *FegToGwRelayServer) TerminateRegistration(
	ctx context.Context, req *protos.RegistrationTerminationRequest) (*protos.RegistrationAnswer, error) {

	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.TerminateRegistrationUnverified(ctx, req)
}

// TerminateRegistrationUnverified relays the RegistrationTerminationRequest sent from HSS->FeG->Access Gateway
// without FeG Identity verification
func (srv *FegToGwRelayServer) TerminateRegistrationUnverified(
	ctx context.Context, req *protos.RegistrationTerminationRequest) (*protos.RegistrationAnswer, error) {

	hwId, err := getHwIDFromIMSI(ctx, req.UserName)
	if err != nil {
		errmsg := fmt.Errorf("unable to get HwID from IMSI %v. err: %v", req.GetUserName(), err)
		log.Print(errmsg)
		return &protos.RegistrationAnswer{}, errmsg
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwAAAService, hwId)
	if err != nil {
		errmsg := fmt.Errorf("unable to connect to GW %s to terminate service for IMSI: %s. err: %v",
			hwId, req.GetUserName(), err)
		log.Print(errmsg)
		return &protos.RegistrationAnswer{}, errmsg
	}
	client := protos.NewSwxGatewayServiceClient(conn)
	return client.TerminateRegistration(ctx, req)
}
