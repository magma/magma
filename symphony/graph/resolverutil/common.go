// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

type AddPropertyArgs struct {
	Context    context.Context
	EntSetter  func(*ent.PropertyCreate)
	IsTemplate *bool
}

// GetPropertyPredicate returns the property predicate for the filter
func GetPropertyPredicate(p models.PropertyTypeInput) (predicate.Property, error) {
	var pred predicate.Property
	switch p.Type {
	case models.PropertyKindString,
		models.PropertyKindEmail,
		models.PropertyKindDate,
		models.PropertyKindEnum,
		models.PropertyKindDatetimeLocal:
		if p.StringValue != nil {
			pred = property.StringVal(*p.StringValue)
		}
	case models.PropertyKindInt:
		if p.IntValue != nil {
			pred = property.IntVal(*p.IntValue)
		}
	case models.PropertyKindBool:
		if p.BooleanValue != nil {
			pred = property.BoolVal(*p.BooleanValue)
		}
	case models.PropertyKindFloat:
		if p.FloatValue != nil {
			pred = property.FloatVal(*p.FloatValue)
		}
	default:
		return nil, errors.Errorf("operator not supported for kind %q", p.Type)
	}
	return pred, nil
}

// GetPropertyTypePredicate returns the propertyType predicate for the filter
func GetPropertyTypePredicate(p models.PropertyTypeInput) (predicate.PropertyType, error) {
	var pred predicate.PropertyType
	switch p.Type {
	case models.PropertyKindString,
		models.PropertyKindEmail,
		models.PropertyKindDate,
		models.PropertyKindEnum,
		models.PropertyKindDatetimeLocal:
		if p.StringValue != nil {
			pred = propertytype.StringVal(*p.StringValue)
		}
	case models.PropertyKindInt:
		if p.IntValue != nil {
			pred = propertytype.IntVal(*p.IntValue)
		}
	case models.PropertyKindBool:
		if p.BooleanValue != nil {
			pred = propertytype.BoolVal(*p.BooleanValue)
		}
	case models.PropertyKindFloat:
		if p.FloatValue != nil {
			pred = propertytype.FloatVal(*p.FloatValue)
		}
	default:
		return nil, errors.Errorf("operator not supported for kind %q", p.Type)
	}
	return pred, nil
}

// GetDatePropertyPred returns the property and propertyType predicate for the date
func GetDatePropertyPred(p models.PropertyTypeInput, operator models.FilterOperator) (predicate.Property, predicate.PropertyType, error) {
	if p.Type != models.PropertyKindDate && p.Type != models.PropertyKindDatetimeLocal {
		return nil, nil, errors.Errorf("property kind should be type")
	}

	propPred := property.StringValGT(*p.StringValue)
	propTypePred := propertytype.StringValGT(*p.StringValue)
	if operator == models.FilterOperatorDateLessThan {
		propPred = property.StringValLT(*p.StringValue)
		propTypePred = propertytype.StringValLT(*p.StringValue)
	}
	return propPred, propTypePred, nil
}
