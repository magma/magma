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

package metrics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"

	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	service_health_metrics "magma/feg/gateway/service_health/metrics"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/gateway/mconfig"
)

const (
	DefaultRequestFailureThresholdPct = 0.50
	DefaultMinimumRequestThreshold    = 1
)

// Prometheus counters are monotonically increasing
// Counters are reset to zero on service restarts
var (
	PcfSmPolicyCreateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcf_sm_policy_create_requests_total",
		Help: "Total number of SM Policy Create requests sent to PCF",
	})
	PcfSmPolicyCreateFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcf_sm_policy_create_failures_total",
		Help: "Total number of SM Policy Create that failed to send to PCF",
	})
	PcfSmPolicyUpdateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcf_sm_policy_update_requests_total",
		Help: "Total number of SM Policy Update sent to PCF",
	})
	PcfSmPolicyUpdateFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcf_sm_policy_update_failures_total",
		Help: "Total number of SM Policy Update that failed to send to PCF",
	})
	PcfSmPolicyDeleteRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcf_sm_policy_delete_requests_total",
		Help: "Total number of SM Policy Delete requests sent to PCF",
	})
	PcfSmPolicyDeleteFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcf_sm_policy_delete_failures_total",
		Help: "Total number of SM Policy Delete requests that failed to send to PCF",
	})

	N7Timeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "n7_timeouts_total",
		Help: "Total number of N7 timeouts",
	})

	N7SuccessTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "n7_success_timestamp",
		Help: "Timestamp of the last successfully completed N7 request",
	})
	N7FailuresSinceLastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "n7_failures_since_last_success",
		Help: "The total number of N7 request failures since the last successful request completed",
	})
)

type SessionHealthTracker struct {
	Metrics                 *SessionHealthMetrics
	RequestFailureThreshold float32
	MinimumRequestThreshold uint32
}

type SessionHealthMetrics struct {
	SmPolicyCreateTotal    int64
	SmpolicyCreateFailures int64
	SmPolicyUpdateTotal    int64
	SmPolicyUpdateFailures int64
	SmPolicyDeleteTotal    int64
	SmPolicyDeleteFailures int64
	N7Timeouts             int64
}

func init() {
	prometheus.MustRegister(
		PcfSmPolicyCreateRequests, PcfSmPolicyCreateFailures, PcfSmPolicyUpdateRequests,
		PcfSmPolicyUpdateFailures, PcfSmPolicyDeleteRequests, PcfSmPolicyDeleteFailures,
		N7Timeouts, N7SuccessTimestamp, N7FailuresSinceLastSuccess,
	)
}

func NewSessionHealthTracker() *SessionHealthTracker {
	initMetrics := &SessionHealthMetrics{
		SmPolicyCreateTotal:    0,
		SmpolicyCreateFailures: 0,
		SmPolicyUpdateTotal:    0,
		SmPolicyUpdateFailures: 0,
		SmPolicyDeleteTotal:    0,
		SmPolicyDeleteFailures: 0,
		N7Timeouts:             0,
	}
	defaultMetrics := &SessionHealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: float32(DefaultRequestFailureThresholdPct),
		MinimumRequestThreshold: uint32(DefaultMinimumRequestThreshold),
	}
	configPtr := &mcfgprotos.N7N40ProxyConfig{}
	err := mconfig.GetServiceConfigs(n7.N7N40ProxyServiceName, configPtr)
	if err != nil {
		return defaultMetrics
	}
	reqFailureThreshold := configPtr.GetRequestFailureThreshold()
	minReqThreshold := configPtr.GetMinimumRequestThreshold()
	if reqFailureThreshold == 0 {
		glog.Info("Request failure threshold cannot be 0; Using default health parameters...")
		return defaultMetrics
	}
	if minReqThreshold == 0 {
		glog.Info("Minimum request threshold cannot be 0; Using default health parameters...")
		return defaultMetrics
	}
	return &SessionHealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: reqFailureThreshold,
		MinimumRequestThreshold: minReqThreshold,
	}
}

func GetCurrentHealthMetrics() (*SessionHealthMetrics, error) {
	smPolicyCreateTotal, err := service_health_metrics.GetInt64("pcf_sm_policy_create_requests_total")
	if err != nil {
		return nil, err
	}
	smPolicyCreateFailures, err := service_health_metrics.GetInt64("pcf_sm_policy_create_failures_total")
	if err != nil {
		return nil, err
	}
	smPolicyUpdateTotal, err := service_health_metrics.GetInt64("pcf_sm_policy_update_requests_total")
	if err != nil {
		return nil, err
	}
	smPolicyUpdateFailures, err := service_health_metrics.GetInt64("pcf_sm_policy_update_failures_total")
	if err != nil {
		return nil, err
	}
	smPolicyDeleteTotal, err := service_health_metrics.GetInt64("pcf_sm_policy_delete_requests_total")
	if err != nil {
		return nil, err
	}
	smPolicyDeleteFailure, err := service_health_metrics.GetInt64("pcf_sm_policy_delete_failures_total")
	if err != nil {
		return nil, err
	}
	n7Timeouts, err := service_health_metrics.GetInt64("pcf_sm_policy_delete_failures_total")
	if err != nil {
		return nil, err
	}
	return &SessionHealthMetrics{
		SmPolicyCreateTotal:    smPolicyCreateTotal,
		SmpolicyCreateFailures: smPolicyCreateFailures,
		SmPolicyUpdateTotal:    smPolicyUpdateTotal,
		SmPolicyUpdateFailures: smPolicyUpdateFailures,
		SmPolicyDeleteTotal:    smPolicyDeleteTotal,
		SmPolicyDeleteFailures: smPolicyDeleteFailure,
		N7Timeouts:             n7Timeouts,
	}, nil
}

func (prevMetrics *SessionHealthMetrics) GetDelta(
	currentMetrics *SessionHealthMetrics,
) (*SessionHealthMetrics, error) {
	if currentMetrics == nil {
		return nil, fmt.Errorf("GetDelta has nil SessionHealthMetrics on N7N40Proxy")
	}
	deltaMetrics := &SessionHealthMetrics{
		SmPolicyCreateTotal:    currentMetrics.SmPolicyCreateTotal - prevMetrics.SmPolicyCreateTotal,
		SmpolicyCreateFailures: currentMetrics.SmpolicyCreateFailures - prevMetrics.SmpolicyCreateFailures,
		SmPolicyUpdateTotal:    currentMetrics.SmPolicyUpdateTotal - prevMetrics.SmPolicyUpdateTotal,
		SmPolicyUpdateFailures: currentMetrics.SmPolicyUpdateFailures - prevMetrics.SmPolicyUpdateFailures,
		SmPolicyDeleteTotal:    currentMetrics.SmPolicyDeleteTotal - prevMetrics.SmPolicyDeleteTotal,
		SmPolicyDeleteFailures: currentMetrics.SmPolicyDeleteFailures - prevMetrics.SmPolicyDeleteFailures,
		N7Timeouts:             currentMetrics.N7Timeouts - prevMetrics.N7Timeouts,
	}
	// Update stored counts to current metric totals
	*prevMetrics = *currentMetrics
	return deltaMetrics, nil
}

// Generic metric functions
func UpdateN7RecentRequestMetrics(err error) {
	if err == nil {
		N7SuccessTimestamp.Set(float64(time.Now().Unix()))
		N7FailuresSinceLastSuccess.Set(0)
	} else {
		N7FailuresSinceLastSuccess.Inc()
	}
}

func ReportCreateSmPolicy(err error) {
	UpdateN7RecentRequestMetrics(err)
	if err != nil {
		PcfSmPolicyCreateFailures.Inc()
		if errors.Is(err, context.DeadlineExceeded) {
			N7Timeouts.Inc()
		}
	}
	PcfSmPolicyCreateRequests.Inc()
}

func ReportUpdateSmPolicy(err error) {
	UpdateN7RecentRequestMetrics(err)
	if err != nil {
		PcfSmPolicyUpdateFailures.Inc()
		if errors.Is(err, context.DeadlineExceeded) {
			N7Timeouts.Inc()
		}
	}
	PcfSmPolicyUpdateRequests.Inc()
}

func ReportDeleteSmPolicy(err error) {
	UpdateN7RecentRequestMetrics(err)
	if err != nil {
		PcfSmPolicyDeleteFailures.Inc()
		if errors.Is(err, context.DeadlineExceeded) {
			N7Timeouts.Inc()
		}
	}
	PcfSmPolicyDeleteRequests.Inc()
}
