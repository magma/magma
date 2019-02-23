/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/logger/exporters"

	"github.com/stretchr/testify/assert"
)

func TestScribeEntryUtils(t *testing.T) {

	logEntries := []*protos.LogEntry{
		{
			Category:  "test",
			NormalMap: map[string]string{"status": "ACTIVE"},
			IntMap:    map[string]int64{"port": 443},
			Time:      12345,
		},
	}
	scribeEntries, err := exporters.ConvertToScribeLogEntries(logEntries)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(scribeEntries))
	assert.Equal(t, logEntries[0].Category, scribeEntries[0].Category)
	expectedMsg := "{\"int\":{\"port\":443,\"time\":12345},\"normal\":{\"status\":\"ACTIVE\"}}"
	assert.Equal(t, expectedMsg, scribeEntries[0].Message)

}
