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

	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
)

type baseOrchestratorStreamProviderServicer struct{}

func NewBaseOrchestratorStreamProviderServicer() streamer_protos.StreamProviderServer {
	return &baseOrchestratorStreamProviderServicer{}
}

func (s *baseOrchestratorStreamProviderServicer) GetUpdates(ctx context.Context, req *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	var streamer providers.StreamProvider
	switch req.GetStreamName() {
	case definitions.MconfigStreamName:
		streamer = &providers.MconfigProvider{}
	default:
		return nil, fmt.Errorf("GetUpdates failed: unknown stream name provided: %s", req.GetStreamName())
	}

	update, err := streamer.GetUpdates(req.GetGatewayId(), req.GetExtraArgs())
	if err != nil {
		return nil, err
	}
	return &protos.DataUpdateBatch{Updates: update}, nil
}
