/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package main implements eap_router service
package main

import (
	"log"

	"golang.org/x/net/context"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/eap/client"
	"magma/feg/gateway/services/eap/protos"
	"magma/orc8r/cloud/go/service"
)

type eapRouter struct {
	supportedMethods []byte
}

func main() {
	// Create the EAP AKA Provider service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.EAP)
	if err != nil {
		log.Fatalf("Error creating EAP Router service: %s", err)
	}

	protos.RegisterEapRouterServer(srv.GrpcServer, &eapRouter{supportedMethods: client.SupportedTypes()})

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running EAP Router service: %s", err)
	}
}

func (s *eapRouter) HandleIdentity(ctx context.Context, in *protos.EapIdentity) (*protos.Eap, error) {
	return client.HandleIdentityResponse(uint8(in.GetMethod()), &protos.Eap{Payload: in.Payload, Ctx: in.Ctx})
}

func (s *eapRouter) Handle(ctx context.Context, in *protos.Eap) (*protos.Eap, error) {
	return client.Handle(in)
}

func (s *eapRouter) SupportedMethods(ctx context.Context, in *protos.Void) (*protos.MethodList, error) {
	return &protos.MethodList{Methods: s.supportedMethods}, nil
}
