/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"time"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"
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
