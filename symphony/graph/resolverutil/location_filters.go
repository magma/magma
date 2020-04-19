// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

func handleLocationFilter(q *ent.LocationQuery, filter *models.LocationFilterInput) (*ent.LocationQuery, error) {
	switch filter.FilterType {
	case models.LocationFilterTypeLocationInst:
		return LocationFilterPredicate(q, filter)
	case models.LocationFilterTypeLocationInstHasEquipment:
		return locationHasEquipmentFilter(q, filter)
	case models.LocationFilterTypeLocationInstName:
		return locationNameFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func locationNameFilter(q *ent.LocationQuery, filter *models.LocationFilterInput) (*ent.LocationQuery, error) {
	if filter.Operator == models.FilterOperatorIs {
		return q.Where(location.NameEqualFold(*filter.StringValue)), nil
	}
	return nil, errors.Errorf("operation %s is not supported", filter.Operator)
}

func locationHasEquipmentFilter(q *ent.LocationQuery, filter *models.LocationFilterInput) (*ent.LocationQuery, error) {
	if filter.Operator == models.FilterOperatorIs {
		var pp predicate.Location
		if *filter.BoolValue {
			pp = location.HasEquipment()
		} else {
			pp = location.Not(location.HasEquipment())
		}
		return q.Where(pp), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func handleLocationTypeFilter(q *ent.LocationQuery, filter *models.LocationFilterInput) (*ent.LocationQuery, error) {
	if filter.FilterType == models.LocationFilterTypeLocationType {
		return locationLocationTypeFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func locationLocationTypeFilter(q *ent.LocationQuery, filter *models.LocationFilterInput) (*ent.LocationQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(location.HasTypeWith(locationtype.IDIn(filter.IDSet...))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

// nolint: dupl
func handleLocationPropertyFilter(q *ent.LocationQuery, filter *models.LocationFilterInput) (*ent.LocationQuery, error) {
	p := filter.PropertyValue
	switch filter.Operator {
	case models.FilterOperatorIs:
		pred, err := GetPropertyPredicate(*p)
		if err != nil {
			return nil, err
		}
		typePred, err := GetPropertyTypePredicate(*p)
		if err != nil {
			return nil, err
		}

		q = q.Where(location.Or(
			location.HasPropertiesWith(
				property.And(
					property.HasTypeWith(
						propertytype.Name(p.Name),
						propertytype.Type(p.Type.String()),
					),
					pred,
				),
			),
			location.And(
				location.HasTypeWith(locationtype.HasPropertyTypesWith(
					propertytype.Name(p.Name),
					propertytype.Type(p.Type.String()),
					typePred,
				)),
				location.Not(location.HasPropertiesWith(
					property.HasTypeWith(
						propertytype.Name(p.Name),
						propertytype.Type(p.Type.String()),
					)),
				))))

		return q, nil
	case models.FilterOperatorDateLessThan, models.FilterOperatorDateGreaterThan:
		propPred, propTypePred, err := GetDatePropertyPred(*p, filter.Operator)
		if err != nil {
			return nil, err
		}
		q = q.Where(location.Or(
			location.HasPropertiesWith(
				property.And(
					property.HasTypeWith(
						propertytype.Name(p.Name),
						propertytype.Type(p.Type.String()),
					),
					propPred,
				),
			),
			location.And(
				location.HasTypeWith(locationtype.HasPropertyTypesWith(
					propertytype.Name(p.Name),
					propertytype.Type(p.Type.String()),
					propTypePred,
				)),
				location.Not(location.HasPropertiesWith(
					property.HasTypeWith(
						propertytype.Name(p.Name),
						propertytype.Type(p.Type.String()),
					)),
				))))
		return q, nil
	default:
		return nil, errors.Errorf("operator %q not supported", filter.Operator)
	}

}
