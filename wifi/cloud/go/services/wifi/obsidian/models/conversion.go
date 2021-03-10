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
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8r_models "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/wifi/cloud/go/wifi"

	"github.com/go-openapi/swag"
)

func (m *WifiNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &WifiNetwork{}
}

func (m *WifiNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        wifi.WifiNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: m.Features,
			wifi.WifiNetworkType:        m.Wifi,
		},
	}
}

func (m *WifiNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: m.Features,
			wifi.WifiNetworkType:        m.Wifi,
		},
	}
}

func (m *WifiNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[wifi.WifiNetworkType]; cfg != nil {
		m.Wifi = cfg.(*NetworkWifiConfigs)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8r_models.NetworkFeatures)
	}
	return m
}

func (m *NetworkWifiConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8r_models.GetNetworkConfigUpdateCriteria(network.ID, wifi.WifiNetworkType, m), nil
}

func (m *NetworkWifiConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8r_models.GetNetworkConfig(network, wifi.WifiNetworkType)
}

func (m *WifiGateway) FromBackendModels(
	magmadGateway, wifiGateway configurator.NetworkEntity,
	device *orc8r_models.GatewayDevice,
	status *orc8r_models.GatewayStatus,
) handlers.GatewayModel {
	// delegate most of the fillin to magmad gateway struct
	mdGW := (&orc8r_models.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status
	if wifiGateway.Config != nil {
		m.Wifi = wifiGateway.Config.(*GatewayWifiConfigs)
	}
	return m
}

func (m *MutableWifiGateway) GetMagmadGateway() *orc8r_models.MagmadGateway {
	return &orc8r_models.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableWifiGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	updates := []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
			Type:        wifi.WifiGatewayType,
			Key:         string(m.ID),
			Name:        string(m.Name),
			Description: string(m.Description),
			Config:      m.Wifi,
		},
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: string(m.ID)}},
		},
	}
	for _, meshUpdate := range GetMeshUpdates(string(m.ID), "", m.Wifi.MeshID) {
		updates = append(updates, meshUpdate)
	}
	return updates
}

func (m *MutableWifiGateway) GetGatewayType() string {
	return wifi.WifiGatewayType
}

func (m *MutableWifiGateway) GetAdditionalLoadsOnLoad(gateway configurator.NetworkEntity) storage.TKs {
	return nil
}

func (m *MutableWifiGateway) GetAdditionalLoadsOnUpdate() storage.TKs {
	return []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: string(m.ID)}}
}

func (m *MutableWifiGateway) GetAdditionalWritesOnUpdate(
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	var ret []configurator.EntityWriteOperation
	existingEnt, ok := loadedEntities[storage.TypeAndKey{Type: wifi.WifiGatewayType, Key: string(m.ID)}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	entUpdate := configurator.EntityUpdateCriteria{
		Type:      wifi.WifiGatewayType,
		Key:       string(m.ID),
		NewConfig: m.Wifi,
	}
	if string(m.Name) != existingEnt.Name {
		entUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingEnt.Description {
		entUpdate.NewDescription = swag.String(string(m.Description))
	}
	ret = append(ret, entUpdate)

	// If the mesh id in the gateway's wifi configs are changing, we have to
	// handle updating corresponding meshes appropriately
	oldMeshID := existingEnt.Config.(*GatewayWifiConfigs).MeshID
	newMeshID := m.Wifi.MeshID
	for _, update := range GetMeshUpdates(string(m.ID), oldMeshID, newMeshID) {
		ret = append(ret, update)
	}

	return ret, nil
}

func (m *GatewayWifiConfigs) FromBackendModels(networkID string, gatewayID string) error {
	wifiConfig, err := configurator.LoadEntityConfig(networkID, wifi.WifiGatewayType, gatewayID, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *wifiConfig.(*GatewayWifiConfigs)
	return nil
}

func (m *GatewayWifiConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	ret := []configurator.EntityUpdateCriteria{}

	ret = append(ret, configurator.EntityUpdateCriteria{
		Type: wifi.WifiGatewayType, Key: gatewayID,
		NewConfig: m,
	})

	existingWifiConfigEnt, err := configurator.LoadEntityConfig(networkID, wifi.WifiGatewayType, gatewayID, EntitySerdes)
	if err != nil {
		return nil, err
	}
	existingWifiConfigs := existingWifiConfigEnt.(*GatewayWifiConfigs)

	// If the mesh id in the gateway's wifi configs are changing, we have to
	// handle updating corresponding meshes appropriately
	oldMeshID := existingWifiConfigs.MeshID
	newMeshID := m.MeshID
	ret = append(ret, GetMeshUpdates(gatewayID, oldMeshID, newMeshID)...)

	return ret, nil
}

func (m *WifiMesh) FromBackendModels(ent configurator.NetworkEntity) *WifiMesh {
	m.Name = MeshName(ent.Name)
	m.ID = MeshID(ent.Key)

	if ent.Config != nil {
		m.Config = ent.Config.(*MeshWifiConfigs)
	}

	gwIds := []models.GatewayID{}
	for _, gwAssoc := range ent.Associations {
		if gwAssoc.Type == orc8r.MagmadGatewayType {
			gwIds = append(gwIds, models.GatewayID(gwAssoc.Key))
		}
	}

	m.GatewayIds = gwIds

	return m
}

func (m *WifiMesh) ToUpdateCriteria() []configurator.EntityUpdateCriteria {
	// TODO: update gateway mesh id if it is added or deleted here. For now,
	// don't allow gatewayids to be updated (this logic is in handlers.go)
	gwIds := []storage.TypeAndKey{}
	for _, gwId := range m.GatewayIds {
		gwIds = append(gwIds, storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: string(gwId)})
	}
	return []configurator.EntityUpdateCriteria{
		{
			Type:              wifi.MeshEntityType,
			Key:               string(m.ID),
			NewName:           swag.String(string(m.Name)),
			AssociationsToSet: gwIds,
			NewConfig:         m.Config,
		},
	}
}

func (m *MeshWifiConfigs) FromBackendModels(networkID string, meshID string) error {
	meshConfig, err := configurator.LoadEntityConfig(networkID, wifi.MeshEntityType, meshID, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *meshConfig.(*MeshWifiConfigs)
	return nil
}

func (m *MeshWifiConfigs) ToUpdateCriteria(networkID string, meshID string) ([]configurator.EntityUpdateCriteria, error) {
	ret := []configurator.EntityUpdateCriteria{}

	ret = append(ret, configurator.EntityUpdateCriteria{
		Type: wifi.MeshEntityType, Key: meshID,
		NewConfig: m,
	})

	return ret, nil
}

func (m *MeshName) FromBackendModels(networkID string, meshID string) error {
	meshEnt, err := configurator.LoadEntity(
		networkID, wifi.MeshEntityType, meshID,
		configurator.EntityLoadCriteria{LoadMetadata: true},
		EntitySerdes,
	)
	if err != nil {
		return err
	}
	*m = MeshName(meshEnt.Name)
	return nil
}

func (m *MeshName) ToUpdateCriteria(networkID string, meshID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: wifi.MeshEntityType, Key: meshID,
			NewName: swag.String(string(*m)),
		},
	}, nil
}
