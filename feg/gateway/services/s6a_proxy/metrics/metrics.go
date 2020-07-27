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
	"time"

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
	AIRRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "air_requests_total",
		Help: "Total number of AIR requests sent to HSS",
	})
	AIRSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "air_send_failures_total",
		Help: "Total number of AIR requests that failed to send to HSS",
	})
	ULRRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ulr_requests_total",
		Help: "Total number of ULR requests sent to HSS",
	})
	ULRSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ulr_send_failures_total",
		Help: "Total number of ULR requests that failed to send to HSS",
	})
	S6aTimeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s6a_timeouts_total",
		Help: "Total number of s6a timeouts",
	})
	S6aUnparseableMsg = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s6a_unparseable_msg_total",
		Help: "Total number of s6a messages received that cannot be parsed",
	})
	S6aResultCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "s6a_result_codes",
			Help: "S6a accumulated result codes",
		},
		[]string{"code"},
	)
	S6aSuccessTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "s6a_success_timestamp",
		Help: "Timestamp of the last successfully completed s6a request",
	})
	S6aFailuresSinceLastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "s6a_failures_since_last_success",
		Help: "The total number of s6a request failures since the last successful request completed",
	})
	// Latencies
	AIRLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "air_lat",
		Help:       "Latency of AIR requests (milliseconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	ULRLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "ulr_lat",
		Help:       "Latency of ULR requests (milliseconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

type S6aHealthMetrics struct {
	AirTotal        int64
	AirSendFailures int64
	UlrTotal        int64
	UlrSendFailures int64
	Timeouts        int64
	UnparseableMsg  int64
}

type S6aHealthTracker struct {
	Metrics                 *S6aHealthMetrics
	RequestFailureThreshold float32
	MinimumRequestThreshold uint32
}

func NewS6aHealthTracker() *S6aHealthTracker {
	initMetrics := &S6aHealthMetrics{
		AirTotal:        0,
		AirSendFailures: 0,
		UlrSendFailures: 0,
		UlrTotal:        0,
		Timeouts:        0,
		UnparseableMsg:  0,
	}
	defaultHealthTracker := &S6aHealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: float32(DefaultRequestFailureThreshold),
		MinimumRequestThreshold: uint32(DefaultMinimumRequiredRequests),
	}
	s6aCfg := &mconfigprotos.S6AConfig{}
	err := mconfig.GetServiceConfigs(strings.ToLower(registry.S6A_PROXY), s6aCfg)
	if err != nil {
		return defaultHealthTracker
	}
	reqFailureThreshold := s6aCfg.GetRequestFailureThreshold()
	minReqThreshold := s6aCfg.GetMinimumRequestThreshold()
	if reqFailureThreshold == 0 {
		glog.Info("Request failure threshold cannot be 0; Using default health parameters...")
		return defaultHealthTracker
	}
	if minReqThreshold == 0 {
		glog.Info("Minimum request threshold cannot be 0; Using default health parameters...")
		return defaultHealthTracker
	}
	return &S6aHealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: reqFailureThreshold,
		MinimumRequestThreshold: minReqThreshold,
	}
}

func init() {
	prometheus.MustRegister(AIRRequests, AIRSendFailures, ULRRequests,
		ULRSendFailures, S6aTimeouts, S6aUnparseableMsg, S6aResultCodes,
		S6aSuccessTimestamp, S6aFailuresSinceLastSuccess, AIRLatency, ULRLatency)
}

func UpdateS6aRecentRequestMetrics(err error) {
	if err == nil {
		S6aSuccessTimestamp.Set(float64(time.Now().Unix()))
		S6aFailuresSinceLastSuccess.Set(0)
	} else {
		S6aFailuresSinceLastSuccess.Inc()
	}
}

func GetCurrentHealthMetrics() (*S6aHealthMetrics, error) {
	airTotal, err := service_health_metrics.GetInt64("air_requests_total")
	if err != nil {
		return nil, err
	}
	airSendFailures, err := service_health_metrics.GetInt64("air_send_failures_total")
	if err != nil {
		return nil, err
	}
	ulrTotal, err := service_health_metrics.GetInt64("ulr_requests_total")
	if err != nil {
		return nil, err
	}
	ulrSendFailures, err := service_health_metrics.GetInt64("ulr_send_failures_total")
	if err != nil {
		return nil, err
	}
	timeouts, err := service_health_metrics.GetInt64("s6a_timeouts_total")
	if err != nil {
		return nil, err
	}
	unparseable, err := service_health_metrics.GetInt64("s6a_unparseable_msg_total")
	if err != nil {
		return nil, err
	}
	return &S6aHealthMetrics{
		AirTotal:        airTotal,
		AirSendFailures: airSendFailures,
		UlrTotal:        ulrTotal,
		UlrSendFailures: ulrSendFailures,
		Timeouts:        timeouts,
		UnparseableMsg:  unparseable,
	}, nil
}

func (prevMetrics *S6aHealthMetrics) GetDelta(currentMetrics *S6aHealthMetrics) (*S6aHealthMetrics, error) {
	if currentMetrics == nil {
		return nil, fmt.Errorf("Nil current s6aHealthMetrics struct provided")
	}
	deltaMetrics := &S6aHealthMetrics{
		AirTotal:        currentMetrics.AirTotal - prevMetrics.AirTotal,
		AirSendFailures: currentMetrics.AirSendFailures - prevMetrics.AirSendFailures,
		UlrTotal:        currentMetrics.UlrTotal - prevMetrics.UlrTotal,
		UlrSendFailures: currentMetrics.UlrSendFailures - prevMetrics.UlrSendFailures,
		Timeouts:        currentMetrics.Timeouts - prevMetrics.Timeouts,
		UnparseableMsg:  currentMetrics.UnparseableMsg - prevMetrics.UnparseableMsg,
	}
	// Update stored counts to current metric totals
	*prevMetrics = *currentMetrics
	return deltaMetrics, nil
}
