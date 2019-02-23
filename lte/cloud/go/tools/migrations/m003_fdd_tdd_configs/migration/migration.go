/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package migration

import (
	cellular_config "magma/lte/cloud/go/services/cellular/config"
	"magma/lte/cloud/go/services/cellular/protos"
	"magma/lte/cloud/go/services/cellular/utils"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/magmad"
)

/**
Performs the a migration from the old format of NetworkRANConfig (where
everything was assumed to be TDD)  to the new format with supports FDD style
networks as well.

After this migration is complete the only valid fields in NetworkRANConfig
will be the .Earfcndl, .BandwidthMhz, and either (.FddConfig.Earfcndl,
.FddConfig.Earfcnul) or (.TddConfig.Earfcndl, .TddConfig.SubframeAssignment,
.TddConfig.SpecialSubframePattern).
*/
func Migrate() error {
	networks, err := magmad.ListNetworks()
	if err != nil {
		return err
	}

	for _, network := range networks {

		conf, err := config.GetConfig(network, cellular_config.CellularNetworkType, network)
		if err != nil {
			return err
		}

		netConf := conf.(*protos.CellularNetworkConfig)
		newRan := protos.NetworkRANConfig{}

		err = convert(&newRan, netConf.Ran)
		if err != nil {
			return err
		}

		netConf.Ran = &newRan

		err = config.UpdateConfig(network, cellular_config.CellularNetworkType, network, netConf)
		if err != nil {
			return err
		}
	}

	return nil
}

func convert(ret *protos.NetworkRANConfig, config *protos.NetworkRANConfig) error {
	band, err := utils.GetBand(config.Earfcndl)
	if err != nil {
		return err
	}

	ret.Earfcndl = config.Earfcndl
	ret.BandwidthMhz = config.BandwidthMhz

	if band.Mode == utils.FDDMode {
		ret.FddConfig = new(protos.NetworkRANConfig_FDDConfig)
		ret.FddConfig.Earfcndl = config.Earfcndl
		ret.FddConfig.Earfcnul = config.Earfcndl - band.StartEarfcnDl + band.StartEarfcnUl
	}
	if band.Mode == utils.TDDMode {
		ret.TddConfig = new(protos.NetworkRANConfig_TDDConfig)
		ret.TddConfig.Earfcndl = config.Earfcndl
		ret.TddConfig.SpecialSubframePattern = config.SpecialSubframePattern
		ret.TddConfig.SubframeAssignment = config.SubframeAssignment
	}

	return nil
}
