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

type stateSerdeRegistry struct {
	sync.RWMutex
	managers map[string]StateSerde
}

var registry = stateSerdeRegistry{
	managers: map[string]StateSerde{},
}

// RegisterStateSerdes registers a collection of State
// serializer/deserializer with the State registry.
// This will return an error if a collision occurs, and will roll back changes
// before returning. I.e. the semantics are all-or-nothing.
// This function is thread-safe.
func RegisterStateSerdes(managers ...StateSerde) error {
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
func MarshalState(stateType string, state interface{}) ([]byte, error) {
	registry.RLock()
	defer registry.RUnlock()

	manager, exists := registry.managers[stateType]
	if !exists {
		return nil, fmt.Errorf("No state serde is registered under %s", stateType)
	}
	return manager.MarshalState(state)
}

func UnmarshalState(stateType string, message []byte) (interface{}, error) {
	if message == nil || len(message) == 0 {
		return nil, nil
	}

	registry.RLock()
	defer registry.RUnlock()
	manager, exists := registry.managers[stateType]
	if !exists {
		return nil, fmt.Errorf("No State serde is registered under %s", stateType)
	}
	return manager.UnmarshalState(message)
}

func registerUnsafe(manager StateSerde) error {
	StateType := manager.GetStateType()
	_, exists := registry.managers[StateType]
	if exists {
		return fmt.Errorf("State serde for type %s already exists", StateType)
	}
	registry.managers[StateType] = manager
	return nil
}

func unregisterUnsafe(managers []StateSerde) {
	for _, man := range managers {
		delete(registry.managers, man.GetStateType())
	}
}
