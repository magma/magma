/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gy

import (
	"magma/lte/cloud/go/protos"
)

func (redirectServer *RedirectServer) ToProto() *protos.RedirectServer {
	if redirectServer == nil {
		return &protos.RedirectServer{}
	}
	return &protos.RedirectServer{
		RedirectAddressType:   protos.RedirectServer_RedirectAddressType(redirectServer.RedirectAddressType),
		RedirectServerAddress: redirectServer.RedirectServerAddress,
	}
}
