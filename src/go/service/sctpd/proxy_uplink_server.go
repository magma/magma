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
	"net"

	"github.com/magma/magma/log"
	pb "github.com/magma/magma/protos/magma/sctpd"
	"github.com/magma/magma/service"
)

// ProxyUplinkServer handles SctpdUplinkServer RPCs by calling out to a
// SctpdUplinkClient.
type ProxyUplinkServer struct {
	service.Router
	log.Logger
	*pb.UnimplementedSctpdUplinkServer
}

// NewProxyUplinkServer creates a new ProxyUplinkServer with the provided
// logger and client connection.
func NewProxyUplinkServer(logger log.Logger, sr service.Router) *ProxyUplinkServer {
	return &ProxyUplinkServer{Logger: logger, Router: sr}
}

// SendUl proxies calls to SctpdUplink.SendUl.
func (p *ProxyUplinkServer) SendUl(ctx context.Context, req *pb.SendUlReq) (*pb.SendUlRes, error) {
	p.Logger.
		With("assoc_id", req.GetAssocId()).
		With("stream", req.GetStream()).
		With("ppid", req.GetPpid()).
		With("payload_size", len(req.GetPayload())).
		Debug().Print("SendUl")
	return p.SctpdUplinkClient().SendUl(ctx, req)
}

// NewAssoc proxies calls to SctpdUplink.NewAssoc.
func (p *ProxyUplinkServer) NewAssoc(ctx context.Context, req *pb.NewAssocReq) (*pb.NewAssocRes, error) {
	p.Logger.
		With("assoc_id", req.GetAssocId()).
		With("instreams", req.GetInstreams()).
		With("outstreams", req.GetOutstreams()).
		With("ppid", req.GetPpid()).
		With("ran_cp_ipaddr", net.IP(req.GetRanCpIpaddr()).String()).
		Debug().Print("NewAssoc")
	return p.SctpdUplinkClient().NewAssoc(ctx, req)
}

// CloseAssoc proxies calls to SctpdUplink.CloseAssoc.
func (p *ProxyUplinkServer) CloseAssoc(ctx context.Context, req *pb.CloseAssocReq) (*pb.CloseAssocRes, error) {
	p.Logger.
		With("assoc_id", req.GetAssocId()).
		With("is_reset", req.GetIsReset()).
		With("ppid", req.GetPpid()).
		Debug().Print("CloseAssoc")
	return p.SctpdUplinkClient().CloseAssoc(ctx, req)
}
