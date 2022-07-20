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

package main

import (
	"github.com/golang/glog"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	lookup_servicers "magma/lte/cloud/go/services/subscriberdb/servicers/protected"
	subscriberdbcloud_servicer "magma/lte/cloud/go/services/subscriberdb/servicers/southbound"
	subscriberdb_storage "magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/obsidian"
	swagger_protos "magma/orc8r/cloud/go/services/obsidian/swagger/protos"
	swagger_servicers "magma/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
	state_protos "magma/orc8r/cloud/go/services/state/protos"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/lib/go/service/config"
)

func main() {
	// Create service
	srv, err := service.NewOrchestratorService(lte.ModuleName, subscriberdb.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %+v", err)
	}

	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %+v", err)
	}
	fact := blobstore.NewSQLStoreFactory(subscriberdb.LookupTableBlobstore, db, sqorc.GetSqlBuilder())
	if err := fact.InitializeFactory(); err != nil {
		glog.Fatalf("Error initializing MSISDN lookup storage: %+v", err)
	}
	ipStore := subscriberdb_storage.NewIPLookup(db, sqorc.GetSqlBuilder())
	if err := ipStore.Initialize(); err != nil {
		glog.Fatalf("Error initializing IP lookup storage: %+v", err)
	}

	syncstoreFact := blobstore.NewSQLStoreFactory(subscriberdb.SyncstoreTableBlobstore, db, sqorc.GetSqlBuilder())
	if err := syncstoreFact.InitializeFactory(); err != nil {
		glog.Fatalf("Error initializing blobstore storage for subscriber syncstore: %+v", err)
	}
	subscriberStore, err := syncstore.NewSyncStoreReader(db, sqorc.GetSqlBuilder(), syncstoreFact, syncstore.Config{TableNamePrefix: subscriberdb.SyncstoreTableNamePrefix})
	if err != nil {
		glog.Fatalf("Error creating new subscriber synsctore reader: %+v", err)
	}
	if err := subscriberStore.Initialize(); err != nil {
		glog.Fatalf("Error initializing subscriber syncstore: %+v", err)
	}

	subscriberStateStore := subscriberdb_storage.NewSubscriberStorage(db, sqorc.GetSqlBuilder())
	if err := subscriberStateStore.Initialize(); err != nil {
		glog.Fatalf("Error initializing subscriber state storage : %+v", err)
	}

	var serviceConfig subscriberdb.Config
	config.MustGetStructuredServiceConfig(lte.ModuleName, subscriberdb.ServiceName, &serviceConfig)
	glog.Infof("Subscriberdb service config %+v", serviceConfig)

	// Attach handlers
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers(subscriberStateStore))
	protos.RegisterSubscriberLookupServer(srv.ProtectedGrpcServer, lookup_servicers.NewLookupServicer(fact, ipStore))
	state_protos.RegisterIndexerServer(srv.ProtectedGrpcServer, lookup_servicers.NewIndexerServicer(subscriberStateStore))
	lte_protos.RegisterSubscriberDBCloudServer(srv.GrpcServer, subscriberdbcloud_servicer.NewSubscriberdbServicer(serviceConfig, subscriberStore))

	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, swagger_servicers.NewSpecServicerFromFile(subscriberdb.ServiceName))

	// Run service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %+v", err)
	}

}
