/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mocks

import (
	"errors"

	"magma/orc8r/cloud/go/services/config/storage"
	mstore "magma/orc8r/cloud/go/storage"
)

type configTable map[mstore.TypeAndKey][]byte

type mapBackedConfigurationStorage struct {
	// {networkId: {(type, key): value}}
	tables map[string]configTable
}

// Returns an implementation of the ConfigurationStorage interface backed by
// in-memory golang maps for use in unit/integration tests where we don't want
// to spin up a sqlite store.
func NewMapBackedConfigurationStorage() storage.ConfigurationStorage {
	return &mapBackedConfigurationStorage{tables: map[string]configTable{}}
}

func (store *mapBackedConfigurationStorage) GetConfig(networkId string, configType string, key string) (*storage.ConfigValue, error) {
	table, tableExists := store.tables[networkId]
	if !tableExists {
		return &storage.ConfigValue{}, nil
	}
	val, exists := table[mstore.TypeAndKey{Type: configType, Key: key}]
	if !exists {
		return &storage.ConfigValue{}, nil
	}
	return &storage.ConfigValue{Value: val, Version: 0}, nil
}

func (store *mapBackedConfigurationStorage) GetConfigs(networkId string, criteria *storage.FilterCriteria) (map[mstore.TypeAndKey]*storage.ConfigValue, error) {
	if len(criteria.Type) == 0 && len(criteria.Key) == 0 {
		return nil, errors.New("Filter criteria not specified")
	}

	ret := map[mstore.TypeAndKey]*storage.ConfigValue{}
	table, tableExists := store.tables[networkId]
	if !tableExists {
		return ret, nil
	}

	for tk, v := range table {
		if doesTypeAndKeyMatchFilter(tk, criteria) {
			ret[tk] = &storage.ConfigValue{Value: v, Version: 0}
		}
	}
	return ret, nil
}

func (store *mapBackedConfigurationStorage) ListKeysForType(networkId string, configType string) ([]string, error) {
	ret := make([]string, 0)
	table, tableExists := store.tables[networkId]
	if !tableExists {
		return ret, nil
	}

	criteria := &storage.FilterCriteria{Type: configType}
	for tk := range table {
		if doesTypeAndKeyMatchFilter(tk, criteria) {
			ret = append(ret, tk.Key)
		}
	}
	return ret, nil
}

func (store *mapBackedConfigurationStorage) CreateConfig(networkId string, configType string, key string, value []byte) error {
	table := store.getTable(networkId)

	tk := mstore.TypeAndKey{Type: configType, Key: key}
	_, exists := table[tk]
	if exists {
		return errors.New("Creating already existing config")
	}

	table[tk] = value
	store.tables[networkId] = table
	return nil
}

func (store *mapBackedConfigurationStorage) UpdateConfig(networkId string, configType string, key string, newValue []byte) error {
	table := store.getTable(networkId)

	tk := mstore.TypeAndKey{Type: configType, Key: key}
	_, exists := table[tk]
	if !exists {
		return errors.New("Updating nonexistent config")
	}

	table[tk] = newValue
	store.tables[networkId] = table
	return nil
}

func (store *mapBackedConfigurationStorage) DeleteConfig(networkId string, configType string, key string) error {
	table := store.getTable(networkId)

	tk := mstore.TypeAndKey{Type: configType, Key: key}
	_, exists := table[tk]
	if !exists {
		return errors.New("Deleting nonexistent config")
	}

	delete(table, tk)
	store.tables[networkId] = table
	return nil
}

func (store *mapBackedConfigurationStorage) DeleteConfigs(networkId string, criteria *storage.FilterCriteria) error {
	table := store.getTable(networkId)

	var tksToDelete []mstore.TypeAndKey
	for tk := range table {
		if doesTypeAndKeyMatchFilter(tk, criteria) {
			tksToDelete = append(tksToDelete, tk)
		}
	}

	for _, tk := range tksToDelete {
		delete(table, tk)
	}
	store.tables[networkId] = table
	return nil
}

func (store *mapBackedConfigurationStorage) DeleteConfigsForNetwork(networkId string) error {
	delete(store.tables, networkId)
	return nil
}

func doesTypeAndKeyMatchFilter(tk mstore.TypeAndKey, filter *storage.FilterCriteria) bool {
	match := true
	if len(filter.Type) > 0 {
		match = match && (tk.Type == filter.Type)
	}
	if len(filter.Key) > 0 {
		match = match && (tk.Key == filter.Key)
	}
	return match
}

func (store *mapBackedConfigurationStorage) getTable(networkId string) configTable {
	table, exists := store.tables[networkId]
	if !exists {
		return configTable{}
	}
	return table
}
