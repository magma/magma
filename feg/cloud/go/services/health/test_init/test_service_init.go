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

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) (*servicers.TestHealthServer, error) {
	srv, lis := test_utils.NewTestService(t, feg.ModuleName, health.ServiceName)
	factory := test_utils.NewSQLBlobstore(t, health.DBTableName)
	servicer, err := servicers.NewTestHealthServer(factory)
	assert.NoError(t, err)
	protos.RegisterHealthServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
	return servicer, nil
}
