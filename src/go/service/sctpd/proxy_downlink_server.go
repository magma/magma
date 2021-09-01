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

package sctpd

import (
	"context"

	"github.com/magma/magma/log"
	pb "github.com/magma/magma/protos/magma/sctpd"
	"github.com/magma/magma/service"
)

// ProxyDownlinkServer handles SctpdDownlinkServer RPCs by calling out to a
// SctpdDownloadClient.
type ProxyDownlinkServer struct {
	service.Router
	log.Logger
	*pb.UnimplementedSctpdDownlinkServer
}

// NewProxyDownlinkServer returns a ProxyDownlinkServer injected with the
// provided logger and service router.
func NewProxyDownlinkServer(logger log.Logger, sr service.Router) *ProxyDownlinkServer {
	return &ProxyDownlinkServer{Logger: logger, Router: sr}
}

// Init proxies calls to SctpdDownlink.Init.
func (p *ProxyDownlinkServer) Init(ctx context.Context, req *pb.InitReq) (*pb.InitRes, error) {
	p.Logger.
		With("use_ipv4", req.GetUseIpv4()).
		With("use_ipv6", req.GetUseIpv6()).
		With("ipv4_addrs", req.GetIpv4Addrs()).
		With("ipv6_addrs", req.GetIpv6Addrs()).
		With("port", req.GetPort()).
		With("ppid", req.GetPpid()).
		With("force_restart", req.GetForceRestart()).
		Debug().Print("Init")
	return p.SctpdDownlinkClient().Init(ctx, req)
}

// SendDl proxies calls to SctpdDownlink.SendDl.
func (p *ProxyDownlinkServer) SendDl(ctx context.Context, req *pb.SendDlReq) (*pb.SendDlRes, error) {
	p.Logger.
		With("assoc_id", req.GetAssocId()).
		With("stream", req.GetStream()).
		With("ppid", req.GetPpid()).
		With("payload_size", len(req.GetPayload())).
		Debug().Print("SendDl")
	return p.SctpdDownlinkClient().SendDl(ctx, req)
}
