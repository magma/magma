/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer

import (
	"sort"

	"magma/lte/cloud/go/lte"
	models2 "magma/lte/cloud/go/plugin/models"
	protos2 "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

type SubscribersProvider struct{}

func (provider *SubscribersProvider) GetStreamName() string {
	return "subscriberdb"
}

func (provider *SubscribersProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	ent, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	subEnts, err := configurator.LoadAllEntitiesInNetwork(ent.NetworkID, lte.SubscriberEntityType, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsToThis: true})
	if err != nil {
		return nil, err
	}

	subProtos := make([]*protos2.SubscriberData, 0, len(subEnts))
	for _, sub := range subEnts {
		subProto := &protos2.SubscriberData{}
		subProto, err = subscriberToMconfig(sub)
		if err != nil {
			return nil, err
		}
		subProto.NetworkId = &protos.NetworkID{Id: ent.NetworkID}
		subProtos = append(subProtos, subProto)
	}
	return subscribersToUpdates(subProtos)
}

func subscribersToUpdates(subs []*protos2.SubscriberData) ([]*protos.DataUpdate, error) {
	ret := make([]*protos.DataUpdate, 0, len(subs))
	for _, sub := range subs {
		marshaledProto, err := proto.Marshal(sub)
		if err != nil {
			return nil, err
		}
		update := &protos.DataUpdate{Key: protos2.SidString(sub.Sid), Value: marshaledProto}
		ret = append(ret, update)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Key < ret[j].Key })
	return ret, nil
}

func subscriberToMconfig(ent configurator.NetworkEntity) (*protos2.SubscriberData, error) {
	sub := &protos2.SubscriberData{}
	t, err := protos2.SidProto(ent.Key)
	if err != nil {
		return nil, err
	}

	sub.Sid = t
	if ent.Config == nil {
		return sub, nil
	}

	cfg := ent.Config.(*models2.LteSubscription)
	sub.Lte = &protos2.LTESubscription{
		State:    protos2.LTESubscription_LTESubscriptionState(protos2.LTESubscription_LTESubscriptionState_value[cfg.State]),
		AuthAlgo: protos2.LTESubscription_LTEAuthAlgo(protos2.LTESubscription_LTEAuthAlgo_value[cfg.AuthAlgo]),
		AuthKey:  cfg.AuthKey,
		AuthOpc:  cfg.AuthOpc,
	}

	if cfg.SubProfile != "" {
		sub.SubProfile = string(cfg.SubProfile)
	} else {
		sub.SubProfile = "default"
	}

	for _, assoc := range ent.ParentAssociations {
		if assoc.Type == lte.BaseNameEntityType {
			sub.Lte.AssignedBaseNames = append(sub.Lte.AssignedBaseNames, assoc.Key)
		} else if assoc.Type == lte.PolicyRuleEntityType {
			sub.Lte.AssignedPolicies = append(sub.Lte.AssignedPolicies, assoc.Key)
		}
	}

	return sub, nil
}
