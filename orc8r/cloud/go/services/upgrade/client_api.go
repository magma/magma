/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package upgrade

import (
	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	upgrade_protos "magma/orc8r/cloud/go/services/upgrade/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	ServiceName = "UPGRADE"
	// type used in configurator to identify network entities that are related to release channels
	UpgradeReleaseChannelEntityType = "upgrade_release_channel"
	// type used in configurator to identify network entities that are related to network tier
	UpgradeTierEntityType = "upgrade_tier"
)

func getUpgradeServiceClient() (upgrade_protos.UpgradeServiceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return upgrade_protos.NewUpgradeServiceClient(conn), err
}

// A channel in this context is a way to partition released packages
// into groups by stability. For example, the "stable" channel will always
// have the most well-tested and stable software, whereas a "beta" or
// "staging" channel will have more recent software that is still
// being tested.

// Create a new global release channel.
func CreateReleaseChannel(channelName string, channel *upgrade_protos.ReleaseChannel) error {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return err
	}

	req := &upgrade_protos.CreateOrUpdateReleaseChannelRequest{
		ChannelName: channelName,
		Channel:     channel,
	}
	_, err = client.CreateReleaseChannel(context.Background(), req)
	return err
}

func GetReleaseChannel(channel string) (*upgrade_protos.ReleaseChannel, error) {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return nil, err
	}

	req := &upgrade_protos.GetReleaseChannelRequest{ChannelName: channel}
	res, err := client.GetReleaseChannel(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, err
}

func ListReleaseChannels() ([]string, error) {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return nil, err
	}

	req := &protos.Void{}
	res, err := client.ListReleaseChannels(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.GetChannelIds(), err
}

func UpdateReleaseChannel(channelName string, updatedChannel *upgrade_protos.ReleaseChannel) error {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return err
	}

	req := &upgrade_protos.CreateOrUpdateReleaseChannelRequest{
		ChannelName: channelName,
		Channel:     updatedChannel,
	}
	_, err = client.UpdateReleaseChannel(context.Background(), req)
	return err
}

func DeleteReleaseChannel(channel string) error {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return err
	}

	req := &upgrade_protos.DeleteReleaseChannelRequest{ChannelName: channel}
	_, err = client.DeleteReleaseChannel(context.Background(), req)
	return err
}

// A tier is a way to partition gateways in a network so they can
// be targeted for software upgrades. Tier membership for gateways
// is defined in the gateway config.

// Get information about some tiers on a network.
// If no tier filter is provided, all tiers in the network will be
// queried.
// If any non-existent tiers are provided in the filter parameter, an
// error will be returned.
func GetTiers(networkId string, tierFilter []string) (map[string]*upgrade_protos.TierInfo, error) {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return map[string]*upgrade_protos.TierInfo{}, err
	}

	req := &upgrade_protos.GetTiersRequest{NetworkId: networkId, TierFilter: tierFilter}
	res, err := client.GetTiers(context.Background(), req)
	if err != nil {
		return map[string]*upgrade_protos.TierInfo{}, err
	}
	return res.GetTiers(), err
}

func CreateTier(networkId string, tierId string, tierInfo *upgrade_protos.TierInfo) error {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return err
	}

	req := &upgrade_protos.CreateTierRequest{
		NetworkId: networkId,
		TierId:    tierId,
		TierInfo:  tierInfo,
	}
	_, err = client.CreateTier(context.Background(), req)
	return err
}

func UpdateTier(networkId string, tierId string, tierInfo *upgrade_protos.TierInfo) error {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return err
	}

	req := &upgrade_protos.UpdateTierRequest{
		NetworkId:   networkId,
		TierId:      tierId,
		UpdatedTier: tierInfo,
	}
	_, err = client.UpdateTier(context.Background(), req)
	return err
}

func DeleteTier(networkId string, tierId string) error {
	client, err := getUpgradeServiceClient()
	if err != nil {
		return err
	}

	req := &upgrade_protos.DeleteTierRequest{
		NetworkId:      networkId,
		TierIdToDelete: tierId,
	}
	_, err = client.DeleteTier(context.Background(), req)
	return err
}
