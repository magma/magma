/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/cwf/cloud/go/cwf"
	cwf_service "magma/cwf/cloud/go/services/cwf"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(cwf.ModuleName, cwf_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating cwf service %s", err)
	}
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
