// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"context"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/pkg/errors"
)

// EquipmentTypes is a helper to bring equipment types
func EquipmentSearch(ctx context.Context, client *ent.Client, filters []*models.EquipmentFilterInput, limit *int) (*models.EquipmentSearchResult, error) {
	var (
		res []*ent.Equipment
		c   int
		err error
	)
	// TODO T46957221 support query Clone
	for i := 0; i < 2; i++ {
		query := client.Equipment.Query()
		for _, f := range filters {
			switch {
			case strings.HasPrefix(f.FilterType.String(), "EQUIPMENT_TYPE"):
				if query, err = handleEquipmentTypeFilter(query, f); err != nil {
					return nil, err
				}
			case strings.HasPrefix(f.FilterType.String(), "EQUIP_INST"), strings.HasPrefix(f.FilterType.String(), "PROPERTY"):
				if query, err = handleEquipmentFilter(query, f); err != nil {
					return nil, err
				}
			case strings.HasPrefix(f.FilterType.String(), "LOCATION_INST"):
				if query, err = handleEquipmentLocationFilter(query, f); err != nil {
					return nil, err
				}
			}
		}
		if i == 0 {
			c, err = query.Count(ctx)
			if err != nil {
				return nil, err
			}
			continue
		}
		if limit != nil {
			query.Limit(*limit)
		}
		res, err = query.Order(ent.Asc(equipment.FieldName)).All(ctx)
		if err != nil {
			return nil, err
		}
	}
	return &models.EquipmentSearchResult{
		Equipment: res,
		Count:     c,
	}, nil
}

// nolint: dupl
func PortSearch(ctx context.Context, client *ent.Client, filters []*models.PortFilterInput, limit *int) (*models.PortSearchResult, error) {
	var (
		query = client.EquipmentPort.Query()
		err   error
	)
	for _, f := range filters {
		switch {
		case strings.HasPrefix(f.FilterType.String(), "PORT_INST"):
			if query, err = handlePortFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "LOCATION_INST"):
			if query, err = handlePortLocationFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "PORT_DEF"):
			if query, err = handlePortDefinitionFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "PROPERTY"):
			if query, err = handlePortPropertyFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "SERVICE_INST"):
			if query, err = handlePortServiceFilter(query, f); err != nil {
				return nil, err
			}
		}
	}
	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Count query failed")
	}
	if limit != nil {
		query.Limit(*limit)
	}
	ports, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying links failed")
	}
	return &models.PortSearchResult{
		Ports: ports,
		Count: count,
	}, nil
}

// nolint: dupl
func LocationSearch(ctx context.Context, client *ent.Client, filters []*models.LocationFilterInput, limit *int) (*models.LocationSearchResult, error) {
	var (
		query = client.Location.Query()
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
	if limit != nil {
		query.Limit(*limit)
	}
	locs, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying locations failed")
	}
	return &models.LocationSearchResult{
		Locations: locs,
		Count:     count,
	}, nil
}

// nolint: dupl
func LinkSearch(ctx context.Context, client *ent.Client, filters []*models.LinkFilterInput, limit *int) (*models.LinkSearchResult, error) {
	var (
		query = client.Link.Query()
		err   error
	)
	for _, f := range filters {
		switch {
		case strings.HasPrefix(f.FilterType.String(), "LINK_"):
			if query, err = handleLinkFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "LOCATION_INST"):
			if query, err = handleLinkLocationFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "EQUIPMENT_"):
			if query, err = handleLinkEquipmentFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "SERVICE_INST"):
			if query, err = handleLinkServiceFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "PROPERTY"):
			if query, err = handleLinkPropertyFilter(query, f); err != nil {
				return nil, err
			}
		}
	}

	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Count query failed")
	}
	if limit != nil {
		query.Limit(*limit)
	}
	links, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying links failed")
	}
	return &models.LinkSearchResult{
		Links: links,
		Count: count,
	}, nil
}

// nolint: dupl
func ServiceSearch(ctx context.Context, client *ent.Client, filters []*models.ServiceFilterInput, limit *int) (*models.ServiceSearchResult, error) {
	var (
		query = client.Service.Query()
		err   error
	)
	for _, f := range filters {
		switch {
		case strings.HasPrefix(f.FilterType.String(), "SERVICE_"):
			if query, err = handleServiceFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "PROPERTY"):
			if query, err = handleServicePropertyFilter(query, f); err != nil {
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
	if limit != nil {
		query.Limit(*limit)
	}
	services, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying services failed")
	}
	return &models.ServiceSearchResult{
		Services: services,
		Count:    count,
	}, nil
}

func WorkOrderSearch(ctx context.Context, client *ent.Client, filters []*models.WorkOrderFilterInput, limit *int, fields []string) (*models.WorkOrderSearchResult, error) {
	var (
		query = client.WorkOrder.Query()
		err   error
	)
	for _, f := range filters {
		switch {
		case strings.HasPrefix(f.FilterType.String(), "WORK_ORDER_"):
			if query, err = handleWorkOrderFilter(query, f); err != nil {
				return nil, err
			}
		case strings.HasPrefix(f.FilterType.String(), "LOCATION_INST"):
			if query, err = handleWOLocationFilter(query, f); err != nil {
				return nil, err
			}
		}
	}

	var woResult models.WorkOrderSearchResult
	for _, field := range fields {
		switch field {
		case "count":
			woResult.Count, err = query.Clone().Count(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "Count query failed")
			}
		case "workOrders":
			if limit != nil {
				query.Limit(*limit)
			}
			woResult.WorkOrders, err = query.All(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "Querying work orders failed")
			}
		}

	}
	return &woResult, nil
}

func UserSearch(ctx context.Context, client *ent.Client, filters []*models.UserFilterInput, limit *int) (*models.UserSearchResult, error) {
	var (
		query = client.User.Query()
		err   error
	)
	for _, f := range filters {
		if strings.HasPrefix(f.FilterType.String(), "USER_") {
			if query, err = handleUserFilter(query, f); err != nil {
				return nil, err
			}
		}
	}

	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Count query failed")
	}
	if limit != nil {
		query.Limit(*limit)
	}
	users, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying users failed")
	}
	return &models.UserSearchResult{
		Users: users,
		Count: count,
	}, nil
}

func UsersGroupSearch(
	ctx context.Context,
	client *ent.Client,
	filters []*models.UsersGroupFilterInput,
	limit *int,
) (*models.UsersGroupSearchResult, error) {
	var (
		query = client.UsersGroup.Query()
		err   error
	)
	for _, f := range filters {
		if strings.HasPrefix(f.FilterType.String(), "GROUP_NAME") {
			if query, err = handleUsersGroupFilter(query, f); err != nil {
				return nil, err
			}
		}
	}

	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Count query failed")
	}
	if limit != nil {
		query.Limit(*limit)
	}
	usersGroups, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Querying usersGroups failed")
	}
	return &models.UsersGroupSearchResult{
		UsersGroups: usersGroups,
		Count:       count,
	}, nil
}
