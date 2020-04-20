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

package test_init

import (
	"fmt"
	"testing"

	"magma/fbinternal/cloud/go/services/testcontroller"
	"magma/fbinternal/cloud/go/services/testcontroller/protos"
	"magma/fbinternal/cloud/go/services/testcontroller/servicers"
	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/definitions"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	StartTestServiceWithDB(t, "testcontroller__test__service__db")
}

func StartTestServiceWithDB(t *testing.T, dbName string) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, testcontroller.ServiceName)

	// Connect to postgres_test
	db := sqorc.OpenCleanForTest(t, dbName, sqorc.PostgresDriver)

	nodeStore := storage.NewSQLNodeLeasorStorage(db, &mockIDGenerator{}, sqorc.GetSqlBuilder())
	err := nodeStore.Init()
	assert.NoError(t, err)
	nodes := servicers.NewNodeLeasorServicer(nodeStore)
	protos.RegisterNodeLeasorServer(srv.GrpcServer, nodes)

	testStore := storage.NewSQLTestcontrollerStorage(db, sqorc.GetSqlBuilder())
	err = testStore.Init()
	assert.NoError(t, err)
	tests := servicers.NewTestControllerServicer(testStore)
	protos.RegisterTestControllerServer(srv.GrpcServer, tests)

	go func() {
		defer db.Close()
		srv.RunTest(lis)
	}()
}

func GetTestTestcontrollerStorage(t *testing.T) storage.TestControllerStorage {
	db, err := sqorc.Open("postgres", definitions.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_test user=magma_test password=magma_test host=postgres_test sslmode=disable"))
	if err != nil {
		t.Fatalf("could not dial potgres_test DB %s", err)
	}
	return storage.NewSQLTestcontrollerStorage(db, sqorc.GetSqlBuilder())
}

type mockIDGenerator struct {
	current uint64
}

func (m *mockIDGenerator) New() string {
	m.current++
	return fmt.Sprintf("%d", m.current)
}
