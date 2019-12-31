// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

func handlePortFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	if filter.FilterType == models.PortFilterTypePortInstEquipment {
		return portEquipmentFilter(q, filter)
	}
	if filter.FilterType == models.PortFilterTypePortInstHasLink {
		return portHasLinkFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func portEquipmentFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	switch filter.Operator {
	case models.FilterOperatorContains:
		return q.Where(equipmentport.HasParentWith(equipment.NameContainsFold(*filter.StringValue))), nil
	case models.FilterOperatorIsOneOf:
		return q.Where(equipmentport.HasParentWith(equipment.IDIn(filter.IDSet...))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func portHasLinkFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	if filter.Operator == models.FilterOperatorIs {
		var pp predicate.EquipmentPort
		if *filter.BoolValue {
			pp = equipmentport.HasLink()
		} else {
			pp = equipmentport.Not(equipmentport.HasLink())
		}
		return q.Where(pp), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func handlePortLocationFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	if filter.FilterType == models.PortFilterTypeLocationInst {
		return portLocationFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func portLocationFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		var pp []predicate.EquipmentPort

		for _, lid := range filter.IDSet {
			pp = append(pp, GetPortLocationPredicate(lid, filter.MaxDepth))
		}
		return q.Where(equipmentport.Or(pp...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func handlePortDefinitionFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	if filter.FilterType == models.PortFilterTypePortDef {
		return portDefFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func portDefFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(equipmentport.HasDefinitionWith(equipmentportdefinition.IDIn(filter.IDSet...))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func handlePortPropertyFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	p := filter.PropertyValue
	switch filter.Operator {
	case models.FilterOperatorIs:
		pred, err := GetPropertyPredicate(*p)
		if err != nil {
			return nil, err
		}
		predForType, err := GetPropertyTypePredicate(*p)
		if err != nil {
			return nil, err
		}

		q = q.Where(
			equipmentport.Or(
				equipmentport.HasPropertiesWith(
					property.And(
						property.HasTypeWith(
							propertytype.Name(p.Name),
							propertytype.Type(p.Type.String()),
						),
						pred,
					),
				),
				equipmentport.And(
					equipmentport.HasDefinitionWith(equipmentportdefinition.HasEquipmentPortTypeWith(
						equipmentporttype.HasPropertyTypesWith(
							propertytype.Name(p.Name),
							propertytype.Type(p.Type.String()),
							predForType,
						))),
					equipmentport.Not(equipmentport.HasPropertiesWith(
						property.HasTypeWith(
							propertytype.Name(p.Name),
							propertytype.Type(p.Type.String()),
						)),
					),
				),
			),
		)
		return q, nil
	case models.FilterOperatorDateLessThan, models.FilterOperatorDateGreaterThan:
		propPred, propTypePred, err := GetDatePropertyPred(*p, filter.Operator)
		if err != nil {
			return nil, err
		}
		q = q.Where(equipmentport.Or(
			equipmentport.HasPropertiesWith(
				property.And(
					property.HasTypeWith(
						propertytype.Name(p.Name),
						propertytype.Type(p.Type.String()),
					),
					propPred,
				),
			),
			equipmentport.And(
				equipmentport.HasDefinitionWith(equipmentportdefinition.HasEquipmentPortTypeWith(equipmentporttype.HasPropertyTypesWith(
					propertytype.Name(p.Name),
					propertytype.Type(p.Type.String()),
					propTypePred,
				))),
				equipmentport.Not(equipmentport.HasPropertiesWith(
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

func handlePortServiceFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	if filter.FilterType == models.PortFilterTypeServiceInst {
		return portServiceFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func portServiceFilter(q *ent.EquipmentPortQuery, filter *models.PortFilterInput) (*ent.EquipmentPortQuery, error) {
	// TODO: add the query
	return q, nil
}
