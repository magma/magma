// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// ProjectTypeUpdate is the builder for updating ProjectType entities.
type ProjectTypeUpdate struct {
	config

	update_time       *time.Time
	name              *string
	description       *string
	cleardescription  bool
	projects          map[string]struct{}
	properties        map[string]struct{}
	work_orders       map[string]struct{}
	removedProjects   map[string]struct{}
	removedProperties map[string]struct{}
	removedWorkOrders map[string]struct{}
	predicates        []predicate.ProjectType
}

// Where adds a new predicate for the builder.
func (ptu *ProjectTypeUpdate) Where(ps ...predicate.ProjectType) *ProjectTypeUpdate {
	ptu.predicates = append(ptu.predicates, ps...)
	return ptu
}

// SetName sets the name field.
func (ptu *ProjectTypeUpdate) SetName(s string) *ProjectTypeUpdate {
	ptu.name = &s
	return ptu
}

// SetDescription sets the description field.
func (ptu *ProjectTypeUpdate) SetDescription(s string) *ProjectTypeUpdate {
	ptu.description = &s
	return ptu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ptu *ProjectTypeUpdate) SetNillableDescription(s *string) *ProjectTypeUpdate {
	if s != nil {
		ptu.SetDescription(*s)
	}
	return ptu
}

// ClearDescription clears the value of description.
func (ptu *ProjectTypeUpdate) ClearDescription() *ProjectTypeUpdate {
	ptu.description = nil
	ptu.cleardescription = true
	return ptu
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptu *ProjectTypeUpdate) AddProjectIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.projects == nil {
		ptu.projects = make(map[string]struct{})
	}
	for i := range ids {
		ptu.projects[ids[i]] = struct{}{}
	}
	return ptu
}

// AddProjects adds the projects edges to Project.
func (ptu *ProjectTypeUpdate) AddProjects(p ...*Project) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptu *ProjectTypeUpdate) AddPropertyIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.properties == nil {
		ptu.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptu.properties[ids[i]] = struct{}{}
	}
	return ptu
}

// AddProperties adds the properties edges to PropertyType.
func (ptu *ProjectTypeUpdate) AddProperties(p ...*PropertyType) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptu *ProjectTypeUpdate) AddWorkOrderIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.work_orders == nil {
		ptu.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		ptu.work_orders[ids[i]] = struct{}{}
	}
	return ptu
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptu *ProjectTypeUpdate) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptu.AddWorkOrderIDs(ids...)
}

// RemoveProjectIDs removes the projects edge to Project by ids.
func (ptu *ProjectTypeUpdate) RemoveProjectIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.removedProjects == nil {
		ptu.removedProjects = make(map[string]struct{})
	}
	for i := range ids {
		ptu.removedProjects[ids[i]] = struct{}{}
	}
	return ptu
}

// RemoveProjects removes projects edges to Project.
func (ptu *ProjectTypeUpdate) RemoveProjects(p ...*Project) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemoveProjectIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to PropertyType by ids.
func (ptu *ProjectTypeUpdate) RemovePropertyIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.removedProperties == nil {
		ptu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		ptu.removedProperties[ids[i]] = struct{}{}
	}
	return ptu
}

// RemoveProperties removes properties edges to PropertyType.
func (ptu *ProjectTypeUpdate) RemoveProperties(p ...*PropertyType) *ProjectTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptu.RemovePropertyIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrderDefinition by ids.
func (ptu *ProjectTypeUpdate) RemoveWorkOrderIDs(ids ...string) *ProjectTypeUpdate {
	if ptu.removedWorkOrders == nil {
		ptu.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		ptu.removedWorkOrders[ids[i]] = struct{}{}
	}
	return ptu
}

// RemoveWorkOrders removes work_orders edges to WorkOrderDefinition.
func (ptu *ProjectTypeUpdate) RemoveWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptu.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ptu *ProjectTypeUpdate) Save(ctx context.Context) (int, error) {
	if ptu.update_time == nil {
		v := projecttype.UpdateDefaultUpdateTime()
		ptu.update_time = &v
	}
	if ptu.name != nil {
		if err := projecttype.NameValidator(*ptu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	return ptu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ptu *ProjectTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := ptu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ptu *ProjectTypeUpdate) Exec(ctx context.Context) error {
	_, err := ptu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptu *ProjectTypeUpdate) ExecX(ctx context.Context) {
	if err := ptu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ptu *ProjectTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(ptu.driver.Dialect())
		selector = builder.Select(projecttype.FieldID).From(builder.Table(projecttype.Table))
	)
	for _, p := range ptu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ptu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := ptu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(projecttype.Table).Where(sql.InInts(projecttype.FieldID, ids...))
	)
	if value := ptu.update_time; value != nil {
		updater.Set(projecttype.FieldUpdateTime, *value)
	}
	if value := ptu.name; value != nil {
		updater.Set(projecttype.FieldName, *value)
	}
	if value := ptu.description; value != nil {
		updater.Set(projecttype.FieldDescription, *value)
	}
	if ptu.cleardescription {
		updater.SetNull(projecttype.FieldDescription)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.removedProjects) > 0 {
		eids := make([]int, len(ptu.removedProjects))
		for eid := range ptu.removedProjects {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(projecttype.ProjectsTable).
			SetNull(projecttype.ProjectsColumn).
			Where(sql.InInts(projecttype.ProjectsColumn, ids...)).
			Where(sql.InInts(project.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.projects) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptu.projects {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(project.FieldID, eid)
			}
			query, args := builder.Update(projecttype.ProjectsTable).
				Set(projecttype.ProjectsColumn, id).
				Where(sql.And(p, sql.IsNull(projecttype.ProjectsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ptu.projects) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"projects\" %v already connected to a different \"ProjectType\"", keys(ptu.projects))})
			}
		}
	}
	if len(ptu.removedProperties) > 0 {
		eids := make([]int, len(ptu.removedProperties))
		for eid := range ptu.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(projecttype.PropertiesTable).
			SetNull(projecttype.PropertiesColumn).
			Where(sql.InInts(projecttype.PropertiesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptu.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(projecttype.PropertiesTable).
				Set(projecttype.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(projecttype.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ptu.properties) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"ProjectType\"", keys(ptu.properties))})
			}
		}
	}
	if len(ptu.removedWorkOrders) > 0 {
		eids := make([]int, len(ptu.removedWorkOrders))
		for eid := range ptu.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(projecttype.WorkOrdersTable).
			SetNull(projecttype.WorkOrdersColumn).
			Where(sql.InInts(projecttype.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorderdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ptu.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptu.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorderdefinition.FieldID, eid)
			}
			query, args := builder.Update(projecttype.WorkOrdersTable).
				Set(projecttype.WorkOrdersColumn, id).
				Where(sql.And(p, sql.IsNull(projecttype.WorkOrdersColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ptu.work_orders) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"ProjectType\"", keys(ptu.work_orders))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// ProjectTypeUpdateOne is the builder for updating a single ProjectType entity.
type ProjectTypeUpdateOne struct {
	config
	id string

	update_time       *time.Time
	name              *string
	description       *string
	cleardescription  bool
	projects          map[string]struct{}
	properties        map[string]struct{}
	work_orders       map[string]struct{}
	removedProjects   map[string]struct{}
	removedProperties map[string]struct{}
	removedWorkOrders map[string]struct{}
}

// SetName sets the name field.
func (ptuo *ProjectTypeUpdateOne) SetName(s string) *ProjectTypeUpdateOne {
	ptuo.name = &s
	return ptuo
}

// SetDescription sets the description field.
func (ptuo *ProjectTypeUpdateOne) SetDescription(s string) *ProjectTypeUpdateOne {
	ptuo.description = &s
	return ptuo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ptuo *ProjectTypeUpdateOne) SetNillableDescription(s *string) *ProjectTypeUpdateOne {
	if s != nil {
		ptuo.SetDescription(*s)
	}
	return ptuo
}

// ClearDescription clears the value of description.
func (ptuo *ProjectTypeUpdateOne) ClearDescription() *ProjectTypeUpdateOne {
	ptuo.description = nil
	ptuo.cleardescription = true
	return ptuo
}

// AddProjectIDs adds the projects edge to Project by ids.
func (ptuo *ProjectTypeUpdateOne) AddProjectIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.projects == nil {
		ptuo.projects = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.projects[ids[i]] = struct{}{}
	}
	return ptuo
}

// AddProjects adds the projects edges to Project.
func (ptuo *ProjectTypeUpdateOne) AddProjects(p ...*Project) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddProjectIDs(ids...)
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (ptuo *ProjectTypeUpdateOne) AddPropertyIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.properties == nil {
		ptuo.properties = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.properties[ids[i]] = struct{}{}
	}
	return ptuo
}

// AddProperties adds the properties edges to PropertyType.
func (ptuo *ProjectTypeUpdateOne) AddProperties(p ...*PropertyType) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.AddPropertyIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (ptuo *ProjectTypeUpdateOne) AddWorkOrderIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.work_orders == nil {
		ptuo.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.work_orders[ids[i]] = struct{}{}
	}
	return ptuo
}

// AddWorkOrders adds the work_orders edges to WorkOrderDefinition.
func (ptuo *ProjectTypeUpdateOne) AddWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptuo.AddWorkOrderIDs(ids...)
}

// RemoveProjectIDs removes the projects edge to Project by ids.
func (ptuo *ProjectTypeUpdateOne) RemoveProjectIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.removedProjects == nil {
		ptuo.removedProjects = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.removedProjects[ids[i]] = struct{}{}
	}
	return ptuo
}

// RemoveProjects removes projects edges to Project.
func (ptuo *ProjectTypeUpdateOne) RemoveProjects(p ...*Project) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemoveProjectIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to PropertyType by ids.
func (ptuo *ProjectTypeUpdateOne) RemovePropertyIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.removedProperties == nil {
		ptuo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.removedProperties[ids[i]] = struct{}{}
	}
	return ptuo
}

// RemoveProperties removes properties edges to PropertyType.
func (ptuo *ProjectTypeUpdateOne) RemoveProperties(p ...*PropertyType) *ProjectTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptuo.RemovePropertyIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrderDefinition by ids.
func (ptuo *ProjectTypeUpdateOne) RemoveWorkOrderIDs(ids ...string) *ProjectTypeUpdateOne {
	if ptuo.removedWorkOrders == nil {
		ptuo.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		ptuo.removedWorkOrders[ids[i]] = struct{}{}
	}
	return ptuo
}

// RemoveWorkOrders removes work_orders edges to WorkOrderDefinition.
func (ptuo *ProjectTypeUpdateOne) RemoveWorkOrders(w ...*WorkOrderDefinition) *ProjectTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return ptuo.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ptuo *ProjectTypeUpdateOne) Save(ctx context.Context) (*ProjectType, error) {
	if ptuo.update_time == nil {
		v := projecttype.UpdateDefaultUpdateTime()
		ptuo.update_time = &v
	}
	if ptuo.name != nil {
		if err := projecttype.NameValidator(*ptuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	return ptuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ptuo *ProjectTypeUpdateOne) SaveX(ctx context.Context) *ProjectType {
	pt, err := ptuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return pt
}

// Exec executes the query on the entity.
func (ptuo *ProjectTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := ptuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptuo *ProjectTypeUpdateOne) ExecX(ctx context.Context) {
	if err := ptuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ptuo *ProjectTypeUpdateOne) sqlSave(ctx context.Context) (pt *ProjectType, err error) {
	var (
		builder  = sql.Dialect(ptuo.driver.Dialect())
		selector = builder.Select(projecttype.Columns...).From(builder.Table(projecttype.Table))
	)
	projecttype.ID(ptuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ptuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		pt = &ProjectType{config: ptuo.config}
		if err := pt.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into ProjectType: %v", err)
		}
		id = pt.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("ProjectType with id: %v", ptuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one ProjectType with the same id: %v", ptuo.id)
	}

	tx, err := ptuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(projecttype.Table).Where(sql.InInts(projecttype.FieldID, ids...))
	)
	if value := ptuo.update_time; value != nil {
		updater.Set(projecttype.FieldUpdateTime, *value)
		pt.UpdateTime = *value
	}
	if value := ptuo.name; value != nil {
		updater.Set(projecttype.FieldName, *value)
		pt.Name = *value
	}
	if value := ptuo.description; value != nil {
		updater.Set(projecttype.FieldDescription, *value)
		pt.Description = value
	}
	if ptuo.cleardescription {
		pt.Description = nil
		updater.SetNull(projecttype.FieldDescription)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.removedProjects) > 0 {
		eids := make([]int, len(ptuo.removedProjects))
		for eid := range ptuo.removedProjects {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(projecttype.ProjectsTable).
			SetNull(projecttype.ProjectsColumn).
			Where(sql.InInts(projecttype.ProjectsColumn, ids...)).
			Where(sql.InInts(project.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.projects) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptuo.projects {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(ptuo.projects) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"projects\" %v already connected to a different \"ProjectType\"", keys(ptuo.projects))})
			}
		}
	}
	if len(ptuo.removedProperties) > 0 {
		eids := make([]int, len(ptuo.removedProperties))
		for eid := range ptuo.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(projecttype.PropertiesTable).
			SetNull(projecttype.PropertiesColumn).
			Where(sql.InInts(projecttype.PropertiesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptuo.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(ptuo.properties) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"ProjectType\"", keys(ptuo.properties))})
			}
		}
	}
	if len(ptuo.removedWorkOrders) > 0 {
		eids := make([]int, len(ptuo.removedWorkOrders))
		for eid := range ptuo.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(projecttype.WorkOrdersTable).
			SetNull(projecttype.WorkOrdersColumn).
			Where(sql.InInts(projecttype.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorderdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ptuo.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ptuo.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(ptuo.work_orders) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"ProjectType\"", keys(ptuo.work_orders))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return pt, nil
}
