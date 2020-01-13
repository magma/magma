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
	"github.com/facebookincubator/symphony/graph/ent/property"
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

// scanValues returns the types for scanning values from sql.Rows.
func (*Property) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullTime{},
		&sql.NullTime{},
		&sql.NullInt64{},
		&sql.NullBool{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullString{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Property fields.
func (pr *Property) assignValues(values ...interface{}) error {
	if m, n := len(values), len(property.Columns); m != n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	pr.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		pr.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		pr.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field int_val", values[2])
	} else if value.Valid {
		pr.IntVal = int(value.Int64)
	}
	if value, ok := values[3].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field bool_val", values[3])
	} else if value.Valid {
		pr.BoolVal = value.Bool
	}
	if value, ok := values[4].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field float_val", values[4])
	} else if value.Valid {
		pr.FloatVal = value.Float64
	}
	if value, ok := values[5].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field latitude_val", values[5])
	} else if value.Valid {
		pr.LatitudeVal = value.Float64
	}
	if value, ok := values[6].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field longitude_val", values[6])
	} else if value.Valid {
		pr.LongitudeVal = value.Float64
	}
	if value, ok := values[7].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field range_from_val", values[7])
	} else if value.Valid {
		pr.RangeFromVal = value.Float64
	}
	if value, ok := values[8].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field range_to_val", values[8])
	} else if value.Valid {
		pr.RangeToVal = value.Float64
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field string_val", values[9])
	} else if value.Valid {
		pr.StringVal = value.String
	}
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

func (pr Properties) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
