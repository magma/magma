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
	"net"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"
)

type echoResponse struct {
	error
}

type S8Proxy struct {
	config      *S8ProxyConfig
	gtpClient   *gtp.Client
	echoChannel chan (error)
}

type S8ProxyConfig struct {
	ClientAddr string
	ServerAddr *net.UDPAddr
}

// NewS8Proxy creates an s8 proxy, but does not checks the PGW is alive
func NewS8Proxy(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewRunningClient(
		context.Background(), config.ClientAddr, gtpv2.IFTypeS5S8SGWGTPC)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	return newS8ProxyImp(gtpCli, config)
}

//NewS8ProxyWithEcho creates an s8 proxy already connected to a server (checks with echo if PGW is alive)
func NewS8ProxyWithEcho(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewConnectedAutoClient(context.Background(), config.ServerAddr.String(), gtpv2.IFTypeS5S8SGWGTPC)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	return newS8ProxyImp(gtpCli, config)
}

func newS8ProxyImp(cli *gtp.Client, config *S8ProxyConfig) (*S8Proxy, error) {
	// TODO: validate config
	s8p := &S8Proxy{
		config:      config,
		gtpClient:   cli,
		echoChannel: make(chan error),
	}
	addS8GtpHandlers(s8p)
	return s8p, nil
}

func (s *S8Proxy) CreateSession(ctx context.Context, req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	cPgwUDPAddr, err := s.getPgwAddress(req.PgwAddrs)
	if err != nil {
		err = fmt.Errorf("Create Session Request failed due to missing server address: %s", err)
		glog.Error(err)
		return nil, err
	}

	// build csReq IE message
	sessionTeids, csReqIEs, err := buildCreateSessionRequestIE(cPgwUDPAddr, req, s.gtpClient)
	if err != nil {
		return nil, err
	}

	// send, register and receive create session (session is created on the gtp client during this process too)
	csRes, err := s.sendAndReceiveCreateSession(cPgwUDPAddr, sessionTeids, csReqIEs)
	if err != nil {
		err = fmt.Errorf("Create Session Request failed: %s", err)
		glog.Error(err)
		return nil, err
	}
	return csRes, nil
}

// TODO: see if ModifyBearerRequest applies to S8 for Magma
func (s *S8Proxy) ModifyBearer(ctx context.Context, req *protos.ModifyBearerRequestPgw) (*protos.ModifyBearerResponsePgw, error) {
	// Todo: delete this condition
	err := fmt.Errorf("ModifyBearer is not completed implemented")
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	session, teid, err := s.gtpClient.GetSessionAndCTeidByIMSI(req.Imsi)
	if err != nil {
		return nil, err
	}
	mbReqIEs := buildModifyBearerRequest(req, session.GetDefaultBearer().EBI)
	mdRes, err := s.sendAndReceiveModifyBearer(teid, session, mbReqIEs)
	if err != nil {
		err = fmt.Errorf("Modify Bearer Request failed: %s", err)
		glog.Error(err)
		return nil, err
	}
	return mdRes, nil
}

func (s *S8Proxy) DeleteSession(ctx context.Context, req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	session, teid, err := s.gtpClient.GetSessionAndCTeidByIMSI(req.Imsi)
	if err != nil {
		return nil, err
	}

	dsReqIE := buildDeleteSessionRequest(session)
	cdRes, err := s.sendAndReceiveDeleteSession(teid, session, dsReqIE)
	if err != nil {
		glog.Errorf("Couldnt delete session for IMSI %s:, %s", req.Imsi, err)
		return nil, err
	}
	// remove session from the s8_proxy client
	s.gtpClient.RemoveSession(session)
	return cdRes, nil
}

func (s *S8Proxy) SendEcho(ctx context.Context, req *protos.EchoRequest) (*protos.EchoResponse, error) {
	cPgwUDPAddr, err := s.getPgwAddress(req.PgwAddrs)
	if err != nil {
		err = fmt.Errorf("SendEcho to %s failed: %s", cPgwUDPAddr, err)
		glog.Error(err)
		return nil, err
	}

	err = s.sendAndReceiveEchoRequest(cPgwUDPAddr)
	if err != nil {
		return nil, err
	}
	return &protos.EchoResponse{}, nil
}

// getPgwAddress returns an UDPAddrs if the passed string corresponds to a valid ip,
// otherwise it uses the server address configured on s8_proxy
func (s *S8Proxy) getPgwAddress(pgwAddrsFromRequest string) (*net.UDPAddr, error) {
	addrs := ParseAddress(pgwAddrsFromRequest)
	if addrs != nil {
		// address comming from string has precednece
		return addrs, nil
	}
	if s.config.ServerAddr != nil {
		return s.config.ServerAddr, nil
	}
	return nil, fmt.Errorf("Neither the request nor s8_proxy has a valid server (pgw) address")
}
