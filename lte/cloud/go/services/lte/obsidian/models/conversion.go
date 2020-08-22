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
	"sort"

	"magma/lte/cloud/go/lte"
	policydbModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/go-openapi/swag"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func (m *LteNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &LteNetwork{}
}

func (m *LteNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        lte.NetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: m.Cellular,
			orc8r.DnsdNetworkType:         m.DNS,
			orc8r.NetworkFeaturesConfig:   m.Features,
		},
	}
}

func (m *LteNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			lte.CellularNetworkConfigType: m.Cellular,
			orc8r.DnsdNetworkType:         m.DNS,
			orc8r.NetworkFeaturesConfig:   m.Features,
		},
	}
}

func (m *LteNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[lte.CellularNetworkConfigType]; cfg != nil {
		m.Cellular = cfg.(*NetworkCellularConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*orc8rModels.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8rModels.NetworkFeatures)
	}
	if cfg := n.Configs[lte.NetworkSubscriberConfigType]; cfg != nil {
		m.SubscriberConfig = cfg.(*policydbModels.NetworkSubscriberConfig)
	}
	return m
}

func (m *NetworkCellularConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkConfigType, m), nil
}

func (m *NetworkCellularConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
}

func (m FegNetworkID) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).FegNetworkID = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkConfigType, iCellularConfig), nil
}

func (m FegNetworkID) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).FegNetworkID
}

func (m *NetworkEpcConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).Epc = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkConfigType, iCellularConfig), nil
}

func (m *NetworkEpcConfigs) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).Epc
}

func (m *NetworkRanConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).Ran = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkConfigType, iCellularConfig), nil
}

func (m *NetworkRanConfigs) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).Ran
}

func (m *LteGateway) FromBackendModels(
	networkID string,
	magmadGateway, cellularGateway configurator.NetworkEntity,
	device *orc8rModels.GatewayDevice,
	status *orc8rModels.GatewayStatus,
) (handlers.GatewayModel, error) {
	m.ConnectedEnodebSerials = EnodebSerials{}
	m.ApnResources = ApnResources{}

	magmadGatewayModel := (&orc8rModels.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	err := copier.Copy(m, magmadGatewayModel)
	if err != nil {
		return nil, err
	}
	if cellularGateway.Config != nil {
		m.Cellular = cellularGateway.Config.(*GatewayCellularConfigs)
	}

	for _, tk := range cellularGateway.Associations {
		switch tk.Type {
		case lte.CellularEnodebEntityType:
			m.ConnectedEnodebSerials = append(m.ConnectedEnodebSerials, tk.Key)
		case lte.APNResourceEntityType:
			r := &ApnResource{}
			err := r.Load(networkID, tk.Key)
			if err != nil {
				return nil, errors.Wrap(err, "error loading apn resource entity")
			}
			m.ApnResources[string(r.ApnName)] = *r
		}
	}
	sort.Strings(m.ConnectedEnodebSerials)

	return m, nil
}

func (m *LteGateway) Load(networkID, gatewayID string) error {
	magmadGateway := &orc8rModels.MagmadGateway{}
	err := magmadGateway.Load(networkID, gatewayID)
	if err != nil {
		return err
	}

	cellularEnt, err := configurator.LoadEntity(
		networkID, lte.CellularGatewayEntityType, gatewayID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
	)
	if err != nil {
		return errors.Wrap(err, "error loading cellular gateway")
	}

	gateway := &LteGateway{
		ID:                     magmadGateway.ID,
		Name:                   magmadGateway.Name,
		Description:            magmadGateway.Description,
		Device:                 magmadGateway.Device,
		Status:                 magmadGateway.Status,
		Tier:                   magmadGateway.Tier,
		Magmad:                 magmadGateway.Magmad,
		ApnResources:           ApnResources{},
		ConnectedEnodebSerials: EnodebSerials{},
	}

	if cellularEnt.Config != nil {
		gateway.Cellular = cellularEnt.Config.(*GatewayCellularConfigs)
	}
	for _, tk := range cellularEnt.Associations {
		switch tk.Type {
		case lte.CellularEnodebEntityType:
			gateway.ConnectedEnodebSerials = append(gateway.ConnectedEnodebSerials, tk.Key)
		case lte.APNResourceEntityType:
			r := &ApnResource{}
			err := r.Load(networkID, tk.Key)
			if err != nil {
				return errors.Wrap(err, "error loading apn resource entity")
			}
			gateway.ApnResources[string(r.ApnName)] = *r
		}
	}

	*m = *gateway
	return nil
}

func (m *MutableLteGateway) Load(networkID, gatewayID string) error {
	gateway := &LteGateway{}
	err := gateway.Load(networkID, gatewayID)
	if err != nil {
		return err
	}
	err = copier.Copy(m, gateway)
	if err != nil {
		return err
	}

	return nil
}

func (m *MutableLteGateway) GetMagmadGateway() *orc8rModels.MagmadGateway {
	return &orc8rModels.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableLteGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	var writes []configurator.EntityWriteOperation

	for _, r := range m.ApnResources {
		writes = append(writes, r.ToEntity())
	}

	cellularGateway := configurator.NetworkEntity{
		Type:        lte.CellularGatewayEntityType,
		Key:         string(m.ID),
		Name:        string(m.Name),
		Description: string(m.Description),
		Config:      m.Cellular,
	}
	for _, s := range m.ConnectedEnodebSerials {
		cellularGateway.Associations = append(cellularGateway.Associations, storage.TypeAndKey{Type: lte.CellularEnodebEntityType, Key: s})
	}
	for _, r := range m.ApnResources {
		cellularGateway.Associations = append(cellularGateway.Associations, storage.TypeAndKey{Type: lte.APNResourceEntityType, Key: r.ID})
	}
	writes = append(writes, cellularGateway)

	linkGateways := configurator.EntityUpdateCriteria{
		Type:              orc8r.MagmadGatewayType,
		Key:               string(m.ID),
		AssociationsToAdd: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: string(m.ID)}},
	}
	writes = append(writes, linkGateways)

	return writes
}

func (m *MutableLteGateway) GetAdditionalLoadsOnUpdate() []storage.TypeAndKey {
	ret := []storage.TypeAndKey{
		{Type: lte.CellularGatewayEntityType, Key: string(m.ID)},
	}
	for _, r := range m.ApnResources {
		ret = append(ret, storage.TypeAndKey{Type: lte.APNResourceEntityType, Key: r.ID})
	}
	return ret
}

func (m *MutableLteGateway) GetAdditionalWritesOnUpdate(
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	var writes []configurator.EntityWriteOperation

	existingGateway, ok := loadedEntities[storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: string(m.ID)}]
	if !ok {
		return writes, merrors.ErrNotFound
	}

	apnResourceWrites, newAPNResourceTKs, err := m.getAPNResourceChanges(existingGateway, loadedEntities)
	if err != nil {
		return nil, err
	}
	writes = append(writes, apnResourceWrites...)

	gatewayUpdate := configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               string(m.ID),
		NewConfig:         m.Cellular,
		AssociationsToAdd: newAPNResourceTKs,
	}
	if string(m.Name) != existingGateway.Name {
		gatewayUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingGateway.Description {
		gatewayUpdate.NewDescription = swag.String(string(m.Description))
	}
	for _, enbSerial := range m.ConnectedEnodebSerials {
		gatewayUpdate.AssociationsToSet = append(gatewayUpdate.AssociationsToSet, storage.TypeAndKey{Type: lte.CellularEnodebEntityType, Key: enbSerial})
	}
	writes = append(writes, gatewayUpdate)

	return writes, nil
}

func (m *MutableLteGateway) GetAdditionalDeletes() []storage.TypeAndKey {
	tks := []storage.TypeAndKey{
		{Type: lte.CellularGatewayEntityType, Key: string(m.ID)},
	}
	for _, r := range m.ApnResources {
		tks = append(tks, r.GetTK())
	}
	return tks
}

// getAPNResourceChanges returns required writes, as well as the TKs of the
// new entities.
func (m *MutableLteGateway) getAPNResourceChanges(
	existingGateway configurator.NetworkEntity,
	loaded map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, []storage.TypeAndKey, error) {
	var writes []configurator.EntityWriteOperation

	oldIDs := storage.GetKeys(storage.Filter(existingGateway.Associations, lte.APNResourceEntityType))

	oldByAPN := ApnResources{}
	err := oldByAPN.Load(existingGateway.NetworkID, oldIDs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error loading existing APN resources")
	}
	oldResources := oldByAPN.GetByID()
	newResources := m.ApnResources.GetByID()

	newIDs := funk.Keys(newResources).([]string)
	newTKs := storage.MakeTKs(lte.APNResourceEntityType, newIDs)

	deletes, creates := funk.DifferenceString(oldIDs, newIDs)
	updates := funk.JoinString(oldIDs, newIDs, funk.InnerJoinString)

	for _, w := range deletes {
		writes = append(writes, oldResources[w].ToDeleteCriteria())
	}
	for _, w := range creates {
		writes = append(writes, newResources[w].ToEntity())
	}
	for _, w := range updates {
		writes = append(writes, newResources[w].ToUpdateCriteria())
	}

	return writes, newTKs, nil
}

func (m *GatewayCellularConfigs) FromBackendModels(networkID string, gatewayID string) error {
	cellularConfig, err := configurator.LoadEntityConfig(networkID, lte.CellularGatewayEntityType, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.(*GatewayCellularConfigs)
	return nil
}

func (m *GatewayCellularConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: lte.CellularGatewayEntityType, Key: gatewayID,
			NewConfig: m,
		},
	}, nil
}

func (m *GatewayEpcConfigs) FromBackendModels(networkID string, gatewayID string) error {
	gatewayConfig := &GatewayCellularConfigs{}
	err := gatewayConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *gatewayConfig.Epc
	return nil
}

func (m *GatewayEpcConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Epc = m
	return cellularConfig.ToUpdateCriteria(networkID, gatewayID)
}

func (m *GatewayRanConfigs) FromBackendModels(networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.Ran
	return nil
}

func (m *GatewayRanConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Ran = m
	return cellularConfig.ToUpdateCriteria(networkID, gatewayID)
}

func (m *GatewayNonEpsConfigs) FromBackendModels(networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.NonEpsService
	return nil
}

func (m *GatewayNonEpsConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.NonEpsService = m
	return cellularConfig.ToUpdateCriteria(networkID, gatewayID)
}

func (m *EnodebSerials) FromBackendModels(networkID string, gatewayID string) error {
	cellularGatewayEntity, err := configurator.LoadEntity(networkID, lte.CellularGatewayEntityType, gatewayID, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
	if err != nil {
		return err
	}
	enodebSerials := EnodebSerials{}
	for _, assoc := range cellularGatewayEntity.Associations {
		if assoc.Type == lte.CellularEnodebEntityType {
			enodebSerials = append(enodebSerials, assoc.Key)
		}
	}
	*m = enodebSerials
	return nil
}

func (m *EnodebSerials) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	enodebSerials := []storage.TypeAndKey{}
	for _, enodebSerial := range *m {
		enodebSerials = append(enodebSerials, storage.TypeAndKey{Type: lte.CellularEnodebEntityType, Key: enodebSerial})
	}
	return []configurator.EntityUpdateCriteria{
		{
			Type:              lte.CellularGatewayEntityType,
			Key:               gatewayID,
			AssociationsToSet: enodebSerials,
		},
	}, nil
}

func (m *EnodebSerials) ToDeleteUpdateCriteria(networkID, gatewayID, enodebID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: lte.CellularGatewayEntityType, Key: gatewayID,
		AssociationsToDelete: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: enodebID}},
	}
}

func (m *EnodebSerials) ToCreateUpdateCriteria(networkID, gatewayID, enodebID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: lte.CellularGatewayEntityType, Key: gatewayID,
		AssociationsToAdd: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: enodebID}},
	}
}

func (m *Enodeb) FromBackendModels(ent configurator.NetworkEntity) *Enodeb {
	m.Name = ent.Name
	m.Description = ent.Description
	m.Serial = ent.Key
	if ent.Config != nil {
		m.Config = ent.Config.(*EnodebConfiguration)
	}
	for _, tk := range ent.ParentAssociations {
		if tk.Type == lte.CellularGatewayEntityType {
			m.AttachedGatewayID = tk.Key
		}
	}
	return m
}

func (m *Enodeb) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:           lte.CellularEnodebEntityType,
		Key:            m.Serial,
		NewName:        swag.String(m.Name),
		NewDescription: swag.String(m.Description),
		NewConfig:      m.Config,
	}
}

func (m *Apn) FromBackendModels(ent configurator.NetworkEntity) *Apn {
	m.ApnName = ApnName(ent.Key)
	m.ApnConfiguration = ent.Config.(*ApnConfiguration)
	return m
}

func (m ApnList) ToAssocs() []storage.TypeAndKey {
	return funk.Map(
		m,
		func(rn string) storage.TypeAndKey {
			return storage.TypeAndKey{Type: lte.APNEntityType, Key: rn}
		},
	).([]storage.TypeAndKey)
}

func (m *ApnResources) Load(networkID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	ents, notFound, err := configurator.LoadEntities(
		networkID,
		nil, nil, nil,
		storage.MakeTKs(lte.APNResourceEntityType, ids),
		configurator.EntityLoadCriteria{LoadConfig: true},
	)
	if err != nil {
		return err
	}
	if len(notFound) != 0 {
		return fmt.Errorf("error loading apn resources: could not find following entities: %v", notFound)
	}

	model := ApnResources{}
	for _, ent := range ents {
		r := (&ApnResource{}).FromEntity(ent)
		model[string(r.ApnName)] = *r
	}

	*m = model
	return nil
}

func (m *ApnResources) GetByID() map[string]*ApnResource {
	byID := map[string]*ApnResource{}
	for _, r := range *m {
		byID[r.ID] = &r
	}
	return byID
}

func (m *ApnResources) GetTKs() []storage.TypeAndKey {
	var tks []storage.TypeAndKey
	for _, r := range *m {
		tks = append(tks, r.GetTK())
	}
	return tks
}

func (m *ApnResource) GetTK() storage.TypeAndKey {
	return storage.TypeAndKey{Type: lte.APNResourceEntityType, Key: m.ID}
}

func (m *ApnResource) Load(networkID string, id string) error {
	ent, err := configurator.LoadEntity(networkID, lte.APNResourceEntityType, id, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return err
	}
	*m = *m.FromEntity(ent)
	return nil
}

func (m *ApnResource) ToEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type:         lte.APNResourceEntityType,
		Key:          m.ID,
		Config:       m,
		Associations: m.getAssocs(),
	}
}

func (m *ApnResource) FromEntity(ent configurator.NetworkEntity) *ApnResource {
	return ent.Config.(*ApnResource)
}

func (m *ApnResource) ToUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:              lte.APNResourceEntityType,
		Key:               m.ID,
		NewConfig:         m,
		AssociationsToSet: m.getAssocs(),
	}
}

func (m *ApnResource) ToDeleteCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:         lte.APNResourceEntityType,
		Key:          m.ID,
		DeleteEntity: true,
	}
}

func (m *ApnResource) getAssocs() []storage.TypeAndKey {
	apnAssoc := []storage.TypeAndKey{{Type: lte.APNEntityType, Key: string(m.ApnName)}}
	return apnAssoc
}
