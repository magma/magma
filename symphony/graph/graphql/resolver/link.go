// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type linkResolver struct{}

func (linkResolver) Ports(ctx context.Context, obj *ent.Link) ([]*ent.EquipmentPort, error) {
	return obj.QueryPorts().All(ctx)
}

func (linkResolver) FutureState(ctx context.Context, obj *ent.Link) (*models.FutureState, error) {
	fs := models.FutureState(obj.FutureState)
	return &fs, nil
}

func (linkResolver) WorkOrder(ctx context.Context, obj *ent.Link) (*ent.WorkOrder, error) {
	wo, err := obj.QueryWorkOrder().Only(ctx)
	return wo, ent.MaskNotFound(err)
}

func (linkResolver) Properties(ctx context.Context, obj *ent.Link) ([]*ent.Property, error) {
	return obj.QueryProperties().All(ctx)
}

func (linkResolver) Services(ctx context.Context, obj *ent.Link) ([]*ent.Service, error) {
	return obj.QueryService().All(ctx)
}
