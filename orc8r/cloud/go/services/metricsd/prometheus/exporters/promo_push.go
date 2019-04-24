/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"fmt"
	"sync"
	"time"

	mxd_exp "magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/golang/glog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	PushJob      = "PrometheusPushGateway"
	PushInterval = time.Second * 15
)

// PrometheusExporter handles registering and updating prometheus metrics
type PrometheusPushExporter struct {
	exporter       mxd_exp.Exporter
	Pusher         *push.Pusher
	exportInterval time.Duration
	pushAddress    string
	sync.Mutex
}

// NewPrometheusExporter create a new PrometheusExporter with own registry
func NewPrometheusPushExporter(pushAddress string) mxd_exp.Exporter {
	config := PrometheusExporterConfig{
		UseHostLabel: false,
	}
	exporter := NewPrometheusExporter(config)
	pusher := push.New(pushAddress, PushJob)
	pusher.Gatherer(exporter.(*PrometheusExporter).Registry.(*prometheus.Registry))

	return &PrometheusPushExporter{
		exporter:       exporter,
		Pusher:         pusher,
		exportInterval: PushInterval,
		pushAddress:    pushAddress,
	}
}

// Submit registers or updates a prometheus metric in the exporter registry.
// All metrics in registry are then pushed to the pushgateway
func (e *PrometheusPushExporter) Submit(metrics []mxd_exp.MetricAndContext) error {
	// Unregister All Metrics in PrometheusExporter, then register new registry
	// with Pusher Before submitting new ones to avoid pushing stale metrics

	e.Lock()
	err := e.exporter.Submit(metrics)
	e.Unlock()
	if err != nil {
		return fmt.Errorf("error pushing metrics: %v\n", err)
	}
	return nil
}

func (e *PrometheusPushExporter) Start() {
	go e.exportEvery()
}

func (e *PrometheusPushExporter) Export() error {
	err := e.Pusher.Push()
	e.resetMetrics()
	return err
}

func (e *PrometheusPushExporter) resetMetrics() {
	e.Lock()
	defer e.Unlock()
	e.exporter.(*PrometheusExporter).clearRegistry()
	e.Pusher = push.New(e.pushAddress, PushJob)
	e.Pusher.Gatherer(e.exporter.(*PrometheusExporter).Registry.(*prometheus.Registry))
}

func (e *PrometheusPushExporter) exportEvery() {
	for range time.Tick(e.exportInterval) {
		err := e.Export()
		if err != nil {
			glog.Errorf("error in pushing to pushgateway: %v", err)
		}
	}
}
