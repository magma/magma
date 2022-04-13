/*
Copyright 2021 The Magma Authors.

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

	"github.com/golang/glog"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay/utils"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	orc8r_protos "magma/orc8r/lib/go/protos"
)

const GTPCauseNotAvailable = 73 // gtp CauseNoResourcesAvailable

// CreateBearer relays the CreateBearerRequest from S8_proxy to a corresponding
// dispatcher service instance, who will in turn relay the request to the
// corresponding AGW gateway
func (srv *FegToGwRelayServer) CreateBearer(
	ctx context.Context,
	req *fegprotos.CreateBearerRequestPgw,
) (*orc8r_protos.Void, error) {
	if req == nil {
		err := fmt.Errorf("feg_relay unable to send CreateBearerPGW, request is nil: ")
		glog.Error(err)
		return nil, err
	}
	// inject user plane TEID
	var err error
	req.UAgwTeid, err = utils.GetUniqueSgwTeid(ctx, utils.UserPlaneTeid)
	if err != nil {
		err = fmt.Errorf("feg_relay S8 CreateBearer couldn't get unique SgwUteid: %v; request: %s", err, req.String())
		glog.Error(err)
		return nil, err
	}

	// get AGW id this cTEID serves
	cTeid := fmt.Sprint(req.CAgwTeid)
	client, ctx, err := getS8ProxyResponderClient(ctx, cTeid)
	if err != nil {
		err = fmt.Errorf("unable to get S8ProxyResponderClient: %s", err)
		glog.Error(err)
		return nil, err
	}
	return client.CreateBearer(ctx, req)
}

// DeleteBearerRequest relays the DeleteBearerRequest from S8_proxy to a corresponding
// dispatcher service instance, who will in turn relay the request to the
// corresponding AGW gateway
func (srv *FegToGwRelayServer) DeleteBearerRequest(
	ctx context.Context,
	req *fegprotos.DeleteBearerRequestPgw,
) (*orc8r_protos.Void, error) {
	if req == nil {
		err := fmt.Errorf("unable to send DeleteBearerPGW, request is nil: ")
		glog.Error(err)
		return nil, err
	}
	teid := fmt.Sprint(req.CAgwTeid)
	client, ctx, err := getS8ProxyResponderClient(ctx, teid)
	if err != nil {
		err = fmt.Errorf("unable to get S8ProxyResponderClient: %s", err)
		glog.Error(err)
		return nil, err
	}
	return client.DeleteBearerRequest(ctx, req)
}

func getS8ProxyResponderClient(ctx context.Context, teid string) (
	fegprotos.S8ProxyResponderClient, context.Context, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, ctx, err
	}
	hwId, err := getHwIDFromTeid(ctx, teid)
	if err != nil {
		return nil, ctx, fmt.Errorf("unable to get HwID from TEID %s. err: %v", teid, err)
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(
		gateway_registry.GwS8Service, hwId)
	if err != nil {
		err = fmt.Errorf("unable to get connection to the gateway ID: %s", hwId)
		return nil, ctx, err
	}
	return fegprotos.NewS8ProxyResponderClient(conn), ctx, nil
}
