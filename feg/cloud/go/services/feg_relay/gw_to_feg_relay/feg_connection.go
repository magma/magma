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

// gw_to_feg_relay is h2c & GRPC server serving requests from AGWs to FeG
package gw_to_feg_relay

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"

	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/registry"
)

// GetFegServiceConnection returns connection to FeG which serves given Gateway network (inCtx) & IMSI
// (in effect it's a connection to the local orc8r SynRPC forwarding service)
// GetFegServiceConnection first looks up in the existing connection cache and returns the cache connection for the
// address if it already exists. If no connection for the address exists - GetFegServiceConnection creates a new
// connection, caches it and returns it
func (rtr *Router) GetFegServiceConnection(
	inCtx context.Context,
	imsi string,
	service gateway_registry.GwServiceType,
) (conn *grpc.ClientConn, ctx context.Context, cancel context.CancelFunc, err error) {

	gwId, err := RetrieveGatewayIdentity(inCtx)
	if err != nil {
		return
	}
	fegHwId, err := FindServingFeGHwId(gwId.GetNetworkId(), imsi)
	if err != nil {
		return
	}
	if rtr == nil {
		// No Router provided, use old GetGatewayConnection API and return
		conn, ctx, err = gateway_registry.GetGatewayConnection(service, fegHwId)
		if err != nil {
			return
		}
		ctx, cancel = context.WithTimeout(ctx, registry.GrpcMaxTimeoutSec*time.Second)
		return conn, ctx, cancel, nil
	}

	// There is a Router with connections cache, use & update it
	var addr string
	addr, err = gateway_registry.GetServiceAddressForGateway(fegHwId)
	if err != nil {
		return
	}
	rtr.RLock() // take read connection cache lock
	var found bool
	serviceAddrKey := connKey{service: service, addr: addr}
	conn, found = rtr.connCache[serviceAddrKey]
	rtr.RUnlock() // release read connection cache lock

	if !found || conn == nil || conn.GetState() != connectivity.Ready {
		// there is either no existing connection for the address or it's in a bad state
		// create a new connection outside of the cache lock
		conn, err = connectToService(addr, service)
		if err != nil {
			return
		}
		// Lock & update the connections cache
		rtr.Lock() // take write connection cache lock
		if len(rtr.connCache) == 0 {
			rtr.connCache = map[connKey]*grpc.ClientConn{serviceAddrKey: conn}
		} else {
			connToCleanUp, ok := rtr.connCache[serviceAddrKey]
			if ok && connToCleanUp != nil {
				if connToCleanUp.GetState() == connectivity.Ready {
					// use the old connection & close the just created connection
					connToCleanUp, conn = conn, connToCleanUp
				} else {
					// use the just created connection & attempt close the old connection (ignore close failures)
					rtr.connCache[serviceAddrKey] = conn
				}
				go connToCleanUp.Close()
			} else {
				// no cached connection for the address, cache the just created connection for future use
				rtr.connCache[serviceAddrKey] = conn
			}
		}
		rtr.Unlock() // release write connection cache lock
	}
	md := metadata.New(map[string]string{gateway_registry.GatewayIdHeaderKey: fegHwId})
	ctx, cancel = context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
	return conn, metadata.NewOutgoingContext(ctx, md), cancel, nil
}

func connectToService(addr string, service gateway_registry.GwServiceType) (*grpc.ClientConn, error) {
	connCtx, connCancel := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
	defer connCancel()
	bckoff := backoff.DefaultConfig
	bckoff.MaxDelay = registry.GrpcMaxDelaySec * time.Second
	return registry.GetClientConnection(
		connCtx, addr,
		grpc.WithConnectParams(
			grpc.ConnectParams{Backoff: bckoff, MinConnectTimeout: registry.GrpcMaxTimeoutSec * time.Second}),
		grpc.WithBlock(),
		grpc.WithAuthority(string(service)))
}
