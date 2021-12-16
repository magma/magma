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
	"time"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers/registration"
	"magma/orc8r/cloud/go/test_utils"
	protos2 "magma/orc8r/lib/go/protos"
)

// StartTestService instantiates a service backed by an in-memory storage
// DOES NOT start a bootstrapper servicer
func StartTestService(t *testing.T) {
	factory := test_utils.NewSQLBlobstore(t, "bootstrapper_test_service_blobstore")
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, bootstrapper.ServiceName)
	store := registration.NewBlobstoreStore(factory)
	cloudRegistrationServicer, err := registration.NewCloudRegistrationServicer(store, "rootCA", 30*time.Minute, false)
	assert.NoError(t, err)
	registrationServicer := registration.NewRegistrationServicer()

	protos2.RegisterCloudRegistrationServer(srv.GrpcServer, cloudRegistrationServicer)
	protos2.RegisterRegistrationServer(srv.GrpcServer, registrationServicer)
	go srv.RunTest(lis)
}
