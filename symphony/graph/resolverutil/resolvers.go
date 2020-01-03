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

// LocationTypes is a helper to bring location types
func LocationTypes(ctx context.Context, client *ent.Client) (*models.LocationTypeConnection, error) {
	lts, err := client.LocationType.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying location types")
	}
	edges := make([]*models.LocationTypeEdge, len(lts))
	for i, lt := range lts {
		edges[i] = &models.LocationTypeEdge{Node: lt}
	}
	return &models.LocationTypeConnection{Edges: edges}, err
}

// EquipmentTypes is a helper to bring equipment types
func EquipmentTypes(ctx context.Context, client *ent.Client) (*models.EquipmentTypeConnection, error) {
	ets, err := client.EquipmentType.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying equipment types")
	}
	edges := make([]*models.EquipmentTypeEdge, len(ets))
	for i, et := range ets {
		edges[i] = &models.EquipmentTypeEdge{Node: et}
	}
	return &models.EquipmentTypeConnection{Edges: edges}, err
}

// ServiceTypes is a helper to bring service types
func ServiceTypes(ctx context.Context, client *ent.Client) (*models.ServiceTypeConnection, error) {
	sts, err := client.ServiceType.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying service types")
	}
	edges := make([]*models.ServiceTypeEdge, len(sts))
	for i, et := range sts {
		edges[i] = &models.ServiceTypeEdge{Node: et}
	}
	return &models.ServiceTypeConnection{Edges: edges}, err
}

// EquipmentPortTypes is a helper to bring equipment port types
func EquipmentPortTypes(ctx context.Context, client *ent.Client) (*models.EquipmentPortTypeConnection, error) {
	ets, err := client.EquipmentPortType.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying equipment types")
	}
	edges := make([]*models.EquipmentPortTypeEdge, len(ets))
	for i, et := range ets {
		edges[i] = &models.EquipmentPortTypeEdge{Node: et}
	}
	return &models.EquipmentPortTypeConnection{Edges: edges}, err
}

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
	}, err
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
	}, err
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
	}, err
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
	}, err
}
