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

	"github.com/prometheus/client_golang/prometheus"

	service_health_metrics "magma/feg/gateway/service_health/metrics"
)

const DefaultRequestFailureThreshold = 0.50
const DefaultMinimumRequiredRequests = 1

// Prometheus counters are monotonically increasing
// Counters reset to zero on service restart
var (
	SessionCreateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "session_create_requests_total",
		Help: "Total number of create session requests.",
	})
	SessionCreateFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "session_create_fails",
		Help: "Total number of create session requests that failed.",
	})
	SessionDeleteRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "session_delete_requests_total",
		Help: "Total number of delete session requests.",
	})
	SessionDeleteFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "session_delete_fails",
		Help: "Total number of delete session requests that failed.",
	})
	BearerCreateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bearer_create_requests_total",
		Help: "Total number of create bearer requests.",
	})
	BearerCreateFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bearer_create_fails",
		Help: "Total number of create bearer requests that failed.",
	})
	BearerDeleteRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bearer_delete_requests_total",
		Help: "Total number of delete bearer requests.",
	})
	BearerDeleteFails = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bearer_delete_fails",
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
	HealthTracker := &S8HealthTracker{
		Metrics:                 initMetrics,
		RequestFailureThreshold: float32(DefaultRequestFailureThreshold),
		MinimumRequestThreshold: uint32(DefaultMinimumRequiredRequests),
	}
	return HealthTracker
}

func GetCurrentHealthMetrics() (*S8HealthMetrics, error) {
	sessionCreateRequests, err := service_health_metrics.GetInt64("session_create_requests_total")
	if err != nil {
		return nil, err
	}
	sessionCreateFailsNoResponse, err := service_health_metrics.GetInt64("session_create_fails")
	if err != nil {
		return nil, err
	}
	sessionDeleteRequests, err := service_health_metrics.GetInt64("session_delete_requests_total")
	if err != nil {
		return nil, err
	}
	sessionDeleteFailsNoResponse, err := service_health_metrics.GetInt64("session_delete_fails")
	if err != nil {
		return nil, err
	}
	bearerCreateRequests, err := service_health_metrics.GetInt64("bearer_create_requests_total")
	if err != nil {
		return nil, err
	}
	bearerCreateFails, err := service_health_metrics.GetInt64("bearer_create_fails")
	if err != nil {
		return nil, err
	}
	bearerDeleteRequests, err := service_health_metrics.GetInt64("bearer_create_delete_total")
	if err != nil {
		return nil, err
	}
	bearerDeleteFails, err := service_health_metrics.GetInt64("bearer_delete_fails")
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
