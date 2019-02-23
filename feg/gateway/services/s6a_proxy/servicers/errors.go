/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Errorf(code codes.Code, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	log.Printf("RPC [%s] %s", code, msg)
	return status.Errorf(code, msg)
}

func Error(code codes.Code, err error) error {
	log.Printf("RPC [%s] %s", code, err)
	return status.Error(code, err.Error())
}
