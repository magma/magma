/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"magma/orc8r/cloud/go/protos"
)

type MconfigStreamer interface {
	// Get a list of config types that this streamer should subscribe to.
	// Any changes to a config that matches any of the types returned from
	// this method will be passed to ApplyMconfigUpdate
	GetSubscribedConfigTypes() []string

	// Seed a new gateway's mconfig with necessary values. For example, you
	// would want to fill in network-level configs here. Modify the`mconfigOut`
	// output parameter in-place here.
	SeedNewGatewayMconfig(
		networkId string,
		gatewayId string,
		mconfigOut *protos.GatewayConfigs, // output parameter
	) error

	// Given a config update and a set of old mconfigs, return a new collection
	// of computed mconfigs after the application of the update.
	ApplyMconfigUpdate(
		update *ConfigUpdate,
		oldMconfigsByGatewayId map[string]*protos.GatewayConfigs,
	) (map[string]*protos.GatewayConfigs, error)
}
