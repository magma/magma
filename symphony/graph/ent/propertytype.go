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
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
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
	// Deleted holds the value of the "deleted" field.
	Deleted bool `json:"deleted,omitempty" gqlgen:"isDeleted"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PropertyType) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullTime{},
		&sql.NullTime{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullInt64{},
		&sql.NullString{},
		&sql.NullInt64{},
		&sql.NullBool{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullString{},
		&sql.NullFloat64{},
		&sql.NullFloat64{},
		&sql.NullBool{},
		&sql.NullBool{},
		&sql.NullBool{},
		&sql.NullBool{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PropertyType fields.
func (pt *PropertyType) assignValues(values ...interface{}) error {
	if m, n := len(values), len(propertytype.Columns); m != n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	pt.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		pt.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		pt.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field type", values[2])
	} else if value.Valid {
		pt.Type = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[3])
	} else if value.Valid {
		pt.Name = value.String
	}
	if value, ok := values[4].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field index", values[4])
	} else if value.Valid {
		pt.Index = int(value.Int64)
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field category", values[5])
	} else if value.Valid {
		pt.Category = value.String
	}
	if value, ok := values[6].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field int_val", values[6])
	} else if value.Valid {
		pt.IntVal = int(value.Int64)
	}
	if value, ok := values[7].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field bool_val", values[7])
	} else if value.Valid {
		pt.BoolVal = value.Bool
	}
	if value, ok := values[8].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field float_val", values[8])
	} else if value.Valid {
		pt.FloatVal = value.Float64
	}
	if value, ok := values[9].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field latitude_val", values[9])
	} else if value.Valid {
		pt.LatitudeVal = value.Float64
	}
	if value, ok := values[10].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field longitude_val", values[10])
	} else if value.Valid {
		pt.LongitudeVal = value.Float64
	}
	if value, ok := values[11].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field string_val", values[11])
	} else if value.Valid {
		pt.StringVal = value.String
	}
	if value, ok := values[12].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field range_from_val", values[12])
	} else if value.Valid {
		pt.RangeFromVal = value.Float64
	}
	if value, ok := values[13].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field range_to_val", values[13])
	} else if value.Valid {
		pt.RangeToVal = value.Float64
	}
	if value, ok := values[14].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field is_instance_property", values[14])
	} else if value.Valid {
		pt.IsInstanceProperty = value.Bool
	}
	if value, ok := values[15].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field editable", values[15])
	} else if value.Valid {
		pt.Editable = value.Bool
	}
	if value, ok := values[16].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field mandatory", values[16])
	} else if value.Valid {
		pt.Mandatory = value.Bool
	}
	if value, ok := values[17].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field deleted", values[17])
	} else if value.Valid {
		pt.Deleted = value.Bool
	}
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
	builder.WriteString(", deleted=")
	builder.WriteString(fmt.Sprintf("%v", pt.Deleted))
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

func (pt PropertyTypes) config(cfg config) {
	for _i := range pt {
		pt[_i].config = cfg
	}
}
