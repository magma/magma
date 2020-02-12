/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/servicers"
	"magma/lte/cloud/go/services/meteringd_records/storage"
	"magma/orc8r/cloud/go/test_utils"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	testNetworkId    = "network"
	testGwHwId1      = "gw1"
	testGwHwId2      = "gw2"
	testGwLogicalId1 = "logical gw1"
	testGwLogicalId2 = "logical gw2"
	testSubId1       = "subscriber1"
	testSubId2       = "subscriber2"
)

func createTestMeteringdRecordsServerController(t *testing.T) *servicers.MeteringdRecordsServer {
	var mockStore = test_utils.NewMockDatastore()
	return servicers.NewMeteringdRecordsServer(storage.GetDatastoreBackedMeteringStorage(mockStore))
}

func TestMeteringdRecords(t *testing.T) {
	srv := createTestMeteringdRecordsServerController(t)

	//
	// Gateway  side write tests
	//

	// Create an identity and context for sending requests as gatway
	id := orcprotos.Identity{}
	idgw := orcprotos.Identity_Gateway{HardwareId: testGwHwId1, NetworkId: testNetworkId, LogicalId: testGwLogicalId1}
	id.SetGateway(&idgw)
	ctx := id.NewContextWithIdentity(context.Background())

	// Create a table
	table := protos.FlowTable{}

	// Update with no records
	_, err := srv.UpdateFlows(ctx, &table)
	assert.NoError(t, err)

	// Add a record to the table
	recordId := protos.FlowRecord_ID{Id: "record1"}
	record := protos.FlowRecord{Id: &recordId, Sid: testSubId1}
	table.Flows = append(table.Flows, &record)

	_, err = srv.UpdateFlows(ctx, &table)
	assert.NoError(t, err)

	// Add another record for same subscriber to the table
	recordId2 := protos.FlowRecord_ID{Id: "record2"}
	record2 := protos.FlowRecord{Id: &recordId2, Sid: testSubId1}
	table.Flows = append(table.Flows, &record2)
	_, err = srv.UpdateFlows(ctx, &table)
	assert.NoError(t, err)

	// Add another record for different subscriber to the table
	recordId3 := protos.FlowRecord_ID{Id: "record3"}
	record3 := protos.FlowRecord{Id: &recordId3, Sid: testSubId2}
	table.Flows = append(table.Flows, &record3)
	_, err = srv.UpdateFlows(ctx, &table)
	assert.NoError(t, err)

	// Update existing flows
	record = protos.FlowRecord{Id: &recordId, Sid: testSubId1, BytesTx: 1}
	record2 = protos.FlowRecord{Id: &recordId2, Sid: testSubId1, BytesTx: 1}
	record3 = protos.FlowRecord{Id: &recordId3, Sid: testSubId2, BytesTx: 1}
	table.Flows = []*protos.FlowRecord{&record, &record2, &record3}
	_, err = srv.UpdateFlows(ctx, &table)

	// Create another gateway identity and context
	id2 := orcprotos.Identity{}
	idgw2 := orcprotos.Identity_Gateway{HardwareId: testGwHwId2, NetworkId: testNetworkId, LogicalId: testGwLogicalId2}
	id2.SetGateway(&idgw2)
	ctx2 := id2.NewContextWithIdentity(context.Background())

	// Add another record for the second gateway
	recordId4 := protos.FlowRecord_ID{Id: "record4"}
	record4 := protos.FlowRecord{Id: &recordId4, Sid: testSubId1}
	table = protos.FlowTable{}
	table.Flows = append(table.Flows, &record4)
	_, err = srv.UpdateFlows(ctx2, &table)
	assert.NoError(t, err)

	//
	// Cloud side read tests
	//

	subQueryParam1 := protos.FlowRecordQuery_SubscriberId{SubscriberId: testSubId1}
	subQuery1 := protos.FlowRecordQuery{NetworkId: testNetworkId, Query: &subQueryParam1}
	records, err := srv.ListSubscriberRecords(ctx, &subQuery1)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(records.GetFlows()))

	subQueryParam2 := protos.FlowRecordQuery_SubscriberId{SubscriberId: testSubId2}
	subQuery2 := protos.FlowRecordQuery{NetworkId: testNetworkId, Query: &subQueryParam2}
	records, err = srv.ListSubscriberRecords(ctx, &subQuery2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records.GetFlows()))
}
