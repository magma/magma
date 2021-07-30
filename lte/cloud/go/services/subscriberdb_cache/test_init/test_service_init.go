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

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	labels, annotations := map[string]string{}, map[string]string{}
	srv, lis := test_utils.NewTestOrchestratorService(
		t, lte.ModuleName, subscriberdb_cache.ServiceName, labels, annotations,
	)

	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory(subscriberdb.SyncstoreBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	store := syncstore.NewSyncStore(db, sqorc.GetSqlBuilder(), fact)
	assert.NoError(t, store.Initialize())

	serviceConfig := subscriberdb_cache.MustGetServiceConfig()
	glog.Infof("Subscriberdb_cache service config %+v", serviceConfig)

	go subscriberdb_cache.MonitorDigests(serviceConfig, store)
	srv.RunTest(lis)
}
