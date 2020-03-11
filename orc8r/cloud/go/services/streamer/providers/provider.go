/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package providers

import (
	"fmt"
	"sync"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/any"
)

// Interface for a streamer policy. Given a gateway hardware ID, return a
// serialized data bundle of updates to stream back to the gateway
type StreamProvider interface {
	// Name of the stream that this provider services. This name must be
	// globally unique
	GetStreamName() string

	// GetUpdates returns updates to stream updates back to a gateway given its hardware ID
	// if GetUpdates returns error, the stream will be closed without sending any updates
	// if GetUpdates returns error == nil, updates will be sent & the stream will be closed after that
	// if GetUpdates returns error == io.EAGAIN - the returned updates will be sent & GetUpdates will be called again
	// on the same stream
	GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error)
}

type providerRegistry struct {
	sync.RWMutex
	providersByStream map[string]StreamProvider
}

var registry = &providerRegistry{providersByStream: map[string]StreamProvider{}}

// RegisterStreamProviders registers a collection of providers with the
// streamer service. This function will roll back changes if any registration
// fails. This function is thread-safe.
func RegisterStreamProviders(provs ...StreamProvider) error {
	registry.Lock()
	defer registry.Unlock()
	for i, provider := range provs {
		if err := registerUnsafe(provider); err != nil {
			unregisterUnsafe(provs[:i])
			return err
		}
	}
	return nil
}

// Register a stream provider to handle streaming requests for a given stream
// name. If a provider is already registered with the stream name given, this
// function will return an error.
func RegisterStreamProvider(provider StreamProvider) error {
	registry.Lock()
	defer registry.Unlock()
	return registerUnsafe(provider)
}

func registerUnsafe(provider StreamProvider) error {
	newName := provider.GetStreamName()
	_, ok := registry.providersByStream[newName]
	if ok {
		return fmt.Errorf("Stream provider already registered for stream name %s", newName)
	}
	registry.providersByStream[newName] = provider
	return nil
}

func unregisterUnsafe(provs []StreamProvider) {
	for _, provider := range provs {
		delete(registry.providersByStream, provider.GetStreamName())
	}
}

// Get the stream provider for a stream name. Returns an error if no provider
// has been registered for the stream.
func GetStreamProvider(streamName string) (StreamProvider, error) {
	registry.RLock()
	provider, ok := registry.providersByStream[streamName]
	registry.RUnlock()

	if !ok {
		return nil, fmt.Errorf("No provider found for stream %s", streamName)
	}
	return provider, nil
}
