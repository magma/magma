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

package servicers

import (
	"flag"
	"fmt"
	"strconv"
	"net"

	"context"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/service/config"
	"magma/orc8r/lib/go/registry"
)

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

func NewService() *Service {
	return &Service{}
}

// GetServiceInfo returns service-level info (name, version, status, etc...)
func (service *Service) GetServiceInfo(ctx context.Context, void *protos.Void) (*protos.ServiceInfo, error) {
	return &protos.ServiceInfo{
		Name:          service.Type,
		Version:       service.Version,
		State:         service.State,
		Health:        service.Health,
		StartTimeSecs: service.StartTimeSecs,
	}, nil
}

// StopService is a request to stop the service gracefully.
func (service *Service) StopService(ctx context.Context, void *protos.Void) (*protos.Void, error) {
	service.State = protos.ServiceInfo_STOPPING
	go service.GrpcServer.GracefulStop()
	service.Health = protos.ServiceInfo_APP_UNHEALTHY
	return new(protos.Void), nil
}

// GetMetrics returns a MetricsContainer with all metrics for the service.
func (service *Service) GetMetrics(ctx context.Context, void *protos.Void) (*protos.MetricsContainer, error) {
	met := &protos.MetricsContainer{}
	metricsFamilies, err := metrics.GetMetrics()
	if err != nil {
		return met, err
	}
	met.Family = metricsFamilies
	return met, nil
}

// SetLogLevel sets the logging level for the service.
func (service *Service) SetLogLevel(ctx context.Context, logLevelMsg *protos.LogLevelMessage) (*protos.Void, error) {
	// TODO: set log level
	return new(protos.Void), fmt.Errorf("SetLogLevel not implemented\n")
}

// SetLogVerbosity sets the glog logging verbosity for the service.
func (service *Service) SetLogVerbosity(ctx context.Context, verbosity *protos.LogVerbosity) (*protos.Void, error) {
	flag.Lookup("v").Value.Set(strconv.Itoa(int(verbosity.Verbosity)))
	return new(protos.Void), nil
}

// ReloadServiceConfig not currently implemented for cloud services
func (service *Service) ReloadServiceConfig(ctx context.Context, void *protos.Void) (*protos.ReloadConfigResponse, error) {
	res := protos.ReloadConfigResponse{}
	return &res, fmt.Errorf("method ReloadServiceConfig  not implemented\n")
}

// GetOperationalStates not currently implemented for go services
func (service *Service) GetOperationalStates(ctx context.Context, void *protos.Void) (*protos.GetOperationalStatesResponse, error) {
	res := protos.GetOperationalStatesResponse{}
	return &res, nil
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
