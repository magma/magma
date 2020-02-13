/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/protos"
)

// MMInformationReq relays the MMInformationRequest sent from VLR->FeG->Access Gateway
func (srv *FegToGwRelayServer) MMInformationReq(
	ctx context.Context,
	req *fegprotos.MMInformationRequest,
) (*protos.Void, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.MMInformationRequestUnverified(ctx, req)
}

func (srv *FegToGwRelayServer) MMInformationRequestUnverified(
	ctx context.Context,
	req *fegprotos.MMInformationRequest,
) (*protos.Void, error) {
	conn, ctx, err := getGWSGSServiceConnCtx(ctx, req.Imsi)
	if err != nil {
		return &protos.Void{}, err
	}
	client := fegprotos.NewCSFBGatewayServiceClient(conn)
	return client.MMInformationReq(ctx, req)
}
