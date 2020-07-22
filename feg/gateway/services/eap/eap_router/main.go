/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package main implements eap_router service
package main

import (
	"github.com/golang/glog"
	"golang.org/x/net/context"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap/client"
	"magma/orc8r/lib/go/service"
)

type eapRouter struct {
	supportedMethods []byte
}

func main() {
	// Create the EAP AKA Provider service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.EAP)
	if err != nil {
		glog.Fatalf("Error creating EAP Router service: %s", err)
	}

	protos.RegisterEapRouterServer(srv.GrpcServer, &eapRouter{supportedMethods: client.SupportedTypes()})

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running EAP Router service: %s", err)
	}
}

func (s *eapRouter) HandleIdentity(ctx context.Context, in *protos.EapIdentity) (*protos.Eap, error) {
	resp, err := client.HandleIdentityResponse(uint8(in.GetMethod()), &protos.Eap{Payload: in.Payload, Ctx: in.Ctx})
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		glog.Errorf("HandleIdentity Error: %v", err)
		err = nil
	}
	return resp, err
}

func (s *eapRouter) Handle(ctx context.Context, in *protos.Eap) (*protos.Eap, error) {
	resp, err := client.Handle(in)
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		glog.Errorf("Handle Error: %v", err)
		err = nil
	}
	return resp, err
}

func (s *eapRouter) SupportedMethods(ctx context.Context, in *protos.Void) (*protos.EapMethodList, error) {
	return &protos.EapMethodList{Methods: s.supportedMethods}, nil
}
