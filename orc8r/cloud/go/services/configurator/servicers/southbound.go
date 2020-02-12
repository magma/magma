/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"
	orc8rStorage "magma/orc8r/cloud/go/storage"
	commonProtos "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type sbConfiguratorServicer struct {
	factory storage.ConfiguratorStorageFactory
}

func NewSouthboundConfiguratorServicer(factory storage.ConfiguratorStorageFactory) (protos.SouthboundConfiguratorServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	return &sbConfiguratorServicer{factory}, nil
}

func (srv *sbConfiguratorServicer) GetMconfig(ctx context.Context, void *commonProtos.Void) (*commonProtos.GatewayConfigs, error) {
	gw := commonProtos.GetClientGateway(ctx)
	if gw == nil {
		return nil, status.Errorf(codes.PermissionDenied, "missing gateway identity")
	}
	if !gw.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "gateway not registered")
	}
	return srv.getMconfigImpl(gw.NetworkId, gw.LogicalId)
}

func (srv *sbConfiguratorServicer) GetMconfigInternal(ctx context.Context, req *protos.GetMconfigRequest) (*protos.GetMconfigResponse, error) {
	store, err := srv.factory.StartTransaction(context.Background(), &orc8rStorage.TxOptions{ReadOnly: true})
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Aborted, "failed to start transaction: %s", err)
	}

	// network ID isn't used in a physical ID query
	loadResult, err := store.LoadEntities("", storage.EntityLoadFilter{PhysicalID: &wrappers.StringValue{Value: req.HardwareID}}, storage.EntityLoadCriteria{})
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "failed to load entity for gateway %s: %s", req.HardwareID, err)
	}
	if funk.IsEmpty(loadResult.Entities) {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.NotFound, "did not find gateway for ID %s", req.HardwareID)
	}

	storage.CommitLogOnError(store)

	ent := loadResult.Entities[0]
	cfg, err := srv.getMconfigImpl(ent.NetworkID, ent.Key)
	if err != nil {
		return nil, err
	}
	return &protos.GetMconfigResponse{Configs: cfg, LogicalID: ent.Key}, nil
}

func (srv *sbConfiguratorServicer) getMconfigImpl(networkID string, gatewayID string) (*commonProtos.GatewayConfigs, error) {
	store, err := srv.factory.StartTransaction(context.Background(), &orc8rStorage.TxOptions{ReadOnly: true})
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Aborted, "failed to start transaction: %s", err)
	}

	graph, err := store.LoadGraphForEntity(
		networkID,
		storage.EntityID{Type: orc8r.MagmadGatewayType, Key: gatewayID},
		storage.FullEntityLoadCriteria,
	)
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "failed to load entity graph: %s", err)
	}

	nwLoad, err := store.LoadNetworks(storage.NetworkLoadFilter{Ids: []string{networkID}}, storage.FullNetworkLoadCriteria)
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "failed to load network: %s", err)
	}
	if !funk.IsEmpty(nwLoad.NetworkIDsNotFound) || funk.IsEmpty(nwLoad.Networks) {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "network %s not found: %s", networkID, err)
	}

	// error on commit is fine for a readonly tx
	storage.CommitLogOnError(store)

	ret, err := configurator.CreateMconfig(networkID, gatewayID, &graph, nwLoad.Networks[0])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build mconfig: %s", err)
	}
	return ret, nil
}
