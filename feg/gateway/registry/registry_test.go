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

package registry_test

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/gateway/service_registry"
	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

type helloServer struct{}

func (srv *helloServer) SayHello(
	ctx context.Context,
	req *protos.HelloRequest,
) (*protos.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("No metadata found for request")
	}
	authority, ok := md[":authority"]
	if !ok || len(authority) == 0 {
		return nil, fmt.Errorf("Authority not found in metadata")
	}
	return &protos.HelloReply{Greeting: authority[0]}, nil
}

func TestCloudConnection(t *testing.T) {
	// Update control_proxy Host to localhost (needed for docker)
	reg := service_registry.Get()
	reg.AddService(platform_registry.ServiceLocation{Name: registry.CONTROL_PROXY, Host: "127.0.0.1", Port: 50053})

	lis, err := net.Listen("tcp", "127.0.0.1:")
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	protos.RegisterHelloServer(grpcServer, &helloServer{})
	serverStarted := make(chan struct{}, 2)
	localPort, _ := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])

	go func() {
		t.Logf("Starting server on addr: %s, port: %d", lis.Addr(), localPort)
		serverStarted <- struct{}{}
		grpcServer.Serve(lis)
	}()
	<-serverStarted
	time.Sleep(time.Millisecond * 7)

	configMap := config.NewConfigMap(map[interface{}]interface{}{
		"local_port": localPort, "cloud_address": "controller.magma.test",
		"cloud_port": 443})

	conn, err := reg.GetCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	client := protos.NewHelloClient(conn)

	reply, err := client.SayHello(context.Background(), &protos.HelloRequest{Greeting: "hi"})
	assert.NoError(t, err)
	assert.Equal(t, "hello-controller.magma.test", reply.Greeting)

	assert.NoError(t, conn.Close())

	conn, err = reg.GetSharedCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	client = protos.NewHelloClient(conn)

	reply, err = client.SayHello(context.Background(), &protos.HelloRequest{Greeting: "hi"})
	assert.NoError(t, err)
	assert.Equal(t, "hello-controller.magma.test", reply.Greeting)

	conn1, err := reg.GetSharedCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	assert.Equal(t, conn, conn1)

	grpcServer.Stop() // kill previous server & all existing connections

	// Start new Hello service & make sure all tests pass without errors &
	// the cached cloud connection is automatically recreated
	lis, err = net.Listen("tcp", "127.0.0.1:")
	assert.NoError(t, err)
	localPort, _ = strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
	grpcServer = grpc.NewServer()
	protos.RegisterHelloServer(grpcServer, &helloServer{})
	go func() {
		t.Logf("Starting second server on addr: %s, port: %d", lis.Addr(), localPort)
		serverStarted <- struct{}{}
		grpcServer.Serve(lis)
	}()
	<-serverStarted
	time.Sleep(time.Millisecond * 7)

	configMap = config.NewConfigMap(map[interface{}]interface{}{
		"local_port": localPort, "cloud_address": "controller.magma.test",
		"cloud_port": 443})
	conn, err = reg.GetSharedCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	assert.NotEqual(t, conn, conn1)

	client = protos.NewHelloClient(conn)
	reply, err = client.SayHello(context.Background(), &protos.HelloRequest{Greeting: "hi"})
	assert.NoError(t, err)
	assert.Equal(t, "hello-controller.magma.test", reply.Greeting)

	conn1, err = reg.GetSharedCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	assert.Equal(t, conn, conn1)

	client = protos.NewHelloClient(conn1)
	reply, err = client.SayHello(context.Background(), &protos.HelloRequest{Greeting: "hi"})
	assert.NoError(t, err)
	assert.Equal(t, "hello-controller.magma.test", reply.Greeting)

	reg.CleanupSharedCloudConnection("hello")
	platform_registry.SetSharedCloudConnectionTTL(time.Millisecond * 100)

	conn, err = reg.GetSharedCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	assert.NotEqual(t, conn, conn1)

	time.Sleep(time.Millisecond * 101)
	conn1, err = reg.GetSharedCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	assert.NotEqual(t, conn1, conn)

	platform_registry.SetSharedCloudConnectionTTL(platform_registry.DefaultSharedCloudConnectionTTL)
}
