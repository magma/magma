/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"errors"
	"fmt"
	"net"
	"regexp"

	"magma/lte/cloud/go/services/cellular/utils"

	"github.com/golang/glog"
)

var mccRe = regexp.MustCompile("^[0-9]{3}$")
var mncRe = regexp.MustCompile("^[0-9]{2,3}$")

func (m *GatewayCellularConfigs) ValidateGatewayConfig() error {
	if m == nil {
		return errors.New("Gateway config is nil")
	}
	if err := m.Ran.validateGatewayRANConfig(); err != nil {
		glog.Errorf("error : %v", err)
		return err
	}
	if err := m.Epc.validateGatewayEPCConfig(); err != nil {
		return err
	}
	return nil
}

func (m *GatewayRanConfigs) validateGatewayRANConfig() error {
	if m == nil {
		return errors.New("Gateway RAN config is nil")
	}
	return nil
}

func (m *GatewayEpcConfigs) validateGatewayEPCConfig() error {
	if m == nil {
		return errors.New("Gateway EPC config is nil")
	}
	if m.IPBlock != "" {
		_, _, err := net.ParseCIDR(m.IPBlock)
		if err != nil {
			return fmt.Errorf("Invalid IP block: %s", err)
		}
	}
	return nil
}

func (m *NetworkCellularConfigs) ValidateNetworkConfig() error {
	if m == nil {
		return errors.New("Network config is nil")
	}
	if err := m.Ran.validateNetworkRANConfig(); err != nil {
		return err
	}
	if err := m.Epc.validateNetworkEPCConfig(); err != nil {
		return err
	}
	return nil
}

func (m *NetworkRanConfigs) validateNetworkRANConfig() error {
	if m == nil {
		return errors.New("Network RAN config is nil")
	}

	// TODO: after data migration and corresponding update to partner portal,
	// we can enforce that exactly one is always set (i.e. none set is invalid)
	fddConfigSet := m.FddConfig != nil
	tddConfigSet := m.TddConfig != nil
	if fddConfigSet && tddConfigSet {
		return errors.New("Only one of TDD or FDD configs can be set")
	}

	earfcnDl := getEarfcnDl(m)
	band, err := utils.GetBand(uint32(earfcnDl))
	if err != nil {
		return err
	}

	if err := validateFDDConfig(earfcnDl, band, m.FddConfig); err != nil {
		return err
	}
	if err := validateTDDConfig(band, m.TddConfig); err != nil {
		return err
	}

	return nil
}

func (m *NetworkEpcConfigs) validateNetworkEPCConfig() error {
	if m == nil {
		return errors.New("Network EPC config is nil")
	}
	if !mccRe.MatchString(m.Mcc) {
		return errors.New("MCC must be in the form of a 3-digit number (leading 0's are allowed).")
	}
	if !mncRe.MatchString(m.Mnc) {
		return errors.New("MNC must be in the form of a 2- or 3-digit number (leading 0's are allowed).")
	}
	tac := m.Tac
	if tac < 1 || tac > 65535 {
		return errors.New("TAC must be between 1 and 65535 inclusive")
	}

	if len(m.LteAuthOp) < 15 || len(m.LteAuthOp) > 16 {
		return errors.New("Auth OP must be between 15 and 16 bytes")
	}
	for name, profile := range m.SubProfiles {
		if name == "" {
			return errors.New("Profile name should be non-empty")
		}
		if profile.MaxDlBitRate == 0 || profile.MaxUlBitRate == 0 {
			return errors.New("Bit rate should be greater than 0")
		}
	}
	return nil
}

func getEarfcnDl(config *NetworkRanConfigs) int32 {
	if config.FddConfig != nil {
		return int32(config.FddConfig.Earfcndl)
	} else if config.TddConfig != nil {
		return int32(config.TddConfig.Earfcndl)
	} else {
		// TODO: after migration, nix this else
		return int32(config.Earfcndl)
	}
}

func validateFDDConfig(earfcnDl int32, band *utils.LTEBand, fddConfig *NetworkRanConfigsFddConfig) error {
	if fddConfig == nil {
		return nil
	}

	if band.Mode != utils.FDDMode {
		return fmt.Errorf("Not a FDD Band: %d", band.ID)
	}
	earfcnUl := fddConfig.Earfcnul
	// Provide default EARFCNUL if not set
	if earfcnUl == 0 {
		fddConfig.Earfcnul = uint32(earfcnDl) - band.StartEarfcnDl + band.StartEarfcnUl
		earfcnUl = fddConfig.Earfcnul
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

func validateTDDConfig(band *utils.LTEBand, tddConfig *NetworkRanConfigsTddConfig) error {
	if tddConfig == nil {
		return nil
	}

	if band.Mode != utils.TDDMode {
		return fmt.Errorf("Not a TDD Band: %d", band.ID)
	}
	return nil
}

func (config *NetworkEnodebConfigs) ValidateEnodebConfig() error {
	if config == nil {
		return errors.New("Gateway config is nil")
	}
	if config.Earfcndl < 0 || config.Earfcndl > 65535 {
		return errors.New("EARFCNDL must be within 0-65535")
	}
	if config.SubframeAssignment < 0 || config.SubframeAssignment > 6 {
		return errors.New("Subframe assignment must be within 0-6")
	}
	if config.SpecialSubframePattern < 0 || config.SpecialSubframePattern > 9 {
		return errors.New("Special subframe pattern must be within 0-9")
	}
	if config.Pci < 0 || config.Pci > 504 {
		return errors.New("PCI must be within 0-504")
	}
	if config.CellID < 0 || config.CellID > 268435455 {
		return errors.New("Cell ID must be within 0-268435455")
	}
	if config.Tac < 0 || config.Tac > 65535 {
		return errors.New("TAC must be within 0-65535")
	}
	switch config.DeviceClass {
	case
		"Baicells Nova-233 G2 OD FDD",
		"Baicells Nova-243 OD TDD",
		"Baicells Neutrino 224 ID FDD",
		"Baicells ID TDD/FDD",
		"NuRAN Cavium OC-LTE":
		break
	default:
		return errors.New("Invalid eNodeB device class")
	}
	switch config.BandwidthMhz {
	case
		3,
		5,
		10,
		15,
		20:
		break
	default:
		return errors.New("Invalid eNodeB bandwidth option")
	}
	return nil
}
