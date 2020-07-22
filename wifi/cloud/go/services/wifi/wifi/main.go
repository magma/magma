/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/service"
	wifi_service "magma/wifi/cloud/go/services/wifi"
	"magma/wifi/cloud/go/services/wifi/obsidian/handlers"
	"magma/wifi/cloud/go/wifi"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(wifi.ModuleName, wifi_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating wifi service %s", err)
	}
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
