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
	"context"
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"

	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
)

func (m *Network) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        string(m.Type),
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
			orc8r.NetworkSentryConfig:   m.SentryConfig,
			orc8r.StateConfig:           m.StateConfig,
		},
	}
}

func (m *Network) FromConfiguratorNetwork(n configurator.Network) *Network {
	m.ID = models.NetworkID(n.ID)
	m.Type = models.NetworkType(n.Type)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		dns := cfg.(*NetworkDNSConfig)
		if dns != nil {
			m.DNS = dns
		}
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		features := cfg.(*NetworkFeatures)
		if features != nil {
			m.Features = features
		}
	}
	if cfg := n.Configs[orc8r.NetworkSentryConfig]; cfg != nil {
		sentryConfig := cfg.(*NetworkSentryConfig)
		if sentryConfig != nil {
			m.SentryConfig = sentryConfig
		}
	}
	if cfg := n.Configs[orc8r.StateConfig]; cfg != nil {
		stateConfig := cfg.(*StateConfig)
		if stateConfig != nil {
			m.StateConfig = stateConfig
		}
	}
	return m
}

func (m *Network) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewType:        swag.String(string(m.Type)),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
			orc8r.NetworkSentryConfig:   m.SentryConfig,
			orc8r.StateConfig:           m.StateConfig,
		},
	}
}

func (m *NetworkFeatures) GetFromNetwork(network configurator.Network) interface{} {
	return GetNetworkConfig(network, orc8r.NetworkFeaturesConfig)
}

func (m *NetworkFeatures) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.NetworkFeaturesConfig, m), nil
}

func (m *NetworkSentryConfig) GetFromNetwork(network configurator.Network) interface{} {
	return GetNetworkConfig(network, orc8r.NetworkSentryConfig)
}

func (m *NetworkSentryConfig) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.NetworkSentryConfig, m), nil
}

func (m *StateConfig) GetFromNetwork(network configurator.Network) interface{} {
	return GetNetworkConfig(network, orc8r.StateConfig)
}

func (m *StateConfig) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.StateConfig, m), nil
}

func (m *NetworkDNSConfig) GetFromNetwork(network configurator.Network) interface{} {
	return GetNetworkConfig(network, orc8r.DnsdNetworkType)
}

func (m *NetworkDNSConfig) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.DnsdNetworkType, m), nil
}

func (m NetworkDNSRecords) GetFromNetwork(network configurator.Network) interface{} {
	iNetworkDnsConfig := GetNetworkConfig(network, orc8r.DnsdNetworkType)
	if iNetworkDnsConfig == nil {
		return nil
	}
	return iNetworkDnsConfig.(*NetworkDNSConfig).Records
}

func (m NetworkDNSRecords) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iNetworkDnsConfig := GetNetworkConfig(network, orc8r.DnsdNetworkType)
	if iNetworkDnsConfig == nil {
		return configurator.NetworkUpdateCriteria{}, errors.New("No DNS Config registered for this network")
	}
	iNetworkDnsConfig.(*NetworkDNSConfig).Records = m
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.DnsdNetworkType, iNetworkDnsConfig), nil
}

func (m *MagmadGateway) GetMagmadGateway() *MagmadGateway {
	return m
}

func (m *MagmadGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	var physicalID string
	if m.Device != nil {
		physicalID = m.Device.HardwareID
	}

	return []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
			Type:        orc8r.MagmadGatewayType,
			Key:         string(m.ID),
			Name:        string(m.Name),
			Description: string(m.Description),
			Config:      m.Magmad,
			PhysicalID:  physicalID,
		},
	}
}

func (m *MagmadGateway) GetGatewayType() string {
	return orc8r.MagmadGatewayType
}

func (m *MagmadGateway) GetAdditionalLoadsOnLoad(gateway configurator.NetworkEntity) storage.TKs {
	return nil
}

func (m *MagmadGateway) GetAdditionalLoadsOnUpdate() storage.TKs {
	return storage.TKs{{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}}
}

func (m *MagmadGateway) GetAdditionalWritesOnUpdate(ctx context.Context, loadedEntities map[storage.TK]configurator.NetworkEntity) ([]configurator.EntityWriteOperation, error) {
	var ret []configurator.EntityWriteOperation
	existingEnt, ok := loadedEntities[storage.TK{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	gatewayUpdate := configurator.EntityUpdateCriteria{
		Type:      orc8r.MagmadGatewayType,
		Key:       string(m.ID),
		NewConfig: m.Magmad,
	}
	if m.Device.HardwareID != existingEnt.PhysicalID {
		gatewayUpdate.NewPhysicalID = swag.String(m.Device.HardwareID)
	}
	if string(m.Name) != existingEnt.Name {
		gatewayUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingEnt.Description {
		gatewayUpdate.NewDescription = swag.String(string(m.Description))
	}

	oldTierTK, _ := existingEnt.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	if oldTierTK.Key != string(m.Tier) {
		ret = append(
			ret,
			configurator.EntityUpdateCriteria{
				Type: orc8r.UpgradeTierEntityType, Key: oldTierTK.Key,
				AssociationsToDelete: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}},
			},
		)

		ret = append(
			ret,
			configurator.EntityUpdateCriteria{
				Type: orc8r.UpgradeTierEntityType, Key: string(m.Tier),
				AssociationsToAdd: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}},
			},
		)
	}

	// do the tier update to delete the old assoc first
	ret = append(ret, gatewayUpdate)
	return ret, nil
}

func (m *MagmadGateway) ToConfiguratorEntities() []configurator.NetworkEntity {
	gatewayEnt := configurator.NetworkEntity{
		Type:        orc8r.MagmadGatewayType,
		Key:         string(m.ID),
		Name:        string(m.Name),
		Description: string(m.Description),
		Config:      m.Magmad,
		PhysicalID:  m.Device.HardwareID,
	}
	return []configurator.NetworkEntity{gatewayEnt}
}

func (m *MagmadGateway) FromBackendModels(ent configurator.NetworkEntity, device *GatewayDevice, status *GatewayStatus) *MagmadGateway {
	m.ID = models.GatewayID(ent.Key)
	m.Name = models.GatewayName(ent.Name)
	m.Description = models.GatewayDescription(ent.Description)
	if ent.Config != nil {
		m.Magmad = ent.Config.(*MagmadGatewayConfigs)
	}
	m.Device = device
	m.Status = status
	tierTK, err := ent.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	if err == nil {
		m.Tier = TierID(tierTK.Key)
	}
	err = PopulateRegistrationInfo(context.Background(), m, ent.NetworkID)
	if err != nil {
		// ignore err and continue returning the rest of the gateway
		glog.V(2).Infof("failed to populate registration info because %v+", err)
	}

	return m
}

func PopulateRegistrationInfos(ctx context.Context, gateways map[string]*MagmadGateway, networkID string) error {
	for _, gateway := range gateways {
		err := PopulateRegistrationInfo(ctx, gateway, networkID)
		if err != nil {
			return err
		}
	}
	return nil
}

// PopulateRegistrationInfo will populate the given gateway's RegistrationInfo if its device is nil
func PopulateRegistrationInfo(ctx context.Context, gateway *MagmadGateway, networkID string) error {
	if gateway.Device != nil {
		return nil
	}

	regToken, err := bootstrapper.GetToken(ctx, networkID, string(gateway.ID), true)
	if err != nil {
		return err
	}

	regInfo, err := bootstrapper.GetGatewayRegistrationInfo(ctx)
	if err != nil {
		return err
	}

	gateway.RegistrationInfo = &models.RegistrationInfo{
		DomainName:        &regInfo.DomainName,
		RegistrationToken: &regToken,
		RootCa:            &regInfo.RootCa,
	}
	return nil
}

func (m *MagmadGateway) ToEntityUpdateCriteria(existingEnt configurator.NetworkEntity) []configurator.EntityUpdateCriteria {
	ret := []configurator.EntityUpdateCriteria{}
	gatewayUpdate := configurator.EntityUpdateCriteria{
		Type:      orc8r.MagmadGatewayType,
		Key:       string(m.ID),
		NewConfig: m.Magmad,
	}

	if m.Device.HardwareID != existingEnt.PhysicalID {
		gatewayUpdate.NewPhysicalID = swag.String(m.Device.HardwareID)
	}
	if string(m.Name) != existingEnt.Name {
		gatewayUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingEnt.Description {
		gatewayUpdate.NewDescription = swag.String(string(m.Description))
	}

	oldTierTK, _ := existingEnt.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	if oldTierTK.Key != string(m.Tier) {
		ret = append(
			ret,
			configurator.EntityUpdateCriteria{
				Type: orc8r.UpgradeTierEntityType, Key: oldTierTK.Key,
				AssociationsToDelete: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}},
			},
		)

		ret = append(
			ret,
			configurator.EntityUpdateCriteria{
				Type: orc8r.UpgradeTierEntityType, Key: string(m.Tier),
				AssociationsToAdd: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}},
			},
		)
	}

	// do the tier update to delete the old assoc first
	ret = append(ret, gatewayUpdate)
	return ret
}

func (m *MagmadGatewayConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	exists, err := configurator.DoesEntityExist(ctx, networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Gateway %s does not exist", gatewayID)
	}

	return []configurator.EntityUpdateCriteria{
		{
			Key:       gatewayID,
			Type:      orc8r.MagmadGatewayType,
			NewConfig: m,
		},
	}, nil
}

func (m *MagmadGatewayConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	config, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.MagmadGatewayType, gatewayID, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *config.(*MagmadGatewayConfigs)
	return nil
}

func (m *TierID) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	entity, err := configurator.LoadEntity(
		ctx,
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true},
		EntitySerdes,
	)
	if err != nil {
		return err
	}
	for _, parentAssoc := range entity.ParentAssociations {
		if parentAssoc.Type == orc8r.UpgradeTierEntityType {
			*m = TierID(parentAssoc.Key)
			return nil
		}
	}
	return nil
}

func (m *TierID) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	tierID := string(*m)
	updateCriteria := []configurator.EntityUpdateCriteria{}

	exists, err := configurator.DoesEntityExist(ctx, networkID, orc8r.UpgradeTierEntityType, tierID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to look up tier")
	}
	if !exists {
		return nil, fmt.Errorf("Tier %s does not exist", tierID)
	}

	// Remove association from old tier
	entity, err := configurator.LoadEntity(
		ctx,
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true},
		EntitySerdes,
	)
	if err != nil {
		return nil, err
	}

	tierTK, err := entity.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	if err != merrors.ErrNotFound {
		if tierTK.Key == tierID {
			// no change
			return []configurator.EntityUpdateCriteria{}, nil
		}
		deleteCurrentTierAssoc := configurator.EntityUpdateCriteria{
			Type:                 tierTK.Type,
			Key:                  tierTK.Key,
			AssociationsToDelete: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
		}
		updateCriteria = append(updateCriteria, deleteCurrentTierAssoc)
	}

	// Add association to new tier
	addNewTierAssoc := configurator.EntityUpdateCriteria{
		Type:              orc8r.UpgradeTierEntityType,
		Key:               tierID,
		AssociationsToAdd: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
	}
	updateCriteria = append(updateCriteria, addNewTierAssoc)
	return updateCriteria, nil
}

func GetNetworkConfig(network configurator.Network, key string) interface{} {
	if network.Configs == nil {
		return nil
	}
	config, exists := network.Configs[key]
	if !exists {
		return nil
	}
	return config
}

func GetNetworkConfigUpdateCriteria(networkID string, key string, iConfig interface{}) configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID: networkID,
		ConfigsToAddOrUpdate: map[string]interface{}{
			key: iConfig,
		},
	}
}

func (m *Tier) ToNetworkEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type: orc8r.UpgradeTierEntityType, Key: string(m.ID),
		Name:         string(m.Name),
		Config:       m,
		Associations: getGatewayTKs(m.Gateways),
	}
}

func (m *Tier) ToUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: orc8r.UpgradeTierEntityType, Key: string(m.ID),
		NewName:           swag.String(string(m.Name)),
		NewConfig:         m,
		AssociationsToSet: getGatewayTKs(m.Gateways),
	}
}

func (m *Tier) FromBackendModel(entity configurator.NetworkEntity) *Tier {
	tier := entity.Config.(*Tier)
	tier.Name = TierName(entity.Name)
	tier.Gateways = getGatewayIDs(entity.Associations)
	return tier
}

func (m *TierName) ToUpdateCriteria(ctx context.Context, networkID string, key string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: orc8r.UpgradeTierEntityType, Key: key, NewName: swag.String(string(*m)),
		},
	}, nil
}

func (m *TierName) FromBackendModels(ctx context.Context, networkID string, key string) error {
	entity, err := configurator.LoadEntity(
		ctx,
		networkID, orc8r.UpgradeTierEntityType, key,
		configurator.EntityLoadCriteria{LoadMetadata: true},
		EntitySerdes,
	)
	if err != nil {
		return err
	}
	*m = TierName(entity.Name)
	return nil
}

func (m *TierVersion) FromBackendModels(ctx context.Context, networkID string, key string) error {
	iConfig, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.UpgradeTierEntityType, key, EntitySerdes)
	if err != nil {
		return err
	}
	*m = iConfig.(*Tier).Version
	return nil
}

func (m *TierVersion) ToUpdateCriteria(ctx context.Context, networkID string, key string) ([]configurator.EntityUpdateCriteria, error) {
	iConfig, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.UpgradeTierEntityType, key, EntitySerdes)
	if err != nil {
		return []configurator.EntityUpdateCriteria{}, err
	}
	tier := iConfig.(*Tier)
	tier.Version = *m
	return []configurator.EntityUpdateCriteria{
		{Type: orc8r.UpgradeTierEntityType, Key: key, NewConfig: tier},
	}, nil
}

func (m *TierVersion) ToString() string {
	return string(*m)
}

func (m *TierImages) FromBackendModels(ctx context.Context, networkID string, key string) error {
	iConfig, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.UpgradeTierEntityType, key, EntitySerdes)
	if err != nil {
		return err
	}
	*m = iConfig.(*Tier).Images
	return nil
}

func (m *TierImages) ToUpdateCriteria(ctx context.Context, networkID string, key string) ([]configurator.EntityUpdateCriteria, error) {
	iConfig, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.UpgradeTierEntityType, key, EntitySerdes)
	if err != nil {
		return []configurator.EntityUpdateCriteria{}, err
	}
	tier := iConfig.(*Tier)
	tier.Images = *m
	return []configurator.EntityUpdateCriteria{
		{
			Type: orc8r.UpgradeTierEntityType, Key: key, NewConfig: tier,
		},
	}, nil
}

func (m *TierGateways) FromBackendModels(ctx context.Context, networkID string, key string) error {
	tierEnt, err := configurator.LoadEntity(
		ctx,
		networkID, orc8r.UpgradeTierEntityType, key,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		EntitySerdes,
	)
	if err != nil {
		return err
	}
	*m = getGatewayIDs(tierEnt.Associations)
	return nil
}

func (m *TierGateways) ToUpdateCriteria(ctx context.Context, networkID string, key string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: orc8r.UpgradeTierEntityType, Key: key,
			AssociationsToSet: getGatewayTKs(*m),
		},
	}, nil
}

func (m *TierGateways) ToAddGatewayUpdateCriteria(tierID, gatewayID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: orc8r.UpgradeTierEntityType, Key: tierID,
		AssociationsToAdd: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
	}
}

func (m *TierGateways) ToDeleteGatewayUpdateCriteria(tierID, gatewayID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: orc8r.UpgradeTierEntityType, Key: tierID,
		AssociationsToDelete: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
	}
}

func (m *TierImage) ToUpdateCriteria(ctx context.Context, networkID string, key string) ([]configurator.EntityUpdateCriteria, error) {
	iConfig, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.UpgradeTierEntityType, key, EntitySerdes)
	if err != nil {
		return []configurator.EntityUpdateCriteria{}, err
	}
	tier := iConfig.(*Tier)
	tier.Images = append(tier.Images, m)
	return []configurator.EntityUpdateCriteria{
		{Type: orc8r.UpgradeTierEntityType, Key: key, NewConfig: tier},
	}, nil
}

func (m *TierImage) ToDeleteImageUpdateCriteria(ctx context.Context, networkID, tierID, imageName string) (configurator.EntityUpdateCriteria, error) {
	iConfig, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.UpgradeTierEntityType, tierID, EntitySerdes)
	if err != nil {
		return configurator.EntityUpdateCriteria{}, err
	}
	tier := iConfig.(*Tier)
	for i, image := range tier.Images {
		if swag.StringValue(image.Name) == imageName {
			if i == len(tier.Images)-1 {
				tier.Images = tier.Images[:i]
			} else {
				tier.Images = append(tier.Images[:i], tier.Images[i+1:]...)
			}
			return configurator.EntityUpdateCriteria{Type: orc8r.UpgradeTierEntityType, Key: tierID, NewConfig: tier}, nil
		}
	}
	return configurator.EntityUpdateCriteria{}, merrors.ErrNotFound
}

func (m *GatewayVpnConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	gatewayConfig := &MagmadGatewayConfigs{}
	err := gatewayConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	gatewayConfig.Vpn = m
	return gatewayConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
}

func (m *GatewayVpnConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	config, err := configurator.LoadEntityConfig(ctx, networkID, orc8r.MagmadGatewayType, gatewayID, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *config.(*MagmadGatewayConfigs).Vpn
	return nil
}

func getGatewayTKs(gateways []models.GatewayID) storage.TKs {
	return funk.Map(
		gateways,
		func(gw models.GatewayID) storage.TK {
			return storage.TK{Type: orc8r.MagmadGatewayType, Key: string(gw)}
		}).([]storage.TK)
}

func getGatewayIDs(gatewayTKs storage.TKs) []models.GatewayID {
	return funk.Map(
		gatewayTKs,
		func(tk storage.TK) models.GatewayID {
			return models.GatewayID(tk.Key)
		}).([]models.GatewayID)
}
