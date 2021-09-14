// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/magma/magma/log"
	sctpdpb "github.com/magma/magma/protos/magma/sctpd"
	"github.com/magma/magma/service"
	"github.com/magma/magma/service/sctpd"
)

func newGrpcConn(target string, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return grpc.DialContext(ctx, target, grpc.WithInsecure())
}

const (
	sctpdDownstreamTarget = "unix:///tmp/mme_sctpd_downstream.sock"
	sctpdGrpcDialTimeout  = time.Second
	mmeUpstreamTarget     = "unix:///tmp/mme_sctpd_upstream.sock"
	mmeGrpcDialTimeout    = time.Second
)

func newServiceRouter() service.Router {
	sctpdDownstreamConn, err := newGrpcConn(sctpdDownstreamTarget, sctpdGrpcDialTimeout)
	if err != nil {
		panic(err)
	}
	mmeGrpcConn, err := newGrpcConn(mmeUpstreamTarget, mmeGrpcDialTimeout)
	if err != nil {
		panic(err)
	}
	return service.NewRouter(
		sctpdpb.NewSctpdDownlinkClient(sctpdDownstreamConn),
		sctpdpb.NewSctpdUplinkClient(mmeGrpcConn),
	)
}

const (
	sctpdDownstreamNetwork = "unix"
	sctpdDownstreamPath    = "/tmp/sctpd_downstream.sock"
)

func startSctpdDownlinkServer(logger log.Logger, sr service.Router) {
	listener, err := net.Listen(sctpdDownstreamNetwork, sctpdDownstreamPath)
	if err != nil {
		panic(errors.Wrapf(
			err,
			"net.Listen(network=%s, address=%s)",
			sctpdDownstreamNetwork,
			sctpdDownstreamPath))
	}

	grpcServer := grpc.NewServer()
	sctpdDownlinkServer := sctpd.NewProxyDownlinkServer(logger, sr)
	sctpdpb.RegisterSctpdDownlinkServer(grpcServer, sctpdDownlinkServer)
	go grpcServer.Serve(listener)
}

const (
	sctpdUpstreamNetwork = "unix"
	sctpdUpstreamPath    = "/tmp/sctpd_upstream.sock"
)

func startSctpdUplinkServer(logger log.Logger, sr service.Router) {
	listener, err := net.Listen(sctpdUpstreamNetwork, sctpdUpstreamPath)
	if err != nil {
		panic(errors.Wrapf(
			err,
			"net.Listen(network=%s, address=%s)",
			sctpdUpstreamNetwork,
			sctpdUpstreamPath))
	}

	grpcServer := grpc.NewServer()
	sctpdUplinkServer := sctpd.NewProxyUplinkServer(logger, sr)
	sctpdpb.RegisterSctpdUplinkServer(grpcServer, sctpdUplinkServer)
	go grpcServer.Serve(listener)
}

func Start(logger log.Logger) {
	sr := newServiceRouter()
	startSctpdDownlinkServer(logger, sr)
	startSctpdUplinkServer(logger, sr)
}
