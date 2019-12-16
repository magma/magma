// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// PropertyType is the model entity for the PropertyType schema.
type PropertyType struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// Category holds the value of the "category" field.
	Category string `json:"category,omitempty"`
	// IntVal holds the value of the "int_val" field.
	IntVal int `json:"int_val,omitempty" gqlgen:"intValue"`
	// BoolVal holds the value of the "bool_val" field.
	BoolVal bool `json:"bool_val,omitempty" gqlgen:"booleanValue"`
	// FloatVal holds the value of the "float_val" field.
	FloatVal float64 `json:"float_val,omitempty" gqlgen:"floatValue"`
	// LatitudeVal holds the value of the "latitude_val" field.
	LatitudeVal float64 `json:"latitude_val,omitempty" gqlgen:"latitudeValue"`
	// LongitudeVal holds the value of the "longitude_val" field.
	LongitudeVal float64 `json:"longitude_val,omitempty" gqlgen:"longitudeValue"`
	// StringVal holds the value of the "string_val" field.
	StringVal string `json:"string_val,omitempty" gqlgen:"stringValue"`
	// RangeFromVal holds the value of the "range_from_val" field.
	RangeFromVal float64 `json:"range_from_val,omitempty" gqlgen:"rangeFromValue"`
	// RangeToVal holds the value of the "range_to_val" field.
	RangeToVal float64 `json:"range_to_val,omitempty" gqlgen:"rangeToValue"`
	// IsInstanceProperty holds the value of the "is_instance_property" field.
	IsInstanceProperty bool `json:"is_instance_property,omitempty" gqlgen:"isInstanceProperty"`
	// Editable holds the value of the "editable" field.
	Editable bool `json:"editable,omitempty" gqlgen:"isEditable"`
	// Mandatory holds the value of the "mandatory" field.
	Mandatory bool `json:"mandatory,omitempty" gqlgen:"isMandatory"`
}

// FromRows scans the sql response data into PropertyType.
func (pt *PropertyType) FromRows(rows *sql.Rows) error {
	var scanpt struct {
		ID                 int
		CreateTime         sql.NullTime
		UpdateTime         sql.NullTime
		Type               sql.NullString
		Name               sql.NullString
		Index              sql.NullInt64
		Category           sql.NullString
		IntVal             sql.NullInt64
		BoolVal            sql.NullBool
		FloatVal           sql.NullFloat64
		LatitudeVal        sql.NullFloat64
		LongitudeVal       sql.NullFloat64
		StringVal          sql.NullString
		RangeFromVal       sql.NullFloat64
		RangeToVal         sql.NullFloat64
		IsInstanceProperty sql.NullBool
		Editable           sql.NullBool
		Mandatory          sql.NullBool
	}
	// the order here should be the same as in the `propertytype.Columns`.
	if err := rows.Scan(
		&scanpt.ID,
		&scanpt.CreateTime,
		&scanpt.UpdateTime,
		&scanpt.Type,
		&scanpt.Name,
		&scanpt.Index,
		&scanpt.Category,
		&scanpt.IntVal,
		&scanpt.BoolVal,
		&scanpt.FloatVal,
		&scanpt.LatitudeVal,
		&scanpt.LongitudeVal,
		&scanpt.StringVal,
		&scanpt.RangeFromVal,
		&scanpt.RangeToVal,
		&scanpt.IsInstanceProperty,
		&scanpt.Editable,
		&scanpt.Mandatory,
	); err != nil {
		return err
	}
	pt.ID = strconv.Itoa(scanpt.ID)
	pt.CreateTime = scanpt.CreateTime.Time
	pt.UpdateTime = scanpt.UpdateTime.Time
	pt.Type = scanpt.Type.String
	pt.Name = scanpt.Name.String
	pt.Index = int(scanpt.Index.Int64)
	pt.Category = scanpt.Category.String
	pt.IntVal = int(scanpt.IntVal.Int64)
	pt.BoolVal = scanpt.BoolVal.Bool
	pt.FloatVal = scanpt.FloatVal.Float64
	pt.LatitudeVal = scanpt.LatitudeVal.Float64
	pt.LongitudeVal = scanpt.LongitudeVal.Float64
	pt.StringVal = scanpt.StringVal.String
	pt.RangeFromVal = scanpt.RangeFromVal.Float64
	pt.RangeToVal = scanpt.RangeToVal.Float64
	pt.IsInstanceProperty = scanpt.IsInstanceProperty.Bool
	pt.Editable = scanpt.Editable.Bool
	pt.Mandatory = scanpt.Mandatory.Bool
	return nil
}

// QueryProperties queries the properties edge of the PropertyType.
func (pt *PropertyType) QueryProperties() *PropertyQuery {
	return (&PropertyTypeClient{pt.config}).QueryProperties(pt)
}

// QueryLocationType queries the location_type edge of the PropertyType.
func (pt *PropertyType) QueryLocationType() *LocationTypeQuery {
	return (&PropertyTypeClient{pt.config}).QueryLocationType(pt)
}

// QueryEquipmentPortType queries the equipment_port_type edge of the PropertyType.
func (pt *PropertyType) QueryEquipmentPortType() *EquipmentPortTypeQuery {
	return (&PropertyTypeClient{pt.config}).QueryEquipmentPortType(pt)
}

// QueryLinkEquipmentPortType queries the link_equipment_port_type edge of the PropertyType.
func (pt *PropertyType) QueryLinkEquipmentPortType() *EquipmentPortTypeQuery {
	return (&PropertyTypeClient{pt.config}).QueryLinkEquipmentPortType(pt)
}

// QueryEquipmentType queries the equipment_type edge of the PropertyType.
func (pt *PropertyType) QueryEquipmentType() *EquipmentTypeQuery {
	return (&PropertyTypeClient{pt.config}).QueryEquipmentType(pt)
}

// QueryServiceType queries the service_type edge of the PropertyType.
func (pt *PropertyType) QueryServiceType() *ServiceTypeQuery {
	return (&PropertyTypeClient{pt.config}).QueryServiceType(pt)
}

// QueryWorkOrderType queries the work_order_type edge of the PropertyType.
func (pt *PropertyType) QueryWorkOrderType() *WorkOrderTypeQuery {
	return (&PropertyTypeClient{pt.config}).QueryWorkOrderType(pt)
}

// QueryProjectType queries the project_type edge of the PropertyType.
func (pt *PropertyType) QueryProjectType() *ProjectTypeQuery {
	return (&PropertyTypeClient{pt.config}).QueryProjectType(pt)
}

// Update returns a builder for updating this PropertyType.
// Note that, you need to call PropertyType.Unwrap() before calling this method, if this PropertyType
// was returned from a transaction, and the transaction was committed or rolled back.
func (pt *PropertyType) Update() *PropertyTypeUpdateOne {
	return (&PropertyTypeClient{pt.config}).UpdateOne(pt)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (pt *PropertyType) Unwrap() *PropertyType {
	tx, ok := pt.config.driver.(*txDriver)
	if !ok {
		panic("ent: PropertyType is not a transactional entity")
	}
	pt.config.driver = tx.drv
	return pt
}

// String implements the fmt.Stringer.
func (pt *PropertyType) String() string {
	var builder strings.Builder
	builder.WriteString("PropertyType(")
	builder.WriteString(fmt.Sprintf("id=%v", pt.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(pt.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(pt.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", type=")
	builder.WriteString(pt.Type)
	builder.WriteString(", name=")
	builder.WriteString(pt.Name)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", pt.Index))
	builder.WriteString(", category=")
	builder.WriteString(pt.Category)
	builder.WriteString(", int_val=")
	builder.WriteString(fmt.Sprintf("%v", pt.IntVal))
	builder.WriteString(", bool_val=")
	builder.WriteString(fmt.Sprintf("%v", pt.BoolVal))
	builder.WriteString(", float_val=")
	builder.WriteString(fmt.Sprintf("%v", pt.FloatVal))
	builder.WriteString(", latitude_val=")
	builder.WriteString(fmt.Sprintf("%v", pt.LatitudeVal))
	builder.WriteString(", longitude_val=")
	builder.WriteString(fmt.Sprintf("%v", pt.LongitudeVal))
	builder.WriteString(", string_val=")
	builder.WriteString(pt.StringVal)
	builder.WriteString(", range_from_val=")
	builder.WriteString(fmt.Sprintf("%v", pt.RangeFromVal))
	builder.WriteString(", range_to_val=")
	builder.WriteString(fmt.Sprintf("%v", pt.RangeToVal))
	builder.WriteString(", is_instance_property=")
	builder.WriteString(fmt.Sprintf("%v", pt.IsInstanceProperty))
	builder.WriteString(", editable=")
	builder.WriteString(fmt.Sprintf("%v", pt.Editable))
	builder.WriteString(", mandatory=")
	builder.WriteString(fmt.Sprintf("%v", pt.Mandatory))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (pt *PropertyType) id() int {
	id, _ := strconv.Atoi(pt.ID)
	return id
}

// PropertyTypes is a parsable slice of PropertyType.
type PropertyTypes []*PropertyType

// FromRows scans the sql response data into PropertyTypes.
func (pt *PropertyTypes) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanpt := &PropertyType{}
		if err := scanpt.FromRows(rows); err != nil {
			return err
		}
		*pt = append(*pt, scanpt)
	}
	return nil
}

func (pt PropertyTypes) config(cfg config) {
	for _i := range pt {
		pt[_i].config = cfg
	}
}
