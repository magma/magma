/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"errors"
	"fmt"
	"net"
	"regexp"

	"magma/lte/cloud/go/services/cellular/utils"
)

var mccRe = regexp.MustCompile("^[0-9]{3}$")
var mncRe = regexp.MustCompile("^[0-9]{2,3}$")

// See cellular_service.proto for full documentation on all config protobuf
// fields and which ones are optional/required.

func ValidateGatewayConfig(config *CellularGatewayConfig) error {
	if config == nil {
		return errors.New("Gateway config is nil")
	}
	if err := validateGatewayRANConfig(config.GetRan()); err != nil {
		return err
	}
	if err := validateGatewayEPCConfig(config.GetEpc()); err != nil {
		return err
	}
	return nil
}

func validateGatewayRANConfig(config *GatewayRANConfig) error {
	if config == nil {
		return errors.New("Gateway RAN config is nil")
	}
	return nil
}

func validateGatewayEPCConfig(config *GatewayEPCConfig) error {
	if config == nil {
		return errors.New("Gateway EPC config is nil")
	}
	if config.GetIpBlock() != "" {
		_, _, err := net.ParseCIDR(config.GetIpBlock())
		if err != nil {
			return fmt.Errorf("Invalid IP block: %s", err)
		}
	}
	return nil
}

func ValidateNetworkConfig(config *CellularNetworkConfig) error {
	if config == nil {
		return errors.New("Network config is nil")
	}
	if err := validateNetworkRANConfig(config.GetRan()); err != nil {
		return err
	}
	if err := validateNetworkEPCConfig(config.GetEpc()); err != nil {
		return err
	}
	return nil
}

func validateNetworkRANConfig(config *NetworkRANConfig) error {
	if config == nil {
		return errors.New("Network RAN config is nil")
	}

	// TODO: after data migration and corresponding update to partner portal,
	// we can enforce that exactly one is always set (i.e. none set is invalid)
	fddConfigSet := config.FddConfig != nil
	tddConfigSet := config.TddConfig != nil
	if fddConfigSet && tddConfigSet {
		return errors.New("Only one of TDD or FDD configs can be set")
	}

	earfcnDl := getEarfcnDl(config)
	band, err := utils.GetBand(earfcnDl)
	if err != nil {
		return err
	}

	if err := validateFDDConfig(earfcnDl, band, config.FddConfig); err != nil {
		return err
	}
	if err := validateTDDConfig(band, config.TddConfig); err != nil {
		return err
	}

	return nil
}

func validateNetworkEPCConfig(config *NetworkEPCConfig) error {
	if config == nil {
		return errors.New("Network EPC config is nil")
	}
	if !mccRe.MatchString(config.GetMcc()) {
		return errors.New("MCC must be in the form of a 3-digit number (leading 0's are allowed).")
	}
	if !mncRe.MatchString(config.GetMnc()) {
		return errors.New("MNC must be in the form of a 2- or 3-digit number (leading 0's are allowed).")
	}
	tac := config.GetTac()
	if tac < 1 || tac > 65535 {
		return errors.New("TAC must be between 1 and 65535 inclusive")
	}

	if len(config.GetLteAuthOp()) < 15 || len(config.GetLteAuthOp()) > 16 {
		return errors.New("Auth OP must be between 15 and 16 bytes")
	}
	for name, profile := range config.GetSubProfiles() {
		if name == "" {
			return errors.New("Profile name should be non-empty")
		}
		if profile.GetMaxDlBitRate() == 0 || profile.GetMaxUlBitRate() == 0 {
			return errors.New("Bit rate should be greater than 0")
		}
	}
	return nil
}

func getEarfcnDl(config *NetworkRANConfig) int32 {
	if config.FddConfig != nil {
		return config.FddConfig.Earfcndl
	} else if config.TddConfig != nil {
		return config.TddConfig.Earfcndl
	} else {
		// TODO: after migration, nix this else
		return config.Earfcndl
	}
}

func validateFDDConfig(earfcnDl int32, band *utils.LTEBand, fddConfig *NetworkRANConfig_FDDConfig) error {
	if fddConfig == nil {
		return nil
	}

	if band.Mode != utils.FDDMode {
		return fmt.Errorf("Not a FDD Band: %d", band.ID)
	}
	earfcnUl := fddConfig.GetEarfcnul()
	// Provide default EARFCNUL if not set
	if earfcnUl == 0 {
		fddConfig.Earfcnul = earfcnDl - band.StartEarfcnDl + band.StartEarfcnUl
		earfcnUl = fddConfig.GetEarfcnul()
	}
	if !band.EarfcnULInRange(earfcnUl) {
		return fmt.Errorf("EARFCNUL=%d invalid for Band %d (%d, %d)",
			earfcnUl,
			band.ID,
			band.StartEarfcnUl,
			band.StartEarfcnUl+band.CountEarfcn)
	}
	return nil
}

func validateTDDConfig(band *utils.LTEBand, tddConfig *NetworkRANConfig_TDDConfig) error {
	if tddConfig == nil {
		return nil
	}

	if band.Mode != utils.TDDMode {
		return fmt.Errorf("Not a TDD Band: %d", band.ID)
	}
	return nil
}
