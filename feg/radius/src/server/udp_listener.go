/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package server

import (
	"context"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fmt"
)

// UDPListener listens to Radius udp packets
type UDPListener struct {
	Server        *radius.PacketServer
	Config        config.ListenerConfig
	Modules       []Module
	HandleRequest modules.Middleware
	dupDropped    uint32
	ready         chan bool
}

// GetModules override
func (l *UDPListener) GetModules() []Module {
	return l.Modules
}

// SetModules override
func (l *UDPListener) SetModules(m []Module) {
	l.Modules = m
}

// AppendModule override
func (l *UDPListener) AppendModule(m *Module) {
	l.Modules = append(l.Modules, *m)
}

// GetConfig override
func (l *UDPListener) GetConfig() config.ListenerConfig {
	return l.Config
}

// SetHandleRequest override
func (l *UDPListener) SetHandleRequest(hr modules.Middleware) {
	l.HandleRequest = hr
}

// Init override
func (l *UDPListener) Init(server *Server, serverConfig config.ServerConfig, listenerConfig config.ListenerConfig) {

	l.SetConfig(listenerConfig)

	// Create packet server
	l.Server = &radius.PacketServer{
		Handler: radius.HandlerFunc(
			generatePacketHandler(l, server),
		),
		SecretSource: radius.StaticSecretSource([]byte(serverConfig.Secret)),
		Addr:         fmt.Sprintf(":%d", l.GetConfig().Port),
		Ready:        make(chan bool),
	}

	l.ready = make(chan bool, 1)
}

// ListenAndServe override
func (l *UDPListener) ListenAndServe() error {
	serverError := make(chan error, 1)
	go func() {
		err := l.Server.ListenAndServe()
		serverError <- err
	}()

	// Wait to see if initialization was successful
	select {
	case _ = <-l.Server.Ready:
		l.ready <- true
		return nil
	case err := <-serverError:
		l.ready <- false
		return err // might be nil if no error
	}
}

// GetHandleRequest override
func (l *UDPListener) GetHandleRequest() modules.Middleware {
	return l.HandleRequest
}

// Shutdown override
func (l *UDPListener) Shutdown(ctx context.Context) error {
	return l.Server.Shutdown(ctx)
}

// GetDupDropped override
func (l *UDPListener) GetDupDropped() *uint32 {
	return &l.dupDropped
}

// Ready override
func (l *UDPListener) Ready() chan bool {
	return l.ready
}

// SetConfig override
func (l *UDPListener) SetConfig(c config.ListenerConfig) {
	l.Config = c
}
