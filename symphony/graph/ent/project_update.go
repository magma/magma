// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// ProjectUpdate is the builder for updating Project entities.
type ProjectUpdate struct {
	config

	update_time       *time.Time
	name              *string
	description       *string
	cleardescription  bool
	creator           *string
	clearcreator      bool
	_type             map[string]struct{}
	location          map[string]struct{}
	comments          map[string]struct{}
	work_orders       map[string]struct{}
	properties        map[string]struct{}
	clearedType       bool
	clearedLocation   bool
	removedComments   map[string]struct{}
	removedWorkOrders map[string]struct{}
	removedProperties map[string]struct{}
	predicates        []predicate.Project
}

// Where adds a new predicate for the builder.
func (pu *ProjectUpdate) Where(ps ...predicate.Project) *ProjectUpdate {
	pu.predicates = append(pu.predicates, ps...)
	return pu
}

// SetName sets the name field.
func (pu *ProjectUpdate) SetName(s string) *ProjectUpdate {
	pu.name = &s
	return pu
}

// SetDescription sets the description field.
func (pu *ProjectUpdate) SetDescription(s string) *ProjectUpdate {
	pu.description = &s
	return pu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (pu *ProjectUpdate) SetNillableDescription(s *string) *ProjectUpdate {
	if s != nil {
		pu.SetDescription(*s)
	}
	return pu
}

// ClearDescription clears the value of description.
func (pu *ProjectUpdate) ClearDescription() *ProjectUpdate {
	pu.description = nil
	pu.cleardescription = true
	return pu
}

// SetCreator sets the creator field.
func (pu *ProjectUpdate) SetCreator(s string) *ProjectUpdate {
	pu.creator = &s
	return pu
}

// SetNillableCreator sets the creator field if the given value is not nil.
func (pu *ProjectUpdate) SetNillableCreator(s *string) *ProjectUpdate {
	if s != nil {
		pu.SetCreator(*s)
	}
	return pu
}

// ClearCreator clears the value of creator.
func (pu *ProjectUpdate) ClearCreator() *ProjectUpdate {
	pu.creator = nil
	pu.clearcreator = true
	return pu
}

// SetTypeID sets the type edge to ProjectType by id.
func (pu *ProjectUpdate) SetTypeID(id string) *ProjectUpdate {
	if pu._type == nil {
		pu._type = make(map[string]struct{})
	}
	pu._type[id] = struct{}{}
	return pu
}

// SetType sets the type edge to ProjectType.
func (pu *ProjectUpdate) SetType(p *ProjectType) *ProjectUpdate {
	return pu.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (pu *ProjectUpdate) SetLocationID(id string) *ProjectUpdate {
	if pu.location == nil {
		pu.location = make(map[string]struct{})
	}
	pu.location[id] = struct{}{}
	return pu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (pu *ProjectUpdate) SetNillableLocationID(id *string) *ProjectUpdate {
	if id != nil {
		pu = pu.SetLocationID(*id)
	}
	return pu
}

// SetLocation sets the location edge to Location.
func (pu *ProjectUpdate) SetLocation(l *Location) *ProjectUpdate {
	return pu.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (pu *ProjectUpdate) AddCommentIDs(ids ...string) *ProjectUpdate {
	if pu.comments == nil {
		pu.comments = make(map[string]struct{})
	}
	for i := range ids {
		pu.comments[ids[i]] = struct{}{}
	}
	return pu
}

// AddComments adds the comments edges to Comment.
func (pu *ProjectUpdate) AddComments(c ...*Comment) *ProjectUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pu.AddCommentIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (pu *ProjectUpdate) AddWorkOrderIDs(ids ...string) *ProjectUpdate {
	if pu.work_orders == nil {
		pu.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		pu.work_orders[ids[i]] = struct{}{}
	}
	return pu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (pu *ProjectUpdate) AddWorkOrders(w ...*WorkOrder) *ProjectUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return pu.AddWorkOrderIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (pu *ProjectUpdate) AddPropertyIDs(ids ...string) *ProjectUpdate {
	if pu.properties == nil {
		pu.properties = make(map[string]struct{})
	}
	for i := range ids {
		pu.properties[ids[i]] = struct{}{}
	}
	return pu
}

// AddProperties adds the properties edges to Property.
func (pu *ProjectUpdate) AddProperties(p ...*Property) *ProjectUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.AddPropertyIDs(ids...)
}

// ClearType clears the type edge to ProjectType.
func (pu *ProjectUpdate) ClearType() *ProjectUpdate {
	pu.clearedType = true
	return pu
}

// ClearLocation clears the location edge to Location.
func (pu *ProjectUpdate) ClearLocation() *ProjectUpdate {
	pu.clearedLocation = true
	return pu
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (pu *ProjectUpdate) RemoveCommentIDs(ids ...string) *ProjectUpdate {
	if pu.removedComments == nil {
		pu.removedComments = make(map[string]struct{})
	}
	for i := range ids {
		pu.removedComments[ids[i]] = struct{}{}
	}
	return pu
}

// RemoveComments removes comments edges to Comment.
func (pu *ProjectUpdate) RemoveComments(c ...*Comment) *ProjectUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pu.RemoveCommentIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (pu *ProjectUpdate) RemoveWorkOrderIDs(ids ...string) *ProjectUpdate {
	if pu.removedWorkOrders == nil {
		pu.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		pu.removedWorkOrders[ids[i]] = struct{}{}
	}
	return pu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (pu *ProjectUpdate) RemoveWorkOrders(w ...*WorkOrder) *ProjectUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return pu.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (pu *ProjectUpdate) RemovePropertyIDs(ids ...string) *ProjectUpdate {
	if pu.removedProperties == nil {
		pu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		pu.removedProperties[ids[i]] = struct{}{}
	}
	return pu
}

// RemoveProperties removes properties edges to Property.
func (pu *ProjectUpdate) RemoveProperties(p ...*Property) *ProjectUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.RemovePropertyIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (pu *ProjectUpdate) Save(ctx context.Context) (int, error) {
	if pu.update_time == nil {
		v := project.UpdateDefaultUpdateTime()
		pu.update_time = &v
	}
	if pu.name != nil {
		if err := project.NameValidator(*pu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(pu._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if pu.clearedType && pu._type == nil {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(pu.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return pu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProjectUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProjectUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProjectUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pu *ProjectUpdate) sqlSave(ctx context.Context) (n int, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: project.FieldID,
			},
		},
	}
	if ps := pu.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := pu.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: project.FieldUpdateTime,
		})
	}
	if value := pu.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldName,
		})
	}
	if value := pu.description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldDescription,
		})
	}
	if pu.cleardescription {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: project.FieldDescription,
		})
	}
	if value := pu.creator; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldCreator,
		})
	}
	if pu.clearcreator {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: project.FieldCreator,
		})
	}
	if pu.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: projecttype.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := pu._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: projecttype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if pu.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := pu.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := pu.removedComments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := pu.comments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := pu.removedWorkOrders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := pu.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := pu.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := pu.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ProjectUpdateOne is the builder for updating a single Project entity.
type ProjectUpdateOne struct {
	config
	id string

	update_time       *time.Time
	name              *string
	description       *string
	cleardescription  bool
	creator           *string
	clearcreator      bool
	_type             map[string]struct{}
	location          map[string]struct{}
	comments          map[string]struct{}
	work_orders       map[string]struct{}
	properties        map[string]struct{}
	clearedType       bool
	clearedLocation   bool
	removedComments   map[string]struct{}
	removedWorkOrders map[string]struct{}
	removedProperties map[string]struct{}
}

// SetName sets the name field.
func (puo *ProjectUpdateOne) SetName(s string) *ProjectUpdateOne {
	puo.name = &s
	return puo
}

// SetDescription sets the description field.
func (puo *ProjectUpdateOne) SetDescription(s string) *ProjectUpdateOne {
	puo.description = &s
	return puo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableDescription(s *string) *ProjectUpdateOne {
	if s != nil {
		puo.SetDescription(*s)
	}
	return puo
}

// ClearDescription clears the value of description.
func (puo *ProjectUpdateOne) ClearDescription() *ProjectUpdateOne {
	puo.description = nil
	puo.cleardescription = true
	return puo
}

// SetCreator sets the creator field.
func (puo *ProjectUpdateOne) SetCreator(s string) *ProjectUpdateOne {
	puo.creator = &s
	return puo
}

// SetNillableCreator sets the creator field if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableCreator(s *string) *ProjectUpdateOne {
	if s != nil {
		puo.SetCreator(*s)
	}
	return puo
}

// ClearCreator clears the value of creator.
func (puo *ProjectUpdateOne) ClearCreator() *ProjectUpdateOne {
	puo.creator = nil
	puo.clearcreator = true
	return puo
}

// SetTypeID sets the type edge to ProjectType by id.
func (puo *ProjectUpdateOne) SetTypeID(id string) *ProjectUpdateOne {
	if puo._type == nil {
		puo._type = make(map[string]struct{})
	}
	puo._type[id] = struct{}{}
	return puo
}

// SetType sets the type edge to ProjectType.
func (puo *ProjectUpdateOne) SetType(p *ProjectType) *ProjectUpdateOne {
	return puo.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (puo *ProjectUpdateOne) SetLocationID(id string) *ProjectUpdateOne {
	if puo.location == nil {
		puo.location = make(map[string]struct{})
	}
	puo.location[id] = struct{}{}
	return puo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableLocationID(id *string) *ProjectUpdateOne {
	if id != nil {
		puo = puo.SetLocationID(*id)
	}
	return puo
}

// SetLocation sets the location edge to Location.
func (puo *ProjectUpdateOne) SetLocation(l *Location) *ProjectUpdateOne {
	return puo.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (puo *ProjectUpdateOne) AddCommentIDs(ids ...string) *ProjectUpdateOne {
	if puo.comments == nil {
		puo.comments = make(map[string]struct{})
	}
	for i := range ids {
		puo.comments[ids[i]] = struct{}{}
	}
	return puo
}

// AddComments adds the comments edges to Comment.
func (puo *ProjectUpdateOne) AddComments(c ...*Comment) *ProjectUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return puo.AddCommentIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (puo *ProjectUpdateOne) AddWorkOrderIDs(ids ...string) *ProjectUpdateOne {
	if puo.work_orders == nil {
		puo.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		puo.work_orders[ids[i]] = struct{}{}
	}
	return puo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (puo *ProjectUpdateOne) AddWorkOrders(w ...*WorkOrder) *ProjectUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return puo.AddWorkOrderIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (puo *ProjectUpdateOne) AddPropertyIDs(ids ...string) *ProjectUpdateOne {
	if puo.properties == nil {
		puo.properties = make(map[string]struct{})
	}
	for i := range ids {
		puo.properties[ids[i]] = struct{}{}
	}
	return puo
}

// AddProperties adds the properties edges to Property.
func (puo *ProjectUpdateOne) AddProperties(p ...*Property) *ProjectUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.AddPropertyIDs(ids...)
}

// ClearType clears the type edge to ProjectType.
func (puo *ProjectUpdateOne) ClearType() *ProjectUpdateOne {
	puo.clearedType = true
	return puo
}

// ClearLocation clears the location edge to Location.
func (puo *ProjectUpdateOne) ClearLocation() *ProjectUpdateOne {
	puo.clearedLocation = true
	return puo
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (puo *ProjectUpdateOne) RemoveCommentIDs(ids ...string) *ProjectUpdateOne {
	if puo.removedComments == nil {
		puo.removedComments = make(map[string]struct{})
	}
	for i := range ids {
		puo.removedComments[ids[i]] = struct{}{}
	}
	return puo
}

// RemoveComments removes comments edges to Comment.
func (puo *ProjectUpdateOne) RemoveComments(c ...*Comment) *ProjectUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return puo.RemoveCommentIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (puo *ProjectUpdateOne) RemoveWorkOrderIDs(ids ...string) *ProjectUpdateOne {
	if puo.removedWorkOrders == nil {
		puo.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		puo.removedWorkOrders[ids[i]] = struct{}{}
	}
	return puo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (puo *ProjectUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *ProjectUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return puo.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (puo *ProjectUpdateOne) RemovePropertyIDs(ids ...string) *ProjectUpdateOne {
	if puo.removedProperties == nil {
		puo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		puo.removedProperties[ids[i]] = struct{}{}
	}
	return puo
}

// RemoveProperties removes properties edges to Property.
func (puo *ProjectUpdateOne) RemoveProperties(p ...*Property) *ProjectUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.RemovePropertyIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (puo *ProjectUpdateOne) Save(ctx context.Context) (*Project, error) {
	if puo.update_time == nil {
		v := project.UpdateDefaultUpdateTime()
		puo.update_time = &v
	}
	if puo.name != nil {
		if err := project.NameValidator(*puo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(puo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if puo.clearedType && puo._type == nil {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(puo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return puo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProjectUpdateOne) SaveX(ctx context.Context) *Project {
	pr, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return pr
}

// Exec executes the query on the entity.
func (puo *ProjectUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProjectUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (puo *ProjectUpdateOne) sqlSave(ctx context.Context) (pr *Project, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  puo.id,
				Type:   field.TypeString,
				Column: project.FieldID,
			},
		},
	}
	if value := puo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: project.FieldUpdateTime,
		})
	}
	if value := puo.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldName,
		})
	}
	if value := puo.description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldDescription,
		})
	}
	if puo.cleardescription {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: project.FieldDescription,
		})
	}
	if value := puo.creator; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldCreator,
		})
	}
	if puo.clearcreator {
		spec.Fields.Clear = append(spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: project.FieldCreator,
		})
	}
	if puo.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: projecttype.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := puo._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TypeTable,
			Columns: []string{project.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: projecttype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if puo.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := puo.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   project.LocationTable,
			Columns: []string{project.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := puo.removedComments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := puo.comments; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.CommentsTable,
			Columns: []string{project.CommentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: comment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := puo.removedWorkOrders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := puo.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.WorkOrdersTable,
			Columns: []string{project.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if nodes := puo.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := puo.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.PropertiesTable,
			Columns: []string{project.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	pr = &Project{config: puo.config}
	spec.Assign = pr.assignValues
	spec.ScanValues = pr.scanValues()
	if err = sqlgraph.UpdateNode(ctx, puo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return pr, nil
}
