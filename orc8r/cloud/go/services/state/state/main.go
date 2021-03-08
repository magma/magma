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

	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
)

// how often to report gateway status
const gatewayStatusReportInterval = time.Second * 60

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f6a54d1a20134c258b1e0b227d4d0982@o529355.ingest.sentry.io/5667116",
	})
	if err != nil {
		glog.Fatalf("sentry.Init: %s", err)
	}
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, state.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating state service %v", err)
	}
	db, err := sqorc.Open(storage.SQLDriver, storage.DatabaseSource)
	if err != nil {
		glog.Fatalf("Error connecting to database: %v", err)
	}
	store := blobstore.NewEntStorage(state.DBTableName, db, sqorc.GetSqlBuilder())
	err = store.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing state database: %v", err)
	}

	stateServicer := newStateServicer(store)
	protos.RegisterStateServiceServer(srv.GrpcServer, stateServicer)
	indexerManagerServer := newIndexerManagerServicer(srv.Config, db, store)
	indexer_protos.RegisterIndexerManagerServer(srv.GrpcServer, indexerManagerServer)

	go metrics.PeriodicallyReportGatewayStatus(gatewayStatusReportInterval)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running state service: %v", err)
	}
}

func newStateServicer(store blobstore.BlobStorageFactory) protos.StateServiceServer {
	servicer, err := servicers.NewStateServicer(store)
	if err != nil {
		glog.Fatalf("Error creating state servicer: %v", err)
	}
	return servicer
}

func newIndexerManagerServicer(cfg *config.ConfigMap, db *sql.DB, store blobstore.BlobStorageFactory) indexer_protos.IndexerManagerServer {
	queue := reindex.NewSQLJobQueue(reindex.DefaultMaxAttempts, db, sqorc.GetSqlBuilder())
	err := queue.Initialize()
	if err != nil {
		glog.Fatal("Error initializing state reindex queue")
	}
	_, err = queue.PopulateJobs()
	if err != nil {
		glog.Fatalf("Unexpected error initializing reindex job queue: %s", err)
	}

	autoReindex := cfg.MustGetBool(state_config.EnableAutomaticReindexing)
	reindexer := reindex.NewReindexer(queue, reindex.NewStore(store))
	servicer := servicers.NewIndexerManagerServicer(reindexer, autoReindex)

	if autoReindex {
		glog.Info("Automatic reindexing enabled for state service")
		go reindexer.Run(context.Background())
	} else {
		glog.Info("Automatic reindexing disabled for state service")
	}

	return servicer
}
