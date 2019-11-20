// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"
)

func (r queryResolver) EquipmentSearch(ctx context.Context, filters []*models.EquipmentFilterInput, limit *int) (*models.EquipmentSearchResult, error) {
	return resolverutil.EquipmentSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

func (r queryResolver) WorkOrderSearch(ctx context.Context, filters []*models.WorkOrderFilterInput, limit *int) ([]*ent.WorkOrder, error) {
	var (
		query = r.ClientFrom(ctx).WorkOrder.Query()
		err   error
	)
	for _, f := range filters {
		switch {
		case strings.HasPrefix(f.FilterType.String(), "WORK_ORDER_"):
			if query, err = r.handleWorkOrderFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "LOCATION_INST"):
			if query, err = r.handleWOLocationFilter(query, f); err != nil {
				return nil, err
			}
		}
	}
	wos, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying work orders failed")
	}
	return wos, nil
}

func (r queryResolver) LinkSearch(ctx context.Context, filters []*models.LinkFilterInput, limit *int) (*models.LinkSearchResult, error) {
	return resolverutil.LinkSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

func (r queryResolver) PortSearch(ctx context.Context, filters []*models.PortFilterInput, limit *int) (*models.PortSearchResult, error) {
	return resolverutil.PortSearch(ctx, r.ClientFrom(ctx), filters, limit)
}

// nolint: dupl
func (r queryResolver) LocationSearch(ctx context.Context, filters []*models.LocationFilterInput, limit *int) (*models.LocationSearchResult, error) {
	var (
		query = r.ClientFrom(ctx).Location.Query()
		err   error
	)
	for _, f := range filters {
		switch {
		case strings.HasPrefix(f.FilterType.String(), "LOCATION_INST"):
			if query, err = handleLocationFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "LOCATION_TYPE"):
			if query, err = handleLocationTypeFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "PROPERTY"):
			if query, err = handleLocationPropertyFilter(query, f); err != nil {
				return nil, err
			}
		}
	}
	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Count query failed")
	}
	locs, err := query.Limit(*limit).All(ctx)

	if err != nil {
		return nil, errors.Wrapf(err, "Querying locations failed")
	}

	return &models.LocationSearchResult{
		Locations: locs,
		Count:     count,
	}, err
}

// nolint: dupl
func (r queryResolver) ServiceSearch(ctx context.Context, filters []*models.ServiceFilterInput, limit *int) (*models.ServiceSearchResult, error) {
	var (
		query = r.ClientFrom(ctx).Service.Query()
		err   error
	)
	for _, f := range filters {
		switch {
		case strings.HasPrefix(f.FilterType.String(), "SERVICE_"):
			if query, err = handleServiceFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "LOCATION_INST"):
			if query, err = handleServiceLocationFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "EQUIPMENT_IN_SERVICE"):
			if query, err = handleEquipmentInServiceFilter(query, f); err != nil {
				return nil, err
			}
		}
	}

	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Count query failed")
	}
	services, err := query.All(ctx)

	if err != nil {
		return nil, errors.Wrapf(err, "Querying services failed")
	}
	return &models.ServiceSearchResult{
		Services: services,
		Count:    count,
	}, err
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
