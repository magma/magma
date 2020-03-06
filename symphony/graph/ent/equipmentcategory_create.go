// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
)

// EquipmentCategoryCreate is the builder for creating a EquipmentCategory entity.
type EquipmentCategoryCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	types       map[int]struct{}
}

// SetCreateTime sets the create_time field.
func (ecc *EquipmentCategoryCreate) SetCreateTime(t time.Time) *EquipmentCategoryCreate {
	ecc.create_time = &t
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
	ecc.update_time = &t
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
	ecc.name = &s
	return ecc
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecc *EquipmentCategoryCreate) AddTypeIDs(ids ...int) *EquipmentCategoryCreate {
	if ecc.types == nil {
		ecc.types = make(map[int]struct{})
	}
	for i := range ids {
		ecc.types[ids[i]] = struct{}{}
	}
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
	if ecc.create_time == nil {
		v := equipmentcategory.DefaultCreateTime()
		ecc.create_time = &v
	}
	if ecc.update_time == nil {
		v := equipmentcategory.DefaultUpdateTime()
		ecc.update_time = &v
	}
	if ecc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	return ecc.sqlSave(ctx)
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
	if value := ecc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentcategory.FieldCreateTime,
		})
		ec.CreateTime = *value
	}
	if value := ecc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentcategory.FieldUpdateTime,
		})
		ec.UpdateTime = *value
	}
	if value := ecc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentcategory.FieldName,
		})
		ec.Name = *value
	}
	if nodes := ecc.types; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
