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
	"context"

	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/lib/go/registry"

	edge_hub "github.com/facebookincubator/prometheus-edge-hub/grpc"
	"github.com/pkg/errors"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"google.golang.org/grpc"
)

const (
	serviceName = "grpc_metrics_exporter"

	grpcMaxMsgSize = 1024 * 1024 * 1024 // 1 Gb
)

var (
	dialOpts = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxMsgSize)),
	}
)

type GRPCPushExporterServicer struct {
	// registry is the servicer's local service registry.
	// Local registry since the gRPC servicer is not a proper Orchestrator
	// service.
	registry *registry.ServiceRegistry
}

// NewGRPCPushExporterServicer returns an exporter pushing metrics to
// prometheus-edge-hubs at the given addresses.
func NewGRPCPushExporterServicer(pushAddr string) protos.MetricsExporterServer {
	srv := &GRPCPushExporterServicer{registry: registry.NewWithMode(registry.YamlRegistryMode)}
	srv.registry.AddService(registry.ServiceLocation{Name: serviceName, Host: pushAddr})
	return srv
}

func (s *GRPCPushExporterServicer) Submit(ctx context.Context, req *protos.SubmitMetricsRequest) (*protos.SubmitMetricsResponse, error) {
	metricsToSend := processMetrics(req.GetMetrics())
	if len(metricsToSend) == 0 {
		return &protos.SubmitMetricsResponse{}, nil
	}
	err := s.pushFamilies(metricsToSend)
	return &protos.SubmitMetricsResponse{}, err
}

func (s *GRPCPushExporterServicer) pushFamilies(families []*io_prometheus_client.MetricFamily) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}
	_, err = client.Collect(context.Background(), &edge_hub.MetricFamilies{Families: families})
	if err != nil {
		return err
	}
	return nil
}

func (s *GRPCPushExporterServicer) getClient() (edge_hub.MetricsControllerClient, error) {
	conn, err := s.registry.GetConnectionWithOptions(serviceName, dialOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "get exporter client connection")
	}
	client := edge_hub.NewMetricsControllerClient(conn)
	return client, nil
}

// EdgeHubServer is an alias of MetricsControllerServer to support
// locally-generated mocks.
type EdgeHubServer interface {
	Collect(ctx context.Context, families *edge_hub.MetricFamilies) (*edge_hub.Void, error)
	ForTestsOnly()
}
