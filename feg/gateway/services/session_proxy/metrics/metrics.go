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

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	service_health_metrics "magma/feg/gateway/service_health/metrics"
	feg_mconfig "magma/gateway/mconfig"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const DefaultRequestFailureThresholdPct = 0.50
const DefaultMinimumRequestThreshold = 1

// Prometheus counters are monotonically increasing
// Counters are reset to zero on service restarts
var (
	PcrfCcrInitRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcrf_ccr_init_requests_total",
		Help: "Total number of CCR Init requests sent to PCRF",
	})
	PcrfCcrInitSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcrf_ccr_init_send_failures_total",
		Help: "Total number of CCR Init requests that failed to send to PCRF",
	})
	PcrfCcrUpdateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcrf_ccr_update_requests_total",
		Help: "Total number of CCR Update requests sent to PCRF",
	})
	PcrfCcrUpdateSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcrf_ccr_update_send_failures_total",
		Help: "Total number of CCR Update requests that failed to send to PCRF",
	})
	PcrfCcrTerminateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcrf_ccr_terminate_requests_total",
		Help: "Total number of CCR terminate requests sent to PCRF",
	})
	PcrfCcrTerminateSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pcrf_ccr_terminate_send_failures_total",
		Help: "Total number of CCR terminate requests that failed to send to PCRF",
	})

	OcsCcrInitRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocs_ccr_init_requests_total",
		Help: "Total number of CCR Init requests sent to OCS",
	})
	OcsCcrInitSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocs_ccr_init_send_failures_total",
		Help: "Total number of CCR Init requests that failed to send to OCS",
	})
	OcsCcrUpdateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocs_ccr_update_requests_total",
		Help: "Total number of CCR Update requests sent to OCS",
	})
	OcsCcrUpdateSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocs_ccr_update_send_failures_total",
		Help: "Total number of CCR Update requests that failed to send to OCS",
	})
	OcsCcrTerminateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocs_ccr_terminate_requests_total",
		Help: "Total number of CCR terminate requests sent to OCS",
	})
	OcsCcrTerminateSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocs_ccr_terminate_send_failures_total",
		Help: "Total number of CCR terminate requests that failed to send to OCS",
	})

	GxUnparseableMsg = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gx_unparseable_msg_total",
		Help: "Total number of gx messages received that cannot be parsed",
	})
	GyUnparseableMsg = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gy_unparseable_msg_total",
		Help: "Total number of gy messages received that cannot be parsed",
	})

	GxTimeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gx_timeouts_total",
		Help: "Total number of gx timeouts",
	})
	GyTimeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gy_timeouts_total",
		Help: "Total number of gy timeouts",
	})

	GxResultCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gx_result_codes",
			Help: "Gx result codes",
		},
		[]string{"code"},
	)
	GyResultCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gy_result_codes",
			Help: "Gy result codes",
		},
		[]string{"code"},
	)

	GxSuccessTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gx_success_timestamp",
		Help: "Timestamp of the last successfully completed gx request",
	})
	GxFailuresSinceLastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gx_failures_since_last_success",
		Help: "The total number of gx request failures since the last successful request completed",
	})

	GySuccessTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gy_success_timestamp",
		Help: "Timestamp of the last successfully completed gy request",
	})
	GyFailuresSinceLastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gy_failures_since_last_success",
		Help: "The total number of gy request failures since the last successful request completed",
	})
)

type SessionHealthTracker struct {
	Metrics                 *SessionHealthMetrics
	RequestFailureThreshold float32
	MinimumRequestThreshold uint32
}

type SessionHealthMetrics struct {
	PcrfInitTotal             int64
	PcrfInitSendFailures      int64
	PcrfUpdateTotal           int64
	PcrfUpdateSendFailures    int64
	PcrfTerminateTotal        int64
	PcrfTerminateSendFailures int64
	OcsInitTotal              int64
	OcsInitSendFailures       int64
	OcsUpdateTotal            int64
	OcsUpdateSendFailures     int64
	OcsTerminateTotal         int64
	OcsTerminateSendFailures  int64
	GxTimeouts                int64
	GxUnparseableMsg          int64
	GyTimeouts                int64
	GyUnparseableMsg          int64
}

func init() {
	prometheus.MustRegister(PcrfCcrInitRequests, PcrfCcrInitSendFailures, PcrfCcrUpdateRequests, PcrfCcrUpdateSendFailures,
		PcrfCcrTerminateRequests, PcrfCcrTerminateSendFailures, OcsCcrInitRequests, OcsCcrInitSendFailures,
		OcsCcrUpdateRequests, OcsCcrUpdateSendFailures, OcsCcrTerminateRequests, OcsCcrTerminateSendFailures,
		GxUnparseableMsg, GyUnparseableMsg, GxTimeouts, GyTimeouts, GxResultCodes, GyResultCodes,
		GxSuccessTimestamp, GxFailuresSinceLastSuccess, GySuccessTimestamp, GyFailuresSinceLastSuccess)
}

func NewSessionHealthTracker() *SessionHealthTracker {
	initMetrics := &SessionHealthMetrics{
		PcrfInitTotal:             0,
		PcrfInitSendFailures:      0,
		PcrfUpdateTotal:           0,
		PcrfUpdateSendFailures:    0,
		PcrfTerminateTotal:        0,
		PcrfTerminateSendFailures: 0,
		OcsInitTotal:              0,
		OcsInitSendFailures:       0,
		OcsUpdateTotal:            0,
		OcsUpdateSendFailures:     0,
		OcsTerminateTotal:         0,
		OcsTerminateSendFailures:  0,
		GxTimeouts:                0,
		GxUnparseableMsg:          0,
		GyTimeouts:                0,
		GyUnparseableMsg:          0,
	}
	defaultMetrics := &SessionHealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: float32(DefaultRequestFailureThresholdPct),
		MinimumRequestThreshold: uint32(DefaultMinimumRequestThreshold),
	}
	spCfg := &mconfig.SessionProxyConfig{}
	err := feg_mconfig.GetServiceConfigs(strings.ToLower(registry.SESSION_PROXY), spCfg)
	if err != nil {
		return defaultMetrics
	}
	reqFailureThreshold := spCfg.GetRequestFailureThreshold()
	minReqThreshold := spCfg.GetMinimumRequestThreshold()
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
	pcrfInitTotal, err := service_health_metrics.GetInt64("pcrf_ccr_init_requests_total")
	if err != nil {
		return nil, err
	}
	pcrfInitSendFailures, err := service_health_metrics.GetInt64("pcrf_ccr_init_send_failures_total")
	if err != nil {
		return nil, err
	}
	pcrfUpdateTotal, err := service_health_metrics.GetInt64("pcrf_ccr_update_requests_total")
	if err != nil {
		return nil, err
	}
	pcrfUpdateSendFailures, err := service_health_metrics.GetInt64("pcrf_ccr_update_send_failures_total")
	if err != nil {
		return nil, err
	}
	pcrfTerminateTotal, err := service_health_metrics.GetInt64("pcrf_ccr_terminate_requests_total")
	if err != nil {
		return nil, err
	}
	pcrfTerminateSendFailures, err := service_health_metrics.GetInt64("pcrf_ccr_terminate_send_failures_total")
	if err != nil {
		return nil, err
	}
	ocsInitTotal, err := service_health_metrics.GetInt64("ocs_ccr_init_requests_total")
	if err != nil {
		return nil, err
	}
	ocsInitSendFailures, err := service_health_metrics.GetInt64("ocs_ccr_init_send_failures_total")
	if err != nil {
		return nil, err
	}
	ocsUpdateTotal, err := service_health_metrics.GetInt64("ocs_ccr_update_requests_total")
	if err != nil {
		return nil, err
	}
	ocsUpdateSendFailures, err := service_health_metrics.GetInt64("ocs_ccr_update_send_failures_total")
	if err != nil {
		return nil, err
	}
	ocsTerminateTotal, err := service_health_metrics.GetInt64("ocs_ccr_terminate_requests_total")
	if err != nil {
		return nil, err
	}
	ocsTerminateSendFailures, err := service_health_metrics.GetInt64("ocs_ccr_terminate_send_failures_total")
	if err != nil {
		return nil, err
	}
	gxTimeouts, err := service_health_metrics.GetInt64("gx_timeouts_total")
	if err != nil {
		return nil, err
	}
	gxUnparseable, err := service_health_metrics.GetInt64("gx_unparseable_msg_total")
	if err != nil {
		return nil, err
	}
	gyTimeouts, err := service_health_metrics.GetInt64("gy_timeouts_total")
	if err != nil {
		return nil, err
	}
	gyUnparseable, err := service_health_metrics.GetInt64("gy_unparseable_msg_total")
	if err != nil {
		return nil, err
	}
	return &SessionHealthMetrics{
		PcrfInitTotal:             pcrfInitTotal,
		PcrfInitSendFailures:      pcrfInitSendFailures,
		PcrfUpdateTotal:           pcrfUpdateTotal,
		PcrfUpdateSendFailures:    pcrfUpdateSendFailures,
		PcrfTerminateTotal:        pcrfTerminateTotal,
		PcrfTerminateSendFailures: pcrfTerminateSendFailures,
		OcsInitTotal:              ocsInitTotal,
		OcsInitSendFailures:       ocsInitSendFailures,
		OcsUpdateTotal:            ocsUpdateTotal,
		OcsUpdateSendFailures:     ocsUpdateSendFailures,
		OcsTerminateTotal:         ocsTerminateTotal,
		OcsTerminateSendFailures:  ocsTerminateSendFailures,
		GxTimeouts:                gxTimeouts,
		GxUnparseableMsg:          gxUnparseable,
		GyTimeouts:                gyTimeouts,
		GyUnparseableMsg:          gyUnparseable,
	}, nil
}

func (prevMetrics *SessionHealthMetrics) GetDelta(
	currentMetrics *SessionHealthMetrics,
) (*SessionHealthMetrics, error) {
	if currentMetrics == nil {
		return nil, fmt.Errorf("Nil current SessionHealthMetrics struct provided")
	}
	deltaMetrics := &SessionHealthMetrics{
		PcrfInitTotal:             currentMetrics.PcrfInitTotal - prevMetrics.PcrfInitTotal,
		PcrfInitSendFailures:      currentMetrics.PcrfInitSendFailures - prevMetrics.PcrfInitSendFailures,
		PcrfUpdateTotal:           currentMetrics.PcrfUpdateTotal - prevMetrics.PcrfUpdateTotal,
		PcrfUpdateSendFailures:    currentMetrics.PcrfUpdateSendFailures - prevMetrics.PcrfUpdateSendFailures,
		PcrfTerminateTotal:        currentMetrics.PcrfTerminateTotal - prevMetrics.PcrfTerminateTotal,
		PcrfTerminateSendFailures: currentMetrics.PcrfTerminateSendFailures - prevMetrics.PcrfTerminateSendFailures,
		OcsInitTotal:              currentMetrics.OcsInitTotal - prevMetrics.OcsInitTotal,
		OcsInitSendFailures:       currentMetrics.OcsInitSendFailures - prevMetrics.OcsInitSendFailures,
		OcsUpdateTotal:            currentMetrics.OcsUpdateTotal - prevMetrics.OcsUpdateTotal,
		OcsUpdateSendFailures:     currentMetrics.OcsUpdateSendFailures - prevMetrics.OcsUpdateSendFailures,
		OcsTerminateTotal:         currentMetrics.OcsTerminateTotal - prevMetrics.OcsTerminateTotal,
		OcsTerminateSendFailures:  currentMetrics.OcsTerminateSendFailures - prevMetrics.OcsTerminateSendFailures,
		GxTimeouts:                currentMetrics.GxTimeouts - prevMetrics.GxTimeouts,
		GxUnparseableMsg:          currentMetrics.GxUnparseableMsg - prevMetrics.GxUnparseableMsg,
		GyTimeouts:                currentMetrics.GyTimeouts - prevMetrics.GyTimeouts,
		GyUnparseableMsg:          currentMetrics.GyUnparseableMsg - prevMetrics.GyUnparseableMsg,
	}
	// Update stored counts to current metric totals
	*prevMetrics = *currentMetrics
	return deltaMetrics, nil
}

// Generic metric functions
func UpdateGxRecentRequestMetrics(err error) {
	if err == nil {
		GxSuccessTimestamp.Set(float64(time.Now().Unix()))
		GxFailuresSinceLastSuccess.Set(0)
	} else {
		GxFailuresSinceLastSuccess.Inc()
	}
}

func UpdateGyRecentRequestMetrics(err error) {
	if err == nil {
		GySuccessTimestamp.Set(float64(time.Now().Unix()))
		GyFailuresSinceLastSuccess.Set(0)
	} else {
		GyFailuresSinceLastSuccess.Inc()
	}
}

// session_controller functions
func ReportCreateGxSession(err error) {
	UpdateGxRecentRequestMetrics(err)
	if err != nil {
		PcrfCcrInitSendFailures.Inc()
	}
	PcrfCcrInitRequests.Inc()
}

func ReportCreateGySession(err error) {
	UpdateGyRecentRequestMetrics(err)
	if err != nil {
		OcsCcrInitSendFailures.Inc()
	}
	OcsCcrInitRequests.Inc()
}

func ReportTerminateGxSession(err error) {
	UpdateGxRecentRequestMetrics(err)
	if err != nil {
		PcrfCcrTerminateSendFailures.Inc()
	} else {
		PcrfCcrTerminateRequests.Inc()
	}
}

func ReportTerminateGySession(err error) {
	UpdateGyRecentRequestMetrics(err)
	if err != nil {
		OcsCcrTerminateSendFailures.Inc()
	} else {
		OcsCcrTerminateRequests.Inc()
	}
}
