/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package migration_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	lteplugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/services/cellular/protos"
	"magma/lte/cloud/go/tools/migrations/m003_fdd_tdd_configs/migration"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_service "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tddNetwork = "tdd_test_network"
	fddNetwork = "fdd_test_network"

	// values taken from lte/services/cellular/test_utils/defaults.go
	tddEarfcndl               int32 = 44590
	tddBandwidthMhz           int32 = 20
	tddSubframeAssignment     int32 = 2
	tddSpecialSubframePattern int32 = 7

	fddEarfcndl     int32 = 1
	fddEarfcnul     int32 = 18001
	fddBandwidthMhz int32 = 20
)

func TestMigrateFddTddConfigs(t *testing.T) {
	// setup networks
	setup(t)

	// perform migration
	err := migration.Migrate()
	require.NoError(t, err)

	// validate
	tddConf, err := config.GetConfig(tddNetwork, lte.CellularNetworkType, tddNetwork)
	require.NoError(t, err, "Failed to get tdd config")
	fddConf, err := config.GetConfig(fddNetwork, lte.CellularNetworkType, fddNetwork)
	require.NoError(t, err, "Failed to get fdd config")

	tddNetConf := tddConf.(*protos.CellularNetworkConfig)
	fddNetConf := fddConf.(*protos.CellularNetworkConfig)

	tdd := tddNetConf.Ran
	fdd := fddNetConf.Ran

	// validate tdd
	err = protos.ValidateNetworkConfig(tddNetConf)
	assert.NoError(t, err)

	assert.Equal(t, tdd.Earfcndl, tdd.TddConfig.Earfcndl)

	assert.Equal(t, tddBandwidthMhz, tdd.BandwidthMhz)
	assert.Equal(t, tddEarfcndl, tdd.TddConfig.Earfcndl)
	assert.Equal(t, tddSubframeAssignment, tdd.TddConfig.SubframeAssignment)
	assert.Equal(t, tddSpecialSubframePattern, tdd.TddConfig.SpecialSubframePattern)

	// validate fdd
	err = protos.ValidateNetworkConfig(fddNetConf)
	assert.NoError(t, err)

	assert.Equal(t, fdd.Earfcndl, fdd.FddConfig.Earfcndl)

	assert.Equal(t, fddBandwidthMhz, fdd.BandwidthMhz)
	assert.Equal(t, fddEarfcndl, fdd.FddConfig.Earfcndl)
	assert.Equal(t, fddEarfcnul, fdd.FddConfig.Earfcnul)
}

func setup(t *testing.T) {
	// start magma test service so nothing gets broken in prod
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	magmad_test_service.StartTestService(t)

	// setup test networks
	_, err := magmad.RegisterNetwork(&magmad_protos.MagmadNetworkRecord{Name: "TDD Test Network"}, tddNetwork)
	require.NoError(t, err)
	_, err = magmad.RegisterNetwork(&magmad_protos.MagmadNetworkRecord{Name: "FDD Test Network"}, fddNetwork)
	require.NoError(t, err)

	// add one TDD/FDD config each
	tdd := protos.NetworkRANConfig{
		Earfcndl:               tddEarfcndl,
		BandwidthMhz:           tddBandwidthMhz,
		SubframeAssignment:     tddSubframeAssignment,
		SpecialSubframePattern: tddSpecialSubframePattern,
	}
	fdd := protos.NetworkRANConfig{
		Earfcndl:     fddEarfcndl,
		BandwidthMhz: fddBandwidthMhz,
	}

	epc := &protos.NetworkEPCConfig{
		Mcc: "001",
		Mnc: "01",
		Tac: 1,
		// 16 bytes of \x11
		LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
		LteAuthAmf: []byte("\x80\x00"),
	}

	tddNetConfig := &protos.CellularNetworkConfig{Ran: &tdd, Epc: epc}
	fddNetConfig := &protos.CellularNetworkConfig{Ran: &fdd, Epc: epc}

	err = config.CreateConfig(tddNetwork, lte.CellularNetworkType, tddNetwork, tddNetConfig)
	require.NoError(t, err)
	err = config.CreateConfig(fddNetwork, lte.CellularNetworkType, fddNetwork, fddNetConfig)
	require.NoError(t, err)
}
