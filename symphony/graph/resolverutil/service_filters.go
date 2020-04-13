// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

func handleServiceFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	switch filter.FilterType {
	case models.ServiceFilterTypeServiceInstName:
		return serviceNameFilter(q, filter)
	case models.ServiceFilterTypeServiceStatus:
		return serviceStatusFilter(q, filter)
	case models.ServiceFilterTypeServiceType:
		return serviceTypeFilter(q, filter)
	case models.ServiceFilterTypeServiceInstExternalID:
		return externalIDFilter(q, filter)
	case models.ServiceFilterTypeServiceInstCustomerName:
		return customerNameFilter(q, filter)
	default:
		return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
	}
}

func serviceNameFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.Operator == models.FilterOperatorContains {
		return q.Where(service.NameContainsFold(*filter.StringValue)), nil
	}
	return nil, errors.Errorf("operation %q not supported", filter.Operator)
}

func serviceStatusFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(service.StatusIn(filter.StringSet...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func serviceTypeFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		return q.Where(service.HasTypeWith(servicetype.IDIn(filter.IDSet...))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func externalIDFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.Operator == models.FilterOperatorIs {
		return q.Where(service.ExternalID(*filter.StringValue)), nil
	}
	return nil, errors.Errorf("operation %q not supported", filter.Operator)
}

func customerNameFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.Operator == models.FilterOperatorContains {
		return q.Where(service.HasCustomerWith(customer.NameContainsFold(*filter.StringValue))), nil
	}
	return nil, errors.Errorf("operation %q not supported", filter.Operator)
}

func handleServicePropertyFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.FilterType == models.ServiceFilterTypeProperty {
		return servicePropertyFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func servicePropertyFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
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
			service.Or(
				service.HasPropertiesWith(
					property.And(
						property.HasTypeWith(
							propertytype.Name(p.Name),
							propertytype.Type(p.Type.String()),
						),
						pred,
					),
				),
				service.And(
					service.HasTypeWith(servicetype.HasPropertyTypesWith(
						propertytype.Name(p.Name),
						propertytype.Type(p.Type.String()),
						predForType,
					)),
					service.Not(service.HasPropertiesWith(
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

func handleServiceLocationFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.FilterType == models.ServiceFilterTypeLocationInst {
		return serviceLocationFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func serviceLocationFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.Operator == models.FilterOperatorIsOneOf {
		var ps []predicate.Service
		for _, lid := range filter.IDSet {
			eqPred := BuildGeneralEquipmentAncestorFilter(
				equipment.HasLocationWith(BuildLocationAncestorFilter(lid, 1, *filter.MaxDepth)),
				1,
				*filter.MaxDepth)
			ps = append(ps, service.HasEndpointsWith(
				serviceendpoint.HasPortWith(equipmentport.HasParentWith(eqPred)),
			))
		}
		return q.Where(service.Or(ps...)), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}

func handleEquipmentInServiceFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.FilterType == models.ServiceFilterTypeEquipmentInService {
		return equipmentInServiceTypeFilter(q, filter)
	}
	return nil, errors.Errorf("filter type is not supported: %s", filter.FilterType)
}

func equipmentInServiceTypeFilter(q *ent.ServiceQuery, filter *models.ServiceFilterInput) (*ent.ServiceQuery, error) {
	if filter.Operator == models.FilterOperatorContains {
		equipmentNameQuery := equipment.NameContainsFold(*filter.StringValue)
		return q.Where(
			service.Or(service.HasLinksWith(
				link.HasPortsWith(equipmentport.HasParentWith(equipmentNameQuery))),
				service.HasEndpointsWith(serviceendpoint.HasPortWith(equipmentport.HasParentWith(equipmentNameQuery))))), nil
	}
	return nil, errors.Errorf("operation is not supported: %s", filter.Operator)
}
