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

	"magma/cwf/cloud/go/services/cwf"
	builder_servicers "magma/cwf/cloud/go/services/cwf/servicers/protected"
	"magma/orc8r/cloud/go/orc8r"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	StartTestServiceInternal(t, builder_servicers.NewBuilderServicer())
}

func StartTestServiceInternal(t *testing.T, builder builder_protos.MconfigBuilderServer) {
	labels := map[string]string{}
	annotations := map[string]string{}

	if builder != nil {
		labels[orc8r.MconfigBuilderLabel] = "true"
	}

	srv, lis, plis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, cwf.ServiceName, labels, annotations)

	if builder != nil {
		builder_protos.RegisterMconfigBuilderServer(srv.ProtectedGrpcServer, builder)
	}

	go srv.RunTest(lis, plis)
}
