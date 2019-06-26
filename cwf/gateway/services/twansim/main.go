/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"flag"

	"magma/cwf/gateway/registry"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.TwanSim)
	if err != nil {
		glog.Fatalf("Error creating TwanSim service: %s", err)
	}

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
