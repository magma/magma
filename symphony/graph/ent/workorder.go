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
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrder is the model entity for the WorkOrder schema.
type WorkOrder struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// Priority holds the value of the "priority" field.
	Priority string `json:"priority,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// OwnerName holds the value of the "owner_name" field.
	OwnerName string `json:"owner_name,omitempty"`
	// InstallDate holds the value of the "install_date" field.
	InstallDate time.Time `json:"install_date,omitempty"`
	// CreationDate holds the value of the "creation_date" field.
	CreationDate time.Time `json:"creation_date,omitempty"`
	// Assignee holds the value of the "assignee" field.
	Assignee string `json:"assignee,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the WorkOrderQuery when eager-loading is set.
	Edges                 WorkOrderEdges `json:"edges"`
	project_work_orders   *int
	work_order_type       *int
	work_order_location   *int
	work_order_technician *int
}

// WorkOrderEdges holds the relations/edges for other nodes in the graph.
type WorkOrderEdges struct {
	// Type holds the value of the type edge.
	Type *WorkOrderType
	// Equipment holds the value of the equipment edge.
	Equipment []*Equipment
	// Links holds the value of the links edge.
	Links []*Link
	// Files holds the value of the files edge.
	Files []*File
	// Hyperlinks holds the value of the hyperlinks edge.
	Hyperlinks []*Hyperlink
	// Location holds the value of the location edge.
	Location *Location
	// Comments holds the value of the comments edge.
	Comments []*Comment
	// Properties holds the value of the properties edge.
	Properties []*Property
	// CheckListCategories holds the value of the check_list_categories edge.
	CheckListCategories []*CheckListCategory
	// CheckListItems holds the value of the check_list_items edge.
	CheckListItems []*CheckListItem
	// Technician holds the value of the technician edge.
	Technician *Technician
	// Project holds the value of the project edge.
	Project *Project
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [12]bool
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkOrderEdges) TypeOrErr() (*WorkOrderType, error) {
	if e.loadedTypes[0] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workordertype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// EquipmentOrErr returns the Equipment value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) EquipmentOrErr() ([]*Equipment, error) {
	if e.loadedTypes[1] {
		return e.Equipment, nil
	}
	return nil, &NotLoadedError{edge: "equipment"}
}

// LinksOrErr returns the Links value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) LinksOrErr() ([]*Link, error) {
	if e.loadedTypes[2] {
		return e.Links, nil
	}
	return nil, &NotLoadedError{edge: "links"}
}

// FilesOrErr returns the Files value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) FilesOrErr() ([]*File, error) {
	if e.loadedTypes[3] {
		return e.Files, nil
	}
	return nil, &NotLoadedError{edge: "files"}
}

// HyperlinksOrErr returns the Hyperlinks value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) HyperlinksOrErr() ([]*Hyperlink, error) {
	if e.loadedTypes[4] {
		return e.Hyperlinks, nil
	}
	return nil, &NotLoadedError{edge: "hyperlinks"}
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkOrderEdges) LocationOrErr() (*Location, error) {
	if e.loadedTypes[5] {
		if e.Location == nil {
			// The edge location was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Location, nil
	}
	return nil, &NotLoadedError{edge: "location"}
}

// CommentsOrErr returns the Comments value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) CommentsOrErr() ([]*Comment, error) {
	if e.loadedTypes[6] {
		return e.Comments, nil
	}
	return nil, &NotLoadedError{edge: "comments"}
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[7] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// CheckListCategoriesOrErr returns the CheckListCategories value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) CheckListCategoriesOrErr() ([]*CheckListCategory, error) {
	if e.loadedTypes[8] {
		return e.CheckListCategories, nil
	}
	return nil, &NotLoadedError{edge: "check_list_categories"}
}

// CheckListItemsOrErr returns the CheckListItems value or an error if the edge
// was not loaded in eager-loading.
func (e WorkOrderEdges) CheckListItemsOrErr() ([]*CheckListItem, error) {
	if e.loadedTypes[9] {
		return e.CheckListItems, nil
	}
	return nil, &NotLoadedError{edge: "check_list_items"}
}

// TechnicianOrErr returns the Technician value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkOrderEdges) TechnicianOrErr() (*Technician, error) {
	if e.loadedTypes[10] {
		if e.Technician == nil {
			// The edge technician was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: technician.Label}
		}
		return e.Technician, nil
	}
	return nil, &NotLoadedError{edge: "technician"}
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkOrderEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[11] {
		if e.Project == nil {
			// The edge project was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*WorkOrder) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // status
		&sql.NullString{}, // priority
		&sql.NullString{}, // description
		&sql.NullString{}, // owner_name
		&sql.NullTime{},   // install_date
		&sql.NullTime{},   // creation_date
		&sql.NullString{}, // assignee
		&sql.NullInt64{},  // index
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*WorkOrder) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // project_work_orders
		&sql.NullInt64{}, // work_order_type
		&sql.NullInt64{}, // work_order_location
		&sql.NullInt64{}, // work_order_technician
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the WorkOrder fields.
func (wo *WorkOrder) assignValues(values ...interface{}) error {
	if m, n := len(values), len(workorder.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	wo.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		wo.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		wo.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		wo.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field status", values[3])
	} else if value.Valid {
		wo.Status = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field priority", values[4])
	} else if value.Valid {
		wo.Priority = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[5])
	} else if value.Valid {
		wo.Description = value.String
	}
	if value, ok := values[6].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field owner_name", values[6])
	} else if value.Valid {
		wo.OwnerName = value.String
	}
	if value, ok := values[7].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field install_date", values[7])
	} else if value.Valid {
		wo.InstallDate = value.Time
	}
	if value, ok := values[8].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field creation_date", values[8])
	} else if value.Valid {
		wo.CreationDate = value.Time
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field assignee", values[9])
	} else if value.Valid {
		wo.Assignee = value.String
	}
	if value, ok := values[10].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field index", values[10])
	} else if value.Valid {
		wo.Index = int(value.Int64)
	}
	values = values[11:]
	if len(values) == len(workorder.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_work_orders", value)
		} else if value.Valid {
			wo.project_work_orders = new(int)
			*wo.project_work_orders = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_type", value)
		} else if value.Valid {
			wo.work_order_type = new(int)
			*wo.work_order_type = int(value.Int64)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_location", value)
		} else if value.Valid {
			wo.work_order_location = new(int)
			*wo.work_order_location = int(value.Int64)
		}
		if value, ok := values[3].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_technician", value)
		} else if value.Valid {
			wo.work_order_technician = new(int)
			*wo.work_order_technician = int(value.Int64)
		}
	}
	return nil
}

// QueryType queries the type edge of the WorkOrder.
func (wo *WorkOrder) QueryType() *WorkOrderTypeQuery {
	return (&WorkOrderClient{wo.config}).QueryType(wo)
}

// QueryEquipment queries the equipment edge of the WorkOrder.
func (wo *WorkOrder) QueryEquipment() *EquipmentQuery {
	return (&WorkOrderClient{wo.config}).QueryEquipment(wo)
}

// QueryLinks queries the links edge of the WorkOrder.
func (wo *WorkOrder) QueryLinks() *LinkQuery {
	return (&WorkOrderClient{wo.config}).QueryLinks(wo)
}

// QueryFiles queries the files edge of the WorkOrder.
func (wo *WorkOrder) QueryFiles() *FileQuery {
	return (&WorkOrderClient{wo.config}).QueryFiles(wo)
}

// QueryHyperlinks queries the hyperlinks edge of the WorkOrder.
func (wo *WorkOrder) QueryHyperlinks() *HyperlinkQuery {
	return (&WorkOrderClient{wo.config}).QueryHyperlinks(wo)
}

// QueryLocation queries the location edge of the WorkOrder.
func (wo *WorkOrder) QueryLocation() *LocationQuery {
	return (&WorkOrderClient{wo.config}).QueryLocation(wo)
}

// QueryComments queries the comments edge of the WorkOrder.
func (wo *WorkOrder) QueryComments() *CommentQuery {
	return (&WorkOrderClient{wo.config}).QueryComments(wo)
}

// QueryProperties queries the properties edge of the WorkOrder.
func (wo *WorkOrder) QueryProperties() *PropertyQuery {
	return (&WorkOrderClient{wo.config}).QueryProperties(wo)
}

// QueryCheckListCategories queries the check_list_categories edge of the WorkOrder.
func (wo *WorkOrder) QueryCheckListCategories() *CheckListCategoryQuery {
	return (&WorkOrderClient{wo.config}).QueryCheckListCategories(wo)
}

// QueryCheckListItems queries the check_list_items edge of the WorkOrder.
func (wo *WorkOrder) QueryCheckListItems() *CheckListItemQuery {
	return (&WorkOrderClient{wo.config}).QueryCheckListItems(wo)
}

// QueryTechnician queries the technician edge of the WorkOrder.
func (wo *WorkOrder) QueryTechnician() *TechnicianQuery {
	return (&WorkOrderClient{wo.config}).QueryTechnician(wo)
}

// QueryProject queries the project edge of the WorkOrder.
func (wo *WorkOrder) QueryProject() *ProjectQuery {
	return (&WorkOrderClient{wo.config}).QueryProject(wo)
}

// Update returns a builder for updating this WorkOrder.
// Note that, you need to call WorkOrder.Unwrap() before calling this method, if this WorkOrder
// was returned from a transaction, and the transaction was committed or rolled back.
func (wo *WorkOrder) Update() *WorkOrderUpdateOne {
	return (&WorkOrderClient{wo.config}).UpdateOne(wo)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (wo *WorkOrder) Unwrap() *WorkOrder {
	tx, ok := wo.config.driver.(*txDriver)
	if !ok {
		panic("ent: WorkOrder is not a transactional entity")
	}
	wo.config.driver = tx.drv
	return wo
}

// String implements the fmt.Stringer.
func (wo *WorkOrder) String() string {
	var builder strings.Builder
	builder.WriteString("WorkOrder(")
	builder.WriteString(fmt.Sprintf("id=%v", wo.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(wo.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(wo.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(wo.Name)
	builder.WriteString(", status=")
	builder.WriteString(wo.Status)
	builder.WriteString(", priority=")
	builder.WriteString(wo.Priority)
	builder.WriteString(", description=")
	builder.WriteString(wo.Description)
	builder.WriteString(", owner_name=")
	builder.WriteString(wo.OwnerName)
	builder.WriteString(", install_date=")
	builder.WriteString(wo.InstallDate.Format(time.ANSIC))
	builder.WriteString(", creation_date=")
	builder.WriteString(wo.CreationDate.Format(time.ANSIC))
	builder.WriteString(", assignee=")
	builder.WriteString(wo.Assignee)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", wo.Index))
	builder.WriteByte(')')
	return builder.String()
}

// WorkOrders is a parsable slice of WorkOrder.
type WorkOrders []*WorkOrder

func (wo WorkOrders) config(cfg config) {
	for _i := range wo {
		wo[_i].config = cfg
	}
}
