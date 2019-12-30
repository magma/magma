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

// Property is the model entity for the Property schema.
type Property struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
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
	// RangeFromVal holds the value of the "range_from_val" field.
	RangeFromVal float64 `json:"range_from_val,omitempty" gqlgen:"rangeFromValue"`
	// RangeToVal holds the value of the "range_to_val" field.
	RangeToVal float64 `json:"range_to_val,omitempty" gqlgen:"rangeToValue"`
	// StringVal holds the value of the "string_val" field.
	StringVal string `json:"string_val,omitempty" gqlgen:"stringValue"`
}

// FromRows scans the sql response data into Property.
func (pr *Property) FromRows(rows *sql.Rows) error {
	var scanpr struct {
		ID           int
		CreateTime   sql.NullTime
		UpdateTime   sql.NullTime
		IntVal       sql.NullInt64
		BoolVal      sql.NullBool
		FloatVal     sql.NullFloat64
		LatitudeVal  sql.NullFloat64
		LongitudeVal sql.NullFloat64
		RangeFromVal sql.NullFloat64
		RangeToVal   sql.NullFloat64
		StringVal    sql.NullString
	}
	// the order here should be the same as in the `property.Columns`.
	if err := rows.Scan(
		&scanpr.ID,
		&scanpr.CreateTime,
		&scanpr.UpdateTime,
		&scanpr.IntVal,
		&scanpr.BoolVal,
		&scanpr.FloatVal,
		&scanpr.LatitudeVal,
		&scanpr.LongitudeVal,
		&scanpr.RangeFromVal,
		&scanpr.RangeToVal,
		&scanpr.StringVal,
	); err != nil {
		return err
	}
	pr.ID = strconv.Itoa(scanpr.ID)
	pr.CreateTime = scanpr.CreateTime.Time
	pr.UpdateTime = scanpr.UpdateTime.Time
	pr.IntVal = int(scanpr.IntVal.Int64)
	pr.BoolVal = scanpr.BoolVal.Bool
	pr.FloatVal = scanpr.FloatVal.Float64
	pr.LatitudeVal = scanpr.LatitudeVal.Float64
	pr.LongitudeVal = scanpr.LongitudeVal.Float64
	pr.RangeFromVal = scanpr.RangeFromVal.Float64
	pr.RangeToVal = scanpr.RangeToVal.Float64
	pr.StringVal = scanpr.StringVal.String
	return nil
}

// QueryType queries the type edge of the Property.
func (pr *Property) QueryType() *PropertyTypeQuery {
	return (&PropertyClient{pr.config}).QueryType(pr)
}

// QueryLocation queries the location edge of the Property.
func (pr *Property) QueryLocation() *LocationQuery {
	return (&PropertyClient{pr.config}).QueryLocation(pr)
}

// QueryEquipment queries the equipment edge of the Property.
func (pr *Property) QueryEquipment() *EquipmentQuery {
	return (&PropertyClient{pr.config}).QueryEquipment(pr)
}

// QueryService queries the service edge of the Property.
func (pr *Property) QueryService() *ServiceQuery {
	return (&PropertyClient{pr.config}).QueryService(pr)
}

// QueryEquipmentPort queries the equipment_port edge of the Property.
func (pr *Property) QueryEquipmentPort() *EquipmentPortQuery {
	return (&PropertyClient{pr.config}).QueryEquipmentPort(pr)
}

// QueryLink queries the link edge of the Property.
func (pr *Property) QueryLink() *LinkQuery {
	return (&PropertyClient{pr.config}).QueryLink(pr)
}

// QueryWorkOrder queries the work_order edge of the Property.
func (pr *Property) QueryWorkOrder() *WorkOrderQuery {
	return (&PropertyClient{pr.config}).QueryWorkOrder(pr)
}

// QueryProject queries the project edge of the Property.
func (pr *Property) QueryProject() *ProjectQuery {
	return (&PropertyClient{pr.config}).QueryProject(pr)
}

// QueryEquipmentValue queries the equipment_value edge of the Property.
func (pr *Property) QueryEquipmentValue() *EquipmentQuery {
	return (&PropertyClient{pr.config}).QueryEquipmentValue(pr)
}

// QueryLocationValue queries the location_value edge of the Property.
func (pr *Property) QueryLocationValue() *LocationQuery {
	return (&PropertyClient{pr.config}).QueryLocationValue(pr)
}

// QueryServiceValue queries the service_value edge of the Property.
func (pr *Property) QueryServiceValue() *ServiceQuery {
	return (&PropertyClient{pr.config}).QueryServiceValue(pr)
}

// Update returns a builder for updating this Property.
// Note that, you need to call Property.Unwrap() before calling this method, if this Property
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Property) Update() *PropertyUpdateOne {
	return (&PropertyClient{pr.config}).UpdateOne(pr)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (pr *Property) Unwrap() *Property {
	tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("ent: Property is not a transactional entity")
	}
	pr.config.driver = tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Property) String() string {
	var builder strings.Builder
	builder.WriteString("Property(")
	builder.WriteString(fmt.Sprintf("id=%v", pr.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(pr.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(pr.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", int_val=")
	builder.WriteString(fmt.Sprintf("%v", pr.IntVal))
	builder.WriteString(", bool_val=")
	builder.WriteString(fmt.Sprintf("%v", pr.BoolVal))
	builder.WriteString(", float_val=")
	builder.WriteString(fmt.Sprintf("%v", pr.FloatVal))
	builder.WriteString(", latitude_val=")
	builder.WriteString(fmt.Sprintf("%v", pr.LatitudeVal))
	builder.WriteString(", longitude_val=")
	builder.WriteString(fmt.Sprintf("%v", pr.LongitudeVal))
	builder.WriteString(", range_from_val=")
	builder.WriteString(fmt.Sprintf("%v", pr.RangeFromVal))
	builder.WriteString(", range_to_val=")
	builder.WriteString(fmt.Sprintf("%v", pr.RangeToVal))
	builder.WriteString(", string_val=")
	builder.WriteString(pr.StringVal)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (pr *Property) id() int {
	id, _ := strconv.Atoi(pr.ID)
	return id
}

// Properties is a parsable slice of Property.
type Properties []*Property

// FromRows scans the sql response data into Properties.
func (pr *Properties) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanpr := &Property{}
		if err := scanpr.FromRows(rows); err != nil {
			return err
		}
		*pr = append(*pr, scanpr)
	}
	return nil
}

func (pr Properties) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
