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

	commonProtos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"
)

type nbConfiguratorServicer struct {
	factory storage.ConfiguratorStorageFactory
}

// NewNorthboundConfiguratorServicer returns a configurator server backed by storage passed in
func NewNorthboundConfiguratorServicer(factory storage.ConfiguratorStorageFactory) (protos.NorthboundConfiguratorServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	return &nbConfiguratorServicer{factory}, nil
}

func (srv *nbConfiguratorServicer) ListNetworks(context context.Context, void *commonProtos.Void) (*protos.ListNetworksResponse, error) {
	return &protos.ListNetworksResponse{}, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) LoadNetworks(context context.Context, req *protos.LoadNetworksRequest) (*protos.LoadNetworksResponse, error) {
	return &protos.LoadNetworksResponse{}, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) CreateNetworks(context context.Context, req *protos.CreateNetworksRequest) (*protos.CreateNetworksResponse, error) {
	return &protos.CreateNetworksResponse{}, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) UpdateNetworks(context context.Context, req *protos.UpdateNetworksRequest) (*protos.UpdateNetworksResponse, error) {
	return &protos.UpdateNetworksResponse{}, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) DeleteNetworks(context context.Context, req *protos.DeleteNetworksRequest) (*commonProtos.Void, error) {
	return &commonProtos.Void{}, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) LoadEntities(context context.Context, req *protos.LoadEntitiesRequest) (*protos.LoadEntitiesResponse, error) {
	res := &protos.LoadEntitiesResponse{}
	return res, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) CreateEntities(context context.Context, req *protos.CreateEntitiesRequest) (*protos.CreateEntitiesResponse, error) {
	return &protos.CreateEntitiesResponse{}, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) UpdateEntities(context context.Context, req *protos.UpdateEntitiesRequest) (*protos.UpdateEntitiesResponse, error) {
	return &protos.UpdateEntitiesResponse{}, fmt.Errorf("Not yet implemented")
}

func (srv *nbConfiguratorServicer) DeleteEntities(context context.Context, req *protos.DeleteEntitiesRequest) (*commonProtos.Void, error) {
	return &commonProtos.Void{}, fmt.Errorf("Not yet implemented")
}
