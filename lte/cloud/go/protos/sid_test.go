/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos_test

import (
	"fmt"
	"testing"

	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestSidProto(t *testing.T) {
	str := "IMSI12345"
	pb, err := protos.SidProto(str)
	assert.NoError(t, err)
	assert.True(t, pb.Id == "12345" && pb.Type == protos.SubscriberID_IMSI)

	_, err = protos.SidProto("BLAH12345")
	assert.Error(t, err)
}

func TestSidString(t *testing.T) {
	str := "IMSI12345"
	pb := protos.SubscriberID{Id: "12345"}
	out := protos.SidString(&pb)
	assert.Equal(t, out, str)
}

func TestParseImsiFromSessionId(t *testing.T) {
	randomSid := "99999"
	IMSI := "123456789"
	prefixedIMSI := fmt.Sprintf("IMSI%s", IMSI)
	magmaSid := fmt.Sprintf("%s-%s", prefixedIMSI, randomSid)

	resultIMSINoprefix, err := protos.ParseIMSIfromSessionIdNoPrefix(magmaSid)
	assert.NoError(t, err)
	assert.Equal(t, resultIMSINoprefix, IMSI)

	resultIMSIWithprefix, err := protos.ParseIMSIfromSessionIdWithPrefix(magmaSid)
	assert.NoError(t, err)
	assert.Equal(t, resultIMSIWithprefix, prefixedIMSI)

}
