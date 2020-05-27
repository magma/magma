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
	"net"
	"strings"

	mconfigprotos "magma/cwf/cloud/go/protos/mconfig"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/gateway_health/health/gre_probe"
	"magma/cwf/gateway/services/gateway_health/health/service_health"
	"magma/cwf/gateway/services/gateway_health/health/system_health"
	"magma/cwf/gateway/services/gateway_health/servicers"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/gateway/mconfig"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

const (
	defaultMemUtilPct       = 0.75
	defaultCpuUtilPct       = 0.75
	defaultGREProbeInterval = 10
	defaultICMPPktCount     = 3
	defaultICMPInterface    = "eth1"
)

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.GatewayHealth)
	if err != nil {
		glog.Fatalf("Error creating %s service: %s", registry.GatewayHealth, err)
	}
	cfg := getHealthMconfig()
	probe := gre_probe.NewICMPProbe(cfg.GrePeers, cfg.GreProbeInterval, int(cfg.IcmpProbePktCount))
	systemHealth, err := system_health.NewCWAGSystemHealthProvider(defaultICMPInterface)
	if err != nil {
		glog.Fatalf("Error creating CWAGServiceHealthProvider: %s", err)
	}
	dockerHealth, err := service_health.NewDockerServiceHealthProvider()
	if err != nil {
		glog.Fatalf("Error creating DockerServiceHealthProvider: %s", err)
	}
	servicer := servicers.NewGatewayHealthServicer(cfg, probe, dockerHealth, systemHealth)
	fegprotos.RegisterServiceHealthServer(srv.GrpcServer, servicer)

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

func getHealthMconfig() *mconfigprotos.CwfGatewayHealthConfig {
	ret := &mconfigprotos.CwfGatewayHealthConfig{}
	err := mconfig.GetServiceConfigs(strings.ToLower(registry.GatewayHealth), ret)
	if err != nil {
		ret.CpuUtilThresholdPct = defaultCpuUtilPct
		ret.MemUtilThresholdPct = defaultMemUtilPct
		ret.GreProbeInterval = defaultGREProbeInterval
		ret.IcmpProbePktCount = defaultICMPPktCount
		glog.Errorf("Could not load mconfig. Using defaults: %v", ret)
		return ret
	}
	if ret.CpuUtilThresholdPct == 0 {
		ret.CpuUtilThresholdPct = defaultCpuUtilPct
	}
	if ret.MemUtilThresholdPct == 0 {
		ret.MemUtilThresholdPct = defaultMemUtilPct
	}
	if ret.GreProbeInterval == 0 {
		ret.GreProbeInterval = defaultGREProbeInterval
	}
	ret.GrePeers = removeCIDREndpoints(ret.GrePeers)
	glog.Infof("Using config: %v", ret)
	return ret
}

// For now, we don't support querying all endpoints defined in a CIDR
// formatted endpoint.
func removeCIDREndpoints(endpoints []*mconfigprotos.CwfGatewayHealthConfigGrePeer) []*mconfigprotos.CwfGatewayHealthConfigGrePeer {
	ret := []*mconfigprotos.CwfGatewayHealthConfigGrePeer{}
	for _, endpoint := range endpoints {
		_, _, err := net.ParseCIDR(endpoint.Ip)
		if err != nil {
			ret = append(ret, endpoint)
			continue
		}
		glog.Infof("Not monitoring CIDR formatted IP: %s. Health service only supports monitoring specific endpoints", endpoint.Ip)
	}
	return ret
}
