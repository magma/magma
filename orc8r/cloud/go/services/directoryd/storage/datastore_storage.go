/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"fmt"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/lib/go/protos"
)

type DirectorydPersistenceServiceImpl struct {
	db datastore.Api
}

func GetDirectorydPersistenceService(db datastore.Api) DirectorydPersistenceService {
	return &DirectorydPersistenceServiceImpl{db: db}
}

func (store *DirectorydPersistenceServiceImpl) GetRecord(tableId protos.TableID, recordId string) (*protos.LocationRecord, error) {
	recordTbl := tableId.String()
	marshaledRecord, _, err := store.db.Get(recordTbl, recordId)
	if err != nil {
		return nil, fmt.Errorf("Error getting location record: %s", err)
	}

	ret := &protos.LocationRecord{}
	if err := protos.Unmarshal(marshaledRecord, ret); err != nil {
		return nil, fmt.Errorf("Error unmarshalling location record: %s", err)
	}
	return ret, nil
}

func (store *DirectorydPersistenceServiceImpl) UpdateOrCreateRecord(tableId protos.TableID, recordId string, record *protos.LocationRecord) error {
	recordTbl := tableId.String()

	value, err := protos.MarshalIntern(record)
	if err != nil {
		return fmt.Errorf("Error marshaling location record: %s", err)
	}

	if err := store.db.Put(recordTbl, recordId, value); err != nil {
		return fmt.Errorf("Error updating new location record: %s", err)
	}
	return nil
}

func (store *DirectorydPersistenceServiceImpl) DeleteRecord(tableId protos.TableID, recordId string) error {
	recordTbl := tableId.String()

	_, _, err := store.db.Get(recordTbl, recordId)
	if err != nil {
		return fmt.Errorf("Error finding location record: %s", err)
	}

	err = store.db.Delete(recordTbl, recordId)
	if err != nil {
		return fmt.Errorf("Error deleting location record: %s", err)
	}
	return err
}
