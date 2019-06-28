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

	ChallengeRequestEapPacket = "\x01\xea\x00\x44\x17\x01\x00\x00\x01\x05\x00\x00\xee\xb3\x53\x6c" +
		"\x2f\xc3\x68\xfe\x3a\xfb\xd5\x5c\xfe\xf9\x6b\x29\x02\x05\x00\x00" +
		"\x94\x73\x37\x74\x82\xbd\x67\x41\x51\x11\x05\x57\x68\x17\xaa\x23" +
		"\x0b\x05\x00\x00\xda\x14\xa9\xce\x0e\x66\xaf\x38\x7b\x9f\xc1\xe6" +
		"\xf0\x31\x5e\x00"
	ChallengeResponseEapPacket = "\x02\xea\x00\x40\x17\x01\x00\x00\x03\x03\x00\x40\xdc\x89\x15\x16" +
		"\x8d\xd2\xeb\x56\x86\x06\x00\x00\x86\xe8\x20\x4d\xc6\xe1\xe3\xd8" +
		"\x94\x44\x3c\x26\xa7\xc6\x5d\xee\x3c\x42\xab\xf8\x0b\x05\x00\x00" +
		"\x13\x00\x7f\xe9\x86\xfc\xc1\x54\xf5\xca\x2b\xa7\x23\x88\x6d\x5b"

	Imsi = "\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x39\x31"
	Key  = "\x8B\xAF\x47\x3F\x2F\x8F\xD0\x94\x87\xCC\xCB\xD7\x09\x7C\x68\x62"
	Opc  = "\x8e\x27\xb6\xaf\x0e\x69\x2e\x75\x0f\x32\x66\x7a\x3b\x14\x60\x5d"
	Sqn  = 32
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

func TestChallengeRequest(t *testing.T) {
	store := blobstore.NewMemoryBlobStorageFactory()

	server, err := servicers.NewUESimServer(store)
	assert.NoError(t, err)

	ue := &protos.UEConfig{Imsi: Imsi, AuthKey: []byte(Key), AuthOpc: []byte(Opc), Seq: Sqn}
	_, err = server.AddUE(context.Background(), ue)
	assert.NoError(t, err)

	res, err := server.Handle(Imsi, eap.Packet(ChallengeRequestEapPacket))
	assert.NoError(t, err)
	assert.True(
		t,
		reflect.DeepEqual([]byte(res), []byte(ChallengeResponseEapPacket)),
		"Actual packet didn't match expected packet\nexpected: %x\nactual:   %x\n",
		ChallengeResponseEapPacket,
		res,
	)
}
