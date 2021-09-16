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
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/magma/magma/agwd/config"
	"github.com/magma/magma/log"
	sctpdpb "github.com/magma/magma/protos/magma/sctpd"
	"github.com/magma/magma/service"
	"github.com/magma/magma/service/sctpd"
)

func newServiceRouter(cfgr config.Configer) service.Router {
	sctpdDownstreamConn, err := grpc.Dial(
		cfgr.Config().GetSctpdDownstreamServiceTarget(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	mmeGrpcConn, err := grpc.Dial(
		cfgr.Config().GetMmeSctpdUpstreamServiceTarget(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return service.NewRouter(
		sctpdpb.NewSctpdDownlinkClient(sctpdDownstreamConn),
		sctpdpb.NewSctpdUplinkClient(mmeGrpcConn),
	)
}

func startSctpdDownlinkServer(
	cfgr config.Configer, logger log.Logger, sr service.Router,
) {
	target := config.ParseTarget(cfgr.Config().GetMmeSctpdDownstreamServiceTarget())
	listener, err := net.Listen(
		target.Scheme, target.Endpoint)
	if err != nil {
		panic(errors.Wrapf(
			err,
			"net.Listen(network=%s, address=%s)",
			target.Scheme,
			target.Endpoint))
	}

	grpcServer := grpc.NewServer()
	sctpdDownlinkServer := sctpd.NewProxyDownlinkServer(logger, sr)
	sctpdpb.RegisterSctpdDownlinkServer(grpcServer, sctpdDownlinkServer)
	go grpcServer.Serve(listener)
}

func startSctpdUplinkServer(
	cfgr config.Configer, logger log.Logger, sr service.Router,
) {
	target := config.ParseTarget(cfgr.Config().GetSctpdUpstreamServiceTarget())
	listener, err := net.Listen(target.Scheme, target.Endpoint)
	if err != nil {
		panic(errors.Wrapf(
			err,
			"net.Listen(network=%s, address=%s)",
			target.Scheme,
			target.Endpoint))
	}

	grpcServer := grpc.NewServer()
	sctpdUplinkServer := sctpd.NewProxyUplinkServer(logger, sr)
	sctpdpb.RegisterSctpdUplinkServer(grpcServer, sctpdUplinkServer)
	go grpcServer.Serve(listener)
}

func Start(cfgr config.Configer, logger log.Logger) {
	sr := newServiceRouter(cfgr)
	startSctpdDownlinkServer(cfgr, logger, sr)
	startSctpdUplinkServer(cfgr, logger, sr)
}
