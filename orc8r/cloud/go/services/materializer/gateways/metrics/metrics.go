/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	UpdatesProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gateway_state_updates_processed",
		Help: "Number of gateway state updates processed by the gateways materializer state recorder",
	})
)

func init() {
	prometheus.MustRegister(UpdatesProcessed)
}
