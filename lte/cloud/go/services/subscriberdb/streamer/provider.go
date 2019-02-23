/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer

import (
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

type SubscribersProvider struct{}

func (provider *SubscribersProvider) GetStreamName() string {
	return "subscriberdb"
}

func (provider *SubscribersProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	networkId, err := magmad.FindGatewayNetworkId(gatewayId)
	if err != nil {
		return nil, err
	}
	subscriberIds, err := subscriberdb.ListSubscribers(networkId)
	if err != nil {
		return nil, err
	}

	ret := make([]*protos.DataUpdate, 0, len(subscriberIds))
	for _, subscriberId := range subscriberIds {
		subscriberData, err := subscriberdb.GetSubscriber(networkId, subscriberId)
		if err != nil {
			return nil, err
		}
		marshaledSubscriber, err := proto.Marshal(subscriberData)
		if err != nil {
			return nil, err
		}

		update := new(protos.DataUpdate)
		update.Key = subscriberId
		update.Value = marshaledSubscriber
		ret = append(ret, update)
	}
	return ret, nil
}
