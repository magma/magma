/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"fmt"
	"reflect"
	"testing"

	"fbc/lib/go/radius/rfc2869"
	"magma/feg/gateway/services/eap"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUplinkTraffic(t *testing.T) {
	fmt.Printf("Running TestAuthenticateUplinkTraffic...\n")
	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	ue := ues[0]
	err = ruleManager.AddDynamicPassAllToPCRF(ue.GetImsi(), "dynamic-pass-all", "mkey1")
	assert.NoError(t, err)
	radiusP, err := tr.Authenticate(ue.GetImsi())
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))

	err = tr.GenULTraffic(ue.GetImsi(), swag.String("100K"))
	assert.NoError(t, err)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}
