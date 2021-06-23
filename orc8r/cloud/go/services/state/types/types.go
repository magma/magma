/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package types contains the state service's types.
//
// State is variously stored in three related objects:
//	- protos.State		- state as it's reported from gateways
//	- SerializedState	- wraps reported state with orc8r-defined metadata
//	- State				- deserializes reported state contents to native object
//
// The state service eagerly converts protos.State to SerializedSTate as soon
// as possible. From there, consumers of the client API provide their own serde
// registry to guide deserializing to differentiated State objects.
package types

import (
	"encoding/json"
	"fmt"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	// SerdeDomain is a copy of the state service's serde domain, from
	// orc8r/cloud/go/services/state/doc.go.
	SerdeDomain = "state"
)

// ID identifies a piece of state.
// A piece of state is uniquely identified by the triplet
// {network ID, device ID, type}.
type ID struct {
	// Type determines how the value is deserialized and validated.
	Type string
	// DeviceID is the ID of the entity with which the state is associated
	// (IMSI, serial number, etc).
	DeviceID string
}

// IDs are a list of state IDs.
type IDs []ID

// MakeIDs converts a list of device IDs to a list of state IDs.
func MakeIDs(typ string, keys ...string) IDs {
	var ids IDs
	for _, key := range keys {
		ids = append(ids, ID{Type: typ, DeviceID: key})
	}
	return ids
}

// IDsByNetwork are a set of state IDs, keyed by network ID.
type IDsByNetwork map[string]IDs

// SerializedState directly wraps serialized reported state with metadata.
type SerializedState struct {
	// SerializedReportedState is the actual state reported by the device,
	// serialized.
	SerializedReportedState []byte

	// Version is the reported version of the state.
	Version uint64
	// ReporterID is the hardware ID of the gateway which reported the state.
	ReporterID string
	// TimeMs is the time the state was received in milliseconds.
	TimeMs uint64
}

// SerializedStatesByID maps state IDs to state.
// A state and ID collectively contain all information for a piece of state.
type SerializedStatesByID map[ID]SerializedState

// State directly wraps deserialized reported state with metadata.
type State struct {
	// ReportedState is the actual state reported by the device, deserialized.
	ReportedState interface{}

	// Version is the reported version of the state.
	Version uint64
	// ReporterID is the hardware ID of the gateway which reported the state.
	ReporterID string
	// TimeMs is the time the state was received in milliseconds.
	TimeMs uint64
	// CertExpirationTime is the expiration time in milliseconds.
	CertExpirationTime int64
}

// StatesByID maps state IDs to state.
// A state and ID collectively contain all information for a piece of state.
type StatesByID map[ID]State

// StateErrors is a mapping of state ID to error experienced handling the state.
type StateErrors map[ID]error

// Filter to the subset that match one of the state types.
func (ids IDs) Filter(types ...string) IDs {
	var ret IDs
	for _, id := range ids {
		if funk.ContainsString(types, id.Type) {
			ret = append(ret, id)
		}
	}
	return ret
}

// Filter to the subset that match one of the state types.
func (s StatesByID) Filter(types ...string) StatesByID {
	ret := StatesByID{}
	for id, st := range s {
		if funk.ContainsString(types, id.Type) {
			ret[id] = st
		}
	}
	return ret
}

// Filter to the subset that match one of the state types.
func (s SerializedStatesByID) Filter(types ...string) SerializedStatesByID {
	ret := SerializedStatesByID{}
	for id, st := range s {
		if funk.ContainsString(types, id.Type) {
			ret[id] = st
		}
	}
	return ret
}

// MakeProtoStateErrors converts state errors to proto state errors.
func MakeProtoStateErrors(errs StateErrors) []*protos.IDAndError {
	var ret []*protos.IDAndError
	for id, e := range errs {
		ret = append(ret, &protos.IDAndError{Type: id.Type, DeviceID: id.DeviceID, Error: e.Error()})
	}
	return ret
}

// MakeStateErrors converts proto state errors to native state errors.
func MakeStateErrors(errs []*protos.IDAndError) StateErrors {
	ret := StateErrors{}
	for _, e := range errs {
		id := ID{Type: e.Type, DeviceID: e.DeviceID}
		ret[id] = fmt.Errorf("%v", e.Error)
	}
	return ret
}

// MakeProtoStates converts states by ID to proto states.
func MakeProtoStates(states SerializedStatesByID) ([]*protos.State, error) {
	var ret []*protos.State
	for id, st := range states {
		p, err := MakeProtoState(id, st)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

// MakeProtoState converts a state and ID to a proto state.
func MakeProtoState(id ID, st SerializedState) (*protos.State, error) {
	bytes, err := json.Marshal(st)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal to json-encoded state proto value")
	}
	p := &protos.State{Type: id.Type, DeviceID: id.DeviceID, Value: bytes, Version: st.Version}
	return p, nil
}

// MakeSerializedStatesByID converts proto states to native serialized states,
// keyed by state ID.
func MakeSerializedStatesByID(states []*protos.State) (SerializedStatesByID, error) {
	byID := SerializedStatesByID{}
	for _, p := range states {
		id := ID{Type: p.Type, DeviceID: p.DeviceID}
		st, err := MakeSerializedState(p)
		if err != nil {
			return nil, err
		}
		byID[id] = st
	}
	return byID, nil
}

// MakeSerializedState converts a proto state to a native serialized state.
func MakeSerializedState(p *protos.State) (SerializedState, error) {
	serialized := SerializedState{}
	err := json.Unmarshal(p.Value, &serialized)
	if err != nil {
		return SerializedState{}, errors.Wrap(err, "error unmarshaling json-encoded state proto value")
	}
	return serialized, nil
}

// MakeStatesByID converts proto states to native states keyed by state ID.
func MakeStatesByID(states []*protos.State, serdes serde.Registry) (StatesByID, error) {
	var errs []*protos.IDAndError
	byID := StatesByID{}
	for _, p := range states {
		id := ID{Type: p.Type, DeviceID: p.DeviceID}
		st, sErr, err := MakeState(p, serdes)
		if err != nil {
			return nil, err
		}
		if sErr != nil {
			errs = append(errs, sErr)
		} else {
			byID[id] = st
		}
	}

	if len(errs) != 0 {
		// Just log, since gateways reporting garbage shouldn't break orc8r
		glog.Errorf("Found malformed reported states: %v", errs)
	}

	return byID, nil
}

// MakeState converts a proto state to a native state.
func MakeState(p *protos.State, serdes serde.Registry) (State, *protos.IDAndError, error) {
	serialized, err := MakeSerializedState(p)
	if err != nil {
		return State{}, nil, err
	}

	st := State{
		ReporterID: serialized.ReporterID,
		TimeMs:     serialized.TimeMs,
		Version:    p.Version,
	}

	// Recover reported state within state struct
	st.ReportedState, err = serde.Deserialize(serialized.SerializedReportedState, p.Type, serdes)
	if err != nil {
		sErr := &protos.IDAndError{Type: p.Type, DeviceID: p.DeviceID, Error: err.Error()}
		return State{}, sErr, nil
	}

	if st.ReportedState == nil {
		err = errors.Errorf("state {type: %s, key: %s} should not have nil SerializedReportedState value", p.Type, p.DeviceID)
		sErr := &protos.IDAndError{Type: p.Type, DeviceID: p.DeviceID, Error: err.Error()}
		return State{}, sErr, nil
	}
	model, ok := st.ReportedState.(serde.ValidateableBinaryConvertible)
	if !ok {
		err = errors.Errorf("could not convert state {type: %s, key: %s} to validateable model", p.Type, p.DeviceID)
		sErr := &protos.IDAndError{Type: p.Type, DeviceID: p.DeviceID, Error: err.Error()}
		return State{}, sErr, nil
	}
	if err := model.ValidateModel(); err != nil {
		sErr := &protos.IDAndError{Type: p.Type, DeviceID: p.DeviceID, Error: err.Error()}
		return State{}, sErr, nil
	}

	return st, nil, nil
}
