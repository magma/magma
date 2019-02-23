/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/protos"
	storage_protos "magma/orc8r/cloud/go/services/config/streaming/storage/protos"
)

type datastoreMconfigStorage struct {
	db datastore.Api
}

func NewDatastoreMconfigStorage(db datastore.Api) MconfigStorage {
	return &datastoreMconfigStorage{db: db}
}

const tableName = "mconfig_views"

func GetMconfigViewTableName(networkId string) string {
	return datastore.GetTableName(networkId, tableName)
}

func (store *datastoreMconfigStorage) GetMconfig(networkId string, gatewayId string) (*StoredMconfig, error) {
	marshaledRecord, _, err := store.db.Get(GetMconfigViewTableName(networkId), gatewayId)
	if err == datastore.ErrNotFound {
		return nil, nil
	}
	return getStorageTypeFromMarshaledProto(networkId, gatewayId, marshaledRecord)
}

func (store *datastoreMconfigStorage) GetMconfigs(networkId string, gatewayIds []string) (map[string]*StoredMconfig, error) {
	marshaledValues, err := store.db.GetMany(GetMconfigViewTableName(networkId), gatewayIds)
	if err != nil {
		return map[string]*StoredMconfig{}, err
	}

	ret := map[string]*StoredMconfig{}
	for gatewayId, marshaledRecord := range marshaledValues {
		storageVal, err := getStorageTypeFromMarshaledProto(networkId, gatewayId, marshaledRecord.Value)
		if err != nil {
			return map[string]*StoredMconfig{}, err
		}
		ret[gatewayId] = storageVal
	}
	return ret, nil
}

func (store *datastoreMconfigStorage) CreateOrUpdateMconfigs(networkId string, updates []*MconfigUpdateCriteria) error {
	putManyUpdates, err := getPutManyInputMapForUpdates(updates)
	if err != nil {
		return err
	}
	_, err = store.db.PutMany(GetMconfigViewTableName(networkId), putManyUpdates)
	return err
}

func (store *datastoreMconfigStorage) DeleteMconfigs(networkId string, gatewayIds []string) error {
	_, err := store.db.DeleteMany(GetMconfigViewTableName(networkId), gatewayIds)
	return err
}

func getStorageTypeFromMarshaledProto(
	networkId string,
	gatewayId string,
	marshaledStorageProto []byte,
) (*StoredMconfig, error) {
	storedMconfigProto := &storage_protos.StoredMconfig{}
	err := protos.Unmarshal(marshaledStorageProto, storedMconfigProto)
	if err != nil {
		return nil, err
	}
	return &StoredMconfig{
		NetworkId: networkId,
		GatewayId: gatewayId,
		Mconfig:   storedMconfigProto.Configs,
		Offset:    storedMconfigProto.Offset,
	}, nil
}

func getPutManyInputMapForUpdates(updates []*MconfigUpdateCriteria) (map[string][]byte, error) {
	ret := map[string][]byte{}
	for _, update := range updates {
		storageProto := &storage_protos.StoredMconfig{Configs: update.NewMconfig, Offset: update.Offset}
		marshaledProto, err := protos.MarshalIntern(storageProto)
		if err != nil {
			return map[string][]byte{}, err
		}
		ret[update.GatewayId] = marshaledProto
	}
	return ret, nil
}
