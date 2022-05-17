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
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	protected_servicers "magma/orc8r/cloud/go/services/state/servicers/protected"
	servicers "magma/orc8r/cloud/go/services/state/servicers/southbound"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

const (
	singleAttempt = 1
)

// StartTestService instantiates a service backed by an in-memory storage.
func StartTestService(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	startService(t, db)
}

// StartTestServiceInternal instantiates a test DB-backed service, returning
// the derived reindexer and job queue for internal usage.
// Supported drivers include: postgres.
func StartTestServiceInternal(t *testing.T, dbName, dbDriver string) (reindex.Reindexer, reindex.JobQueue) {
	db := sqorc.OpenCleanForTest(t, dbName, dbDriver)
	return startService(t, db)
}

func startService(t *testing.T, db *sql.DB) (reindex.Reindexer, reindex.JobQueue) {
	srv, lis, plis := test_utils.NewTestService(t, orc8r.ModuleName, state.ServiceName)

	factory := blobstore.NewSQLStoreFactory(state.DBTableName, db, sqorc.GetSqlBuilder())
	require.NoError(t, factory.InitializeFactory())

	stateServicer, err := servicers.NewStateServicer(factory)
	require.NoError(t, err)
	protos.RegisterStateServiceServer(srv.GrpcServer, stateServicer)

	cloudStateServicer, err := protected_servicers.NewCloudStateServicer(factory)
	require.NoError(t, err)
	protos.RegisterCloudStateServiceServer(srv.ProtectedGrpcServer, cloudStateServicer)

	queue := reindex.NewSQLJobQueue(singleAttempt, db, sqorc.GetSqlBuilder())
	require.NoError(t, queue.Initialize())
	reindexer := reindex.NewReindexerQueue(queue, reindex.NewStore(factory))
	indexerServicer := protected_servicers.NewIndexerManagerServicer(reindexer, false)
	indexer_protos.RegisterIndexerManagerServer(srv.ProtectedGrpcServer, indexerServicer)

	go srv.RunTest(lis, plis)
	return reindexer, queue
}

// StartTestSingletonServiceInternal instantiates a test DB-backed singleton service, returning
// the derived reindexer and job queue for internal usage.
// Supported drivers include: postgres.
func StartTestSingletonServiceInternal(t *testing.T) reindex.Reindexer {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	return startSingletonService(t, db)
}

func startSingletonService(t *testing.T, db *sql.DB) reindex.Reindexer {
	srv, lis, plis := test_utils.NewTestService(t, orc8r.ModuleName, state.ServiceName)

	factory := blobstore.NewSQLStoreFactory(state.DBTableName, db, sqorc.GetSqlBuilder())
	require.NoError(t, factory.InitializeFactory())

	stateServicer, err := servicers.NewStateServicer(factory)
	require.NoError(t, err)
	protos.RegisterStateServiceServer(srv.GrpcServer, stateServicer)

	cloudStateServicer, err := protected_servicers.NewCloudStateServicer(factory)
	require.NoError(t, err)
	protos.RegisterCloudStateServiceServer(srv.ProtectedGrpcServer, cloudStateServicer)

	versioner := reindex.NewVersioner(db, sqorc.GetSqlBuilder())
	versioner.Initialize()
	reindexer := reindex.NewReindexerSingleton(reindex.NewStore(factory), versioner)
	indexerServicer := protected_servicers.NewIndexerManagerServicer(reindexer, false)
	indexer_protos.RegisterIndexerManagerServer(srv.ProtectedGrpcServer, indexerServicer)

	go srv.RunTest(lis, plis)
	return reindexer
}
