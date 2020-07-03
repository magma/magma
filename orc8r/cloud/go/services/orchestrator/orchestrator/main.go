/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/orchestrator"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, orchestrator.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating orchestrator service %s", err)
	}
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
