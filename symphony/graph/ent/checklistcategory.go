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
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
)

// CheckListCategory is the model entity for the CheckListCategory schema.
type CheckListCategory struct {
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
	// The values are being populated by the CheckListCategoryQuery when eager-loading is set.
	Edges                                 CheckListCategoryEdges `json:"edges"`
	work_order_check_list_categories      *int
	work_order_type_check_list_categories *int
}

// CheckListCategoryEdges holds the relations/edges for other nodes in the graph.
type CheckListCategoryEdges struct {
	// CheckListItems holds the value of the check_list_items edge.
	CheckListItems []*CheckListItem
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// CheckListItemsOrErr returns the CheckListItems value or an error if the edge
// was not loaded in eager-loading.
func (e CheckListCategoryEdges) CheckListItemsOrErr() ([]*CheckListItem, error) {
	if e.loadedTypes[0] {
		return e.CheckListItems, nil
	}
	return nil, &NotLoadedError{edge: "check_list_items"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CheckListCategory) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // title
		&sql.NullString{}, // description
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*CheckListCategory) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // work_order_check_list_categories
		&sql.NullInt64{}, // work_order_type_check_list_categories
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CheckListCategory fields.
func (clc *CheckListCategory) assignValues(values ...interface{}) error {
	if m, n := len(values), len(checklistcategory.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	clc.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		clc.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		clc.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field title", values[2])
	} else if value.Valid {
		clc.Title = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[3])
	} else if value.Valid {
		clc.Description = value.String
	}
	values = values[4:]
	if len(values) == len(checklistcategory.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_check_list_categories", value)
		} else if value.Valid {
			clc.work_order_check_list_categories = new(int)
			*clc.work_order_check_list_categories = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_type_check_list_categories", value)
		} else if value.Valid {
			clc.work_order_type_check_list_categories = new(int)
			*clc.work_order_type_check_list_categories = int(value.Int64)
		}
	}
	return nil
}

// QueryCheckListItems queries the check_list_items edge of the CheckListCategory.
func (clc *CheckListCategory) QueryCheckListItems() *CheckListItemQuery {
	return (&CheckListCategoryClient{config: clc.config}).QueryCheckListItems(clc)
}

// Update returns a builder for updating this CheckListCategory.
// Note that, you need to call CheckListCategory.Unwrap() before calling this method, if this CheckListCategory
// was returned from a transaction, and the transaction was committed or rolled back.
func (clc *CheckListCategory) Update() *CheckListCategoryUpdateOne {
	return (&CheckListCategoryClient{config: clc.config}).UpdateOne(clc)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (clc *CheckListCategory) Unwrap() *CheckListCategory {
	tx, ok := clc.config.driver.(*txDriver)
	if !ok {
		panic("ent: CheckListCategory is not a transactional entity")
	}
	clc.config.driver = tx.drv
	return clc
}

// String implements the fmt.Stringer.
func (clc *CheckListCategory) String() string {
	var builder strings.Builder
	builder.WriteString("CheckListCategory(")
	builder.WriteString(fmt.Sprintf("id=%v", clc.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(clc.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(clc.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", title=")
	builder.WriteString(clc.Title)
	builder.WriteString(", description=")
	builder.WriteString(clc.Description)
	builder.WriteByte(')')
	return builder.String()
}

// CheckListCategories is a parsable slice of CheckListCategory.
type CheckListCategories []*CheckListCategory

func (clc CheckListCategories) config(cfg config) {
	for _i := range clc {
		clc[_i].config = cfg
	}
}
