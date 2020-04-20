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

package authstate

import (
	"errors"
	"fbc/cwf/radius/modules/eap/packet"
	"fbc/cwf/radius/monitoring"
	"fmt"
	"sync"

	"go.opencensus.io/tag"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
)

// Manager an interface for EAP state management storage
type memoryManager struct {
	getOpCounter   monitoring.Operation
	setOpCounter   monitoring.Operation
	resetOpCounter monitoring.Operation
	storage        sync.Map
}

// Set sets a value in the state manager for the given auth request and eap packet type
//
// state should be restorable with any (radius.Packet, packet.EAPType) which are
// received on the same auth session, that is - storing must use value which will be
// sent for each EAP packet
func (m *memoryManager) Set(authReq *radius.Packet, eaptype packet.EAPType, state Container) error {
	m.setOpCounter.Start()
	m.storage.Store(getKey(authReq), state)
	m.setOpCounter.Success()
	return nil
}

// Get gets a value from the state manager for the given auth request and eap packet type
func (m *memoryManager) Get(authReq *radius.Packet, eaptype packet.EAPType) (*Container, error) {
	m.getOpCounter.Start()
	state, ok := m.storage.Load(getKey(authReq))
	if !ok {
		m.getOpCounter.Failure("not_found")
		return nil, errors.New("eap state not found")
	}

	result, ok := state.(Container)
	if !ok {
		m.getOpCounter.Failure("deserialize_failed")
		return nil, errors.New("eap state failed to deserialize")
	}

	m.getOpCounter.Success()
	return &result, nil
}

// Reset resets the value stored in auth state manager for the given auth request and eap packet type
func (m *memoryManager) Reset(authReq *radius.Packet, eapType packet.EAPType) error {
	m.resetOpCounter.Start()
	m.storage.Delete(getKey(authReq))
	m.resetOpCounter.Success()
	return nil
}

// getKey Composes a storage key to store and access the EAP state under
// at this point, a state is stored at device level (that is a unique device
// may only engage in one auth flow at any given moment through this system)
func getKey(r *radius.Packet) string {
	clientID := r.Get(rfc2865.CallingStationID_Type)
	nasID := r.Get(rfc2865.CalledStationID_Type)
	return fmt.Sprintf("eap__%s__%s", string(clientID), string(nasID))
}

// NewMemoryManager Create a new EAP Auth State Manager which uses
// local memory for (transient) storage
func NewMemoryManager() Manager {
	return &memoryManager{
		getOpCounter:   monitoring.NewOperation("eap_state_get", tag.Upsert(monitoring.StorageTag, "memory")),
		setOpCounter:   monitoring.NewOperation("eap_state_set", tag.Upsert(monitoring.StorageTag, "memory")),
		resetOpCounter: monitoring.NewOperation("eap_state_reset", tag.Upsert(monitoring.StorageTag, "memory")),
		storage:        sync.Map{},
	}
}
