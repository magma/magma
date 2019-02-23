/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry

import (
	"fmt"
	"sync"
)

type configRegistry struct {
	sync.RWMutex
	managers map[string]ConfigManager
}

var registry = configRegistry{
	managers: map[string]ConfigManager{},
}

// RegisterConfigManagers registers a collection of configuration managers
// with the config registry.
// This will return an error if a collision occurs, and will roll back changes
// before returning. I.e. the semantics are all-or-nothing.
// This function is thread-safe.
func RegisterConfigManagers(managers ...ConfigManager) error {
	registry.Lock()
	defer registry.Unlock()

	for i, man := range managers {
		if err := registerUnsafe(man); err != nil {
			unregisterUnsafe(managers[:i])
			return err
		}
	}
	return nil
}

// Register a config manager with the config registry. This will return an error
// if there already is a config manager registered for the type defined by the
// manager being registered.
// All config registry functions are thread-safe via a read-write lock.
func RegisterConfigManager(manager ConfigManager) error {
	registry.Lock()
	defer registry.Unlock()
	return registerUnsafe(manager)
}

func registerUnsafe(manager ConfigManager) error {
	configType := manager.GetConfigType()
	_, exists := registry.managers[configType]
	if exists {
		return fmt.Errorf("Config manager for type %s already exists", configType)
	}
	registry.managers[configType] = manager
	return nil
}

func unregisterUnsafe(managers []ConfigManager) {
	for _, man := range managers {
		delete(registry.managers, man.GetConfigType())
	}
}

func MarshalConfig(configType string, config interface{}) ([]byte, error) {
	registry.RLock()
	defer registry.RUnlock()

	manager, exists := registry.managers[configType]
	if !exists {
		return nil, fmt.Errorf("No config manager is registered under %s", configType)
	}
	return manager.MarshalConfig(config)
}

func UnmarshalConfig(configType string, message []byte) (interface{}, error) {
	if message == nil || len(message) == 0 {
		return nil, nil
	}

	registry.RLock()
	defer registry.RUnlock()
	manager, exists := registry.managers[configType]
	if !exists {
		return nil, fmt.Errorf("No config manager is registered under %s", configType)
	}
	return manager.UnmarshalConfig(message)
}

func GetGatewayIdsForConfig(configType string, networkId string, configKey string) ([]string, error) {
	registry.RLock()
	defer registry.RUnlock()
	manager, exists := registry.managers[configType]
	if !exists {
		return nil, fmt.Errorf("No config manager is registered under %s", configType)
	}
	return manager.GetGatewayIdsForConfig(networkId, configKey)
}

// ONLY USE FOR TESTING
func ClearRegistryForTesting() {
	registry.Lock()
	defer registry.Unlock()
	registry.managers = map[string]ConfigManager{}
}
