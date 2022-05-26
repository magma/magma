/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package configurator

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/thoas/go-funk"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"
	storage2 "magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
	commonProtos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

// ListNetworkIDs loads a list of all networkIDs registered
func ListNetworkIDs(ctx context.Context) ([]string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	idsWrapper, err := client.ListNetworkIDs(ctx, &commonProtos.Void{})
	if err != nil {
		return nil, err
	}
	return idsWrapper.NetworkIDs, nil
}

// ListNetworksOfType returns a list of all network IDs which match the given
// type
func ListNetworksOfType(ctx context.Context, networkType string) ([]string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	networks, err := client.LoadNetworks(
		ctx,
		&protos.LoadNetworksRequest{
			Criteria: &storage.NetworkLoadCriteria{},
			Filter: &storage.NetworkLoadFilter{
				TypeFilter: strPtrToWrapper(&networkType),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return funk.Map(networks.Networks, func(n *storage.Network) string { return n.ID }).([]string), nil
}

func CreateNetwork(ctx context.Context, network Network, serdes serde.Registry) error {
	_, err := CreateNetworks(ctx, []Network{network}, serdes)
	return err
}

// CreateNetworks registers the given list of Networks and returns the created networks
func CreateNetworks(ctx context.Context, networks []Network, serdes serde.Registry) ([]Network, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	req := &protos.CreateNetworksRequest{Networks: make([]*storage.Network, 0, len(networks))}
	for _, n := range networks {
		pNet, err := n.ToProto(serdes)
		if err != nil {
			return nil, err
		}
		req.Networks = append(req.Networks, pNet)
	}
	res, err := client.CreateNetworks(ctx, req)
	if err != nil {
		return nil, err
	}

	ret := make([]Network, len(res.CreatedNetworks))
	for i, protoNet := range res.CreatedNetworks {
		ent, err := ret[i].FromProto(protoNet, serdes)
		if err != nil {
			return nil, err
		}
		ret[i] = ent
	}
	return ret, nil
}

// UpdateNetworks updates the specified networks and returns the updated networks
func UpdateNetworks(ctx context.Context, updates []NetworkUpdateCriteria, serdes serde.Registry) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}

	req := &protos.UpdateNetworksRequest{Updates: make([]*storage.NetworkUpdateCriteria, 0, len(updates))}
	for _, update := range updates {
		protoUpdate, err := update.toProto(serdes)
		if err != nil {
			return err
		}
		req.Updates = append(req.Updates, protoUpdate)
	}
	_, err = client.UpdateNetworks(ctx, req)
	return err
}

// DeleteNetworks deletes the network specified by networkID
func DeleteNetworks(ctx context.Context, networkIDs []string) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNetworks(ctx, &protos.DeleteNetworksRequest{NetworkIDs: networkIDs})
	return err
}

// DeleteNetwork deletes a network.
func DeleteNetwork(ctx context.Context, networkID string) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNetworks(
		ctx,
		&protos.DeleteNetworksRequest{NetworkIDs: []string{networkID}},
	)
	return err
}

// DoesNetworkExist returns true iff the network exists.
func DoesNetworkExist(ctx context.Context, networkID string) (bool, error) {
	loaded, _, err := LoadNetworks(ctx, []string{networkID}, true, false, nil)
	if err != nil {
		return false, err
	}
	if len(loaded) == 0 {
		return false, nil
	}
	return true, nil
}

// LoadNetworks loads networks networks according to specified criteria.
func LoadNetworks(ctx context.Context, networks []string, loadMetadata bool, loadConfigs bool, serdes serde.Registry) ([]Network, []string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, nil, err
	}
	req := &protos.LoadNetworksRequest{
		Filter: &storage.NetworkLoadFilter{
			Ids: networks,
		},
		Criteria: &storage.NetworkLoadCriteria{
			LoadMetadata: loadMetadata,
			LoadConfigs:  loadConfigs,
		},
	}
	res, err := client.LoadNetworks(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	ret := make([]Network, len(res.Networks))
	for i, n := range res.Networks {
		r, err := ret[i].FromProto(n, serdes)
		if err != nil {
			return nil, nil, err
		}
		ret[i] = r
	}
	return ret, res.NetworkIDsNotFound, nil
}

// LoadNetworksOfType loads all networks of the passed type.
func LoadNetworksOfType(ctx context.Context, typeVal string, loadMetadata bool, loadConfigs bool, serdes serde.Registry) ([]Network, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	req := &protos.LoadNetworksRequest{
		Filter: &storage.NetworkLoadFilter{
			TypeFilter: strPtrToWrapper(&typeVal),
		},
		Criteria: &storage.NetworkLoadCriteria{
			LoadMetadata: loadMetadata,
			LoadConfigs:  loadConfigs,
		},
	}
	res, err := client.LoadNetworks(ctx, req)
	if err != nil {
		return nil, err
	}

	ret := make([]Network, len(res.Networks))
	for i, n := range res.Networks {
		retNet, err := ret[i].FromProto(n, serdes)
		if err != nil {
			return nil, err
		}
		ret[i] = retNet
	}
	return ret, nil
}

// LoadNetwork loads the network identified by the network ID.
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func LoadNetwork(ctx context.Context, networkID string, loadMetadata bool, loadConfigs bool, serdes serde.Registry) (Network, error) {
	networks, _, err := LoadNetworks(ctx, []string{networkID}, loadMetadata, loadConfigs, serdes)
	if err != nil {
		return Network{}, err
	}
	if len(networks) == 0 {
		return Network{}, merrors.ErrNotFound
	}
	return networks[0], nil
}

// LoadNetworkConfig loads network config of type configType registered under the network ID.
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func LoadNetworkConfig(ctx context.Context, networkID, configType string, serdes serde.Registry) (interface{}, error) {
	network, err := LoadNetwork(ctx, networkID, false, true, serdes)
	if err != nil {
		return nil, err
	}
	if network.Configs == nil {
		return nil, merrors.ErrNotFound
	}
	if _, exists := network.Configs[configType]; !exists {
		return nil, merrors.ErrNotFound
	}
	return network.Configs[configType], nil
}

func UpdateNetworkConfig(ctx context.Context, networkID, configType string, config interface{}, serdes serde.Registry) error {
	updateCriteria := NetworkUpdateCriteria{
		ID:                   networkID,
		ConfigsToAddOrUpdate: map[string]interface{}{configType: config},
	}
	return UpdateNetworks(ctx, []NetworkUpdateCriteria{updateCriteria}, serdes)
}

// WriteEntities executes a series of entity writes (creation or update) to be
// executed in order within a single transaction.
// This function is all-or-nothing - any failure or error encountered during
// any operation will rollback the entire batch.
func WriteEntities(ctx context.Context, networkID string, writes []EntityWriteOperation, serdes serde.Registry) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}

	req := &protos.WriteEntitiesRequest{NetworkID: networkID}
	for _, write := range writes {
		switch op := write.(type) {
		case NetworkEntity:
			protoEnt, err := op.toProto(serdes)
			if err != nil {
				return err
			}
			req.Writes = append(req.Writes, &protos.WriteEntityRequest{Request: &protos.WriteEntityRequest_Create{Create: protoEnt}})
		case EntityUpdateCriteria:
			protoEuc, err := op.toProto(serdes)
			if err != nil {
				return err
			}
			req.Writes = append(req.Writes, &protos.WriteEntityRequest{Request: &protos.WriteEntityRequest_Update{Update: protoEuc}})
		default:
			return fmt.Errorf("unrecognized entity write operation %T", op)
		}
	}

	_, err = client.WriteEntities(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// CreateEntity creates a network entity.
func CreateEntity(ctx context.Context, networkID string, entity NetworkEntity, serdes serde.Registry) (NetworkEntity, error) {
	ret, err := CreateEntities(ctx, networkID, NetworkEntities{entity}, serdes)
	if err != nil {
		return NetworkEntity{}, err
	}
	return ret[0], nil
}

// CreateEntities registers the given entities and returns the created network
// entities.
func CreateEntities(ctx context.Context, networkID string, entities NetworkEntities, serdes serde.Registry) (NetworkEntities, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	req := &protos.CreateEntitiesRequest{NetworkID: networkID, Entities: make([]*storage.NetworkEntity, 0, len(entities))}
	for _, ent := range entities {
		protoEnt, err := ent.toProto(serdes)
		if err != nil {
			return nil, err
		}
		req.Entities = append(req.Entities, protoEnt)
	}
	res, err := client.CreateEntities(ctx, req)
	if err != nil {
		return nil, err
	}

	created := (NetworkEntities{}).fromProtosSerialized(res.CreatedEntities)
	return created, nil
}

// CreateInternalEntity is a loose wrapper around CreateEntity to create an
// entity in the internal network structure
func CreateInternalEntity(ctx context.Context, entity NetworkEntity, serdes serde.Registry) (NetworkEntity, error) {
	return CreateEntity(ctx, storage.InternalNetworkID, entity, serdes)
}

// UpdateEntity updates a network entity.
func UpdateEntity(ctx context.Context, networkID string, update EntityUpdateCriteria, serdes serde.Registry) (NetworkEntity, error) {
	updates, err := UpdateEntities(ctx, networkID, []EntityUpdateCriteria{update}, serdes)
	if err != nil {
		return NetworkEntity{}, err
	}
	for _, e := range updates {
		return e, nil
	}
	return NetworkEntity{}, merrors.ErrNotFound
}

// UpdateEntities updates the registered entities and returns the updated entities
func UpdateEntities(ctx context.Context, networkID string, updates []EntityUpdateCriteria, serdes serde.Registry) (NetworkEntities, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	req := &protos.UpdateEntitiesRequest{NetworkID: networkID, Updates: make([]*storage.EntityUpdateCriteria, 0, len(updates))}
	for _, update := range updates {
		upProto, err := update.toProto(serdes)
		if err != nil {
			return nil, err
		}
		req.Updates = append(req.Updates, upProto)
	}
	res, err := client.UpdateEntities(ctx, req)
	if err != nil {
		return nil, err
	}

	updatedEnts := funk.Values(res.UpdatedEntities).([]*storage.NetworkEntity)
	ret := (NetworkEntities{}).fromProtosSerialized(updatedEnts)

	return ret, nil
}

// UpdateInternalEntity is a loose wrapper around UpdateEntity to update an
// entity in the internal network structure.
func UpdateInternalEntity(ctx context.Context, update EntityUpdateCriteria, serdes serde.Registry) (NetworkEntity, error) {
	return UpdateEntity(ctx, storage.InternalNetworkID, update, serdes)
}

func CreateOrUpdateEntityConfig(ctx context.Context, networkID string, entityType string, entityKey string, config interface{}, serdes serde.Registry) error {
	updateCriteria := EntityUpdateCriteria{
		Key:       entityKey,
		Type:      entityType,
		NewConfig: config,
	}
	_, err := UpdateEntities(ctx, networkID, []EntityUpdateCriteria{updateCriteria}, serdes)
	return err
}

func DeleteEntity(ctx context.Context, networkID string, entityType string, entityKey string) error {
	return DeleteEntities(ctx, networkID, storage2.TKs{{Type: entityType, Key: entityKey}})
}

// DeleteEntities deletes the entities specified by networkID and tks.
// We also have cascading deletes to delete foreign keys for assocs.
func DeleteEntities(ctx context.Context, networkID string, ids storage2.TKs) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteEntities(
		ctx,
		&protos.DeleteEntitiesRequest{
			NetworkID: networkID,
			ID:        tksToEntIDs(ids),
		},
	)
	return err
}

// DeleteInternalEntity is a loose wrapper around DeleteEntities to delete an
// entity in the internal network structure
func DeleteInternalEntity(ctx context.Context, entityType, entityKey string) error {
	return DeleteEntity(ctx, storage.InternalNetworkID, entityType, entityKey)
}

// GetPhysicalIDOfEntity gets the physicalID associated with the entity identified by (networkID, entityType, entityKey)
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func GetPhysicalIDOfEntity(ctx context.Context, networkID, entityType, entityKey string) (string, error) {
	entity, err := LoadSerializedEntity(ctx, networkID, entityType, entityKey, EntityLoadCriteria{})
	if err != nil {
		return "", err
	}
	return entity.PhysicalID, nil
}

// ListEntityKeys returns all keys for an entity type in a network.
func ListEntityKeys(ctx context.Context, networkID string, entityType string) ([]string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return []string{}, err
	}

	networkExists, _ := DoesNetworkExist(ctx, networkID)
	if !networkExists {
		return []string{}, merrors.ErrNotFound
	}

	res, err := client.LoadEntities(
		ctx,
		&protos.LoadEntitiesRequest{
			NetworkID: networkID,
			Filter: &storage.EntityLoadFilter{
				TypeFilter: &wrappers.StringValue{Value: entityType},
			},
			Criteria: EntityLoadCriteria{}.toProto(),
		},
	)
	if err != nil {
		return []string{}, err
	}

	return funk.Map(res.Entities, func(ent *storage.NetworkEntity) string { return ent.Key }).([]string), nil
}

// ListInternalEntityKeys calls ListEntityKeys with the internal networkID
func ListInternalEntityKeys(ctx context.Context, entityType string) ([]string, error) {
	return ListEntityKeys(ctx, storage.InternalNetworkID, entityType)
}

// LoadEntity loads the network entity identified by (network ID, entity type, entity key).
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func LoadEntity(ctx context.Context, networkID string, entityType string, entityKey string, criteria EntityLoadCriteria, serdes serde.Registry) (NetworkEntity, error) {
	ret := NetworkEntity{}
	loaded, notFound, err := LoadEntities(
		ctx,
		networkID,
		nil, nil, nil,
		storage2.TKs{{Type: entityType, Key: entityKey}},
		criteria,
		serdes,
	)
	if err != nil {
		return ret, err
	}
	if len(notFound) != 0 || len(loaded) == 0 {
		return ret, merrors.ErrNotFound
	}
	return loaded[0], nil
}

// LoadEntityConfig loads the config for the entity identified by (network ID, entity type, entity key).
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func LoadEntityConfig(ctx context.Context, networkID, entityType, entityKey string, serdes serde.Registry) (interface{}, error) {
	entity, err := LoadEntity(ctx, networkID, entityType, entityKey, EntityLoadCriteria{LoadConfig: true}, serdes)
	if err != nil {
		return nil, err
	}
	if entity.Config == nil {
		return nil, merrors.ErrNotFound
	}
	return entity.Config, nil
}

// LoadEntityForPhysicalID loads the network entity identified by the physical ID.
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func LoadEntityForPhysicalID(ctx context.Context, physicalID string, criteria EntityLoadCriteria, serdes serde.Registry) (NetworkEntity, error) {
	ret := NetworkEntity{}
	loaded, _, err := LoadEntities(
		ctx,
		"placeholder",
		nil, nil, &physicalID,
		nil,
		criteria,
		serdes,
	)
	if err != nil {
		return ret, err
	}
	if len(loaded) == 0 {
		return ret, merrors.ErrNotFound
	}
	if len(loaded) > 1 {
		return ret, fmt.Errorf("expected one entity from query, found %d", len(loaded))
	}
	return loaded[0], nil
}

// LoadEntities loads entities specified by the parameters.
// typeFilter, keyFilter, physicalID, and ids are all used to define a filter to
// filter out results - if they are all nil, it will return all network entities
// If ids is empty, all entities will be returned
func LoadEntities(ctx context.Context, networkID string, typeFilter *string, keyFilter *string, physicalID *string, ids storage2.TKs, criteria EntityLoadCriteria, serdes serde.Registry) (NetworkEntities, storage2.TKs, error) {
	protoEnts, notFound, err := loadEntities(ctx, networkID, typeFilter, keyFilter, physicalID, ids, criteria)
	if err != nil {
		return nil, nil, err
	}
	ret, err := (NetworkEntities{}).fromProtos(protoEnts, serdes)
	if err != nil {
		return nil, nil, fmt.Errorf("request succeeded but deserialization failed: %w", err)
	}
	return ret, entIDsToTKs(notFound), nil
}

// LoadSerializedEntity is same as LoadEntity, but doesn't deserialize the
// loaded entity.
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func LoadSerializedEntity(ctx context.Context, networkID string, entityType string, entityKey string, criteria EntityLoadCriteria) (NetworkEntity, error) {
	ret := NetworkEntity{}
	loaded, notFound, err := LoadSerializedEntities(
		ctx,
		networkID,
		nil, nil, nil,
		storage2.TKs{{Type: entityType, Key: entityKey}},
		criteria,
	)
	if err != nil {
		return ret, err
	}
	if len(notFound) != 0 || len(loaded) == 0 {
		return ret, merrors.ErrNotFound
	}
	return loaded[0], nil
}

// LoadSerializedEntities is same as LoadEntities, but doesn't deserialize
// the loaded entities.
func LoadSerializedEntities(
	ctx context.Context,
	networkID string,
	typeFilter *string,
	keyFilter *string,
	physicalID *string,
	ids storage2.TKs,
	criteria EntityLoadCriteria,
) (NetworkEntities, storage2.TKs, error) {
	protoEnts, notFound, err := loadEntities(ctx, networkID, typeFilter, keyFilter, physicalID, ids, criteria)
	if err != nil {
		return nil, nil, err
	}
	ents := (NetworkEntities{}).fromProtosSerialized(protoEnts)
	return ents, entIDsToTKs(notFound), nil
}

func loadEntities(
	ctx context.Context,
	networkID string,
	typeFilter *string,
	keyFilter *string,
	physicalID *string,
	ids storage2.TKs,
	criteria EntityLoadCriteria,
) ([]*storage.NetworkEntity, []*storage.EntityID, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, nil, err
	}

	res, err := client.LoadEntities(
		ctx,
		&protos.LoadEntitiesRequest{
			NetworkID: networkID,
			Filter: &storage.EntityLoadFilter{
				TypeFilter: protos.GetStringWrapper(typeFilter),
				KeyFilter:  protos.GetStringWrapper(keyFilter),
				PhysicalID: protos.GetStringWrapper(physicalID),
				IDs:        tksToEntIDs(ids),
			},
			Criteria: criteria.toProto(),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return res.Entities, res.EntitiesNotFound, nil
}

// LoadInternalEntity calls LoadEntity with the internal network ID.
// If not found, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func LoadInternalEntity(ctx context.Context, entityType string, entityKey string, criteria EntityLoadCriteria, serdes serde.Registry) (NetworkEntity, error) {
	return LoadEntity(ctx, storage.InternalNetworkID, entityType, entityKey, criteria, serdes)
}

// DoesEntityExist returns a boolean that indicated whether the entity specified
// exists in the network
func DoesEntityExist(ctx context.Context, networkID, entityType, entityKey string) (bool, error) {
	found, _, err := loadEntities(
		ctx,
		networkID,
		nil, nil, nil,
		storage2.TKs{{Type: entityType, Key: entityKey}},
		EntityLoadCriteria{},
	)
	if err != nil {
		return false, err
	}
	if len(found) != 1 {
		return false, nil
	}
	return true, nil
}

// DoEntitiesExist returns a boolean that indicated whether all entities
// specified exist in the network
func DoEntitiesExist(ctx context.Context, networkID string, ids storage2.TKs) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	found, _, err := loadEntities(
		ctx,
		networkID,
		nil, nil, nil,
		ids,
		EntityLoadCriteria{},
	)
	if err != nil {
		return false, err
	}
	if len(found) != len(ids) {
		return false, nil
	}
	return true, nil
}

// DoesInternalEntityExist calls DoesEntityExist with the internal networkID
func DoesInternalEntityExist(ctx context.Context, entityType, entityKey string) (bool, error) {
	return DoesEntityExist(ctx, storage.InternalNetworkID, entityType, entityKey)
}

// LoadAllEntitiesOfType fetches all entities of specified type in a network.
// Loads can be paginated by specifying a page size and token in the entity
// load criteria. To exhaustively read all pages, clients must continue
// querying until an empty page token is received in the load result.
func LoadAllEntitiesOfType(ctx context.Context, networkID string, entityType string, criteria EntityLoadCriteria, serdes serde.Registry) (NetworkEntities, string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, "", err
	}

	res, err := client.LoadEntities(
		ctx,
		&protos.LoadEntitiesRequest{
			NetworkID: networkID,
			Filter: &storage.EntityLoadFilter{
				TypeFilter: &wrappers.StringValue{Value: entityType},
			},
			Criteria: criteria.toProto(),
		},
	)
	if err != nil {
		return nil, "", err
	}

	ret, err := (NetworkEntities{}).fromProtos(res.Entities, serdes)
	if err != nil {
		return nil, "", fmt.Errorf("request succeeded but deserialization failed: %w", err)
	}

	return ret, res.NextPageToken, nil
}

// CountEntitiesOfType provides total count of entities of this type
func CountEntitiesOfType(ctx context.Context, networkID string, entityType string) (uint64, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return 0, err
	}
	res, err := client.CountEntities(
		ctx,
		&protos.LoadEntitiesRequest{
			NetworkID: networkID,
			Filter: &storage.EntityLoadFilter{
				TypeFilter: &wrappers.StringValue{Value: entityType},
			},
		},
	)
	if err != nil {
		return 0, err
	}
	return res.Count, nil
}

func getNBConfiguratorClient() (protos.NorthboundConfiguratorClient, error) {
	conn, err := registry.GetConnection(ServiceName, commonProtos.ServiceType_PROTECTED)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewNorthboundConfiguratorClient(conn), err
}

func GetMconfigFor(ctx context.Context, hardwareID string) (*protos.GetMconfigResponse, error) {
	client, err := getSBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	return client.GetMconfigInternal(ctx, &protos.GetMconfigRequest{HardwareID: hardwareID})
}

func getSBConfiguratorClient() (protos.SouthboundConfiguratorClient, error) {
	conn, err := registry.GetConnection(ServiceName, commonProtos.ServiceType_SOUTHBOUND)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewSouthboundConfiguratorClient(conn), err
}
