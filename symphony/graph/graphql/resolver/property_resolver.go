// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
)

type propertyTypeResolver struct{}

func (propertyTypeResolver) RawValue(ctx context.Context, propertyType *ent.PropertyType) (*string, error) {
	raw, err := resolverutil.PropertyValue(ctx, propertyType.Type, propertyType.NodeType, propertyType)
	return &raw, err
}

type propertyResolver struct{}

func (propertyResolver) RawValue(ctx context.Context, property *ent.Property) (*string, error) {
	propertyType, err := property.QueryType().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying property type %w", err)
	}
	raw, err := resolverutil.PropertyValue(ctx, propertyType.Type, propertyType.NodeType, property)
	return &raw, err
}

func (propertyResolver) PropertyType(ctx context.Context, obj *ent.Property) (*ent.PropertyType, error) {
	typ, err := obj.Edges.TypeOrErr()
	if ent.IsNotLoaded(err) {
		return obj.QueryType().Only(ctx)
	}
	return typ, err
}

func (propertyResolver) NodeValue(ctx context.Context, property *ent.Property) (models.NamedNode, error) {
	propertyType, err := property.QueryType().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying property type %w", err)
	}
	switch propertyType.NodeType {
	case resolverutil.NodeTypeLocation:
		l, err := property.QueryLocationValue().Only(ctx)
		return l, ent.MaskNotFound(err)
	case resolverutil.NodeTypeEquipment:
		e, err := property.QueryEquipmentValue().Only(ctx)
		return e, ent.MaskNotFound(err)
	case resolverutil.NodeTypeService:
		s, err := property.QueryServiceValue().Only(ctx)
		return s, ent.MaskNotFound(err)
	case resolverutil.NodeTypeWorkOrder:
		s, err := property.QueryWorkOrderValue().Only(ctx)
		return s, ent.MaskNotFound(err)
	case resolverutil.NodeTypeUser:
		s, err := property.QueryUserValue().Only(ctx)
		return s, ent.MaskNotFound(err)
	default:
		return nil, nil
	}
}
