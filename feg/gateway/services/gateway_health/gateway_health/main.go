/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"time"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/gateway_health/health_manager"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

const healthIntervalSec = 10

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.HEALTH)
	if err != nil {
		glog.Fatalf("Error creating HEALTH service: %s", err)
	}

	cloudReg := registry.NewCloudRegistry()
	healthManager := health_manager.NewHealthManager(cloudReg)
	// Run Health Collection Loop
	go func() {
		for {
			<-time.After(healthIntervalSec * time.Second)
			healthManager.SendHealthUpdate()
		}
	}()

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
