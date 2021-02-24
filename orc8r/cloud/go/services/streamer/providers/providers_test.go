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

package providers_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/configurator/mconfig/mocks"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/servicers"
	orchestrator_test_init "magma/orc8r/cloud/go/services/orchestrator/test_init"
	"magma/orc8r/cloud/go/services/streamer"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/cloud/go/services/streamer/test_utils/mconfig/test_protos"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestMconfigStreamer_Configurator(t *testing.T) {
	configurator_test_init.StartTestService(t)
	streamer_test_init.StartTestService(t)
	orchestrator_test_init.StartTestServiceInternal(t, nil, nil, servicers.NewProviderServicer())

	vals := map[string]proto.Message{
		"some_service": &test_protos.Message1{Field: "hello"},
	}
	marshaledConfigs, err := mconfig.MarshalConfigs(vals)
	assert.NoError(t, err)

	mockBuilder := &mocks.Builder{}
	mockBuilder.On("Build", mock.Anything, mock.Anything, "gw1").Return(marshaledConfigs, nil)
	configurator_test_init.StartNewTestBuilder(t, mockBuilder)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1", PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)

	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)
	streamerClient := protos.NewStreamerClient(conn)

	// TODO(#2310): add back test for matching digest => empty data updates

	t.Run("normal stream update", func(t *testing.T) {
		stream, err := streamerClient.GetUpdates(context.Background(), &protos.StreamRequest{GatewayId: "hw1", StreamName: "configs"})
		assert.NoError(t, err)

		expectedProtos := map[string]proto.Message{
			"some_service": &test_protos.Message1{Field: "hello"},
		}
		expected := make(map[string]*any.Any, len(expectedProtos))
		for k, v := range expectedProtos {
			anyV, err := ptypes.MarshalAny(v)
			assert.NoError(t, err)
			expected[k] = anyV
		}

		actualMarshaled, err := stream.Recv()
		assert.NoError(t, err)
		actual := &protos.GatewayConfigs{}
		err = protos.Unmarshal(actualMarshaled.Updates[0].Value, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual.ConfigsByKey)
	})
}
