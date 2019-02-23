/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package collection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_model/go"
)

const running = "RUNNING"

// KafkaConnectCollector is a MetricCollector which uses Kafka Connect's REST
// API to query for the status of a specific Kafka connector. It will return a
// single metric named "kafka_connector_status" with a label "connectorName"
// corresponding to the monitored connector. The value will be 1 if the
// connector is healthy and 0 otherwise.
type KafkaConnectCollector struct {
	connectorName string
	client        HttpClient
}

// HttpClient is a interface to wrap the http.Client methods we depend on in
// order to mock out that dependency for testing.
type HttpClient interface {
	Get(url string) (*http.Response, error)
}

// NewKafkaConnectCollector instantiates a collector instance for the specified
// connector. For testing purposes, this supports injecting an http Client.
// If `client` is nil, this defaults to `http.DefaultClient`.
func NewKafkaConnectCollector(connectorName string, client HttpClient) MetricCollector {
	if client == nil {
		return &KafkaConnectCollector{connectorName: connectorName, client: http.DefaultClient}
	}
	return &KafkaConnectCollector{connectorName: connectorName, client: client}
}

func (kcc *KafkaConnectCollector) GetMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	url := fmt.Sprintf("http://localhost:8083/connectors/%s/status", kcc.connectorName)
	resp, err := kcc.client.Get(url)
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []*io_prometheus_client.MetricFamily{}, fmt.Errorf("Status code %d returned from Kafka Connect", resp.StatusCode)
	}

	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, fmt.Errorf("Error while reading response body: %s", err)
	}

	unmarshaledBody := &kafkaConnectStatus{}
	err = json.Unmarshal(bodyData, unmarshaledBody)
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, fmt.Errorf("Error while unmarshaling Kafka Connect response: %s", err)
	}

	return []*io_prometheus_client.MetricFamily{
		makeConnectorStatusMetric(kcc.connectorName, unmarshaledBody),
	}, nil
}

func makeConnectorStatusMetric(connectorName string, status *kafkaConnectStatus) *io_prometheus_client.MetricFamily {
	name := "kafka_connector_status"
	help := "1 if kafka connector is healthy, 0 otherwise"

	labelName := "connectorName"
	return MakeSingleGaugeFamily(
		name,
		help,
		&MetricLabel{Name: labelName, Value: connectorName},
		getGaugeValue(status),
	)
}

func getGaugeValue(status *kafkaConnectStatus) float64 {
	if status.Connector != nil && status.Connector.State == running {
		return 1
	} else {
		return 0
	}
}

type kafkaConnectStatus struct {
	Connector *kafkaConnectorStatus `json:"connector"`
}

type kafkaConnectorStatus struct {
	State string `json:"state"`
}
