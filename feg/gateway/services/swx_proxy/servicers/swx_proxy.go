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

// Package servicers implements Swx GRPC proxy service which sends MAR/SAR messages over
// diameter connection, waits (blocks) for diameter's MAA/SAAs returns their RPC representation
package servicers

import (
	"context"
	"fmt"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/hlr_proxy"
	"magma/feg/gateway/services/swx_proxy/cache"
	"magma/feg/gateway/services/swx_proxy/metrics"
	orcprotos "magma/orc8r/lib/go/protos"
)

const (
	TIMEOUT_SECONDS  = 10
	MAX_DIAM_RETRIES = 1
)

type Relay interface {
	RelayRTR(*RTR) (protos.ErrorCode, error)
	RelayASR(*diameter.ASR) (protos.ErrorCode, error)
}

type swxProxy struct {
	config         *SwxProxyConfig
	smClient       *sm.Client
	connMan        *diameter.ConnectionManager
	requestTracker *diameter.RequestTracker
	originStateID  uint32
	cache          *cache.Impl
	Relay          Relay
	healthTracker  *metrics.SwxHealthTracker
}

type SwxProxyConfig struct {
	ClientCfg             *diameter.DiameterClientConfig
	ServerCfg             *diameter.DiameterServerConfig
	VerifyAuthorization   bool // should we verify non-3gpp IP access is enabled for user
	RegisterOnAuth        bool // should we send SAR REGISTER on every MAR/A
	DeriveUnregisterRealm bool // use returned maa.AAAServerName to derive Origin Realm from
	CacheTTLSeconds       uint32
	HlrPlmnIds            plmn_filter.PlmnIdVals
}

// NewSwxProxy creates a new instance of the proxy with configured cache TTL
func NewSwxProxy(config *SwxProxyConfig) (*swxProxy, error) {
	cache := createCache(config)
	return NewSwxProxyWithCache(config, cache)
}

// NewSwxProxyWithCache creates a new instance of the proxy with given cache implementation
func NewSwxProxyWithCache(config *SwxProxyConfig, cache *cache.Impl) (*swxProxy, error) {
	err := ValidateSwxProxyConfig(config)
	if err != nil {
		return nil, err
	}
	config.ClientCfg = config.ClientCfg.FillInDefaults()

	originStateID := uint32(time.Now().Unix())

	mux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(config.ClientCfg.Host),
		OriginRealm:      datatype.DiameterIdentity(config.ClientCfg.Realm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(config.ClientCfg.ProductName),
		OriginStateID:    datatype.Unsigned32(originStateID),
		FirmwareRevision: 1,
	})

	mux.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		if m != nil {
			glog.Infof("Unhandled SWx message: %s", m)
		}
	}) // Catch all.

	if config.ClientCfg.WatchdogInterval == 0 {
		config.ClientCfg.WatchdogInterval = diameter.DefaultWatchdogIntervalSeconds
	}

	smClient := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     config.ClientCfg.Retransmits,
		RetransmitInterval: time.Second,
		EnableWatchdog:     config.ClientCfg.WatchdogInterval > 0,
		WatchdogInterval:   time.Second * time.Duration(config.ClientCfg.WatchdogInterval),
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_SWX_APP_ID)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
				},
			}),
		},
	}

	connMan := diameter.NewConnectionManager()
	// create connection in connection map
	connMan.GetConnection(smClient, config.ServerCfg)

	proxy := &swxProxy{
		config:         config,
		smClient:       smClient,
		connMan:        connMan,
		healthTracker:  metrics.NewSwxHealthTracker(),
		requestTracker: diameter.NewRequestTracker(),
		originStateID:  originStateID,
		cache:          cache,
		Relay:          &fegRelayClient{registry: registry.Get()},
	}
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_SWX_APP_ID, Code: diam.MultimediaAuthentication, Request: false},
		handleMAA(proxy))
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_SWX_APP_ID, Code: diam.ServerAssignment, Request: false},
		handleSAA(proxy))
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_SWX_APP_ID, Code: diam.RegistrationTermination, Request: true},
		handleRTR(proxy))

	return proxy, nil
}

// SwxProxyServer implementation
//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s *swxProxy) Authenticate(
	ctx context.Context,
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {
	var (
		res *protos.AuthenticationAnswer
		err error
	)
	authStartTime := time.Now()
	if s.IsHlrClient(req.GetUserName()) {
		res, err = hlr_proxy.Authenticate(ctx, req)
	} else {
		res, err = s.AuthenticateImpl(req)
	}
	if err == nil {
		metrics.AuthLatency.Observe(time.Since(authStartTime).Seconds())
	}
	return res, err
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s *swxProxy) Register(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	var (
		res *protos.RegistrationAnswer
		err error
	)
	registerStartTime := time.Now()
	if s.IsHlrClient(req.GetUserName()) {
		res, err = hlr_proxy.Register(ctx, req)
	} else {
		res, err = s.RegisterImpl(req, ServerAssignmentType_REGISTRATION)
	}
	if err == nil {
		metrics.RegisterLatency.Observe(time.Since(registerStartTime).Seconds())
	}
	return res, err
}

// Deregister sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s *swxProxy) Deregister(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	var (
		res *protos.RegistrationAnswer
		err error
	)
	deregisterStartTime := time.Now()
	if s.IsHlrClient(req.GetUserName()) {
		res, err = hlr_proxy.Register(ctx, req)
	} else {
		res, err = s.RegisterImpl(req, ServerAssignnmentType_USER_DEREGISTRATION)
	}
	if err == nil {
		metrics.DeregisterLatency.Observe(time.Since(deregisterStartTime).Seconds())
	}
	return res, err
}

// Disable closes all existing diameter connections and disables
// connection creation for the time specified in the request
func (s *swxProxy) Disable(ctx context.Context, req *protos.DisableMessage) (*orcprotos.Void, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil Disable Request")
	}
	s.cache.ClearAll()
	s.connMan.DisableFor(time.Duration(req.DisablePeriodSecs) * time.Second)
	return &orcprotos.Void{}, nil
}

// Enable enables diameter connection creation and gets a connection to the
// diameter server. If creation is already enabled and a connection already
// exists, Enable has no effect
func (s *swxProxy) Enable(ctx context.Context, req *orcprotos.Void) (*orcprotos.Void, error) {
	s.connMan.Enable()
	_, err := s.connMan.GetConnection(s.smClient, s.config.ServerCfg)
	return &orcprotos.Void{}, err
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
func (s *swxProxy) GetHealthStatus(ctx context.Context, req *orcprotos.Void) (*protos.HealthStatus, error) {
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
	reqTotal := deltaMetrics.MarTotal + deltaMetrics.SarTotal +
		deltaMetrics.MarSendFailures + deltaMetrics.SarSendFailures
	failureTotal := deltaMetrics.MarSendFailures + deltaMetrics.SarSendFailures +
		deltaMetrics.Timeouts + deltaMetrics.UnparseableMsg

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

func (s *swxProxy) genSID(imsi string) string {
	return s.config.ClientCfg.GenSessionIdImsi("swx", imsi)
}

// CreateCache creates a cache initialized with SWx config parameters
func createCache(config *SwxProxyConfig) *cache.Impl {
	fixConfigCacheMinTTL(config)
	cch, _ := cache.NewExt(cache.DefaultGcInterval, time.Second*time.Duration(config.CacheTTLSeconds))
	return cch
}

// fixConfigCacheMinTTL changes CacheTTLSeconds on the config if this is smaller than the default value
func fixConfigCacheMinTTL(config *SwxProxyConfig) {
	if config.CacheTTLSeconds < uint32(cache.DefaultGcInterval.Seconds()) {
		config.CacheTTLSeconds = uint32(cache.DefaultTtl.Seconds())
	}
}
