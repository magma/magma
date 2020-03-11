/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/directoryd"

	"github.com/golang/glog"
)

// NOTE: directoryd service currently attaches no handlers.
// The service is preserved for future plans related to custom indexers.
func main() {
	// Create Magma micro-service
	directoryService, err := service.NewOrchestratorService(orc8r.ModuleName, directoryd.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating directory service: %s", err)
	}

	// Run the service
	glog.V(2).Info("Starting Directory Service...")
	err = directoryService.Run()
	if err != nil {
		glog.Fatalf("Error running directory service: %s", err)
	}
}
