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

	"magma/orc8r/cloud/go/orc8r"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/services/orchestrator"
	"magma/orc8r/cloud/go/services/orchestrator/servicers"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	streamer_servicers "magma/orc8r/cloud/go/services/streamer/servicers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
)

type testStreamerServer struct {
	protos.StreamerServer
}

func (srv *testStreamerServer) GetUpdates(req *protos.StreamRequest, stream protos.Streamer_GetUpdatesServer) error {
	return streamer_servicers.GetUpdatesUnverified(req, stream)
}

func StartTestService(t *testing.T) {
	StartTestServiceInternal(t, servicers.NewBuilderServicer(), servicers.NewIndexerServicer(), servicers.NewProviderServicer())
}

func StartTestServiceInternal(
	t *testing.T,
	builder builder_protos.MconfigBuilderServer,
	indexer indexer_protos.IndexerServer,
	provider streamer_protos.StreamProviderServer,
) {
	labels := map[string]string{}
	annotations := map[string]string{}

	if builder != nil {
		labels[orc8r.MconfigBuilderLabel] = "true"
	}
	if provider != nil {
		labels[orc8r.StreamProviderLabel] = "true"
		annotations[orc8r.StreamProviderStreamsAnnotation] = definitions.MconfigStreamName
	}

	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, orchestrator.ServiceName, labels, annotations)
	protos.RegisterStreamerServer(srv.GrpcServer, &testStreamerServer{})

	if builder != nil {
		builder_protos.RegisterMconfigBuilderServer(srv.GrpcServer, builder)
	}
	if indexer != nil {
		indexer_protos.RegisterIndexerServer(srv.GrpcServer, indexer)
	}
	if provider != nil {
		streamer_protos.RegisterStreamProviderServer(srv.GrpcServer, provider)
	}

	go srv.RunTest(lis)
}
