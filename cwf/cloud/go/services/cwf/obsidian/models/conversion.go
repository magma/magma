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

	"magma/cwf/cloud/go/cwf"
	"magma/feg/cloud/go/feg"
	fegModels "magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/lte/cloud/go/lte"
	policyModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func (m *CwfNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &CwfNetwork{}
}

func (m *CwfNetwork) ToConfiguratorNetwork() configurator.Network {
	network := configurator.Network{
		ID:          string(m.ID),
		Type:        cwf.CwfNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			cwf.CwfNetworkType:          m.CarrierWifi,
			feg.FederatedNetworkType:    m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
	if m.SubscriberConfig != nil {
		network.Configs[lte.NetworkSubscriberConfigType] = m.SubscriberConfig
	}
	return network
}

func (m *CwfNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	update := configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			cwf.CwfNetworkType:          m.CarrierWifi,
			feg.FederatedNetworkType:    m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
	if m.SubscriberConfig != nil {
		update.ConfigsToAddOrUpdate[lte.NetworkSubscriberConfigType] = m.SubscriberConfig
	}
	return update
}

func (m *CwfNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[cwf.CwfNetworkType]; cfg != nil {
		m.CarrierWifi = cfg.(*NetworkCarrierWifiConfigs)
	}
	if cfg := n.Configs[feg.FederatedNetworkType]; cfg != nil {
		m.Federation = cfg.(*fegModels.FederatedNetworkConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*orc8rModels.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8rModels.NetworkFeatures)
	}
	if cfg := n.Configs[lte.NetworkSubscriberConfigType]; cfg != nil {
		m.SubscriberConfig = cfg.(*policyModels.NetworkSubscriberConfig)
	}
	return m
}

func (m *CwfGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *CwfGateway) FromBackendModels(
	magmadGateway, cwfGateway configurator.NetworkEntity,
	device *orc8rModels.GatewayDevice,
	status *orc8rModels.GatewayStatus,
) handlers.GatewayModel {
	mdGW := (&orc8rModels.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status
	if cwfGateway.Config != nil {
		m.CarrierWifi = cwfGateway.Config.(*GatewayCwfConfigs)
	}

	return m
}

func (m *MutableCwfGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableCwfGateway) GetMagmadGateway() *orc8rModels.MagmadGateway {
	return &orc8rModels.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableCwfGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	return []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
			Type:        cwf.CwfGatewayType,
			Key:         string(m.ID),
			Name:        string(m.Name),
			Description: string(m.Description),
			Config:      m.CarrierWifi,
		},
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: cwf.CwfGatewayType, Key: string(m.ID)}},
		},
	}
}

func (m *MutableCwfGateway) GetGatewayType() string {
	return cwf.CwfGatewayType
}

func (m *MutableCwfGateway) GetAdditionalLoadsOnLoad(gateway configurator.NetworkEntity) storage.TKs {
	return nil
}

func (m *MutableCwfGateway) GetAdditionalLoadsOnUpdate() storage.TKs {
	return []storage.TypeAndKey{{Type: cwf.CwfGatewayType, Key: string(m.ID)}}
}

func (m *MutableCwfGateway) GetAdditionalWritesOnUpdate(
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	var ret []configurator.EntityWriteOperation
	existingEnt, ok := loadedEntities[storage.TypeAndKey{Type: cwf.CwfGatewayType, Key: string(m.ID)}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	entUpdate := configurator.EntityUpdateCriteria{
		Type:      cwf.CwfGatewayType,
		Key:       string(m.ID),
		NewConfig: m.CarrierWifi,
	}
	if string(m.Name) != existingEnt.Name {
		entUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingEnt.Description {
		entUpdate.NewDescription = swag.String(string(m.Description))
	}

	ret = append(ret, entUpdate)
	return ret, nil
}

func (m *GatewayCwfConfigs) FromBackendModels(networkID string, gatewayID string) error {
	carrierWifi, err := configurator.LoadEntityConfig(networkID, cwf.CwfGatewayType, gatewayID, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *carrierWifi.(*GatewayCwfConfigs)
	return nil
}

func (m *GatewayCwfConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: cwf.CwfGatewayType, Key: gatewayID,
			NewConfig: m,
		},
	}, nil
}

func (m *NetworkCarrierWifiConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, cwf.CwfNetworkType, m), nil
}

func (m *NetworkCarrierWifiConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8rModels.GetNetworkConfig(network, cwf.CwfNetworkType)
}

func (m *LiUes) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	networkConfig := orc8rModels.GetNetworkConfig(network, cwf.CwfNetworkType)
	if networkConfig == nil {
		return configurator.NetworkUpdateCriteria{}, merrors.ErrNotFound
	}
	networkConfig.(*NetworkCarrierWifiConfigs).LiUes = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, cwf.CwfNetworkType, networkConfig), nil
}

func (m *LiUes) GetFromNetwork(network configurator.Network) interface{} {
	networkConfig := orc8rModels.GetNetworkConfig(network, cwf.CwfNetworkType)
	if networkConfig == nil {
		return nil
	}
	return networkConfig.(*NetworkCarrierWifiConfigs).LiUes
}

func (m *LiUes) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *CwfHaPair) ToEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type:   cwf.CwfHAPairType,
		Key:    m.HaPairID,
		Config: m.Config,
		Associations: []storage.TypeAndKey{
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID1,
			},
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID2,
			},
		},
	}
}

func (m *CwfHaPair) FromBackendModels(ent configurator.NetworkEntity) error {
	gatewayIDs := []string{}
	for _, assoc := range ent.Associations {
		if assoc.Type == cwf.CwfGatewayType {
			gatewayIDs = append(gatewayIDs, assoc.Key)
		}
	}
	if len(gatewayIDs) != 2 {
		return fmt.Errorf("could not convert entity to CwfHaPair; could not parse gateway pair IDs")
	}
	if ent.Config == nil {
		return fmt.Errorf("could not convert entity to CwfHaPair; config was nil")
	}
	cfg, ok := ent.Config.(*CwfHaPairConfigs)
	if !ok {
		return fmt.Errorf("could not convert entity config type %T to CwfHaPair", ent.Config)
	}
	m.HaPairID = ent.Key
	m.GatewayID1 = gatewayIDs[0]
	m.GatewayID2 = gatewayIDs[1]
	m.Config = cfg
	return nil
}

func (m *CwfHaPair) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:      cwf.CwfHAPairType,
		Key:       m.HaPairID,
		NewConfig: m.Config,
		AssociationsToSet: []storage.TypeAndKey{
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID1,
			},
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID2,
			},
		},
	}
	return ret
}

func (m *MutableCwfHaPair) ToEntityUpdateCriteria(haPairID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:      cwf.CwfHAPairType,
		Key:       haPairID,
		NewConfig: m.Config,
		AssociationsToSet: []storage.TypeAndKey{
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID1,
			},
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID2,
			},
		},
	}
	return ret
}

func (m *MutableCwfHaPair) ToEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type:   cwf.CwfHAPairType,
		Key:    m.HaPairID,
		Config: m.Config,
		Associations: []storage.TypeAndKey{
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID1,
			},
			{
				Type: cwf.CwfGatewayType,
				Key:  m.GatewayID2,
			},
		},
	}
}
