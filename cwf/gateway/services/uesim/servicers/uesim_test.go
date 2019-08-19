/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers_test

import (
	"context"
	"testing"

	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/orc8r/cloud/go/blobstore"

	"github.com/stretchr/testify/assert"
)

func TestUESimulator_AddUE(t *testing.T) {
	store := blobstore.NewMemoryBlobStorageFactory()

	server, err := servicers.NewUESimServer(store)
	assert.NoError(t, err)

	expectedIMSI1 := "1234567890"
	expectedIMSI2 := "2345678901"
	ue1 := &protos.UEConfig{Imsi: expectedIMSI1, AuthKey: make([]byte, 16), AuthOpc: make([]byte, 16), Seq: 0}
	ue2 := &protos.UEConfig{Imsi: expectedIMSI2, AuthKey: make([]byte, 16), AuthOpc: make([]byte, 16), Seq: 0}

	_, err = server.AddUE(context.Background(), ue1)
	assert.NoError(t, err)

	_, err = server.AddUE(context.Background(), ue2)
	assert.NoError(t, err)
}
