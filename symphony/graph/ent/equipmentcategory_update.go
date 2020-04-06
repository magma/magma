// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentCategoryUpdate is the builder for updating EquipmentCategory entities.
type EquipmentCategoryUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentCategoryMutation
	predicates []predicate.EquipmentCategory
}

// Where adds a new predicate for the builder.
func (ecu *EquipmentCategoryUpdate) Where(ps ...predicate.EquipmentCategory) *EquipmentCategoryUpdate {
	ecu.predicates = append(ecu.predicates, ps...)
	return ecu
}

// SetName sets the name field.
func (ecu *EquipmentCategoryUpdate) SetName(s string) *EquipmentCategoryUpdate {
	ecu.mutation.SetName(s)
	return ecu
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecu *EquipmentCategoryUpdate) AddTypeIDs(ids ...int) *EquipmentCategoryUpdate {
	ecu.mutation.AddTypeIDs(ids...)
	return ecu
}

// AddTypes adds the types edges to EquipmentType.
func (ecu *EquipmentCategoryUpdate) AddTypes(e ...*EquipmentType) *EquipmentCategoryUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.AddTypeIDs(ids...)
}

// RemoveTypeIDs removes the types edge to EquipmentType by ids.
func (ecu *EquipmentCategoryUpdate) RemoveTypeIDs(ids ...int) *EquipmentCategoryUpdate {
	ecu.mutation.RemoveTypeIDs(ids...)
	return ecu
}

// RemoveTypes removes types edges to EquipmentType.
func (ecu *EquipmentCategoryUpdate) RemoveTypes(e ...*EquipmentType) *EquipmentCategoryUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.RemoveTypeIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ecu *EquipmentCategoryUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := ecu.mutation.UpdateTime(); !ok {
		v := equipmentcategory.UpdateDefaultUpdateTime()
		ecu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(ecu.hooks) == 0 {
		affected, err = ecu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ecu.mutation = mutation
			affected, err = ecu.sqlSave(ctx)
			return affected, err
		})
		for i := len(ecu.hooks); i > 0; i-- {
			mut = ecu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, ecu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ecu *EquipmentCategoryUpdate) SaveX(ctx context.Context) int {
	affected, err := ecu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ecu *EquipmentCategoryUpdate) Exec(ctx context.Context) error {
	_, err := ecu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecu *EquipmentCategoryUpdate) ExecX(ctx context.Context) {
	if err := ecu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ecu *EquipmentCategoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentcategory.Table,
			Columns: equipmentcategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentcategory.FieldID,
			},
		},
	}
	if ps := ecu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ecu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentcategory.FieldUpdateTime,
		})
	}
	if value, ok := ecu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentcategory.FieldName,
		})
	}
	if nodes := ecu.mutation.RemovedTypesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecu.mutation.TypesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ecu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentcategory.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentCategoryUpdateOne is the builder for updating a single EquipmentCategory entity.
type EquipmentCategoryUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentCategoryMutation
}

// SetName sets the name field.
func (ecuo *EquipmentCategoryUpdateOne) SetName(s string) *EquipmentCategoryUpdateOne {
	ecuo.mutation.SetName(s)
	return ecuo
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecuo *EquipmentCategoryUpdateOne) AddTypeIDs(ids ...int) *EquipmentCategoryUpdateOne {
	ecuo.mutation.AddTypeIDs(ids...)
	return ecuo
}

// AddTypes adds the types edges to EquipmentType.
func (ecuo *EquipmentCategoryUpdateOne) AddTypes(e ...*EquipmentType) *EquipmentCategoryUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.AddTypeIDs(ids...)
}

// RemoveTypeIDs removes the types edge to EquipmentType by ids.
func (ecuo *EquipmentCategoryUpdateOne) RemoveTypeIDs(ids ...int) *EquipmentCategoryUpdateOne {
	ecuo.mutation.RemoveTypeIDs(ids...)
	return ecuo
}

// RemoveTypes removes types edges to EquipmentType.
func (ecuo *EquipmentCategoryUpdateOne) RemoveTypes(e ...*EquipmentType) *EquipmentCategoryUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.RemoveTypeIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ecuo *EquipmentCategoryUpdateOne) Save(ctx context.Context) (*EquipmentCategory, error) {
	if _, ok := ecuo.mutation.UpdateTime(); !ok {
		v := equipmentcategory.UpdateDefaultUpdateTime()
		ecuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *EquipmentCategory
	)
	if len(ecuo.hooks) == 0 {
		node, err = ecuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ecuo.mutation = mutation
			node, err = ecuo.sqlSave(ctx)
			return node, err
		})
		for i := len(ecuo.hooks); i > 0; i-- {
			mut = ecuo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, ecuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ecuo *EquipmentCategoryUpdateOne) SaveX(ctx context.Context) *EquipmentCategory {
	ec, err := ecuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ec
}

// Exec executes the query on the entity.
func (ecuo *EquipmentCategoryUpdateOne) Exec(ctx context.Context) error {
	_, err := ecuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecuo *EquipmentCategoryUpdateOne) ExecX(ctx context.Context) {
	if err := ecuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ecuo *EquipmentCategoryUpdateOne) sqlSave(ctx context.Context) (ec *EquipmentCategory, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentcategory.Table,
			Columns: equipmentcategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentcategory.FieldID,
			},
		},
	}
	id, ok := ecuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing EquipmentCategory.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := ecuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentcategory.FieldUpdateTime,
		})
	}
	if value, ok := ecuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentcategory.FieldName,
		})
	}
	if nodes := ecuo.mutation.RemovedTypesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecuo.mutation.TypesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	ec = &EquipmentCategory{config: ecuo.config}
	_spec.Assign = ec.assignValues
	_spec.ScanValues = ec.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ecuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentcategory.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ec, nil
}
