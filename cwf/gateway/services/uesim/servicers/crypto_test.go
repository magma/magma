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
	"reflect"
	"testing"

	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/feg/gateway/services/eap"
	"magma/orc8r/cloud/go/blobstore"

	"github.com/stretchr/testify/assert"
)

// EAP packets taken from cwf_2_aps.pcap
const (
	IdentityRequestEapPacket  = "\x01\xe8\x00\x0c\x17\x05\x00\x00\x0a\x01\x00\x00"
	IdentityResponseEapPacket = "\x02\xe8\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31" +
		"\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x39\x31\x40\x77\x6c\x61" +
		"\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e" +
		"\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00"

	Imsi = "\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x39\x31"
)

func TestIdentityRequest(t *testing.T) {
	store := blobstore.NewMemoryBlobStorageFactory()

	server, err := servicers.NewUESimServer(store)
	assert.NoError(t, err)

	ue := &protos.UEConfig{Imsi: Imsi, AuthKey: make([]byte, 16), AuthOpc: make([]byte, 16), Seq: 0}
	_, err = server.AddUE(context.Background(), ue)
	assert.NoError(t, err)

	res, err := server.Handle(Imsi, eap.Packet(IdentityRequestEapPacket))
	assert.NoError(t, err)
	assert.True(
		t,
		reflect.DeepEqual([]byte(res), []byte(IdentityResponseEapPacket)),
		"Actual packet didn't match expected packet\nexpected: %x\nactual:   %x\n",
		IdentityResponseEapPacket,
		res,
	)
}
