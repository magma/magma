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

package streamer_test

import (
	"fmt"
	"testing"
	"time"

	streamer_client "magma/gateway/streamer"
	"magma/orc8r/cloud/go/services/streamer"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	testStreamName = "mock1"
)

// Mock Cloud Streamer
type mockStreamProvider struct {
	retVal []*protos.DataUpdate
	extra  *any.Any
	retErr error
}

func (m *mockStreamProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	m.extra = extraArgs
	return m.retVal, m.retErr
}

var expected = []*protos.DataUpdate{
	{Key: "a", Value: []byte("123")},
	{Key: "b", Value: []byte("456")},
}

// Mock Client Streamer Listener
type testListener struct {
	err       chan error
	updateErr chan error
}

func (l testListener) GetName() string {
	return testStreamName
}

func (l testListener) GetExtraArgs() *any.Any {
	extra, _ := ptypes.MarshalAny(expected[0])
	return extra
}

func (l testListener) ReportError(e error) error {
	l.err <- e
	return nil // continue listener
}

func (l testListener) Update(ub *protos.DataUpdateBatch) bool {
	if len(expected) != len(ub.GetUpdates()) {
		l.updateErr <- fmt.Errorf("updates # %d != expected # %d", len(ub.GetUpdates()), len(expected))
		return false
	}
	for i, u := range ub.GetUpdates() {
		if protos.TestMarshal(expected[i]) != protos.TestMarshal(u) {
			l.updateErr <- fmt.Errorf(
				"update %s != expected %s", protos.TestMarshal(u), protos.TestMarshal(expected[i]))
			return false
		}
	}
	l.updateErr <- nil
	return true
}

// Mock GW Cloud Service registry
type mockedCloudRegistry struct {
	*platform_registry.ServiceRegistry
}

func (cr mockedCloudRegistry) GetCloudConnection(service string) (*grpc.ClientConn, error) {
	if service != definitions.StreamerServiceName {
		return nil, fmt.Errorf("not Implemented")
	}
	return platform_registry.GetConnection(streamer.ServiceName)
}

func (cr mockedCloudRegistry) GetCloudConnectionFromServiceConfig(serviceConfig *config.ConfigMap, service string) (*grpc.ClientConn, error) {
	return nil, fmt.Errorf("not Implemented")

}

// Test
func TestStreamerClient(t *testing.T) {
	streamer_test_init.StartTestService(t)

	streamerClient := streamer_client.NewStreamerClient(mockedCloudRegistry{})
	mockProvider := &mockStreamProvider{retVal: expected}
	streamer_test_init.StartNewTestProvider(t, mockProvider, testStreamName)

	l := testListener{}
	l.err = make(chan error)
	l.updateErr = make(chan error)
	assert.NoError(t, streamerClient.AddListener(l))
	go streamerClient.Stream(l)

	select {
	case e := <-l.err:
		assert.NoError(t, e)
	case e := <-l.updateErr:
		assert.NoError(t, e)
		var extra protos.DataUpdate
		err := ptypes.UnmarshalAny(mockProvider.extra, &extra)
		assert.NoError(t, err)
		assert.Equal(t, protos.TestMarshal(expected[0]), protos.TestMarshal(&extra))
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Test Timeout")
	}
}
