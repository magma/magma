/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/radius/metrics"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.RADIUS)
	if err != nil {
		glog.Fatalf("Error creating RADIUS service: %s", err)
	}

	// The actual service TBD ...

	// Example of incrementing a metric counter
	metrics.TotalRequests.Inc()

	// Run the service - for Radius has only built in FB303 Service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running RADIUS service: %s", err)
	}
}
