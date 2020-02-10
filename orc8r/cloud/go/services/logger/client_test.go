/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package logger_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/services/logger"
	"magma/orc8r/cloud/go/services/logger/test_init"
	"magma/orc8r/lib/go/protos"
)

func TestLoggingServiceClientMethods(t *testing.T) {
	mockExporter := test_init.StartTestServiceWithMockExporterExposed(t)
	scribeData := map[string]string{"url": "test_url", "ret_code": "test return code", "method": "GET(test)", "operator_info": "Txia dev test",
		"client_ip": "123.456.789.123"}

	entry := &protos.LogEntry{Category: "perfpipe_magma_rest_api_stats", NormalMap: scribeData, Time: int64(time.Now().Unix())}
	matchEntries := []*protos.LogEntry{entry}
	entries := []*protos.LogEntry{proto.Clone(entry).(*protos.LogEntry)}

	mockExporter.On("Submit", mock.MatchedBy(logEntriesMatcher(matchEntries))).Return(nil)
	err := logger.LogEntriesToDest(entries, protos.LoggerDestination_SCRIBE, 0)
	assert.NoError(t, err)
	mockExporter.AssertNotCalled(t, "Submit", entries)

	err = logger.LogEntriesToDest(entries, 1, 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = LoggerDestination 1 not supported")
	mockExporter.AssertNotCalled(t, "Submit", entries)

	mockExporter.On("Submit", mock.MatchedBy(logEntriesMatcher(matchEntries))).Return(nil)
	err = logger.LogEntriesToDest(entries, protos.LoggerDestination_SCRIBE, 1)
	assert.NoError(t, err)
	mockExporter.AssertCalled(t, "Submit", mock.AnythingOfType("[]*protos.LogEntry"))
}

func logEntriesMatcher(expected []*protos.LogEntry) interface{} {
	return func(entries []*protos.LogEntry) bool {
		cleanupEntries(expected)
		cleanupEntries(entries)
		return reflect.DeepEqual(entries, expected)
	}
}

func cleanupEntries(entries []*protos.LogEntry) {
	for i, e := range entries {
		if e != nil {
			b, err := protos.Marshal(e)
			if err == nil {
				ce := &protos.LogEntry{}
				err = protos.Unmarshal(b, ce)
				if err == nil {
					entries[i] = ce
				}
			}
		}
	}
}
