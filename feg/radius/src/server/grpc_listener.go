/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package server

import (
	"context"
	"errors"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/coa/protos"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// GRPCListener listens to gRpc
type GRPCListener struct {
	Server        *grpc.Server
	Config        config.ListenerConfig
	Modules       []Module
	HandleRequest modules.Middleware
	dupDropped    uint32
	ready         chan bool
}

// GetModules override
func (l *GRPCListener) GetModules() []Module {
	return l.Modules
}

// SetModules override
func (l *GRPCListener) SetModules(m []Module) {
	l.Modules = m
}

// AppendModule override
func (l *GRPCListener) AppendModule(m *Module) {
	l.Modules = append(l.Modules, *m)
}

// GetConfig override
func (l *GRPCListener) GetConfig() config.ListenerConfig {
	return l.Config
}

// SetHandleRequest override
func (l *GRPCListener) SetHandleRequest(hr modules.Middleware) {
	l.HandleRequest = hr
}

// Init override
func (l *GRPCListener) Init(server *Server, serverConfig config.ServerConfig, listenerConfig config.ListenerConfig) {
	l.ready = make(chan bool, 1)
}

// ListenAndServe override
func (l *GRPCListener) ListenAndServe() error {
	// Start listenning
	listenAddress := fmt.Sprintf(":%d", l.GetConfig().Port)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		l.ready <- false
		return errors.New("grpc listener: failed to open tcp connection" + listenAddress)
	}

	// Start serving
	l.Server = grpc.NewServer()
	protos.RegisterAuthorizationServer(l.Server, &authorizationServer{Listener: l})
	go func() {
		l.Server.Serve(lis)
	}()

	// Signal listener is ready
	go func() {
		l.ready <- true
	}()
	return nil
}

// GetHandleRequest override
func (l *GRPCListener) GetHandleRequest() modules.Middleware {
	return l.HandleRequest
}

// Shutdown override
func (l *GRPCListener) Shutdown(ctx context.Context) error {
	return nil
}

// GetDupDropped override
func (l *GRPCListener) GetDupDropped() *uint32 {
	return &l.dupDropped
}

// Ready override
func (l *GRPCListener) Ready() chan bool {
	return l.ready
}

// SetConfig override
func (l *GRPCListener) SetConfig(c config.ListenerConfig) {
	l.Config = c
}

type authorizationServer struct {
	Listener *GRPCListener
}

func (s *authorizationServer) Change(ctx context.Context, request *protos.ChangeRequest) (*protos.CoaResponse, error) {

	//@todo #T45624133 - parse packet and do something similar to generatePacketHandler
	return &protos.CoaResponse{CoaResponseType: protos.CoaResponse_NAK, Ctx: request.Ctx}, nil
}

func (s *authorizationServer) Disconnect(ctx context.Context, request *protos.DisconnectRequest) (*protos.CoaResponse, error) {

	//@todo - #T45624133 parse packet and do something similar to generatePacketHandler
	return &protos.CoaResponse{CoaResponseType: protos.CoaResponse_ACK, Ctx: request.Ctx}, nil
}
