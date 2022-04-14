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

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay"
	"magma/feg/cloud/go/services/feg_relay/utils"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	orc8r_protos "magma/orc8r/lib/go/protos"
)

const FegS8Proxy gateway_registry.GwServiceType = "s8_proxy"

type S8RelayRouter struct {
	*gw_to_feg_relay.Router
}

// NewRelayRouter creates & returns a new RelayRouter
func NewS8RelayRouter(router *gw_to_feg_relay.Router) S8RelayRouter {
	if router == nil {
		router = &gw_to_feg_relay.Router{}
	}
	return S8RelayRouter{Router: router}
}

func (s S8RelayRouter) CreateSession(
	c context.Context, req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	err := validateCreateSessionRequest(req)
	if err != nil {
		return nil, err
	}
	// inject control plane TEID
	req.CAgwTeid, err = utils.GetUniqueSgwTeid(c, utils.ControlPlaneTeid)
	if err != nil {
		err = fmt.Errorf("feg_relay S8 Create Session couldn't get unite SgwCteid: %v; request: %s", err, req.String())
		glog.Error(err)
		return nil, err
	}
	// inject user plane TEID
	req.BearerContext.UserPlaneFteid.Teid, err = utils.GetUniqueSgwTeid(c, utils.UserPlaneTeid)
	if err != nil {
		err = fmt.Errorf("feg_relay S8 Create Session couldn't get unique SgwUteid: %v; request: %s", err, req.String())
		glog.Error(err)
		return nil, err
	}
	client, ctx, cancel, err := s.getS8Client(c, req.GetImsi())
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := client.CreateSession(ctx, req)
	if err != nil {
		err = fmt.Errorf("S8 Create Session failure: %v; request: %s", err, req.String())
		glog.Error(err)
	}
	return res, err
}

func (s S8RelayRouter) DeleteSession(c context.Context, req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	client, ctx, cancel, err := s.getS8Client(c, req.GetImsi())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.DeleteSession(ctx, req)
}

func (s S8RelayRouter) SendEcho(c context.Context, req *protos.EchoRequest) (*protos.EchoResponse, error) {
	client, ctx, cancel, err := s.getS8Client(c, req.GetImsi())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.SendEcho(ctx, req)
}

func (s S8RelayRouter) CreateBearerResponse(c context.Context, req *protos.CreateBearerResponsePgw) (*orc8r_protos.Void, error) {
	client, ctx, cancel, err := s.getS8Client(c, req.GetImsi())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.CreateBearerResponse(ctx, req)
}

func (s S8RelayRouter) DeleteBearerResponse(c context.Context, req *protos.DeleteBearerResponsePgw) (*orc8r_protos.Void, error) {
	client, ctx, cancel, err := s.getS8Client(c, req.GetImsi())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.DeleteBearerResponse(ctx, req)
}

func (s S8RelayRouter) getS8Client(c context.Context, imsi string) (protos.S8ProxyClient, context.Context, context.CancelFunc, error) {

	conn, ctx, cancel, err := s.GetFegServiceConnection(c, imsi, FegS8Proxy)
	if err != nil {
		glog.V(1).Infof("failed to get FeG S8 service connection for IMSI %s: %v", imsi, err)
		return nil, nil, nil, err
	}
	return protos.NewS8ProxyClient(conn), ctx, cancel, nil
}

func validateCreateSessionRequest(req *protos.CreateSessionRequestPgw) error {
	if req == nil || req.BearerContext == nil || req.BearerContext.UserPlaneFteid == nil {
		return fmt.Errorf("s8 Create Bearere request is missing some parameters: %+v", req)
	}
	return nil
}
