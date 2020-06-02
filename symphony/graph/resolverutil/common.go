// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/facebookincubator/symphony/pkg/ent/user"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/pkg/errors"
)

const (
	boolVal          = "bool"
	emailVal         = "email"
	stringVal        = "string"
	dateVal          = "date"
	intVal           = "int"
	floatVal         = "float"
	gpsLocationVal   = "gps_location"
	rangeVal         = "range"
	enum             = "enum"
	datetimeLocalVal = "datetime_local"
	nodeVal          = "node"
)

type AddPropertyArgs struct {
	Context    context.Context
	EntSetter  func(*ent.PropertyCreate)
	IsTemplate *bool
}

func PropertyValue(ctx context.Context, typ string, v interface{}) (string, error) {
	switch v.(type) {
	case *ent.PropertyType, *ent.Property:
	default:
		return "", errors.Errorf("invalid type: %T", v)
	}
	vo := reflect.ValueOf(v).Elem()
	switch typ {
	case emailVal, stringVal, dateVal, enum, datetimeLocalVal:
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
	case nodeVal:
		p, ok := v.(*ent.Property)
		if !ok {
			return "", nil
		}
		var id int
		if i, err := p.QueryEquipmentValue().OnlyID(ctx); err == nil {
			id = i
		}
		if i, err := p.QueryLocationValue().OnlyID(ctx); err == nil {
			id = i
		}
		if i, err := p.QueryServiceValue().OnlyID(ctx); err == nil {
			id = i
		}
		if i, err := p.QueryWorkOrderValue().OnlyID(ctx); err == nil {
			id = i
		}
		if i, err := p.QueryUserValue().OnlyID(ctx); err == nil {
			id = i
		}
		return strconv.Itoa(id), nil
	default:
		return "", errors.Errorf("type not supported %s", typ)
	}
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
	if operator == models.FilterOperatorDateLessThan {
		return property.StringValLT(*p.StringValue), propertytype.StringValLT(*p.StringValue), nil
	}
	return property.StringValGT(*p.StringValue), propertytype.StringValGT(*p.StringValue), nil
}

func GetUserID(ctx context.Context, userID *int, userName *string) (*int, error) {
	if userName != nil && *userName != "" {
		c := ent.FromContext(ctx)
		id, err := c.User.Query().Where(user.AuthID(*userName)).OnlyID(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching assignee user: %w", err)
		}
		return &id, nil
	}
	return userID, nil
}
