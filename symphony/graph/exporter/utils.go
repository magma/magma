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

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/pkg/errors"
)

const (
	maxEquipmentParents = 3
	boolVal             = "bool"
	emailVal            = "email"
	stringVal           = "string"
	dateVal             = "date"
	intVal              = "int"
	floatVal            = "float"
	gpsLocationVal      = "gps_location"
	rangeVal            = "range"
	enum                = "enum"
	equipmentVal        = "equipment"
	locationVal         = "location"
)

func index(a []string, x string) int {
	for i, n := range a {
		if strings.EqualFold(x, n) {
			return i
		}
	}
	return -1
}

func locationTypeHierarchy(ctx context.Context, c *ent.Client) ([]string, error) {
	locTypeResult, err := resolverutil.LocationTypes(ctx, c)
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

func locationHierarchy(ctx context.Context, equipment *ent.Equipment, orderedLocTypes []string) ([]string, error) {
	var parents = make([]string, len(orderedLocTypes))
	firstEquipmentWithLocation := equipment
	var err error
	for {
		exist, err := firstEquipmentWithLocation.QueryLocation().Exist(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying location parent for equipment: %s, ID: %s", firstEquipmentWithLocation.Name, firstEquipmentWithLocation.ID)
		}
		if exist {
			break
		}
		// switch to parent equipment
		position, err := firstEquipmentWithLocation.QueryParentPosition().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "no location and equipment parent for equipment %s, ID: %s", firstEquipmentWithLocation.Name, firstEquipmentWithLocation.ID)
		}
		firstEquipmentWithLocation = position.QueryParent().OnlyX(ctx)
	}
	currLoc := firstEquipmentWithLocation.QueryLocation().OnlyX(ctx)
	for {
		typeName := currLoc.QueryType().OnlyX(ctx).Name
		idx := index(orderedLocTypes, typeName)
		if idx == -1 {
			return nil, errors.Errorf("Location  type does not exist : %s", typeName)
		}
		parents[idx] = currLoc.Name
		currLoc, err = currLoc.QueryParent().Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				break
			}
			return nil, errors.Wrapf(err, "error querying parent location for location: %s", currLoc.Name)
		}
	}
	return parents, nil
}

func propertyTypesSlice(ctx context.Context, ids []string, c *ent.Client, entity models.PropertyEntity) ([]string, error) {
	var (
		propTypes       []string
		alreadyAppended = map[string]string{}
	)

	if entity == models.PropertyEntityEquipment {
		var equipTypesWithEquipment []ent.EquipmentType
		equipTypes, err := resolverutil.EquipmentTypes(ctx, c)
		if err != nil {
			return nil, err
		}

		for _, typ := range equipTypes.Edges {
			equipType := typ.Node
			if equipType.QueryEquipment().Where(equipment.IDIn(ids...)).ExistX(ctx) {
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
	}
	if entity == models.PropertyEntityPort {
		var portTypesWithPorts []ent.EquipmentPortType
		portTypes, err := resolverutil.EquipmentPortTypes(ctx, c)
		if err != nil {
			return nil, err
		}

		for _, typ := range portTypes.Edges {
			portType := typ.Node
			if portType.QueryPortDefinitions().QueryPorts().Where(equipmentport.IDIn(ids...)).ExistX(ctx) {
				portTypesWithPorts = append(portTypesWithPorts, *portType)
			}
		}
		for _, portType := range portTypesWithPorts {
			pts, err := portType.QueryPropertyTypes().All(ctx)
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
	}
	return propTypes, nil
}

func propertiesSlice(ctx context.Context, instance interface{}, propertyTypes []string, entityType models.PropertyEntity) ([]string, error) {
	var ret = make([]string, len(propertyTypes))
	var typs []*ent.PropertyType
	var props []*ent.Property

	if entityType == models.PropertyEntityEquipment {
		entity := instance.(*ent.Equipment)
		var err error
		typs, err = entity.QueryType().QueryPropertyTypes().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "can't query property types for equipment %s (id=%s)", entity.Name, entity.ID)
		}
		props = entity.QueryProperties().AllX(ctx)
	} else {
		entity := instance.(*ent.EquipmentPort)
		var err error
		typs, err = entity.QueryDefinition().QueryEquipmentPortType().QueryPropertyTypes().All(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "can't query property types for port (id=%s)", entity.ID)
		}
		props = entity.QueryProperties().AllX(ctx)
	}

	for _, typ := range typs {
		idx := index(propertyTypes, typ.Name)
		val, err := propertyValue(ctx, typ.Type, typ)
		if err != nil {
			return nil, err
		}
		ret[idx] = val
	}

	for _, p := range props {
		propTypeName := p.QueryType().OnlyX(ctx).Name
		idx := index(propertyTypes, propTypeName)
		if idx == -1 {
			return nil, errors.Errorf("Property type does not exist : %s", propTypeName)
		}
		typ := p.QueryType().OnlyX(ctx).Type
		val, err := propertyValue(ctx, typ, p)
		if err != nil {
			return nil, err
		}
		ret[idx] = val
	}
	return ret, nil
}

func propertyValue(ctx context.Context, typ string, v interface{}) (string, error) {
	switch v.(type) {
	case *ent.PropertyType, *ent.Property:
	default:
		return "", errors.Errorf("invalid type: %T", v)
	}
	vo := reflect.ValueOf(v).Elem()
	switch typ {
	case emailVal, stringVal, dateVal, enum:
		return vo.FieldByName("StringVal").String(), nil
	case intVal:
		i := vo.FieldByName("IntVal").Int()
		return strconv.Itoa(int(i)), nil
	case floatVal:
		return fmt.Sprintf("%.3f", vo.FieldByName("FloatVal").Float()), nil
	case gpsLocationVal:
		la, lo := vo.FieldByName("LatitudeVal").Float(), vo.FieldByName("LongitudeVal").Float()
		return fmt.Sprintf("%f", la) + ", " + fmt.Sprintf("%f", lo), nil
	case rangeVal:
		rf, rt := vo.FieldByName("RangeFromVal").Float(), vo.FieldByName("RangeToVal").Float()
		return fmt.Sprintf("%.3f", rf) + " - " + fmt.Sprintf("%.3f", rt), nil
	case boolVal:
		return strconv.FormatBool(vo.FieldByName("BoolVal").Bool()), nil
	case equipmentVal, locationVal:
		property, ok := v.(*ent.Property)
		if !ok {
			return "", nil
		}
		var id string
		if typ == equipmentVal {
			id, _ = property.QueryEquipmentValue().OnlyID(ctx)
		} else {
			id, _ = property.QueryLocationValue().OnlyID(ctx)
		}
		return id, nil
	default:
		return "", errors.Errorf("type not supported %s", typ)
	}
}
