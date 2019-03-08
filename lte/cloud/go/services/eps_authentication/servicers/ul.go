/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos"
)

func (srv *EPSAuthServer) UpdateLocation(ctx context.Context, ulr *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error) {
	return nil, status.Errorf(codes.Unimplemented, "update location not implemented")
}
