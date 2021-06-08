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

// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
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

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/services/s6a_proxy/metrics"
	orcprotos "magma/orc8r/lib/go/protos"
)

// Flag definitions
type FlagBit int

const (
	EmptyFlagBit FlagBit = 0
	FlagBit1     FlagBit = 1 << 1
	FlagBit5     FlagBit = 1 << 5
	FlagBit8     FlagBit = 1 << 8
	FlagBit9     FlagBit = 1 << 9
	FlagBit27    FlagBit = 1 << 27
)

const (
	ULR_RAT_TYPE     = 1004
	ULR_FLAGS        = FlagBit1 | FlagBit5 // 29.272 Table 7.3.7/1: ULR-Flags S6a/S6d-Indicator (bit 1), and Initial-AttachIndicator (bit 5)
	TIMEOUT_SECONDS  = 10
	MAX_DIAM_RETRIES = 1
)

type S6aProxyConfig struct {
	ClientCfg *diameter.DiameterClientConfig
	ServerCfg *diameter.DiameterServerConfig
	PlmnIds   plmn_filter.PlmnIdVals
}

type s6aProxy struct {
	config         *S6aProxyConfig
	smClient       *sm.Client
	connMan        *diameter.ConnectionManager
	requestTracker *diameter.RequestTracker
	healthTracker  *metrics.S6aHealthTracker
	originStateID  uint32
}

func NewS6aProxy(
	config *S6aProxyConfig,
) (*s6aProxy, error) {
	if config == nil {
		return nil, fmt.Errorf("S6aProxyConfig is nil")
	}
	clientCfg, serverCfg := config.ClientCfg, config.ServerCfg

	err := clientCfg.Validate()
	if err != nil {
		return nil, err
	}
	clientCfg = clientCfg.FillInDefaults()

	err = serverCfg.Validate()
	if err != nil {
		return nil, err
	}
	originStateID := uint32(time.Now().Unix())

	mux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(clientCfg.Host),
		OriginRealm:      datatype.DiameterIdentity(clientCfg.Realm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(clientCfg.ProductName),
		OriginStateID:    datatype.Unsigned32(originStateID),
		FirmwareRevision: 1,
	})

	mux.HandleFunc("ALL", func(diam.Conn, *diam.Message) {}) // Catch all.

	if clientCfg.WatchdogInterval == 0 {
		clientCfg.WatchdogInterval = diameter.DefaultWatchdogIntervalSeconds
	}

	smClient := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     clientCfg.Retransmits,
		RetransmitInterval: time.Second,
		EnableWatchdog:     clientCfg.WatchdogInterval > 0,
		WatchdogInterval:   time.Second * time.Duration(clientCfg.WatchdogInterval),
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_S6A_APP_ID)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
				},
			}),
		},
	}

	connMan := diameter.NewConnectionManager()
	// create connection in connection map
	connMan.GetConnection(smClient, serverCfg)

	proxy := &s6aProxy{
		config:         config,
		smClient:       smClient,
		connMan:        connMan,
		requestTracker: diameter.NewRequestTracker(),
		healthTracker:  metrics.NewS6aHealthTracker(),
		originStateID:  originStateID,
	}
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: false},
		handleAIA(proxy))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.UpdateLocation, Request: false},
		handleULA(proxy))

	mux.HandleIdx(diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.CancelLocation, Request: true},
		handleCLR(proxy))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.PurgeUE, Request: false},
		handlePUA(proxy))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.Reset, Request: true},
		handleRSR(proxy))

	return proxy, nil
}

// S6AProxyServer implementation
//
// AuthenticationInformation sends AIR over diameter connection,
// waits (blocks) for AIA & returns its RPC representation
func (s *s6aProxy) AuthenticationInformation(
	ctx context.Context, req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error,
) {
	airStartTime := time.Now()
	res, err := s.AuthenticationInformationImpl(req)
	if err == nil {
		metrics.AIRLatency.Observe(float64(time.Since(airStartTime)) / float64(time.Millisecond))
	}
	metrics.UpdateS6aRecentRequestMetrics(err)
	return res, err
}

// UpdateLocation sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULA & returns its RPC representation
func (s *s6aProxy) UpdateLocation(
	ctx context.Context, req *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error,
) {
	ulrStartTime := time.Now()
	res, err := s.UpdateLocationImpl(req)
	if err == nil {
		metrics.ULRLatency.Observe(float64(time.Since(ulrStartTime)) / float64(time.Millisecond))
	}
	metrics.UpdateS6aRecentRequestMetrics(err)
	return res, err
}

// PurgeUE sends PUR (Code 321) over diameter connection,
// waits (blocks) for PUA & returns its RPC representation
func (s *s6aProxy) PurgeUE(ctx context.Context, req *protos.PurgeUERequest) (*protos.PurgeUEAnswer, error) {
	res, err := s.PurgeUEImpl(req)
	metrics.UpdateS6aRecentRequestMetrics(err)
	return res, err
}

// Disable closes all existing diameter connections and disables
// connection creation for the time specified in the request
func (s *s6aProxy) Disable(ctx context.Context, req *protos.DisableMessage) (*orcprotos.Void, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil Disable Request")
	}
	s.connMan.DisableFor(time.Duration(req.DisablePeriodSecs) * time.Second)
	return &orcprotos.Void{}, nil
}

// Enable enables diameter connection creation and gets a connection to the
// diameter server. If creation is already enabled and a connection already
// exists, Enable has no effect
func (s *s6aProxy) Enable(ctx context.Context, req *orcprotos.Void) (*orcprotos.Void, error) {
	s.connMan.Enable()
	_, err := s.connMan.GetConnection(s.smClient, s.config.ServerCfg)
	return &orcprotos.Void{}, err
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
func (s *s6aProxy) GetHealthStatus(ctx context.Context, req *orcprotos.Void) (*protos.HealthStatus, error) {
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
	reqTotal := deltaMetrics.AirTotal + deltaMetrics.UlrTotal +
		deltaMetrics.AirSendFailures + deltaMetrics.UlrSendFailures
	failureTotal := deltaMetrics.AirSendFailures + deltaMetrics.UlrSendFailures +
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

func (s *s6aProxy) genSID() string {
	return s.config.ClientCfg.GenSessionID("s6a")
}
