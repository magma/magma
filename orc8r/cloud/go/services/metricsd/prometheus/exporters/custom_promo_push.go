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
	"strings"
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
	prometheusNameRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
	nonPromoChars       = regexp.MustCompile("[^a-zA-Z\\d_]")
)

// CustomPushExporter pushes metrics to one or more custom prometheus pushgateways
type CustomPushExporter struct {
	familiesByName map[string]*io_prometheus_client.MetricFamily
	exportInterval time.Duration
	pushAddresses  []string
	sync.Mutex
}

// NewCustomPushExporter creates a new exporter to a custom pushgateway
func NewCustomPushExporter(pushAddresses []string) mxd_exp.Exporter {
	for i, addr := range pushAddresses {
		if !strings.HasPrefix(addr, "http") {
			pushAddresses[i] = fmt.Sprintf("http://%s", addr)
		}
	}
	return &CustomPushExporter{
		familiesByName: make(map[string]*io_prometheus_client.MetricFamily),
		exportInterval: pushInterval,
		pushAddresses:  pushAddresses,
	}
}

// Submit takes in a MetricAndContext, adds labels and timestamps to the metrics
// and stores them to be pushed later
func (e *CustomPushExporter) Submit(metrics []mxd_exp.MetricAndContext) error {
	e.Lock()
	defer e.Unlock()

	for _, metricAndContext := range metrics {
		// Don't register family if it has 0 metrics. Would cause prometheus scrape
		// to fail.
		if len(metricAndContext.Family.Metric) == 0 {
			continue
		}
		originalFamily := metricAndContext.Family
		originalFamily.Name = sanitizePrometheusName(metricAndContext.Context.MetricName)
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
				if metric.TimestampMs == nil || *metric.TimestampMs == 0 {
					timeStamp := time.Now().Unix() * 1000
					metric.TimestampMs = &timeStamp
				}
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

func (e *CustomPushExporter) Start() {
	go e.exportEvery()
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

func validateLabels(metric *io_prometheus_client.Metric) error {
	for _, label := range metric.Label {
		if !prometheusNameRegex.MatchString(label.GetName()) {
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

func (e *CustomPushExporter) exportEvery() {
	for range time.Tick(e.exportInterval) {
		errs := e.pushFamilies()
		e.resetFamilies()
		if len(errs) > 0 {
			glog.Errorf("error in pushing to pushgateway: %v", errs)
		}
	}
}

func (e *CustomPushExporter) pushFamilies() []error {
	var errs []error
	if len(e.familiesByName) == 0 {
		return []error{}
	}
	builder := strings.Builder{}

	e.Lock()
	for _, fam := range e.familiesByName {
		familyString, err := familyToString(fam)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		builder.WriteString(familyString)
		builder.WriteString("\n")
	}
	e.Unlock()

	body := builder.String()
	client := http.Client{}
	for _, address := range e.pushAddresses {
		resp, err := client.Post(address, "text/plain", bytes.NewBufferString(body))
		if err != nil {
			errs = append(errs, fmt.Errorf("error sending request to pushgateway %s: %v", address, err))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			respBody, _ := ioutil.ReadAll(resp.Body)
			errs = append(errs, fmt.Errorf("non-200 response code from pushgateway %s: %s", address, respBody))
			continue
		}
	}
	return errs
}

func (e *CustomPushExporter) resetFamilies() {
	e.familiesByName = make(map[string]*io_prometheus_client.MetricFamily)
}

func makeStringPointer(str string) *string {
	return &str
}

func sanitizePrometheusName(name string) *string {
	sanitizedName := nonPromoChars.ReplaceAllString(name, "_")
	// If still doesn't match, must be because digit is first character.
	if !prometheusNameRegex.MatchString(sanitizedName) {
		sanitizedName = "_" + sanitizedName
	}
	return &sanitizedName
}
