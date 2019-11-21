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
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// ProjectTypeCreate is the builder for creating a ProjectType entity.
type ProjectTypeCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	description *string
	projects    map[string]struct{}
	properties  map[string]struct{}
	work_orders map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (ptc *ProjectTypeCreate) SetCreateTime(t time.Time) *ProjectTypeCreate {
	ptc.create_time = &t
	return ptc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ptc *ProjectTypeCreate) SetNillableCreateTime(t *time.Time) *ProjectTypeCreate {
	if t != nil {
		ptc.SetCreateTime(*t)
	}
	return ptc
}

// SetUpdateTime sets the update_time field.
func (ptc *ProjectTypeCreate) SetUpdateTime(t time.Time) *ProjectTypeCreate {
	ptc.update_time = &t
	return ptc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ptc *ProjectTypeCreate) SetNillableUpdateTime(t *time.Time) *ProjectTypeCreate {
	if t != nil {
		ptc.SetUpdateTime(*t)
	}
	return ptc
}

// SetName sets the name field.
func (ptc *ProjectTypeCreate) SetName(s string) *ProjectTypeCreate {
	ptc.name = &s
	return ptc
}

// SetDescription sets the description field.
func (ptc *ProjectTypeCreate) SetDescription(s string) *ProjectTypeCreate {
	ptc.description = &s
	return ptc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ptc *ProjectTypeCreate) SetNillableDescription(s *string) *ProjectTypeCreate {
	if s != nil {
		ptc.SetDescription(*s)
	}
	return ptc
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptc *ProjectTypeCreate) AddProjectIDs(ids ...string) *ProjectTypeCreate {
	if ptc.projects == nil {
		ptc.projects = make(map[string]struct{})
	}
	for i := range ids {
		ptc.projects[ids[i]] = struct{}{}
	}
	return ptc
}

// AddProjects adds the projects edges to Project.
func (ptc *ProjectTypeCreate) AddProjects(p ...*Project) *ProjectTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptc *ProjectTypeCreate) AddPropertyIDs(ids ...string) *ProjectTypeCreate {
	if ptc.properties == nil {
		ptc.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptc.properties[ids[i]] = struct{}{}
	}
	return ptc
}

// AddProperties adds the properties edges to PropertyType.
func (ptc *ProjectTypeCreate) AddProperties(p ...*PropertyType) *ProjectTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptc *ProjectTypeCreate) AddWorkOrderIDs(ids ...string) *ProjectTypeCreate {
	if ptc.work_orders == nil {
		ptc.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		ptc.work_orders[ids[i]] = struct{}{}
	}
	return ptc
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptc *ProjectTypeCreate) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptc.AddWorkOrderIDs(ids...)
}

// Save creates the ProjectType in the database.
func (ptc *ProjectTypeCreate) Save(ctx context.Context) (*ProjectType, error) {
	if ptc.create_time == nil {
		v := projecttype.DefaultCreateTime()
		ptc.create_time = &v
	}
	if ptc.update_time == nil {
		v := projecttype.DefaultUpdateTime()
		ptc.update_time = &v
	}
	if ptc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := projecttype.NameValidator(*ptc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	return ptc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (ptc *ProjectTypeCreate) SaveX(ctx context.Context) *ProjectType {
	v, err := ptc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ptc *ProjectTypeCreate) sqlSave(ctx context.Context) (*ProjectType, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ptc.driver.Dialect())
		pt      = &ProjectType{config: ptc.config}
	)
	tx, err := ptc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(projecttype.Table).Default()
	if value := ptc.create_time; value != nil {
		insert.Set(projecttype.FieldCreateTime, *value)
		pt.CreateTime = *value
	}
	if value := ptc.update_time; value != nil {
		insert.Set(projecttype.FieldUpdateTime, *value)
		pt.UpdateTime = *value
	}
	if value := ptc.name; value != nil {
		insert.Set(projecttype.FieldName, *value)
		pt.Name = *value
	}
	if value := ptc.description; value != nil {
		insert.Set(projecttype.FieldDescription, *value)
		pt.Description = value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(projecttype.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	pt.ID = strconv.FormatInt(id, 10)
	if len(ptc.projects) > 0 {
		p := sql.P()
		for eid := range ptc.projects {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(project.FieldID, eid)
		}
		query, args := builder.Update(projecttype.ProjectsTable).
			Set(projecttype.ProjectsColumn, id).
			Where(sql.And(p, sql.IsNull(projecttype.ProjectsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ptc.projects) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"projects\" %v already connected to a different \"ProjectType\"", keys(ptc.projects))})
		}
	}
	if len(ptc.properties) > 0 {
		p := sql.P()
		for eid := range ptc.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(propertytype.FieldID, eid)
		}
		query, args := builder.Update(projecttype.PropertiesTable).
			Set(projecttype.PropertiesColumn, id).
			Where(sql.And(p, sql.IsNull(projecttype.PropertiesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ptc.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"ProjectType\"", keys(ptc.properties))})
		}
	}
	if len(ptc.work_orders) > 0 {
		p := sql.P()
		for eid := range ptc.work_orders {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(workorderdefinition.FieldID, eid)
		}
		query, args := builder.Update(projecttype.WorkOrdersTable).
			Set(projecttype.WorkOrdersColumn, id).
			Where(sql.And(p, sql.IsNull(projecttype.WorkOrdersColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ptc.work_orders) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"ProjectType\"", keys(ptc.work_orders))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return pt, nil
}
