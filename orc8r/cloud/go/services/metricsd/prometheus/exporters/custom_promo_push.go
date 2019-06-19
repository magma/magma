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
	"io/ioutil"
	"net/http"
	"regexp"
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

var (
	prometheusLabelRegex = regexp.MustCompile("[a-zA-Z_][a-zA-Z0-9_]*")
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
		originalFamily := metricAndContext.Family
		originalFamily.Name = makeStringPointer(metricAndContext.Context.MetricName)
		// Convert all families to gauges to avoid name collisions of different
		// types.
		convertedFamilies := convertFamilyToGauges(originalFamily)
		for _, fam := range convertedFamilies {
			familyName := fam.GetName()
			fam.Metric = dropInvalidMetrics(fam.Metric, familyName)
			// if all metrics from this family were dropped, don't submit it
			if len(fam.Metric) == 0 {
				continue
			}
			for _, metric := range fam.Metric {
				addContextLabelsToMetric(metric, metricAndContext.Context)
				metric.TimestampMs = &timeStamp
			}
			if baseFamily, ok := e.familiesByName[familyName]; ok {
				addMetricsToFamily(baseFamily, fam)
			} else {
				e.familiesByName[familyName] = fam
			}
		}
	}
	return nil
}

// dropInvalidMetrics because invalid label names would cause the entire scrape
// to fail. Drop them here and log to allow good metrics through
func dropInvalidMetrics(metrics []*io_prometheus_client.Metric, familyName string) []*io_prometheus_client.Metric {
	validMetrics := make([]*io_prometheus_client.Metric, 0, len(metrics))
	for _, metric := range metrics {
		if err := validateLabels(metric); err != nil {
			glog.Errorf("Dropping metric %s because of invalid label: %v", familyName, err)
		} else {
			validMetrics = append(validMetrics, metric)
		}
	}
	return validMetrics
}

func addContextLabelsToMetric(metric *io_prometheus_client.Metric, ctx mxd_exp.MetricsContext) {
	networkAdded, gatewayAdded := false, false
	for _, label := range metric.Label {
		if label.GetName() == NetworkLabelNetwork {
			label.Value = makeStringPointer(ctx.NetworkID)
			networkAdded = true
		}
		if label.GetName() == NetworkLabelGateway {
			label.Value = makeStringPointer(ctx.GatewayID)
			gatewayAdded = true
		}
	}
	if !networkAdded {
		metric.Label = append(metric.Label,
			&io_prometheus_client.LabelPair{Name: makeStringPointer(NetworkLabelNetwork), Value: &ctx.NetworkID},
		)
	}
	if !gatewayAdded {
		metric.Label = append(metric.Label,
			&io_prometheus_client.LabelPair{Name: makeStringPointer(NetworkLabelGateway), Value: &ctx.GatewayID},
		)
	}
}

func validateLabels(metric *io_prometheus_client.Metric) error {
	for _, label := range metric.Label {
		if !prometheusLabelRegex.MatchString(label.GetName()) {
			return fmt.Errorf("label %s invalid", label.GetName())
		}
	}
	return nil
}

func addMetricsToFamily(baseFamily *io_prometheus_client.MetricFamily, newFamily *io_prometheus_client.MetricFamily) {
	baseFamily.Metric = append(baseFamily.Metric, newFamily.Metric...)
}

func familyToString(family *io_prometheus_client.MetricFamily) (string, error) {
	var buf bytes.Buffer
	_, err := expfmt.MetricFamilyToText(&buf, family)
	if err != nil {
		return "", fmt.Errorf("error writing family string: %v", err)
	}
	return buf.String(), nil
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
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error pushing to pushgateway: %v", string(respBody))
	}
	return nil
}

func (e *CustomPushExporter) resetFamilies() {
	e.familiesByName = make(map[string]*io_prometheus_client.MetricFamily)
}

func makeStringPointer(str string) *string {
	return &str
}
