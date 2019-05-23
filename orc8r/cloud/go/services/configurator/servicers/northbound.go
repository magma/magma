/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"
	"fmt"

	commonProtos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"
)

type nbConfiguratorServicer struct {
	factory storage.ConfiguratorStorageFactory
}

// NewNorthboundConfiguratorServicer returns a configurator server backed by storage passed in
func NewNorthboundConfiguratorServicer(factory storage.ConfiguratorStorageFactory) (protos.NorthboundConfiguratorServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	return &nbConfiguratorServicer{factory}, nil
}

func (srv *nbConfiguratorServicer) LoadNetworks(context context.Context, req *protos.LoadNetworksRequest) (*protos.LoadNetworksResponse, error) {
	res := &protos.LoadNetworksResponse{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: true})
	if err != nil {
		return res, err
	}

	result, err := store.LoadNetworks(req.Networks, req.Criteria.ToNetworkLoadCriteria())
	if err != nil {
		store.Rollback()
		return res, err
	}

	res.Networks = map[string]*protos.Network{}
	for _, network := range result.Networks {
		pNetwork := protos.FromStorageNetwork(network)
		res.Networks[network.ID] = pNetwork
	}
	res.NotFound = result.NetworkIDsNotFound
	return res, store.Commit()
}

func (srv *nbConfiguratorServicer) ListNetworkIDs(context context.Context, void *commonProtos.Void) (*protos.ListNetworkIDsResponse, error) {
	res := &protos.ListNetworkIDsResponse{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: true})
	if err != nil {
		return res, err
	}

	networks, err := store.LoadAllNetworks(storage.FullNetworkLoadCriteria)
	if err != nil {
		store.Rollback()
		return res, err
	}
	res.NetworkIDs = []string{}
	for _, network := range networks {
		res.NetworkIDs = append(res.NetworkIDs, network.ID)
	}
	return res, store.Commit()
}

func (srv *nbConfiguratorServicer) CreateNetworks(context context.Context, req *protos.CreateNetworksRequest) (*protos.CreateNetworksResponse, error) {
	emptyRes := &protos.CreateNetworksResponse{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: false})
	if err != nil {
		return emptyRes, err
	}

	createdNetworks := []*protos.Network{}
	for _, network := range req.Networks {
		err = networkConfigsAreValid(network.Configs)
		if err != nil {
			return emptyRes, err
		}
		createdNetwork, err := store.CreateNetwork(network.ToNetwork())
		if err != nil {
			store.Rollback()
			return emptyRes, err
		}
		createdNetworks = append(createdNetworks, protos.FromStorageNetwork(createdNetwork))
	}
	return &protos.CreateNetworksResponse{CreatedNetworks: createdNetworks}, store.Commit()
}

func (srv *nbConfiguratorServicer) UpdateNetworks(context context.Context, req *protos.UpdateNetworksRequest) (*commonProtos.Void, error) {
	void := &commonProtos.Void{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: false})
	if err != nil {
		return void, err
	}

	updates := []storage.NetworkUpdateCriteria{}
	for _, pUpdate := range req.Updates {
		err = networkConfigsAreValid(pUpdate.ConfigsToAddOrUpdate)
		if err != nil {
			return void, err
		}
		updates = append(updates, pUpdate.ToNetworkUpdateCriteria())
	}
	_, err = store.UpdateNetworks(updates)
	if err != nil {
		store.Rollback()
		return void, err
	}
	return void, store.Commit()
}

func (srv *nbConfiguratorServicer) DeleteNetworks(context context.Context, req *protos.DeleteNetworksRequest) (*commonProtos.Void, error) {
	void := &commonProtos.Void{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: false})
	if err != nil {
		return void, err
	}

	deleteRequests := []storage.NetworkUpdateCriteria{}
	for _, networkID := range req.NetworkIDs {
		deleteRequests = append(deleteRequests, storage.NetworkUpdateCriteria{ID: networkID, DeleteNetwork: true})
	}
	_, err = store.UpdateNetworks(deleteRequests)
	if err != nil {
		store.Rollback()
		return void, err
	}
	return void, store.Commit()
}

func (srv *nbConfiguratorServicer) LoadEntities(context context.Context, req *protos.LoadEntitiesRequest) (*protos.LoadEntitiesResponse, error) {
	emptyRes := &protos.LoadEntitiesResponse{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: false})
	if err != nil {
		return emptyRes, err
	}

	loadFilter := protos.ToEntityLoadFilter(req.TypeFilter, req.KeyFilter, req.EntityIDs)
	loadResult, err := store.LoadEntities(req.NetworkID, loadFilter, req.Criteria.ToEntityLoadCriteria())
	if err != nil {
		store.Rollback()
		return emptyRes, err
	}
	return &protos.LoadEntitiesResponse{
		Entities: protos.FromStorageNetworkEntities(loadResult.Entities),
		NotFound: protos.FromTKs(loadResult.EntitiesNotFound),
	}, store.Commit()
}

func (srv *nbConfiguratorServicer) CreateEntities(context context.Context, req *protos.CreateEntitiesRequest) (*protos.CreateEntitiesResponse, error) {
	emptyRes := &protos.CreateEntitiesResponse{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: false})
	if err != nil {
		return emptyRes, err
	}

	createdEntities := []*protos.NetworkEntity{}
	for _, entity := range req.Entities {
		if err := entityConfigIsValid(entity.Type, entity.Config); err != nil {
			return emptyRes, err
		}
		createdEntity, err := store.CreateEntity(req.NetworkID, entity.ToNetworkEntity())
		if err != nil {
			store.Rollback()
			return emptyRes, err
		}
		createdEntities = append(createdEntities, protos.FromStorageNetworkEntity(createdEntity))
	}
	return &protos.CreateEntitiesResponse{CreatedEntities: createdEntities}, store.Commit()
}

func (srv *nbConfiguratorServicer) UpdateEntities(context context.Context, req *protos.UpdateEntitiesRequest) (*protos.UpdateEntitiesResponse, error) {
	emptyRes := &protos.UpdateEntitiesResponse{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: false})
	if err != nil {
		return emptyRes, err
	}

	updatedEntities := map[string]*protos.NetworkEntity{}
	for _, update := range req.Updates {
		if update.NewConfig != nil {
			if err := entityConfigIsValid(update.Type, update.NewConfig.Value); err != nil {
				return emptyRes, err
			}
		}

		updatedEntity, err := store.UpdateEntity(req.NetworkID, update.ToEntityUpdateCriteria())
		if err != nil {
			store.Rollback()
			return emptyRes, err
		}
		updatedEntities[update.Key] = protos.FromStorageNetworkEntity(updatedEntity)
	}
	return &protos.UpdateEntitiesResponse{UpdatedEntities: updatedEntities}, store.Commit()
}

func (srv *nbConfiguratorServicer) DeleteEntities(context context.Context, req *protos.DeleteEntitiesRequest) (*commonProtos.Void, error) {
	void := &commonProtos.Void{}
	store, err := srv.factory.StartTransaction(context, &storage.TxOptions{ReadOnly: false})
	if err != nil {
		return void, err
	}

	for _, entityID := range req.ID {
		request := storage.EntityUpdateCriteria{
			Type:         entityID.Type,
			Key:          entityID.Id,
			DeleteEntity: true,
		}
		_, err = store.UpdateEntity(req.NetworkID, request)
		if err != nil {
			store.Rollback()
			return void, err
		}
	}
	return void, store.Commit()
}

func networkConfigsAreValid(configs map[string][]byte) error {
	for typeVal, config := range configs {
		_, err := serde.Deserialize(configurator.SerdeDomain, typeVal, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func entityConfigIsValid(typeVal string, config []byte) error {
	_, err := serde.Deserialize(configurator.SerdeDomain, typeVal, config)
	if err != nil {
		return err
	}
	return nil
}
