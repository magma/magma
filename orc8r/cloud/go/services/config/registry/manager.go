/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry

// ConfigManager interface is used by the config service to marshal and
// unmarshal configs based on the type of the config stored. There should be
// one implementation registered with the registry per config type.

// To begin storing a new config type in the config service, provide an
// implementation of this interface for each config that you want to store
// and register that implementation with the config service's registry.
type ConfigManager interface {

	// Returns the config type that this manager is responsible for.
	// This key is expected to be unique across the whole system.
	GetConfigType() string

	// Returns the gateway IDs that a specified config applies to. For example,
	// A network-level config would return a list of all gateway IDs in the
	// network and a gateway-level config would return a list with just
	// the configKey parameter.
	GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error)

	// Marshal a config object into a byte array to be persisted by the
	// config service.
	MarshalConfig(config interface{}) ([]byte, error)
	// Unmarshal a byte array representing a serialized value of this config
	// type into a concrete config.
	UnmarshalConfig(message []byte) (interface{}, error)
}
