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
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/lib/go/service/config"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, subscriberdb_cache.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %+v", err)
	}

	var serviceConfig subscriberdb_cache.Config
	config.MustGetStructuredServiceConfig(lte.ModuleName, subscriberdb_cache.ServiceName, &serviceConfig)
	if err := serviceConfig.Validate(); err != nil {
		glog.Fatalf("Invalid subscriberdb_cache service configs: %+v", err)
	}
	glog.Infof("Subscriberdb_cache service config %+v", serviceConfig)

	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %+v", err)
	}
	fact := blobstore.NewSQLStoreFactory(subscriberdb.SyncstoreTableBlobstore, db, sqorc.GetSqlBuilder())
	if err := fact.InitializeFactory(); err != nil {
		glog.Fatalf("Error initializing blobstore storage for subscriber syncstore: %+v", err)
	}
	// Garbage collection interval for syncstore cache writers is enforced to be half the time for the service worker's
	// update interval, to prevent cache writers from outliving update cycles
	store, err := syncstore.NewSyncStore(db, sqorc.GetSqlBuilder(), fact, syncstore.Config{
		TableNamePrefix:              subscriberdb.SyncstoreTableNamePrefix,
		CacheWriterValidIntervalSecs: int64(serviceConfig.SleepIntervalSecs / 2),
	})
	if err != nil {
		glog.Fatalf("Error creating new subscriber syncstore: %+v", err)
	}
	if err := store.Initialize(); err != nil {
		glog.Fatalf("Error initializing subscriber syncstore")
	}

	go subscriberdb_cache.MonitorDigests(serviceConfig, store)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %+v", err)
	}
}
