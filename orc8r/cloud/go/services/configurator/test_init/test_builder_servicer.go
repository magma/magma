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
	"context"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/golang/protobuf/ptypes/any"
)

type builderServicer struct {
	builder mconfig.Builder
}

// StartNewTestBuilder starts a new mconfig builder service which forwards
// calls to the passed builder.
func StartNewTestBuilder(t *testing.T, builder mconfig.Builder) {
	labels := map[string]string{
		orc8r.MconfigBuilderLabel: "true",
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_mconfig_builder_service", labels, nil)
	servicer := &builderServicer{builder: builder}
	builder_protos.RegisterMconfigBuilderServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}

func (s *builderServicer) Build(ctx context.Context, request *builder_protos.BuildRequest) (ret *builder_protos.BuildResponse, err error) {
	ret = &builder_protos.BuildResponse{ConfigsByKey: map[string]*any.Any{}, JsonConfigsByKey: map[string][]byte{}}

	// TODO(T71525030): revert defer, above changes, and fn signature changes
	defer func() { err = ret.FillJSONConfigs(err) }()

	ret.ConfigsByKey, err = s.builder.Build(request.Network, request.Graph, request.GatewayId)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
