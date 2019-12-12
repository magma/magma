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
	"github.com/facebookincubator/symphony/graph/ent/project"
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
		res     sql.Result
		builder = sql.Dialect(pc.driver.Dialect())
		pr      = &Project{config: pc.config}
	)
	tx, err := pc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(project.Table).Default()
	if value := pc.create_time; value != nil {
		insert.Set(project.FieldCreateTime, *value)
		pr.CreateTime = *value
	}
	if value := pc.update_time; value != nil {
		insert.Set(project.FieldUpdateTime, *value)
		pr.UpdateTime = *value
	}
	if value := pc.name; value != nil {
		insert.Set(project.FieldName, *value)
		pr.Name = *value
	}
	if value := pc.description; value != nil {
		insert.Set(project.FieldDescription, *value)
		pr.Description = value
	}
	if value := pc.creator; value != nil {
		insert.Set(project.FieldCreator, *value)
		pr.Creator = value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(project.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	pr.ID = strconv.FormatInt(id, 10)
	if len(pc._type) > 0 {
		for eid := range pc._type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(project.TypeTable).
				Set(project.TypeColumn, eid).
				Where(sql.EQ(project.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.location) > 0 {
		for eid := range pc.location {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(project.LocationTable).
				Set(project.LocationColumn, eid).
				Where(sql.EQ(project.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(pc.work_orders) > 0 {
		p := sql.P()
		for eid := range pc.work_orders {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(workorder.FieldID, eid)
		}
		query, args := builder.Update(project.WorkOrdersTable).
			Set(project.WorkOrdersColumn, id).
			Where(sql.And(p, sql.IsNull(project.WorkOrdersColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(pc.work_orders) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"Project\"", keys(pc.work_orders))})
		}
	}
	if len(pc.properties) > 0 {
		p := sql.P()
		for eid := range pc.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(property.FieldID, eid)
		}
		query, args := builder.Update(project.PropertiesTable).
			Set(project.PropertiesColumn, id).
			Where(sql.And(p, sql.IsNull(project.PropertiesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(pc.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Project\"", keys(pc.properties))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return pr, nil
}
