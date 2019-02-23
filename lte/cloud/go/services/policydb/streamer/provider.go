/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer

import (
	"magma/lte/cloud/go/services/policydb"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

type PoliciesProvider struct{}

func (provider *PoliciesProvider) GetStreamName() string {
	return "policydb"
}

func (provider *PoliciesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	networkId, err := magmad.FindGatewayNetworkId(gatewayId)
	if err != nil {
		return nil, err
	}
	policies, err := policydb.GetAllRules(networkId)
	if err != nil {
		return nil, err
	}

	ret := make([]*protos.DataUpdate, 0, len(policies))
	for _, policy := range policies {
		marshaledPolicy, err := proto.Marshal(policy)
		if err != nil {
			return nil, err
		}

		update := new(protos.DataUpdate)
		update.Key = policy.Id
		update.Value = marshaledPolicy
		ret = append(ret, update)
	}
	return ret, nil
}

type BaseNamesProvider struct{}

func (provider *BaseNamesProvider) GetStreamName() string {
	return "base_names"
}

func (provider *BaseNamesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	networkId, err := magmad.FindGatewayNetworkId(gatewayId)
	if err != nil {
		return nil, err
	}
	baseNameRecords, err := policydb.GetAllBaseNames(networkId)
	if err != nil {
		return nil, err
	}

	ret := make([]*protos.DataUpdate, 0, len(baseNameRecords))
	for _, baseNameRecord := range baseNameRecords {
		marshaledBaseNameSet, err := proto.Marshal(baseNameRecord.GetRuleNamesSet())
		if err != nil {
			return nil, err
		}

		update := new(protos.DataUpdate)
		update.Key = baseNameRecord.GetName()
		update.Value = marshaledBaseNameSet
		ret = append(ret, update)
	}
	return ret, nil
}
