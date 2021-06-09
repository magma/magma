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

package state

import (
	"context"

	"magma/orc8r/cloud/go/serde"
	state_types "magma/orc8r/cloud/go/services/state/types"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

// GetStateClient returns a client to the state service.
func GetStateClient() (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewStateServiceClient(conn), nil
}

// GetState returns the state specified by the networkID, typeVal, and hwID.
func GetState(networkID string, typ string, hwID string, serdes serde.Registry) (state_types.State, error) {
	id := state_types.ID{Type: typ, DeviceID: hwID}

	states, err := GetStates(networkID, state_types.IDs{id}, serdes)
	if err != nil {
		return state_types.State{}, err
	}

	st, ok := states[id]
	if !ok {
		return state_types.State{}, merrors.ErrNotFound
	}
	return st, nil
}

// GetStates returns a map of states specified by the networkID and a list of
// type and key.
func GetStates(networkID string, stateIDs state_types.IDs, serdes serde.Registry) (state_types.StatesByID, error) {
	if len(stateIDs) == 0 {
		return state_types.StatesByID{}, nil
	}

	client, err := GetStateClient()
	if err != nil {
		return nil, err
	}

	res, err := client.GetStates(
		context.Background(), &protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       makeProtoIDs(stateIDs),
		},
	)
	if err != nil {
		return nil, err
	}
	return state_types.MakeStatesByID(res.States, serdes)
}

// SearchStates returns all states matching the filter arguments.
// typeFilter and keyFilter are both OR clauses, and the final predicate
// applied to the search will be the AND of both filters.
// If keyPrefix is defined (non-nil and non-empty), it will take precedence
// the keyFilter argument.
// e.g.: ["t1", "t2"], ["k1", "k2"] => (t1 OR t2) AND (k1 OR k2)
func SearchStates(networkID string, typeFilter []string, keyFilter []string, keyPrefix *string, serdes serde.Registry) (state_types.StatesByID, error) {
	client, err := GetStateClient()
	if err != nil {
		return nil, err
	}

	req := &protos.GetStatesRequest{
		NetworkID:  networkID,
		TypeFilter: typeFilter,
		IdFilter:   keyFilter,
		LoadValues: true,
	}
	if !funk.IsEmpty(keyPrefix) {
		req.IdPrefix = *keyPrefix
		req.IdFilter = nil
	}
	res, err := client.GetStates(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return state_types.MakeStatesByID(res.States, serdes)
}

// DeleteStates deletes states specified by the networkID and a list of
// type and key.
func DeleteStates(networkID string, stateIDs state_types.IDs) error {
	client, err := GetStateClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteStates(
		context.Background(),
		&protos.DeleteStatesRequest{
			NetworkID: networkID,
			Ids:       makeProtoIDs(stateIDs),
		},
	)
	return err
}

// GetSerializedStates returns a map of states specified by the networkID and
// a list of type and key.
func GetSerializedStates(networkID string, stateIDs state_types.IDs) (state_types.SerializedStatesByID, error) {
	if len(stateIDs) == 0 {
		return state_types.SerializedStatesByID{}, nil
	}

	client, err := GetStateClient()
	if err != nil {
		return nil, err
	}

	res, err := client.GetStates(
		context.Background(), &protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       makeProtoIDs(stateIDs),
		},
	)
	if err != nil {
		return nil, err
	}
	return state_types.MakeSerializedStatesByID(res.States)
}

func makeProtoIDs(stateIDs state_types.IDs) []*protos.StateID {
	var ids []*protos.StateID
	for _, st := range stateIDs {
		ids = append(ids, &protos.StateID{Type: st.Type, DeviceID: st.DeviceID})
	}
	return ids
}
