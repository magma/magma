/*
Copyright 2021 The Magma Authors.

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
	"flag"
	"sync"

	"github.com/golang/glog"
	anpb "github.com/magma/augmented-networks/accounting/protos"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/basic_acct"
	"magma/orc8r/lib/go/service/config"
)

const (
	remoteAddrFlagName = "remote_addr"
	caCertFlagName     = "ca_crt"
	clientCertFlagName = "client_crt"
	clientKeyFlagName  = "client_key"
	notlsFlagName      = "notls"
	insecureFlagName   = "insecure"
)

var (
	// Remote address
	remoteAddrFlag = flag.String(remoteAddrFlagName, "", "remote accounting service address")

	// Certificate files.
	caCertFlag     = flag.String(caCertFlagName, "", "CA certificate file. Used to verify server TLS certificate.")
	clientCertFlag = flag.String(clientCertFlagName, "", "Client certificate file. Used for client certificate-based authentication.")
	clientKeyFlag  = flag.String(clientKeyFlagName, "", "Client private key file. Used for client certificate-based authentication.")

	// notls flag, use to disable TLS
	notlsFlag = flag.Bool(notlsFlagName, false, `Disable TLS when set`)

	// insecure flag, use to disable server TLS cert verification
	insecureFlag = flag.Bool(insecureFlagName, false, `Disable TLS server certificate verification`)
)

// BaseAccService Configuration
type Config struct {
	RemoteAddr   string `yaml:"remote_address"`
	RootCaCert   string `yaml:"root_ca_cert"`
	ClientCrt    string `yaml:"client_cert"`
	ClientCrtKey string `yaml:"client_cert_key"`
	NoTls        bool   `yaml:"no_tls"`
	Insecure     bool   `yaml:"insecure"`
}

// BaseAccService
type BaseAccService struct {
	sync.RWMutex
	cfg        *Config
	remoteConn *grpc.ClientConn
}

// GetConfig returns a copy of current server configs
func (s *BaseAccService) GetConfig() *Config {
	if s != nil {
		s.RLock()
		defer s.RUnlock()
		if s.cfg != nil {
			cfgCopy := *s.cfg
			return &cfgCopy
		}
	}
	return nil
}

// Set Config updates current server configs
func (s *BaseAccService) SetConfig(cfg *Config) {
	if s != nil && cfg != nil {
		s.Lock()
		s.cfg = cfg
		s.Unlock()
	}
}

// NewBaseAcctService creates a new BaseAccService and initializes its configs
func NewBaseAcctService() *BaseAccService {
	cfg := &Config{
		RemoteAddr:   *remoteAddrFlag,
		RootCaCert:   *caCertFlag,
		ClientCrt:    *clientCertFlag,
		ClientCrtKey: *clientKeyFlag,
		NoTls:        *notlsFlag,
		Insecure:     *insecureFlag,
	}
	_, _, err := config.GetStructuredServiceConfig(feg.ModuleName, basic_acct.ServiceName, cfg)
	if err != nil {
		glog.Warningf("Failed reading '%s' service config: %v; Using the following defaults: %+v",
			basic_acct.ServiceName, err, *cfg)
	}
	// Update configs from command flags, they should overwrite .yml based configs
	cfg.updateFromFlags()
	return &BaseAccService{cfg: cfg}
}

// Start will be called at the end of every new user session creation
// start is responsible for verification & initiation of an accounting contract
// between the user identity provider/MNO and service provider (ISP/WISP/PLTE)
// A non-error return will indicate successful contract establishment and will
// result in the beginning of service for the user
func (s *BaseAccService) Start(ctx context.Context, req *protos.AcctSession) (*protos.AcctSessionResp, error) {
	provider, consumer, gw, err := RetrieveParticipants(ctx, req)
	if err != nil {
		return nil, err
	}
	client, outCtx, cancel, err := s.GetAcctClient()
	if err != nil {
		return nil, err
	}
	defer cancel()

	anSession := ToAnAcctSession(req)
	anSession.ConsumerId = consumer
	anSession.ProviderId = provider
	anSession.ProviderGatewayId = gw

	resp, err := client.Start(outCtx, anSession)
	if err != nil {
		return nil, err
	}
	return &protos.AcctSessionResp{
		ReportingAdvisory: &protos.AcctSessionResp_ReportLimits{
			OctetsIn:       resp.GetReportingAdvisory().GetOctetsIn(),
			OctetsOut:      resp.GetReportingAdvisory().GetOctetsOut(),
			ElapsedTimeSec: resp.GetReportingAdvisory().GetElapsedTimeSec(),
		},
		MinAcceptableQos: &protos.AcctQoS{
			DownloadMbps: resp.GetMinAcceptableQos().GetDownloadMbps(),
			UploadMbps:   resp.GetMinAcceptableQos().GetUploadMbps(),
		},
	}, nil
}

// Update should be continuously called for every ongoing service session to update
// the user bandwidth usage as well as current quality of provided service.
// If update returns error the session should be terminated and the user disconnected,
// In the case of unsuccessful update completion, service provider is suppose to follow up
// with final Stop call
func (s *BaseAccService) Update(ctx context.Context, req *protos.AcctUpdateReq) (*protos.AcctSessionResp, error) {
	provider, consumer, gw, err := RetrieveParticipants(ctx, req.GetSession())
	if err != nil {
		return nil, err
	}
	client, outCtx, cancel, err := s.GetAcctClient()
	if err != nil {
		return nil, err
	}
	defer cancel()
	anSession := ToAnAcctSession(req.GetSession())
	anSession.ConsumerId = consumer
	anSession.ProviderId = provider
	anSession.ProviderGatewayId = gw
	anReq := &anpb.UpdateReq{
		Session:     anSession,
		OctetsIn:    req.GetOctetsIn(),
		OctetsOut:   req.GetOctetsOut(),
		SessionTime: req.GetSessionTime(),
	}
	resp, err := client.Update(outCtx, anReq)
	if err != nil {
		return nil, err
	}
	return &protos.AcctSessionResp{
		ReportingAdvisory: &protos.AcctSessionResp_ReportLimits{
			OctetsIn:       resp.GetReportingAdvisory().GetOctetsIn(),
			OctetsOut:      resp.GetReportingAdvisory().GetOctetsOut(),
			ElapsedTimeSec: resp.GetReportingAdvisory().GetElapsedTimeSec(),
		},
		MinAcceptableQos: &protos.AcctQoS{
			DownloadMbps: resp.GetMinAcceptableQos().GetDownloadMbps(),
			UploadMbps:   resp.GetMinAcceptableQos().GetUploadMbps(),
		},
	}, nil
}

// Stop is a notification call to communicate to identity provider
// user/network  initiated service termination.
// stop will provide final used bandwidth count. stop call is issued
// after the user session was terminated.
func (s *BaseAccService) Stop(ctx context.Context, req *protos.AcctUpdateReq) (*protos.AcctStopResp, error) {
	provider, consumer, gw, err := RetrieveParticipants(ctx, req.GetSession())
	if err != nil {
		return nil, err
	}
	client, outCtx, cancel, err := s.GetAcctClient()
	if err != nil {
		return nil, err
	}
	defer cancel()
	anSession := ToAnAcctSession(req.GetSession())
	anSession.ConsumerId = consumer
	anSession.ProviderId = provider
	anSession.ProviderGatewayId = gw
	anReq := &anpb.UpdateReq{
		Session:     anSession,
		OctetsIn:    req.GetOctetsIn(),
		OctetsOut:   req.GetOctetsOut(),
		SessionTime: req.GetSessionTime(),
	}
	_, err = client.Stop(outCtx, anReq)
	return &protos.AcctStopResp{}, err
}

func ToAnAcctSession(session *protos.AcctSession) *anpb.Session {
	res := &anpb.Session{
		SessionId:   session.GetSessionId(),
		ProviderApn: session.GetServingApn(),
	}
	switch usr := session.GetUser().(type) {
	case *protos.AcctSession_IMSI:
		res.User = &anpb.Session_IMSI{IMSI: usr.IMSI}
	case *protos.AcctSession_CertificateSerialNumber:
		res.User = &anpb.Session_CertificateSerialNumber{CertificateSerialNumber: usr.CertificateSerialNumber}
	case *protos.AcctSession_Name:
		res.User = &anpb.Session_Name{Name: usr.Name}
	case *protos.AcctSession_HardwareAddr:
		res.User = &anpb.Session_HardwareAddr{HardwareAddr: usr.HardwareAddr}
	}
	res.AvailableQos = &anpb.QoS{
		DownloadMbps: session.GetAvailableQos().GetDownloadMbps(),
		UploadMbps:   session.GetAvailableQos().GetUploadMbps(),
	}
	return res
}
