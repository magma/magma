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
	"sync"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	grpc_proto "google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/keepalive"

	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/service/config"
	"magma/orc8r/lib/go/util"
	"magma/orc8r/lib/go/service/servicers/protected"
)

const (
	PrintGrpcPayloadFlag = "print-grpc-payload"
	PrintGrpcPayloadEnv  = "MAGMA_PRINT_GRPC_PAYLOAD"
)

var (
	printGrpcPayload           bool
	currentlyRunningServices   = make(map[string]*servicers.Service)
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
	flag.BoolVar(&printGrpcPayload, PrintGrpcPayloadFlag, false, "Enable GRPC Payload Printout")
}

// NewServiceWithOptions returns a new GRPC orchestrator service implementing
// service303 with the specified grpc server options.
// It will not instantiate the service with the identity checking middleware.
func NewServiceWithOptions(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*servicers.Service, error) {
	return NewServiceWithOptionsImpl(moduleName, serviceName, serverOptions...)
}

// NewServiceWithOptionsImpl returns a new GRPC service implementing
// service303 with the specified grpc server options. This will not instantiate
// the service with the identity checking middleware.
func NewServiceWithOptionsImpl(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*servicers.Service, error) {
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
	service := &servicers.Service{
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

// GetDefaultKeepaliveParameters returns the default keepalive server parameters.
func GetDefaultKeepaliveParameters() keepalive.ServerParameters {
	return defaultKeepaliveParams
}
