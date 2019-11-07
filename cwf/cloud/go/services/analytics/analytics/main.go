/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"time"

	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/services/analytics"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/confignames"

	"github.com/golang/glog"
	promAPI "github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ServiceName = "ANALYTICS"

	activeUsersMetricName = "active_users_over_time"
	analysisSchedule      = "0 */12 * * *" // Every 12 hours
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(cwf.ModuleName, ServiceName)
	if err != nil {
		glog.Fatalf("Error creating CWF Analytics service: %s", err)
	}

	calculations := getAnalyticsCalculations()
	promAPIClient := getPrometheusClient()
	runAnalyses(promAPIClient, calculations, analysisSchedule)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}

func getAnalyticsCalculations() []analytics.Calculation {
	xapGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: activeUsersMetricName}, []string{"days"})
	prometheus.MustRegister(xapGauge)

	return []analytics.Calculation{
		// MAP
		&analytics.XAPCalculation{
			Days:            30,
			ThresholdBytes:  300000, //300kb
			QueryStepSize:   5 * time.Minute,
			Labels:          prometheus.Labels{"days": "30"},
			RegisteredGauge: xapGauge,
		},
		// WAP
		&analytics.XAPCalculation{
			Days:            7,
			ThresholdBytes:  70000, //70kb
			QueryStepSize:   time.Minute,
			Labels:          prometheus.Labels{"days": "7"},
			RegisteredGauge: xapGauge,
		},
		// DAP
		&analytics.XAPCalculation{
			Days:            1,
			ThresholdBytes:  10000, // 10kb
			QueryStepSize:   15 * time.Second,
			Labels:          prometheus.Labels{"days": "1"},
			RegisteredGauge: xapGauge,
		},
	}
}

func getPrometheusClient() v1.API {
	metricsConfig, err := config.GetServiceConfig(orc8r.ModuleName, metricsd.ServiceName)
	if err != nil {
		glog.Fatalf("Could not retrieve metricsd configuration: %s", err)
	}
	promClient, err := promAPI.NewClient(promAPI.Config{Address: metricsConfig.GetRequiredStringParam(confignames.PrometheusQueryAddress)})
	if err != nil {
		glog.Fatalf("Error creating prometheus client: %s", promClient)
	}
	return v1.NewAPI(promClient)
}

func runAnalyses(promAPIClient v1.API, calculations []analytics.Calculation, schedule string) {
	analyzer := analytics.NewPrometheusAnalyzer(promAPIClient, calculations)
	err := analyzer.Schedule(schedule)
	if err != nil {
		glog.Fatalf("Error scheduling analyzer: %s", err)
	}

	for _, c := range calculations {
		err = c.Calculate(promAPIClient)
		if err != nil {
			glog.Error(err)
		}
	}
	go analyzer.Run()
}
