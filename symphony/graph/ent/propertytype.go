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
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PropertyTypeQuery when eager-loading is set.
	Edges                       PropertyTypeEdges `json:"edges"`
	equipment_port_type_id      *string
	link_equipment_port_type_id *string
	equipment_type_id           *string
	location_type_id            *string
	project_type_id             *string
	service_type_id             *string
	work_order_type_id          *string
}

// PropertyTypeEdges holds the relations/edges for other nodes in the graph.
type PropertyTypeEdges struct {
	// Properties holds the value of the properties edge.
	Properties []*Property
	// LocationType holds the value of the location_type edge.
	LocationType *LocationType
	// EquipmentPortType holds the value of the equipment_port_type edge.
	EquipmentPortType *EquipmentPortType
	// LinkEquipmentPortType holds the value of the link_equipment_port_type edge.
	LinkEquipmentPortType *EquipmentPortType
	// EquipmentType holds the value of the equipment_type edge.
	EquipmentType *EquipmentType
	// ServiceType holds the value of the service_type edge.
	ServiceType *ServiceType
	// WorkOrderType holds the value of the work_order_type edge.
	WorkOrderType *WorkOrderType
	// ProjectType holds the value of the project_type edge.
	ProjectType *ProjectType
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [8]bool
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e PropertyTypeEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[0] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// LocationTypeOrErr returns the LocationType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyTypeEdges) LocationTypeOrErr() (*LocationType, error) {
	if e.loadedTypes[1] {
		if e.LocationType == nil {
			// The edge location_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: locationtype.Label}
		}
		return e.LocationType, nil
	}
	return nil, &NotLoadedError{edge: "location_type"}
}

// EquipmentPortTypeOrErr returns the EquipmentPortType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyTypeEdges) EquipmentPortTypeOrErr() (*EquipmentPortType, error) {
	if e.loadedTypes[2] {
		if e.EquipmentPortType == nil {
			// The edge equipment_port_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmentporttype.Label}
		}
		return e.EquipmentPortType, nil
	}
	return nil, &NotLoadedError{edge: "equipment_port_type"}
}

// LinkEquipmentPortTypeOrErr returns the LinkEquipmentPortType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyTypeEdges) LinkEquipmentPortTypeOrErr() (*EquipmentPortType, error) {
	if e.loadedTypes[3] {
		if e.LinkEquipmentPortType == nil {
			// The edge link_equipment_port_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmentporttype.Label}
		}
		return e.LinkEquipmentPortType, nil
	}
	return nil, &NotLoadedError{edge: "link_equipment_port_type"}
}

// EquipmentTypeOrErr returns the EquipmentType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyTypeEdges) EquipmentTypeOrErr() (*EquipmentType, error) {
	if e.loadedTypes[4] {
		if e.EquipmentType == nil {
			// The edge equipment_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmenttype.Label}
		}
		return e.EquipmentType, nil
	}
	return nil, &NotLoadedError{edge: "equipment_type"}
}

// ServiceTypeOrErr returns the ServiceType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyTypeEdges) ServiceTypeOrErr() (*ServiceType, error) {
	if e.loadedTypes[5] {
		if e.ServiceType == nil {
			// The edge service_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: servicetype.Label}
		}
		return e.ServiceType, nil
	}
	return nil, &NotLoadedError{edge: "service_type"}
}

// WorkOrderTypeOrErr returns the WorkOrderType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyTypeEdges) WorkOrderTypeOrErr() (*WorkOrderType, error) {
	if e.loadedTypes[6] {
		if e.WorkOrderType == nil {
			// The edge work_order_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workordertype.Label}
		}
		return e.WorkOrderType, nil
	}
	return nil, &NotLoadedError{edge: "work_order_type"}
}

// ProjectTypeOrErr returns the ProjectType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyTypeEdges) ProjectTypeOrErr() (*ProjectType, error) {
	if e.loadedTypes[7] {
		if e.ProjectType == nil {
			// The edge project_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: projecttype.Label}
		}
		return e.ProjectType, nil
	}
	return nil, &NotLoadedError{edge: "project_type"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PropertyType) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},   // id
		&sql.NullTime{},    // create_time
		&sql.NullTime{},    // update_time
		&sql.NullString{},  // type
		&sql.NullString{},  // name
		&sql.NullInt64{},   // index
		&sql.NullString{},  // category
		&sql.NullInt64{},   // int_val
		&sql.NullBool{},    // bool_val
		&sql.NullFloat64{}, // float_val
		&sql.NullFloat64{}, // latitude_val
		&sql.NullFloat64{}, // longitude_val
		&sql.NullString{},  // string_val
		&sql.NullFloat64{}, // range_from_val
		&sql.NullFloat64{}, // range_to_val
		&sql.NullBool{},    // is_instance_property
		&sql.NullBool{},    // editable
		&sql.NullBool{},    // mandatory
		&sql.NullBool{},    // deleted
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*PropertyType) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // equipment_port_type_id
		&sql.NullInt64{}, // link_equipment_port_type_id
		&sql.NullInt64{}, // equipment_type_id
		&sql.NullInt64{}, // location_type_id
		&sql.NullInt64{}, // project_type_id
		&sql.NullInt64{}, // service_type_id
		&sql.NullInt64{}, // work_order_type_id
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PropertyType fields.
func (pt *PropertyType) assignValues(values ...interface{}) error {
	if m, n := len(values), len(propertytype.Columns); m < n {
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
	values = values[18:]
	if len(values) == len(propertytype.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_port_type_id", value)
		} else if value.Valid {
			pt.equipment_port_type_id = new(string)
			*pt.equipment_port_type_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field link_equipment_port_type_id", value)
		} else if value.Valid {
			pt.link_equipment_port_type_id = new(string)
			*pt.link_equipment_port_type_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_type_id", value)
		} else if value.Valid {
			pt.equipment_type_id = new(string)
			*pt.equipment_type_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[3].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_type_id", value)
		} else if value.Valid {
			pt.location_type_id = new(string)
			*pt.location_type_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[4].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_type_id", value)
		} else if value.Valid {
			pt.project_type_id = new(string)
			*pt.project_type_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[5].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_type_id", value)
		} else if value.Valid {
			pt.service_type_id = new(string)
			*pt.service_type_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[6].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_type_id", value)
		} else if value.Valid {
			pt.work_order_type_id = new(string)
			*pt.work_order_type_id = strconv.FormatInt(value.Int64, 10)
		}
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
