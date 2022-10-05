/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package unary_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service/middleware/unary/test/protos"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	state_test "magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/test_utils"
	lib_protos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

func TestInterceptorsHappyPath(t *testing.T) {
	servicer := newTestService(func(req string) (res string, err error) { return req, nil })
	srv, lis, _ := test_utils.NewTestService(t, orc8r.ModuleName, "test_service")
	protos.RegisterTestServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis, nil)

	c := getClient(t)
	ctx := registerGateway(t)

	res, err := c.Get(ctx, &protos.GetRequest{Req: "pass"})
	assert.NoError(t, err)
	assert.Equal(t, "pass", res.Res)
}

// TestInterceptorsBadIdentity ensures incorrect certificate info causes error.
//
// NOTE: there's another case we'd like to test as well: missing certificate
// info. However, that's hard to test right now because the access middleware
// always permits localhost connections when identity is missing.
func TestInterceptorsBadIdentity(t *testing.T) {
	servicer := newTestService(func(req string) (res string, err error) { return req, nil })
	srv, lis, _ := test_utils.NewTestService(t, orc8r.ModuleName, "test_service")
	protos.RegisterTestServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis, nil)

	c := getClient(t)
	ctx := state_test.GetContextWithCertificate(t, "bad_test_hwid")

	res, err := c.Get(ctx, &protos.GetRequest{Req: "pass"})
	assert.Error(t, err)
	assert.NotEqual(t, "pass", res.GetRes())
}

func TestInterceptorsUnregisteredGateway(t *testing.T) {
	servicer := newTestService(func(req string) (res string, err error) { return req, nil })
	srv, lis, _ := test_utils.NewTestService(t, orc8r.ModuleName, "test_service")
	protos.RegisterTestServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis, nil)

	c := getClient(t)
	ctx := registerGateway(t)
	removeGateway(t)

	res, err := c.Get(ctx, &protos.GetRequest{Req: "pass"})
	assert.Error(t, err)
	assert.NotEqual(t, "pass", res.GetRes())
}

func TestInterceptorHandlerPanic(t *testing.T) {
	servicer := newTestService(func(req string) (res string, err error) { panic("can we panic now!") })
	srv, lis, _ := test_utils.NewTestService(t, orc8r.ModuleName, "test_service")
	protos.RegisterTestServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis, nil)

	c := getClient(t)
	ctx := registerGateway(t)

	res, err := c.Get(ctx, &protos.GetRequest{Req: "pass"})
	assert.Error(t, err)
	assert.NotEqual(t, "pass", res.GetRes())
}

func getClient(t *testing.T) protos.TestServiceClient {
	conn, err := registry.GetConnection("test_service", lib_protos.ServiceType_SOUTHBOUND)
	assert.NoError(t, err)
	return protos.NewTestServiceClient(conn)
}

type tfunc func(req string) (res string, err error)

type testService struct {
	f tfunc
}

func newTestService(f tfunc) protos.TestServiceServer {
	return &testService{f: f}
}
func (t *testService) Get(ctx context.Context, req *protos.GetRequest) (*protos.GetResponse, error) {
	res, err := t.f(req.Req)
	return &protos.GetResponse{Res: res}, err
}

func registerGateway(t *testing.T) context.Context {
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)

	configurator_test.RegisterNetwork(t, "test_network", "Test Network")
	configurator_test.RegisterGateway(t, "test_network", "test_gw", &models.GatewayDevice{HardwareID: "test_hwid"})
	ctx := state_test.GetContextWithCertificate(t, "test_hwid")
	return ctx
}

func removeGateway(t *testing.T) {
	configurator_test.RemoveGateway(t, "test_network", "test_gw")
}
