/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/directoryd/servicers"
	magmaerrors "magma/orc8r/lib/go/errors"
)

type MockStore struct {
	previous_seen string
	Hits          int
}

func NewMockStore() MockStore {
	return MockStore{previous_seen: "0", Hits: 0}
}

func (m *MockStore) getIDOnDatabaseWithExistingMatch(networkId string, id string) (string, error) {
	m.Hits += 1
	if m.previous_seen == "0" {
		m.previous_seen = id
		return "hw_id_1", nil
	}
	return "", magmaerrors.ErrNotFound
}

func (m *MockStore) getIDOnDatabaseNotExistingMatch(networkId string, id string) (string, error) {
	m.Hits += 1
	m.previous_seen = id
	return "", magmaerrors.ErrNotFound
}

func (m *MockStore) getIDOnDatabaseAlwaysError(networkId string, id string) (string, error) {
	m.Hits += 1
	m.previous_seen = id
	return "", fmt.Errorf("unknkwn error")
}

func TestGetUniqueUint32Id_NoExistingId(t *testing.T) {
	gen := servicers.NewIdGenerator()
	mockStore := NewMockStore()
	id, err := gen.GetUniqueUint32Id("any", mockStore.getIDOnDatabaseNotExistingMatch)
	assert.NoError(t, err)
	assert.NotEqual(t, uint32(0), id)
	assert.Equal(t, 1, mockStore.Hits)
}

func TestGetUniqueUint32Id_AlreadyExistingId(t *testing.T) {
	gen := servicers.NewIdGenerator()
	mockStore := NewMockStore()
	id, err := gen.GetUniqueUint32Id("any", mockStore.getIDOnDatabaseWithExistingMatch)
	assert.NoError(t, err)
	assert.NotEqual(t, uint32(0), id)
	assert.Equal(t, 2, mockStore.Hits)
}

func TestGetUniqueUint32Id_AlwaysErrors(t *testing.T) {
	gen := servicers.NewIdGeneratorWithAttempts(3)
	mockStore := NewMockStore()
	id, err := gen.GetUniqueUint32Id("any", mockStore.getIDOnDatabaseAlwaysError)
	assert.Error(t, err)
	assert.Equal(t, uint32(0), id)
	assert.Equal(t, 3, mockStore.Hits)
}
