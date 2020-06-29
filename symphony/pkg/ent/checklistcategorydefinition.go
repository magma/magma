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
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/workordertemplate"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
)

// CheckListCategoryDefinition is the model entity for the CheckListCategoryDefinition schema.
type CheckListCategoryDefinition struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CheckListCategoryDefinitionQuery when eager-loading is set.
	Edges                                               CheckListCategoryDefinitionEdges `json:"edges"`
	work_order_template_check_list_category_definitions *int
	work_order_type_check_list_category_definitions     *int
}

// CheckListCategoryDefinitionEdges holds the relations/edges for other nodes in the graph.
type CheckListCategoryDefinitionEdges struct {
	// CheckListItemDefinitions holds the value of the check_list_item_definitions edge.
	CheckListItemDefinitions []*CheckListItemDefinition
	// WorkOrderType holds the value of the work_order_type edge.
	WorkOrderType *WorkOrderType
	// WorkOrderTemplate holds the value of the work_order_template edge.
	WorkOrderTemplate *WorkOrderTemplate
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// CheckListItemDefinitionsOrErr returns the CheckListItemDefinitions value or an error if the edge
// was not loaded in eager-loading.
func (e CheckListCategoryDefinitionEdges) CheckListItemDefinitionsOrErr() ([]*CheckListItemDefinition, error) {
	if e.loadedTypes[0] {
		return e.CheckListItemDefinitions, nil
	}
	return nil, &NotLoadedError{edge: "check_list_item_definitions"}
}

// WorkOrderTypeOrErr returns the WorkOrderType value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CheckListCategoryDefinitionEdges) WorkOrderTypeOrErr() (*WorkOrderType, error) {
	if e.loadedTypes[1] {
		if e.WorkOrderType == nil {
			// The edge work_order_type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workordertype.Label}
		}
		return e.WorkOrderType, nil
	}
	return nil, &NotLoadedError{edge: "work_order_type"}
}

// WorkOrderTemplateOrErr returns the WorkOrderTemplate value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CheckListCategoryDefinitionEdges) WorkOrderTemplateOrErr() (*WorkOrderTemplate, error) {
	if e.loadedTypes[2] {
		if e.WorkOrderTemplate == nil {
			// The edge work_order_template was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workordertemplate.Label}
		}
		return e.WorkOrderTemplate, nil
	}
	return nil, &NotLoadedError{edge: "work_order_template"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CheckListCategoryDefinition) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // title
		&sql.NullString{}, // description
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*CheckListCategoryDefinition) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // work_order_template_check_list_category_definitions
		&sql.NullInt64{}, // work_order_type_check_list_category_definitions
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CheckListCategoryDefinition fields.
func (clcd *CheckListCategoryDefinition) assignValues(values ...interface{}) error {
	if m, n := len(values), len(checklistcategorydefinition.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	clcd.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		clcd.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		clcd.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field title", values[2])
	} else if value.Valid {
		clcd.Title = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[3])
	} else if value.Valid {
		clcd.Description = value.String
	}
	values = values[4:]
	if len(values) == len(checklistcategorydefinition.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_template_check_list_category_definitions", value)
		} else if value.Valid {
			clcd.work_order_template_check_list_category_definitions = new(int)
			*clcd.work_order_template_check_list_category_definitions = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_type_check_list_category_definitions", value)
		} else if value.Valid {
			clcd.work_order_type_check_list_category_definitions = new(int)
			*clcd.work_order_type_check_list_category_definitions = int(value.Int64)
		}
	}
	return nil
}

// QueryCheckListItemDefinitions queries the check_list_item_definitions edge of the CheckListCategoryDefinition.
func (clcd *CheckListCategoryDefinition) QueryCheckListItemDefinitions() *CheckListItemDefinitionQuery {
	return (&CheckListCategoryDefinitionClient{config: clcd.config}).QueryCheckListItemDefinitions(clcd)
}

// QueryWorkOrderType queries the work_order_type edge of the CheckListCategoryDefinition.
func (clcd *CheckListCategoryDefinition) QueryWorkOrderType() *WorkOrderTypeQuery {
	return (&CheckListCategoryDefinitionClient{config: clcd.config}).QueryWorkOrderType(clcd)
}

// QueryWorkOrderTemplate queries the work_order_template edge of the CheckListCategoryDefinition.
func (clcd *CheckListCategoryDefinition) QueryWorkOrderTemplate() *WorkOrderTemplateQuery {
	return (&CheckListCategoryDefinitionClient{config: clcd.config}).QueryWorkOrderTemplate(clcd)
}

// Update returns a builder for updating this CheckListCategoryDefinition.
// Note that, you need to call CheckListCategoryDefinition.Unwrap() before calling this method, if this CheckListCategoryDefinition
// was returned from a transaction, and the transaction was committed or rolled back.
func (clcd *CheckListCategoryDefinition) Update() *CheckListCategoryDefinitionUpdateOne {
	return (&CheckListCategoryDefinitionClient{config: clcd.config}).UpdateOne(clcd)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (clcd *CheckListCategoryDefinition) Unwrap() *CheckListCategoryDefinition {
	tx, ok := clcd.config.driver.(*txDriver)
	if !ok {
		panic("ent: CheckListCategoryDefinition is not a transactional entity")
	}
	clcd.config.driver = tx.drv
	return clcd
}

// String implements the fmt.Stringer.
func (clcd *CheckListCategoryDefinition) String() string {
	var builder strings.Builder
	builder.WriteString("CheckListCategoryDefinition(")
	builder.WriteString(fmt.Sprintf("id=%v", clcd.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(clcd.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(clcd.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", title=")
	builder.WriteString(clcd.Title)
	builder.WriteString(", description=")
	builder.WriteString(clcd.Description)
	builder.WriteByte(')')
	return builder.String()
}

// CheckListCategoryDefinitions is a parsable slice of CheckListCategoryDefinition.
type CheckListCategoryDefinitions []*CheckListCategoryDefinition

func (clcd CheckListCategoryDefinitions) config(cfg config) {
	for _i := range clcd {
		clcd[_i].config = cfg
	}
}
