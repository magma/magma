/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Errorf(code codes.Code, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	glog.Errorf("GRPC [%s] %s", code, msg)
	return status.Errorf(code, msg)
}

func NewGrpcValidationError(err error) error {
	glog.Error(err)
	return status.Errorf(codes.FailedPrecondition, "%s", err)
}
