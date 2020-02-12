/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/logger/exporters"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (client *MockClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	args := client.Called(url, data)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestScribeExporter_Submit(t *testing.T) {
	exporter := exporters.NewScribeExporter(
		"",
		"",
		"",
		5,
		time.Second*10,
	)
	logEntries := []*protos.LogEntry{{Category: "test"}, {Category: "test2"}}
	err := exporter.Submit(logEntries)
	assert.EqualError(t, err, fmt.Sprintf("ScribeLogEntry %v doesn't have time field set", logEntries[0]))
	logEntries = []*protos.LogEntry{{Category: "test1", Time: 12345}, {Category: "test2", Time: 23456}}
	err = exporter.Submit(logEntries)
	assert.NoError(t, err)
	err = exporter.Submit(logEntries)
	assert.NoError(t, err)
	// submiting when queue is full should give error
	err = exporter.Submit(logEntries)
	assert.NoError(t, err) // queue is cleared
	logEntries = []*protos.LogEntry{
		{Category: "test1", Time: 12345},
		{Category: "test2", Time: 23456},
		{Category: "test3", Time: 23457},
		{Category: "test4", Time: 23458},
		{Category: "test5", Time: 23459},
		{Category: "test6", Time: 23460},
	}
	err = exporter.Submit(logEntries)
	assert.EqualError(t, err, fmt.Sprintf("dropping %v logEntries as it exceeds max queue length", len(logEntries)))
}

func TestScribeExporter_Export(t *testing.T) {
	exporter := exporters.NewScribeExporter(
		"",
		"",
		"",
		2,
		time.Second*10,
	)
	client := new(MockClient)
	resp := &http.Response{StatusCode: 200}
	logEntries := []*protos.LogEntry{{Category: "test1", Time: 12345}, {Category: "test2", Time: 23456}}
	scribeEntries, err := exporters.ConvertToScribeLogEntries(logEntries)
	assert.NoError(t, err)
	logJson, err := json.Marshal(scribeEntries)
	assert.NoError(t, err)
	client.On("PostForm", mock.AnythingOfType("string"), url.Values{"access_token": {"|"}, "logs": {string(logJson)}}).Return(resp, nil)
	err = exporter.Export(client)
	assert.NoError(t, err)
	client.AssertNotCalled(t, "PostForm", mock.AnythingOfType("string"), mock.AnythingOfType("url.Values"))
	err = exporter.Submit(logEntries)
	assert.NoError(t, err)
	err = exporter.Export(client)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}
