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
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/cloud/go/services/tenants/servicers"
	"magma/orc8r/cloud/go/services/tenants/servicers/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

// StartTestService instantiates a service backed by an in-memory storage
func StartTestService(t *testing.T) {
	factory := test_utils.NewSQLBlobstore(t, "device_test_service_blobstore")
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, tenants.ServiceName)
	store := storage.NewBlobstoreStore(factory)
	server, err := servicers.NewTenantsServicer(store)
	assert.NoError(t, err)
	protos.RegisterTenantsServiceServer(srv.GrpcServer, server)
	go srv.RunTest(lis)
}
