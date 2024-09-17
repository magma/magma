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

package session

import (
	"fbc/cwf/radius/monitoring"
	"fmt"
	"sync"

	"go.opencensus.io/tag"
)

type memoryStorage struct {
	data sync.Map
}

func (m *memoryStorage) Get(sessionID string) (*State, error) {
	counter := ReadSessionState.Start(
		tag.Upsert(monitoring.SessionIDTag, sessionID),
		tag.Upsert(monitoring.StorageTag, "memory"),
	)
	data, ok := m.data.Load(sessionID)
	if !ok {
		counter.Failure("not_found")
		return nil, fmt.Errorf("session %s no found in storage", sessionID)
	}

	shapedData, ok := data.(State)
	if !ok {
		counter.Failure("corrupted")
		return nil, ErrInvalidDataFormat
	}

	counter.Success()
	return &shapedData, nil
}

func (m *memoryStorage) Set(sessionID string, state State) error {
	counter := WriteSessionState.Start(
		tag.Upsert(monitoring.SessionIDTag, sessionID),
		tag.Upsert(monitoring.StorageTag, "memory"),
	)
	m.data.Store(sessionID, state)
	counter.Success()
	return nil
}

func (m *memoryStorage) Reset(sessionID string) error {
	counter := ResetSessionState.Start(
		tag.Upsert(monitoring.SessionIDTag, sessionID),
		tag.Upsert(monitoring.StorageTag, "memory"),
	)
	m.data.Delete(sessionID)
	counter.Success()
	return nil
}

// NewMultiSessionMemoryStorage Returns a new memory-stored session state storage
func NewMultiSessionMemoryStorage() GlobalStorage {
	return &memoryStorage{
		data: sync.Map{},
	}
}
