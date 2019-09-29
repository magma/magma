/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package counters

import (
	"fbc/cwf/radius/monitoring/counters/census"
	"net/http"

	"go.uber.org/zap"
)

// Init Initialize views and Prometheus exporter
func Init(config census.Config, logger *zap.Logger) {
	// Create metrics server
	census, err := config.Build(logger)
	if err != nil {
		logger.Error("Failed building census", zap.Error(err))
		return
	}
	http.Handle("/metrics", census.StatsHandler)
	go func() {
		defer census.Close()
		http.ListenAndServe(":9100", nil)
	}()
}
