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

package servicers_test

import (
	"errors"
	"testing"

	"magma/orc8r/cloud/go/services/streamer"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type mockStreamProvider struct {
	retVal []*protos.DataUpdate
	retErr error
}

func (m *mockStreamProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	return m.retVal, m.retErr
}

func TestStreamingServer_GetUpdates(t *testing.T) {
	streamer_test_init.StartTestService(t)
	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)
	grpcClient := protos.NewStreamerClient(conn)

	expected := []*protos.DataUpdate{
		{Key: "a", Value: []byte("123")},
		{Key: "b", Value: []byte("456")},
	}
	streamer_test_init.StartNewTestProvider(t, &mockStreamProvider{retVal: expected}, "mock1")

	streamerClient, err := grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hwId", StreamName: "mock1"},
	)
	assert.NoError(t, err)

	actual, err := streamerClient.Recv()
	assert.NoError(t, err)
	updates := actual.GetUpdates()
	assert.Equal(t, len(expected), len(updates))

	for i, u := range updates {
		assert.Equal(t, protos.TestMarshal(expected[i]), protos.TestMarshal(u))
	}

	// Error in provider
	streamer_test_init.StartNewTestProvider(t, &mockStreamProvider{retVal: nil, retErr: errors.New("MOCK")}, "mock2")
	streamerClient, err = grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hwId", StreamName: "mock2"},
	)
	assert.NoError(t, err)
	_, err = streamerClient.Recv()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MOCK")

	// Provider does not exist
	streamerClient, err = grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hwId", StreamName: "stream_dne"},
	)
	assert.NoError(t, err)
	_, err = streamerClient.Recv()
	assert.Error(t, err, "Stream stream_dne does not exist", codes.Unavailable)
}
