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
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/robfig/cron"
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
	PrometheusClient v1.API
	Calculations     []Calculation
}

func NewPrometheusAnalyzer(prometheusClient v1.API, calculations []Calculation) Analyzer {
	cronJob := cron.New()
	return &PrometheusAnalyzer{
		Cron:             cronJob,
		PrometheusClient: prometheusClient,
		Calculations:     calculations,
	}
}

func (a *PrometheusAnalyzer) Schedule(schedule string) error {
	a.Cron = cron.New()

	err := a.Cron.AddFunc(schedule, a.Analyze)
	if err != nil {
		return err
	}
	return nil
}

func (a *PrometheusAnalyzer) Analyze() {
	for _, xap := range a.Calculations {
		err := xap.Calculate(a.PrometheusClient)
		if err != nil {
			glog.Errorf("Error calculating XAP metric: %s", err)
		}
	}
}

func (a *PrometheusAnalyzer) Run() {
	a.Cron.Run()
}

type Calculation interface {
	Calculate(prometheusClient v1.API) error
}

// XAPCalculation holds the parameters needed to run a XAP query and the registered
// prometheus gauge that the resulting value should be stored in
type XAPCalculation struct {
	Days            int
	ThresholdBytes  int
	QueryStepSize   time.Duration
	RegisteredGauge *prometheus.GaugeVec
	Labels          prometheus.Labels
}

// Calculate returns the number of unique users who have had a session in the
// past X days and have used over `thresholdBytes` data in that time
func (x *XAPCalculation) Calculate(prometheusClient v1.API) error {
	// List the users who have had an active session over the last X days
	uniqueUsersQuery := fmt.Sprintf(`count(max_over_time(active_sessions[%dd]) >= 1) by (imsi)`, x.Days)
	// List the users who have used at least x.ThresholdBytes of data in the last X days
	usersOverThresholdQuery := fmt.Sprintf(`count(sum(increase(octets_in[%dd])) by (imsi) > %d)`, x.Days, x.ThresholdBytes)
	// Count the users who match both conditions
	intersectionQuery := fmt.Sprintf(`count(%s and %s)`, uniqueUsersQuery, usersOverThresholdQuery)

	val, err := prometheusClient.Query(context.Background(), intersectionQuery, time.Now())
	if err != nil {
		return err
	}
	vec, ok := val.(model.Vector)
	if !ok {
		x.RegisteredGauge.With(x.Labels).Set(float64(-1))
		return fmt.Errorf("XAP query returned unexpected ValueType: %v", val.Type())
	}
	if len(vec) > 1 {
		x.RegisteredGauge.With(x.Labels).Set(float64(-1))
		return fmt.Errorf("XAP query returned more than 1 sample: %v", vec)
	}

	// If prometheus returns "no data", the value is actually 0
	if len(vec) == 0 {
		x.RegisteredGauge.With(x.Labels).Set(float64(0))
		return nil
	}
	x.RegisteredGauge.With(x.Labels).Set(float64(vec[0].Value))
	return nil
}
