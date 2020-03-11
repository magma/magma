/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package configurator

import (
	"context"
	"fmt"

	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"
	storage2 "magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	commonProtos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func getNBConfiguratorClient() (protos.NorthboundConfiguratorClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewNorthboundConfiguratorClient(conn), err
}

// ListNetworkIDs loads a list of all networkIDs registered
func ListNetworkIDs() ([]string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	idsWrapper, err := client.ListNetworkIDs(context.Background(), &commonProtos.Void{})
	if err != nil {
		return nil, err
	}
	return idsWrapper.NetworkIDs, nil
}

// ListNetworksOfType returns a list of all network IDs which match the given
// type
func ListNetworksOfType(networkType string) ([]string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	networks, err := client.LoadNetworks(
		context.Background(),
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

func CreateNetwork(network Network) error {
	_, err := CreateNetworks([]Network{network})
	return err
}

// CreateNetworks registers the given list of Networks and returns the created networks
func CreateNetworks(networks []Network) ([]Network, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	req := &protos.CreateNetworksRequest{Networks: make([]*storage.Network, 0, len(networks))}
	for _, n := range networks {
		pNet, err := n.toStorageProto()
		if err != nil {
			return nil, err
		}
		req.Networks = append(req.Networks, pNet)
	}
	result, err := client.CreateNetworks(context.Background(), req)
	if err != nil {
		return nil, err
	}

	ret := make([]Network, len(result.CreatedNetworks))
	for i, protoNet := range result.CreatedNetworks {
		ent, err := ret[i].fromStorageProto(protoNet)
		if err != nil {
			return nil, err
		}
		ret[i] = ent
	}
	return ret, nil
}

// UpdateNetworks updates the specified networks and returns the updated networks
func UpdateNetworks(updates []NetworkUpdateCriteria) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}

	request := &protos.UpdateNetworksRequest{Updates: make([]*storage.NetworkUpdateCriteria, 0, len(updates))}
	for _, update := range updates {
		protoUpdate, err := update.toStorageProto()
		if err != nil {
			return err
		}
		request.Updates = append(request.Updates, protoUpdate)
	}
	_, err = client.UpdateNetworks(context.Background(), request)
	return err
}

// DeleteNetworks deletes the network specified by networkID
func DeleteNetworks(networkIDs []string) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNetworks(context.Background(), &protos.DeleteNetworksRequest{NetworkIDs: networkIDs})
	return err
}

// DeleteNetwork deletes a network.
func DeleteNetwork(networkID string) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNetworks(
		context.Background(),
		&protos.DeleteNetworksRequest{NetworkIDs: []string{networkID}},
	)
	return err
}

// DoesNetworkExist returns a boolean that indicates whether the networkID
func DoesNetworkExist(networkID string) (bool, error) {
	loaded, _, err := LoadNetworks([]string{networkID}, true, false)
	if err != nil {
		return false, err
	}
	if len(loaded) == 0 {
		return false, nil
	}
	return true, nil
}

// LoadNetworks loads networks specified by networks according to criteria specified and
// returns the result
func LoadNetworks(networks []string, loadMetadata bool, loadConfigs bool) ([]Network, []string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, nil, err
	}
	request := &protos.LoadNetworksRequest{
		Filter: &storage.NetworkLoadFilter{
			Ids: networks,
		},
		Criteria: &storage.NetworkLoadCriteria{
			LoadMetadata: loadMetadata,
			LoadConfigs:  loadConfigs,
		},
	}
	result, err := client.LoadNetworks(context.Background(), request)
	if err != nil {
		return nil, nil, err
	}

	ret := make([]Network, len(result.Networks))
	for i, n := range result.Networks {
		retNet, err := ret[i].fromStorageProto(n)
		if err != nil {
			return nil, nil, err
		}
		ret[i] = retNet
	}
	return ret, result.NetworkIDsNotFound, nil
}

func LoadNetworksByType(typeVal string, loadMetadata bool, loadConfigs bool) ([]Network, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	request := &protos.LoadNetworksRequest{
		Filter: &storage.NetworkLoadFilter{
			TypeFilter: strPtrToWrapper(&typeVal),
		},
		Criteria: &storage.NetworkLoadCriteria{
			LoadMetadata: loadMetadata,
			LoadConfigs:  loadConfigs,
		},
	}
	result, err := client.LoadNetworks(context.Background(), request)
	if err != nil {
		return nil, err
	}

	ret := make([]Network, len(result.Networks))
	for i, n := range result.Networks {
		retNet, err := ret[i].fromStorageProto(n)
		if err != nil {
			return nil, err
		}
		ret[i] = retNet
	}
	return ret, nil
}

func LoadNetwork(networkID string, loadMetadata bool, loadConfigs bool) (Network, error) {
	networks, _, err := LoadNetworks([]string{networkID}, loadMetadata, loadConfigs)
	if err != nil {
		return Network{}, err
	}
	if len(networks) == 0 {
		return Network{}, merrors.ErrNotFound
	}
	return networks[0], nil
}

// LoadNetworkConfig loads network config of type configType registered under the networkID
func LoadNetworkConfig(networkID, configType string) (interface{}, error) {
	network, err := LoadNetwork(networkID, false, true)
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

func UpdateNetworkConfig(networkID, configType string, config interface{}) error {
	updateCriteria := NetworkUpdateCriteria{
		ID:                   networkID,
		ConfigsToAddOrUpdate: map[string]interface{}{configType: config},
	}
	return UpdateNetworks([]NetworkUpdateCriteria{updateCriteria})
}

func DeleteNetworkConfig(networkID, configType string) error {
	updateCriteria := NetworkUpdateCriteria{
		ID:              networkID,
		ConfigsToDelete: []string{configType},
	}
	return UpdateNetworks([]NetworkUpdateCriteria{updateCriteria})
}

func GetNetworkConfigsByType(networkID string, configType string) (interface{}, error) {
	networks, _, err := LoadNetworks([]string{networkID}, false, true)
	if err != nil {
		return nil, err
	}
	if len(networks) == 0 {
		return nil, fmt.Errorf("Network %s not found", networkID)
	}
	return networks[0].Configs[configType], nil
}

// WriteEntities executes a series of entity writes (creation or update) to be
// executed in order within a single transaction.
// This function is all-or-nothing - any failure or error encountered during
// any operation will rollback the entire batch.
func WriteEntities(networkID string, writes ...EntityWriteOperation) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}

	req := &protos.WriteEntitiesRequest{NetworkID: networkID}
	for _, write := range writes {
		switch op := write.(type) {
		case NetworkEntity:
			protoEnt, err := op.toStorageProto()
			if err != nil {
				return err
			}
			req.Writes = append(req.Writes, &protos.WriteEntityRequest{Request: &protos.WriteEntityRequest_Create{Create: protoEnt}})
		case EntityUpdateCriteria:
			protoEuc, err := op.toStorageProto()
			if err != nil {
				return err
			}
			req.Writes = append(req.Writes, &protos.WriteEntityRequest{Request: &protos.WriteEntityRequest_Update{Update: protoEuc}})
		default:
			return errors.Errorf("unrecognized entity write operation %T", op)
		}
	}

	_, err = client.WriteEntities(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func CreateEntity(networkID string, entity NetworkEntity) (NetworkEntity, error) {
	ret, err := CreateEntities(networkID, []NetworkEntity{entity})
	if err != nil {
		return NetworkEntity{}, err
	}
	return ret[0], nil
}

// CreateEntities registers the given entities and returns the created network entities
func CreateEntities(networkID string, entities []NetworkEntity) ([]NetworkEntity, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	request := &protos.CreateEntitiesRequest{NetworkID: networkID, Entities: make([]*storage.NetworkEntity, 0, len(entities))}
	for _, ent := range entities {
		protoEnt, err := ent.toStorageProto()
		if err != nil {
			return nil, err
		}
		request.Entities = append(request.Entities, protoEnt)
	}
	response, err := client.CreateEntities(context.Background(), request)
	if err != nil {
		return nil, err
	}

	ret := make([]NetworkEntity, len(response.CreatedEntities))
	for i, protoEnt := range response.CreatedEntities {
		ent, err := ret[i].fromStorageProto(protoEnt)
		if err != nil {
			return nil, errors.Wrap(err, "request succeeded but deserialization failed")
		}
		ret[i] = ent
	}
	return ret, err
}

// CreateInternalEntity is a loose wrapper around CreateEntity to create an
// entity in the internal network structure
func CreateInternalEntity(entity NetworkEntity) (NetworkEntity, error) {
	return CreateEntity(storage.InternalNetworkID, entity)
}

func UpdateEntity(networkID string, update EntityUpdateCriteria) (NetworkEntity, error) {
	retMap, err := UpdateEntities(networkID, []EntityUpdateCriteria{update})
	if err != nil {
		return NetworkEntity{}, err
	}
	for _, v := range retMap {
		return v, nil
	}
	return NetworkEntity{}, merrors.ErrNotFound
}

// UpdateEntities updates the registered entities and returns the updated entities
func UpdateEntities(networkID string, updates []EntityUpdateCriteria) (map[string]NetworkEntity, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	request := &protos.UpdateEntitiesRequest{NetworkID: networkID, Updates: make([]*storage.EntityUpdateCriteria, 0, len(updates))}
	for _, update := range updates {
		upProto, err := update.toStorageProto()
		if err != nil {
			return nil, err
		}
		request.Updates = append(request.Updates, upProto)
	}
	response, err := client.UpdateEntities(context.Background(), request)
	if err != nil {
		return nil, err
	}

	ret := map[string]NetworkEntity{}
	for id, protoEnt := range response.UpdatedEntities {
		ent, err := (NetworkEntity{}).fromStorageProto(protoEnt)
		if err != nil {
			return nil, errors.Wrap(err, "request succeeded but response deserialization failed")
		}
		ret[id] = ent
	}
	return ret, err
}

// UpdateInternalEntity is a loose wrapper around UpdateEntity to update an
// entity in the internal network structure
func UpdateInternalEntity(update EntityUpdateCriteria) (NetworkEntity, error) {
	return UpdateEntity(storage.InternalNetworkID, update)
}

func CreateOrUpdateEntityConfig(networkID string, entityType string, entityKey string, config interface{}) error {
	updateCriteria := EntityUpdateCriteria{
		Key:       entityKey,
		Type:      entityType,
		NewConfig: config,
	}
	_, err := UpdateEntities(networkID, []EntityUpdateCriteria{updateCriteria})
	return err
}

func CreateOrUpdateEntityConfigAndAssoc(networkID string, entityType string, entityKey string, config interface{}, updatedAssoc []storage2.TypeAndKey) error {
	// first delete old associations
	updateCriteria := EntityUpdateCriteria{
		Key:               entityKey,
		Type:              entityType,
		NewConfig:         config,
		AssociationsToSet: updatedAssoc,
	}
	_, err := UpdateEntities(networkID, []EntityUpdateCriteria{updateCriteria})
	return err
}

func DeleteEntityConfig(networkID, entityType, entityKey string) error {
	updateCriteria := EntityUpdateCriteria{
		Key:          entityKey,
		Type:         entityType,
		DeleteConfig: true,
	}
	_, err := UpdateEntities(networkID, []EntityUpdateCriteria{updateCriteria})
	return err
}

func DeleteEntity(networkID string, entityType string, entityKey string) error {
	return DeleteEntities(networkID, []storage2.TypeAndKey{{Type: entityType, Key: entityKey}})
}

// DeleteEntity deletes the entity specified by networkID, type, key
// We also have cascading deletes to delete foreign keys for assocs
func DeleteEntities(networkID string, ids []storage2.TypeAndKey) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteEntities(
		context.Background(),
		&protos.DeleteEntitiesRequest{
			NetworkID: networkID,
			ID:        tksToEntIDs(ids),
		},
	)
	return err
}

// DeleteInternalEntity is a loose wrapper around DeleteEntities to delete an
// entity in the internal network structure
func DeleteInternalEntity(entityType, entityKey string) error {
	return DeleteEntity(storage.InternalNetworkID, entityType, entityKey)
}

// GetPhysicalIDOfEntity gets the physicalID associated with the entity
// identified by (networkID, entityType, entityKey)
func GetPhysicalIDOfEntity(networkID, entityType, entityKey string) (string, error) {
	entities, _, err := LoadEntities(
		networkID,
		nil, nil, nil,
		[]storage2.TypeAndKey{
			{Type: entityType, Key: entityKey},
		},
		EntityLoadCriteria{},
	)
	if err != nil || len(entities) != 1 {
		return "", err
	}
	return entities[0].PhysicalID, nil
}

// ListEntityKeys returns all keys for an entity type in a network.
func ListEntityKeys(networkID string, entityType string) ([]string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return []string{}, err
	}
	networkExists, _ := DoesNetworkExist(networkID)
	if !networkExists {
		return []string{}, merrors.ErrNotFound
	}

	resp, err := client.LoadEntities(
		context.Background(),
		&protos.LoadEntitiesRequest{
			NetworkID: networkID,
			Filter: &storage.EntityLoadFilter{
				TypeFilter: &wrappers.StringValue{Value: entityType},
			},
			Criteria: EntityLoadCriteria{}.toStorageProto(),
		},
	)
	if err != nil {
		return []string{}, err
	}

	return funk.Map(resp.Entities, func(ent *storage.NetworkEntity) string { return ent.Key }).([]string), nil
}

// ListInternalEntityKeys calls ListEntityKeys with the internal networkID
func ListInternalEntityKeys(entityType string) ([]string, error) {
	return ListEntityKeys(storage.InternalNetworkID, entityType)
}

func LoadEntity(networkID string, entityType string, entityKey string, criteria EntityLoadCriteria) (NetworkEntity, error) {
	ret := NetworkEntity{}
	loaded, notFound, err := LoadEntities(
		networkID,
		nil, nil, nil,
		[]storage2.TypeAndKey{{Type: entityType, Key: entityKey}},
		criteria,
	)
	if err != nil {
		return ret, err
	}
	if !funk.IsEmpty(notFound) || funk.IsEmpty(loaded) {
		return ret, merrors.ErrNotFound
	}
	return loaded[0], nil
}

func LoadEntityConfig(networkID, entityType, entityKey string) (interface{}, error) {
	entity, err := LoadEntity(networkID, entityType, entityKey, EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return nil, err
	}
	if entity.Config == nil {
		return nil, merrors.ErrNotFound
	}
	return entity.Config, nil
}

func LoadEntityForPhysicalID(physicalID string, criteria EntityLoadCriteria) (NetworkEntity, error) {
	ret := NetworkEntity{}
	loaded, _, err := LoadEntities(
		"placeholder",
		nil, nil, &physicalID, nil,
		criteria,
	)
	if err != nil {
		return ret, err
	}
	if funk.IsEmpty(loaded) {
		return ret, merrors.ErrNotFound
	}
	if len(loaded) > 1 {
		return ret, errors.Errorf("expected one entity from query, found %d", len(loaded))
	}
	return loaded[0], nil
}

func GetNetworkAndEntityIDForPhysicalID(physicalID string) (string, string, error) {
	if len(physicalID) == 0 {
		return "", "", errors.New("Empty Hardware ID")
	}
	entity, err := LoadEntityForPhysicalID(physicalID, EntityLoadCriteria{})
	if err != nil {
		return "", "", err
	}
	return entity.NetworkID, entity.Key, nil
}

// LoadEntities loads entities specified by the parameters.
// typeFilter, keyFilter, physicalID, and ids are all used to define a filter to
// filter out results - if they are all nil, it will return all network entities
func LoadEntities(
	networkID string,
	typeFilter *string,
	keyFilter *string,
	physicalID *string,
	ids []storage2.TypeAndKey,
	criteria EntityLoadCriteria,
) (NetworkEntities, []storage2.TypeAndKey, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.LoadEntities(
		context.Background(),
		&protos.LoadEntitiesRequest{
			NetworkID: networkID,
			Filter: &storage.EntityLoadFilter{
				TypeFilter: protos.GetStringWrapper(typeFilter),
				KeyFilter:  protos.GetStringWrapper(keyFilter),
				PhysicalID: protos.GetStringWrapper(physicalID),
				IDs:        tksToEntIDs(ids),
			},
			Criteria: criteria.toStorageProto(),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	ret := make([]NetworkEntity, len(resp.Entities))
	for i, protoEnt := range resp.Entities {
		ent, err := ret[i].fromStorageProto(protoEnt)
		if err != nil {
			return nil, nil, errors.Wrap(err, "request succeeded but deserialization failed")
		}
		ret[i] = ent
	}
	return ret, entIDsToTKs(resp.EntitiesNotFound), nil
}

// LoadInternalEntity calls LoadEntity with the internal networkID
func LoadInternalEntity(entityType string, entityKey string, criteria EntityLoadCriteria) (NetworkEntity, error) {
	return LoadEntity(storage.InternalNetworkID, entityType, entityKey, criteria)
}

// DoesEntityExist returns a boolean that indicated whether the entity specified
// exists in the network
func DoesEntityExist(networkID, entityType, entityKey string) (bool, error) {
	found, _, err := LoadEntities(
		networkID,
		nil, nil, nil,
		[]storage2.TypeAndKey{{Type: entityType, Key: entityKey}},
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
func DoEntitiesExist(networkID string, ids []storage2.TypeAndKey) (bool, error) {
	found, _, err := LoadEntities(
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
func DoesInternalEntityExist(entityType, entityKey string) (bool, error) {
	return DoesEntityExist(storage.InternalNetworkID, entityType, entityKey)
}

// LoadAllEntitiesInNetwork fetches all entities of specified type in a network
func LoadAllEntitiesInNetwork(networkID string, entityType string, criteria EntityLoadCriteria) ([]NetworkEntity, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.LoadEntities(
		context.Background(),
		&protos.LoadEntitiesRequest{
			NetworkID: networkID,
			Filter: &storage.EntityLoadFilter{
				TypeFilter: &wrappers.StringValue{Value: entityType},
			},
			Criteria: criteria.toStorageProto(),
		},
	)
	if err != nil {
		return nil, err
	}

	ret := make([]NetworkEntity, len(resp.Entities))
	for i, protoEnt := range resp.Entities {
		ent, err := ret[i].fromStorageProto(protoEnt)
		if err != nil {
			return nil, errors.Wrapf(err, "request succeeded but deserialization failed")
		}
		ret[i] = ent
	}
	return ret, nil
}

func getSBConfiguratorClient() (protos.SouthboundConfiguratorClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewSouthboundConfiguratorClient(conn), err
}

func GetMconfigFor(hardwareID string) (*protos.GetMconfigResponse, error) {
	client, err := getSBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	return client.GetMconfigInternal(context.Background(), &protos.GetMconfigRequest{HardwareID: hardwareID})
}
