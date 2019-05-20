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

type sbConfiguratorServicer struct {
	factory storage.ConfiguratorStorageFactory
}

func NewSouthboundConfiguratorServicer(factory storage.ConfiguratorStorageFactory) (protos.SouthboundConfiguratorServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	return &sbConfiguratorServicer{factory}, nil
}

func (srv *sbConfiguratorServicer) GetMconfig(context context.Context, void *commonProtos.Void) (*commonProtos.GatewayConfigs, error) {
	return &commonProtos.GatewayConfigs{}, nil
}
