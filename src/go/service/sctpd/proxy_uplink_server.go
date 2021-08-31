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

	"google.golang.org/grpc"

	"github.com/magma/magma/log"
	"github.com/magma/magma/protos/magma/sctpd"
)

type ProxyUplinkServer struct {
	client sctpd.SctpdUplinkClient

	log.Logger
	*sctpd.UnimplementedSctpdUplinkServer
}

func NewProxyUplinkServer(logger log.Logger, cc *grpc.ClientConn) *ProxyUplinkServer {
	return &ProxyUplinkServer{Logger: logger, client: sctpd.NewSctpdUplinkClient(cc)}
}

func (p *ProxyUplinkServer) SendUl(ctx context.Context, req *sctpd.SendUlReq) (*sctpd.SendUlRes, error) {
	p.Logger.
		With("assoc_id", req.GetAssocId()).
		With("stream", req.GetStream()).
		With("ppid", req.GetPpid()).
		With("payload_size", len(req.GetPayload())).
		Debug().Print("SendUl")
	return p.client.SendUl(ctx, req)
}

func (p *ProxyUplinkServer) NewAssoc(ctx context.Context, req *sctpd.NewAssocReq) (*sctpd.NewAssocRes, error) {
	p.Logger.
		With("assoc_id", req.GetAssocId()).
		With("instreams", req.GetInstreams()).
		With("outstreams", req.GetOutstreams()).
		With("ppid", req.GetPpid()).
		With("ran_cp_ipaddr", net.IP(req.GetRanCpIpaddr()).String()).
		Debug().Print("NewAssoc")
	return p.client.NewAssoc(ctx, req)
}

func (p *ProxyUplinkServer) CloseAssoc(ctx context.Context, req *sctpd.CloseAssocReq) (*sctpd.CloseAssocRes, error) {
	p.Logger.
		With("assoc_id", req.GetAssocId()).
		With("is_reset", req.GetIsReset()).
		With("ppid", req.GetPpid()).
		Debug().Print("CloseAssoc")
	return p.client.CloseAssoc(ctx, req)
}
