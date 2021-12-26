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

// Package service outlines the Magma microservices framework.
// The framework helps to create a microservice easily, and provides
// the common service logic like service303, config, etc.
package service

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

const (
	PrintGrpcPayloadFlag = "print-grpc-payload"
	PrintGrpcPayloadEnv  = "MAGMA_PRINT_GRPC_PAYLOAD"
)

var (
	printGrpcPayload           int
	currentlyRunningServices   = make(map[string]*Service)
	currentlyRunningServicesMu sync.RWMutex
)

var defaultKeepaliveParams = keepalive.ServerParameters{
	MaxConnectionIdle:     2 * time.Minute,
	MaxConnectionAge:      10 * time.Minute,
	MaxConnectionAgeGrace: 5 * time.Minute,
	// first ping will be sent after 2*Time from server to client
	// Note nghttpx proxy has a backend-read-timeout defaults to 1m so we
	// want to send a ping within 1m to keep the connection alive
	// https://nghttp2.org/documentation/nghttpx.1.html#cmdoption-nghttpx--backend-read-timeout
	Time:    20 * time.Second,
	Timeout: 10 * time.Second,
}

func init() {
	flag.IntVar(&printGrpcPayload, PrintGrpcPayloadFlag, int(GRPCLOG_DISABLED),
		"Enable GRPC Payload Printout (0: disabled 1: enabled 2: hide verbose")
}

type Service struct {
	// Type identifies the service
	Type string

	// GrpcServer runs on the port specified in the registry.
	// Services can attach different servicers to the GrpcServer.
	GrpcServer *grpc.Server

	// Version of the service
	Version string

	// State of the service
	State protos.ServiceInfo_ServiceState

	// Health of the service
	Health protos.ServiceInfo_ApplicationHealth

	// Start time of the service
	StartTimeSecs uint64

	// Config of the service
	Config *config.Map
}

// NewServiceWithOptions returns a new GRPC orchestrator service implementing
// service303 with the specified grpc server options.
// It will not instantiate the service with the identity checking middleware.
func NewServiceWithOptions(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*Service, error) {
	return NewServiceWithOptionsImpl(moduleName, serviceName, serverOptions...)
}

// NewServiceWithOptionsImpl returns a new GRPC service implementing
// service303 with the specified grpc server options. This will not instantiate
// the service with the identity checking middleware.
func NewServiceWithOptionsImpl(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*Service, error) {
	// Load config, in case it does not exist, log
	configMap, err := config.GetServiceConfig(moduleName, serviceName)
	if err != nil {
		glog.Warningf("Failed to load config for service %s: %s", serviceName, err)
		configMap = nil
	}

	// Registers new logger in case print-grpc-payload flag or MAGMA_PRINT_GRPC_PAYLOAD env is set
	registerPrintGrpcPayloadLogCodecIfRequired()

	// Use keepalive options to proactively reinit http2 connections and
	// mitigate flow control issues
	opts := []grpc.ServerOption{grpc.KeepaliveParams(defaultKeepaliveParams)}
	opts = append(opts, serverOptions...) // keepalive is prepended so serverOptions can override if requested

	grpcServer := grpc.NewServer(opts...)
	service := &Service{
		Type:          serviceName,
		GrpcServer:    grpcServer,
		Version:       "0.0.0",
		State:         protos.ServiceInfo_STARTING,
		Health:        protos.ServiceInfo_APP_UNHEALTHY,
		StartTimeSecs: uint64(time.Now().Unix()),
		Config:        configMap,
	}
	protos.RegisterService303Server(service.GrpcServer, service)

	// Store into global for future access
	currentlyRunningServicesMu.Lock()
	currentlyRunningServices[serviceName] = service
	currentlyRunningServicesMu.Unlock()

	return service, nil
}

// Run the service. This function blocks until its interrupted
// by a signal or until the gRPC server is stopped.
func (service *Service) Run() error {
	port, err := registry.GetServicePort(service.Type)
	if err != nil {
		return fmt.Errorf("get service port: %v", err)
	}

	// Create the server socket for gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen on port %d: %v", port, err)
	}
	service.State = protos.ServiceInfo_ALIVE
	service.Health = protos.ServiceInfo_APP_HEALTHY
	return service.GrpcServer.Serve(lis)
}

// RunTest runs the test service on a given Listener. This function blocks
// by a signal or until the gRPC server is stopped.
func (service *Service) RunTest(lis net.Listener) {
	service.State = protos.ServiceInfo_ALIVE
	service.Health = protos.ServiceInfo_APP_HEALTHY
	err := service.GrpcServer.Serve(lis)
	if err != nil {
		glog.Fatal("Failed to run test service")
	}
}

// GetDefaultKeepaliveParameters returns the default keepalive server parameters.
func GetDefaultKeepaliveParameters() keepalive.ServerParameters {
	return defaultKeepaliveParams
}
