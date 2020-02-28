/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package streamer provides streamer client Go implementation for golang based gateways
package streamer

// Streamer Client Interface
// The package implememntation provides NewStreamerClient(cr registry.CloudRegistry) Client method to create
// New streamer clients
type Client interface {
	// AddListener registers a new streaming updates listener for the
	// listener.GetName() stream.
	// The stream name must be unique and AddListener will error out if a listener
	// for the same stream is already registered.
	AddListener(l Listener) error
	// Stream starts streaming loop for a registered by AddListener listener
	// If successful, Stream never return and should be called in it's own go routine or main()
	// If the provided Listener is not registered, Stream will try to register it prior to starting streaming
	Stream(l Listener) error
	// RemoveListener removes currently registered listener. It returns true is the
	// listener with provided l.Name() exists and was unregistered successfully
	// RemoveListener is the only way to terminate streaming loop
	RemoveListener(l Listener) bool
}
