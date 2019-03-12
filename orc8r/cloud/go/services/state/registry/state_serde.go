/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry

// StateSerde interface is used by the state service to marshal and
// unmarshal states based on the type of the state stored. There should be
// one implementation registered with the registry per state type.

// To begin storing a new state type in the state service, provide an
// implementation of this interface for each state that you want to store
// and register that implementation with the state service's registry.
type StateSerde interface {

	// Returns the state type that this serde is responsible for.
	// This key is expected to be unique across the whole system.
	GetStateType() string

	// Marshal a state object into a byte array to be persisted by the
	// state service.
	MarshalState(state interface{}) ([]byte, error)

	// Unmarshal a byte array representing a serialized value of this state
	// type into a concrete state.
	UnmarshalState(message []byte) (interface{}, error)
}
