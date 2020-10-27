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

import "github.com/prometheus/client_golang/prometheus"

// Prometheus counters are monotonically increasing
// Counters reset to zero on service restart
var (
	// Generic service counters
	Requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of EAP-AKA Handle requests",
	})
	FailedRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_requests_total",
		Help: "Total number of failed EAP-AKA Handle requests",
	})
	FailureNotifications = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failure_notifications_total",
		Help: "Total number of Notification Failures Returned to peers",
	})
	SwxRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_requests_total",
		Help: "Total number of SWx Proxy RPC Requiests sent",
	})
	SwxFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_failures_total",
		Help: "Total number of SWx Proxy RPC Failures",
	})
	SessionTimeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "session_timeouts_total",
		Help: "Total number of EAP-AKA Session Timeouts",
	})

	// Method Handlers metrics
	IdentityRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "identity_requests_total",
		Help: "Total number of calls to AKA Identity Handler",
	})
	FailedIdentityRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_identity_requests_total",
		Help: "Total number of failed calls to AKA Identity Handler",
	})
	ChallengeRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "challenge_requests_total",
		Help: "Total number of calls to AKA Challenge Handler",
	})
	FailedChallengeRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_challenge_requests_total",
		Help: "Total number of failed calls to AKA Challenge Handler",
	})
	ResyncRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "resync_requests_total",
		Help: "Total number of calls to AKA Resync Handler",
	})
	FailedResyncRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_resync_requests_total",
		Help: "Total number of failed calls to AKA Resync Handler",
	})
	S6aRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s6a_requests_total",
		Help: "Total number of s6a Proxy RPC Requiests sent",
	})
	S6aFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s6a_failures_total",
		Help: "Total number of s6a Proxy RPC Failures",
	})
	S6aULRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s6a_ul_requests_total",
		Help: "Total number of s6a Proxy RPC Requiests sent",
	})
	S6aULFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "s6a_ul_failures_total",
		Help: "Total number of s6a Proxy RPC Failures",
	})

	// Peer initiated failures
	PeerAuthReject = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "peer_auth_regect_total",
		Help: "Total number of AKA SubtypeAuthenticationReject calls from peer",
	})
	PeerClientError = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "peer_client_errors_total",
		Help: "Total number of AKA SubtypeClientError calls from peer",
	})
	PeerNotification = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "peer_notifications_total",
		Help: "Total number of AKA SubtypeNotification from peer",
	})
	PeerFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "peer_failures_total",
		Help: "Total number of AKA Errors/Failures originated from peers",
	})

	// Latencies
	SWxLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "swx_proxy_lat",
		Help:       "Latency of SWx Proxy requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	AuthLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "auth_lat",
		Help:       "Latency of EAP-AKA Authentication round (seconds). Only calculated for completed authentications.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	S6aLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "s6a_ai_lat",
		Help:       "Latency of s6a Proxy requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	S6aULLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "s6a_ul_lat",
		Help:       "Latency of s6a Proxy Update-Location requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func init() {
	prometheus.MustRegister(Requests, FailedRequests, FailureNotifications,
		SwxFailures, SessionTimeouts, IdentityRequests, FailedIdentityRequests,
		ChallengeRequests, FailedChallengeRequests, ResyncRequests, FailedResyncRequests,
		PeerAuthReject, PeerClientError, PeerNotification, PeerFailures, SWxLatency, AuthLatency)
}
