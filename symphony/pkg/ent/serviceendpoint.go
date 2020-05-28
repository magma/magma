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
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentport"
	"github.com/facebookincubator/symphony/pkg/ent/service"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"
)

// ServiceEndpoint is the model entity for the ServiceEndpoint schema.
type ServiceEndpoint struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceEndpointQuery when eager-loading is set.
	Edges                                 ServiceEndpointEdges `json:"edges"`
	service_endpoints                     *int
	service_endpoint_port                 *int
	service_endpoint_equipment            *int
	service_endpoint_definition_endpoints *int
}

// ServiceEndpointEdges holds the relations/edges for other nodes in the graph.
type ServiceEndpointEdges struct {
	// Port holds the value of the port edge.
	Port *EquipmentPort
	// Equipment holds the value of the equipment edge.
	Equipment *Equipment
	// Service holds the value of the service edge.
	Service *Service
	// Definition holds the value of the definition edge.
	Definition *ServiceEndpointDefinition
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// PortOrErr returns the Port value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEndpointEdges) PortOrErr() (*EquipmentPort, error) {
	if e.loadedTypes[0] {
		if e.Port == nil {
			// The edge port was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmentport.Label}
		}
		return e.Port, nil
	}
	return nil, &NotLoadedError{edge: "port"}
}

// EquipmentOrErr returns the Equipment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEndpointEdges) EquipmentOrErr() (*Equipment, error) {
	if e.loadedTypes[1] {
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
func (e ServiceEndpointEdges) ServiceOrErr() (*Service, error) {
	if e.loadedTypes[2] {
		if e.Service == nil {
			// The edge service was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: service.Label}
		}
		return e.Service, nil
	}
	return nil, &NotLoadedError{edge: "service"}
}

// DefinitionOrErr returns the Definition value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEndpointEdges) DefinitionOrErr() (*ServiceEndpointDefinition, error) {
	if e.loadedTypes[3] {
		if e.Definition == nil {
			// The edge definition was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: serviceendpointdefinition.Label}
		}
		return e.Definition, nil
	}
	return nil, &NotLoadedError{edge: "definition"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ServiceEndpoint) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // id
		&sql.NullTime{},  // create_time
		&sql.NullTime{},  // update_time
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*ServiceEndpoint) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // service_endpoints
		&sql.NullInt64{}, // service_endpoint_port
		&sql.NullInt64{}, // service_endpoint_equipment
		&sql.NullInt64{}, // service_endpoint_definition_endpoints
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ServiceEndpoint fields.
func (se *ServiceEndpoint) assignValues(values ...interface{}) error {
	if m, n := len(values), len(serviceendpoint.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	se.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		se.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		se.UpdateTime = value.Time
	}
	values = values[2:]
	if len(values) == len(serviceendpoint.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_endpoints", value)
		} else if value.Valid {
			se.service_endpoints = new(int)
			*se.service_endpoints = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_endpoint_port", value)
		} else if value.Valid {
			se.service_endpoint_port = new(int)
			*se.service_endpoint_port = int(value.Int64)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_endpoint_equipment", value)
		} else if value.Valid {
			se.service_endpoint_equipment = new(int)
			*se.service_endpoint_equipment = int(value.Int64)
		}
		if value, ok := values[3].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field service_endpoint_definition_endpoints", value)
		} else if value.Valid {
			se.service_endpoint_definition_endpoints = new(int)
			*se.service_endpoint_definition_endpoints = int(value.Int64)
		}
	}
	return nil
}

// QueryPort queries the port edge of the ServiceEndpoint.
func (se *ServiceEndpoint) QueryPort() *EquipmentPortQuery {
	return (&ServiceEndpointClient{config: se.config}).QueryPort(se)
}

// QueryEquipment queries the equipment edge of the ServiceEndpoint.
func (se *ServiceEndpoint) QueryEquipment() *EquipmentQuery {
	return (&ServiceEndpointClient{config: se.config}).QueryEquipment(se)
}

// QueryService queries the service edge of the ServiceEndpoint.
func (se *ServiceEndpoint) QueryService() *ServiceQuery {
	return (&ServiceEndpointClient{config: se.config}).QueryService(se)
}

// QueryDefinition queries the definition edge of the ServiceEndpoint.
func (se *ServiceEndpoint) QueryDefinition() *ServiceEndpointDefinitionQuery {
	return (&ServiceEndpointClient{config: se.config}).QueryDefinition(se)
}

// Update returns a builder for updating this ServiceEndpoint.
// Note that, you need to call ServiceEndpoint.Unwrap() before calling this method, if this ServiceEndpoint
// was returned from a transaction, and the transaction was committed or rolled back.
func (se *ServiceEndpoint) Update() *ServiceEndpointUpdateOne {
	return (&ServiceEndpointClient{config: se.config}).UpdateOne(se)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (se *ServiceEndpoint) Unwrap() *ServiceEndpoint {
	tx, ok := se.config.driver.(*txDriver)
	if !ok {
		panic("ent: ServiceEndpoint is not a transactional entity")
	}
	se.config.driver = tx.drv
	return se
}

// String implements the fmt.Stringer.
func (se *ServiceEndpoint) String() string {
	var builder strings.Builder
	builder.WriteString("ServiceEndpoint(")
	builder.WriteString(fmt.Sprintf("id=%v", se.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(se.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(se.UpdateTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// ServiceEndpoints is a parsable slice of ServiceEndpoint.
type ServiceEndpoints []*ServiceEndpoint

func (se ServiceEndpoints) config(cfg config) {
	for _i := range se {
		se[_i].config = cfg
	}
}
