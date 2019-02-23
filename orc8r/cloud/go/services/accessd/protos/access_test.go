/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/accessd/protos"
)

func TestAccessControlDefinitions(t *testing.T) {
	for n, v := range protos.AccessControl_Permission_value {
		if v&(v-1) != 0 {
			t.Fatalf(
				"Invalid AccessControl Permission definition: %s = %d (B%b). "+
					"AccessControl Permissions must be powers of 2.", n, v, v)
		}
	}

	assert.Equal(t, "READ", protos.AccessControl_READ.ToString())
	assert.Equal(t, "WRITE", protos.AccessControl_WRITE.ToString())
	assert.Equal(t, "NONE", protos.AccessControl_NONE.ToString())

	rwStr := (protos.AccessControl_READ | protos.AccessControl_WRITE).ToString()
	assert.True(t, "READ|WRITE" == rwStr || "WRITE|READ" == rwStr)

	assert.Equal(t, "NONE", protos.AccessControl_Permission(16).ToString())
}
