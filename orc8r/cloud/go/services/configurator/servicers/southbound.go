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
	commonProtos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"

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

	store, err := srv.factory.StartTransaction(context.Background(), &storage.TxOptions{ReadOnly: true})
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Aborted, "failed to start transaction: %s", err)
	}

	graph, err := store.LoadGraphForEntity(
		gw.NetworkId,
		storage.EntityID{Type: orc8r.MagmadGatewayType, Key: gw.LogicalId},
		storage.FullEntityLoadCriteria,
	)
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "failed to load entity graph: %s", err)
	}

	nwLoad, err := store.LoadNetworks([]string{gw.NetworkId}, storage.FullNetworkLoadCriteria)
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "failed to load network: %s", err)
	}
	if !funk.IsEmpty(nwLoad.NetworkIDsNotFound) || funk.IsEmpty(nwLoad.Networks) {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "network %s not found: %s", gw.NetworkId, err)
	}

	// error on commit is fine for a readonly tx
	storage.CommitLogOnError(store)

	ret, err := configurator.CreateMconfig(gw.NetworkId, gw.LogicalId, &graph, nwLoad.Networks[0])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build mconfig: %s", err)
	}
	return ret, nil
}
