/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"fmt"

	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/datastore"

	"github.com/golang/protobuf/proto"
)

const (
	FLOWS_TABLE          = "flows"
	SUBSCRIBER_FLOWS_IDX = "subscriber_to_flow_id"
)

func GetFlowsTableName(networkId string) string {
	return datastore.GetTableName(networkId, FLOWS_TABLE)
}

func GetSubscriberIndexTableName(networkId string) string {
	return datastore.GetTableName(networkId, SUBSCRIBER_FLOWS_IDX)
}

type datastoreStorage struct {
	db datastore.Api
}

func GetDatastoreBackedMeteringStorage(db datastore.Api) MeteringRecordsStorage {
	return &datastoreStorage{db: db}
}

func (*datastoreStorage) InitTables() error {
	// No need for table initialization using datastore API
	return nil
}

func (store *datastoreStorage) UpdateOrCreateRecords(networkId string, flows []*protos.FlowRecord) error {
	marshaledRecordsById, err := getMarshaledFlowsById(flows)
	if err != nil {
		return err
	}

	recordTbl := datastore.GetTableName(networkId, FLOWS_TABLE)
	if _, err := store.db.PutMany(recordTbl, marshaledRecordsById); err != nil {
		return fmt.Errorf("Error writing new flows: %s", err)
	}
	if err := store.updateSubscriberFlowsIndex(networkId, flows); err != nil {
		return fmt.Errorf("Error updating subscriber flows index: %s", err)
	}
	return nil
}

func (store *datastoreStorage) GetRecord(networkId string, recordId string) (*protos.FlowRecord, error) {
	recordTbl := datastore.GetTableName(networkId, FLOWS_TABLE)
	marshaledRecord, _, err := store.db.Get(recordTbl, recordId)
	if err != nil {
		return nil, fmt.Errorf("Error getting flow record: %s", err)
	}
	ret := &protos.FlowRecord{}
	if err := proto.Unmarshal(marshaledRecord, ret); err != nil {
		return nil, fmt.Errorf("Error unmarshalling flow record: %s", err)
	}
	return ret, nil
}

func (store *datastoreStorage) GetRecordsForSubscriber(networkId string, sid string) ([]*protos.FlowRecord, error) {
	exists, err := store.db.DoesKeyExist(GetSubscriberIndexTableName(networkId), sid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return []*protos.FlowRecord{}, nil
	}

	recordSet, err := store.getRecordSet(GetSubscriberIndexTableName(networkId), networkId, sid)
	if err != nil {
		return nil, err
	}

	// Cross-reference the flow record table
	marshaledFlowsById, err := store.db.GetMany(GetFlowsTableName(networkId), recordSet.GetRecordIds())
	if err != nil {
		return nil, err
	}
	return unmarshalFlowsById(marshaledFlowsById)
}

func (store *datastoreStorage) DeleteRecordsForSubscriber(networkId string, sid string) error {
	subIdxTblName := GetSubscriberIndexTableName(networkId)
	exists, err := store.db.DoesKeyExist(subIdxTblName, sid)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	// Use the sid->[flow_id] mapping to get all flow_id for this sid
	recordSet, err := store.getRecordSet(subIdxTblName, networkId, sid)
	if err != nil {
		return err
	}
	// Delete entry of sid in the index table
	err = store.db.Delete(subIdxTblName, sid)
	if err != nil {
		return err
	}
	// Delete all flow_id in the main table
	_, err = store.db.DeleteMany(GetFlowsTableName(networkId), recordSet.GetRecordIds())

	return err
}

// ==========
// Helpers
// ==========

func getMarshaledFlowsById(flows []*protos.FlowRecord) (map[string][]byte, error) {
	ret := make(map[string][]byte, len(flows))
	for _, record := range flows {
		value, err := proto.Marshal(record)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling flow record: %s", err)
		}
		ret[record.GetId().GetId()] = value
	}
	return ret, nil
}

func (store *datastoreStorage) updateSubscriberFlowsIndex(networkId string, flows []*protos.FlowRecord) error {
	// Partition new flows by SID and fetch existing record sets for all SIDs
	flowsBySid := getFlowsBySid(flows)
	existingRecordSetsBySid, err := store.getExistingRecordSetsForSubscribers(networkId, flowsBySid)
	if err != nil {
		return fmt.Errorf("Failed to update subscriber flows index: %s", err)
	}

	// Union records and write to DB
	recordSetsToPersist := getUnionedRecordSetsBySid(flowsBySid, existingRecordSetsBySid)
	marshaledRecordSets, err := marshalRecordSets(recordSetsToPersist)
	if err != nil {
		return fmt.Errorf("Failed to update subscriber flows index: %s", err)
	}
	subIdxTblName := GetSubscriberIndexTableName(networkId)
	if _, err := store.db.PutMany(subIdxTblName, marshaledRecordSets); err != nil {
		return fmt.Errorf("Failed to write subscriber flow set: %s", err)
	}
	return nil
}

func getFlowsBySid(flows []*protos.FlowRecord) map[string][]*protos.FlowRecord {
	flowsBySid := make(map[string][]*protos.FlowRecord)
	for _, record := range flows {
		sid := record.GetSid()
		if _, exists := flowsBySid[sid]; !exists {
			flowsBySid[sid] = []*protos.FlowRecord{record}
		} else {
			flowsBySid[sid] = append(flowsBySid[sid], record)
		}
	}
	return flowsBySid
}

func (store *datastoreStorage) getExistingRecordSetsForSubscribers(
	networkId string,
	flowsBySid map[string][]*protos.FlowRecord,
) (map[string]*protos.FlowRecordSet, error) {
	sids := getAllSids(flowsBySid)
	subscriberRecordTbl := GetSubscriberIndexTableName(networkId)
	subscriberTblFlows, err := store.db.GetMany(subscriberRecordTbl, sids)
	if err != nil {
		return nil, fmt.Errorf("Error getting subscriber flow set: %s", err)
	}

	ret := make(map[string]*protos.FlowRecordSet, len(sids))
	for sid, marshaledRecordSet := range subscriberTblFlows {
		recordSet := &protos.FlowRecordSet{}
		if err := proto.Unmarshal(marshaledRecordSet.Value, recordSet); err != nil {
			return nil, fmt.Errorf("failed to unmarshal subscriber flow set: %s", err)
		}
		ret[sid] = recordSet
	}
	return ret, nil
}

func getAllSids(flowsBySid map[string][]*protos.FlowRecord) []string {
	ret := make([]string, 0, len(flowsBySid))
	for sid := range flowsBySid {
		ret = append(ret, sid)
	}
	return ret
}

func getUnionedRecordSetsBySid(
	flowsBySid map[string][]*protos.FlowRecord,
	recordSetsBySid map[string]*protos.FlowRecordSet,
) map[string]*protos.FlowRecordSet {
	emptyRecordSet := &protos.FlowRecordSet{RecordIds: []string{}}
	ret := make(map[string]*protos.FlowRecordSet, len(flowsBySid))

	for sid, newFlows := range flowsBySid {
		existingRecordSet, exists := recordSetsBySid[sid]
		if exists {
			ret[sid] = unionFlowTableWithRecordSet(newFlows, existingRecordSet)
		} else {
			ret[sid] = unionFlowTableWithRecordSet(newFlows, emptyRecordSet)
		}
	}

	return ret
}

func marshalRecordSets(recordSetsBySid map[string]*protos.FlowRecordSet) (map[string][]byte, error) {
	ret := make(map[string][]byte, len(recordSetsBySid))
	for sid, recordSet := range recordSetsBySid {
		marshaledRecordSet, err := proto.Marshal(recordSet)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal subscriber record set: %s", err)
		}
		ret[sid] = marshaledRecordSet
	}
	return ret, nil
}

func unionFlowTableWithRecordSet(flows []*protos.FlowRecord, recordSet *protos.FlowRecordSet) *protos.FlowRecordSet {
	ret := &protos.FlowRecordSet{RecordIds: make([]string, 0, len(flows)+len(recordSet.GetRecordIds()))}
	existingRecordIds := map[string]struct{}{}
	for _, recordId := range recordSet.GetRecordIds() {
		existingRecordIds[recordId] = struct{}{}
		ret.RecordIds = append(ret.GetRecordIds(), recordId)
	}
	for _, newFlow := range flows {
		flowId := newFlow.GetId().GetId()
		if _, exists := existingRecordIds[flowId]; !exists {
			existingRecordIds[flowId] = struct{}{}
			ret.RecordIds = append(ret.GetRecordIds(), flowId)
		}
	}
	return ret
}

func (store *datastoreStorage) getRecordSet(tableName string, networkId string, key string) (*protos.FlowRecordSet, error) {
	marshaledRecordSet, _, err := store.db.Get(tableName, key)
	if err != nil {
		return nil, err
	}

	ret := &protos.FlowRecordSet{}
	if err = proto.Unmarshal(marshaledRecordSet, ret); err != nil {
		return nil, fmt.Errorf("Could not unmarshal record set: %s", err)
	}
	return ret, nil
}

func unmarshalFlowsById(marshaledFlowsById map[string]datastore.ValueWrapper) ([]*protos.FlowRecord, error) {
	ret := make([]*protos.FlowRecord, 0, len(marshaledFlowsById))
	for _, marshaledFlow := range marshaledFlowsById {
		flowRecord := &protos.FlowRecord{}
		err := proto.Unmarshal(marshaledFlow.Value, flowRecord)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling flow record: %s", err)
		}
		ret = append(ret, flowRecord)
	}
	return ret, nil
}
