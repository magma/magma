/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils_test

import (
	"testing"

	"magma/lte/cloud/go/services/cellular/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetBand(t *testing.T) {
	expected := map[uint32]uint32{
		0:     1,
		599:   1,
		600:   2,
		749:   2,
		38650: 40,
		43590: 43,
		45589: 43,
	}

	for earfcndl, bandExpected := range expected {
		band, err := utils.GetBand(earfcndl)
		assert.NoError(t, err)
		assert.Equal(t, bandExpected, band.ID)
	}
}

func TestGetBandError(t *testing.T) {
	expectedErr := [...]uint32{60255, 61250}

	for _, earfcndl := range expectedErr {
		_, err := utils.GetBand(earfcndl)
		assert.Error(t, err, "Invalid EARFCNDL: no matching band")
	}
}
