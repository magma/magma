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
	"time"

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/state"
	state_config "magma/orc8r/cloud/go/services/state/config"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	"magma/orc8r/cloud/go/services/state/metrics"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	"magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/service/config"
)

// how often to report gateway status
const gatewayStatusReportInterval = time.Second * 60

const nonPostgresDriverMessage = `Configuration warning:

This deployment has automatic state reindexing enabled, but is targeting a
database driver other than Postgres. This will cause the state service
to log a (harmless) DB syntax error, due to its use of Postgres-specific
syntax for automatic reindexing.

(Option 1) Continue using non-Postgres driver. To clear this warning, update
the state.yml cloud config to set enable_automatic_reindexing to false.
Keep in mind that, for this option, you will have to perform manual state
reindexing on every Orc8r upgrade. We provide a CLI to manage this, and will
provide directions in the upgrade notes.

(Option 2) Switch to a Postgres driver.
`

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, "singleton_state")
	if err != nil {
		glog.Fatalf("Error creating singleton_state service %v", err)
	}

	singletonReindex := srv.Config.MustGetBool(state_config.EnableSingletonReindex)
	if singletonReindex {
		db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
		if err != nil {
			glog.Fatalf("Error connecting to database: %v", err)
		}
		store := blobstore.NewSQLStoreFactory(state.DBTableName, db, sqorc.GetSqlBuilder())
		err = store.InitializeFactory()
		if err != nil {
			glog.Fatalf("Error initializing state database: %v", err)
		}

		stateServicer := newStateServicer(store)
		protos.RegisterStateServiceServer(srv.GrpcServer, stateServicer)
		glog.Info("srv.Config %s", srv.Config)

		indexerManagerServer := newSingletonIndexerManagerServicer(srv.Config, db, store)

		indexer_protos.RegisterIndexerManagerServer(srv.GrpcServer, indexerManagerServer)
	}

	go metrics.PeriodicallyReportGatewayStatus(gatewayStatusReportInterval)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running state service: %v", err)
	}
}

func newStateServicer(store blobstore.StoreFactory) protos.StateServiceServer {
	servicer, err := servicers.NewStateServicer(store)
	if err != nil {
		glog.Fatalf("Error creating state servicer: %v", err)
	}
	return servicer
}

func newSingletonIndexerManagerServicer(cfg *config.Map, db *sql.DB, store blobstore.StoreFactory) indexer_protos.IndexerManagerServer {
	glog.Info("newSingletonIndexerManagerServicer")

	versioner := reindex.NewVersioner(db, sqorc.GetSqlBuilder())
	err := versioner.Initialize()
	if err != nil {
		glog.Fatal("Error initializing state reindex queue")
	}

	autoReindex := cfg.MustGetBool(state_config.EnableAutomaticReindexing)
	reindexer := reindex.NewReindexerSingleton(reindex.NewStore(store), versioner)
	servicer := servicers.NewIndexerManagerServicer(reindexer, autoReindex)

	if autoReindex {
		glog.Info("Automatic reindexing enabled for state service")
		go reindexer.Run(context.Background())
	} else {
		glog.Info("Automatic reindexing disabled for state service")
	}

	return servicer
}
