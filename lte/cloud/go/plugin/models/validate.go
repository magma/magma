/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"fmt"
	"net"

	"magma/lte/cloud/go/services/cellular/utils"
	"magma/orc8r/cloud/go/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"

	errors2 "github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

func (m *NetworkCellularConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if err := m.FegNetworkID.ValidateModel(); err != nil {
		return err
	}
	if err := m.Epc.ValidateModel(); err != nil {
		return err
	}
	if err := m.Ran.ValidateModel(); err != nil {
		return err
	}
	return nil
}

func (m FegNetworkID) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if !swag.IsZero(m) {
		exists, err := configurator.DoesNetworkExist(string(m))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to search for network %s", string(m)))
		}
		if !exists {
			return errors.New(fmt.Sprintf("Network: %s does not exist", string(m)))
		}
	}
	return nil
}

func (m *NetworkEpcConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	if m.Mobility != nil {
		if err := m.Mobility.validateMobility(); err != nil {
			return err
		}
	}

	for name := range m.SubProfiles {
		if name == "" {
			return errors.New("profile name should be non-empty")
		}
	}
	return nil
}

func (m *NetworkRanConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	tddConfigSet := m.TddConfig != nil
	fddConfigSet := m.FddConfig != nil

	if tddConfigSet && fddConfigSet {
		return errors.New("only one of TDD or FDD configs can be set")
	} else if !tddConfigSet && !fddConfigSet {
		return errors.New("either TDD or FDD configs must be set")
	}

	earfcnDl := m.getEarfcnDl()
	band, err := utils.GetBand(earfcnDl)
	if err != nil {
		return err
	}

	if tddConfigSet && band.Mode != utils.TDDMode {
		return errors.Errorf("band %d not a TDD band", band.ID)
	}
	if fddConfigSet {
		if band.Mode != utils.FDDMode {
			return errors.Errorf("band %d not a FDD band", band.ID)
		}
		if !band.EarfcnULInRange(m.FddConfig.Earfcnul) {
			return errors.Errorf("EARFCNUL=%d invalid for band %d (%d, %d)", m.FddConfig.Earfcnul, band.ID, band.StartEarfcnUl, band.StartEarfcnDl)
		}
	}

	return nil
}

func (m *NetworkRanConfigs) getEarfcnDl() uint32 {
	if m.TddConfig != nil {
		return m.TddConfig.Earfcndl
	}
	if m.FddConfig != nil {
		return m.FddConfig.Earfcndl
	}
	// This should truly be unreachable
	return 0
}

func (m *NetworkEpcConfigsMobility) validateMobility() error {
	mobilityNatConfigSet := m.Nat != nil
	mobilityStaticConfigSet := m.Static != nil
	// TODO: Add validation for DHCP once is added to EPC config

	if mobilityNatConfigSet && mobilityStaticConfigSet {
		return errors.New("only one of the mobility IP allocation modes can be set")
	}

	if mobilityNatConfigSet {
		if m.IPAllocationMode != NATAllocationMode {
			return errors.New("invalid config set for NAT allocation mode")
		}

		if err := validateIPBlocks(m.Nat.IPBlocks); err != nil {
			return errors.New("invalid IP block on config")
		}
	}
	if mobilityStaticConfigSet {
		if m.IPAllocationMode != StaticAllocationMode {
			return errors.New("invalid config set for STATIC allocation mode")
		}

		for _, ipBlocks := range m.Static.IPBlocksByTac {
			if err := validateIPBlocks(ipBlocks); err != nil {
				return errors.New("invalid IP block on config")
			}
		}
	}

	return nil
}

func (m *GatewayCellularConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	// Custom validation only exists for the non EPS configs, EPC and RAN
	// validation are handled by the above call the Validate()
	if m.NonEpsService == nil {
		return nil
	}
	if err := m.NonEpsService.ValidateModel(); err != nil {
		return err
	}
	return nil
}

func (m *GatewayRanConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayEpcConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayNonEpsConfigs) ValidateModel() error {
	// Don't validate sub-fields if Non-EPS control is off
	if swag.Uint32Value(m.NonEpsServiceControl) == 0 {
		return nil
	}

	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	// Sub-fields are actually required if non-EPS control is on - swagger
	// doesn't support dependent validation, so we'll do this in code here

	var res []error
	if err := validate.Required("arfcn_2g", "body", m.Arfcn2g); err != nil {
		res = append(res, err)
	}
	if err := validate.RequiredString("csfb_mcc", "body", string(m.CsfbMcc)); err != nil {
		res = append(res, err)
	}
	if err := validate.RequiredString("csfb_mnc", "body", string(m.CsfbMnc)); err != nil {
		res = append(res, err)
	}
	if err := validate.Required("csfb_rat", "body", m.CsfbRat); err != nil {
		res = append(res, err)
	}
	if err := validate.Required("lac", "body", m.Lac); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors2.CompositeValidationError(res...)
	}
	return nil
}

func (m *EnodebSerials) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *Enodeb) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *EnodebConfiguration) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

const (
	lteAuthKeyLength = 16
	lteAuthOpcLength = 16
)

func (m *Subscriber) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if err := m.Lte.ValidateModel(); err != nil {
		return err
	}
	return nil
}

func (m *LteSubscription) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	authKeyLen := len([]byte(m.AuthKey))
	if authKeyLen != lteAuthKeyLength {
		return models.ValidateErrorf("expected lte auth key to be %d bytes but got %d bytes", lteAuthKeyLength, authKeyLen)
	}

	// OPc is optional, but if it's provided it should be 16 bytes
	authOpcLen := len([]byte(m.AuthOpc))
	if authOpcLen > 0 && authOpcLen != lteAuthOpcLength {
		return models.ValidateErrorf("expected lte auth opc to be %d bytes but got %d bytes", lteAuthOpcLength, authOpcLen)
	}

	return nil
}

func (m *EnodebState) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

// validateIPBlocks parses and validates IP networks containing subnet masks.
// Returns an error in case any IP network in list is invalid.
func validateIPBlocks(ipBlocks []string) error {
	for _, ipBlock := range ipBlocks {
		_, _, err := net.ParseCIDR(ipBlock)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateModel does standard swagger validation and any custom validation
func (m *PolicyRule) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

// ValidateModel does standard swagger validation and any custom validation
func (m *RatingGroup) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

// ValidateModel does standard swagger validation and any custom validation
func (m *MutableRatingGroup) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
