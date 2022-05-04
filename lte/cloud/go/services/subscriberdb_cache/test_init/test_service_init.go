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

package test_init

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	labels, annotations := map[string]string{}, map[string]string{}
	srv, lis, plis := test_utils.NewTestOrchestratorService(
		t, lte.ModuleName, subscriberdb_cache.ServiceName, labels, annotations,
	)

	serviceConfig := subscriberdb_cache.Config{
		UpdateIntervalSecs: 300,
		SleepIntervalSecs:  120,
	}

	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	fact := blobstore.NewSQLStoreFactory(subscriberdb.SyncstoreTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	store, err := syncstore.NewSyncStore(db, sqorc.GetSqlBuilder(), fact, syncstore.Config{
		TableNamePrefix:              subscriberdb.SyncstoreTableNamePrefix,
		CacheWriterValidIntervalSecs: int64(serviceConfig.SleepIntervalSecs / 2),
	})
	assert.NoError(t, err)
	assert.NoError(t, store.Initialize())

	go subscriberdb_cache.MonitorDigests(serviceConfig, store)
	srv.RunTest(lis, plis)
}
