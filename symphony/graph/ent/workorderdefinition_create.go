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

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
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
		wod  = &WorkOrderDefinition{config: wodc.config}
		spec = &sqlgraph.CreateSpec{
			Table: workorderdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: workorderdefinition.FieldID,
			},
		}
	)
	if value := wodc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorderdefinition.FieldCreateTime,
		})
		wod.CreateTime = *value
	}
	if value := wodc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: workorderdefinition.FieldUpdateTime,
		})
		wod.UpdateTime = *value
	}
	if value := wodc.index; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: workorderdefinition.FieldIndex,
		})
		wod.Index = *value
	}
	if nodes := wodc._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   workorderdefinition.TypeTable,
			Columns: []string{workorderdefinition.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workordertype.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := wodc.project_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workorderdefinition.ProjectTypeTable,
			Columns: []string{workorderdefinition.ProjectTypeColumn},
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
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, wodc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	wod.ID = strconv.FormatInt(id, 10)
	return wod, nil
}
