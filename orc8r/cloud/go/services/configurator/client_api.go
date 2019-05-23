/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package configurator

import (
	"context"
	"sync"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/configurator/protos"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

var connSingleton = (*grpc.ClientConn)(nil)
var connGuard = sync.Mutex{}

// Todo use grpc conn from registry once it's landed
func getNBConfiguratorClient() (protos.NorthboundConfiguratorClient, error) {
	if connSingleton == nil {
		connGuard.Lock()
		if connSingleton == nil {
			conn, err := registry.GetConnection(ServiceName)
			if err != nil {
				initErr := errors.NewInitError(err, ServiceName)
				glog.Error(initErr)
				connGuard.Unlock()
				return nil, initErr
			}
			connSingleton = conn
		}
		connGuard.Unlock()
	}
	return protos.NewNorthboundConfiguratorClient(connSingleton), nil
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
