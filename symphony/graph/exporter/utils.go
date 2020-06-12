// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"

	"github.com/facebookincubator/symphony/pkg/ent/location"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentport"
	"github.com/facebookincubator/symphony/pkg/ent/link"
	"github.com/facebookincubator/symphony/pkg/ent/service"

	"github.com/pkg/errors"
)

const (
	maxEquipmentParents = 3
)

func toIntSlice(a []string) ([]int, error) {
	var intSlice []int
	for _, i := range a {
		j, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, j)
	}
	return intSlice, nil
}

func index(a []string, x string) int {
	for i, n := range a {
		if strings.EqualFold(x, n) {
			return i
		}
	}
	return -1
}

func getQueryFields(e ExportEntity) []string {
	var v reflect.Value
	switch e {
	case ExportEntityWorkOrders:
		model := models.WorkOrderSearchResult{}
		v = reflect.ValueOf(&model).Elem()
	case ExportEntityLocation:
		model := models.LocationSearchResult{}
		v = reflect.ValueOf(&model).Elem()
	case ExportEntityPort:
		model := models.PortSearchResult{}
		v = reflect.ValueOf(&model).Elem()
	case ExportEntityEquipment:
		model := models.EquipmentSearchResult{}
		v = reflect.ValueOf(&model).Elem()
	case ExportEntityLink:
		model := models.LinkSearchResult{}
		v = reflect.ValueOf(&model).Elem()
	case ExportEntityService:
		model := models.ServiceSearchResult{}
		v = reflect.ValueOf(&model).Elem()
	default:
		return []string{}
	}

	fields := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		a := []rune(v.Type().Field(i).Name)
		a[0] = unicode.ToLower(a[0])
		fields[i] = string(a)
	}
	return fields
}

func locationTypeHierarchy(ctx context.Context, c *ent.Client) ([]string, error) {
	locTypeResult, err := c.LocationType.Query().
		Paginate(ctx, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	sortedEnts := locTypeResult.Edges
	sort.Slice(sortedEnts, func(i, j int) bool {
		return sortedEnts[i].Node.Index < sortedEnts[j].Node.Index
	})

	var hierarchy = make([]string, len(sortedEnts))
	for i, loc := range sortedEnts {
		name := loc.Node.Name
		if index(hierarchy, name) != -1 {
			return nil, errors.Errorf("duplicate location type names %s", name)
		}
		hierarchy[i] = name
	}
	return hierarchy, nil
}

func parentHierarchy(ctx context.Context, equipment ent.Equipment) []string {
	var parents = make([]string, maxEquipmentParents)
	pos, _ := equipment.QueryParentPosition().Only(ctx)
	for i := maxEquipmentParents - 1; i >= 0; i-- {
		if pos == nil {
			break
		}
		parentEquipment := pos.QueryParent().OnlyX(ctx)
		parents[i] = parentEquipment.Name
		pos, _ = parentEquipment.QueryParentPosition().Only(ctx)
	}
	return parents
}

func parentHierarchyWithAllPositions(ctx context.Context, equipment ent.Equipment) []string {
	var parents = make([]string, 2*maxEquipmentParents)
	pos, _ := equipment.QueryParentPosition().Only(ctx)
	for i := (2 * maxEquipmentParents) - 1; i >= 1; i -= 2 {
		if pos == nil {
			break
		}
		parentEquipment := pos.QueryParent().OnlyX(ctx)
		parents[i] = pos.QueryDefinition().OnlyX(ctx).Name
		parents[i-1] = parentEquipment.Name
		pos, _ = parentEquipment.QueryParentPosition().Only(ctx)
	}
	return parents
}

func locationHierarchyForEquipment(ctx context.Context, equipment *ent.Equipment, orderedLocTypes []string) ([]string, error) {
	firstEquipmentWithLocation := equipment
	for {
		exist, err := firstEquipmentWithLocation.QueryLocation().Exist(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying location parent for equipment: %s, ID: %d", firstEquipmentWithLocation.Name, firstEquipmentWithLocation.ID)
		}
		if exist {
			break
		}
		// switch to parent equipment
		position, err := firstEquipmentWithLocation.QueryParentPosition().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "no location and equipment parent for equipment %s, ID: %d", firstEquipmentWithLocation.Name, firstEquipmentWithLocation.ID)
		}
		firstEquipmentWithLocation = position.QueryParent().OnlyX(ctx)
	}
	currLoc := firstEquipmentWithLocation.QueryLocation().OnlyX(ctx)
	return locationHierarchy(ctx, currLoc, orderedLocTypes)
}

func locationHierarchy(ctx context.Context, location *ent.Location, orderedLocTypes []string) ([]string, error) {
	var parents = make([]string, len(orderedLocTypes))
	currLoc := location
	for {
		typ, err := currLoc.QueryType().Only(ctx)
		if err != nil {
			return nil, errors.Errorf("getting location type for location : %s (id:%d)", currLoc.Name, currLoc.ID)
		}
		typeName := typ.Name
		idx := index(orderedLocTypes, typeName)
		if idx == -1 {
			return nil, errors.Errorf("location type does not exist: %s", typeName)
		}
		parents[idx] = currLoc.Name
		currLoc, err = currLoc.QueryParent().Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				break
			}
			return nil, errors.Wrapf(err, "querying parent location for location: %s", parents[idx])
		}
	}
	return parents, nil
}

func getLastLocations(ctx context.Context, e *ent.Equipment, level int) (*string, error) {
	ppos, err := e.QueryParentPosition().Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("querying parent position: %w", err)
	}

	var lastPosition *ent.EquipmentPosition
	for ppos != nil {
		lastPosition = ppos
		ppos, err = ppos.QueryParent().QueryParentPosition().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, fmt.Errorf("querying parent position: %w", err)
		}
	}
	var query *ent.LocationQuery
	if lastPosition != nil {
		query = lastPosition.QueryParent().QueryLocation()
	} else {
		query = e.QueryLocation()
	}
	loc, err := query.Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying equipemnt location: %w", err)
	}
	locations := loc.Name

	for i := 0; i < level-1; i++ {
		loc, err = loc.QueryParent().Only(ctx)
		if ent.MaskNotFound(err) != nil {
			return nil, fmt.Errorf("querying location for equipment: %w", err)
		}
		if ent.IsNotFound(err) || loc == nil {
			break
		}
		locations = loc.Name + "; " + locations
	}
	return &locations, nil
}

// nolint: funlen
func propertyTypesSlice(ctx context.Context, ids []int, c *ent.Client, entity models.PropertyEntity) ([]string, error) {
	var (
		propTypes       []string
		alreadyAppended = map[string]string{}
	)

	switch entity {
	case models.PropertyEntityEquipment:
		var equipTypesWithEquipment []ent.EquipmentType
		equipTypes, err := c.EquipmentType.Query().
			Paginate(ctx, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, typ := range equipTypes.Edges {
			equipType := typ.Node
			// TODO (T59268484) solve the case where there are too many IDs to check (trying to optimize)
			if len(ids) < 50 {
				switch exist, err := equipType.QueryEquipment().Where(equipment.IDIn(ids...)).Exist(ctx); {
				case err != nil:
					return nil, errors.Wrapf(err, "checking equipment instance existence for type: %s", equipType.Name)
				case exist:
					equipTypesWithEquipment = append(equipTypesWithEquipment, *equipType)
				}
			} else {
				equipTypesWithEquipment = append(equipTypesWithEquipment, *equipType)
			}
		}
		for _, equipType := range equipTypesWithEquipment {
			pts, err := equipType.QueryPropertyTypes().All(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "querying property types")
			}
			for _, ptype := range pts {
				if _, ok := alreadyAppended[ptype.Name]; !ok {
					alreadyAppended[ptype.Name] = ""
					propTypes = append(propTypes, ptype.Name)
				}
			}
		}
	case models.PropertyEntityLocation:
		var locTypesWithInstances []ent.LocationType
		locTypes, err := c.LocationType.Query().
			Paginate(ctx, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, typ := range locTypes.Edges {
			locType := typ.Node
			// TODO (T59268484) solve the case where there are too many IDs to check (trying to optimize)
			if len(ids) < 50 {
				switch exist, err := locType.QueryLocations().Where(location.IDIn(ids...)).Exist(ctx); {
				case err != nil:
					return nil, errors.Wrapf(err, "checking location instance existence for type: %s", locType.Name)
				case exist:
					locTypesWithInstances = append(locTypesWithInstances, *locType)
				}
			} else {
				locTypesWithInstances = append(locTypesWithInstances, *locType)
			}
		}
		for _, locType := range locTypesWithInstances {
			pts, err := locType.QueryPropertyTypes().All(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "querying property types")
			}
			for _, ptype := range pts {
				if _, ok := alreadyAppended[ptype.Name]; !ok {
					alreadyAppended[ptype.Name] = ""
					propTypes = append(propTypes, ptype.Name)
				}
			}
		}
	case models.PropertyEntityPort, models.PropertyEntityLink:
		var relevantPortTypes []ent.EquipmentPortType
		portTypes, err := c.EquipmentPortType.Query().
			Paginate(ctx, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, typ := range portTypes.Edges {
			portType := typ.Node
			if entity == models.PropertyEntityLink {
				// TODO (T59268484) solve the case where there are too many IDs to check (trying to optimize)
				if len(ids) < 50 {
					switch exist, err := portType.QueryPortDefinitions().QueryPorts().QueryLink().Where(link.IDIn(ids...)).Exist(ctx); {
					case err != nil:
						return nil, errors.Wrapf(err, "checking port instance existence for type: %s", portType.Name)
					case exist:
						relevantPortTypes = append(relevantPortTypes, *portType)
					}
				} else {
					relevantPortTypes = append(relevantPortTypes, *portType)
				}
			} else if entity == models.PropertyEntityPort {
				// TODO (T59268484) solve the case where there are too many IDs to check (trying to optimize)
				if len(ids) < 50 {
					switch exist, err := portType.QueryPortDefinitions().QueryPorts().Where(equipmentport.IDIn(ids...)).Exist(ctx); {
					case err != nil:
						return nil, errors.Wrapf(err, "checking port instance existence for type: %s", portType.Name)
					case exist:
						relevantPortTypes = append(relevantPortTypes, *portType)
					}
				} else {
					relevantPortTypes = append(relevantPortTypes, *portType)
				}
			}
		}
		for _, portType := range relevantPortTypes {
			var pts []*ent.PropertyType
			if entity == models.PropertyEntityPort {
				pts, err = portType.QueryPropertyTypes().All(ctx)
			} else if entity == models.PropertyEntityLink {
				pts, err = portType.QueryLinkPropertyTypes().All(ctx)
			}
			if err != nil {
				return nil, errors.Wrapf(err, "querying property types for %s", entity.String())
			}
			for _, pType := range pts {
				if _, ok := alreadyAppended[pType.Name]; !ok {
					alreadyAppended[pType.Name] = ""
					propTypes = append(propTypes, pType.Name)
				}
			}
		}
	case models.PropertyEntityService:
		var serviceTypesWithServices []ent.ServiceType
		serviceTypes, err := c.ServiceType.Query().
			Paginate(ctx, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, typ := range serviceTypes.Edges {
			serviceType := typ.Node
			// TODO (T59268484) solve the case where there are too many IDs to check (trying to optimize)
			if len(ids) < 50 {
				switch exist, err := serviceType.QueryServices().Where(service.IDIn(ids...)).Exist(ctx); {
				case err != nil:
					return nil, errors.Wrapf(err, "checking service instance existence for type: %s", serviceType.Name)
				case exist:
					serviceTypesWithServices = append(serviceTypesWithServices, *serviceType)
				}
			} else {
				serviceTypesWithServices = append(serviceTypesWithServices, *serviceType)
			}
		}
		for _, serviceType := range serviceTypesWithServices {
			pts, err := serviceType.QueryPropertyTypes().All(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "querying property types")
			}
			for _, ptype := range pts {
				if _, ok := alreadyAppended[ptype.Name]; !ok {
					alreadyAppended[ptype.Name] = ""
					propTypes = append(propTypes, ptype.Name)
				}
			}
		}
	case models.PropertyEntityWorkOrder:
		types, err := c.PropertyType.Query().
			Where(propertytype.HasPropertiesWith(property.HasWorkOrderWith(workorder.IDIn(ids...)))).
			GroupBy(propertytype.FieldName).Strings(ctx)
		if err != nil {
			return nil, err
		}
		return types, nil
	default:
		return nil, errors.Errorf("entity not supported %s", entity)
	}
	return propTypes, nil
}

// nolint: funlen
func propertiesSlice(ctx context.Context, instance interface{}, propertyTypes []string, entityType models.PropertyEntity) ([]string, error) {
	var ret = make([]string, len(propertyTypes))
	var typs []*ent.PropertyType
	var props []*ent.Property

	switch entityType {
	case models.PropertyEntityEquipment:
		entity := instance.(*ent.Equipment)
		var err error
		typs, err = entity.QueryType().QueryPropertyTypes().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying property types for equipment %s (id=%d)", entity.Name, entity.ID)
		}
		props, err = entity.QueryProperties().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying equipment properties (id=%d)", entity.ID)
		}
	case models.PropertyEntityPort:
		entity := instance.(*ent.EquipmentPort)
		var err error
		typs, err = entity.QueryDefinition().QueryEquipmentPortType().QueryPropertyTypes().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying property types for port (id=%d)", entity.ID)
		}
		props, err = entity.QueryProperties().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying port properties (id=%d)", entity.ID)
		}
	case models.PropertyEntityLink:
		entity := instance.(*ent.Link)
		ports, err := entity.QueryPorts().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying link for ports (id=%d)", entity.ID)
		}
		for _, port := range ports {
			var err error
			portTypeLinkProperties, err := port.QueryDefinition().QueryEquipmentPortType().QueryLinkPropertyTypes().All(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "querying property types for port (id=%d)", entity.ID)
			}
			typs = append(typs, portTypeLinkProperties...)
		}
		props, err = entity.QueryProperties().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying link properties (id=%d)", entity.ID)
		}
	case models.PropertyEntityService:
		entity := instance.(*ent.Service)
		var err error
		typs, err = entity.QueryType().QueryPropertyTypes().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying property types for service (id=%d)", entity.ID)
		}
		props, err = entity.QueryProperties().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying services properties (id=%d)", entity.ID)
		}
	case models.PropertyEntityLocation:
		entity := instance.(*ent.Location)
		var err error
		typs, err = entity.QueryType().QueryPropertyTypes().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "can't query property types for location (id=%d)", entity.ID)
		}
		props, err = entity.QueryProperties().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying location properties (id=%d)", entity.ID)
		}
	case models.PropertyEntityWorkOrder:
		entity := instance.(*ent.WorkOrder)
		var err error
		props, err = entity.QueryProperties().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying property types for work order %s (id=%d)", entity.Name, entity.ID)
		}
	default:
		return nil, errors.Errorf("entityType not supported %s", entityType)
	}

	for _, typ := range typs {
		idx := index(propertyTypes, typ.Name)
		if idx == -1 {
			continue
		}
		val, err := resolverutil.PropertyValue(ctx, typ.Type, typ)
		if err != nil {
			return nil, err
		}
		ret[idx] = val
	}

	for _, p := range props {
		propTyp, err := p.QueryType().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying type of property (id=%d)", p.ID)
		}
		propTypeName := propTyp.Name
		idx := index(propertyTypes, propTypeName)
		if idx == -1 {
			return nil, errors.Errorf("Property type does not exist in header: %s", propTypeName)
		}
		typ := propTyp.Type
		val, err := resolverutil.PropertyValue(ctx, typ, p)
		if err != nil {
			return nil, err
		}
		ret[idx] = val
	}
	return ret, nil
}
