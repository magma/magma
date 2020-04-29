/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers_test

import (
	"reflect"
	"testing"

	"fbc/lib/go/radius/rfc2869"

	"github.com/stretchr/testify/assert"
)

const (
	EapIdentityResponseMessage = "\x02\x00\x00\x38\x01\x30\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30" +
		"\x30\x30\x30\x39\x31\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30" +
		"\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74" +
		"\x77\x6f\x72\x6b\x2e\x6f\x72\x67"
	CalledStationID2 = "99-FE-FF-84-B5-46:CWF-TP-LINK_B547_5G"
)

func TestCreateEAPIdentityRequest(t *testing.T) {
	server, _, err := setupTest()
	assert.NoError(t, err)

	radiusP, err := server.CreateEAPIdentityRequest(Imsi, CalledStationID2)
	assert.NoError(t, err)

	eapMessage := []byte(radiusP.Get(rfc2869.EAPMessage_Type))
	assert.True(t, reflect.DeepEqual(eapMessage, []byte(EapIdentityResponseMessage)))
}
