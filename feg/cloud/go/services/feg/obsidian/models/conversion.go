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
	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/lte/cloud/go/lte"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	lteModels "magma/lte/cloud/go/services/lte/obsidian/models"
	policyModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func (m *FegNetwork) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *FegNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &FegNetwork{}
}

func (m *FegNetwork) ToConfiguratorNetwork() configurator.Network {
	network := configurator.Network{
		ID:          string(m.ID),
		Type:        feg.FederationNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			feg.FegNetworkType:          m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
	if m.SubscriberConfig != nil {
		network.Configs[lte.NetworkSubscriberConfigType] = m.SubscriberConfig
	}
	return network
}

func (m *FegNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	update := configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			feg.FegNetworkType:          m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
	if m.SubscriberConfig != nil {
		update.ConfigsToAddOrUpdate[lte.NetworkSubscriberConfigType] = m.SubscriberConfig
	}
	return update
}

func (m *FegNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[feg.FegNetworkType]; cfg != nil {
		m.Federation = cfg.(*NetworkFederationConfigs)
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

func (m *FegLteNetwork) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *FegLteNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &FegLteNetwork{}
}

func (m *FegLteNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        feg.FederatedLteNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: m.Cellular,
			feg.FederatedNetworkType:      m.Federation,
			orc8r.DnsdNetworkType:         m.DNS,
			orc8r.NetworkFeaturesConfig:   m.Features,
		},
	}
}

func (m *FegLteNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			lte.CellularNetworkConfigType: m.Cellular,
			feg.FederatedNetworkType:      m.Federation,
			orc8r.DnsdNetworkType:         m.DNS,
			orc8r.NetworkFeaturesConfig:   m.Features,
		},
	}
}

func (m *FegLteNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[feg.FederatedNetworkType]; cfg != nil {
		m.Federation = cfg.(*FederatedNetworkConfigs)
	}
	if cfg := n.Configs[lte.CellularNetworkConfigType]; cfg != nil {
		m.Cellular = cfg.(*lteModels.NetworkCellularConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*orc8rModels.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8rModels.NetworkFeatures)
	}
	return m
}

func (m *NetworkFederationConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8rModels.GetNetworkConfig(network, feg.FegNetworkType)
}

func (m *NetworkFederationConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, feg.FegNetworkType, m), nil
}

func (m *FederationGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *FederationGateway) FromBackendModels(
	magmadGateway, federationGateway configurator.NetworkEntity,
	device *orc8rModels.GatewayDevice,
	status *orc8rModels.GatewayStatus,
) handlers.GatewayModel {
	// delegate most of the fillin to magmad gateway struct
	mdGW := (&orc8rModels.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status
	m.Federation = federationGateway.Config.(*GatewayFederationConfigs)
	return m
}

func (m *MutableFederationGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableFederationGateway) GetMagmadGateway() *orc8rModels.MagmadGateway {
	return &orc8rModels.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableFederationGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	return []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
			Type:        feg.FegGatewayType,
			Key:         string(m.ID),
			Name:        string(m.Name),
			Description: string(m.Description),
			Config:      m.Federation,
		},
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: string(m.ID)}},
		},
	}
}

func (m *MutableFederationGateway) GetGatewayType() string {
	return feg.FegGatewayType
}

func (m *MutableFederationGateway) GetAdditionalLoadsOnLoad(gateway configurator.NetworkEntity) storage.TKs {
	return nil
}

func (m *MutableFederationGateway) GetAdditionalLoadsOnUpdate() storage.TKs {
	return []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: string(m.ID)}}
}

func (m *MutableFederationGateway) GetAdditionalWritesOnUpdate(
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	var ret []configurator.EntityWriteOperation
	existingEnt, ok := loadedEntities[storage.TypeAndKey{Type: feg.FegGatewayType, Key: string(m.ID)}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	entUpdate := configurator.EntityUpdateCriteria{
		Type:      feg.FegGatewayType,
		Key:       string(m.ID),
		NewConfig: m.Federation,
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

func (m *FederatedNetworkConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8rModels.GetNetworkConfig(network, feg.FederatedNetworkType)
}

func (m *FederatedNetworkConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, feg.FederatedNetworkType, m), nil
}

func (m *GatewayFederationConfigs) FromBackendModels(networkID string, gatewayID string) error {
	federationConfig, err := configurator.LoadEntityConfig(networkID, feg.FegGatewayType, gatewayID, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *federationConfig.(*GatewayFederationConfigs)
	return nil
}

func (m *GatewayFederationConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: feg.FegGatewayType, Key: gatewayID,
			NewConfig: m,
		},
	}, nil
}

func (m *DiameterClientConfigs) ToMconfig() *mconfig.DiamClientConfig {
	res := &mconfig.DiamClientConfig{}
	protos.FillIn(m, res)
	return res
}

// TODO: remove this once backwards compatibility is not needed for the field server
func ToMultipleServersMconfig(server *DiameterClientConfigs, servers []*DiameterClientConfigs) []*mconfig.DiamClientConfig {
	diamClientMconfigs := make([]*mconfig.DiamClientConfig, 0, len(servers)+1)
	if server != nil {
		// prepend server to Servers
		tmpSrv := append([]*DiameterClientConfigs{server}, servers...)
		servers = tmpSrv
	}
	for _, diamClientProto := range servers {
		diamClientConf := &mconfig.DiamClientConfig{}
		protos.FillIn(diamClientProto, diamClientConf)
		diamClientMconfigs = append(diamClientMconfigs, diamClientConf)
	}
	return diamClientMconfigs
}

func (m *DiameterServerConfigs) ToMconfig() *mconfig.DiamServerConfig {
	res := &mconfig.DiamServerConfig{}
	protos.FillIn(m, res)
	return res
}

func (m *SubscriptionProfile) ToMconfig() *mconfig.HSSConfig_SubscriptionProfile {
	res := &mconfig.HSSConfig_SubscriptionProfile{}
	protos.FillIn(m, res)
	return res
}

func ToVirtualApnRuleMconfig(rules []*VirtualApnRule) []*mconfig.VirtualApnRule {
	virtualApnRuleConfigs := make([]*mconfig.VirtualApnRule, 0, len(rules)+1)
	for _, ruleProto := range rules {
		apnConf := &mconfig.VirtualApnRule{}
		protos.FillIn(ruleProto, apnConf)
		virtualApnRuleConfigs = append(virtualApnRuleConfigs, apnConf)
	}
	return virtualApnRuleConfigs
}

func ToFederatedModesMap(modesMap *FederatedModeMap) *lte_mconfig.FederatedModeMap {
	if modesMap == nil {
		return &lte_mconfig.FederatedModeMap{}
	}
	res := &lte_mconfig.FederatedModeMap{}
	protos.FillIn(modesMap, res)
	res.Mapping = ToModesMap(modesMap.Mapping)
	return res
}

func ToModesMap(model_modes []*ModeMapItem) []*lte_mconfig.ModeMapItem {
	if model_modes == nil {
		return []*lte_mconfig.ModeMapItem{}
	}
	proto_modes := make([]*lte_mconfig.ModeMapItem, len(model_modes))
	for i, model_mode := range model_modes {
		proto_mode := &lte_mconfig.ModeMapItem{}
		protos.FillIn(model_mode, proto_mode)
		proto_modes[i] = proto_mode
		// translate the mode
		proto_modes[i].Mode = ToFederatedMode(model_mode.Mode)
	}
	return proto_modes
}

func ToFederatedMode(mode string) lte_mconfig.ModeMapItem_FederatedMode {
	switch mode {
	case "local_subscriber":
		return lte_mconfig.ModeMapItem_LOCAL_SUBSCRIBER
	case "s8_subscriber":
		return lte_mconfig.ModeMapItem_S8_SUBSCRIBER
	}
	// default case
	return lte_mconfig.ModeMapItem_SPGW_SUBSCRIBER
}
