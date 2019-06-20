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
)

// Listener encapsulates runtime for concreate listeners
type Listener interface {
	//
	// Configuration methods
	//
	GetModules() []Module
	SetModules(m []Module)
	AppendModule(m *Module)
	GetConfig() config.ListenerConfig
	SetConfig(config config.ListenerConfig)
	SetHandleRequest(hr modules.Middleware)
	GetHandleRequest() modules.Middleware

	//
	// Server and listenning methods
	//
	Init(server *Server, serverConfig config.ServerConfig, listenerConfig config.ListenerConfig)

	// Blocking call to shutting down a listener
	Shutdown(ctx context.Context) error

	// A channel that indicates whether the listener is ready
	// to server requests. The channel MUST be sent either `true` or `false`
	// values, otherwise initialization may freeze indefineitly
	Ready() chan bool

	// Starts listenning and serving requests.
	// The method MUST return (as opposed to serving in a loop). If a loop is
	// needed, it should be spawned in a separate go routine from within this
	// method. Notice that when listener is ready, the channel returned from
	// Ready() must be sent a `true` value (or `false` upon failure)
	ListenAndServe() error

	//
	// Stats methods
	//
	GetDupDropped() *uint32
}
