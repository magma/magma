/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"sync"

	"magma/orc8r/cloud/go/protos"

	"github.com/golang/protobuf/ptypes/any"
)

type mconfigStreamerRegistry struct {
	sync.RWMutex
	streamersByConfigType map[string][]MconfigStreamer
	allStreamers          []MconfigStreamer
}

var mconfigRegistry = &mconfigStreamerRegistry{
	streamersByConfigType: map[string][]MconfigStreamer{},
	allStreamers:          []MconfigStreamer{},
}

// RegisterMconfigStreamers registers a collection of mconfig streamers.
// This function is thread-safe.
func RegisterMconfigStreamers(streamers ...MconfigStreamer) {
	mconfigRegistry.Lock()
	defer mconfigRegistry.Unlock()
	for _, streamer := range streamers {
		registerUnsafe(streamer)
	}
}

// Register an mconfig streamer.
func RegisterMconfigStreamer(streamer MconfigStreamer) {
	mconfigRegistry.Lock()
	defer mconfigRegistry.Unlock()
	registerUnsafe(streamer)
}

func registerUnsafe(streamer MconfigStreamer) {
	subscribedTypes := streamer.GetSubscribedConfigTypes()
	for _, t := range subscribedTypes {
		existingCollection, ok := mconfigRegistry.streamersByConfigType[t]
		if !ok {
			mconfigRegistry.streamersByConfigType[t] = []MconfigStreamer{streamer}
		} else {
			mconfigRegistry.streamersByConfigType[t] = append(existingCollection, streamer)
		}
	}

	mconfigRegistry.allStreamers = append(mconfigRegistry.allStreamers, streamer)
}

func GetMconfigForNewGateway(networkId string, gatewayId string) (*protos.GatewayConfigs, error) {
	mconfigRegistry.RLock()
	defer mconfigRegistry.RUnlock()

	ret := &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	for _, streamer := range mconfigRegistry.allStreamers {
		err := streamer.SeedNewGatewayMconfig(networkId, gatewayId, ret)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

// Get mconfig updates for a ConfigUpdate event
func ApplyMconfigUpdate(
	update *ConfigUpdate,
	oldMconfigs map[string]*protos.GatewayConfigs,
) (map[string]*protos.GatewayConfigs, error) {
	mconfigRegistry.RLock()
	defer mconfigRegistry.RUnlock()

	streamers, ok := mconfigRegistry.streamersByConfigType[update.ConfigType]
	if !ok {
		// No streamers registered for this config so no op needed
		return map[string]*protos.GatewayConfigs{}, nil
	}

	previousStateInput := oldMconfigs
	for _, streamer := range streamers {
		updatedMconfigs, err := streamer.ApplyMconfigUpdate(update, previousStateInput)
		if err != nil {
			return map[string]*protos.GatewayConfigs{}, err
		}
		previousStateInput = updatedMconfigs
	}
	return previousStateInput, nil
}

// ONLY USE FOR TESTING
func ClearRegistryForTesting() {
	mconfigRegistry.streamersByConfigType = map[string][]MconfigStreamer{}
}
