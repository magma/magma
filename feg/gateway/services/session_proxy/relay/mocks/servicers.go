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

package mocks

import (
	"log"
	"net"
	"testing"
	"time"

	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// Create and start a mock SessionProxyResponder as a service. This returns the
// mock object that is registered as the servicer to configure in test code and
// an instance of a `MockCloudRegistry` bound to the server address allocated
// to the mock local session manager.
func StartMockSessionProxyResponder(t *testing.T) (*SessionProxyResponderServer, *MockCloudRegistry) {
	lis, err := net.Listen("tcp", "")
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	sm := &SessionProxyResponderServer{}
	protos.RegisterSessionProxyResponderServer(grpcServer, sm)
	serverStarted := make(chan struct{})
	go func() {
		log.Printf("Starting server")
		serverStarted <- struct{}{}
		grpcServer.Serve(lis)
	}()
	<-serverStarted
	time.Sleep(time.Millisecond)

	return sm, &MockCloudRegistry{ServerAddr: lis.Addr().String()}
}
