/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package legacy_stream_providers

import (
	"context"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/any"
)

// LegacyProviderFactor defines a factory for LegacyProviders
type LegacyProviderFactory struct{}

// CreateLegacyProvider creates a LegacyProvider with the provided stream name
// and servicer implementation
func (f *LegacyProviderFactory) CreateLegacyProvider(stream string, provider protos.StreamProviderServer) *LegacyProvider {
	return &LegacyProvider{
		streamName:             stream,
		streamProviderServicer: provider,
	}
}

// LegacyProvider implements the legacy StreamProvider plugin interface
type LegacyProvider struct {
	streamName             string
	streamProviderServicer protos.StreamProviderServer
}

// GetStreamName returns the stream name of the legacy provider
func (provider *LegacyProvider) GetStreamName() string {
	return provider.streamName
}

// GetUpdates implements the GetUpdates method of the legacy StreamProvider
// plugin interface by calling the provider's servicer implementation
func (provider *LegacyProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	streamReq := &protos.StreamRequest{
		GatewayId:  gatewayId,
		StreamName: provider.GetStreamName(),
		ExtraArgs:  extraArgs,
	}
	updateRes, err := provider.streamProviderServicer.GetUpdates(context.Background(), streamReq)
	if err != nil {
		return []*protos.DataUpdate{}, err
	}
	return updateRes.GetUpdates(), nil
}
