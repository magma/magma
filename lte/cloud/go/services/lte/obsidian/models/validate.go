/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"fmt"
	"net"

	"magma/lte/cloud/go/lte"
	"magma/orc8r/cloud/go/services/configurator"

	oerrors "github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

func (m *LteNetwork) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	var res []error
	if err := m.Cellular.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if err := m.DNS.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if err := m.Features.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if m.SubscriberConfig != nil {
		if err := m.SubscriberConfig.ValidateModel(); err != nil {
			res = append(res, err)
		}
	}

	if len(res) > 0 {
		return oerrors.CompositeValidationError(res...)
	}
	return nil
}

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
	band, err := lte.GetBand(earfcnDl)
	if err != nil {
		return err
	}

	if tddConfigSet && band.Mode != lte.TDDMode {
		return errors.Errorf("band %d not a TDD band", band.ID)
	}
	if fddConfigSet {
		if band.Mode != lte.FDDMode {
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

func (m *LteGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableLteGateway) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	// Custom validation only for cellular and device
	var res []error
	if err := m.Cellular.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if err := m.Device.ValidateModel(); err != nil {
		res = append(res, err)
	}

	resourceIDs := map[string]struct{}{}
	for apnName, resource := range m.ApnResources {
		if apnName != string(resource.ApnName) {
			return fmt.Errorf("APN resources key (%s) and APN name (%s) must match", apnName, resource.ApnName)
		}
		if _, ok := resourceIDs[resource.ID]; ok {
			return fmt.Errorf("duplicate APN resource ID in request: %s", resource.ID)
		}
		resourceIDs[resource.ID] = struct{}{}
	}

	if len(res) > 0 {
		return oerrors.CompositeValidationError(res...)
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
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	if m.DNSPrimary != "" {
		ip := net.ParseIP(m.DNSPrimary)
		if ip == nil {
			return errors.New("Invalid primary DNS address")
		} else if ip.To4() == nil {
			return errors.New("Only IPv4 is supported currently for DNS")
		}
	}

	if m.DNSSecondary != "" {
		secIp := net.ParseIP(m.DNSSecondary)
		if secIp == nil {
			return errors.New("Invalid secondary DNS address")
		} else if secIp.To4() == nil {
			return errors.New("Only IPv4 is supported currently for DNS")
		}
	}

	if m.IPV6DNSAddr != "" {
		ip := net.ParseIP(string(m.IPV6DNSAddr))
		if ip == nil {
			return errors.New("Invalid IPV6 DNS address")
		}
	}
	return nil
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
	if err := validate.RequiredString("csfb_mcc", "body", m.CsfbMcc); err != nil {
		res = append(res, err)
	}
	if err := validate.RequiredString("csfb_mnc", "body", m.CsfbMnc); err != nil {
		res = append(res, err)
	}
	if err := validate.Required("csfb_rat", "body", m.CsfbRat); err != nil {
		res = append(res, err)
	}
	if err := validate.Required("lac", "body", m.Lac); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return oerrors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GatewayDNSConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayDNSRecords) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *EnodebSerials) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayHeConfig) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *Enodeb) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	if m.EnodebConfig != nil {
		if err := m.EnodebConfig.validateEnodebConfig(); err != nil {
			return err
		}
	}

	return nil
}

func (m *EnodebConfiguration) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *UnmanagedEnodebConfiguration) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *EnodebConfig) validateEnodebConfig() error {
	managedConfigSet := m.ManagedConfig != nil
	unmanagedConfigSet := m.UnmanagedConfig != nil

	if managedConfigSet && unmanagedConfigSet {
		return errors.New("only one of the eNodeb config types can be set")
	}

	if managedConfigSet {
		if m.ConfigType != ManagedConfigType {
			return errors.New("invalid type set for managed config")
		}

	}
	if unmanagedConfigSet {
		if m.ConfigType != UnmanagedConfigType {
			return errors.New("invalid type set for unmanaged config")
		}
	}

	return nil
}

func (m *EnodebState) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *Apn) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *CellularGatewayPool) ValidateModel() error {
	err := m.Validate(strfmt.Default)
	if err != nil {
		return err
	}
	return m.Config.ValidateModel()
}

func (m *CellularGatewayPoolConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableCellularGatewayPool) ValidateModel() error {
	err := m.Validate(strfmt.Default)
	if err != nil {
		return err
	}
	return m.Config.ValidateModel()
}

func (m *CellularGatewayPoolRecords) ValidateModel() error {
	err := m.Validate(strfmt.Default)
	if err != nil {
		return err
	}
	uniquePool := make(map[GatewayPoolID]bool, len(*m))
	for _, record := range *m {
		if !uniquePool[record.GatewayPoolID] {
			uniquePool[record.GatewayPoolID] = true
		} else {
			return fmt.Errorf("All pool records must have unique pool IDs")
		}
	}
	if len(*m) == 0 {
		return nil
	}
	relCapacity := (*m)[0].MmeRelativeCapacity
	mmeCode := (*m)[0].MmeCode
	for _, record := range *m {
		if record.MmeRelativeCapacity != relCapacity {
			return fmt.Errorf("Setting different MME relative capacities for the same gateway is currently unsupported")
		}
		if record.MmeCode != mmeCode {
			return fmt.Errorf("Setting different MME codes for the same gateway is currently unsupported")
		}
	}
	return nil
}
