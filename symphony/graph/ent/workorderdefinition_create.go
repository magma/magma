// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
)

// WorkOrderDefinitionCreate is the builder for creating a WorkOrderDefinition entity.
type WorkOrderDefinitionCreate struct {
	config
	create_time  *time.Time
	update_time  *time.Time
	index        *int
	_type        map[string]struct{}
	project_type map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (wodc *WorkOrderDefinitionCreate) SetCreateTime(t time.Time) *WorkOrderDefinitionCreate {
	wodc.create_time = &t
	return wodc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableCreateTime(t *time.Time) *WorkOrderDefinitionCreate {
	if t != nil {
		wodc.SetCreateTime(*t)
	}
	return wodc
}

// SetUpdateTime sets the update_time field.
func (wodc *WorkOrderDefinitionCreate) SetUpdateTime(t time.Time) *WorkOrderDefinitionCreate {
	wodc.update_time = &t
	return wodc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableUpdateTime(t *time.Time) *WorkOrderDefinitionCreate {
	if t != nil {
		wodc.SetUpdateTime(*t)
	}
	return wodc
}

// SetIndex sets the index field.
func (wodc *WorkOrderDefinitionCreate) SetIndex(i int) *WorkOrderDefinitionCreate {
	wodc.index = &i
	return wodc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableIndex(i *int) *WorkOrderDefinitionCreate {
	if i != nil {
		wodc.SetIndex(*i)
	}
	return wodc
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (wodc *WorkOrderDefinitionCreate) SetTypeID(id string) *WorkOrderDefinitionCreate {
	if wodc._type == nil {
		wodc._type = make(map[string]struct{})
	}
	wodc._type[id] = struct{}{}
	return wodc
}

// SetNillableTypeID sets the type edge to WorkOrderType by id if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableTypeID(id *string) *WorkOrderDefinitionCreate {
	if id != nil {
		wodc = wodc.SetTypeID(*id)
	}
	return wodc
}

// SetType sets the type edge to WorkOrderType.
func (wodc *WorkOrderDefinitionCreate) SetType(w *WorkOrderType) *WorkOrderDefinitionCreate {
	return wodc.SetTypeID(w.ID)
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (wodc *WorkOrderDefinitionCreate) SetProjectTypeID(id string) *WorkOrderDefinitionCreate {
	if wodc.project_type == nil {
		wodc.project_type = make(map[string]struct{})
	}
	wodc.project_type[id] = struct{}{}
	return wodc
}

// SetNillableProjectTypeID sets the project_type edge to ProjectType by id if the given value is not nil.
func (wodc *WorkOrderDefinitionCreate) SetNillableProjectTypeID(id *string) *WorkOrderDefinitionCreate {
	if id != nil {
		wodc = wodc.SetProjectTypeID(*id)
	}
	return wodc
}

// SetProjectType sets the project_type edge to ProjectType.
func (wodc *WorkOrderDefinitionCreate) SetProjectType(p *ProjectType) *WorkOrderDefinitionCreate {
	return wodc.SetProjectTypeID(p.ID)
}

// Save creates the WorkOrderDefinition in the database.
func (wodc *WorkOrderDefinitionCreate) Save(ctx context.Context) (*WorkOrderDefinition, error) {
	if wodc.create_time == nil {
		v := workorderdefinition.DefaultCreateTime()
		wodc.create_time = &v
	}
	if wodc.update_time == nil {
		v := workorderdefinition.DefaultUpdateTime()
		wodc.update_time = &v
	}
	if len(wodc._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if len(wodc.project_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"project_type\"")
	}
	return wodc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (wodc *WorkOrderDefinitionCreate) SaveX(ctx context.Context) *WorkOrderDefinition {
	v, err := wodc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wodc *WorkOrderDefinitionCreate) sqlSave(ctx context.Context) (*WorkOrderDefinition, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(wodc.driver.Dialect())
		wod     = &WorkOrderDefinition{config: wodc.config}
	)
	tx, err := wodc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(workorderdefinition.Table).Default()
	if value := wodc.create_time; value != nil {
		insert.Set(workorderdefinition.FieldCreateTime, *value)
		wod.CreateTime = *value
	}
	if value := wodc.update_time; value != nil {
		insert.Set(workorderdefinition.FieldUpdateTime, *value)
		wod.UpdateTime = *value
	}
	if value := wodc.index; value != nil {
		insert.Set(workorderdefinition.FieldIndex, *value)
		wod.Index = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(workorderdefinition.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	wod.ID = strconv.FormatInt(id, 10)
	if len(wodc._type) > 0 {
		for eid := range wodc._type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(workorderdefinition.TypeTable).
				Set(workorderdefinition.TypeColumn, eid).
				Where(sql.EQ(workorderdefinition.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(wodc.project_type) > 0 {
		for eid := range wodc.project_type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(workorderdefinition.ProjectTypeTable).
				Set(workorderdefinition.ProjectTypeColumn, eid).
				Where(sql.EQ(workorderdefinition.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return wod, nil
}
