/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
Checkind service provides the gRPC interface for the Gateways, REST server to
collect/report and retreive Gateway runtime statuses.

The service relies on it's own DBs & tables for storage Gateway statuses and
magmad API for mapping Gateway IDs & networks.

It's not intended for storage of Gateway configs or any long term persistent
data/configuration .
*/
package servicers

import (
	"fmt"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind/store"
	"magma/orc8r/cloud/go/services/magmad"

	"golang.org/x/net/context"
)

// A little Go "polymorphism" magic for testing
type testCheckindServer struct {
	checkindServer
}

// Checkin receiver for testCheckindServer injects GW Identity into CTX if it's
// missing for testing without heavy mock of Certifier & certificate addition
func (srv *testCheckindServer) Checkin(
	ctx context.Context,
	req *protos.CheckinRequest) (*protos.CheckinResponse, error) {

	// See if there is an Identity in the CTX and add one if it was missing
	// but, leave it alone if it's present
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		networkId, err := magmad.FindGatewayNetworkId(req.GetGatewayId())
		if err != nil {
			return nil, fmt.Errorf("ID Lookup Error for Gateway '%s': %s",
				req.GetGatewayId(), err)
		}
		logicalId, err := magmad.FindGatewayId(networkId, req.GetGatewayId())
		if err != nil {
			return nil, err
		}
		ctx = protos.NewGatewayIdentity(
			req.GetGatewayId(), networkId, logicalId).NewContextWithIdentity(ctx)
	}
	return srv.checkindServer.Checkin(ctx, req)
}

func NewTestCheckindServer(store *store.CheckinStore) (*testCheckindServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Cannot initialize Test Checkin Server with Nil store")
	}
	return &testCheckindServer{checkindServer{store}}, nil
}
