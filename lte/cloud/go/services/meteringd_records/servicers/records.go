/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
MeteringD flows servicer provides the gRPC interface for the REST and
services to interact with traffic flows records.

The servicer require a backing Datastore (which is typically Postgres)
for storing and retrieving the data and access to Magmad to resolve the network.
*/
package servicers

import (
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/storage"
	orcprotos "magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements a gRPC Server interface for Metering Records
type MeteringdRecordsServer struct {
	storage storage.MeteringRecordsStorage
}

// Convenience method to instantiate a new gRPC server instance
func NewMeteringdRecordsServer(storage storage.MeteringRecordsStorage) *MeteringdRecordsServer {
	srv := &MeteringdRecordsServer{storage: storage}
	return srv
}

// Servicer for synchronizing the flows in the datastore with the gateway
// Given a Flow Table which is a list of Flow Records, each element is put
// into the record table.
func (srv *MeteringdRecordsServer) UpdateFlows(ctx context.Context, tbl *protos.FlowTable) (*orcprotos.Void, error) {
	ret := &orcprotos.Void{}
	id, err := orcprotos.GetGatewayIdentity(ctx)
	if err != nil {
		return ret, err
	}

	fillFlowsWithGatewayId(tbl, id.GetLogicalId())
	err = srv.storage.UpdateOrCreateRecords(id.GetNetworkId(), tbl.GetFlows())
	return ret, err
}

// Lists the ids of all usage records for a subscriber on the network
func (srv *MeteringdRecordsServer) ListSubscriberRecords(
	ctx context.Context,
	query *protos.FlowRecordQuery,
) (*protos.FlowTable, error) {
	if query.GetNetworkId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Missing Network identity")
	}
	if query.GetSubscriberId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Missing subscriber id")
	}

	subscriberFlows, err := srv.storage.GetRecordsForSubscriber(query.GetNetworkId(), query.GetSubscriberId())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}
	return &protos.FlowTable{Flows: subscriberFlows}, nil
}

// Gets a flow data record by ID
func (srv *MeteringdRecordsServer) GetRecord(ctx context.Context, query *protos.FlowRecordQuery) (*protos.FlowRecord, error) {
	if query.GetNetworkId() == "" {
		return &protos.FlowRecord{}, status.Errorf(codes.InvalidArgument, "Missing Network identity")
	}
	if query.GetRecordId() == "" {
		return &protos.FlowRecord{}, status.Errorf(codes.InvalidArgument, "Missing record id for query")
	}

	return srv.storage.GetRecord(query.GetNetworkId(), query.GetRecordId())
}

func fillFlowsWithGatewayId(tbl *protos.FlowTable, gatewayId string) *protos.FlowTable {
	for _, record := range tbl.GetFlows() {
		record.GatewayId = gatewayId
	}
	return tbl
}
