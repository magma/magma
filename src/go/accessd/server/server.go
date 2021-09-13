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
	"github.com/magma/magma/service/sctpd"
)

const (
	mmeUpstreamTarget  = "unix:///tmp/mme_sctpd_upstream.sock"
	mmeGrpcDialTimeout = time.Second
)

func newMmeGrpcConn() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mmeGrpcDialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, mmeUpstreamTarget, grpc.WithInsecure(), grpc.WithBlock())
}

func newServer(logger log.Logger) *grpc.Server {
	grpcServer := grpc.NewServer()

	mmeGrpcConn, err := newMmeGrpcConn()
	if err != nil {
		panic(err)
	}
	sctpdUplinkServer := sctpd.NewProxyUplinkServer(logger, mmeGrpcConn)
	sctpdpb.RegisterSctpdUplinkServer(grpcServer, sctpdUplinkServer)

	return grpcServer
}

const (
	sctpdUpstreamNetwork = "unix"
	sctpdUpstreamPath    = "/tmp/sctpd_upstream.sock"
)

func Start(logger log.Logger) error {
	listener, err := net.Listen(sctpdUpstreamNetwork, sctpdUpstreamPath)
	if err != nil {
		return errors.Wrapf(
			err,
			"net.Listen(network=%s, address=%s)",
			sctpdUpstreamNetwork,
			sctpdUpstreamPath)
	}
	return errors.WithStack(newServer(logger).Serve(listener))
}
