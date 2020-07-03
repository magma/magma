/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, lte_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating lte service %s", err)
	}
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
