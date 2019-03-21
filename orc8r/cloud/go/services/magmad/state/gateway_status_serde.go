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

type CheckinRequestSerde struct{}

func (*CheckinRequestSerde) GetDomain() string {
	return "state"
}

func (s *CheckinRequestSerde) GetType() string {
	return "checkin_request"
}

func (s *CheckinRequestSerde) Serialize(in interface{}) ([]byte, error) {
	castedState, ok := in.(*protos.CheckinRequest)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid gateway state type. Expected *CheckinRequest, received %s",
			reflect.TypeOf(in),
		)
	}
	return protos.MarshalIntern(castedState)
}

func (s *CheckinRequestSerde) Deserialize(in []byte) (interface{}, error) {
	response := &protos.CheckinRequest{}
	err := protos.Unmarshal(in, response)
	return response, err
}
