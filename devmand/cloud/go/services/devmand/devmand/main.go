/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/orc8r/cloud/go/service"
	"orc8r/devmand/cloud/go/devmand"
	devmand_service "orc8r/devmand/cloud/go/services/devmand"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(devmand.ModuleName, devmand_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating devmand service %s", err)
	}
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
