// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
)

// EquipmentCategoryCreate is the builder for creating a EquipmentCategory entity.
type EquipmentCategoryCreate struct {
	config
	mutation *EquipmentCategoryMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (ecc *EquipmentCategoryCreate) SetCreateTime(t time.Time) *EquipmentCategoryCreate {
	ecc.mutation.SetCreateTime(t)
	return ecc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ecc *EquipmentCategoryCreate) SetNillableCreateTime(t *time.Time) *EquipmentCategoryCreate {
	if t != nil {
		ecc.SetCreateTime(*t)
	}
	return ecc
}

// SetUpdateTime sets the update_time field.
func (ecc *EquipmentCategoryCreate) SetUpdateTime(t time.Time) *EquipmentCategoryCreate {
	ecc.mutation.SetUpdateTime(t)
	return ecc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ecc *EquipmentCategoryCreate) SetNillableUpdateTime(t *time.Time) *EquipmentCategoryCreate {
	if t != nil {
		ecc.SetUpdateTime(*t)
	}
	return ecc
}

// SetName sets the name field.
func (ecc *EquipmentCategoryCreate) SetName(s string) *EquipmentCategoryCreate {
	ecc.mutation.SetName(s)
	return ecc
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecc *EquipmentCategoryCreate) AddTypeIDs(ids ...int) *EquipmentCategoryCreate {
	ecc.mutation.AddTypeIDs(ids...)
	return ecc
}

// AddTypes adds the types edges to EquipmentType.
func (ecc *EquipmentCategoryCreate) AddTypes(e ...*EquipmentType) *EquipmentCategoryCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecc.AddTypeIDs(ids...)
}

// Save creates the EquipmentCategory in the database.
func (ecc *EquipmentCategoryCreate) Save(ctx context.Context) (*EquipmentCategory, error) {
	if _, ok := ecc.mutation.CreateTime(); !ok {
		v := equipmentcategory.DefaultCreateTime()
		ecc.mutation.SetCreateTime(v)
	}
	if _, ok := ecc.mutation.UpdateTime(); !ok {
		v := equipmentcategory.DefaultUpdateTime()
		ecc.mutation.SetUpdateTime(v)
	}
	if _, ok := ecc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	var (
		err  error
		node *EquipmentCategory
	)
	if len(ecc.hooks) == 0 {
		node, err = ecc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ecc.mutation = mutation
			node, err = ecc.sqlSave(ctx)
			return node, err
		})
		for i := len(ecc.hooks) - 1; i >= 0; i-- {
			mut = ecc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ecc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ecc *EquipmentCategoryCreate) SaveX(ctx context.Context) *EquipmentCategory {
	v, err := ecc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ecc *EquipmentCategoryCreate) sqlSave(ctx context.Context) (*EquipmentCategory, error) {
	var (
		ec    = &EquipmentCategory{config: ecc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: equipmentcategory.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentcategory.FieldID,
			},
		}
	)
	if value, ok := ecc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentcategory.FieldCreateTime,
		})
		ec.CreateTime = value
	}
	if value, ok := ecc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentcategory.FieldUpdateTime,
		})
		ec.UpdateTime = value
	}
	if value, ok := ecc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentcategory.FieldName,
		})
		ec.Name = value
	}
	if nodes := ecc.mutation.TypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentcategory.TypesTable,
			Columns: []string{equipmentcategory.TypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ecc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	ec.ID = int(id)
	return ec, nil
}
