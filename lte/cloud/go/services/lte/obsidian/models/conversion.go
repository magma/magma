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
	"sort"

	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	policydbModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/models"
	commonModels "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
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

func (m *NetworkNgcConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).Ngc = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkConfigType, iCellularConfig), nil
}

func (m *NetworkNgcConfigs) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkConfigType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).Ngc
}

func (m *NetworkNgcConfigs) ConvertSuciEntsToProtos(ent *SuciProfile) *protos.SuciProfile {
	suciData := &protos.SuciProfile{
		HomeNetPublicKeyId: ent.HomeNetworkPublicKeyIdentifier,
		HomeNetPublicKey:   ent.HomeNetworkPublicKey,
		HomeNetPrivateKey:  ent.HomeNetworkPrivateKey,
		ProtectionScheme:   protos.SuciProfile_ECIESProtectionScheme(protos.SuciProfile_ECIESProtectionScheme_value[ent.ProtectionScheme]),
	}
	return suciData
}

func (m *LteGateway) FromBackendModels(
	magmadGateway, cellularGateway configurator.NetworkEntity,
	loadedEntsByTK configurator.NetworkEntitiesByTK,
	device *orc8rModels.GatewayDevice,
	status *orc8rModels.GatewayStatus,
) handlers.GatewayModel {
	m.ConnectedEnodebSerials = EnodebSerials{}
	m.ApnResources = ApnResources{}

	// delegate most of the fillin to magmad gateway struct
	mdGW := (&orc8rModels.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status
	if cellularGateway.Config != nil {
		m.Cellular = cellularGateway.Config.(*GatewayCellularConfigs)
	}

	for _, tk := range cellularGateway.Associations.Filter(lte.APNResourceEntityType) {
		r := (&ApnResource{}).FromEntity(loadedEntsByTK[tk])
		m.ApnResources[string(r.ApnName)] = *r
	}

	for _, tk := range cellularGateway.Associations {
		if tk.Type == lte.CellularEnodebEntityType {
			m.ConnectedEnodebSerials = append(m.ConnectedEnodebSerials, tk.Key)
		}
	}
	sort.Strings(m.ConnectedEnodebSerials)

	return m
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
		cellularGateway.Associations = append(cellularGateway.Associations, storage.TK{Type: lte.CellularEnodebEntityType, Key: s})
	}
	for _, r := range m.ApnResources {
		cellularGateway.Associations = append(cellularGateway.Associations, storage.TK{Type: lte.APNResourceEntityType, Key: r.ID})
	}
	writes = append(writes, cellularGateway)

	linkGateways := configurator.EntityUpdateCriteria{
		Type:              orc8r.MagmadGatewayType,
		Key:               string(m.ID),
		AssociationsToAdd: storage.TKs{{Type: lte.CellularGatewayEntityType, Key: string(m.ID)}},
	}
	writes = append(writes, linkGateways)

	return writes
}

func (m *MutableLteGateway) GetGatewayType() string {
	return lte.CellularGatewayEntityType
}

func (m *MutableLteGateway) GetAdditionalLoadsOnLoad(gateway configurator.NetworkEntity) storage.TKs {
	return gateway.Associations.Filter(lte.APNResourceEntityType)
}

func (m *MutableLteGateway) GetAdditionalLoadsOnUpdate() storage.TKs {
	var loads storage.TKs
	loads = append(loads, storage.TK{Type: lte.CellularGatewayEntityType, Key: string(m.ID)})
	loads = append(loads, m.ApnResources.ToTKs()...)
	return loads
}

func (m *MutableLteGateway) GetAdditionalWritesOnUpdate(ctx context.Context, loadedEntities map[storage.TK]configurator.NetworkEntity) ([]configurator.EntityWriteOperation, error) {
	var writes []configurator.EntityWriteOperation

	existingGateway, ok := loadedEntities[storage.TK{Type: lte.CellularGatewayEntityType, Key: string(m.ID)}]
	if !ok {
		return writes, merrors.ErrNotFound
	}

	apnResourceWrites, newAPNResourceTKs, err := m.getAPNResourceChanges(ctx, existingGateway, loadedEntities)
	if err != nil {
		return nil, err
	}
	writes = append(writes, apnResourceWrites...)

	gatewayUpdate := configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               string(m.ID),
		NewConfig:         m.Cellular,
		AssociationsToSet: newAPNResourceTKs,
	}
	if string(m.Name) != existingGateway.Name {
		gatewayUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingGateway.Description {
		gatewayUpdate.NewDescription = swag.String(string(m.Description))
	}
	for _, enbSerial := range m.ConnectedEnodebSerials {
		gatewayUpdate.AssociationsToSet = append(gatewayUpdate.AssociationsToSet, storage.TK{Type: lte.CellularEnodebEntityType, Key: enbSerial})
	}
	writes = append(writes, gatewayUpdate)

	return writes, nil
}

// getAPNResourceChanges returns required writes, as well as the TKs of the
// new entities.
func (m *MutableLteGateway) getAPNResourceChanges(
	ctx context.Context,
	existingGateway configurator.NetworkEntity,
	loaded map[storage.TK]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, storage.TKs, error) {
	var writes []configurator.EntityWriteOperation

	oldIDs := existingGateway.Associations.Filter(lte.APNResourceEntityType).Keys()
	oldByAPN, err := LoadAPNResources(ctx, existingGateway.NetworkID, oldIDs)
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

func (m *GatewayCellularConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	cellularConfig, err := configurator.LoadEntityConfig(ctx, networkID, lte.CellularGatewayEntityType, gatewayID, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *cellularConfig.(*GatewayCellularConfigs)
	return nil
}

func (m *GatewayCellularConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: lte.CellularGatewayEntityType, Key: gatewayID,
			NewConfig: m,
		},
	}, nil
}

func (m *GatewayEpcConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	gatewayConfig := &GatewayCellularConfigs{}
	err := gatewayConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *gatewayConfig.Epc
	return nil
}

func (m *GatewayEpcConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Epc = m
	return cellularConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
}

func (m *GatewayRanConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.Ran
	return nil
}

func (m *GatewayRanConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Ran = m
	return cellularConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
}

func (m *GatewayNgcConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	gatewayConfig := &GatewayCellularConfigs{}
	err := gatewayConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *gatewayConfig.Ngc
	return nil
}

func (m *GatewayNgcConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Ngc = m
	return cellularConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
}

func (m *GatewayNonEpsConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.NonEpsService
	return nil
}

func (m *GatewayNonEpsConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.NonEpsService = m
	return cellularConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
}

func (m *GatewayDNSConfigs) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return err
	}
	if cellularConfig.DNS != nil {
		*m = *cellularConfig.DNS
	}
	return nil
}

func (m *GatewayDNSConfigs) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.DNS = m
	return cellularConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
}

func (m *GatewayDNSRecords) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return err
	}
	if cellularConfig.DNS != nil {
		*m = cellularConfig.DNS.Records
	}
	return nil
}

func (m *GatewayDNSRecords) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.DNS.Records = *m
	return cellularConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
}

func (m *EnodebSerials) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	cellularGatewayEntity, err := configurator.LoadEntity(
		ctx,
		networkID, lte.CellularGatewayEntityType, gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		EntitySerdes,
	)
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

func (m *EnodebSerials) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	enodebSerials := storage.TKs{}
	for _, enodebSerial := range *m {
		enodebSerials = append(enodebSerials, storage.TK{Type: lte.CellularEnodebEntityType, Key: enodebSerial})
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
		AssociationsToDelete: storage.TKs{{Type: lte.CellularEnodebEntityType, Key: enodebID}},
	}
}

func (m *EnodebSerials) ToCreateUpdateCriteria(networkID, gatewayID, enodebID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: lte.CellularGatewayEntityType, Key: gatewayID,
		AssociationsToAdd: storage.TKs{{Type: lte.CellularEnodebEntityType, Key: enodebID}},
	}
}

func (m *Enodeb) FromBackendModels(ent configurator.NetworkEntity) *Enodeb {
	m.Name = ent.Name
	m.Description = ent.Description
	m.Serial = ent.Key
	if ent.Config != nil {
		// TODO(v1.4.0+): For backwards compatibility we maintain the 'config'
		// field previously reserved for managed enb configs.
		//  We can remove this after the next minor version
		config := ent.Config.(*EnodebConfig)
		m.Config = config.ManagedConfig
		m.EnodebConfig = config
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
		NewConfig:      m.EnodebConfig,
	}
}

func (m *Apn) FromBackendModels(ent configurator.NetworkEntity) *Apn {
	m.ApnName = ApnName(ent.Key)
	m.ApnConfiguration = ent.Config.(*ApnConfiguration)
	return m
}

func LoadAPNResources(ctx context.Context, networkID string, ids []string) (ApnResources, error) {
	ret := ApnResources{}
	if len(ids) == 0 {
		return ret, nil
	}

	ents, notFound, err := configurator.LoadEntities(
		ctx,
		networkID,
		nil, nil, nil,
		storage.MakeTKs(lte.APNResourceEntityType, ids),
		configurator.EntityLoadCriteria{LoadConfig: true},
		EntitySerdes,
	)
	if err != nil {
		return ret, err
	}
	if len(notFound) != 0 {
		return ret, fmt.Errorf("error loading apn resources: could not find following entities: %v", notFound)
	}

	model := ApnResources{}
	for _, ent := range ents {
		r := (&ApnResource{}).FromEntity(ent)
		model[string(r.ApnName)] = *r
	}

	return model, nil
}

func (m *ApnResources) GetByID() map[string]*ApnResource {
	byID := map[string]*ApnResource{}
	for i, r := range *m {
		var apnr = (*m)[i]
		byID[r.ID] = &apnr
	}
	return byID
}

func (m *ApnResources) ToTKs() storage.TKs {
	var tks storage.TKs
	for _, r := range *m {
		tks = append(tks, r.ToTK())
	}
	return tks
}

func (m *ApnResources) ToProto() map[string]*protos.APNConfiguration_APNResource {
	byAPN := map[string]*protos.APNConfiguration_APNResource{}
	if m == nil {
		return nil
	}
	for _, r := range *m {
		byAPN[string(r.ApnName)] = r.ToProto()
	}
	return byAPN
}

func (m *ApnResource) ToTK() storage.TK {
	return storage.TK{Type: lte.APNResourceEntityType, Key: m.ID}
}

func (m *ApnResource) ToEntity() configurator.NetworkEntity {
	cfg := *m // make explicit copy
	return configurator.NetworkEntity{
		Type:         lte.APNResourceEntityType,
		Key:          m.ID,
		Config:       &cfg,
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

func (m *ApnResource) ToProto() *protos.APNConfiguration_APNResource {
	if m == nil {
		return nil
	}
	proto := &protos.APNConfiguration_APNResource{
		ApnName:    string(m.ApnName),
		GatewayIp:  m.GatewayIP.String(),
		GatewayMac: m.GatewayMac.String(),
		VlanId:     m.VlanID,
	}
	return proto
}

func (m *ApnResource) getAssocs() storage.TKs {
	apnAssoc := storage.TKs{{Type: lte.APNEntityType, Key: string(m.ApnName)}}
	return apnAssoc
}

func (m *CellularGatewayPool) FromBackendModels(ent configurator.NetworkEntity) error {
	m.GatewayPoolName = ent.Name
	m.GatewayPoolID = GatewayPoolID(ent.Key)
	cfg, ok := ent.Config.(*CellularGatewayPoolConfigs)
	if !ok {
		return fmt.Errorf("could not convert entity config type %T to GateawyPool", ent.Config)
	}
	m.Config = cfg
	m.GatewayIds = []models.GatewayID{}
	for _, gwID := range ent.Associations {
		m.GatewayIds = append(m.GatewayIds, commonModels.GatewayID(gwID.Key))
	}
	return nil
}

func (m *CellularGatewayPool) ToEntity() configurator.NetworkEntity {
	assocs := storage.TKs{}
	for _, id := range m.GatewayIds {
		tk := storage.TK{Type: lte.CellularGatewayEntityType, Key: string(id)}
		assocs = append(assocs, tk)
	}
	ent := configurator.NetworkEntity{
		Key:          string(m.GatewayPoolID),
		Type:         lte.CellularGatewayPoolEntityType,
		Config:       m.Config,
		Name:         m.GatewayPoolName,
		Associations: assocs,
	}
	return ent
}

func (m *CellularGatewayPool) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	update := configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayPoolEntityType,
		Key:               string(m.GatewayPoolID),
		NewName:           &m.GatewayPoolName,
		NewConfig:         m.Config,
		AssociationsToSet: m.getAssocs(),
	}
	return update
}

func (m *CellularGatewayPool) getAssocs() storage.TKs {
	assocs := storage.TKs{}
	for _, gwID := range m.GatewayIds {
		gateway := storage.TK{
			Type: lte.CellularGatewayEntityType,
			Key:  string(gwID),
		}
		assocs = append(assocs, gateway)
	}
	return assocs
}

func (m *MutableCellularGatewayPool) ToEntity() configurator.NetworkEntity {
	ent := configurator.NetworkEntity{
		Key:    string(m.GatewayPoolID),
		Type:   lte.CellularGatewayPoolEntityType,
		Config: m.Config,
		Name:   m.GatewayPoolName,
	}
	return ent
}

func (m *MutableCellularGatewayPool) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	update := configurator.EntityUpdateCriteria{
		Type:      lte.CellularGatewayPoolEntityType,
		Key:       string(m.GatewayPoolID),
		NewName:   &m.GatewayPoolName,
		NewConfig: m.Config,
	}
	return update
}

func (m *CellularGatewayPoolRecords) FromBackendModels(ctx context.Context, networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = cellularConfig.Pooling
	return nil
}

func (m *CellularGatewayPoolRecords) ToUpdateCriteria(ctx context.Context, networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	updates := []configurator.EntityUpdateCriteria{}
	gatewayEnt, err := configurator.LoadEntity(
		ctx,
		networkID, lte.CellularGatewayEntityType, gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true},
		EntitySerdes,
	)
	if err != nil {
		return nil, err
	}
	var oldPoolIds []GatewayPoolID
	for _, parentAssoc := range gatewayEnt.ParentAssociations.Filter(lte.CellularGatewayPoolEntityType) {
		oldPoolIds = append(oldPoolIds, GatewayPoolID(parentAssoc.Key))
	}
	var newPoolIds []GatewayPoolID
	for _, record := range *m {
		newPoolIds = append(newPoolIds, record.GatewayPoolID)
	}
	err = validateNewGatewayPools(ctx, networkID, newPoolIds)
	if err != nil {
		return nil, err
	}

	idsToDelete, idsToAdd := funk.Difference(oldPoolIds, newPoolIds)
	for _, idToDelete := range idsToDelete.([]GatewayPoolID) {
		deleteCurrentPoolAssoc := configurator.EntityUpdateCriteria{
			Type:                 lte.CellularGatewayPoolEntityType,
			Key:                  string(idToDelete),
			AssociationsToDelete: storage.TKs{{Type: lte.CellularGatewayEntityType, Key: gatewayID}},
		}
		updates = append(updates, deleteCurrentPoolAssoc)
	}
	for _, idToAdd := range idsToAdd.([]GatewayPoolID) {
		addNewPoolAssoc := configurator.EntityUpdateCriteria{
			Type:              lte.CellularGatewayPoolEntityType,
			Key:               string(idToAdd),
			AssociationsToAdd: storage.TKs{{Type: lte.CellularGatewayEntityType, Key: gatewayID}},
		}
		updates = append(updates, addNewPoolAssoc)
	}
	cellularConfig := &GatewayCellularConfigs{}
	err = cellularConfig.FromBackendModels(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Pooling = *m
	configUpdates, err := cellularConfig.ToUpdateCriteria(context.Background(), networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	updates = append(updates, configUpdates...)
	return updates, nil
}

func validateNewGatewayPools(ctx context.Context, networkID string, ids []GatewayPoolID) error {
	var mmeGroupID uint32
	for i, id := range ids {
		ent, err := configurator.LoadEntity(ctx, networkID, lte.CellularGatewayPoolEntityType, string(id),
			configurator.EntityLoadCriteria{LoadConfig: true}, EntitySerdes)
		if err == merrors.ErrNotFound {
			return fmt.Errorf("Gateway pool %s does not exist", id)
		}
		if err != nil {
			return err
		}
		cfg, ok := ent.Config.(*CellularGatewayPoolConfigs)
		if !ok {
			return fmt.Errorf("Unable to add gateway to pool %s; pool has invalid config", id)
		}
		if i == 0 {
			mmeGroupID = cfg.MmeGroupID
		}
		if cfg.MmeGroupID != mmeGroupID {
			return fmt.Errorf("Adding a gateway to pools with different MME group ID's (%d), (%d) is currently unsupported", cfg.MmeGroupID, mmeGroupID)
		}
	}
	return nil
}
