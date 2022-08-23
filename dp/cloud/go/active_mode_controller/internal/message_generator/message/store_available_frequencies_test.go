/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package message_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/message"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestStoreAvailableFrequenciesMessageString(t *testing.T) {
	m := message.NewStoreAvailableFrequenciesMessage(id, freqs)
	msg := "store available frequencies: %d (1110, 1100, 1100, 1000)"
	expected := fmt.Sprintf(msg, id)
	assert.Equal(t, expected, m.String())
}

func TestStoreAvailableFrequenciesMessageSend(t *testing.T) {
	client := &stubStoreClient{}
	provider := &stubStoreClientProvider{client: client}

	m := message.NewStoreAvailableFrequenciesMessage(id, freqs)
	require.NoError(t, m.Send(context.Background(), provider))

	expected := &active_mode.StoreAvailableFrequenciesRequest{
		Id:                   id,
		AvailableFrequencies: freqs,
	}
	assert.Equal(t, expected, client.req)
}

var freqs = []uint32{0b1110, 0b1100, 0b1100, 0b1000}

type stubStoreClientProvider struct {
	message.ClientProvider
	client *stubStoreClient
}

func (s *stubStoreClientProvider) GetActiveModeClient() active_mode.ActiveModeControllerClient {
	return s.client
}

type stubStoreClient struct {
	active_mode.ActiveModeControllerClient
	req *active_mode.StoreAvailableFrequenciesRequest
}

func (s *stubStoreClient) StoreAvailableFrequencies(_ context.Context, in *active_mode.StoreAvailableFrequenciesRequest, _ ...grpc.CallOption) (*empty.Empty, error) {
	s.req = in
	return &empty.Empty{}, nil
}
