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
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/magma/magma/src/go/capture"
	capturepb "github.com/magma/magma/src/go/protos/magma/capture"
	configpb "github.com/magma/magma/src/go/protos/magma/config"
	service_capture "github.com/magma/magma/src/go/service/capture"
	service_config "github.com/magma/magma/src/go/service/config"
	"google.golang.org/grpc"

	"github.com/magma/magma/src/go/agwd/config"
	"github.com/magma/magma/src/go/log"
	pipelinedpb "github.com/magma/magma/src/go/protos/magma/pipelined"
	sctpdpb "github.com/magma/magma/src/go/protos/magma/sctpd"
	"github.com/magma/magma/src/go/service"
	"github.com/magma/magma/src/go/service/pipelined"
	"github.com/magma/magma/src/go/service/sctpd"
)

func newServiceRouter(cfgr config.Configer) service.Router {
	sctpdDownstreamConn, err := grpc.Dial(cfgr.Config().GetSctpdDownstreamServiceTarget(), grpc.WithInsecure())
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

func cleanupUnixSocket(
	logger log.Logger,
	osStat func(string) (os.FileInfo, error),
	osRemove func(string) error,
	netDialTimeout func(network, address string, timeout time.Duration) (net.Conn, error),
	path string) error {
	_, err := osStat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("os.Stat(%s): %w", path, err)
	}

	// attempt to connect to see if someone is still bound to the socket.
	if runtime.GOOS == "windows" {
		return fmt.Errorf("cleanupUnixSocket(logger=_,path=%q) found pre-existing socket at path, but does not support cleanup in GOOS == windows", path)
	}
	_, err = netDialTimeout("unix", path, time.Second)
	if err == nil {
		return fmt.Errorf("os.Stat(%s): existing listener on socket file", path)
	}
	logger.Warning().Printf("Removing existing socket file; previous unclean shutdown?")
	if err := osRemove(path); err != nil {
		return fmt.Errorf("os.Stat(%s): %w", path, err)
	}

	return nil
}

func cleanupUnixSocketOrDie(logger log.Logger, path string) {
	if err := cleanupUnixSocket(logger, os.Stat, os.Remove, net.DialTimeout, path); err != nil {
		panic(fmt.Errorf(
			"cleanupUnixSocket(logger=_, target.Endpoint=%s): %w",
			path,
			err))
	}
}

func startSctpdDownlinkServer(
	cfgr config.Configer, logger log.Logger, sr service.Router, so ...grpc.ServerOption,
) {
	target := config.ParseTarget(cfgr.Config().GetMmeSctpdDownstreamServiceTarget())
	if target.Scheme == "unix" {
		cleanupUnixSocketOrDie(logger, target.Endpoint)
	}

	listener, err := net.Listen(
		target.Scheme, target.Endpoint)
	if err != nil {
		panic(fmt.Errorf(
			"net.Listen(network=%s, address=%s): %w",
			target.Scheme,
			target.Endpoint,
			err))
	}

	grpcServer := grpc.NewServer(so...)
	sctpdDownlinkServer := sctpd.NewProxyDownlinkServer(logger, sr)
	sctpdpb.RegisterSctpdDownlinkServer(grpcServer, sctpdDownlinkServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(fmt.Errorf(
				"startSctpdDownlinkServer(network=%s, address=%s): %w",
				target.Scheme,
				target.Endpoint,
				err))
		}
	}()
}

func startSctpdUplinkServer(
	cfgr config.Configer, logger log.Logger, sr service.Router, so ...grpc.ServerOption,
) {
	target := config.ParseTarget(cfgr.Config().GetSctpdUpstreamServiceTarget())
	if target.Scheme == "unix" {
		cleanupUnixSocketOrDie(logger, target.Endpoint)
	}

	listener, err := net.Listen(target.Scheme, target.Endpoint)
	if err != nil {
		panic(fmt.Errorf(
			"net.Listen(network=%s, address=%s): %w",
			target.Scheme,
			target.Endpoint,
			err))
	}

	grpcServer := grpc.NewServer(so...)
	sctpdUplinkServer := sctpd.NewProxyUplinkServer(logger, sr)
	sctpdpb.RegisterSctpdUplinkServer(grpcServer, sctpdUplinkServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(fmt.Errorf(
				"startSctpdUplinkServer(network=%s, address=%s): %w",
				target.Scheme,
				target.Endpoint,
				err))
		}
	}()
}

func startPipelinedServer(cfgr config.Configer, logger log.Logger) {
	target := config.ParseTarget(cfgr.Config().GetPipelinedServiceTarget())
	listener, err := net.Listen(target.Scheme, target.Endpoint)

	if err != nil {
		panic(fmt.Errorf(
			"net.Listen(network=%s, address=%s): %w",
			target.Scheme,
			target.Endpoint,
			err))
	}
	grpcServer := grpc.NewServer()
	pipelinedServer := pipelined.NewPipelinedServer(logger)
	pipelinedpb.RegisterPipelinedServer(grpcServer, pipelinedServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(fmt.Errorf(
				"startPipelinedServer(network=%s, address=%s): %w",
				target.Scheme,
				target.Endpoint,
				err))
		}
	}()
}

func startCaptureServer(
	cfgr config.Configer, logger log.Logger, buf *capture.Buffer,
) {
	address := fmt.Sprintf(":%s", cfgr.Config().GetCaptureServicePort())
	listener, err := net.Listen(config.TCP, address)
	if err != nil {
		panic(fmt.Errorf(
			"net.Listen(network=%s, address=%s): %w",
			config.TCP,
			address,
			err))
	}

	grpcServer := grpc.NewServer()
	captureServer := service_capture.NewCaptureServer(logger, buf)
	capturepb.RegisterCaptureServer(grpcServer, captureServer)
	go grpcServer.Serve(listener)
}

func startConfigServer(
	cfgr config.Configer, logger log.Logger,
) {
	address := fmt.Sprintf(":%s", cfgr.Config().GetConfigServicePort())

	listener, err := net.Listen(config.TCP, address)
	if err != nil {
		panic(fmt.Errorf(
			"net.Listen(network=%s, address=%s): %w",
			config.TCP,
			address,
			err))
	}

	grpcServer := grpc.NewServer()
	configServer := service_config.NewConfigServer(logger, cfgr)
	configpb.RegisterConfigServer(grpcServer, configServer)
	go grpcServer.Serve(listener)
}

func Start(cfgr config.Configer, logger log.Logger) {
	buf := capture.NewBuffer()
	mw := capture.NewMiddleware(cfgr, buf)
	serverOptions := mw.GetServerOptions()
	sr := newServiceRouter(cfgr)
	startConfigServer(cfgr, logger)
	startCaptureServer(cfgr, logger, buf)
	startSctpdDownlinkServer(cfgr, logger, sr, serverOptions...)
	startSctpdUplinkServer(cfgr, logger, sr, serverOptions...)
	startPipelinedServer(cfgr, logger)
}
