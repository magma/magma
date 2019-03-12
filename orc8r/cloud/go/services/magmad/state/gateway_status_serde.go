/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package state

import (
	"fmt"
	"reflect"

	"magma/orc8r/cloud/go/protos"
)

// GatewayStatusSerde manages how to marshal / unmarshal checkin
// response received
type GatewayStatusSerde struct{}

// GetStateType returns the type of state
func (*GatewayStatusSerde) GetStateType() string {
	return "gateway_state"
}

// MarshalState calls the appropriate marshalling method
func (*GatewayStatusSerde) MarshalState(state interface{}) ([]byte, error) {
	castedState, ok := state.(*protos.ServiceStatus)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid magmad gateway state type. Expected *MagmadGatewayState, received %s",
			reflect.TypeOf(state),
		)
	}
	return protos.MarshalIntern(castedState)
}

// UnmarshalState calls the appropriate un-marshalling method
func (*GatewayStatusSerde) UnmarshalState(message []byte) (interface{}, error) {
	response := &protos.ServiceStatus{}
	err := protos.Unmarshal(message, response)
	return response, err
}
