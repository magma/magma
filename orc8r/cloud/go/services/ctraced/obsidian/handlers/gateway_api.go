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

package handlers

import (
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type GwCtracedClient interface {
	StartCallTrace(networkId string, gatewayId string, req *protos.StartTraceRequest) (*protos.StartTraceResponse, error)

	EndCallTrace(networkId string, gatewayId string, req *protos.EndTraceRequest) (*protos.EndTraceResponse, error)
}

type gwCtracedClientImpl struct{}

func NewGwCtracedClient() GwCtracedClient {
	return gwCtracedClientImpl{}
}

// getGWCtracedClient gets a GRPC client to the ctraced service running on the gateway specified by (network ID, gateway ID).
// If gateway not found by configurator, returns ErrNotFound from magma/orc8r/lib/go/errors.
func getGWCtracedClient(networkID string, gatewayID string) (protos.CallTraceServiceClient, context.Context, error) {
	hwID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("gateway not found, network-id: %s, gateway-id: %s", networkID, gatewayID))
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwCtraced, hwID)
	if err != nil {
		errMsg := fmt.Sprintf("gateway ctraced client initialization error: %s", err)
		glog.Errorf(errMsg, err)
		return nil, nil, errors.New(errMsg)
	}
	return protos.NewCallTraceServiceClient(conn), ctx, nil
}

// StartCallTrace starts a call trace on the specified gateway
// If gateway not registered, returns ErrNotFound from magma/orc8r/lib/go/errors.
func (c gwCtracedClientImpl) StartCallTrace(networkId string, gatewayId string, req *protos.StartTraceRequest) (*protos.StartTraceResponse, error) {
	client, ctx, err := getGWCtracedClient(networkId, gatewayId)
	if err != nil {
		return nil, err
	}
	resp, err := client.StartCallTrace(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start call trace, gRPC client received error")
	}
	return resp, err
}

// EndCallTrace ends a call trace on the specified gateway
// If gateway not registered, returns ErrNotFound from magma/orc8r/lib/go/errors.
func (c gwCtracedClientImpl) EndCallTrace(networkId string, gatewayId string, req *protos.EndTraceRequest) (*protos.EndTraceResponse, error) {
	client, ctx, err := getGWCtracedClient(networkId, gatewayId)
	if err != nil {
		return nil, err
	}
	resp, err := client.EndCallTrace(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to end call trace, gRPC client received error")
	}
	return resp, err
}
