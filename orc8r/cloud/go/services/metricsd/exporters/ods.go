/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Implements the MetricExporter interface for saving to ODS
package exporters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	dto "github.com/prometheus/client_model/go"
)

const SERVICE_LABEL_NAME = "service"
const TAGS_LABEL_NAME = "tags"

type HttpClient interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type ODSMetricExporter struct {
	odsUrl         string
	queue          []*Sample
	queueMutex     sync.Mutex
	queueLen       int
	exportInterval time.Duration
	metricsPrefix  string
}

func NewODSExporter(
	baseUrl string,
	appId string,
	appSecret string,
	metricsPrefix string,
	queueLen int,
	exportInterval time.Duration,
) *ODSMetricExporter {
	e := new(ODSMetricExporter)
	e.odsUrl = fmt.Sprintf("%s?access_token=%s|%s", baseUrl, appId,
		appSecret)
	e.queueLen = queueLen
	e.exportInterval = exportInterval
	e.metricsPrefix = metricsPrefix
	return e
}

// Submit a Metric for writing
func (e *ODSMetricExporter) Submit(family *dto.MetricFamily, context MetricsContext) error {
	for _, metric := range family.GetMetric() {
		for _, s := range GetSamplesForMetrics(context.DecodedName, family.GetType(), metric, context.OriginatingEntity) {
			err := e.submitSample(&s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *ODSMetricExporter) submitSample(sample *Sample) error {
	e.queueMutex.Lock()
	defer e.queueMutex.Unlock()
	if len(e.queue) >= e.queueLen {
		return fmt.Errorf("ODS queue is full, dropping sample: %v", sample)
	}
	e.queue = append(e.queue, sample)
	return nil
}

func (e *ODSMetricExporter) Start() {
	go e.exportEvery()
}

func (e *ODSMetricExporter) exportEvery() {
	client := &http.Client{Timeout: 30 * time.Second}
	for range time.Tick(e.exportInterval) {
		err := e.Export(client)
		if err != nil {
			glog.Errorf("Error in syncing to ods: %v\n", err)
		}
	}
}

//Sync metrics in SampleQueue to MetricsController
func (e *ODSMetricExporter) Export(client HttpClient) error {
	e.queueMutex.Lock()
	samples := e.queue
	e.queue = []*Sample{}
	e.queueMutex.Unlock()

	if len(samples) != 0 {
		err := e.write(client, samples)
		if err != nil {
			return fmt.Errorf("Failed to sync to ODS: %v\n", err)
		}
	}
	return nil
}

// Write to ODS from queued samples or error
func (e *ODSMetricExporter) write(client HttpClient, samples []*Sample) error {
	datapoints := []ODSDatapoint{}
	for _, s := range samples {
		key := e.FormatKey(s)
		datapoints = append(datapoints, ODSDatapoint{
			Entity: fmt.Sprintf("%s.%s", e.metricsPrefix, s.entity),
			Key:    key,
			Value:  s.value,
			Tags:   e.GetTags(s),
			Time:   int(s.timestampMs)})
	}

	datapointsJson, err := json.Marshal(datapoints)
	if err != nil {
		return err
	}
	resp, err := client.PostForm(e.odsUrl, url.Values{"datapoints": {string(datapointsJson)}})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errMsg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		glog.Errorf("ODS status code: %d: %s", resp.StatusCode, errMsg)
		// Don't return an error here so we still clear these samples
	}
	return err
}

// FormatKey generates an entity name for the Sample for use with ODS.
// The entity name is a dot separated concatenation of sorted label
// key value pairs.
func (e *ODSMetricExporter) FormatKey(s *Sample) string {
	var keyBuffer bytes.Buffer
	var prefixBuffer bytes.Buffer // stores the service when found
	for _, labelPair := range s.labels {
		if strings.Compare(labelPair.GetName(), TAGS_LABEL_NAME) == 0 {
			continue
		}
		if strings.Compare(labelPair.GetName(), SERVICE_LABEL_NAME) == 0 {
			prefixBuffer.WriteString(labelPair.GetValue())
			prefixBuffer.WriteString(".")
		} else {
			keyBuffer.WriteString(".")
			keyBuffer.WriteString(labelPair.GetName())
			keyBuffer.WriteString("-")
			keyBuffer.WriteString(labelPair.GetValue())
		}
	}
	// return combined strings
	prefixBuffer.WriteString(s.name)
	prefixBuffer.Write(keyBuffer.Bytes())
	return prefixBuffer.String()
}

// GetTags parse label for tags and appends them to a comma-separated string.
func (e *ODSMetricExporter) GetTags(s *Sample) []string {
	for _, labelPair := range s.labels {
		if strings.Compare(labelPair.GetName(), TAGS_LABEL_NAME) == 0 {
			return strings.Split(labelPair.GetValue(), ",")
		}
	}
	return []string{}
}

// ODSDatapoint is used to Marshal JSON encoding for ODS data submission
// see https://phabricator.intern.facebook.com/diffusion/E/browse/tfb/trunk/www/flib/platform/graph/resources/ods/metrics/GraphOdsMetricsPost.php
// for types accepted
type ODSDatapoint struct {
	Entity string   `json:"entity"`
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	Time   int      `json:"time"`
	Tags   []string `json:"tags"`
}
