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

	"magma/lte/cloud/go/lte"
	policydb_streamer "magma/lte/cloud/go/services/policydb/streamer"
	subscriber_streamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/protos"
)

type lteStreamProviderServicer struct{}

func NewLTEStreamProviderServicer() streamer_protos.StreamProviderServer {
	return &lteStreamProviderServicer{}
}

func (s *lteStreamProviderServicer) GetUpdates(ctx context.Context, req *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	var streamer providers.StreamProvider
	switch req.GetStreamName() {
	case lte.SubscriberStreamName:
		streamer = &subscriber_streamer.SubscribersProvider{}
	case lte.PolicyStreamName:
		streamer = &policydb_streamer.PoliciesProvider{}
	case lte.BaseNameStreamName:
		streamer = &policydb_streamer.BaseNamesProvider{}
	case lte.MappingsStreamName:
		streamer = &policydb_streamer.RuleMappingsProvider{}
	case lte.NetworkWideRulesStreamName:
		streamer = &policydb_streamer.NetworkWideRulesProvider{}
	case lte.RatingGroupStreamName:
		streamer = &policydb_streamer.RatingGroupsProvider{}
	default:
		return nil, fmt.Errorf("GetUpdates failed: unknown stream name provided: %s", req.GetStreamName())
	}

	updates, err := streamer.GetUpdates(req.GetGatewayId(), req.GetExtraArgs())
	if err != nil {
		return &protos.DataUpdateBatch{}, err
	}
	return &protos.DataUpdateBatch{Updates: updates}, nil
}
