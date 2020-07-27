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
	"net"
	"testing"

	"magma/orc8r/cloud/go/services/dispatcher/broker/mocks"
	"magma/orc8r/cloud/go/services/dispatcher/httpserver"
)

func StartTestHttpServer(t *testing.T) (net.Addr, *mocks.GatewayRPCBroker) {
	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("net.Listen err: %v\n", err)
	}

	broker := new(mocks.GatewayRPCBroker)
	server := httpserver.NewSyncRPCHttpServer(broker)
	go server.Serve(lis)
	return lis.Addr(), broker
}
