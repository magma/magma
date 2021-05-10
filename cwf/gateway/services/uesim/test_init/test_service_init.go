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

	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	factory := test_utils.NewSQLBlobstore(t, "uesim_test_service_blobstore")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.UeSim)
	server, err := servicers.NewUESimServer(factory)
	assert.NoError(t, err)
	protos.RegisterUESimServer(srv.GrpcServer, server)
	go srv.RunTest(lis)
}
