/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package collection_test

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/collection"

	"github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHttpClient struct {
	mock.Mock
}

func (_m *mockHttpClient) Get(url string) (*http.Response, error) {
	ret := _m.Called(url)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(string) *http.Response); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Get(0).(*http.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Wraps bytes.Buffer to implement ReadCloser
type closeableBuffer struct {
	*bytes.Buffer
}

func (*closeableBuffer) Close() error {
	return nil
}

func TestKafkaConnectCollector_GetMetrics(t *testing.T) {
	mockClient := &mockHttpClient{}
	coll := collection.NewKafkaConnectCollector("foo", mockClient)

	// Happy path
	responseBodyStr := `
	{
		"name": "foo-connector",
		"connector": {
			"state": "RUNNING",
			"worker_id": "172.16.48.12:8083"
		},
		"tasks": [
			{
				"state": "RUNNING",
				"id": 0,
				"worker_id": "172.16.16.10:8083"
			}
		],
		"type": "source"
	}
	`
	mockResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       &closeableBuffer{bytes.NewBufferString(responseBodyStr)},
	}
	mockClient.On("Get", "http://localhost:8083/connectors/foo/status").
		Return(mockResponse, nil)

	res, err := coll.GetMetrics()
	assert.NoError(t, err)
	assert.Equal(t, getExpectedMetrics(true), res)

	// Failed connector
	responseBodyStr = `
	{
		"name": "foo-connector",
		"connector": {
			"state": "FAILED"
		},
		"tasks": [
			{
				"state": "FAILED",
				"id": 0,
				"trace": "FOO\nBAR"
			}
		],
		"type": "source"
	}
	`
	mockResponse.Body = &closeableBuffer{bytes.NewBufferString(responseBodyStr)}

	res, err = coll.GetMetrics()
	assert.NoError(t, err)
	assert.Equal(t, getExpectedMetrics(false), res)

	// Bad status code
	mockResponse.Status = "500 Error"
	mockResponse.StatusCode = 500
	mockResponse.Body = &closeableBuffer{bytes.NewBufferString("foobar")}

	_, err = coll.GetMetrics()
	assert.Error(t, err, "Status code 500 returned from Kafka Connect")

	// Http get error
	mockClient.On("Get", "http://localhost:8083/connectors/foo/status").
		Return(nil, errors.New("foo error"))
	_, err = coll.GetMetrics()
	assert.Error(t, err, "foo error")
}

func getExpectedMetrics(success bool) []*io_prometheus_client.MetricFamily {
	name := "kafka_connector_status"
	help := "1 if kafka connector is healthy, 0 otherwise"
	mtype := io_prometheus_client.MetricType_GAUGE

	labelName := "connectorName"
	labelValue := "foo"

	var val float64
	if success {
		val = 1
	} else {
		val = 0
	}

	return []*io_prometheus_client.MetricFamily{
		{
			Name: &name,
			Help: &help,
			Type: &mtype,
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						{Name: &labelName, Value: &labelValue},
					},
					Gauge: &io_prometheus_client.Gauge{Value: &val},
				},
			},
		},
	}
}
