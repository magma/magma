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

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	hss "magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/service"
)

func StartTestService(t *testing.T) (*service.Service, error) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.MOCK_HSS)

	service := hss.NewTestHomeSubscriberServer(t)

	protos.RegisterHSSConfiguratorServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	return srv, nil
}
