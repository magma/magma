/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"time"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/radiusd/collection"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.RADIUSD)
	if err != nil {
		glog.Fatalf("Error creating RADIUSD service: %s", err)
	}
	metricAggregateRegistry := collection.NewMetricAggregateRegistry()
	metricsRequester, err := collection.NewMetricsRequester()
	if err != nil {
		glog.Fatalf("Error getting metrics requester: %s", err)
	}

	radiusdCfg := collection.GetRadiusdConfig()
	// Run Radius metrics collection Loop
	go func() {
		for {
			<-time.After(time.Duration(radiusdCfg.GetUpdateIntervalSecs()) * time.Second)
			prometheusText, err := metricsRequester.FetchMetrics()
			if err != nil {
				glog.Fatalf("Error getting metrics from server: %s", err)
			}
			metricFamilies, err := collection.ParsePrometheusText(prometheusText)
			if err != nil {
				glog.Fatalf("Unable to parse prometheus text: %s", err)
			}
			metricAggregateRegistry.Update(metricFamilies)
		}
	}()

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
