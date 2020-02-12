/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package mconfig provides gateway Go support for cloud managed configuration (mconfig)
package mconfig

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

func GetServiceConfigs(service string, result proto.Message) error {
	current := GetGatewayConfigs()
	anyCfg, found := current.ConfigsByKey[service]
	if !found {
		cfgMu.Lock()
		defer cfgMu.Unlock()
		return fmt.Errorf("No configs found for service: '%s' in %s", service, lastFilePath)
	}

	return ptypes.UnmarshalAny(anyCfg, result)
}

func GetGatewayConfigs() *protos.GatewayConfigs {
	current := (*protos.GatewayConfigs)(atomic.LoadPointer(&localConfig))
	if current == nil {
		// initial refresh, only do it one time
		RefreshConfigs()
		// Swap with an empty configs obj if localConfig is still nil, use CompareAndSwap to not to overwrite
		// the result of concurrent, successful refresh
		atomic.CompareAndSwapPointer(&localConfig, nil, (unsafe.Pointer)(&protos.GatewayConfigs{}))
		// Return the latest value of localConfig
		return (*protos.GatewayConfigs)(atomic.LoadPointer(&localConfig))
	}
	return current
}
