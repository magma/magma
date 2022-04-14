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
	"github.com/prometheus/client_golang/prometheus"

	"magma/orc8r/lib/go/metrics"
)

var (
	AIRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ai_requests_total",
		Help: "Total number of AIRs received",
	})
	ULRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ul_requests_total",
		Help: "Total number of ULRs received",
	})
	PURequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pu_requests_total",
		Help: "Total number of PURs received",
	})
	InvalidRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "invalid_request_total",
		Help: "Total number of requests which did not contain the correct data",
	})
	NetworkIDErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "network_id_error_total",
		Help: "Total number of times the network ID could not be retrieved from the gRPC context",
	})
	ConfigErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "config_error_total",
		Help: "Total number of times the config could not be found",
	})
	UnknownSubscribers = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "unknown_subscriber_total",
		Help: "Total number of requests with unknown subscribers",
	})
	UnknownSubProfiles = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "unknown_sub_profile_total",
		Help: "Total number of requests with an unknown subscriber profile",
	})
	AuthErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_error_totals",
		Help: "Total number of times authentication is rejected",
	})
	ResyncAuthErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "resync_auth_error_total",
		Help: "Total number of times that resync requests fail",
	})
	StorageErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "storage_error_total",
		Help: "Total number of times storing a value in the database fails",
	})
	AuthErrorsByNetwork = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_failures_total",
		Help: "Total number of auth failures by network"},
		[]string{metrics.NetworkLabelName},
	)
	AuthSuccessesByNetwork = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_successess_total",
		Help: "Total number of auth successess by network"},
		[]string{metrics.NetworkLabelName},
	)
	UnknowSubscribersByNetwork = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "unknown_subscribers_total",
		Help: "Total number of unknown subscribers by network"},
		[]string{metrics.NetworkLabelName},
	)
)

func init() {
	prometheus.MustRegister(
		AIRequests,
		ULRequests,
		PURequests,
		InvalidRequests,
		NetworkIDErrors,
		ConfigErrors,
		UnknownSubscribers,
		UnknownSubProfiles,
		AuthErrors,
		ResyncAuthErrors,
		StorageErrors,
		UnknowSubscribersByNetwork,
		AuthSuccessesByNetwork,
		AuthErrorsByNetwork,
	)
}
