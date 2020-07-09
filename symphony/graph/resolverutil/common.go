// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/pkg/errors"
)

const (
	NodeTypeLocation  = "location"
	NodeTypeEquipment = "equipment"
	NodeTypeService   = "service"
	NodeTypeWorkOrder = "work_order"
	NodeTypeUser      = "user"
)

type AddPropertyArgs struct {
	Context    context.Context
	EntSetter  func(*ent.PropertyCreate)
	IsTemplate *bool
}

func NodePropertyValue(ctx context.Context, p *ent.Property, nodeType string) string {
	var id *int
	switch nodeType {
	case NodeTypeLocation:
		if i, err := p.QueryLocationValue().OnlyID(ctx); err == nil {
			id = &i
		}
	case NodeTypeEquipment:
		if i, err := p.QueryEquipmentValue().OnlyID(ctx); err == nil {
			id = &i
		}
	case NodeTypeService:
		if i, err := p.QueryServiceValue().OnlyID(ctx); err == nil {
			id = &i
		}
	case NodeTypeWorkOrder:
		if i, err := p.QueryWorkOrderValue().OnlyID(ctx); err == nil {
			id = &i
		}
	case NodeTypeUser:
		if i, err := p.QueryUserValue().OnlyID(ctx); err == nil {
			id = &i
		}
	default:
		return ""
	}
	if id == nil {
		return ""
	}
	return strconv.Itoa(*id)
}

func PropertyValue(ctx context.Context, typ propertytype.Type, nodeType string, v interface{}) (string, error) {
	switch v.(type) {
	case *ent.PropertyType, *ent.Property:
	default:
		return "", errors.Errorf("invalid type: %T", v)
	}
	vo := reflect.ValueOf(v).Elem()
	switch typ {
	case propertytype.TypeEmail, propertytype.TypeString, propertytype.TypeDate,
		propertytype.TypeEnum, propertytype.TypeDatetimeLocal:
		strValue := vo.FieldByName("StringVal")
		if strValue.IsNil() {
			return "", nil
		}
		return reflect.Indirect(strValue).String(), nil
	case propertytype.TypeInt:
		intValue := vo.FieldByName("IntVal")
		if intValue.IsNil() {
			return "", nil
		}
		return strconv.Itoa(int(reflect.Indirect(intValue).Int())), nil
	case propertytype.TypeFloat:
		floatValue := vo.FieldByName("FloatVal")
		if floatValue.IsNil() {
			return "", nil
		}
		return fmt.Sprintf("%.3f", reflect.Indirect(floatValue).Float()), nil
	case propertytype.TypeGpsLocation:
		latitudeValue := vo.FieldByName("LatitudeVal")
		longitudeValue := vo.FieldByName("LongitudeVal")
		if latitudeValue.IsNil() || longitudeValue.IsNil() {
			return "", nil
		}
		la, lo := reflect.Indirect(latitudeValue).Float(), reflect.Indirect(longitudeValue).Float()
		return fmt.Sprintf("%f", la) + ", " + fmt.Sprintf("%f", lo), nil
	case propertytype.TypeRange:
		rangeFromValue := vo.FieldByName("RangeFromVal")
		rangeToValue := vo.FieldByName("RangeToVal")
		if rangeFromValue.IsNil() || rangeToValue.IsNil() {
			return "", nil
		}
		rf, rt := reflect.Indirect(rangeFromValue).Float(), reflect.Indirect(rangeToValue).Float()
		return fmt.Sprintf("%.3f", rf) + " - " + fmt.Sprintf("%.3f", rt), nil
	case propertytype.TypeBool:
		boolValue := vo.FieldByName("BoolVal")
		if boolValue.IsNil() {
			return "", nil
		}
		return strconv.FormatBool(reflect.Indirect(boolValue).Bool()), nil
	case propertytype.TypeNode:
		p, ok := v.(*ent.Property)
		if !ok {
			return "", nil
		}
		return NodePropertyValue(ctx, p, nodeType), nil
	default:
		return "", errors.Errorf("type not supported %s", typ)
	}
}

// GetPropertyPredicate returns the property predicate for the filter
func GetPropertyPredicate(p models.PropertyTypeInput) (predicate.Property, error) {
	var pred predicate.Property
	switch p.Type {
	case propertytype.TypeString,
		propertytype.TypeEmail,
		propertytype.TypeDate,
		propertytype.TypeEnum,
		propertytype.TypeDatetimeLocal:
		if p.StringValue != nil {
			pred = property.StringVal(*p.StringValue)
		}
	case propertytype.TypeInt:
		if p.IntValue != nil {
			pred = property.IntVal(*p.IntValue)
		}
	case propertytype.TypeBool:
		if p.BooleanValue != nil {
			pred = property.BoolVal(*p.BooleanValue)
		}
	case propertytype.TypeFloat:
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
	case propertytype.TypeString,
		propertytype.TypeEmail,
		propertytype.TypeDate,
		propertytype.TypeEnum,
		propertytype.TypeDatetimeLocal:
		if p.StringValue != nil {
			pred = propertytype.StringVal(*p.StringValue)
		}
	case propertytype.TypeInt:
		if p.IntValue != nil {
			pred = propertytype.IntVal(*p.IntValue)
		}
	case propertytype.TypeBool:
		if p.BooleanValue != nil {
			pred = propertytype.BoolVal(*p.BooleanValue)
		}
	case propertytype.TypeFloat:
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
	if p.Type != propertytype.TypeDate && p.Type != propertytype.TypeDatetimeLocal {
		return nil, nil, errors.Errorf("property kind should be type")
	}
	if operator == models.FilterOperatorDateLessThan {
		return property.StringValLT(*p.StringValue), propertytype.StringValLT(*p.StringValue), nil
	}
	return property.StringValGT(*p.StringValue), propertytype.StringValGT(*p.StringValue), nil
}
