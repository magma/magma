/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package diameter

import (
	"fmt"

	"magma/feg/cloud/go/protos"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TranslateDiamResultCode maps Diameter Result Codes (both Base and Experimental)
// to GRPC Status Error. Diam success codes result in a nil error returned
func TranslateDiamResultCode(diamResult uint32) error {
	if diamResult == uint32(protos.ErrorCode_UNDEFINED) { // diamResult was not set (default will be 0)
		return nil
	}
	// diam result code is 2xxx
	if diamResult >= uint32(protos.ErrorCode_SUCCESS) && diamResult < uint32(protos.ErrorCode_COMMAND_UNSUPORTED) {
		return nil
	}
	errName, ok := diamCodeToNameMap[diamResult]
	if !ok {
		errName = "BASE_DIAMETER"
	}

	msg := fmt.Sprintf("Diameter Error: %d (%s)", diamResult, errName)
	glog.Errorf("RPC Result: %s", msg)
	return status.Errorf(codes.Code(diamResult), msg)
}
