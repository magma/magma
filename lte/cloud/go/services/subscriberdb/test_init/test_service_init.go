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
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	state_protos "magma/orc8r/cloud/go/services/state/protos"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"

	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
)

const (
	// testMaxProtosLoadSize is the maxProtosLoadSize used in a test
	// subscriberdb service.
	testMaxProtosLoadSize = 10
)

func StartTestService(t *testing.T) {
	// Create service
	labels := map[string]string{
		orc8r.StateIndexerLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StateIndexerVersionAnnotation: "1",
		orc8r.StateIndexerTypesAnnotation:   lte.MobilitydStateType,
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, subscriberdb.ServiceName, labels, annotations)

	// Init storage
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory(subscriberdb.LookupTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	ipStore := storage.NewIPLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, ipStore.Initialize())
	digestStore := storage.NewDigestStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, digestStore.Initialize())
	perSubDigestFact := blobstore.NewSQLBlobStorageFactory(subscriberdb.PerSubDigestTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, perSubDigestFact.InitializeFactory())
	perSubDigestStore := storage.NewPerSubDigestStore(perSubDigestFact)
	subStore := storage.NewSubStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, subStore.Initialize())

	// Load service configs
	var serviceConfig subscriberdb.Config
	config.MustGetStructuredServiceConfig(lte.ModuleName, subscriberdb.ServiceName, &serviceConfig)
	glog.Infof("Subscriberdb service config %+v", serviceConfig)
	serviceConfig.MaxProtosLoadSize = testMaxProtosLoadSize

	// Add servicers
	protos.RegisterSubscriberLookupServer(srv.GrpcServer, servicers.NewLookupServicer(fact, ipStore))
	state_protos.RegisterIndexerServer(srv.GrpcServer, servicers.NewIndexerServicer())
	lte_protos.RegisterSubscriberDBCloudServer(srv.GrpcServer, servicers.NewSubscriberdbServicer(serviceConfig, digestStore, perSubDigestStore, subStore))

	// Run service
	go srv.RunTest(lis)
}
