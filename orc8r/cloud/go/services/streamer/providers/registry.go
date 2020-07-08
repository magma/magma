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
)

type remoteRegistry struct {
	providersByStream map[string]StreamProvider
	sync.RWMutex
}

var reg = &remoteRegistry{providersByStream: map[string]StreamProvider{}}

// RegisterStreamProviders registers a collection of providers with the
// streamer service. This function will roll back changes if any registration
// fails. This function is thread-safe.
func RegisterStreamProviders(provs ...StreamProvider) error {
	reg.Lock()
	defer reg.Unlock()
	for i, provider := range provs {
		if err := registerUnsafe(provider); err != nil {
			unregisterUnsafe(provs[:i])
			return err
		}
	}
	return nil
}

// RegisterStreamProvider registers a stream provider to handle streaming
// requests for a given stream name. If a provider is already registered with
// the stream name given, this function will return an error.
func RegisterStreamProvider(provider StreamProvider) error {
	reg.Lock()
	defer reg.Unlock()
	return registerUnsafe(provider)
}

// GetStreamProvider gets the stream provider for a stream name. Returns an
// error if no provider has been registered for the stream.
func GetStreamProvider(streamName string) (StreamProvider, error) {
	reg.RLock()
	provider, ok := reg.providersByStream[streamName]
	reg.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no provider found for stream %s", streamName)
	}
	return provider, nil
}

func registerUnsafe(provider StreamProvider) error {
	newName := provider.GetStreamName()
	_, ok := reg.providersByStream[newName]
	if ok {
		return fmt.Errorf("stream provider already registered for stream name %s", newName)
	}
	reg.providersByStream[newName] = provider
	return nil
}

func unregisterUnsafe(provs []StreamProvider) {
	for _, provider := range provs {
		delete(reg.providersByStream, provider.GetStreamName())
	}
}
