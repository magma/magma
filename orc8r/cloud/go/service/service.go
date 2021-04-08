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

// Package service outlines the Magma microservices framework in the cloud.
// The framework helps to create a microservice easily, and provides
// the common service logic like service303, config, etc.
package service

import (
	"context"
	"flag"
	"fmt"
	"net"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service/middleware/unary"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
	platform_service "magma/orc8r/lib/go/service"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	RunEchoServerFlag = "run_echo_server"
)

var (
	runEchoServer bool
)

func init() {
	flag.BoolVar(&runEchoServer, RunEchoServerFlag, false, "Run echo HTTP server with service")
}

// OrchestratorService defines a service which extends the generic platform
// service with an optional HTTP server.
type OrchestratorService struct {
	*platform_service.Service

	// EchoServer runs on the echo_port specified in the registry.
	// This field will be nil for services that don't specify the
	// 'run_echo_server' flag.
	EchoServer *echo.Echo
}

// NewOrchestratorService returns a new gRPC orchestrator service
// implementing service303. If configured, it will also initialize an HTTP echo
// server as a part of the service. This service will implement a middleware
// interceptor to perform identity check. If your service does not or can not
// perform identity checks, (e.g., federation), use NewServiceWithOptions.
func NewOrchestratorService(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*OrchestratorService, error) {
	flag.Parse()

	err := registry.PopulateServices()
	if err != nil {
		return nil, err
	}

	sharedConfig, err := getSharedConfig()
	if err != nil {
		return nil, err
	}
	maxGRPCMsgSize := sharedConfig.MaxGRPCMessageSizeMB * 1024 * 1024

	// Set max gRPC message size to receive when acting as the client
	opts := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxGRPCMsgSize))
	registry.SetDialOpts(opts)
	// Set max gRPC message size to receive when acting as the server
	serverOptions = append(serverOptions, grpc.MaxRecvMsgSize(maxGRPCMsgSize))

	serverOptions = append(serverOptions, grpc.UnaryInterceptor(unary.MiddlewareHandler))
	platformService, err := platform_service.NewServiceWithOptionsImpl(moduleName, serviceName, serverOptions...)
	if err != nil {
		return nil, err
	}

	echoSrv, err := getEchoServerForOrchestratorService(serviceName)
	if err != nil {
		return nil, err
	}

	return &OrchestratorService{Service: platformService, EchoServer: echoSrv}, nil
}

// Run runs the service. If the echo HTTP server is non-nil, both the HTTP
// server and gRPC server are run, blocking until an error occurs or a server
// stopped. If the HTTP server is nil, only the gRPC server is run, blocking
// until its interrupted by a signal or until the gRPC server is stopped.
func (s *OrchestratorService) Run() error {
	if s.EchoServer == nil {
		return s.Service.Run()
	}
	serverErr := make(chan error, 1)
	go func() {
		err := s.Service.Run()
		shutdownErr := s.EchoServer.Shutdown(context.Background())
		if shutdownErr != nil {
			glog.Errorf("Error shutting down echo server: %v", shutdownErr)
		}
		serverErr <- err
	}()
	go func() {
		err := s.EchoServer.StartServer(s.EchoServer.Server)
		_, shutdownErr := s.Service.StopService(context.Background(), &protos.Void{})
		if shutdownErr != nil {
			glog.Errorf("Error shutting down orc8r service: %v", shutdownErr)
		}
		serverErr <- err
	}()
	return <-serverErr
}

// RunTest runs the test service on a given Listener and the HTTP on it's
// configured addr if exists. This function blocks by a signal or until a
// server is stopped.
func (s *OrchestratorService) RunTest(lis net.Listener) {
	s.State = protos.ServiceInfo_ALIVE
	s.Health = protos.ServiceInfo_APP_HEALTHY
	serverErr := make(chan error, 1)
	go func() {
		err := s.GrpcServer.Serve(lis)
		serverErr <- err
	}()
	if s.EchoServer != nil {
		go func() {
			err := s.EchoServer.StartServer(s.EchoServer.Server)
			serverErr <- err
		}()
	}
	err := <-serverErr
	if err != nil {
		glog.Fatal(err)
	}
}

func getEchoServerForOrchestratorService(serviceName string) (*echo.Echo, error) {
	if !runEchoServer {
		return nil, nil
	}
	echoPort, err := registry.GetEchoServerPort(serviceName)
	if err != nil {
		return nil, err
	}
	portStr := fmt.Sprintf(":%d", echoPort)
	e := echo.New()
	e.Server.Addr = portStr
	e.HideBanner = true
	return e, nil
}

type Config struct {
	// MaxGRPCMessageSizeMB is the maximum message size, in megabytes, allowed
	// by this service's gRPC servicer.
	//
	// Defaults:
	// - Server receive max:	4mb
	// - Server send max:		1gb
	// - Client receive max:	4mb
	// - Client send max:		1gb
	//
	// For simplicity, this config sets the receive max for both server and
	// client, leaving the send max unchanged.
	MaxGRPCMessageSizeMB int `yaml:"maxGRPCMessageSizeMB"`
}

func getSharedConfig() (*Config, error) {
	c := &Config{}

	_, _, err := config.GetStructuredServiceConfig(orc8r.ModuleName, orc8r.SharedService, c)
	if err != nil {
		return nil, err
	}

	if c.MaxGRPCMessageSizeMB == 0 {
		return nil, errors.New("parsed shared.yml and didn't find a max gRPC message size")
	}

	return c, nil
}
