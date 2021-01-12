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
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2"
)

type s8Proxy struct {
	config    *S8ProxyConfig
	gtpClient *gtp.Client
}

type S8ProxyConfig struct {
	ClientAddr string
	ServerAddr string
}

func NewS8Proxy(config *S8ProxyConfig) (*s8Proxy, error) {
	// TODO: validate config
	gtpCli, err := gtp.NewConnectedAutoClient(context.Background(), config.ServerAddr, gtpv2.IFTypeS5S8SGWGTPC)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	addS8GtpHandlers(gtpCli)
	return &s8Proxy{
		config:    config,
		gtpClient: gtpCli,
	}, nil
}

func (s *s8Proxy) CreateSession(ctx context.Context, req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	// build csReq IE message
	csReq, sessionTeids := buildCreateSessionRequestIE(req, s.gtpClient.Conn, s.gtpClient.GetServerAddress())

	// send, register and receive create session
	csRes, err := s.sendAndReceiveCreateSession(csReq, sessionTeids)
	if err != nil {
		err = fmt.Errorf("Create Session Request failed: %s", err)
		glog.Error(err)
		return &protos.CreateSessionResponsePgw{}, err
	}

	// TODO: build grpc CreateSessionResponsePgw message
	glog.V(2).Infof("This is session response %+v", csRes)

	return &protos.CreateSessionResponsePgw{}, nil
}

// TODO
func (s *s8Proxy) ModifyBearer(ctx context.Context, req *protos.ModifyBearerRequestPgw) (*protos.ModifyBearerResponsePgw, error) {
	return &protos.ModifyBearerResponsePgw{}, nil
}

// TODO
func (s *s8Proxy) DeleteSession(ctx context.Context, req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	return &protos.DeleteSessionResponsePgw{}, nil
}

// TODO
func (s *s8Proxy) SendEcho(ctx context.Context, req *orcprotos.Void) (*protos.EchoResponse, error) {
	return &protos.EchoResponse{}, nil
}
