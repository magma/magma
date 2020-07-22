/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/devmand/cloud/go/devmand"
	devmand_service "magma/devmand/cloud/go/services/devmand"
	"magma/devmand/cloud/go/services/devmand/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(devmand.ModuleName, devmand_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating devmand service %s", err)
	}
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
