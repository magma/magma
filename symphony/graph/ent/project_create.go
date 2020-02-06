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

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// ProjectCreate is the builder for creating a Project entity.
type ProjectCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	description *string
	creator     *string
	_type       map[string]struct{}
	location    map[string]struct{}
	comments    map[string]struct{}
	work_orders map[string]struct{}
	properties  map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (pc *ProjectCreate) SetCreateTime(t time.Time) *ProjectCreate {
	pc.create_time = &t
	return pc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (pc *ProjectCreate) SetNillableCreateTime(t *time.Time) *ProjectCreate {
	if t != nil {
		pc.SetCreateTime(*t)
	}
	return pc
}

// SetUpdateTime sets the update_time field.
func (pc *ProjectCreate) SetUpdateTime(t time.Time) *ProjectCreate {
	pc.update_time = &t
	return pc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (pc *ProjectCreate) SetNillableUpdateTime(t *time.Time) *ProjectCreate {
	if t != nil {
		pc.SetUpdateTime(*t)
	}
	return pc
}

// SetName sets the name field.
func (pc *ProjectCreate) SetName(s string) *ProjectCreate {
	pc.name = &s
	return pc
}

// SetDescription sets the description field.
func (pc *ProjectCreate) SetDescription(s string) *ProjectCreate {
	pc.description = &s
	return pc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (pc *ProjectCreate) SetNillableDescription(s *string) *ProjectCreate {
	if s != nil {
		pc.SetDescription(*s)
	}
	return pc
}

// SetCreator sets the creator field.
func (pc *ProjectCreate) SetCreator(s string) *ProjectCreate {
	pc.creator = &s
	return pc
}

// SetNillableCreator sets the creator field if the given value is not nil.
func (pc *ProjectCreate) SetNillableCreator(s *string) *ProjectCreate {
	if s != nil {
		pc.SetCreator(*s)
	}
	return pc
}

// SetTypeID sets the type edge to ProjectType by id.
func (pc *ProjectCreate) SetTypeID(id string) *ProjectCreate {
	if pc._type == nil {
		pc._type = make(map[string]struct{})
	}
	pc._type[id] = struct{}{}
	return pc
}

// SetType sets the type edge to ProjectType.
func (pc *ProjectCreate) SetType(p *ProjectType) *ProjectCreate {
	return pc.SetTypeID(p.ID)
}

// SetLocationID sets the location edge to Location by id.
func (pc *ProjectCreate) SetLocationID(id string) *ProjectCreate {
	if pc.location == nil {
		pc.location = make(map[string]struct{})
	}
	pc.location[id] = struct{}{}
	return pc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (pc *ProjectCreate) SetNillableLocationID(id *string) *ProjectCreate {
	if id != nil {
		pc = pc.SetLocationID(*id)
	}
	return pc
}

// SetLocation sets the location edge to Location.
func (pc *ProjectCreate) SetLocation(l *Location) *ProjectCreate {
	return pc.SetLocationID(l.ID)
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (pc *ProjectCreate) AddCommentIDs(ids ...string) *ProjectCreate {
	if pc.comments == nil {
		pc.comments = make(map[string]struct{})
	}
	for i := range ids {
		pc.comments[ids[i]] = struct{}{}
	}
	return pc
}

// AddComments adds the comments edges to Comment.
func (pc *ProjectCreate) AddComments(c ...*Comment) *ProjectCreate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pc.AddCommentIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (pc *ProjectCreate) AddWorkOrderIDs(ids ...string) *ProjectCreate {
	if pc.work_orders == nil {
		pc.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		pc.work_orders[ids[i]] = struct{}{}
	}
	return pc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (pc *ProjectCreate) AddWorkOrders(w ...*WorkOrder) *ProjectCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return pc.AddWorkOrderIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (pc *ProjectCreate) AddPropertyIDs(ids ...string) *ProjectCreate {
	if pc.properties == nil {
		pc.properties = make(map[string]struct{})
	}
	for i := range ids {
		pc.properties[ids[i]] = struct{}{}
	}
	return pc
}

// AddProperties adds the properties edges to Property.
func (pc *ProjectCreate) AddProperties(p ...*Property) *ProjectCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pc.AddPropertyIDs(ids...)
}

// Save creates the Project in the database.
func (pc *ProjectCreate) Save(ctx context.Context) (*Project, error) {
	if pc.create_time == nil {
		v := project.DefaultCreateTime()
		pc.create_time = &v
	}
	if pc.update_time == nil {
		v := project.DefaultUpdateTime()
		pc.update_time = &v
	}
	if pc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := project.NameValidator(*pc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if len(pc._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if pc._type == nil {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	if len(pc.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return pc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *ProjectCreate) SaveX(ctx context.Context) *Project {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pc *ProjectCreate) sqlSave(ctx context.Context) (*Project, error) {
	var (
		pr    = &Project{config: pc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: project.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: project.FieldID,
			},
		}
	)
	if value := pc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: project.FieldCreateTime,
		})
		pr.CreateTime = *value
	}
	if value := pc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: project.FieldUpdateTime,
		})
		pr.UpdateTime = *value
	}
	if value := pc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldName,
		})
		pr.Name = *value
	}
	if value := pc.description; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldDescription,
		})
		pr.Description = value
	}
	if value := pc.creator; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: project.FieldCreator,
		})
		pr.Creator = value
	}
	if nodes := pc._type; len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.location; len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.comments; len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.work_orders; len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.properties; len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pr.ID = strconv.FormatInt(id, 10)
	return pr, nil
}
