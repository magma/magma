/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test

import (
	"testing"

	"magma/feg/gateway/services/csfb/servicers/decode"
	decodeIE "magma/feg/gateway/services/csfb/servicers/decode/ie"
	encodeIE "magma/feg/gateway/services/csfb/servicers/encode/ie"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeIMSI(t *testing.T) {
	// odd number of digits
	imsi := "11111111"
	encodedIMSI, err := encodeIE.EncodeIMSI(imsi)
	assert.NoError(t, err)
	decodedIMSI, ieLength, err := decodeIE.DecodeIMSI(encodedIMSI)
	assert.NoError(t, err)
	assert.Equal(t, 7, ieLength)
	assert.Equal(t, imsi, decodedIMSI)

	// odd number of digits
	imsi = "1111111"
	encodedIMSI, err = encodeIE.EncodeIMSI(imsi)
	assert.NoError(t, err)
	decodedIMSI, ieLength, err = decodeIE.DecodeIMSI(encodedIMSI)
	assert.NoError(t, err)
	assert.Equal(t, 6, ieLength)
	assert.Equal(t, imsi, decodedIMSI)
}

func TestEncodeDecodeFixedLengthIE(t *testing.T) {
	tmsi := []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)}
	encodedTMSI, err := encodeIE.EncodeFixedLengthIE(decode.IEITMSI, decode.IELengthTMSI, tmsi)
	assert.NoError(t, err)
	decodedTMSI, err := decodeIE.DecodeFixedLengthIE(encodedTMSI, decode.IELengthTMSI, decode.IEITMSI)
	assert.NoError(t, err)
	assert.Equal(t, tmsi, decodedTMSI)
}

func TestEncodeDecodeLimitedLengthIE(t *testing.T) {
	nasMessageContainer := []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)}
	encodedNASMessageContainer, err := encodeIE.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		nasMessageContainer,
	)
	assert.NoError(t, err)
	decodedNASMessageContainer, ieLength, err := decodeIE.DecodeLimitedLengthIE(
		encodedNASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		decode.IEINASMessageContainer,
	)
	assert.NoError(t, err)
	assert.Equal(t, 6, ieLength)
	assert.Equal(t, nasMessageContainer, decodedNASMessageContainer)
}

func TestEncodeDecodeVariableLengthIE(t *testing.T) {
	vlrName := "www.facebook.com"
	encodedVLRName, err := encodeIE.EncodeVariableLengthIE(decode.IEIVLRName, decode.IELengthVLRNameMin, []byte(vlrName))
	assert.NoError(t, err)
	decodedVLRName, ieLength, err := decodeIE.DecodeVariableLengthIE(
		encodedVLRName,
		decode.IELengthVLRNameMin,
		decode.IEIVLRName,
	)
	assert.NoError(t, err)
	assert.Equal(t, len(vlrName)+decode.LengthIEI+decode.LengthLengthIndicator, ieLength)
	assert.Equal(t, vlrName, string(decodedVLRName))
}
