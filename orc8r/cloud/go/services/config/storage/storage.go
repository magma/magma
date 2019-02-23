/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

// ConfigurationStorage defines the persistence interface for config service.
// Implementors are expected to manage configurations given as arbritrary
// bytes, uniquely keyed by a string tuple of (type, key).
type ConfigurationStorage interface {
	// Fetch a specific config. If no such config value is found, this will
	// return an empty ConfigValue.
	GetConfig(networkId string, configType string, key string) (*ConfigValue, error)

	// Fetch configs given some filter criteria. At least one field of the
	// input filter criteria must be specified.
	GetConfigs(networkId string, criteria *FilterCriteria) (map[TypeAndKey]*ConfigValue, error)

	// List all keys for a given config type
	ListKeysForType(networkId string, configType string) ([]string, error)

	CreateConfig(networkId string, configType string, key string, value []byte) error
	UpdateConfig(networkId string, configType string, key string, newValue []byte) error
	DeleteConfig(networkId string, configType string, key string) error

	// Delete all configs matching a filter criteria. At least one field of the
	// input filter criteria must be specified.
	DeleteConfigs(networkId string, criteria *FilterCriteria) error

	// Delete all configs for a network (drop the table)
	DeleteConfigsForNetwork(networkId string) error
}

type TypeAndKey struct {
	Type string
	Key  string
}

type ConfigValue struct {
	Value   []byte
	Version uint64
}

// FilterCriteria specifies a matching for a configuration's (type, key)
// identifier.
type FilterCriteria struct {
	Type string
	Key  string
}
