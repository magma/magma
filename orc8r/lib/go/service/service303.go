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

package service

import (
	"flag"
	"fmt"
	"strconv"

	"golang.org/x/net/context"

	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"
)

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
