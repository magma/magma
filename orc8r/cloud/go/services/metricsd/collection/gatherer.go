/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package collection

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	prometheus_proto "github.com/prometheus/client_model/go"
)

// MetricsGatherer wraps a set of MetricCollectors, polling each collector
// at the configured interval and putting the results onto an output channel.
type MetricsGatherer struct {
	Collectors      []MetricCollector
	CollectInterval time.Duration
	OutputChan      chan *prometheus_proto.MetricFamily
}

// NewMetricsGatherer validates params and returns a new metrics gatherer.
func NewMetricsGatherer(
	collectors []MetricCollector,
	collectInterval time.Duration,
	outputChan chan *prometheus_proto.MetricFamily,
) (*MetricsGatherer, error) {
	if collectors == nil || len(collectors) == 0 {
		return nil, fmt.Errorf("MetricsGatherer must be initialized with at least one MetricCollector")
	}
	if collectInterval < 0 {
		return nil, fmt.Errorf("collectInterval should be positive")
	}

	return &MetricsGatherer{
		Collectors:      collectors,
		CollectInterval: collectInterval,
		OutputChan:      outputChan,
	}, nil
}

func (gatherer *MetricsGatherer) Run() {
	glog.V(2).Info("Running metrics gatherer")

	// Gather metrics from each collector periodically in separate goroutines
	// so a hanging collector doesn't block other collectors
	for _, collector := range gatherer.Collectors {
		go gatherer.gatherEvery(collector)
	}
}

func (gatherer *MetricsGatherer) gatherEvery(collector MetricCollector) {
	for range time.Tick(gatherer.CollectInterval) {
		fams, err := collector.GetMetrics()
		if err != nil {
			glog.Errorf("Metric collector error: %s", err)
		}

		for _, fam := range fams {
			gatherer.OutputChan <- fam
		}
	}
}
