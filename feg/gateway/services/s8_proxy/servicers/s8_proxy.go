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
	ServerAddr string
}

// NewS8Proxy creates an s8 proxy already connected to a server (checks with echo if PGW is alive)
func NewS8Proxy(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewConnectedAutoClient(context.Background(), config.ServerAddr, gtpv2.IFTypeS5S8SGWGTPC)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	return newS8ProxyImp(gtpCli, config)
}

//NewS8ProxyNoFirstEcho creates an s8 proxy, but does not checks the PGW is alive
func NewS8ProxyNoFirstEcho(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewRunningAutoClient(context.Background(), config.ServerAddr, gtpv2.IFTypeS5S8SGWGTPC)
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
	// build csReq IE message
	csReqIEs, sessionTeids, err := buildCreateSessionRequestIE(req, s.gtpClient)
	if err != nil {
		return nil, err
	}

	// send, register and receive create session (session is created on the gtp client during this process too)
	csRes, err := s.sendAndReceiveCreateSession(csReqIEs, sessionTeids)
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

	session, teid, err := getSessionAndCTeid(s.gtpClient, req.Imsi)
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
	session, teid, err := getSessionAndCTeid(s.gtpClient, req.Imsi)
	if err != nil {
		return nil, err
	}

	cdRes, err := s.sendAndReceiveDeleteSession(teid, session)
	if err != nil {
		glog.Errorf("Couldnt delete session for IMSI %s:, %s", req.Imsi, err)
		return nil, err
	}

	// remove session from the s8_proxy client
	s.gtpClient.RemoveSession(session)

	return cdRes, nil
}

func (s *S8Proxy) SendEcho(ctx context.Context, req *protos.EchoRequest) (*protos.EchoResponse, error) {
	err := s.sendAndReceiveEchoRequest()
	if err != nil {
		return nil, err
	}
	return &protos.EchoResponse{}, nil
}

// TODO: this is a function to expose WaitUntilClientIsReady. That function is only used
// as a hack for testing and will be removed.
func (s *S8Proxy) WaitUntilClientIsReady() {
	s.gtpClient.WaitUntilClientIsReady(0)
}

func getSessionAndCTeid(cli *gtp.Client, imsi string) (*gtpv2.Session, uint32, error) {
	session, err := cli.GetSessionByIMSI(imsi)
	if err != nil {
		glog.Errorf("Couldnt delete session. Couldnt find a session for IMSI %s:, %s", imsi, err)
		return nil, 0, err
	}
	teid, err := session.GetTEID(gtpv2.IFTypeS5S8PGWGTPC)
	if err != nil {
		glog.Errorf("Couldnt delete session. Couldnt find control TEID for IMSI %s:, %s", imsi, err)
		return nil, 0, err
	}
	return session, teid, nil
}
