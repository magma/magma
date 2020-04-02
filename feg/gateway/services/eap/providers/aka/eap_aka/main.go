/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package main implements Magma EAP AKA Service
package main

import (
	"flag"
	"log"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
	managed_configs "magma/gateway/mconfig"
	"magma/orc8r/lib/go/service"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the EAP AKA Provider service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.EAP_AKA)
	if err != nil {
		log.Fatalf("Error creating EAP AKA service: %s", err)
	}

	akaConfigs := &mconfig.EapAkaConfig{}
	err = managed_configs.GetServiceConfigs(aka.EapAkaServiceName, akaConfigs)
	if err != nil {
		log.Printf("Error getting EAP AKA service configs: %s", err)
		akaConfigs = nil
	}
	servicer, err := servicers.NewEapAkaService(akaConfigs)
	if err != nil {
		log.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	protos.RegisterEapServiceServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running EAP AKA service: %s", err)
	}
}
