/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessQuery(t *testing.T) {
	testQuery := "up"
	networkID := "network1"
	preprocessedQuery, err := preprocessQuery(testQuery, networkID)
	assert.NoError(t, err)
	expectedQuery := fmt.Sprintf("%s{%s=\"%s\"}", testQuery, exporters.NetworkLabelInstance, networkID)
	assert.Equal(t, expectedQuery, preprocessedQuery)
}

func TestParseTime(t *testing.T) {
	exampleUnixTimeString := "1547751342"
	exampleUnixFloatString := "1547751342.23"
	exampleUnixTime := time.Unix(1547751342, 0)

	exampleRFCTimeString := "2018-07-01T20:10:30.781Z"
	exampleBadRFCTimeString := "2018-07-01T20:10.781Z"
	exampleRFCTime, err := time.Parse(time.RFC3339, "2018-07-01T20:10:30.781Z")

	defaultTime := time.Time{}

	time, err := parseTime(exampleUnixTimeString, &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, exampleUnixTime, time)

	time, err = parseTime(exampleUnixFloatString, &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, exampleUnixTime, time)

	time, err = parseTime(exampleRFCTimeString, &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, exampleRFCTime, time)

	time, err = parseTime(exampleBadRFCTimeString, &defaultTime)
	assert.Error(t, err)

	time, err = parseTime("", &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, time, defaultTime)

	time, err = parseTime("", nil)
	assert.Error(t, err)
}
