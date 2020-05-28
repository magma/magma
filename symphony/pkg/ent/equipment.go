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
	"github.com/facebookincubator/symphony/pkg/ent/equipmentposition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

// Equipment is the model entity for the Equipment schema.
type Equipment struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// FutureState holds the value of the "future_state" field.
	FutureState string `json:"future_state,omitempty"`
	// DeviceID holds the value of the "device_id" field.
	DeviceID string `json:"device_id,omitempty"`
	// ExternalID holds the value of the "external_id" field.
	ExternalID string `json:"external_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EquipmentQuery when eager-loading is set.
	Edges                         EquipmentEdges `json:"edges"`
	equipment_type                *int
	equipment_work_order          *int
	equipment_position_attachment *int
	location_equipment            *int
}

// EquipmentEdges holds the relations/edges for other nodes in the graph.
type EquipmentEdges struct {
	// Type holds the value of the type edge.
	Type *EquipmentType `gqlgen:"equipmentType"`
	// Location holds the value of the location edge.
	Location *Location `gqlgen:"parentLocation"`
	// ParentPosition holds the value of the parent_position edge.
	ParentPosition *EquipmentPosition `gqlgen:"parentPosition"`
	// Positions holds the value of the positions edge.
	Positions []*EquipmentPosition `gqlgen:"positions"`
	// Ports holds the value of the ports edge.
	Ports []*EquipmentPort `gqlgen:"ports"`
	// WorkOrder holds the value of the work_order edge.
	WorkOrder *WorkOrder `gqlgen:"workOrder"`
	// Properties holds the value of the properties edge.
	Properties []*Property `gqlgen:"properties"`
	// Files holds the value of the files edge.
	Files []*File `gqlgen:"files"`
	// Hyperlinks holds the value of the hyperlinks edge.
	Hyperlinks []*Hyperlink `gqlgen:"hyperlinks"`
	// Endpoints holds the value of the endpoints edge.
	Endpoints []*ServiceEndpoint
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [10]bool
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EquipmentEdges) TypeOrErr() (*EquipmentType, error) {
	if e.loadedTypes[0] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmenttype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EquipmentEdges) LocationOrErr() (*Location, error) {
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

// ParentPositionOrErr returns the ParentPosition value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EquipmentEdges) ParentPositionOrErr() (*EquipmentPosition, error) {
	if e.loadedTypes[2] {
		if e.ParentPosition == nil {
			// The edge parent_position was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipmentposition.Label}
		}
		return e.ParentPosition, nil
	}
	return nil, &NotLoadedError{edge: "parent_position"}
}

// PositionsOrErr returns the Positions value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentEdges) PositionsOrErr() ([]*EquipmentPosition, error) {
	if e.loadedTypes[3] {
		return e.Positions, nil
	}
	return nil, &NotLoadedError{edge: "positions"}
}

// PortsOrErr returns the Ports value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentEdges) PortsOrErr() ([]*EquipmentPort, error) {
	if e.loadedTypes[4] {
		return e.Ports, nil
	}
	return nil, &NotLoadedError{edge: "ports"}
}

// WorkOrderOrErr returns the WorkOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EquipmentEdges) WorkOrderOrErr() (*WorkOrder, error) {
	if e.loadedTypes[5] {
		if e.WorkOrder == nil {
			// The edge work_order was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workorder.Label}
		}
		return e.WorkOrder, nil
	}
	return nil, &NotLoadedError{edge: "work_order"}
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[6] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// FilesOrErr returns the Files value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentEdges) FilesOrErr() ([]*File, error) {
	if e.loadedTypes[7] {
		return e.Files, nil
	}
	return nil, &NotLoadedError{edge: "files"}
}

// HyperlinksOrErr returns the Hyperlinks value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentEdges) HyperlinksOrErr() ([]*Hyperlink, error) {
	if e.loadedTypes[8] {
		return e.Hyperlinks, nil
	}
	return nil, &NotLoadedError{edge: "hyperlinks"}
}

// EndpointsOrErr returns the Endpoints value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentEdges) EndpointsOrErr() ([]*ServiceEndpoint, error) {
	if e.loadedTypes[9] {
		return e.Endpoints, nil
	}
	return nil, &NotLoadedError{edge: "endpoints"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Equipment) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // future_state
		&sql.NullString{}, // device_id
		&sql.NullString{}, // external_id
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Equipment) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // equipment_type
		&sql.NullInt64{}, // equipment_work_order
		&sql.NullInt64{}, // equipment_position_attachment
		&sql.NullInt64{}, // location_equipment
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Equipment fields.
func (e *Equipment) assignValues(values ...interface{}) error {
	if m, n := len(values), len(equipment.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	e.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		e.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		e.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		e.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field future_state", values[3])
	} else if value.Valid {
		e.FutureState = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field device_id", values[4])
	} else if value.Valid {
		e.DeviceID = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field external_id", values[5])
	} else if value.Valid {
		e.ExternalID = value.String
	}
	values = values[6:]
	if len(values) == len(equipment.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_type", value)
		} else if value.Valid {
			e.equipment_type = new(int)
			*e.equipment_type = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_work_order", value)
		} else if value.Valid {
			e.equipment_work_order = new(int)
			*e.equipment_work_order = int(value.Int64)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_position_attachment", value)
		} else if value.Valid {
			e.equipment_position_attachment = new(int)
			*e.equipment_position_attachment = int(value.Int64)
		}
		if value, ok := values[3].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_equipment", value)
		} else if value.Valid {
			e.location_equipment = new(int)
			*e.location_equipment = int(value.Int64)
		}
	}
	return nil
}

// QueryType queries the type edge of the Equipment.
func (e *Equipment) QueryType() *EquipmentTypeQuery {
	return (&EquipmentClient{config: e.config}).QueryType(e)
}

// QueryLocation queries the location edge of the Equipment.
func (e *Equipment) QueryLocation() *LocationQuery {
	return (&EquipmentClient{config: e.config}).QueryLocation(e)
}

// QueryParentPosition queries the parent_position edge of the Equipment.
func (e *Equipment) QueryParentPosition() *EquipmentPositionQuery {
	return (&EquipmentClient{config: e.config}).QueryParentPosition(e)
}

// QueryPositions queries the positions edge of the Equipment.
func (e *Equipment) QueryPositions() *EquipmentPositionQuery {
	return (&EquipmentClient{config: e.config}).QueryPositions(e)
}

// QueryPorts queries the ports edge of the Equipment.
func (e *Equipment) QueryPorts() *EquipmentPortQuery {
	return (&EquipmentClient{config: e.config}).QueryPorts(e)
}

// QueryWorkOrder queries the work_order edge of the Equipment.
func (e *Equipment) QueryWorkOrder() *WorkOrderQuery {
	return (&EquipmentClient{config: e.config}).QueryWorkOrder(e)
}

// QueryProperties queries the properties edge of the Equipment.
func (e *Equipment) QueryProperties() *PropertyQuery {
	return (&EquipmentClient{config: e.config}).QueryProperties(e)
}

// QueryFiles queries the files edge of the Equipment.
func (e *Equipment) QueryFiles() *FileQuery {
	return (&EquipmentClient{config: e.config}).QueryFiles(e)
}

// QueryHyperlinks queries the hyperlinks edge of the Equipment.
func (e *Equipment) QueryHyperlinks() *HyperlinkQuery {
	return (&EquipmentClient{config: e.config}).QueryHyperlinks(e)
}

// QueryEndpoints queries the endpoints edge of the Equipment.
func (e *Equipment) QueryEndpoints() *ServiceEndpointQuery {
	return (&EquipmentClient{config: e.config}).QueryEndpoints(e)
}

// Update returns a builder for updating this Equipment.
// Note that, you need to call Equipment.Unwrap() before calling this method, if this Equipment
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Equipment) Update() *EquipmentUpdateOne {
	return (&EquipmentClient{config: e.config}).UpdateOne(e)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (e *Equipment) Unwrap() *Equipment {
	tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Equipment is not a transactional entity")
	}
	e.config.driver = tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Equipment) String() string {
	var builder strings.Builder
	builder.WriteString("Equipment(")
	builder.WriteString(fmt.Sprintf("id=%v", e.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(e.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(e.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(e.Name)
	builder.WriteString(", future_state=")
	builder.WriteString(e.FutureState)
	builder.WriteString(", device_id=")
	builder.WriteString(e.DeviceID)
	builder.WriteString(", external_id=")
	builder.WriteString(e.ExternalID)
	builder.WriteByte(')')
	return builder.String()
}

// EquipmentSlice is a parsable slice of Equipment.
type EquipmentSlice []*Equipment

func (e EquipmentSlice) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}
