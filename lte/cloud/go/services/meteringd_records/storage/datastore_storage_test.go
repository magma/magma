/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"errors"
	"sort"
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/storage"
	"magma/orc8r/cloud/go/test_utils"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestDatastoreStorage_GetRecord(t *testing.T) {
	ds := test_utils.NewMockDatastore()
	store := storage.GetDatastoreBackedMeteringStorage(ds)
	networkId := "network"

	// Fixtures
	flow1 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow1"},
		Sid:       "sid1",
		GatewayId: "gw1",
	}
	flow2 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow2"},
		Sid:       "sid2",
		GatewayId: "gw2",
	}
	setupTestFixtures(t, ds, networkId, []*protos.FlowRecord{flow1, flow2})

	actual, err := store.GetRecord(networkId, "flow1")
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(flow1), orcprotos.TestMarshal(actual))

	actual, err = store.GetRecord(networkId, "flow2")
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(flow2), orcprotos.TestMarshal(actual))

	_, err = store.GetRecord(networkId, "flow3")
	assert.Error(t, err)
}

func TestDatastoreStorage_UpdateOrCreateRecords(t *testing.T) {
	ds := test_utils.NewMockDatastore()
	store := storage.GetDatastoreBackedMeteringStorage(ds)
	networkId := "network"

	// Fixtures
	flow1 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow1"},
		Sid:       "sid1",
		GatewayId: "gw1",
	}
	flow2 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow2"},
		Sid:       "sid2",
		GatewayId: "gw2",
	}
	setupTestFixtures(t, ds, networkId, []*protos.FlowRecord{flow1, flow2})

	// Add a new flow
	flow3 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow3"},
		Sid:       "sid3",
		GatewayId: "gw3",
	}
	err := store.UpdateOrCreateRecords(networkId, []*protos.FlowRecord{flow3})
	assert.NoError(t, err)
	assertDatastoreWritesSucceeded(t, ds, networkId, []*protos.FlowRecord{flow1, flow2, flow3})

	// Update an existing flow
	flow2.GatewayId = "gw4"
	flow2.BytesTx = 1
	err = store.UpdateOrCreateRecords(networkId, []*protos.FlowRecord{flow2})
	assert.NoError(t, err)
	assertDatastoreWritesSucceeded(t, ds, networkId, []*protos.FlowRecord{flow1, flow2, flow3})

	// Update and add flows
	flow3.GatewayId = "gw5"
	flow3.BytesTx = 2
	flow4 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow4"},
		Sid:       "sid4",
		GatewayId: "gw6",
	}
	err = store.UpdateOrCreateRecords(networkId, []*protos.FlowRecord{flow3, flow4})
	assert.NoError(t, err)
	assertDatastoreWritesSucceeded(t, ds, networkId, []*protos.FlowRecord{flow1, flow2, flow3, flow4})

	// Add flows from scratch
	err = store.UpdateOrCreateRecords("network2", []*protos.FlowRecord{flow1})
	assert.NoError(t, err)
	assertDatastoreWritesSucceeded(t, ds, "network2", []*protos.FlowRecord{flow1})
}

func TestDatastoreStorage_GetRecordsForSubscriber(t *testing.T) {
	ds := test_utils.NewMockDatastore()
	store := storage.GetDatastoreBackedMeteringStorage(ds)
	networkId := "network"

	// Fixtures
	flow1 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow1"},
		Sid:       "sid1",
		GatewayId: "gw1",
	}
	flow2 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow2"},
		Sid:       "sid2",
		GatewayId: "gw2",
	}
	flow3 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow3"},
		Sid:       "sid1",
		GatewayId: "gw3",
	}
	setupTestFixtures(t, ds, networkId, []*protos.FlowRecord{flow1, flow2, flow3})

	actual, err := store.GetRecordsForSubscriber(networkId, "sid1")
	assert.NoError(t, err)
	sort.Slice(actual, func(i, j int) bool { return actual[i].GetId().GetId() < actual[j].GetId().GetId() })

	assert.Equal(t, orcprotos.TestMarshal(flow1), orcprotos.TestMarshal(actual[0]))
	assert.Equal(t, orcprotos.TestMarshal(flow3), orcprotos.TestMarshal(actual[1]))

	actual, err = store.GetRecordsForSubscriber(networkId, "sid2")
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(flow2), orcprotos.TestMarshal(actual[0]))
}

func TestDatastoreStorage_DeleteRecordsForSubscriber(t *testing.T) {
	ds := test_utils.NewMockDatastore()
	store := storage.GetDatastoreBackedMeteringStorage(ds)
	networkId := "network"

	// Fixtures
	flow1 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow1"},
		Sid:       "sid1",
		GatewayId: "gw1",
	}
	flow2 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow2"},
		Sid:       "sid1",
		GatewayId: "gw2",
	}
	flow3 := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "flow3"},
		Sid:       "sid2",
		GatewayId: "gw3",
	}
	setupTestFixtures(t, ds, networkId, []*protos.FlowRecord{flow1, flow2, flow3})

	// Delete flows of sid1
	err := store.DeleteRecordsForSubscriber(networkId, "sid1")
	assert.NoError(t, err)

	// Check if the sid1 is deleted
	test_utils.AssertDatastoreDoesNotHaveRow(t, ds, storage.GetSubscriberIndexTableName(networkId), "sid1")

	// Check if flows of sid1 are deleted
	test_utils.AssertDatastoreDoesNotHaveRow(t, ds, storage.GetFlowsTableName(networkId), "flow1")
	test_utils.AssertDatastoreDoesNotHaveRow(t, ds, storage.GetFlowsTableName(networkId), "flow2")

	// Check if flow of sid2 still exists
	assertDatastoreWritesSucceeded(t, ds, networkId, []*protos.FlowRecord{flow3})
}

func assertDatastoreWritesSucceeded(t *testing.T, store *test_utils.MockDatastore, networkId string, flows []*protos.FlowRecord) {
	flowsById := getFlowInterfacesById(flows)
	flowSetsBySid := getExpectedSubscriberFlowSetsBySid(flows)

	test_utils.AssertDatastoreHasRows(
		t, store,
		storage.GetFlowsTableName(networkId),
		flowsById,
		deserializeFlowRecord,
	)
	test_utils.AssertDatastoreHasRows(
		t, store,
		storage.GetSubscriberIndexTableName(networkId),
		flowSetsBySid,
		deserializeFlowSet,
	)
}

// Given a list of flow record fixtures, set up the datastore including the
// subscriber flow index
func setupTestFixtures(t *testing.T, store *test_utils.MockDatastore, networkId string, flows []*protos.FlowRecord) {
	flowsById := getFlowInterfacesById(flows)
	flowSetsBySid := getExpectedSubscriberFlowSetsBySid(flows)

	test_utils.SetupTestFixtures(
		t, store,
		storage.GetFlowsTableName(networkId),
		flowsById,
		serializeFlowRecord,
	)

	test_utils.SetupTestFixtures(
		t, store,
		storage.GetSubscriberIndexTableName(networkId),
		flowSetsBySid,
		serializeFlowSet,
	)
}

func getFlowInterfacesById(flows []*protos.FlowRecord) map[string]interface{} {
	flowsById := map[string]interface{}{}
	for _, flow := range flows {
		flowsById[flow.GetId().GetId()] = flow
	}
	return flowsById
}

func getExpectedSubscriberFlowSetsBySid(flows []*protos.FlowRecord) map[string]interface{} {
	flowSetsBySid := map[string]interface{}{}
	for _, flow := range flows {
		iCurrentSet, exists := flowSetsBySid[flow.GetSid()]
		if exists {
			currentSet := iCurrentSet.(*protos.FlowRecordSet)
			currentSet.RecordIds = append(currentSet.RecordIds, flow.GetId().GetId())
		} else {
			flowSetsBySid[flow.GetSid()] = &protos.FlowRecordSet{RecordIds: []string{flow.GetId().GetId()}}
		}
	}
	return flowSetsBySid
}

func serializeFlowRecord(flow interface{}) ([]byte, error) {
	flowCasted, ok := flow.(*protos.FlowRecord)
	if !ok {
		return nil, errors.New("Expected *protos.FlowRecord")
	}
	return proto.Marshal(flowCasted)
}

func serializeFlowSet(flowSet interface{}) ([]byte, error) {
	flowSetCasted, ok := flowSet.(*protos.FlowRecordSet)
	if !ok {
		return nil, errors.New("Expected *protos.FlowRecordSet")
	}
	return proto.Marshal(flowSetCasted)
}

func deserializeFlowRecord(marshaled []byte) (interface{}, error) {
	ret := &protos.FlowRecord{}
	err := proto.Unmarshal(marshaled, ret)
	return ret, err
}

func deserializeFlowSet(marshaled []byte) (interface{}, error) {
	ret := &protos.FlowRecordSet{}
	err := proto.Unmarshal(marshaled, ret)
	return ret, err
}
