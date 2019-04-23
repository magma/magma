/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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
)

func init() {
	prometheus.MustRegister(Requests, FailedRequests, FailureNotifications,
		SwxFailures, SessionTimeouts, IdentityRequests, FailedIdentityRequests,
		ChallengeRequests, FailedChallengeRequests, ResyncRequests, FailedResyncRequests,
		PeerAuthReject, PeerClientError, PeerNotification, PeerFailures)
}
