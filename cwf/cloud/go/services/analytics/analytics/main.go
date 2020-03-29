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
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/confignames"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	promAPI "github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ServiceName = "ANALYTICS"

	activeUsersMetricName     = "active_users_over_time"
	userThroughputMetricName  = "user_throughput"
	userConsumptionMetricName = "user_consumption"
	apThroughputMetricName    = "throughput_per_ap"

	defaultAnalysisSchedule = "0 */12 * * *" // Every 12 hours
)

var (
	// Map from number of days to query to size the step should be to get best granularity
	// without causes prometheus to reject the query for having too many datapoints
	daysToQueryStepSize = map[int]time.Duration{1: 15 * time.Second, 7: time.Minute, 30: 5 * time.Minute}

	daysToCalculate = []int{1, 7, 30}
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(cwf.ModuleName, ServiceName)
	if err != nil {
		glog.Fatalf("Error creating CWF Analytics service: %s", err)
	}

	analysisSchedule := defaultAnalysisSchedule
	providedSchedule, _ := srv.Config.GetStringParam("analysisSchedule")
	if providedSchedule != "" {
		analysisSchedule = providedSchedule
	}

	calculations := getAnalyticsCalculations()
	promAPIClient := getPrometheusClient()
	shouldExportData, _ := srv.Config.GetBoolParam("exportMetrics")
	var exporter analytics.Exporter
	if shouldExportData {
		glog.Errorf("Creating CWF Analytics Exporter")
		exporter = analytics.NewWWWExporter(
			srv.Config.GetRequiredStringParam("metricsPrefix"),
			srv.Config.GetRequiredStringParam("appSecret"),
			srv.Config.GetRequiredStringParam("appID"),
			srv.Config.GetRequiredStringParam("metricExportURL"),
			srv.Config.GetRequiredStringParam("categoryName"),
		)
	}
	analyzer := analytics.NewPrometheusAnalyzer(promAPIClient, calculations, exporter)
	err = analyzer.Schedule(analysisSchedule)
	if err != nil {
		glog.Fatalf("Error scheduling analyzer: %s", err)
	}

	go analyzer.Run()

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}

func getAnalyticsCalculations() []analytics.Calculation {
	xapGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: activeUsersMetricName}, []string{analytics.DaysLabel, metrics.NetworkLabelName})
	userThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userThroughputMetricName}, []string{analytics.DaysLabel, metrics.NetworkLabelName, analytics.DirectionLabel})
	userConsumptionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userConsumptionMetricName}, []string{analytics.DaysLabel, metrics.NetworkLabelName, analytics.DirectionLabel})
	apThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: apThroughputMetricName}, []string{analytics.DaysLabel, metrics.NetworkLabelName, analytics.DirectionLabel, analytics.APNLabel})
	prometheus.MustRegister(xapGauge, userThroughputGauge, userConsumptionGauge, apThroughputGauge)

	allCalculations := make([]analytics.Calculation, 0)

	// MAP, WAP, DAP Calculations
	allCalculations = append(allCalculations, getXAPCalculations(daysToCalculate, xapGauge, activeUsersMetricName)...)

	// User Throughput Calculations
	allCalculations = append(allCalculations, getUserThroughputCalculations(daysToCalculate, userThroughputGauge, userThroughputMetricName)...)

	// AP Throughput Calculations
	allCalculations = append(allCalculations, getAPThroughputCalculations(daysToCalculate, apThroughputGauge, apThroughputMetricName)...)

	// User Consumption Calculations
	allCalculations = append(allCalculations, getUserConsumptionCalculations(daysToCalculate, userConsumptionGauge, userConsumptionMetricName)...)

	return allCalculations
}

func getXAPCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string) []analytics.Calculation {
	calcs := make([]analytics.Calculation, 0)
	for _, dayParam := range daysList {
		calcs = append(calcs, &analytics.XAPCalculation{
			CalculationParams: analytics.CalculationParams{
				Days:            dayParam,
				RegisteredGauge: gauge,
				Labels:          prometheus.Labels{analytics.DaysLabel: string(dayParam)},
				Name:            metricName,
			},
		})
	}
	return calcs
}

func getUserThroughputCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string) []analytics.Calculation {
	calcs := make([]analytics.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []analytics.ConsumptionDirection{analytics.ConsumptionIn, analytics.ConsumptionOut} {
			calcs = append(calcs, &analytics.UserThroughputCalculation{
				CalculationParams: analytics.CalculationParams{
					Days:            dayParam,
					RegisteredGauge: gauge,
					Labels:          prometheus.Labels{analytics.DaysLabel: string(dayParam)},
					Name:            metricName,
				},
				Direction:     dir,
				QueryStepSize: daysToQueryStepSize[dayParam],
			})
		}
	}
	return calcs
}

func getAPThroughputCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string) []analytics.Calculation {
	calcs := make([]analytics.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []analytics.ConsumptionDirection{analytics.ConsumptionIn, analytics.ConsumptionOut} {
			calcs = append(calcs, &analytics.APThroughputCalculation{
				CalculationParams: analytics.CalculationParams{
					Days:            dayParam,
					RegisteredGauge: gauge,
					Labels:          prometheus.Labels{analytics.DaysLabel: string(dayParam)},
					Name:            metricName,
				},
				Direction:     dir,
				QueryStepSize: daysToQueryStepSize[dayParam],
			})
		}
	}
	return calcs
}

func getUserConsumptionCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string) []analytics.Calculation {
	calcs := make([]analytics.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []analytics.ConsumptionDirection{analytics.ConsumptionIn, analytics.ConsumptionOut} {
			calcs = append(calcs, &analytics.UserConsumptionCalculation{
				CalculationParams: analytics.CalculationParams{
					Days:            dayParam,
					RegisteredGauge: gauge,
					Labels:          prometheus.Labels{analytics.DaysLabel: string(dayParam)},
					Name:            metricName,
				},
				Direction: dir,
			})
		}
	}
	return calcs
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
