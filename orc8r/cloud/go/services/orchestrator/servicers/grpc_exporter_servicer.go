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
	"time"

	"magma/orc8r/cloud/go/services/metricsd/protos"

	"github.com/golang/glog"
	"github.com/prometheus/client_model/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	edgeHubPB "github.com/facebookincubator/prometheus-edge-hub/grpc"
)

const (
	grpcMaxTimeoutSec = 60
	grpcMaxDelaySec   = 20
	grpcMaxMsgSize    = 1024 * 1024 * 1024 // 1 Gb
)

type GRPCPushExporterServicer struct {
	PushAddress string
	GrpcClient  EdgeHubClient
}

type EdgeHubClient interface {
	Collect(ctx context.Context, in *edgeHubPB.MetricFamilies, opts ...grpc.CallOption) (*edgeHubPB.Void, error)
}

// NewGRPCPushExporterServicer returns an exporter pushing metrics to
// prometheus-edge-hubs at the given addresses.
func NewGRPCPushExporterServicer(pushAddr string) protos.MetricsExporterServer {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, pushAddr, getDialOptions()...)
	if err != nil {
		glog.Fatalf("Error creating grpc metrics exporter: %v", err)
	}
	conn.ResetConnectBackoff()
	client := edgeHubPB.NewMetricsControllerClient(conn)

	srv := &GRPCPushExporterServicer{
		GrpcClient:  client,
		PushAddress: pushAddr,
	}
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
	_, err := s.GrpcClient.Collect(context.Background(), &edgeHubPB.MetricFamilies{
		Families: families,
	})
	return err
}

func getDialOptions() []grpc.DialOption {
	bckoff := backoff.DefaultConfig
	bckoff.MaxDelay = grpcMaxDelaySec * time.Second

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxMsgSize)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           bckoff,
			MinConnectTimeout: grpcMaxTimeoutSec * time.Second},
		),
	}
	return opts
}
