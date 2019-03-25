/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package main implements eap_router service
package main

import (
	"flag"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"magma/feg/gateway/services/eap/client"
	"magma/feg/gateway/services/eap/protos"
)

type eapRouter struct {
	supportedMethods []byte
}

func main() {
	addr := flag.String("addr", ":11111", "Server address (host:port)")
	flag.Parse()

	log.Printf("Starting EAP Router on tcp: %s", *addr)
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to create listener on address '%s': %v", *addr, err)
	}

	grpcServer := grpc.NewServer()
	protos.RegisterEapRouterServer(grpcServer, &eapRouter{supportedMethods: client.SupportedTypes()})
	grpcServer.Serve(lis)
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
