/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mock_feg_relay

import (
	"context"
	"fmt"
	"testing"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/orc8r/cloud/go/test_utils"
	orc8r_protos "magma/orc8r/lib/go/protos"
)

type TestS8ProxyResponderServer struct {
	protos.S8ProxyResponderServer
	ReceivedCreateBearerRequest *protos.CreateBearerRequestPgw
	ReceivedDeleteBearerRequest *protos.DeleteBearerRequestPgw
	ListAddr                    string
	Ready                       chan struct{}
}

func NewTestS8ProxyResponderServer() *TestS8ProxyResponderServer {
	return &TestS8ProxyResponderServer{
		Ready: make(chan struct{}),
	}
}

func (ts *TestS8ProxyResponderServer) CreateBearer(
	_ context.Context,
	cbReq *protos.CreateBearerRequestPgw) (*orc8r_protos.Void, error) {
	defer func() {
		// comunicate through the channel that we are done processing the call
		ts.Ready <- struct{}{}
	}()
	ts.ReceivedCreateBearerRequest = cbReq
	if cbReq == nil || cbReq.BearerContext == nil || cbReq.CAgwTeid == 0 {
		return nil, fmt.Errorf("mock feg_relay Create Bearer Request missing Bearer Contexct or TEID")
	}
	return &orc8r_protos.Void{}, nil
}

func (ts *TestS8ProxyResponderServer) DeleteBearerRequest(
	_ context.Context,
	dbReq *protos.DeleteBearerRequestPgw) (*orc8r_protos.Void, error) {
	defer func() {
		// comunicate through the channel that we are done processing the call
		ts.Ready <- struct{}{}
	}()
	ts.ReceivedDeleteBearerRequest = dbReq
	if dbReq == nil || dbReq.CAgwTeid == 0 {
		return nil, fmt.Errorf("mock feg_relay Delete Bearer Request missing Bearer Contexct or TEID")
	}
	return &orc8r_protos.Void{}, nil
}

// StartFegRelayTestService starts a grpc test service
func StartFegRelayTestService(t *testing.T) (*TestS8ProxyResponderServer, string) {
	labels := map[string]string{}
	annotations := map[string]string{}
	srv, lis, tempDir := test_utils.NewTestOrchestratorServiceWithControlProxy(
		t, feg.ModuleName, feg_relay.ServiceName, labels, annotations)
	// responder mocks feg relay service
	testResponderSrv := NewTestS8ProxyResponderServer()
	testResponderSrv.ListAddr = lis.Addr().String()
	go srv.RunTest(lis, nil)
	protos.RegisterS8ProxyResponderServer(srv.GrpcServer, testResponderSrv)
	fmt.Printf("Starting Mock Feg Relay service at %s", lis.Addr().String())
	// Remember to delete tempDir with defer os.RemoveAll(dir) once test is done
	return testResponderSrv, tempDir
}
