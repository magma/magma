/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	mxd_exp "magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/golang/glog"
	"github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

const (
	pushInterval = time.Second * 30
)

// CustomPushExporter pushes metrics to a custom prometheus pushgateway
type CustomPushExporter struct {
	familiesByName map[string]*io_prometheus_client.MetricFamily
	exportInterval time.Duration
	pushAddress    string
	sync.Mutex
}

// NewCustomPushExporter creates a new exporter to a custom pushgateway
func NewCustomPushExporter(pushAddress string) mxd_exp.Exporter {
	return &CustomPushExporter{
		familiesByName: make(map[string]*io_prometheus_client.MetricFamily),
		exportInterval: pushInterval,
		pushAddress:    pushAddress,
	}
}

// Submit takes in a MetricAndContext, adds labels and timestamps to the metrics
// and stores them to be pushed later
func (e *CustomPushExporter) Submit(metrics []mxd_exp.MetricAndContext) error {
	e.Lock()
	defer e.Unlock()

	timeStamp := time.Now().Unix() * 1000
	for _, metricAndContext := range metrics {
		// Don't register family if it has 0 metrics. Would cause prometheus scrape
		// to fail.
		if len(metricAndContext.Family.Metric) == 0 {
			continue
		}
		familyName := metricAndContext.Context.MetricName
		for _, metric := range metricAndContext.Family.Metric {
			metricType := getMetricType(metric)
			familyType := *metricAndContext.Family.Type
			// Metrics must be of the same type as their family. Otherwise prometheus
			// scrape fails.
			if metricType != familyType {
				glog.Errorf("metric type %s not same as family %s: %s\n", metricType, familyType, familyName)
				continue
			}
			addContextLabelsToMetric(metric, metricAndContext.Context)
			metric.TimestampMs = &timeStamp
		}
		if baseFamily, ok := e.familiesByName[familyName]; ok {
			addMetricsToFamily(baseFamily, metricAndContext.Family)
		} else {
			e.familiesByName[familyName] = metricAndContext.Family
		}
	}
	return nil
}

func addContextLabelsToMetric(metric *io_prometheus_client.Metric, ctx mxd_exp.MetricsContext) {
	metric.Label = append(
		metric.Label,
		&io_prometheus_client.LabelPair{Name: makeStringPointer(NetworkLabelGateway), Value: &ctx.GatewayID},
		&io_prometheus_client.LabelPair{Name: makeStringPointer(NetworkLabelNetwork), Value: &ctx.NetworkID},
	)
}

func addMetricsToFamily(baseFamily *io_prometheus_client.MetricFamily, newFamily *io_prometheus_client.MetricFamily) {
	for _, metric := range newFamily.GetMetric() {
		baseFamily.Metric = append(baseFamily.Metric, metric)
	}
}

func familyToString(family *io_prometheus_client.MetricFamily) (string, error) {
	var buf bytes.Buffer
	_, err := expfmt.MetricFamilyToText(&buf, family)
	if err != nil {
		return "", fmt.Errorf("error writing family string: %v", err)
	}
	return buf.String(), nil
}

func getMetricType(metric *io_prometheus_client.Metric) io_prometheus_client.MetricType {
	if metric.Counter != nil {
		return io_prometheus_client.MetricType_COUNTER
	} else if metric.Gauge != nil {
		return io_prometheus_client.MetricType_GAUGE
	} else if metric.Summary != nil {
		return io_prometheus_client.MetricType_SUMMARY
	} else if metric.Histogram != nil {
		return io_prometheus_client.MetricType_HISTOGRAM
	}
	return io_prometheus_client.MetricType_UNTYPED
}

// Start runs exportEvery() in a goroutine to continuously push metrics at every
// push interval
func (e *CustomPushExporter) Start() {
	go e.exportEvery()
}

func (e *CustomPushExporter) exportEvery() {
	for range time.Tick(e.exportInterval) {
		err := e.export()
		if err != nil {
			glog.Errorf("error in pushing to pushgateway: %v", err)
		}
	}
}

func (e *CustomPushExporter) export() error {
	err := e.pushFamilies()
	e.resetFamilies()
	return err
}

func (e *CustomPushExporter) pushFamilies() error {
	if len(e.familiesByName) == 0 {
		return nil
	}
	body := bytes.Buffer{}
	for _, fam := range e.familiesByName {
		familyString, err := familyToString(fam)
		if err != nil {
			return err
		}
		body.WriteString(familyString)
		body.WriteString("\n")
	}
	client := http.Client{}
	resp, err := client.Post(e.pushAddress, "text/plain", &body)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error pushing to pushgateway: %v", err)
	}
	return nil
}

func (e *CustomPushExporter) resetFamilies() {
	e.familiesByName = make(map[string]*io_prometheus_client.MetricFamily)
}

func makeStringPointer(str string) *string {
	return &str
}
