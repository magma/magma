// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"strconv"
	"time"

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

	update_time  *time.Time
	name         *string
	types        map[string]struct{}
	removedTypes map[string]struct{}
	predicates   []predicate.EquipmentCategory
}

// Where adds a new predicate for the builder.
func (ecu *EquipmentCategoryUpdate) Where(ps ...predicate.EquipmentCategory) *EquipmentCategoryUpdate {
	ecu.predicates = append(ecu.predicates, ps...)
	return ecu
}

// SetName sets the name field.
func (ecu *EquipmentCategoryUpdate) SetName(s string) *EquipmentCategoryUpdate {
	ecu.name = &s
	return ecu
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecu *EquipmentCategoryUpdate) AddTypeIDs(ids ...string) *EquipmentCategoryUpdate {
	if ecu.types == nil {
		ecu.types = make(map[string]struct{})
	}
	for i := range ids {
		ecu.types[ids[i]] = struct{}{}
	}
	return ecu
}

// AddTypes adds the types edges to EquipmentType.
func (ecu *EquipmentCategoryUpdate) AddTypes(e ...*EquipmentType) *EquipmentCategoryUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.AddTypeIDs(ids...)
}

// RemoveTypeIDs removes the types edge to EquipmentType by ids.
func (ecu *EquipmentCategoryUpdate) RemoveTypeIDs(ids ...string) *EquipmentCategoryUpdate {
	if ecu.removedTypes == nil {
		ecu.removedTypes = make(map[string]struct{})
	}
	for i := range ids {
		ecu.removedTypes[ids[i]] = struct{}{}
	}
	return ecu
}

// RemoveTypes removes types edges to EquipmentType.
func (ecu *EquipmentCategoryUpdate) RemoveTypes(e ...*EquipmentType) *EquipmentCategoryUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.RemoveTypeIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ecu *EquipmentCategoryUpdate) Save(ctx context.Context) (int, error) {
	if ecu.update_time == nil {
		v := equipmentcategory.UpdateDefaultUpdateTime()
		ecu.update_time = &v
	}
	return ecu.sqlSave(ctx)
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
				Type:   field.TypeString,
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
	if value := ecu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentcategory.FieldUpdateTime,
		})
	}
	if value := ecu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentcategory.FieldName,
		})
	}
	if nodes := ecu.removedTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentcategory.TypesTable,
			Columns: []string{equipmentcategory.TypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecu.types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentcategory.TypesTable,
			Columns: []string{equipmentcategory.TypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ecu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentCategoryUpdateOne is the builder for updating a single EquipmentCategory entity.
type EquipmentCategoryUpdateOne struct {
	config
	id string

	update_time  *time.Time
	name         *string
	types        map[string]struct{}
	removedTypes map[string]struct{}
}

// SetName sets the name field.
func (ecuo *EquipmentCategoryUpdateOne) SetName(s string) *EquipmentCategoryUpdateOne {
	ecuo.name = &s
	return ecuo
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecuo *EquipmentCategoryUpdateOne) AddTypeIDs(ids ...string) *EquipmentCategoryUpdateOne {
	if ecuo.types == nil {
		ecuo.types = make(map[string]struct{})
	}
	for i := range ids {
		ecuo.types[ids[i]] = struct{}{}
	}
	return ecuo
}

// AddTypes adds the types edges to EquipmentType.
func (ecuo *EquipmentCategoryUpdateOne) AddTypes(e ...*EquipmentType) *EquipmentCategoryUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.AddTypeIDs(ids...)
}

// RemoveTypeIDs removes the types edge to EquipmentType by ids.
func (ecuo *EquipmentCategoryUpdateOne) RemoveTypeIDs(ids ...string) *EquipmentCategoryUpdateOne {
	if ecuo.removedTypes == nil {
		ecuo.removedTypes = make(map[string]struct{})
	}
	for i := range ids {
		ecuo.removedTypes[ids[i]] = struct{}{}
	}
	return ecuo
}

// RemoveTypes removes types edges to EquipmentType.
func (ecuo *EquipmentCategoryUpdateOne) RemoveTypes(e ...*EquipmentType) *EquipmentCategoryUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.RemoveTypeIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ecuo *EquipmentCategoryUpdateOne) Save(ctx context.Context) (*EquipmentCategory, error) {
	if ecuo.update_time == nil {
		v := equipmentcategory.UpdateDefaultUpdateTime()
		ecuo.update_time = &v
	}
	return ecuo.sqlSave(ctx)
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
				Value:  ecuo.id,
				Type:   field.TypeString,
				Column: equipmentcategory.FieldID,
			},
		},
	}
	if value := ecuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentcategory.FieldUpdateTime,
		})
	}
	if value := ecuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentcategory.FieldName,
		})
	}
	if nodes := ecuo.removedTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentcategory.TypesTable,
			Columns: []string{equipmentcategory.TypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmenttype.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecuo.types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentcategory.TypesTable,
			Columns: []string{equipmentcategory.TypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmenttype.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	ec = &EquipmentCategory{config: ecuo.config}
	_spec.Assign = ec.assignValues
	_spec.ScanValues = ec.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ecuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ec, nil
}
