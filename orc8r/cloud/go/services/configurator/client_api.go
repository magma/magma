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
	"magma/orc8r/cloud/go/services/configurator/protos"

	"github.com/golang/glog"
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

// CreateNetworks registers the given list of Networks and returns the created networks
func CreateNetworks(networks []*protos.Network) ([]*protos.Network, error) {
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
func UpdateNetworks(updates []*protos.NetworkUpdateCriteria) error {
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

// LoadNetworks loads networks specified by networks according to criteria specified and
// returns the result
func LoadNetworks(networks []string, loadMetadata bool, loadConfigs bool) (map[string]*protos.Network, []string, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, nil, err
	}
	request := &protos.LoadNetworksRequest{
		Networks: networks,
		Criteria: &protos.NetworkLoadCriteria{
			LoadMetadata: loadMetadata,
			LoadConfigs:  loadConfigs,
		},
	}
	result, err := client.LoadNetworks(context.Background(), request)
	if err != nil {
		return nil, nil, err
	}
	return result.Networks, result.NotFound, nil
}

func GetNetworkConfigsByType(networkID string, configType string) ([]byte, error) {
	networks, _, err := LoadNetworks([]string{networkID}, false, true)
	if err != nil {
		return nil, err
	}
	if len(networks) == 0 {
		return nil, fmt.Errorf("Network %s not found", networkID)
	}
	return networks[networkID].Configs[configType], nil
}

// CreateEntities registers the given entities and returns the created network entities
func CreateEntities(networkID string, entities []*protos.NetworkEntity) ([]*protos.NetworkEntity, error) {
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

// UpdateEntities updates the registered entities and returns the updated entities
func UpdateEntities(networkID string, updates []*protos.EntityUpdateCriteria) (map[string]*protos.NetworkEntity, error) {
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

// DeleteEntity deletes the entity specified by networkID, type, key
func DeleteEntities(networkID string, ids []*protos.EntityID) error {
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

// LoadEntities loads entities specified by the parameters.
func LoadEntities(networkID string, typeFilter *string, keyFilter *string, ids []*protos.EntityID,
	criteria *protos.EntityLoadCriteria) ([]*protos.NetworkEntity, []*protos.EntityID, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.LoadEntities(
		context.Background(),
		&protos.LoadEntitiesRequest{
			NetworkID:  networkID,
			TypeFilter: protos.GetStringWrapper(typeFilter),
			KeyFilter:  protos.GetStringWrapper(keyFilter),
			EntityIDs:  ids,
			Criteria:   criteria,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return resp.Entities, resp.NotFound, err
}

func LoadAllEntitiesInNetwork(networkID string, typeVal string, criteria *protos.EntityLoadCriteria) ([]*protos.NetworkEntity, error) {
	client, err := getNBConfiguratorClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.LoadEntities(
		context.Background(),
		&protos.LoadEntitiesRequest{
			NetworkID:  networkID,
			TypeFilter: protos.GetStringWrapper(&typeVal),
			KeyFilter:  nil,
			EntityIDs:  nil,
			Criteria:   criteria,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Entities, err
}
