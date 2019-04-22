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
