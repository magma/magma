/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package stream_providers

import (
	"context"
	"fmt"

	"magma/orc8r/cloud/go/services/streamer/mconfig"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
)

// BaseOrchestratorStreamProviderServicer implements the StreamProvider RPC
// servicer for the BaseOrchestrator module
type BaseOrchestratorStreamProviderServicer struct{}

// GetUpdates fetches updates for a given stream by calling the associated
// stream provider for that stream name
func (s *BaseOrchestratorStreamProviderServicer) GetUpdates(context context.Context, request *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	switch request.GetStreamName() {
	case definitions.MconfigStreamName:
		mconfigStreamer := mconfig.ConfigProvider{}
		updateRes, err := mconfigStreamer.GetUpdatesImpl(request.GetGatewayId(), request.GetExtraArgs())
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
