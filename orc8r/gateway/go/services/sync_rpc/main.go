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
)

func main() {
	// start sync RPC client with default configuration
	syncRpcService := service.NewClient(nil)
	syncRpcService.Run()
}
