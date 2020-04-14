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
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceEndpointDefinition is the model entity for the ServiceEndpointDefinition schema.
type ServiceEndpointDefinition struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Role holds the value of the "role" field.
	Role string `json:"role,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceEndpointDefinitionQuery when eager-loading is set.
	Edges                                       ServiceEndpointDefinitionEdges `json:"edges"`
	equipment_type_service_endpoint_definitions *int
	service_type_endpoint_definitions           *int
}

// ServiceEndpointDefinitionEdges holds the relations/edges for other nodes in the graph.
type ServiceEndpointDefinitionEdges struct {
	// Endpoints holds the value of the endpoints edge.
	Endpoints []*ServiceEndpoint
	// ServiceType holds the value of the service_type edge.
	ServiceType *ServiceType
	// EquipmentType holds the value of the equipment_type edge.
	EquipmentType *EquipmentType
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// EndpointsOrErr returns the Endpoints value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEndpointDefinitionEdges) EndpointsOrErr() ([]*ServiceEndpoint, error) {
	if e.loadedTypes[0] {
		return e.Endpoints, nil
	}
	return nil, &NotLoadedError{edge: "endpoints"}
}

// ServiceTypeOrErr returns the ServiceType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEndpointDefinitionEdges) ServiceTypeOrErr() (*ServiceType, error) {
	if e.loadedTypes[1] {
		if e.ServiceType == nil {
			// The edge service_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: servicetype.Label}
		}
		return e.ServiceType, nil
	}
	return nil, &NotLoadedError{edge: "service_type"}
}

// EquipmentTypeOrErr returns the EquipmentType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEndpointDefinitionEdges) EquipmentTypeOrErr() (*EquipmentType, error) {
	if e.loadedTypes[2] {
		if e.EquipmentType == nil {
			// The edge equipment_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmenttype.Label}
		}
		return e.EquipmentType, nil
	}
	return nil, &NotLoadedError{edge: "equipment_type"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ServiceEndpointDefinition) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // role
		&sql.NullString{}, // name
		&sql.NullInt64{},  // index
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*ServiceEndpointDefinition) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // equipment_type_service_endpoint_definitions
		&sql.NullInt64{}, // service_type_endpoint_definitions
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ServiceEndpointDefinition fields.
func (sed *ServiceEndpointDefinition) assignValues(values ...interface{}) error {
	if m, n := len(values), len(serviceendpointdefinition.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	sed.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		sed.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		sed.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field role", values[2])
	} else if value.Valid {
		sed.Role = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[3])
	} else if value.Valid {
		sed.Name = value.String
	}
	if value, ok := values[4].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field index", values[4])
	} else if value.Valid {
		sed.Index = int(value.Int64)
	}
	values = values[5:]
	if len(values) == len(serviceendpointdefinition.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_type_service_endpoint_definitions", value)
		} else if value.Valid {
			sed.equipment_type_service_endpoint_definitions = new(int)
			*sed.equipment_type_service_endpoint_definitions = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_type_endpoint_definitions", value)
		} else if value.Valid {
			sed.service_type_endpoint_definitions = new(int)
			*sed.service_type_endpoint_definitions = int(value.Int64)
		}
	}
	return nil
}

// QueryEndpoints queries the endpoints edge of the ServiceEndpointDefinition.
func (sed *ServiceEndpointDefinition) QueryEndpoints() *ServiceEndpointQuery {
	return (&ServiceEndpointDefinitionClient{config: sed.config}).QueryEndpoints(sed)
}

// QueryServiceType queries the service_type edge of the ServiceEndpointDefinition.
func (sed *ServiceEndpointDefinition) QueryServiceType() *ServiceTypeQuery {
	return (&ServiceEndpointDefinitionClient{config: sed.config}).QueryServiceType(sed)
}

// QueryEquipmentType queries the equipment_type edge of the ServiceEndpointDefinition.
func (sed *ServiceEndpointDefinition) QueryEquipmentType() *EquipmentTypeQuery {
	return (&ServiceEndpointDefinitionClient{config: sed.config}).QueryEquipmentType(sed)
}

// Update returns a builder for updating this ServiceEndpointDefinition.
// Note that, you need to call ServiceEndpointDefinition.Unwrap() before calling this method, if this ServiceEndpointDefinition
// was returned from a transaction, and the transaction was committed or rolled back.
func (sed *ServiceEndpointDefinition) Update() *ServiceEndpointDefinitionUpdateOne {
	return (&ServiceEndpointDefinitionClient{config: sed.config}).UpdateOne(sed)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (sed *ServiceEndpointDefinition) Unwrap() *ServiceEndpointDefinition {
	tx, ok := sed.config.driver.(*txDriver)
	if !ok {
		panic("ent: ServiceEndpointDefinition is not a transactional entity")
	}
	sed.config.driver = tx.drv
	return sed
}

// String implements the fmt.Stringer.
func (sed *ServiceEndpointDefinition) String() string {
	var builder strings.Builder
	builder.WriteString("ServiceEndpointDefinition(")
	builder.WriteString(fmt.Sprintf("id=%v", sed.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(sed.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(sed.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", role=")
	builder.WriteString(sed.Role)
	builder.WriteString(", name=")
	builder.WriteString(sed.Name)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", sed.Index))
	builder.WriteByte(')')
	return builder.String()
}

// ServiceEndpointDefinitions is a parsable slice of ServiceEndpointDefinition.
type ServiceEndpointDefinitions []*ServiceEndpointDefinition

func (sed ServiceEndpointDefinitions) config(cfg config) {
	for _i := range sed {
		sed[_i].config = cfg
	}
}
