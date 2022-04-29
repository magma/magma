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

package service_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

func init() {
	// _ = flag.Set("logtostderr", "true") // uncomment to view logs during test

	_ = flag.Set("run_echo_server", "true")
}

func TestServiceRun(t *testing.T) {
	testStartTime := time.Now().Unix()
	allowedStartRange := 15.0
	serviceName := state.ServiceName

	// Create the service
	srv, lis, _ := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, serviceName, nil, nil)
	assert.Equal(t, protos.ServiceInfo_STARTING, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_UNHEALTHY, srv.Health)
	assert.NotNil(t, srv.EchoServer)

	// Start the service
	go srv.RunTest(lis, nil)

	// Wait for the service to be started and check its state and health
	time.Sleep(time.Second)
	assert.Equal(t, protos.ServiceInfo_ALIVE, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_HEALTHY, srv.Health)

	// Create a rpc stub and query the Service303 interface
	conn, err := registry.GetConnection(serviceName, protos.ServiceType_SOUTHBOUND)
	assert.NoError(t, err, "err in getting connection to service")
	client := protos.NewService303Client(conn)

	actualServiceInfo, err := client.GetServiceInfo(context.Background(), new(protos.Void))
	assert.NoError(t, err)

	expectedServiceInfo := &protos.ServiceInfo{
		Name:          "state",
		Version:       "0.0.0",
		State:         protos.ServiceInfo_ALIVE,
		Health:        protos.ServiceInfo_APP_HEALTHY,
		StartTimeSecs: actualServiceInfo.StartTimeSecs,
	}
	assert.NoError(t, err, "err in getting service info after srv started")
	test_utils.AssertMessagesEqual(t, expectedServiceInfo, actualServiceInfo)
	assert.InDelta(t, testStartTime, actualServiceInfo.StartTimeSecs, allowedStartRange)

	// check StopService rpc call.
	// this will have a connection error, which is expected.
	client.StopService(context.Background(), &protos.Void{})

	assert.Equal(t, protos.ServiceInfo_STOPPING, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_UNHEALTHY, srv.Health)
}

func TestRun_WithProtectedGrpcServer(t *testing.T) {
	//testStartTime := time.Now().Unix()
	//allowedStartRange := 15.0
	serviceName := state.ServiceName

	// Create the service
	srv, lis, plis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, serviceName, nil, nil, protos.ServiceType_PROTECTED)
	assert.Equal(t, protos.ServiceInfo_STARTING, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_UNHEALTHY, srv.Health)
	assert.NotNil(t, srv.EchoServer)

	// Start the service
	go srv.RunTest(lis, plis)

	// Wait for the service to be started and check its state and health
	time.Sleep(time.Second)
	assert.Equal(t, protos.ServiceInfo_ALIVE, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_HEALTHY, srv.Health)

	// Create a rpc stub and query the Service303 interface
	_, err := registry.GetConnection(serviceName, protos.ServiceType_SOUTHBOUND)
	assert.NoError(t, err, "err in getting connection to service")

	// Create a protected rpc stub to make sure that our client can connect to service
	_, err = registry.GetConnection(serviceName, protos.ServiceType_PROTECTED)
	assert.NoError(t, err, "err in getting protected connection to service")
}
