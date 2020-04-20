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

package metrics

import (
	"fmt"
	"strings"

	mconfigprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	service_health_metrics "magma/feg/gateway/service_health/metrics"
	"magma/gateway/mconfig"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const DefaultRequestFailureThreshold = 0.50
const DefaultMinimumRequiredRequests = 1

// Prometheus counters are monotonically increasing
// Counters reset to zero on service restart
var (
	MARRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mar_requests_total",
		Help: "Total number of MAR requests sent to HSS",
	})
	MARSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mar_send_failures_total",
		Help: "Total number of MAR requests that failed to send to HSS",
	})
	SARRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sar_requests_total",
		Help: "Total number of SAR requests sent to HSS",
	})
	SARSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sar_send_failures_total",
		Help: "Total number of SAR requests that failed to send to HSS",
	})
	SwxTimeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_timeouts_total",
		Help: "Total number of swx timeouts",
	})
	SwxUnparseableMsg = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_unparseable_msg_total",
		Help: "Total number of swx messages received that cannot be parsed",
	})
	SwxInvalidSessions = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_invalid_sessions_total",
		Help: "Total number of swx responses received with invalid sids",
	})
	SwxResultCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "swx_result_codes",
			Help: "swx accumulated result codes",
		},
		[]string{"code"},
	)
	SwxExperimentalResultCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "swx_experimental_result_codes",
			Help: "swx accumulated experimental result codes",
		},
		[]string{"code"},
	)
	UnauthorizedAuthAttempts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "unauthorized_auth_requests_total",
		Help: "Total number of authentication requests for un-authorized users",
	})

	// Latency Metrics
	MARLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "mar_latency",
		Help:       "Latency of MAR Diameter requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	SARLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "sar_latency",
		Help:       "Latency of SAR Diameter requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	AuthLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "auth_latency",
		Help:       "Latency of Authenticate GRPC requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	RegisterLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "register_latency",
		Help:       "Latency of Register GRPC requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	DeregisterLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "deregister_latency",
		Help:       "Latency of Deregister GRPC requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func init() {
	prometheus.MustRegister(MARRequests, MARSendFailures, SARRequests,
		SARSendFailures, SwxTimeouts, SwxUnparseableMsg, SwxInvalidSessions,
		SwxResultCodes, SwxExperimentalResultCodes, UnauthorizedAuthAttempts,
		MARLatency, SARLatency, AuthLatency, RegisterLatency, DeregisterLatency)
}

type SwxHealthMetrics struct {
	MarTotal        int64
	MarSendFailures int64
	SarTotal        int64
	SarSendFailures int64
	Timeouts        int64
	UnparseableMsg  int64
}

type SwxHealthTracker struct {
	Metrics                 *SwxHealthMetrics
	RequestFailureThreshold float32
	MinimumRequestThreshold uint32
}

func NewSwxHealthTracker() *SwxHealthTracker {
	initMetrics := &SwxHealthMetrics{
		MarTotal:        0,
		MarSendFailures: 0,
		SarTotal:        0,
		SarSendFailures: 0,
		Timeouts:        0,
		UnparseableMsg:  0,
	}
	defaultHealthTracker := &SwxHealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: float32(DefaultRequestFailureThreshold),
		MinimumRequestThreshold: uint32(DefaultMinimumRequiredRequests),
	}
	swxCfg := &mconfigprotos.SwxConfig{}
	err := mconfig.GetServiceConfigs(strings.ToLower(registry.SWX_PROXY), swxCfg)
	if err != nil {
		return defaultHealthTracker
	}
	reqFailureThreshold := swxCfg.GetRequestFailureThreshold()
	minReqThreshold := swxCfg.GetMinimumRequestThreshold()
	if reqFailureThreshold == 0 {
		glog.Info("Request failure threshold cannot be 0; Using default health parameters...")
		return defaultHealthTracker
	}
	if minReqThreshold == 0 {
		glog.Info("Minimum request threshold cannot be 0; Using default health parameters...")
		return defaultHealthTracker
	}
	return &SwxHealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: reqFailureThreshold,
		MinimumRequestThreshold: minReqThreshold,
	}
}

func GetCurrentHealthMetrics() (*SwxHealthMetrics, error) {
	marTotal, err := service_health_metrics.GetInt64("mar_requests_total")
	if err != nil {
		return nil, err
	}
	marSendFailures, err := service_health_metrics.GetInt64("mar_send_failures_total")
	if err != nil {
		return nil, err
	}
	sarTotal, err := service_health_metrics.GetInt64("sar_requests_total")
	if err != nil {
		return nil, err
	}
	sarSendFailures, err := service_health_metrics.GetInt64("sar_send_failures_total")
	if err != nil {
		return nil, err
	}
	timeouts, err := service_health_metrics.GetInt64("swx_timeouts_total")
	if err != nil {
		return nil, err
	}
	unparseable, err := service_health_metrics.GetInt64("swx_unparseable_msg_total")
	if err != nil {
		return nil, err
	}
	return &SwxHealthMetrics{
		MarTotal:        marTotal,
		MarSendFailures: marSendFailures,
		SarTotal:        sarTotal,
		SarSendFailures: sarSendFailures,
		Timeouts:        timeouts,
		UnparseableMsg:  unparseable,
	}, nil
}

func (prevMetrics *SwxHealthMetrics) GetDelta(currentMetrics *SwxHealthMetrics) (*SwxHealthMetrics, error) {
	if currentMetrics == nil {
		return nil, fmt.Errorf("Nil current swxHealthMetrics struct provided")
	}
	deltaMetrics := &SwxHealthMetrics{
		MarTotal:        currentMetrics.MarTotal - prevMetrics.MarTotal,
		MarSendFailures: currentMetrics.MarSendFailures - prevMetrics.MarSendFailures,
		SarTotal:        currentMetrics.SarTotal - prevMetrics.SarTotal,
		SarSendFailures: currentMetrics.SarSendFailures - prevMetrics.SarSendFailures,
		Timeouts:        currentMetrics.Timeouts - prevMetrics.Timeouts,
		UnparseableMsg:  currentMetrics.UnparseableMsg - prevMetrics.UnparseableMsg,
	}
	// Update stored counts to current metric totals
	*prevMetrics = *currentMetrics
	return deltaMetrics, nil
}
