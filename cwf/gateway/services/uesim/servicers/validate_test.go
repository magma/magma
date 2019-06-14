/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"errors"
	"testing"

	"magma/cwf/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestValidateUEData(t *testing.T) {
	err := validateUEData(nil)
	assert.Exactly(t, errors.New("Invalid Argument: UE data cannot be nil"), err)

	ue := &protos.UEConfig{Imsi: "0123456789", AuthKey: make([]byte, 32), AuthOpc: make([]byte, 32)}
	err = validateUEData(ue)
	assert.NoError(t, err)
}

func TestValidateUEIMSI(t *testing.T) {
	err := validateUEIMSI("")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must be between 5 and 15 digits long"), err)

	err = validateUEIMSI("0123")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must be between 5 and 15 digits long"), err)

	err = validateUEIMSI("0123456789012345")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must be between 5 and 15 digits long"), err)

	err = validateUEIMSI("0ABCDEF")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must only be digits"), err)

	err = validateUEIMSI("ABCDEF0")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must only be digits"), err)

	err = validateUEIMSI("0123456789")
	assert.NoError(t, err)
}

func TestValidateUEKey(t *testing.T) {
	err := validateUEKey(nil)
	assert.Exactly(t, errors.New("Invalid Argument: key cannot be nil"), err)

	err = validateUEKey(make([]byte, 5))
	assert.Exactly(t, errors.New("Invalid Argument: key must be 32 bytes"), err)

	err = validateUEKey(make([]byte, 32))
	assert.NoError(t, err)
}

func TestValidateUEOpc(t *testing.T) {
	err := validateUEOpc(nil)
	assert.Exactly(t, errors.New("Invalid Argument: opc cannot be nil"), err)

	err = validateUEOpc(make([]byte, 5))
	assert.Exactly(t, errors.New("Invalid Argument: opc must be 32 bytes"), err)

	err = validateUEOpc(make([]byte, 32))
	assert.NoError(t, err)
}
