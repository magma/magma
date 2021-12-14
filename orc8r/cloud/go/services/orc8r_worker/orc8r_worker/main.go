/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"database/sql"

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/orc8r_worker"
	"magma/orc8r/cloud/go/services/state"
	state_config "magma/orc8r/cloud/go/services/state/config"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	"magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/service/config"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, orc8r_worker.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating orc8r_worker service %v", err)
	}

	singletonReindex := srv.Config.MustGetBool(state_config.EnableSingletonReindex)
	if singletonReindex {
		startSingletonReindexer(srv)
	}

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running orc8r_worker service: %v", err)
	}
}

func startSingletonReindexer(srv *service.OrchestratorService) {
	glog.Info("Running singleton reindexer")

	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error connecting to database: %v", err)
	}
	store := blobstore.NewSQLStoreFactory(state.DBTableName, db, sqorc.GetSqlBuilder())
	err = store.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing state database: %v", err)
	}

	indexerManagerServer := newSingletonIndexerManagerServicer(srv.Config, db, store)
	indexer_protos.RegisterIndexerManagerServer(srv.GrpcServer, indexerManagerServer)
}

func newSingletonIndexerManagerServicer(cfg *config.Map, db *sql.DB, store blobstore.StoreFactory) indexer_protos.IndexerManagerServer {
	versioner := reindex.NewVersioner(db, sqorc.GetSqlBuilder())
	err := versioner.Initialize()
	if err != nil {
		glog.Fatal("Error initializing orc8r_worker reindex versioner")
	}

	autoReindex := cfg.MustGetBool(state_config.EnableAutomaticReindexing)
	reindexer := reindex.NewReindexerSingleton(reindex.NewStore(store), versioner)
	servicer := servicers.NewIndexerManagerServicer(reindexer, autoReindex)

	if autoReindex {
		glog.Info("Automatic reindexing enabled for orc8r_worker service")
		go reindexer.Run(context.Background())
	} else {
		glog.Info("Automatic reindexing disabled for orc8r_worker service")
	}

	return servicer
}
