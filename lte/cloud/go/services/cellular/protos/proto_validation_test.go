/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos_test

import (
	"testing"

	"magma/lte/cloud/go/services/cellular/protos"
	"magma/lte/cloud/go/services/cellular/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestValidateGatewayConfig(t *testing.T) {
	config := test_utils.NewDefaultGatewayConfig()
	err := protos.ValidateGatewayConfig(config)
	assert.NoError(t, err)

	// Test nils
	config.Epc = nil
	err = protos.ValidateGatewayConfig(config)
	assert.Error(t, err, "Gateway EPC config is nil")

	config = test_utils.NewDefaultGatewayConfig()
	config.Ran = nil
	err = protos.ValidateGatewayConfig(config)
	assert.Error(t, err, "Gateway RAN config is nil")

	err = protos.ValidateGatewayConfig(nil)
	assert.Error(t, err, "Gateway config is nil")

	// IP block parsing
	config = test_utils.NewDefaultGatewayConfig()
	config.Epc.IpBlock = "20.20.20.0/24"
	assert.NoError(t, protos.ValidateGatewayConfig(config))
	config.Epc.IpBlock = "12345"
	assert.Error(t, protos.ValidateGatewayConfig(config))
}

func TestValidateNetworkConfig(t *testing.T) {
	config := test_utils.NewDefaultTDDNetworkConfig()
	err := protos.ValidateNetworkConfig(config)
	assert.NoError(t, err)

	positiveTestCases := []struct {
		mcc string
		mnc string
	}{
		// MNC test cases
		{"123", "001"},
		{"123", "01"},
		{"123", "123"},
		{"123", "12"},
		// MCC test cases
		{"001", "123"},
		{"012", "123"},
		{"000", "123"},
	}

	negativeTestCases := []struct {
		mcc                  string
		mnc                  string
		expectedErrorMessage string
	}{
		// MNC test cases
		{"123", "1", "MNC must be in the form of a 2- or 3-digit number (leading 0's are allowed)."},
		{"123", "0", "MNC must be in the form of a 2- or 3-digit number (leading 0's are allowed)."},
		{"123", "ab", "MNC must be in the form of a 2- or 3-digit number (leading 0's are allowed)."},
		// MCC test cases
		{"01", "123", "MCC must be in the form of a 3-digit number (leading 0's are allowed)."},
		{"12", "123", "MCC must be in the form of a 3-digit number (leading 0's are allowed)."},
		{"0", "123", "MCC must be in the form of a 3-digit number (leading 0's are allowed)."},
		{"1", "123", "MCC must be in the form of a 3-digit number (leading 0's are allowed)."},
		{"abc", "123", "MCC must be in the form of a 3-digit number (leading 0's are allowed)."},
	}

	for _, testCase := range positiveTestCases {
		config.Epc.Mcc = testCase.mcc
		config.Epc.Mnc = testCase.mnc
		err := protos.ValidateNetworkConfig(config)
		assert.NoError(t, err)
	}

	for _, testCase := range negativeTestCases {
		config.Epc.Mcc = testCase.mcc
		config.Epc.Mnc = testCase.mnc
		err := protos.ValidateNetworkConfig(config)
		assert.Error(t, err)
		assert.Equal(t, testCase.expectedErrorMessage, err.Error())
	}
}

func TestValidateNetworkRANConfig(t *testing.T) {
	config := test_utils.NewDefaultTDDNetworkConfig()
	tddConf := &protos.NetworkRANConfig_TDDConfig{
		Earfcndl:               43950,
		SubframeAssignment:     2,
		SpecialSubframePattern: 7,
	}

	fddConf := &protos.NetworkRANConfig_FDDConfig{
		Earfcndl: 0,
		Earfcnul: 18000,
	}

	positiveTestCases := []struct {
		tdd *protos.NetworkRANConfig_TDDConfig
		fdd *protos.NetworkRANConfig_FDDConfig
	}{
		{tddConf, nil},
		{nil, fddConf},
		{nil, nil},
	}

	negativeTestCases := []struct {
		tdd                  *protos.NetworkRANConfig_TDDConfig
		fdd                  *protos.NetworkRANConfig_FDDConfig
		expectedErrorMessage string
	}{
		// TODO: ensure that one must be set after migration
		{tddConf, fddConf, "Only one of TDD or FDD configs can be set"},
	}

	for _, testCase := range positiveTestCases {
		config.Ran.FddConfig = testCase.fdd
		config.Ran.TddConfig = testCase.tdd
		err := protos.ValidateNetworkConfig(config)
		assert.NoError(t, err)
	}

	for _, testCase := range negativeTestCases {
		config.Ran.FddConfig = testCase.fdd
		config.Ran.TddConfig = testCase.tdd
		err := protos.ValidateNetworkConfig(config)
		assert.Error(t, err)
		assert.Equal(t, testCase.expectedErrorMessage, err.Error())
	}
}

func TestValidateNetworkTDDRANConfig(t *testing.T) {
	config := test_utils.NewDefaultTDDNetworkConfig()

	positiveTestCases := []struct {
		earfcndl int32
	}{
		// TDD EARFCNDLs
		{43950},
	}

	negativeTestCases := []struct {
		earfcndl             int32
		expectedErrorMessage string
	}{
		// FDD EARFCNDLs
		{0, "Not a TDD Band: 1"},
		{600, "Not a TDD Band: 2"},
		{1200, "Not a TDD Band: 3"},
		{-1, "Invalid EARFCNDL: no matching band"},
	}

	config.Ran.TddConfig = &protos.NetworkRANConfig_TDDConfig{
		SubframeAssignment:     2,
		SpecialSubframePattern: 7,
	}

	for _, testCase := range positiveTestCases {
		config.Ran.TddConfig.Earfcndl = testCase.earfcndl
		err := protos.ValidateNetworkConfig(config)
		assert.NoError(t, err)
	}

	for _, testCase := range negativeTestCases {
		config.Ran.TddConfig.Earfcndl = testCase.earfcndl
		err := protos.ValidateNetworkConfig(config)
		assert.Error(t, err)
		assert.Equal(t, testCase.expectedErrorMessage, err.Error())
	}
}

func TestValidateNetworkFDDRANConfig(t *testing.T) {
	config := test_utils.NewDefaultTDDNetworkConfig()

	positiveTestCases := []struct {
		earfcndl int32
		earfcnul int32
	}{
		// FDD EARFCNDLs
		{0, 0},
		{0, 18000},
		{1, 0},
		{1, 18001},
	}

	negativeTestCases := []struct {
		earfcndl             int32
		earfcnul             int32
		expectedErrorMessage string
	}{
		{0, 17999, "EARFCNUL=17999 invalid for Band 1 (18000, 18600)"},
		{0, 18600, "EARFCNUL=18600 invalid for Band 1 (18000, 18600)"},
		{43950, 43950, "Not a FDD Band: 43"},
	}

	config.Ran.TddConfig = nil
	config.Ran.FddConfig = &protos.NetworkRANConfig_FDDConfig{}
	for _, testCase := range positiveTestCases {
		config.Ran.FddConfig.Earfcndl = testCase.earfcndl
		config.Ran.FddConfig.Earfcnul = testCase.earfcnul
		err := protos.ValidateNetworkConfig(config)
		assert.NoError(t, err)
	}

	for _, testCase := range negativeTestCases {
		config.Ran.FddConfig.Earfcndl = testCase.earfcndl
		config.Ran.FddConfig.Earfcnul = testCase.earfcnul
		err := protos.ValidateNetworkConfig(config)
		assert.Error(t, err)
		assert.Equal(t, testCase.expectedErrorMessage, err.Error())
	}
}

func TestValidateSubProfile(t *testing.T) {
	config := test_utils.NewDefaultTDDNetworkConfig()
	config.Epc.SubProfiles = make(
		map[string]*protos.NetworkEPCConfig_SubscriptionProfile)
	config.Epc.SubProfiles["test"] = &protos.NetworkEPCConfig_SubscriptionProfile{
		MaxUlBitRate: 0, MaxDlBitRate: 0,
	}
	err := protos.ValidateNetworkConfig(config)
	assert.Error(t, err)

	config.Epc.SubProfiles["test"] = &protos.NetworkEPCConfig_SubscriptionProfile{
		MaxUlBitRate: 100, MaxDlBitRate: 100,
	}
	err = protos.ValidateNetworkConfig(config)
	assert.NoError(t, err)

	config.Epc.SubProfiles[""] = &protos.NetworkEPCConfig_SubscriptionProfile{
		MaxUlBitRate: 100, MaxDlBitRate: 100,
	}
	err = protos.ValidateNetworkConfig(config)
	assert.Error(t, err)
}

func TestValidateEnodebConfig(t *testing.T) {
	config := test_utils.NewDefaultEnodebConfig()
	err := protos.ValidateEnodebConfig(config)
	assert.NoError(t, err)

	config = test_utils.NewDefaultEnodebConfig()
	config.Earfcndl = -1
	err = protos.ValidateEnodebConfig(config)
	assert.Error(t, err)

	config = test_utils.NewDefaultEnodebConfig()
	config.SubframeAssignment = 7
	err = protos.ValidateEnodebConfig(config)
	assert.Error(t, err)

	config = test_utils.NewDefaultEnodebConfig()
	config.SpecialSubframePattern = 10
	err = protos.ValidateEnodebConfig(config)
	assert.Error(t, err)

	config = test_utils.NewDefaultEnodebConfig()
	config.Pci = 505
	err = protos.ValidateEnodebConfig(config)
	assert.Error(t, err)

	config = test_utils.NewDefaultEnodebConfig()
	config.DeviceClass = "Some unsupported device"
	err = protos.ValidateEnodebConfig(config)
	assert.Error(t, err)
}
