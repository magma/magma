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

package servicers_test

import (
	"context"
	"encoding/json"
	"testing"

	"magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/lte/protos"
	"magma/lte/cloud/go/services/lte/servicers"
	"magma/lte/cloud/go/services/lte/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
)

func TestLookupServicer_EnodebState(t *testing.T) {
	ctx := context.Background()
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	store := storage.NewEnodebStateLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, err)
	err = store.Initialize()
	assert.NoError(t, err)
	l := servicers.NewLookupServicer(store)

	t.Run("basic", func(t *testing.T) {
		_, err := l.GetEnodebState(ctx, &protos.GetEnodebStateRequest{
			NetworkId: "nid0",
			GatewayId: "g1",
			EnodebSn:  "123",
		})
		assert.Error(t, err)

		enbState := models.NewDefaultEnodebStatus()
		serializedState, err := json.Marshal(enbState)
		assert.NoError(t, err)
		_, err = l.SetEnodebState(ctx, &protos.SetEnodebStateRequest{
			NetworkId:       "nid0",
			GatewayId:       "g1",
			EnodebSn:        "123",
			SerializedState: serializedState,
		})
		assert.NoError(t, err)

		got, err := l.GetEnodebState(ctx, &protos.GetEnodebStateRequest{
			NetworkId: "nid0",
			GatewayId: "g1",
			EnodebSn:  "123",
		})
		assert.NoError(t, err)
		assert.Equal(t, got.GetSerializedState(), serializedState)
	})
}
