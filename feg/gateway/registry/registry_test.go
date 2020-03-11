/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry_test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	platform_registry.AddService(platform_registry.ServiceLocation{Name: registry.CONTROL_PROXY, Host: "127.0.0.1", Port: 50053})

	lis, err := net.Listen("tcp", ":44444")
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	protos.RegisterHelloServer(grpcServer, &helloServer{})
	serverStarted := make(chan struct{})
	go func() {
		log.Printf("Starting server")
		serverStarted <- struct{}{}
		grpcServer.Serve(lis)
	}()
	<-serverStarted
	time.Sleep(time.Millisecond)

	configMap := config.NewConfigMap(map[interface{}]interface{}{
		"local_port": 44444, "cloud_address": "controller.magma.test",
		"cloud_port": 443})

	conn, err := registry.NewCloudRegistry().GetCloudConnectionFromServiceConfig(configMap, "hello")
	assert.NoError(t, err)
	client := protos.NewHelloClient(conn)

	reply, err := client.SayHello(context.Background(), &protos.HelloRequest{Greeting: "hi"})
	assert.NoError(t, err)
	assert.Equal(t, "hello-controller.magma.test", reply.Greeting)
}
