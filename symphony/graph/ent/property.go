// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// Property is the model entity for the Property schema.
type Property struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// IntVal holds the value of the "int_val" field.
	IntVal int `json:"intValue" gqlgen:"intValue"`
	// BoolVal holds the value of the "bool_val" field.
	BoolVal bool `json:"booleanValue" gqlgen:"booleanValue"`
	// FloatVal holds the value of the "float_val" field.
	FloatVal float64 `json:"floatValue" gqlgen:"floatValue"`
	// LatitudeVal holds the value of the "latitude_val" field.
	LatitudeVal float64 `json:"latitudeValue" gqlgen:"latitudeValue"`
	// LongitudeVal holds the value of the "longitude_val" field.
	LongitudeVal float64 `json:"longitudeValue" gqlgen:"longitudeValue"`
	// RangeFromVal holds the value of the "range_from_val" field.
	RangeFromVal float64 `json:"rangeFromValue" gqlgen:"rangeFromValue"`
	// RangeToVal holds the value of the "range_to_val" field.
	RangeToVal float64 `json:"rangeToValue" gqlgen:"rangeToValue"`
	// StringVal holds the value of the "string_val" field.
	StringVal string `json:"stringValue" gqlgen:"stringValue"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PropertyQuery when eager-loading is set.
	Edges                     PropertyEdges `json:"edges"`
	equipment_properties      *int
	equipment_port_properties *int
	link_properties           *int
	location_properties       *int
	project_properties        *int
	property_type             *int
	property_equipment_value  *int
	property_location_value   *int
	property_service_value    *int
	service_properties        *int
	work_order_properties     *int
}

// PropertyEdges holds the relations/edges for other nodes in the graph.
type PropertyEdges struct {
	// Type holds the value of the type edge.
	Type *PropertyType `gqlgen:"propertyType"`
	// Location holds the value of the location edge.
	Location *Location `gqlgen:"locationValue"`
	// Equipment holds the value of the equipment edge.
	Equipment *Equipment `gqlgen:"equipmentValue"`
	// Service holds the value of the service edge.
	Service *Service `gqlgen:"serviceValue"`
	// EquipmentPort holds the value of the equipment_port edge.
	EquipmentPort *EquipmentPort
	// Link holds the value of the link edge.
	Link *Link
	// WorkOrder holds the value of the work_order edge.
	WorkOrder *WorkOrder
	// Project holds the value of the project edge.
	Project *Project
	// EquipmentValue holds the value of the equipment_value edge.
	EquipmentValue *Equipment
	// LocationValue holds the value of the location_value edge.
	LocationValue *Location
	// ServiceValue holds the value of the service_value edge.
	ServiceValue *Service
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [11]bool
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) TypeOrErr() (*PropertyType, error) {
	if e.loadedTypes[0] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: propertytype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) LocationOrErr() (*Location, error) {
	if e.loadedTypes[1] {
		if e.Location == nil {
			// The edge location was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Location, nil
	}
	return nil, &NotLoadedError{edge: "location"}
}

// EquipmentOrErr returns the Equipment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) EquipmentOrErr() (*Equipment, error) {
	if e.loadedTypes[2] {
		if e.Equipment == nil {
			// The edge equipment was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipment.Label}
		}
		return e.Equipment, nil
	}
	return nil, &NotLoadedError{edge: "equipment"}
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) ServiceOrErr() (*Service, error) {
	if e.loadedTypes[3] {
		if e.Service == nil {
			// The edge service was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: service.Label}
		}
		return e.Service, nil
	}
	return nil, &NotLoadedError{edge: "service"}
}

// EquipmentPortOrErr returns the EquipmentPort value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) EquipmentPortOrErr() (*EquipmentPort, error) {
	if e.loadedTypes[4] {
		if e.EquipmentPort == nil {
			// The edge equipment_port was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmentport.Label}
		}
		return e.EquipmentPort, nil
	}
	return nil, &NotLoadedError{edge: "equipment_port"}
}

// LinkOrErr returns the Link value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) LinkOrErr() (*Link, error) {
	if e.loadedTypes[5] {
		if e.Link == nil {
			// The edge link was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: link.Label}
		}
		return e.Link, nil
	}
	return nil, &NotLoadedError{edge: "link"}
}

// WorkOrderOrErr returns the WorkOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) WorkOrderOrErr() (*WorkOrder, error) {
	if e.loadedTypes[6] {
		if e.WorkOrder == nil {
			// The edge work_order was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workorder.Label}
		}
		return e.WorkOrder, nil
	}
	return nil, &NotLoadedError{edge: "work_order"}
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[7] {
		if e.Project == nil {
			// The edge project was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// EquipmentValueOrErr returns the EquipmentValue value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) EquipmentValueOrErr() (*Equipment, error) {
	if e.loadedTypes[8] {
		if e.EquipmentValue == nil {
			// The edge equipment_value was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipment.Label}
		}
		return e.EquipmentValue, nil
	}
	return nil, &NotLoadedError{edge: "equipment_value"}
}

// LocationValueOrErr returns the LocationValue value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) LocationValueOrErr() (*Location, error) {
	if e.loadedTypes[9] {
		if e.LocationValue == nil {
			// The edge location_value was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.LocationValue, nil
	}
	return nil, &NotLoadedError{edge: "location_value"}
}

// ServiceValueOrErr returns the ServiceValue value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PropertyEdges) ServiceValueOrErr() (*Service, error) {
	if e.loadedTypes[10] {
		if e.ServiceValue == nil {
			// The edge service_value was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: service.Label}
		}
		return e.ServiceValue, nil
	}
	return nil, &NotLoadedError{edge: "service_value"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Property) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},   // id
		&sql.NullTime{},    // create_time
		&sql.NullTime{},    // update_time
		&sql.NullInt64{},   // int_val
		&sql.NullBool{},    // bool_val
		&sql.NullFloat64{}, // float_val
		&sql.NullFloat64{}, // latitude_val
		&sql.NullFloat64{}, // longitude_val
		&sql.NullFloat64{}, // range_from_val
		&sql.NullFloat64{}, // range_to_val
		&sql.NullString{},  // string_val
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Property) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // equipment_properties
		&sql.NullInt64{}, // equipment_port_properties
		&sql.NullInt64{}, // link_properties
		&sql.NullInt64{}, // location_properties
		&sql.NullInt64{}, // project_properties
		&sql.NullInt64{}, // property_type
		&sql.NullInt64{}, // property_equipment_value
		&sql.NullInt64{}, // property_location_value
		&sql.NullInt64{}, // property_service_value
		&sql.NullInt64{}, // service_properties
		&sql.NullInt64{}, // work_order_properties
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Property fields.
func (pr *Property) assignValues(values ...interface{}) error {
	if m, n := len(values), len(property.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	pr.ID = int(value.Int64)
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
	values = values[10:]
	if len(values) == len(property.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_properties", value)
		} else if value.Valid {
			pr.equipment_properties = new(int)
			*pr.equipment_properties = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_port_properties", value)
		} else if value.Valid {
			pr.equipment_port_properties = new(int)
			*pr.equipment_port_properties = int(value.Int64)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field link_properties", value)
		} else if value.Valid {
			pr.link_properties = new(int)
			*pr.link_properties = int(value.Int64)
		}
		if value, ok := values[3].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_properties", value)
		} else if value.Valid {
			pr.location_properties = new(int)
			*pr.location_properties = int(value.Int64)
		}
		if value, ok := values[4].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_properties", value)
		} else if value.Valid {
			pr.project_properties = new(int)
			*pr.project_properties = int(value.Int64)
		}
		if value, ok := values[5].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field property_type", value)
		} else if value.Valid {
			pr.property_type = new(int)
			*pr.property_type = int(value.Int64)
		}
		if value, ok := values[6].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field property_equipment_value", value)
		} else if value.Valid {
			pr.property_equipment_value = new(int)
			*pr.property_equipment_value = int(value.Int64)
		}
		if value, ok := values[7].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field property_location_value", value)
		} else if value.Valid {
			pr.property_location_value = new(int)
			*pr.property_location_value = int(value.Int64)
		}
		if value, ok := values[8].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field property_service_value", value)
		} else if value.Valid {
			pr.property_service_value = new(int)
			*pr.property_service_value = int(value.Int64)
		}
		if value, ok := values[9].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_properties", value)
		} else if value.Valid {
			pr.service_properties = new(int)
			*pr.service_properties = int(value.Int64)
		}
		if value, ok := values[10].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_properties", value)
		} else if value.Valid {
			pr.work_order_properties = new(int)
			*pr.work_order_properties = int(value.Int64)
		}
	}
	return nil
}

// QueryType queries the type edge of the Property.
func (pr *Property) QueryType() *PropertyTypeQuery {
	return (&PropertyClient{config: pr.config}).QueryType(pr)
}

// QueryLocation queries the location edge of the Property.
func (pr *Property) QueryLocation() *LocationQuery {
	return (&PropertyClient{config: pr.config}).QueryLocation(pr)
}

// QueryEquipment queries the equipment edge of the Property.
func (pr *Property) QueryEquipment() *EquipmentQuery {
	return (&PropertyClient{config: pr.config}).QueryEquipment(pr)
}

// QueryService queries the service edge of the Property.
func (pr *Property) QueryService() *ServiceQuery {
	return (&PropertyClient{config: pr.config}).QueryService(pr)
}

// QueryEquipmentPort queries the equipment_port edge of the Property.
func (pr *Property) QueryEquipmentPort() *EquipmentPortQuery {
	return (&PropertyClient{config: pr.config}).QueryEquipmentPort(pr)
}

// QueryLink queries the link edge of the Property.
func (pr *Property) QueryLink() *LinkQuery {
	return (&PropertyClient{config: pr.config}).QueryLink(pr)
}

// QueryWorkOrder queries the work_order edge of the Property.
func (pr *Property) QueryWorkOrder() *WorkOrderQuery {
	return (&PropertyClient{config: pr.config}).QueryWorkOrder(pr)
}

// QueryProject queries the project edge of the Property.
func (pr *Property) QueryProject() *ProjectQuery {
	return (&PropertyClient{config: pr.config}).QueryProject(pr)
}

// QueryEquipmentValue queries the equipment_value edge of the Property.
func (pr *Property) QueryEquipmentValue() *EquipmentQuery {
	return (&PropertyClient{config: pr.config}).QueryEquipmentValue(pr)
}

// QueryLocationValue queries the location_value edge of the Property.
func (pr *Property) QueryLocationValue() *LocationQuery {
	return (&PropertyClient{config: pr.config}).QueryLocationValue(pr)
}

// QueryServiceValue queries the service_value edge of the Property.
func (pr *Property) QueryServiceValue() *ServiceQuery {
	return (&PropertyClient{config: pr.config}).QueryServiceValue(pr)
}

// Update returns a builder for updating this Property.
// Note that, you need to call Property.Unwrap() before calling this method, if this Property
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Property) Update() *PropertyUpdateOne {
	return (&PropertyClient{config: pr.config}).UpdateOne(pr)
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

// Properties is a parsable slice of Property.
type Properties []*Property

func (pr Properties) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
