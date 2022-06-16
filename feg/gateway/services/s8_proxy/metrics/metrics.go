/*
Copyright 2022 The Magma Authors.

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

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"

	mconfigprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	service_health_metrics "magma/feg/gateway/service_health/metrics"
	"magma/gateway/mconfig"
)

const DefaultRequestFailureThreshold = 0.50
const DefaultMinimumRequiredRequests = 1

// Prometheus counters are monotonically increasing
// Counters reset to zero on service restart
var (
	SessionCreateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_session_create_requests_total",
		Help: "Total number of create session requests.",
	})
	SessionCreateFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_session_create_failures",
		Help: "Total number of create session requests that failed.",
	})
	SessionDeleteRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_session_delete_requests_total",
		Help: "Total number of delete session requests.",
	})
	SessionDeleteFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_session_delete_failures",
		Help: "Total number of delete session requests that failed.",
	})
	BearerCreateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_bearer_create_requests_total",
		Help: "Total number of create bearer requests.",
	})
	BearerCreateFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_bearer_create_failures",
		Help: "Total number of create bearer requests that failed.",
	})
	BearerDeleteRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_bearer_delete_requests_total",
		Help: "Total number of delete bearer requests.",
	})
	BearerDeleteFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s8_bearer_delete_failures",
		Help: "Total number of delete bearer requests that failed.",
	})
)

type S8HealthTracker struct {
	Metrics                 *S8HealthMetrics
	RequestFailureThreshold float32
	MinimumRequestThreshold uint32
}

type S8HealthMetrics struct {
	SessionCreateRequests int64
	SessionCreateFails    int64
	SessionDeleteRequests int64
	SessionDeleteFails    int64
	BearerCreateRequests  int64
	BearerCreateFails     int64
	BearerDeleteRequests  int64
	BearerDeleteFails     int64
}

func init() {
	prometheus.MustRegister(SessionCreateRequests, SessionCreateFails,
		SessionDeleteRequests, SessionDeleteFails, BearerCreateRequests,
		BearerCreateFails, BearerDeleteRequests, BearerDeleteFails)
}

func NewS8HealthTracker() *S8HealthTracker {
	initMetrics := &S8HealthMetrics{
		SessionCreateRequests: 0,
		SessionCreateFails:    0,
		SessionDeleteRequests: 0,
		SessionDeleteFails:    0,
		BearerCreateRequests:  0,
		BearerCreateFails:     0,
		BearerDeleteRequests:  0,
		BearerDeleteFails:     0,
	}
	defaultHealthTracker := &S8HealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: float32(DefaultRequestFailureThreshold),
		MinimumRequestThreshold: uint32(DefaultMinimumRequiredRequests),
	}
	s8Cfg := &mconfigprotos.S8Config{}
	err := mconfig.GetServiceConfigs(strings.ToLower(registry.S8_PROXY), s8Cfg)
	if err != nil {
		return defaultHealthTracker
	}

	reqFailureThreshold := s8Cfg.GetRequestFailureThreshold()
	minReqThreshold := s8Cfg.GetMinimumRequestThreshold()

	if reqFailureThreshold == 0 {
		glog.Info("Request failure threshold cannot be 0; Using default health parameters...")
		return defaultHealthTracker
	}
	if minReqThreshold == 0 {
		glog.Info("Minimum request threshold cannot be 0; Using default health parameters...")
		return defaultHealthTracker
	}
	return &S8HealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: reqFailureThreshold,
		MinimumRequestThreshold: minReqThreshold,
	}
}

func GetCurrentHealthMetrics() (*S8HealthMetrics, error) {
	sessionCreateRequests, err := service_health_metrics.GetInt64("s8_session_create_requests_total")
	if err != nil {
		return nil, err
	}
	sessionCreateFailsNoResponse, err := service_health_metrics.GetInt64("s8_session_create_failures")
	if err != nil {
		return nil, err
	}
	sessionDeleteRequests, err := service_health_metrics.GetInt64("s8_session_delete_requests_total")
	if err != nil {
		return nil, err
	}
	sessionDeleteFailsNoResponse, err := service_health_metrics.GetInt64("s8_session_delete_failures")
	if err != nil {
		return nil, err
	}
	bearerCreateRequests, err := service_health_metrics.GetInt64("s8_bearer_create_requests_total")
	if err != nil {
		return nil, err
	}
	bearerCreateFails, err := service_health_metrics.GetInt64("s8_bearer_create_failures")
	if err != nil {
		return nil, err
	}
	bearerDeleteRequests, err := service_health_metrics.GetInt64("s8_bearer_create_delete_total")
	if err != nil {
		return nil, err
	}
	bearerDeleteFails, err := service_health_metrics.GetInt64("s8_bearer_delete_failures")
	if err != nil {
		return nil, err
	}

	return &S8HealthMetrics{
		SessionCreateRequests: sessionCreateRequests,
		SessionCreateFails:    sessionCreateFailsNoResponse,
		SessionDeleteRequests: sessionDeleteRequests,
		SessionDeleteFails:    sessionDeleteFailsNoResponse,
		BearerCreateRequests:  bearerCreateRequests,
		BearerCreateFails:     bearerCreateFails,
		BearerDeleteRequests:  bearerDeleteRequests,
		BearerDeleteFails:     bearerDeleteFails,
	}, nil
}

func (prevMetrics *S8HealthMetrics) GetDelta(currentMetrics *S8HealthMetrics) (*S8HealthMetrics, error) {
	if currentMetrics == nil {
		return nil, fmt.Errorf("nil current S8HealthMetrics struct provided")
	}
	deltaMetrics := &S8HealthMetrics{
		SessionCreateRequests: currentMetrics.SessionCreateRequests - prevMetrics.SessionCreateRequests,
		SessionCreateFails:    currentMetrics.SessionCreateFails - prevMetrics.SessionCreateFails,
		SessionDeleteRequests: currentMetrics.SessionDeleteRequests - prevMetrics.SessionDeleteRequests,
		SessionDeleteFails:    currentMetrics.SessionDeleteFails - prevMetrics.SessionDeleteFails,
		BearerCreateRequests:  currentMetrics.BearerCreateRequests - prevMetrics.BearerCreateRequests,
		BearerCreateFails:     currentMetrics.BearerCreateFails - prevMetrics.BearerCreateFails,
		BearerDeleteRequests:  currentMetrics.BearerDeleteRequests - prevMetrics.BearerDeleteRequests,
		BearerDeleteFails:     currentMetrics.BearerDeleteFails - prevMetrics.BearerDeleteFails,
	}
	// Update stored counts to current metric totals
	*prevMetrics = *currentMetrics
	return deltaMetrics, nil
}
