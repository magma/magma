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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderDefinitionUpdate is the builder for updating WorkOrderDefinition entities.
type WorkOrderDefinitionUpdate struct {
	config

	update_time        *time.Time
	index              *int
	addindex           *int
	clearindex         bool
	_type              map[string]struct{}
	project_type       map[string]struct{}
	clearedType        bool
	clearedProjectType bool
	predicates         []predicate.WorkOrderDefinition
}

// Where adds a new predicate for the builder.
func (wodu *WorkOrderDefinitionUpdate) Where(ps ...predicate.WorkOrderDefinition) *WorkOrderDefinitionUpdate {
	wodu.predicates = append(wodu.predicates, ps...)
	return wodu
}

// SetIndex sets the index field.
func (wodu *WorkOrderDefinitionUpdate) SetIndex(i int) *WorkOrderDefinitionUpdate {
	wodu.index = &i
	wodu.addindex = nil
	return wodu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (wodu *WorkOrderDefinitionUpdate) SetNillableIndex(i *int) *WorkOrderDefinitionUpdate {
	if i != nil {
		wodu.SetIndex(*i)
	}
	return wodu
}

// AddIndex adds i to index.
func (wodu *WorkOrderDefinitionUpdate) AddIndex(i int) *WorkOrderDefinitionUpdate {
	if wodu.addindex == nil {
		wodu.addindex = &i
	} else {
		*wodu.addindex += i
	}
	return wodu
}

// ClearIndex clears the value of index.
func (wodu *WorkOrderDefinitionUpdate) ClearIndex() *WorkOrderDefinitionUpdate {
	wodu.index = nil
	wodu.clearindex = true
	return wodu
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wodu *WorkOrderDefinitionUpdate) SetTypeID(id string) *WorkOrderDefinitionUpdate {
	if wodu._type == nil {
		wodu._type = make(map[string]struct{})
	}
	wodu._type[id] = struct{}{}
	return wodu
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wodu *WorkOrderDefinitionUpdate) SetNillableTypeID(id *string) *WorkOrderDefinitionUpdate {
	if id != nil {
		wodu = wodu.SetTypeID(*id)
	}
	return wodu
}

// SetType sets the type edge to WorkOrderType.
func (wodu *WorkOrderDefinitionUpdate) SetType(w *WorkOrderType) *WorkOrderDefinitionUpdate {
	return wodu.SetTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (wodu *WorkOrderDefinitionUpdate) SetProjectTypeID(id string) *WorkOrderDefinitionUpdate {
	if wodu.project_type == nil {
		wodu.project_type = make(map[string]struct{})
	}
	wodu.project_type[id] = struct{}{}
	return wodu
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (wodu *WorkOrderDefinitionUpdate) SetNillableProjectTypeID(id *string) *WorkOrderDefinitionUpdate {
	if id != nil {
		wodu = wodu.SetProjectTypeID(*id)
	}
	return wodu
}

// SetProjectType sets the project_type edge to ProjectType.
func (wodu *WorkOrderDefinitionUpdate) SetProjectType(p *ProjectType) *WorkOrderDefinitionUpdate {
	return wodu.SetProjectTypeID(p.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (wodu *WorkOrderDefinitionUpdate) ClearType() *WorkOrderDefinitionUpdate {
	wodu.clearedType = true
	return wodu
}

// ClearProjectType clears the project_type edge to ProjectType.
func (wodu *WorkOrderDefinitionUpdate) ClearProjectType() *WorkOrderDefinitionUpdate {
	wodu.clearedProjectType = true
	return wodu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (wodu *WorkOrderDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if wodu.update_time == nil {
		v := workorderdefinition.UpdateDefaultUpdateTime()
		wodu.update_time = &v
	}
	if len(wodu._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if len(wodu.project_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"project_type\"")
	}
	return wodu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (wodu *WorkOrderDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := wodu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wodu *WorkOrderDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := wodu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wodu *WorkOrderDefinitionUpdate) ExecX(ctx context.Context) {
	if err := wodu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wodu *WorkOrderDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(wodu.driver.Dialect())
		selector = builder.Select(workorderdefinition.FieldID).From(builder.Table(workorderdefinition.Table))
	)
	for _, p := range wodu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = wodu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := wodu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(workorderdefinition.Table).Where(sql.InInts(workorderdefinition.FieldID, ids...))
	)
	if value := wodu.update_time; value != nil {
		updater.Set(workorderdefinition.FieldUpdateTime, *value)
	}
	if value := wodu.index; value != nil {
		updater.Set(workorderdefinition.FieldIndex, *value)
	}
	if value := wodu.addindex; value != nil {
		updater.Add(workorderdefinition.FieldIndex, *value)
	}
	if wodu.clearindex {
		updater.SetNull(workorderdefinition.FieldIndex)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if wodu.clearedType {
		query, args := builder.Update(workorderdefinition.TypeTable).
			SetNull(workorderdefinition.TypeColumn).
			Where(sql.InInts(workordertype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(wodu._type) > 0 {
		for eid := range wodu._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(workorderdefinition.TypeTable).
				Set(workorderdefinition.TypeColumn, eid).
				Where(sql.InInts(workorderdefinition.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if wodu.clearedProjectType {
		query, args := builder.Update(workorderdefinition.ProjectTypeTable).
			SetNull(workorderdefinition.ProjectTypeColumn).
			Where(sql.InInts(projecttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(wodu.project_type) > 0 {
		for eid := range wodu.project_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(workorderdefinition.ProjectTypeTable).
				Set(workorderdefinition.ProjectTypeColumn, eid).
				Where(sql.InInts(workorderdefinition.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// WorkOrderDefinitionUpdateOne is the builder for updating a single WorkOrderDefinition entity.
type WorkOrderDefinitionUpdateOne struct {
	config
	id string

	update_time        *time.Time
	index              *int
	addindex           *int
	clearindex         bool
	_type              map[string]struct{}
	project_type       map[string]struct{}
	clearedType        bool
	clearedProjectType bool
}

// SetIndex sets the index field.
func (woduo *WorkOrderDefinitionUpdateOne) SetIndex(i int) *WorkOrderDefinitionUpdateOne {
	woduo.index = &i
	woduo.addindex = nil
	return woduo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (woduo *WorkOrderDefinitionUpdateOne) SetNillableIndex(i *int) *WorkOrderDefinitionUpdateOne {
	if i != nil {
		woduo.SetIndex(*i)
	}
	return woduo
}

// AddIndex adds i to index.
func (woduo *WorkOrderDefinitionUpdateOne) AddIndex(i int) *WorkOrderDefinitionUpdateOne {
	if woduo.addindex == nil {
		woduo.addindex = &i
	} else {
		*woduo.addindex += i
	}
	return woduo
}

// ClearIndex clears the value of index.
func (woduo *WorkOrderDefinitionUpdateOne) ClearIndex() *WorkOrderDefinitionUpdateOne {
	woduo.index = nil
	woduo.clearindex = true
	return woduo
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (woduo *WorkOrderDefinitionUpdateOne) SetTypeID(id string) *WorkOrderDefinitionUpdateOne {
	if woduo._type == nil {
		woduo._type = make(map[string]struct{})
	}
	woduo._type[id] = struct{}{}
	return woduo
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (woduo *WorkOrderDefinitionUpdateOne) SetNillableTypeID(id *string) *WorkOrderDefinitionUpdateOne {
	if id != nil {
		woduo = woduo.SetTypeID(*id)
	}
	return woduo
}

// SetType sets the type edge to WorkOrderType.
func (woduo *WorkOrderDefinitionUpdateOne) SetType(w *WorkOrderType) *WorkOrderDefinitionUpdateOne {
	return woduo.SetTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (woduo *WorkOrderDefinitionUpdateOne) SetProjectTypeID(id string) *WorkOrderDefinitionUpdateOne {
	if woduo.project_type == nil {
		woduo.project_type = make(map[string]struct{})
	}
	woduo.project_type[id] = struct{}{}
	return woduo
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (woduo *WorkOrderDefinitionUpdateOne) SetNillableProjectTypeID(id *string) *WorkOrderDefinitionUpdateOne {
	if id != nil {
		woduo = woduo.SetProjectTypeID(*id)
	}
	return woduo
}

// SetProjectType sets the project_type edge to ProjectType.
func (woduo *WorkOrderDefinitionUpdateOne) SetProjectType(p *ProjectType) *WorkOrderDefinitionUpdateOne {
	return woduo.SetProjectTypeID(p.ID)
}

// ClearType clears the type edge to WorkOrderType.
func (woduo *WorkOrderDefinitionUpdateOne) ClearType() *WorkOrderDefinitionUpdateOne {
	woduo.clearedType = true
	return woduo
}

// ClearProjectType clears the project_type edge to ProjectType.
func (woduo *WorkOrderDefinitionUpdateOne) ClearProjectType() *WorkOrderDefinitionUpdateOne {
	woduo.clearedProjectType = true
	return woduo
}

// Save executes the query and returns the updated entity.
func (woduo *WorkOrderDefinitionUpdateOne) Save(ctx context.Context) (*WorkOrderDefinition, error) {
	if woduo.update_time == nil {
		v := workorderdefinition.UpdateDefaultUpdateTime()
		woduo.update_time = &v
	}
	if len(woduo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if len(woduo.project_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project_type\"")
	}
	return woduo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (woduo *WorkOrderDefinitionUpdateOne) SaveX(ctx context.Context) *WorkOrderDefinition {
	wod, err := woduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return wod
}

// Exec executes the query on the entity.
func (woduo *WorkOrderDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := woduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (woduo *WorkOrderDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := woduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (woduo *WorkOrderDefinitionUpdateOne) sqlSave(ctx context.Context) (wod *WorkOrderDefinition, err error) {
	var (
		builder  = sql.Dialect(woduo.driver.Dialect())
		selector = builder.Select(workorderdefinition.Columns...).From(builder.Table(workorderdefinition.Table))
	)
	workorderdefinition.ID(woduo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = woduo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		wod = &WorkOrderDefinition{config: woduo.config}
		if err := wod.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into WorkOrderDefinition: %v", err)
		}
		id = wod.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("WorkOrderDefinition with id: %v", woduo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one WorkOrderDefinition with the same id: %v", woduo.id)
	}

	tx, err := woduo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(workorderdefinition.Table).Where(sql.InInts(workorderdefinition.FieldID, ids...))
	)
	if value := woduo.update_time; value != nil {
		updater.Set(workorderdefinition.FieldUpdateTime, *value)
		wod.UpdateTime = *value
	}
	if value := woduo.index; value != nil {
		updater.Set(workorderdefinition.FieldIndex, *value)
		wod.Index = *value
	}
	if value := woduo.addindex; value != nil {
		updater.Add(workorderdefinition.FieldIndex, *value)
		wod.Index += *value
	}
	if woduo.clearindex {
		var value int
		wod.Index = value
		updater.SetNull(workorderdefinition.FieldIndex)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if woduo.clearedType {
		query, args := builder.Update(workorderdefinition.TypeTable).
			SetNull(workorderdefinition.TypeColumn).
			Where(sql.InInts(workordertype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(woduo._type) > 0 {
		for eid := range woduo._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(workorderdefinition.TypeTable).
				Set(workorderdefinition.TypeColumn, eid).
				Where(sql.InInts(workorderdefinition.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if woduo.clearedProjectType {
		query, args := builder.Update(workorderdefinition.ProjectTypeTable).
			SetNull(workorderdefinition.ProjectTypeColumn).
			Where(sql.InInts(projecttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(woduo.project_type) > 0 {
		for eid := range woduo.project_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(workorderdefinition.ProjectTypeTable).
				Set(workorderdefinition.ProjectTypeColumn, eid).
				Where(sql.InInts(workorderdefinition.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return wod, nil
}
