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

	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	// Create service
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, subscriberdb.ServiceName)

	// Init storage
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory(subscriberdb.LookupTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()

	// Add servicers
	servicer := servicers.NewLookupServicer(fact)
	assert.NoError(t, err)
	protos.RegisterSubscriberLookupServer(srv.GrpcServer, servicer)

	// Run service
	go srv.RunTest(lis)
}
