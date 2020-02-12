/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Magma microservices framework.
// The framework helps to create a microservice easily, and provides
// the common service logic like service303, config, etc.
package service

import (
	"flag"
	"fmt"
	"net"
	"time"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/service/middleware/unary"
	"magma/orc8r/cloud/go/util"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	grpc_proto "google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/keepalive"
)

const (
	PrintGrpcPayloadFlag = "print-grpc-payload"
	PrintGrpcPayloadEnv  = "MAGMA_PRINT_GRPC_PAYLOAD"
)

var printGrpcPayload bool
var currentlyRunningServices = make(map[string]*Service)

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
	flag.BoolVar(&printGrpcPayload, PrintGrpcPayloadFlag, false, "Enable GRPC Payload Printout")
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
	Config *config.ConfigMap
}

// NewOrchestratorService returns a new GRPC orchestrator service
// implementing service303. This service will implement a middleware
// interceptor to perform identity check. If your service does not or can not
// perform identity checks, (e.g. federation), use NewServiceWithOptions.
func NewOrchestratorService(moduleName string, serviceName string) (*Service, error) {
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})
	return NewServiceWithOptions(moduleName, serviceName, grpc.UnaryInterceptor(unary.MiddlewareHandler))
}

// NewOrchestratorServiceWithOptions returns a new GRPC orchestrator service
// implementing service303 with the specified grpc server options. This service
// will implement a middleware interceptor to perform identity check.
func NewOrchestratorServiceWithOptions(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*Service, error) {
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})
	serverOptions = append(serverOptions, grpc.UnaryInterceptor(unary.MiddlewareHandler))
	return NewServiceWithOptions(moduleName, serviceName, serverOptions...)
}

// NewServiceWithOptions returns a new GRPC orchestrator service implementing
// service303 with the specified grpc server options. This will not instantiate
// the service with the identity checking middleware.
//
// This function will also load all orchestrator plugins to populate registries
// which the service may use. Errors during plugin loading will result in a
// fatal. It also will load the config specified by [service name].yml. Since
// not all services have configs, it will only log in case it does not exist
func NewServiceWithOptions(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*Service, error) {
	// Parse the command line flags
	flag.Parse()

	// Load config, in case it does not exist, log
	configMap, err := config.GetServiceConfig(moduleName, serviceName)
	if err != nil {
		glog.Warningf("Failed to load config for service %s: %s", serviceName, err)
		configMap = nil
	}

	// Check if service was started with print-grpc-payload flag or MAGMA_PRINT_GRPC_PAYLOAD env is set
	if printGrpcPayload || util.IsTruthyEnv(PrintGrpcPayloadEnv) {
		ls := logCodec{encoding.GetCodec(grpc_proto.Name)}
		if ls.protoCodec != nil {
			glog.Errorf("Adding Debug Codec for service %s", serviceName)
			encoding.RegisterCodec(ls)
		}
	}

	// Use keepalive options to proactively reinit http2 connections and
	// mitigate flow control issues
	opts := []grpc.ServerOption{grpc.KeepaliveParams(defaultKeepaliveParams)}
	opts = append(opts, serverOptions...) // keepalive is prepended so serverOptions can override if requested

	grpcServer := grpc.NewServer(opts...)
	service := Service{
		Type:          serviceName,
		GrpcServer:    grpcServer,
		Version:       "0.0.0",
		State:         protos.ServiceInfo_STARTING,
		Health:        protos.ServiceInfo_APP_UNHEALTHY,
		StartTimeSecs: uint64(time.Now().Unix()),
		Config:        configMap,
	}
	protos.RegisterService303Server(service.GrpcServer, &service)

	// Store into global for future access
	currentlyRunningServices[serviceName] = &service

	return &service, nil
}

// Run the service. This function blocks until its interrupted
// by a signal or until the gRPC server is stopped.
func (service *Service) Run() error {
	port, err := registry.GetServicePort(service.Type)
	if err != nil {
		return fmt.Errorf("Failed to get service port: %s", err)
	}

	// Create the server socket for gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("Failed to listen on port %d: %s", port, err)
	}
	service.State = protos.ServiceInfo_ALIVE
	service.Health = protos.ServiceInfo_APP_HEALTHY
	return service.GrpcServer.Serve(lis)
}

// Run the test service on a given Listener. This function blocks
// by a signal or until the gRPC server is stopped.
func (service *Service) RunTest(lis net.Listener) error {
	service.State = protos.ServiceInfo_ALIVE
	service.Health = protos.ServiceInfo_APP_HEALTHY
	return service.GrpcServer.Serve(lis)
}

// GetDefaultKeepaliveParameters returns the default keepalive server parameters.
func GetDefaultKeepaliveParameters() keepalive.ServerParameters {
	return defaultKeepaliveParams
}
