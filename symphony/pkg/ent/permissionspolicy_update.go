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
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"
)

// PermissionsPolicyUpdate is the builder for updating PermissionsPolicy entities.
type PermissionsPolicyUpdate struct {
	config
	hooks      []Hook
	mutation   *PermissionsPolicyMutation
	predicates []predicate.PermissionsPolicy
}

// Where adds a new predicate for the builder.
func (ppu *PermissionsPolicyUpdate) Where(ps ...predicate.PermissionsPolicy) *PermissionsPolicyUpdate {
	ppu.predicates = append(ppu.predicates, ps...)
	return ppu
}

// SetName sets the name field.
func (ppu *PermissionsPolicyUpdate) SetName(s string) *PermissionsPolicyUpdate {
	ppu.mutation.SetName(s)
	return ppu
}

// SetDescription sets the description field.
func (ppu *PermissionsPolicyUpdate) SetDescription(s string) *PermissionsPolicyUpdate {
	ppu.mutation.SetDescription(s)
	return ppu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ppu *PermissionsPolicyUpdate) SetNillableDescription(s *string) *PermissionsPolicyUpdate {
	if s != nil {
		ppu.SetDescription(*s)
	}
	return ppu
}

// ClearDescription clears the value of description.
func (ppu *PermissionsPolicyUpdate) ClearDescription() *PermissionsPolicyUpdate {
	ppu.mutation.ClearDescription()
	return ppu
}

// SetIsGlobal sets the is_global field.
func (ppu *PermissionsPolicyUpdate) SetIsGlobal(b bool) *PermissionsPolicyUpdate {
	ppu.mutation.SetIsGlobal(b)
	return ppu
}

// SetNillableIsGlobal sets the is_global field if the given value is not nil.
func (ppu *PermissionsPolicyUpdate) SetNillableIsGlobal(b *bool) *PermissionsPolicyUpdate {
	if b != nil {
		ppu.SetIsGlobal(*b)
	}
	return ppu
}

// ClearIsGlobal clears the value of is_global.
func (ppu *PermissionsPolicyUpdate) ClearIsGlobal() *PermissionsPolicyUpdate {
	ppu.mutation.ClearIsGlobal()
	return ppu
}

// SetInventoryPolicy sets the inventory_policy field.
func (ppu *PermissionsPolicyUpdate) SetInventoryPolicy(mpi *models.InventoryPolicyInput) *PermissionsPolicyUpdate {
	ppu.mutation.SetInventoryPolicy(mpi)
	return ppu
}

// ClearInventoryPolicy clears the value of inventory_policy.
func (ppu *PermissionsPolicyUpdate) ClearInventoryPolicy() *PermissionsPolicyUpdate {
	ppu.mutation.ClearInventoryPolicy()
	return ppu
}

// SetWorkforcePolicy sets the workforce_policy field.
func (ppu *PermissionsPolicyUpdate) SetWorkforcePolicy(mpi *models.WorkforcePolicyInput) *PermissionsPolicyUpdate {
	ppu.mutation.SetWorkforcePolicy(mpi)
	return ppu
}

// ClearWorkforcePolicy clears the value of workforce_policy.
func (ppu *PermissionsPolicyUpdate) ClearWorkforcePolicy() *PermissionsPolicyUpdate {
	ppu.mutation.ClearWorkforcePolicy()
	return ppu
}

// AddGroupIDs adds the groups edge to UsersGroup by ids.
func (ppu *PermissionsPolicyUpdate) AddGroupIDs(ids ...int) *PermissionsPolicyUpdate {
	ppu.mutation.AddGroupIDs(ids...)
	return ppu
}

// AddGroups adds the groups edges to UsersGroup.
func (ppu *PermissionsPolicyUpdate) AddGroups(u ...*UsersGroup) *PermissionsPolicyUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ppu.AddGroupIDs(ids...)
}

// RemoveGroupIDs removes the groups edge to UsersGroup by ids.
func (ppu *PermissionsPolicyUpdate) RemoveGroupIDs(ids ...int) *PermissionsPolicyUpdate {
	ppu.mutation.RemoveGroupIDs(ids...)
	return ppu
}

// RemoveGroups removes groups edges to UsersGroup.
func (ppu *PermissionsPolicyUpdate) RemoveGroups(u ...*UsersGroup) *PermissionsPolicyUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ppu.RemoveGroupIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ppu *PermissionsPolicyUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := ppu.mutation.UpdateTime(); !ok {
		v := permissionspolicy.UpdateDefaultUpdateTime()
		ppu.mutation.SetUpdateTime(v)
	}
	if v, ok := ppu.mutation.Name(); ok {
		if err := permissionspolicy.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(ppu.hooks) == 0 {
		affected, err = ppu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PermissionsPolicyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ppu.mutation = mutation
			affected, err = ppu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ppu.hooks) - 1; i >= 0; i-- {
			mut = ppu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ppu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ppu *PermissionsPolicyUpdate) SaveX(ctx context.Context) int {
	affected, err := ppu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ppu *PermissionsPolicyUpdate) Exec(ctx context.Context) error {
	_, err := ppu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppu *PermissionsPolicyUpdate) ExecX(ctx context.Context) {
	if err := ppu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ppu *PermissionsPolicyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   permissionspolicy.Table,
			Columns: permissionspolicy.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: permissionspolicy.FieldID,
			},
		},
	}
	if ps := ppu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ppu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: permissionspolicy.FieldUpdateTime,
		})
	}
	if value, ok := ppu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: permissionspolicy.FieldName,
		})
	}
	if value, ok := ppu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: permissionspolicy.FieldDescription,
		})
	}
	if ppu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: permissionspolicy.FieldDescription,
		})
	}
	if value, ok := ppu.mutation.IsGlobal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: permissionspolicy.FieldIsGlobal,
		})
	}
	if ppu.mutation.IsGlobalCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: permissionspolicy.FieldIsGlobal,
		})
	}
	if value, ok := ppu.mutation.InventoryPolicy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: permissionspolicy.FieldInventoryPolicy,
		})
	}
	if ppu.mutation.InventoryPolicyCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: permissionspolicy.FieldInventoryPolicy,
		})
	}
	if value, ok := ppu.mutation.WorkforcePolicy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: permissionspolicy.FieldWorkforcePolicy,
		})
	}
	if ppu.mutation.WorkforcePolicyCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: permissionspolicy.FieldWorkforcePolicy,
		})
	}
	if nodes := ppu.mutation.RemovedGroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   permissionspolicy.GroupsTable,
			Columns: permissionspolicy.GroupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersgroup.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ppu.mutation.GroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   permissionspolicy.GroupsTable,
			Columns: permissionspolicy.GroupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersgroup.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ppu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{permissionspolicy.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// PermissionsPolicyUpdateOne is the builder for updating a single PermissionsPolicy entity.
type PermissionsPolicyUpdateOne struct {
	config
	hooks    []Hook
	mutation *PermissionsPolicyMutation
}

// SetName sets the name field.
func (ppuo *PermissionsPolicyUpdateOne) SetName(s string) *PermissionsPolicyUpdateOne {
	ppuo.mutation.SetName(s)
	return ppuo
}

// SetDescription sets the description field.
func (ppuo *PermissionsPolicyUpdateOne) SetDescription(s string) *PermissionsPolicyUpdateOne {
	ppuo.mutation.SetDescription(s)
	return ppuo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ppuo *PermissionsPolicyUpdateOne) SetNillableDescription(s *string) *PermissionsPolicyUpdateOne {
	if s != nil {
		ppuo.SetDescription(*s)
	}
	return ppuo
}

// ClearDescription clears the value of description.
func (ppuo *PermissionsPolicyUpdateOne) ClearDescription() *PermissionsPolicyUpdateOne {
	ppuo.mutation.ClearDescription()
	return ppuo
}

// SetIsGlobal sets the is_global field.
func (ppuo *PermissionsPolicyUpdateOne) SetIsGlobal(b bool) *PermissionsPolicyUpdateOne {
	ppuo.mutation.SetIsGlobal(b)
	return ppuo
}

// SetNillableIsGlobal sets the is_global field if the given value is not nil.
func (ppuo *PermissionsPolicyUpdateOne) SetNillableIsGlobal(b *bool) *PermissionsPolicyUpdateOne {
	if b != nil {
		ppuo.SetIsGlobal(*b)
	}
	return ppuo
}

// ClearIsGlobal clears the value of is_global.
func (ppuo *PermissionsPolicyUpdateOne) ClearIsGlobal() *PermissionsPolicyUpdateOne {
	ppuo.mutation.ClearIsGlobal()
	return ppuo
}

// SetInventoryPolicy sets the inventory_policy field.
func (ppuo *PermissionsPolicyUpdateOne) SetInventoryPolicy(mpi *models.InventoryPolicyInput) *PermissionsPolicyUpdateOne {
	ppuo.mutation.SetInventoryPolicy(mpi)
	return ppuo
}

// ClearInventoryPolicy clears the value of inventory_policy.
func (ppuo *PermissionsPolicyUpdateOne) ClearInventoryPolicy() *PermissionsPolicyUpdateOne {
	ppuo.mutation.ClearInventoryPolicy()
	return ppuo
}

// SetWorkforcePolicy sets the workforce_policy field.
func (ppuo *PermissionsPolicyUpdateOne) SetWorkforcePolicy(mpi *models.WorkforcePolicyInput) *PermissionsPolicyUpdateOne {
	ppuo.mutation.SetWorkforcePolicy(mpi)
	return ppuo
}

// ClearWorkforcePolicy clears the value of workforce_policy.
func (ppuo *PermissionsPolicyUpdateOne) ClearWorkforcePolicy() *PermissionsPolicyUpdateOne {
	ppuo.mutation.ClearWorkforcePolicy()
	return ppuo
}

// AddGroupIDs adds the groups edge to UsersGroup by ids.
func (ppuo *PermissionsPolicyUpdateOne) AddGroupIDs(ids ...int) *PermissionsPolicyUpdateOne {
	ppuo.mutation.AddGroupIDs(ids...)
	return ppuo
}

// AddGroups adds the groups edges to UsersGroup.
func (ppuo *PermissionsPolicyUpdateOne) AddGroups(u ...*UsersGroup) *PermissionsPolicyUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ppuo.AddGroupIDs(ids...)
}

// RemoveGroupIDs removes the groups edge to UsersGroup by ids.
func (ppuo *PermissionsPolicyUpdateOne) RemoveGroupIDs(ids ...int) *PermissionsPolicyUpdateOne {
	ppuo.mutation.RemoveGroupIDs(ids...)
	return ppuo
}

// RemoveGroups removes groups edges to UsersGroup.
func (ppuo *PermissionsPolicyUpdateOne) RemoveGroups(u ...*UsersGroup) *PermissionsPolicyUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ppuo.RemoveGroupIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ppuo *PermissionsPolicyUpdateOne) Save(ctx context.Context) (*PermissionsPolicy, error) {
	if _, ok := ppuo.mutation.UpdateTime(); !ok {
		v := permissionspolicy.UpdateDefaultUpdateTime()
		ppuo.mutation.SetUpdateTime(v)
	}
	if v, ok := ppuo.mutation.Name(); ok {
		if err := permissionspolicy.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err  error
		node *PermissionsPolicy
	)
	if len(ppuo.hooks) == 0 {
		node, err = ppuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PermissionsPolicyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ppuo.mutation = mutation
			node, err = ppuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ppuo.hooks) - 1; i >= 0; i-- {
			mut = ppuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ppuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ppuo *PermissionsPolicyUpdateOne) SaveX(ctx context.Context) *PermissionsPolicy {
	pp, err := ppuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return pp
}

// Exec executes the query on the entity.
func (ppuo *PermissionsPolicyUpdateOne) Exec(ctx context.Context) error {
	_, err := ppuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppuo *PermissionsPolicyUpdateOne) ExecX(ctx context.Context) {
	if err := ppuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ppuo *PermissionsPolicyUpdateOne) sqlSave(ctx context.Context) (pp *PermissionsPolicy, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   permissionspolicy.Table,
			Columns: permissionspolicy.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: permissionspolicy.FieldID,
			},
		},
	}
	id, ok := ppuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing PermissionsPolicy.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := ppuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: permissionspolicy.FieldUpdateTime,
		})
	}
	if value, ok := ppuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: permissionspolicy.FieldName,
		})
	}
	if value, ok := ppuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: permissionspolicy.FieldDescription,
		})
	}
	if ppuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: permissionspolicy.FieldDescription,
		})
	}
	if value, ok := ppuo.mutation.IsGlobal(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: permissionspolicy.FieldIsGlobal,
		})
	}
	if ppuo.mutation.IsGlobalCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: permissionspolicy.FieldIsGlobal,
		})
	}
	if value, ok := ppuo.mutation.InventoryPolicy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: permissionspolicy.FieldInventoryPolicy,
		})
	}
	if ppuo.mutation.InventoryPolicyCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: permissionspolicy.FieldInventoryPolicy,
		})
	}
	if value, ok := ppuo.mutation.WorkforcePolicy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: permissionspolicy.FieldWorkforcePolicy,
		})
	}
	if ppuo.mutation.WorkforcePolicyCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: permissionspolicy.FieldWorkforcePolicy,
		})
	}
	if nodes := ppuo.mutation.RemovedGroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   permissionspolicy.GroupsTable,
			Columns: permissionspolicy.GroupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersgroup.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ppuo.mutation.GroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   permissionspolicy.GroupsTable,
			Columns: permissionspolicy.GroupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersgroup.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	pp = &PermissionsPolicy{config: ppuo.config}
	_spec.Assign = pp.assignValues
	_spec.ScanValues = pp.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ppuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{permissionspolicy.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return pp, nil
}
