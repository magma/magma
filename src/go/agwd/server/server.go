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

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/magma/magma/src/go/agwd/config"
	"github.com/magma/magma/src/go/log"
	sctpdpb "github.com/magma/magma/src/go/protos/magma/sctpd"
	"github.com/magma/magma/src/go/service"
	"github.com/magma/magma/src/go/service/sctpd"
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
		return errors.Wrapf(err, "os.Stat(%s)", path)
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
		return errors.Wrapf(err, "os.Stat(%s)", path)
	}

	return nil
}

func cleanupUnixSocketOrDie(logger log.Logger, path string) {
	if err := cleanupUnixSocket(logger, os.Stat, os.Remove, net.DialTimeout, path); err != nil {
		panic(errors.Wrapf(
			err,
			"cleanupUnixSocket(logger=_, target.Endpoint=%s)",
			path))
	}
}

func startSctpdDownlinkServer(
	cfgr config.Configer, logger log.Logger, sr service.Router,
) {
	target := config.ParseTarget(cfgr.Config().GetMmeSctpdDownstreamServiceTarget())
	if target.Scheme == "unix" {
		cleanupUnixSocketOrDie(logger, target.Endpoint)
	}

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
