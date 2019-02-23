/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package providers

import (
	"fmt"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/streamer/providers"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

// GetProvider returns the StreamProvider for materialized mconfig views
func GetProvider() providers.StreamProvider {
	db, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		glog.Fatalf("Could not initialize DB connection: %s", err)
	}
	store := storage.NewDatastoreMconfigStorage(db)
	return NewMconfigStreamProvider(store)
}

type MconfigStreamProvider struct {
	store storage.MconfigStorage
}

func NewMconfigStreamProvider(store storage.MconfigStorage) *MconfigStreamProvider {
	return &MconfigStreamProvider{store: store}
}

func (*MconfigStreamProvider) GetStreamName() string {
	return "mconfig_views"
}

func (msp *MconfigStreamProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	if extraArgs == nil {
		return nil, fmt.Errorf("Offset must be specified in stream request")
	}

	offsetProto := &protos.MconfigStreamRequest{}
	err := ptypes.UnmarshalAny(extraArgs, offsetProto)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal offset extraArgs: %s", err)
	}

	networkId, err := magmad.FindGatewayNetworkId(gatewayId)
	if err != nil {
		return nil, err
	}
	logicalId, err := magmad.FindGatewayId(networkId, gatewayId)
	if err != nil {
		return nil, err
	}

	storedMconfig, err := msp.store.GetMconfig(networkId, logicalId)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve mconfig: %s", err)
	}

	// Short circuit if no mconfig computed or if the mconfig has not been
	// recomputed since the last time it was requested.
	if storedMconfig == nil || storedMconfig.Offset <= offsetProto.Offset {
		return []*protos.DataUpdate{}, nil
	}

	response := &protos.OffsetGatewayConfigs{
		Configs: storedMconfig.Mconfig,
		Offset:  storedMconfig.Offset,
	}
	marshaledResponse, err := protos.MarshalIntern(response)
	if err != nil {
		return nil, err
	}
	update := &protos.DataUpdate{Key: logicalId, Value: marshaledResponse}
	return []*protos.DataUpdate{update}, nil
}
