/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Magma's Swx Proxy Service converts gRPC requests into Swx protocol over diameter
package main

import (
	"flag"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.SWX_PROXY)
	if err != nil {
		glog.Fatalf("Error creating Swx Proxy service: %s", err)
	}

	// TODO: remove when multiple config is supported
	configs := []*servicers.SwxProxyConfig{servicers.GetSwxProxyConfig()}

	servicer, err := servicers.NewSwxProxiesWithHealthAndDefaultMultiplexor(configs)
	if err != nil {
		glog.Fatalf("Failed to create SwxProxy: %v", err)
	}
	protos.RegisterSwxProxyServer(srv.GrpcServer, servicer)
	protos.RegisterServiceHealthServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
