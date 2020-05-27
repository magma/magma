// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
)

type (
	serviceEndpointTypeResolver struct{}
)

func (s serviceEndpointTypeResolver) Endpoints(ctx context.Context, obj *ent.ServiceEndpointDefinition) ([]*ent.ServiceEndpoint, error) {
	return obj.QueryEndpoints().All(ctx)
}

func (s serviceEndpointTypeResolver) EquipmentType(ctx context.Context, obj *ent.ServiceEndpointDefinition) (*ent.EquipmentType, error) {
	et, err := obj.QueryEquipmentType().Only(ctx)
	return et, ent.MaskNotFound(err)
}

func (s serviceEndpointTypeResolver) ServiceType(ctx context.Context, obj *ent.ServiceEndpointDefinition) (*ent.ServiceType, error) {
	return obj.QueryServiceType().Only(ctx)
}
