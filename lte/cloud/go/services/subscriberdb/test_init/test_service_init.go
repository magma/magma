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
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	lookup_servicers "magma/lte/cloud/go/services/subscriberdb/servicers/protected"
	subscriberdbcloud_servicer "magma/lte/cloud/go/services/subscriberdb/servicers/southbound"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	state_protos "magma/orc8r/cloud/go/services/state/protos"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) storage.SubscriberStorage {
	// Create service
	labels := map[string]string{
		orc8r.StateIndexerLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StateIndexerVersionAnnotation: "1",
		orc8r.StateIndexerTypesAnnotation:   lte.MobilitydStateType + "," + lte.GatewaySubscriberStateType,
	}
	srv, lis, plis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, subscriberdb.ServiceName, labels, annotations)

	// Init storage
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLStoreFactory(subscriberdb.LookupTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	ipStore := storage.NewIPLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, ipStore.Initialize())
	syncstoreFact := blobstore.NewSQLStoreFactory(subscriberdb.SyncstoreTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, syncstoreFact.InitializeFactory())
	subscriberStore, err := syncstore.NewSyncStoreReader(db, sqorc.GetSqlBuilder(), syncstoreFact, syncstore.Config{TableNamePrefix: subscriberdb.SyncstoreTableNamePrefix})
	assert.NoError(t, err)
	assert.NoError(t, subscriberStore.Initialize())
	subscriberStateStore := storage.NewSubscriberStorage(db, sqorc.GetSqlBuilder())
	assert.NoError(t, subscriberStateStore.Initialize())

	// Sane default service configs
	serviceConfig := subscriberdb.Config{
		DigestsEnabled:         true,
		ChangesetSizeThreshold: 500,
		MaxProtosLoadSize:      10,
		ResyncIntervalSecs:     86400,
	}

	// Add servicers
	protos.RegisterSubscriberLookupServer(srv.ProtectedGrpcServer, lookup_servicers.NewLookupServicer(fact, ipStore))
	state_protos.RegisterIndexerServer(srv.ProtectedGrpcServer, lookup_servicers.NewIndexerServicer(subscriberStateStore))
	lte_protos.RegisterSubscriberDBCloudServer(srv.GrpcServer, subscriberdbcloud_servicer.NewSubscriberdbServicer(serviceConfig, subscriberStore))

	// Run service
	go srv.RunTest(lis, plis)

	return subscriberStateStore
}
