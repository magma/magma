/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package meteringd_records_test

import (
	"sort"
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records"
	meteringd_test_init "magma/lte/cloud/go/services/meteringd_records/test_init"
	orcprotos "magma/orc8r/cloud/go/protos"
	mw_tests "magma/orc8r/cloud/go/service/middleware/unary/interceptors/tests"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	testAgHwId1 = "Test-AGW-Hw-Id1"
	testAgHwId2 = "Test-AGW-Hw-Id2"
	testSubId1  = "sub1"
	testSubId2  = "sub2"
)

// Update flows in tables as a gateway would
// NOTE: This endpoint exists for testing ONLY
// Real clients will use gRPC directly
func UpdateFlowsTest(csn string, tbl *protos.FlowTable) error {
	client, conn, err := meteringd_records.GetMeteringdRecordsClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Hack in the identity context
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", csn))
	_, err = client.UpdateFlows(ctx, tbl)
	return err
}

func TestMeteringdRecordsControllerClientMethods(t *testing.T) {
	magmad_test_init.StartTestService(t)
	meteringd_test_init.StartTestService(t)
	csns := mw_tests.StartMockGwAccessControl(t, []string{testAgHwId1, testAgHwId2})

	//
	// Build fake network
	//
	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"},
		"meteringd_records_test_network")
	if err != nil {
		t.Fatalf("Magmad Register Network '%s' Error: %s", testNetworkId, err)
	}

	t.Logf("New Registered Network: %s", testNetworkId)

	hwId1 := orcprotos.AccessGatewayID{Id: testAgHwId1}
	logicalId1, err := magmad.RegisterGateway(testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &hwId1, Name: "Test GW Name"})

	if err != nil || logicalId1 == "" {
		t.Fatalf(
			"Magmad Register Error: %s, logical ID: %#v",
			err, logicalId1)
	}
	hwId2 := orcprotos.AccessGatewayID{Id: testAgHwId2}
	logicalId2, err := magmad.RegisterGateway(testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &hwId2, Name: "Test GW Name"})

	if err != nil || logicalId2 == "" {
		t.Fatalf(
			"Magmad Register Error: %s, logical ID: %#v",
			err, logicalId1)
	}

	// Create fake gateway context ids
	id1 := &orcprotos.Identity{}
	idgw1 := orcprotos.Identity_Gateway{HardwareId: hwId1.Id, NetworkId: testNetworkId, LogicalId: logicalId1}
	id1.SetGateway(&idgw1)

	id2 := &orcprotos.Identity{}
	idgw2 := orcprotos.Identity_Gateway{HardwareId: hwId2.Id, NetworkId: testNetworkId, LogicalId: logicalId2}
	id1.SetGateway(&idgw1)
	id2.SetGateway(&idgw2)

	//
	// Generate some fake flows
	//

	// Ensure there are no flows to start
	_, err = meteringd_records.GetRecord(testNetworkId, "doesn't exist")
	assert.Error(t, err)
	actualRecordSet, err := meteringd_records.ListSubscriberRecords(testNetworkId, testSubId1)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actualRecordSet))

	// Create two flows from two subs on gateway 1
	recordId1 := &protos.FlowRecord_ID{Id: "test1"}
	record1 := &protos.FlowRecord{Id: recordId1, Sid: testSubId1}
	recordId2 := &protos.FlowRecord_ID{Id: "test2"}
	record2 := &protos.FlowRecord{Id: recordId2, Sid: testSubId2}
	tbl1 := &protos.FlowTable{}
	tbl1.Flows = append(tbl1.Flows, record1)
	tbl1.Flows = append(tbl1.Flows, record2)
	err = UpdateFlowsTest(csns[0], tbl1)
	assert.NoError(t, err)

	// Create one flow for subscriber 2 on gateway 2
	recordId3 := &protos.FlowRecord_ID{Id: "test3"}
	record3 := &protos.FlowRecord{Id: recordId3, Sid: testSubId2}
	tbl2 := &protos.FlowTable{}
	tbl2.Flows = append(tbl2.Flows, record3)
	err = UpdateFlowsTest(csns[1], tbl2)
	assert.NoError(t, err)

	//
	// Read back the flows
	//

	// One for the first subscriber
	actualRecordSet, err = meteringd_records.ListSubscriberRecords(testNetworkId, testSubId1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actualRecordSet))
	// Fill in gateway ID for expected records
	record1.GatewayId = logicalId1
	assert.Equal(t, orcprotos.TestMarshal(record1), orcprotos.TestMarshal(actualRecordSet[0]))

	// Two for the second subscriber
	actualRecordSet, err = meteringd_records.ListSubscriberRecords(testNetworkId, testSubId2)
	assert.NoError(t, err)
	sort.Slice(actualRecordSet, func(i, j int) bool { return actualRecordSet[i].GetId().GetId() < actualRecordSet[j].GetId().GetId() })

	assert.Equal(t, 2, len(actualRecordSet))
	// Fill in gateway ID for expected records
	record2.GatewayId = logicalId1
	record3.GatewayId = logicalId2
	assert.Equal(t, orcprotos.TestMarshal(record2), orcprotos.TestMarshal(actualRecordSet[0]))
	assert.Equal(t, orcprotos.TestMarshal(record3), orcprotos.TestMarshal(actualRecordSet[1]))
}
