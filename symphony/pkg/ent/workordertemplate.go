// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/ent/workordertemplate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
)

// WorkOrderTemplate is the model entity for the WorkOrderTemplate schema.
type WorkOrderTemplate struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the WorkOrderTemplateQuery when eager-loading is set.
	Edges                    WorkOrderTemplateEdges `json:"edges"`
	work_order_template_type *int
}

// WorkOrderTemplateEdges holds the relations/edges for other nodes in the graph.
type WorkOrderTemplateEdges struct {
	// PropertyTypes holds the value of the property_types edge.
	PropertyTypes []*PropertyType
	// CheckListCategoryDefinitions holds the value of the check_list_category_definitions edge.
	CheckListCategoryDefinitions []*CheckListCategoryDefinition
	// Type holds the value of the type edge.
	Type *WorkOrderType
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// PropertyTypesOrErr returns the PropertyTypes value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderTemplateEdges) PropertyTypesOrErr() ([]*PropertyType, error) {
	if e.loadedTypes[0] {
		return e.PropertyTypes, nil
	}
	return nil, &NotLoadedError{edge: "property_types"}
}

// CheckListCategoryDefinitionsOrErr returns the CheckListCategoryDefinitions value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderTemplateEdges) CheckListCategoryDefinitionsOrErr() ([]*CheckListCategoryDefinition, error) {
	if e.loadedTypes[1] {
		return e.CheckListCategoryDefinitions, nil
	}
	return nil, &NotLoadedError{edge: "check_list_category_definitions"}
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkOrderTemplateEdges) TypeOrErr() (*WorkOrderType, error) {
	if e.loadedTypes[2] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workordertype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*WorkOrderTemplate) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullString{}, // name
		&sql.NullString{}, // description
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*WorkOrderTemplate) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // work_order_template_type
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the WorkOrderTemplate fields.
func (wot *WorkOrderTemplate) assignValues(values ...interface{}) error {
	if m, n := len(values), len(workordertemplate.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	wot.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[0])
	} else if value.Valid {
		wot.Name = value.String
	}
	if value, ok := values[1].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[1])
	} else if value.Valid {
		wot.Description = value.String
	}
	values = values[2:]
	if len(values) == len(workordertemplate.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_template_type", value)
		} else if value.Valid {
			wot.work_order_template_type = new(int)
			*wot.work_order_template_type = int(value.Int64)
		}
	}
	return nil
}

// QueryPropertyTypes queries the property_types edge of the WorkOrderTemplate.
func (wot *WorkOrderTemplate) QueryPropertyTypes() *PropertyTypeQuery {
	return (&WorkOrderTemplateClient{config: wot.config}).QueryPropertyTypes(wot)
}

// QueryCheckListCategoryDefinitions queries the check_list_category_definitions edge of the WorkOrderTemplate.
func (wot *WorkOrderTemplate) QueryCheckListCategoryDefinitions() *CheckListCategoryDefinitionQuery {
	return (&WorkOrderTemplateClient{config: wot.config}).QueryCheckListCategoryDefinitions(wot)
}

// QueryType queries the type edge of the WorkOrderTemplate.
func (wot *WorkOrderTemplate) QueryType() *WorkOrderTypeQuery {
	return (&WorkOrderTemplateClient{config: wot.config}).QueryType(wot)
}

// Update returns a builder for updating this WorkOrderTemplate.
// Note that, you need to call WorkOrderTemplate.Unwrap() before calling this method, if this WorkOrderTemplate
// was returned from a transaction, and the transaction was committed or rolled back.
func (wot *WorkOrderTemplate) Update() *WorkOrderTemplateUpdateOne {
	return (&WorkOrderTemplateClient{config: wot.config}).UpdateOne(wot)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (wot *WorkOrderTemplate) Unwrap() *WorkOrderTemplate {
	tx, ok := wot.config.driver.(*txDriver)
	if !ok {
		panic("ent: WorkOrderTemplate is not a transactional entity")
	}
	wot.config.driver = tx.drv
	return wot
}

// String implements the fmt.Stringer.
func (wot *WorkOrderTemplate) String() string {
	var builder strings.Builder
	builder.WriteString("WorkOrderTemplate(")
	builder.WriteString(fmt.Sprintf("id=%v", wot.ID))
	builder.WriteString(", name=")
	builder.WriteString(wot.Name)
	builder.WriteString(", description=")
	builder.WriteString(wot.Description)
	builder.WriteByte(')')
	return builder.String()
}

// WorkOrderTemplates is a parsable slice of WorkOrderTemplate.
type WorkOrderTemplates []*WorkOrderTemplate

func (wot WorkOrderTemplates) config(cfg config) {
	for _i := range wot {
		wot[_i].config = cfg
	}
}
