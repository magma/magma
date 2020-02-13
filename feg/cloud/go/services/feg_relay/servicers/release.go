/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	fegprotos "magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
)

// ReleaseReq relays the ReleaseRequest sent from VLR->FeG->Access Gateway
func (srv *FegToGwRelayServer) ReleaseReq(
	ctx context.Context,
	req *fegprotos.ReleaseRequest,
) (*protos.Void, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.ReleaseRequestUnverified(ctx, req)
}

func (srv *FegToGwRelayServer) ReleaseRequestUnverified(
	ctx context.Context,
	req *fegprotos.ReleaseRequest,
) (*protos.Void, error) {
	conn, ctx, err := getGWSGSServiceConnCtx(ctx, req.Imsi)
	if err != nil {
		return &protos.Void{}, err
	}
	client := fegprotos.NewCSFBGatewayServiceClient(conn)
	return client.ReleaseReq(ctx, req)
}
