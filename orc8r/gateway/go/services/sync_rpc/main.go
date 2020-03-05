/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package main - implementation of a stand alone syncrpcclient
package main

import (
	"magma/gateway/services/sync_rpc/service"
	"time"
)

func main() {
	// hardcode initial values similar to what we have currently in
	// sync_rpc_client.py. Replace it later with yaml config file
	cfg := service.Config{
		SyncRpcHeartbeatInterval: 30 * time.Second,
		GatewayKeepaliveInterval: 10 * time.Second,
		GatewayResponseTimeout:   120 * time.Second,
	}

	syncRpcService := service.NewSyncRpcClient(&cfg)
	syncRpcService.Run()
}
