// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type propertyTypeResolver struct{}

func (propertyTypeResolver) Type(_ context.Context, obj *ent.PropertyType) (models.PropertyKind, error) {
	return models.PropertyKind(obj.Type), nil
}

type propertyResolver struct{}

func (propertyResolver) PropertyType(ctx context.Context, obj *ent.Property) (*ent.PropertyType, error) {
	return obj.QueryType().Only(ctx)
}

func (propertyResolver) EquipmentValue(ctx context.Context, obj *ent.Property) (*ent.Equipment, error) {
	e, err := obj.QueryEquipmentValue().Only(ctx)
	return e, ent.MaskNotFound(err)
}

func (propertyResolver) LocationValue(ctx context.Context, obj *ent.Property) (*ent.Location, error) {
	e, err := obj.QueryLocationValue().Only(ctx)
	return e, ent.MaskNotFound(err)
}

func (propertyResolver) ServiceValue(ctx context.Context, obj *ent.Property) (*ent.Service, error) {
	e, err := obj.QueryServiceValue().Only(ctx)
	return e, ent.MaskNotFound(err)
}
