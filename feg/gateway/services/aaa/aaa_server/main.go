/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package main implements WiFi AAA server
package main

import (
	"log"

	"github.com/golang/protobuf/proto"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/servicers"
	"magma/feg/gateway/services/aaa/store"
	managed_configs "magma/gateway/mconfig"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/service"
)

const (
	AAAServiceName = "aaa_server"
	Version        = "0.1"
)

func main() {
	// Create a shared Session Table
	sessions := store.NewMemorySessionTable()

	// Create the EAP AKA Provider service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.AAA_SERVER)
	if err != nil {
		log.Fatalf("Error creating AAA service: %s", err)
	}
	aaaConfigs := &mconfig.AAAConfig{}
	err = managed_configs.GetServiceConfigs(AAAServiceName, aaaConfigs)
	if err != nil {
		log.Printf("Error getting AAA Server service configs: %s", err)
		aaaConfigs = nil
	}
	acct, _ := servicers.NewAccountingService(sessions, proto.Clone(aaaConfigs).(*mconfig.AAAConfig))
	protos.RegisterAccountingServer(srv.GrpcServer, acct)
	lteprotos.RegisterAbortSessionResponderServer(srv.GrpcServer, acct)
	fegprotos.RegisterSwxGatewayServiceServer(srv.GrpcServer, acct)

	auth, _ := servicers.NewEapAuthenticator(sessions, aaaConfigs, acct)
	protos.RegisterAuthenticatorServer(srv.GrpcServer, auth)

	log.Printf("Starting AAA Service v%s.", Version)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running AAA service: %s", err)
	}
}
