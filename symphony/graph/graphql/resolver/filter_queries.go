// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/99designs/gqlgen/graphql"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"
)

func (r queryResolver) EquipmentSearch(ctx context.Context, filters []*models.EquipmentFilterInput, limit *int) (*models.EquipmentSearchResult, error) {
	return resolverutil.EquipmentSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

func (r queryResolver) WorkOrderSearch(ctx context.Context, filters []*models.WorkOrderFilterInput, limit *int) (*models.WorkOrderSearchResult, error) {
	return resolverutil.WorkOrderSearch(ctx, r.ClientFrom(ctx), filters, limit, graphql.CollectAllFields(ctx))
}

func (r queryResolver) LinkSearch(ctx context.Context, filters []*models.LinkFilterInput, limit *int) (*models.LinkSearchResult, error) {
	return resolverutil.LinkSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

func (r queryResolver) PortSearch(ctx context.Context, filters []*models.PortFilterInput, limit *int) (*models.PortSearchResult, error) {
	return resolverutil.PortSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

func (r queryResolver) LocationSearch(ctx context.Context, filters []*models.LocationFilterInput, limit *int) (*models.LocationSearchResult, error) {
	return resolverutil.LocationSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

func (r queryResolver) ServiceSearch(ctx context.Context, filters []*models.ServiceFilterInput, limit *int) (*models.ServiceSearchResult, error) {
	return resolverutil.ServiceSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

func (r queryResolver) ProjectSearch(ctx context.Context, filters []*models.ProjectFilterInput, limit *int) ([]*ent.Project, error) {
	var (
		query = r.ClientFrom(ctx).Project.Query()
		err   error
	)
	pros, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying projects failed")
	}
	return pros, nil
}

func (r queryResolver) CustomerSearch(ctx context.Context, limit *int) ([]*ent.Customer, error) {
	var (
		query = r.ClientFrom(ctx).Customer.Query()
		err   error
	)
	pros, err := query.Limit(*limit).All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying customers failed")
	}
	return pros, nil
}
