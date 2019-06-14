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

	"magma/orc8r/cloud/go/errors"
	commonProtos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func getNBConfiguratorClient() (protos.NorthboundConfiguratorClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
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

// CreateNetworks registers the given list of Networks and returns the created networks
func CreateNetworks(networks []*storage.Network) ([]*storage.Network, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}
	request := &protos.CreateNetworksRequest{Networks: networks}
	result, err := client.CreateNetworks(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return result.CreatedNetworks, err
}

// UpdateNetworks updates the specified networks and returns the updated networks
func UpdateNetworks(updates []*storage.NetworkUpdateCriteria) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	request := &protos.UpdateNetworksRequest{Updates: updates}
	_, err = client.UpdateNetworks(context.Background(), request)
	return err
}

// DeleteNetwork deletes the network specified by networkID
func DeleteNetworks(networkIDs []string) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNetworks(context.Background(), &protos.DeleteNetworksRequest{NetworkIDs: networkIDs})
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
func LoadNetworks(networks []string, loadMetadata bool, loadConfigs bool) ([]*storage.Network, []string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, nil, err
	}
	request := &protos.LoadNetworksRequest{
		Networks: networks,
		Criteria: &storage.NetworkLoadCriteria{
			LoadMetadata: loadMetadata,
			LoadConfigs:  loadConfigs,
		},
	}
	result, err := client.LoadNetworks(context.Background(), request)
	if err != nil {
		return nil, nil, err
	}
	return result.Networks, result.NetworkIDsNotFound, nil
}

func UpdateNetworkConfig(networkID, configType string, config interface{}) error {
	serializedConfig, err := serde.Serialize(NetworkConfigSerdeDomain, configType, config)
	if err != nil {
		return err
	}
	configMap := map[string][]byte{}
	configMap[configType] = serializedConfig
	updateCriteria := &storage.NetworkUpdateCriteria{
		ID:                   networkID,
		ConfigsToAddOrUpdate: configMap,
	}
	return UpdateNetworks([]*storage.NetworkUpdateCriteria{updateCriteria})
}

func DeleteNetworkConfig(networkID, configType string) error {
	updateCriteria := &storage.NetworkUpdateCriteria{
		ID:              networkID,
		ConfigsToDelete: []string{configType},
	}
	return UpdateNetworks([]*storage.NetworkUpdateCriteria{updateCriteria})
}

func GetNetworkConfigsByType(networkID string, configType string) (interface{}, error) {
	networks, _, err := LoadNetworks([]string{networkID}, false, true)
	if err != nil {
		return nil, err
	}
	if len(networks) == 0 {
		return nil, fmt.Errorf("Network %s not found", networkID)
	}
	serializedConfig := networks[0].Configs[configType]
	model, err := serde.Deserialize(NetworkConfigSerdeDomain, configType, serializedConfig)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// CreateEntities registers the given entities and returns the created network entities
func CreateEntities(networkID string, entities []*storage.NetworkEntity) ([]*storage.NetworkEntity, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	request := &protos.CreateEntitiesRequest{NetworkID: networkID, Entities: entities}
	response, err := client.CreateEntities(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response.CreatedEntities, err
}

// CreateInternalEntity is a loose wrapper around CreateEntities to create an
// entity in the internal network structure
func CreateInternalEntities(entities []*storage.NetworkEntity) ([]*storage.NetworkEntity, error) {
	return CreateEntities(storage.InternalNetworkID, entities)
}

// UpdateEntities updates the registered entities and returns the updated entities
func UpdateEntities(networkID string, updates []*storage.EntityUpdateCriteria) (map[string]*storage.NetworkEntity, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	request := &protos.UpdateEntitiesRequest{NetworkID: networkID, Updates: updates}
	response, err := client.UpdateEntities(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response.UpdatedEntities, err
}

// UpdateInternalEntity is a loose wrapper around UpdateEntities to update an
// entity in the internal network structure
func UpdateInternalEntity(updates []*storage.EntityUpdateCriteria) (map[string]*storage.NetworkEntity, error) {
	return UpdateEntities(storage.InternalNetworkID, updates)
}

func UpdateEntityConfig(networkID string, entityType string, entityKey string, config interface{}) error {
	serializedConfig, err := serde.Serialize(NetworkEntitySerdeDomain, entityType, config)
	if err != nil {
		return err
	}
	updateCriteria := &storage.EntityUpdateCriteria{
		Key:       entityKey,
		Type:      entityType,
		NewConfig: &wrappers.BytesValue{Value: serializedConfig},
	}
	_, err = UpdateEntities(networkID, []*storage.EntityUpdateCriteria{updateCriteria})
	return err
}

func DeleteEntityConfig(networkID, entityType, entityKey string) error {
	updateCriteria := &storage.EntityUpdateCriteria{
		Key:       entityKey,
		Type:      entityType,
		NewConfig: &wrappers.BytesValue{Value: []byte("")},
	}
	_, err := UpdateEntities(networkID, []*storage.EntityUpdateCriteria{updateCriteria})
	return err
}

// DeleteEntity deletes the entity specified by networkID, type, key
func DeleteEntities(networkID string, ids []*storage.EntityID) error {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteEntities(
		context.Background(),
		&protos.DeleteEntitiesRequest{
			NetworkID: networkID,
			ID:        ids,
		},
	)
	return err
}

// DeleteInternalEntity is a loose wrapper around DeleteEntities to delete an
// entity in the internal network structure
func DeleteInternalEntities(ids []*storage.EntityID) error {
	return DeleteEntities(storage.InternalNetworkID, ids)
}

// GetPhysicalIDOfEntity gets the physicalID associated with the entity
// identified by (networkID, entityType, entityKey)
func GetPhysicalIDOfEntity(networkID, entityType, entityKey string) (string, error) {
	entities, _, err := LoadEntities(
		networkID,
		nil,
		nil,
		[]*storage.EntityID{
			{Type: entityType, Key: entityKey},
		},
		&storage.EntityLoadCriteria{LoadMetadata: true},
	)
	if err != nil || len(entities) != 1 {
		return "", err
	}
	return entities[0].PhysicalID, nil
}

// LoadEntities loads entities specified by the parameters.
func LoadEntities(
	networkID string,
	typeFilter *string,
	keyFilter *string,
	IDs []*storage.EntityID,
	criteria *storage.EntityLoadCriteria,
) ([]*storage.NetworkEntity, []*storage.EntityID, error) {
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
				IDs:        IDs,
			},
			Criteria: criteria,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return resp.Entities, resp.EntitiesNotFound, err
}

// DoesEntityExist returns a boolean that indicated whether the entity specified
// exists in the network
func DoesEntityExist(networkID, entityType, entityKey string) (bool, error) {
	found, _, err := LoadEntities(
		networkID,
		nil,
		nil,
		[]*storage.EntityID{{Type: entityType, Key: entityKey}},
		&storage.EntityLoadCriteria{},
	)
	if err != nil {
		return false, err
	}
	if len(found) != 1 {
		return false, nil
	}
	return true, nil
}

// LoadAllEntitiesInNetwork fetches all entities of specified type in a network
func LoadAllEntitiesInNetwork(networkID string, entityType string, criteria *storage.EntityLoadCriteria) ([]*storage.NetworkEntity, error) {
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
			Criteria: criteria,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Entities, err
}
