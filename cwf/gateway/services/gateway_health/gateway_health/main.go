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
	"os"
	"strings"
	"time"

	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/gateway_health/health/gre_probe"
	"magma/cwf/gateway/services/gateway_health/health/service_health"
	"magma/cwf/gateway/services/gateway_health/health/system_health"
	"magma/cwf/gateway/services/gateway_health/servicers"
	"magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

const (
	defaultMemUtilPct       = 0.75
	defaultCpuUtilPct       = 0.75
	defaultGREProbeInterval = 10 * time.Second
)

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.GatewayHealth)
	if err != nil {
		glog.Fatalf("Error creating %s service: %s", registry.GatewayHealth, err)
	}
	//TODO: Update with mconfig values
	endpointStr := os.Getenv("GRE_ENDPOINTS")
	greEndpoints := strings.Split(endpointStr, ",")
	probe := &gre_probe.DummyGREProbe{}
	systemHealth := &system_health.DummySystemStatsProvider{}
	serviceHealth := &service_health.DummyServiceHealthProvider{}
	cfg := &servicers.HealthConfig{
		GrePeers:      greEndpoints,
		MaxCpuUtilPct: defaultMemUtilPct,
		MaxMemUtilPct: defaultCpuUtilPct,
	}
	servicer := servicers.NewGatewayHealthServicer(cfg, probe, serviceHealth, systemHealth)
	protos.RegisterServiceHealthServer(srv.GrpcServer, servicer)

	// Start GRE probe
	err = probe.Start()
	if err != nil {
		glog.Fatalf("Error running GRE health probe: %s", err)
	}
	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running %s service: %s", registry.GatewayHealth, err)
	}
}
