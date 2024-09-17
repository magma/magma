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

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_protos "magma/orc8r/cloud/go/services/directoryd/protos"
	servicers "magma/orc8r/cloud/go/services/directoryd/servicers/protected"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

func StartTestService(t *testing.T) {
	// Create service
	srv, lis, plis := test_utils.NewTestService(t, orc8r.ModuleName, directoryd.ServiceName)

	// Init storage
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLStoreFactory(storage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	store := storage.NewDirectorydBlobstore(fact)

	// Add servicers
	directoryServicer, err := servicers.NewDirectoryLookupServicer(store)
	assert.NoError(t, err)
	directoryd_protos.RegisterDirectoryLookupServer(srv.ProtectedGrpcServer, directoryServicer)
	protos.RegisterGatewayDirectoryServiceServer(srv.GrpcServer, servicers.NewDirectoryUpdateServicer())

	// Run service
	go srv.RunTest(lis, plis)
}
