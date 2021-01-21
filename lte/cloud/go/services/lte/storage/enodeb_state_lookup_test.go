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

package storage_test

import (
	"encoding/json"
	"testing"

	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/lte/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestEnodebStateLookup(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	s := storage.NewEnodebStateLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, s.Initialize())

	t.Run("empty initially", func(t *testing.T) {
		_, err := s.GetEnodebState("n0", "g1", "10101")
		assert.Error(t, err)
	})
	enbState := lte_models.NewDefaultEnodebStatus()
	enbState2 := lte_models.NewDefaultEnodebStatus()
	enbState2.MmeConnected = swag.Bool(false)
	enbState2.EnodebConnected = swag.Bool(false)
	serializedState1, err := json.Marshal(enbState)
	assert.NoError(t, err)
	serializedState2, err := json.Marshal(enbState2)
	assert.NoError(t, err)

	t.Run("basic insert", func(t *testing.T) {
		err := s.SetEnodebState("n0", "g1", "123", serializedState1)
		assert.NoError(t, err)

		got, err := s.GetEnodebState("n0", "g1", "123")
		assert.NoError(t, err)
		assert.Equal(t, serializedState1, got)

		err = s.SetEnodebState("n0", "g1", "555", serializedState2)
		assert.NoError(t, err)

		got, err = s.GetEnodebState("n0", "g1", "555")
		assert.NoError(t, err)
		assert.Equal(t, serializedState2, got)
	})

	t.Run("upsert", func(t *testing.T) {
		enbState.EnodebConnected = swag.Bool(false)
		serializedState1, err = json.Marshal(enbState)
		assert.NoError(t, err)
		err := s.SetEnodebState("n0", "g1", "123", serializedState1)
		assert.NoError(t, err)

		got, err := s.GetEnodebState("n0", "g1", "123")
		assert.NoError(t, err)
		assert.Equal(t, serializedState1, got)

		enbState2.IPAddress = "123.123.123.123"
		serializedState2, err = json.Marshal(enbState2)
		assert.NoError(t, err)

		err = s.SetEnodebState("n0", "g1", "555", serializedState2)
		assert.NoError(t, err)

		got, err = s.GetEnodebState("n0", "g1", "555")
		assert.NoError(t, err)
		assert.Equal(t, serializedState2, got)
	})

	t.Run("second gateway", func(t *testing.T) {
		enbState.EnodebConnected = swag.Bool(true)
		serializedState1, err = json.Marshal(enbState)
		assert.NoError(t, err)
		err := s.SetEnodebState("n0", "g2", "123", serializedState1)
		assert.NoError(t, err)

		got, err := s.GetEnodebState("n0", "g2", "123")
		assert.NoError(t, err)
		assert.Equal(t, serializedState1, got)

		enbState2.RfTxOn = swag.Bool(false)
		serializedState2, err = json.Marshal(enbState2)
		assert.NoError(t, err)
		err = s.SetEnodebState("n0", "g2", "555", serializedState2)
		assert.NoError(t, err)

		got, err = s.GetEnodebState("n0", "g2", "555")
		assert.NoError(t, err)
		assert.Equal(t, serializedState2, got)
	})
}
