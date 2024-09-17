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

package servicers

import (
	"magma/orc8r/cloud/go/services/dispatcher/broker"
	"magma/orc8r/lib/go/protos"
)

// A little Go "polymorphism" magic for testing
type testSyncRPCServer struct {
	SyncRPCService
}

const TestSyncRPCAgHwId = "Test-AGW-Hw-Id"

func (srv *testSyncRPCServer) EstablishSyncRPCStream(stream protos.SyncRPCService_EstablishSyncRPCStreamServer) error {
	// See if there is an Identity in the CTX and if not, use default TestSyncRPCAgHwId
	gw := protos.GetClientGateway(stream.Context())
	if gw == nil {
		return srv.serveGwId(stream, TestSyncRPCAgHwId)
	}
	return srv.SyncRPCService.EstablishSyncRPCStream(stream)
}

func NewTestSyncRPCServer(hostName string, broker broker.GatewayRPCBroker) (*testSyncRPCServer, error) {
	return &testSyncRPCServer{SyncRPCService{hostName, broker}}, nil
}
