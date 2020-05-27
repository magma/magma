/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package stream_provider

import (
	"context"
	"fmt"

	"magma/lte/cloud/go/lte"
	policydbstreamer "magma/lte/cloud/go/services/policydb/streamer"
	subscriberstreamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/lib/go/protos"
)

// LteStreamProviderServicer implements the StreamProvider RPC servicer for the
// LTE module
type LteStreamProviderServicer struct{}

// GetUpdates fetches updates for a given stream by calling the associated
// stream provider for that stream name
func (s *LteStreamProviderServicer) GetUpdates(context context.Context, request *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	var updateRes []*protos.DataUpdate
	var err error
	switch request.GetStreamName() {
	case lte.SubscriberStreamName:
		subscriberStreamer := subscriberstreamer.SubscribersProvider{}
		updateRes, err = subscriberStreamer.GetUpdatesImpl(request.GetGatewayId(), request.GetExtraArgs())
	case lte.PolicyStreamName:
		policyStreamer := policydbstreamer.PoliciesProvider{}
		updateRes, err = policyStreamer.GetUpdatesImpl(request.GetGatewayId(), request.GetExtraArgs())
	case lte.BaseNameStreamName:
		baseNameStreamer := policydbstreamer.BaseNamesProvider{}
		updateRes, err = baseNameStreamer.GetUpdatesImpl(request.GetGatewayId(), request.GetExtraArgs())
	case lte.MappingsStreamName:
		ruleMappingsStreamer := policydbstreamer.RuleMappingsProvider{}
		updateRes, err = ruleMappingsStreamer.GetUpdatesImpl(request.GetGatewayId(), request.GetExtraArgs())
	case lte.NetworkWideRules:
		networkWideRulesStreamer := policydbstreamer.NetworkWideRulesProvider{}
		updateRes, err = networkWideRulesStreamer.GetUpdatesImpl(request.GetGatewayId(), request.GetExtraArgs())
	default:
		return nil, fmt.Errorf("GetUpdates failed: unknown stream name provided: %s", request.GetStreamName())
	}
	if err != nil {
		return &protos.DataUpdateBatch{}, err
	}
	return &protos.DataUpdateBatch{
		Updates: updateRes,
	}, nil
}
