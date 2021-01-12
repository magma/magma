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
	// EAP (Authenticator) related metrics
	Auth = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eap_auth",
			Help: "EAP Auth responses processed, partitioned by authenticator's EapMsgCode code, " +
				"supplicant's EAP method & APN (Called Station ID). Attach failures will have code: Failure",
		},
		[]string{"code", "method", "apn", "imsi"},
	)

	// Sessions
	Sessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_sessions",
			Help: "Number of active user sessions partitioned by APN, IMSI, SessionID, MSISDN",
		},
		[]string{"apn", "imsi", "id", "msisdn"},
	)

	SessionTimeouts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_timeouts",
			Help: "Session timeouts, partitioned by APN, IMSI, MSISDN",
		},
		[]string{"apn", "imsi", "msisdn"},
	)

	SessionStart = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_start",
			Help: "Session start partitioned by APN, IMSI, SessionID, MSISDN",
		},
		[]string{"apn", "imsi", "id", "msisdn"},
	)
	SessionStop = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_stop",
			Help: "Session stop partitioned by APN, IMSI, SessionID, MSISDN",
		},
		[]string{"apn", "imsi", "id", "msisdn"},
	)

	// Latencies
	CreateSessionLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "create_session_lat",
		Help:       "Latency of accounting.CreateSession requests (seconds).",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})

	// Data usage
	OctetsIn = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "octets_in",
			Help: "Inboud data usage, partitioned by APN, IMSI, MSISDN",
		},
		[]string{"apn", "imsi", "msisdn"},
	)
	OctetsOut = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "octets_out",
			Help: "Outbound data usage, partitioned by APN, IMSI, MSISDN",
		},
		[]string{"apn", "imsi", "msisdn"},
	)

	// Acct
	AcctStop = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "accounting_stop",
			Help: "Accounting Stop Calls, partitioned by APN, IMSI, MSISDN",
		},
		[]string{"apn", "imsi", "msisdn"},
	)

	SessionTerminate = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_manager_terminate",
			Help: "Terminate Session Calls by Local Session Manager, partitioned by APN, IMSI, MSISDN",
		},
		[]string{"apn", "imsi", "msisdn"},
	)
	EndSession = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "end_session",
			Help: "EndSession Calls to Local Session Manager, partitioned by APN, IMSI, MSISDN",
		},
		[]string{"apn", "imsi", "msisdn"},
	)
)

func init() {
	prometheus.MustRegister(Auth, Sessions, SessionStart,
		SessionStop, CreateSessionLatency, OctetsIn, OctetsOut,
		SessionTimeouts, AcctStop, SessionTerminate, EndSession)
}

const imsiPrefix = "IMSI"

// DecorateIMSI prepends "IMSI" to 'clean' IMSI
func DecorateIMSI(imsi string) string {
	return imsiPrefix + imsi
}
