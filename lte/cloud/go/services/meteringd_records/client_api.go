/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package meteringd_records provides a client API for interacting with the
// meteringd_records service, which manages flow record entities to track
// data usage by subscribers.
package meteringd_records

import (
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const ServiceName = "METERINGD_RECORDS"

// Get a thin RPC client to the stats service.
func GetMeteringdRecordsClient() (protos.MeteringdRecordsControllerClient, *grpc.ClientConn, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, nil, initErr
	}
	return protos.NewMeteringdRecordsControllerClient(conn), conn, err
}

// Get a Record from a network
func GetRecord(networkId string, recordId string) (*protos.FlowRecord, error) {
	client, conn, err := GetMeteringdRecordsClient()
	if err != nil {
		return &protos.FlowRecord{}, err
	}
	defer conn.Close()

	req := &protos.FlowRecordQuery{
		NetworkId: networkId,
		Query: &protos.FlowRecordQuery_RecordId{
			RecordId: recordId,
		},
	}
	ctx := context.Background()
	return client.GetRecord(ctx, req)
}

// List Records for a subscriber
func ListSubscriberRecords(networkId string, sid string) ([]*protos.FlowRecord, error) {
	client, conn, err := GetMeteringdRecordsClient()
	if err != nil {
		return []*protos.FlowRecord{}, err
	}
	defer conn.Close()

	req := &protos.FlowRecordQuery{
		NetworkId: networkId,
		Query: &protos.FlowRecordQuery_SubscriberId{
			SubscriberId: sid,
		},
	}

	ctx := context.Background()
	res, err := client.ListSubscriberRecords(ctx, req)
	if err != nil {
		return []*protos.FlowRecord{}, err
	}
	return res.GetFlows(), nil
}
