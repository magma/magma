/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"testing"

	"fbc/lib/go/radius/rfc2869"
	"magma/feg/gateway/services/eap"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	tr := NewTestRunner()
	ues, err := tr.ConfigUEs(3)
	assert.NoError(t, err)

	for _, ue := range ues {
		radiusP, err := tr.Authenticate(ue.GetImsi())
		assert.NoError(t, err)

		eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
		assert.NotNil(t, eapMessage)
		assert.Equal(t, int(eapMessage[0]), eap.SuccessCode)
	}
}
