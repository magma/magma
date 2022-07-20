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
	"time"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"
	"magma/feg/gateway/services/s8_proxy/metrics"
	orc8r_protos "magma/orc8r/lib/go/protos"
)

type S8Proxy struct {
	config        *S8ProxyConfig
	gtpClient     *gtp.Client
	healthTracker *metrics.S8HealthTracker
}

type S8ProxyConfig struct {
	GtpTimeout        time.Duration
	ClientAddr        string
	ServerAddr        *net.UDPAddr
	ApnOperatorSuffix string
}

// NewS8Proxy creates an s8 proxy, but does not checks the PGW is alive
func NewS8Proxy(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewRunningClient(
		context.Background(), config.ClientAddr,
		gtp.SGWControlPlaneIfType, config.GtpTimeout)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	return newS8ProxyImp(gtpCli, config)
}

// NewS8ProxyWithEcho creates an s8 proxy already connected to a server (checks with echo if PGW is alive)
// Used mainly for testing with s8_cli
func NewS8ProxyWithEcho(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewConnectedAutoClient(
		context.Background(), config.ServerAddr.String(),
		gtp.SGWControlPlaneIfType, config.GtpTimeout)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	return newS8ProxyImp(gtpCli, config)
}

func newS8ProxyImp(cli *gtp.Client, config *S8ProxyConfig) (*S8Proxy, error) {
	// TODO: validate config
	s8p := &S8Proxy{
		config:        config,
		gtpClient:     cli,
		healthTracker: metrics.NewS8HealthTracker(),
	}
	addS8GtpHandlers(s8p)
	return s8p, nil
}

func (s *S8Proxy) CreateSession(ctx context.Context, req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	metrics.SessionCreateRequests.Inc()
	err := validateCreateSessionRequest(req)
	if err != nil {
		metrics.SessionCreateFails.Inc()
		err = fmt.Errorf("Create Session failed for IMSI %s:, couldn't validate request: %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}

	cPgwUDPAddr, err := s.configOrRequestedPgwAddress(req.PgwAddrs)
	if err != nil {
		metrics.SessionCreateFails.Inc()
		err = fmt.Errorf("Create Session failed for IMSI %s: %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}
	// build csReq IE message
	csReqMsg, err := buildCreateSessionRequestMsg(cPgwUDPAddr, s.config.ApnOperatorSuffix, req)
	if err != nil {
		metrics.SessionCreateFails.Inc()
		err = fmt.Errorf("Create Session failed to build IEs for IMSI %s: %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}

	// send, register and receive create session (session is created on the gtp client during this process too)
	csRes, err := s.sendAndReceiveCreateSession(req, cPgwUDPAddr, csReqMsg)
	if err != nil {
		metrics.SessionCreateFails.Inc()
		err = fmt.Errorf("Create Session failed for IMSI %s: %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}
	return csRes, nil
}

func (s *S8Proxy) DeleteSession(ctx context.Context, req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	metrics.SessionDeleteRequests.Inc()
	err := validateDeleteSessionRequest(req)
	if err != nil {
		metrics.SessionDeleteFails.Inc()
		err = fmt.Errorf("Delete Session failed for IMSI %s:, couldn't validate request: %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}
	cPgwUDPAddr, err := s.configOrRequestedPgwAddress(req.PgwAddrs)
	if err != nil {
		metrics.SessionDeleteFails.Inc()
		err = fmt.Errorf("Delete Session failed for IMSI %s: %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}
	dsReqMsg, err := buildDeleteSessionRequestMsg(cPgwUDPAddr, req)
	if err != nil {
		metrics.SessionDeleteFails.Inc()
		err = fmt.Errorf("Delete Session failed to build IEs for IMSI %s: %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}

	cdRes, err := s.sendAndReceiveDeleteSession(req, cPgwUDPAddr, dsReqMsg)
	if err != nil {
		metrics.SessionDeleteFails.Inc()
		err = fmt.Errorf("Delete Session failed for IMSI %s:, %s", req.Imsi, err)
		glog.Error(err)
		return nil, err
	}
	return cdRes, nil
}

func (s *S8Proxy) SendEcho(_ context.Context, req *protos.EchoRequest) (*protos.EchoResponse, error) {
	cPgwUDPAddr, err := s.configOrRequestedPgwAddress(req.PgwAddrs)
	if err != nil {
		err = fmt.Errorf("SendEcho to %s failed: %s", cPgwUDPAddr, err)
		glog.Error(err)
		return nil, err
	}
	err = s.gtpClient.SendEchoRequest(cPgwUDPAddr)
	if err != nil {
		return nil, err
	}
	return &protos.EchoResponse{}, nil
}

func (s *S8Proxy) CreateBearerResponse(_ context.Context, res *protos.CreateBearerResponsePgw) (*orc8r_protos.Void, error) {
	metrics.BearerCreateRequests.Inc()
	cPgwUDPAddr := ParseAddress(res.PgwAddrs)
	if cPgwUDPAddr == nil {
		metrics.BearerCreateFails.Inc()
		err := fmt.Errorf("CreateBearerResponse to %s failed: couldnt paarse address", res.PgwAddrs)
		glog.Error(err)
		return nil, err
	}

	cbResMsg, err := buildCreateBearerResMsg(res)
	if err != nil {
		metrics.BearerCreateFails.Inc()
		return nil, err
	}

	_, err = s.sendAndReceiveCreateBearerResponse(res, cPgwUDPAddr, cbResMsg)
	if err != nil {
		metrics.BearerCreateFails.Inc()
		err = fmt.Errorf("Create Bearer Response failed for IMSI %s:, %s", res.Imsi, err)
		glog.Error(err)
		return nil, err
	}

	return &orc8r_protos.Void{}, nil
}

func (s *S8Proxy) DeleteBearerResponse(_ context.Context, res *protos.DeleteBearerResponsePgw) (*orc8r_protos.Void, error) {
	metrics.BearerDeleteRequests.Inc()
	cPgwUDPAddr := ParseAddress(res.PgwAddrs)
	if cPgwUDPAddr == nil {
		metrics.BearerDeleteFails.Inc()
		err := fmt.Errorf("DeleteBearerResponse to %s failed: couldnt paarse address", res.PgwAddrs)
		glog.Error(err)
		return nil, err
	}

	dbResMsg, err := buildDeleteBearerResMsg(res)
	if err != nil {
		metrics.BearerDeleteFails.Inc()
		return nil, err
	}

	_, err = s.sendAndReceiveDeleteBearerResponse(res, cPgwUDPAddr, dbResMsg)
	if err != nil {
		metrics.BearerDeleteFails.Inc()
		err = fmt.Errorf("Create Bearer Response failed for IMSI %s:, %s", res.Imsi, err)
		glog.Error(err)
		return nil, err
	}

	return &orc8r_protos.Void{}, nil
}

// configOrRequestedPgwAddress returns an UDPAddrs if the passed string corresponds to a valid ip,
// otherwise it uses the server address configured on s8_proxy
func (s *S8Proxy) configOrRequestedPgwAddress(pgwAddrsFromRequest string) (*net.UDPAddr, error) {
	addrs := ParseAddress(pgwAddrsFromRequest)
	if addrs != nil {
		// address coming from string has precedence
		return addrs, nil
	}
	if s.config.ServerAddr != nil {
		return s.config.ServerAddr, nil
	}
	return nil, fmt.Errorf("Neither the request nor s8_proxy has a valid server (pgw) address")
}

func validateCreateSessionRequest(csr *protos.CreateSessionRequestPgw) error {
	if csr.BearerContext == nil || csr.BearerContext.UserPlaneFteid == nil || csr.BearerContext.Id == 0 ||
		csr.BearerContext.Qos == nil || csr.Uli == nil || csr.ServingNetwork == nil {
		return fmt.Errorf("CreateSessionRequest missing fields %+v", csr)
	}
	return nil
}

func validateDeleteSessionRequest(dsr *protos.DeleteSessionRequestPgw) error {
	if dsr.Imsi == "" || dsr.Uli == nil || dsr.ServingNetwork == nil {
		return fmt.Errorf("DeleteSessionRequest missing fields %+v", dsr)
	}
	return nil
}

func (s *S8Proxy) Disable(ctx context.Context, req *protos.DisableMessage) (*orc8r_protos.Void, error) {
	return &orc8r_protos.Void{}, nil
}

func (s *S8Proxy) Enable(ctx context.Context, req *orc8r_protos.Void) (*orc8r_protos.Void, error) {
	return &orc8r_protos.Void{}, nil
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
func (s *S8Proxy) GetHealthStatus(ctx context.Context, req *orc8r_protos.Void) (*protos.HealthStatus, error) {
	currentMetrics, err := metrics.GetCurrentHealthMetrics()
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("Error occurred while retrieving health metrics: %s", err),
		}, err
	}
	deltaMetrics, err := s.healthTracker.Metrics.GetDelta(currentMetrics)
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: err.Error(),
		}, err
	}

	reqTotal := deltaMetrics.SessionCreateRequests + deltaMetrics.SessionDeleteRequests +
		deltaMetrics.BearerCreateRequests + deltaMetrics.BearerDeleteRequests
	failureTotal := deltaMetrics.SessionCreateFails +
		deltaMetrics.SessionCreateFails +
		deltaMetrics.BearerCreateFails + deltaMetrics.BearerDeleteFails

	exceedsThreshold := reqTotal >= int64(s.healthTracker.MinimumRequestThreshold) &&
		float32(failureTotal)/float32(reqTotal) >= s.healthTracker.RequestFailureThreshold
	if exceedsThreshold {
		unhealthyMsg := fmt.Sprintf("Metric Request Failure Ratio >= threshold %f; %d / %d",
			s.healthTracker.RequestFailureThreshold,
			failureTotal,
			reqTotal,
		)
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: unhealthyMsg,
		}, nil
	}
	return &protos.HealthStatus{
		Health:        protos.HealthStatus_HEALTHY,
		HealthMessage: "All metrics appear healthy",
	}, nil
}
