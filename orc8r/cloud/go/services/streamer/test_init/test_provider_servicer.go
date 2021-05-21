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
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

type providerServicer struct {
	provider providers.StreamProvider
}

// StartNewTestProvider starts a new stream provider service which forwards
// calls to the passed provider.
func StartNewTestProvider(t *testing.T, provider providers.StreamProvider, streamName string) {
	labels := map[string]string{
		orc8r.StreamProviderLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StreamProviderStreamsAnnotation: streamName,
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, streamName, labels, annotations)
	servicer := &providerServicer{provider: provider}
	streamer_protos.RegisterStreamProviderServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}

func (p *providerServicer) GetUpdates(ctx context.Context, req *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	updates, err := p.provider.GetUpdates(req.GatewayId, req.ExtraArgs)
	res := &protos.DataUpdateBatch{Updates: updates}
	return res, err
}
