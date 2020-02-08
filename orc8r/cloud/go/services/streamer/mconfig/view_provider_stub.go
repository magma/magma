/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mconfig

import (
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/any"
)

// GetViewProvider returns a StreamProvider under the name "mconfig_views"
// which returns an empty OffsetGatewayConfigs. This provider exists to
// maintain backwards-compatibility with gateways which still try to access
// this streamer provider.
func GetViewProvider() providers.StreamProvider {
	return &ViewProviderStub{}
}

type ViewProviderStub struct{}

func (*ViewProviderStub) GetStreamName() string {
	return "mconfig_views"
}

func (*ViewProviderStub) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	ret := &protos.OffsetGatewayConfigs{}
	marshaledRet, err := protos.MarshalIntern(ret)
	if err != nil {
		return nil, err
	}

	update := &protos.DataUpdate{
		Key:   gatewayId,
		Value: marshaledRet,
	}
	return []*protos.DataUpdate{update}, nil
}
