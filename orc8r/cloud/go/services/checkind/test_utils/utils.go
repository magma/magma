/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/checkind"

	"github.com/stretchr/testify/assert"
)

func Checkin(t *testing.T, req *protos.CheckinRequest) *protos.CheckinResponse {
	conn, err := registry.GetConnection(checkind.ServiceName)
	assert.NoError(t, err)
	client := protos.NewCheckindClient(conn)
	resp, err := client.Checkin(context.Background(), req)
	assert.NoError(t, err)
	return resp
}
