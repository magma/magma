/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetIdleSessionTimeout returns Idle Session Timeout Duration if set in mconfigs or DefaultSessionTimeout otherwise
func GetIdleSessionTimeout(cfg *mconfig.AAAConfig) time.Duration {
	if cfg != nil {
		if tout := time.Millisecond * time.Duration(cfg.GetIdleSessionTimeoutMs()); tout > 0 {
			return tout
		}
	}
	return aaa.DefaultSessionTimeout
}

func isThruthy(value string) bool {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return false
	}
	value = strings.ToLower(value)
	if value == "0" || strings.HasPrefix(value, "false") || strings.HasPrefix(value, "n") {
		return false
	}
	return true
}

func Errorf(code codes.Code, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	log.Printf("%s; [RPC: %s]", msg, code.String())
	return status.Errorf(code, msg)
}

func Error(code codes.Code, err error) error {
	if err != nil {
		if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
			code = se.GRPCStatus().Code()
		}
		log.Printf("%v; [RPC: %s]", err, code.String())
		return status.Error(code, err.Error())
	}
	return nil
}
