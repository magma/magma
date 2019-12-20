/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mconfig

import (
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/streamer/providers"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

// GetProvider returns the StreamProvider for on demand mconfigs.
func GetProvider() providers.StreamProvider {
	return &ConfigProvider{}
}

type ConfigProvider struct{}

func (provider *ConfigProvider) GetStreamName() string {
	return "configs"
}

func (provider *ConfigProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	resp, err := configurator.GetMconfigFor(gatewayId)
	if err != nil {
		return nil, err
	}

	if extraArgs != nil {
		// Currently, only use of extraArgs is mconfig hash
		receivedDigest := &protos.GatewayConfigsDigest{}
		if err := ptypes.UnmarshalAny(extraArgs, receivedDigest); err == nil {
			glog.V(2).Infof("Received, generated config digests: %v, %v\n",
				receivedDigest,
				resp.Configs.Metadata.Digest.Md5HexDigest,
			)
			return mconfigToUpdate(resp.Configs, resp.LogicalID, receivedDigest.Md5HexDigest)
		}
	}

	return mconfigToUpdate(resp.Configs, resp.LogicalID, "")
}

func mconfigToUpdate(configs *protos.GatewayConfigs, key string, digest string) ([]*protos.DataUpdate, error) {
	// Early/empty return if gateway already has config that would be sent here
	if digest == configs.Metadata.Digest.Md5HexDigest {
		return []*protos.DataUpdate{}, nil
	}

	marshaledConfig, err := protos.MarshalIntern(configs)
	if err != nil {
		return nil, err
	}
	update := new(protos.DataUpdate)
	update.Key = key
	update.Value = marshaledConfig
	return []*protos.DataUpdate{update}, nil
}
