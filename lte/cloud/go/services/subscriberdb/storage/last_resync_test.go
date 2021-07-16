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
	"testing"
	"time"

	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestLastResyncTimeStore(t *testing.T) {
	fact := test_utils.NewSQLBlobstore(t, subscriberdb.LastResyncTimeTableBlobstore)
	assert.NoError(t, fact.InitializeFactory())
	s := storage.NewLastResyncTimeStore(fact)

	t.Run("initially empty", func(t *testing.T) {
		lastResyncTime, err := s.Get("n0", "g0")
		assert.NoError(t, err)
		assert.Empty(t, lastResyncTime)
	})

	t.Run("set and get", func(t *testing.T) {
		expectedLastResyncTime := uint64(time.Now().Unix())
		err := s.Set("n0", "g0", expectedLastResyncTime)
		assert.NoError(t, err)
		lastResyncTime, err := s.Get("n0", "g0")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime, lastResyncTime)
	})

	t.Run("multiple set and get", func(t *testing.T) {
		expectedLastResyncTime := uint64(time.Now().Unix())

		err := s.Set("n0", "g0", expectedLastResyncTime+1)
		assert.NoError(t, err)
		err = s.Set("n0", "g1", expectedLastResyncTime+2)
		assert.NoError(t, err)
		err = s.Set("n1", "g0", expectedLastResyncTime+3)
		assert.NoError(t, err)
		err = s.Set("n1", "g1", expectedLastResyncTime+4)
		assert.NoError(t, err)

		lastResyncTime1, err := s.Get("n0", "g0")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+1, lastResyncTime1)
		lastResyncTime2, err := s.Get("n0", "g1")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+2, lastResyncTime2)
		lastResyncTime3, err := s.Get("n1", "g0")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+3, lastResyncTime3)
		lastResyncTime4, err := s.Get("n1", "g1")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+4, lastResyncTime4)

		// Test upserting value to the store
		err = s.Set("n0", "g0", expectedLastResyncTime+5)
		assert.NoError(t, err)
		lastResyncTime5, err := s.Get("n0", "g0")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+5, lastResyncTime5)
	})
}
