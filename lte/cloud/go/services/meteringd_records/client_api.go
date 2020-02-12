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
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const ServiceName = "METERINGD_RECORDS"

// GetMeteringdRecordsClient get a thin RPC client to the stats service.
func GetMeteringdRecordsClient() (protos.MeteringdRecordsControllerClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewMeteringdRecordsControllerClient(conn), err
}

// GetRecord get a Record from a network
func GetRecord(networkID string, recordID string) (*protos.FlowRecord, error) {
	client, err := GetMeteringdRecordsClient()
	if err != nil {
		return &protos.FlowRecord{}, err
	}

	req := &protos.FlowRecordQuery{
		NetworkId: networkID,
		Query: &protos.FlowRecordQuery_RecordId{
			RecordId: recordID,
		},
	}
	ctx := context.Background()
	return client.GetRecord(ctx, req)
}

// ListSubscriberRecords list Records for a subscriber
func ListSubscriberRecords(networkID string, sid string) ([]*protos.FlowRecord, error) {
	client, err := GetMeteringdRecordsClient()
	if err != nil {
		return []*protos.FlowRecord{}, err
	}

	req := &protos.FlowRecordQuery{
		NetworkId: networkID,
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
