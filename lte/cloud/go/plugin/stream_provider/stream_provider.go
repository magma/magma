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
	subscriberstreamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/lib/go/protos"
)

// LteStreamProviderServicer implements the StreamProvider RPC servicer for the
// LTE module
type LteStreamProviderServicer struct{}

// GetUpdates fetches updates for a given stream by calling the associated
// stream provider for that stream name
func (s *LteStreamProviderServicer) GetUpdates(context context.Context, request *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	switch request.GetStreamName() {
	case lte.SubscriberStreamName:
		subscriberStreamer := subscriberstreamer.SubscribersProvider{}
		updateRes, err := subscriberStreamer.GetUpdatesImpl(request.GetGatewayId(), request.GetExtraArgs())
		if err != nil {
			return &protos.DataUpdateBatch{}, err
		}
		return &protos.DataUpdateBatch{
			Updates: updateRes,
		}, nil
	default:
		return nil, fmt.Errorf("GetUpdates failed: unknown stream name provided: %s", request.GetStreamName())
	}
}
