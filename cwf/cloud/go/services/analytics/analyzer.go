/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package analytics

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/robfig/cron/v3"
)

type Analyzer interface {
	// Schedule the analyzer to run calculations periodically based on the
	// cron expression format schedule parameter
	Schedule(schedule string) error

	// Run triggers the analyzer's cronjob to start running. This function
	// blocks.
	Run()
}

// PrometheusAnalyzer accesses prometheus metrics and performs
// queries/aggregations to calculate various metrics
type PrometheusAnalyzer struct {
	Cron             *cron.Cron
	PrometheusClient PrometheusAPI
	Calculations     []Calculation
	Exporter         Exporter
}

type PrometheusAPI interface {
	Query(ctx context.Context, query string, ts time.Time) (model.Value, api.Warnings, error)
	QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, api.Warnings, error)
}

func NewPrometheusAnalyzer(prometheusClient v1.API, calculations []Calculation, exporter Exporter) Analyzer {
	cronJob := cron.New()
	return &PrometheusAnalyzer{
		Cron:             cronJob,
		PrometheusClient: prometheusClient,
		Calculations:     calculations,
		Exporter:         exporter,
	}
}

func (a *PrometheusAnalyzer) Schedule(schedule string) error {
	a.Cron = cron.New()

	_, err := a.Cron.AddFunc(schedule, a.Analyze)
	if err != nil {
		return err
	}
	return nil
}

func (a *PrometheusAnalyzer) Analyze() {
	for _, calc := range a.Calculations {
		results, err := calc.Calculate(a.PrometheusClient)
		if err != nil {
			glog.Errorf("Error calculating metric: %s", err)
			continue
		}
		if a.Exporter == nil {
			continue
		}
		for _, res := range results {
			err = a.Exporter.Export(res, http.DefaultClient)
			if err != nil {
				glog.Errorf("Error exporting result: %v", err)
			} else {
				glog.Infof("Exported %s, %s, %f", res.metricName, res.labels, res.value)
			}
		}
	}
}

func (a *PrometheusAnalyzer) Run() {
	a.Cron.Run()
}
