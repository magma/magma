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
	"github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
)

// PermissionsPolicyCreate is the builder for creating a PermissionsPolicy entity.
type PermissionsPolicyCreate struct {
	config
	mutation *PermissionsPolicyMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (ppc *PermissionsPolicyCreate) SetCreateTime(t time.Time) *PermissionsPolicyCreate {
	ppc.mutation.SetCreateTime(t)
	return ppc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ppc *PermissionsPolicyCreate) SetNillableCreateTime(t *time.Time) *PermissionsPolicyCreate {
	if t != nil {
		ppc.SetCreateTime(*t)
	}
	return ppc
}

// SetUpdateTime sets the update_time field.
func (ppc *PermissionsPolicyCreate) SetUpdateTime(t time.Time) *PermissionsPolicyCreate {
	ppc.mutation.SetUpdateTime(t)
	return ppc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ppc *PermissionsPolicyCreate) SetNillableUpdateTime(t *time.Time) *PermissionsPolicyCreate {
	if t != nil {
		ppc.SetUpdateTime(*t)
	}
	return ppc
}

// SetName sets the name field.
func (ppc *PermissionsPolicyCreate) SetName(s string) *PermissionsPolicyCreate {
	ppc.mutation.SetName(s)
	return ppc
}

// SetDescription sets the description field.
func (ppc *PermissionsPolicyCreate) SetDescription(s string) *PermissionsPolicyCreate {
	ppc.mutation.SetDescription(s)
	return ppc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ppc *PermissionsPolicyCreate) SetNillableDescription(s *string) *PermissionsPolicyCreate {
	if s != nil {
		ppc.SetDescription(*s)
	}
	return ppc
}

// SetIsGlobal sets the is_global field.
func (ppc *PermissionsPolicyCreate) SetIsGlobal(b bool) *PermissionsPolicyCreate {
	ppc.mutation.SetIsGlobal(b)
	return ppc
}

// SetNillableIsGlobal sets the is_global field if the given value is not nil.
func (ppc *PermissionsPolicyCreate) SetNillableIsGlobal(b *bool) *PermissionsPolicyCreate {
	if b != nil {
		ppc.SetIsGlobal(*b)
	}
	return ppc
}

// SetInventoryPolicy sets the inventory_policy field.
func (ppc *PermissionsPolicyCreate) SetInventoryPolicy(mpi *models.InventoryPolicyInput) *PermissionsPolicyCreate {
	ppc.mutation.SetInventoryPolicy(mpi)
	return ppc
}

// SetWorkforcePolicy sets the workforce_policy field.
func (ppc *PermissionsPolicyCreate) SetWorkforcePolicy(mpi *models.WorkforcePolicyInput) *PermissionsPolicyCreate {
	ppc.mutation.SetWorkforcePolicy(mpi)
	return ppc
}

// AddGroupIDs adds the groups edge to UsersGroup by ids.
func (ppc *PermissionsPolicyCreate) AddGroupIDs(ids ...int) *PermissionsPolicyCreate {
	ppc.mutation.AddGroupIDs(ids...)
	return ppc
}

// AddGroups adds the groups edges to UsersGroup.
func (ppc *PermissionsPolicyCreate) AddGroups(u ...*UsersGroup) *PermissionsPolicyCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ppc.AddGroupIDs(ids...)
}

// Save creates the PermissionsPolicy in the database.
func (ppc *PermissionsPolicyCreate) Save(ctx context.Context) (*PermissionsPolicy, error) {
	if _, ok := ppc.mutation.CreateTime(); !ok {
		v := permissionspolicy.DefaultCreateTime()
		ppc.mutation.SetCreateTime(v)
	}
	if _, ok := ppc.mutation.UpdateTime(); !ok {
		v := permissionspolicy.DefaultUpdateTime()
		ppc.mutation.SetUpdateTime(v)
	}
	if _, ok := ppc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := ppc.mutation.Name(); ok {
		if err := permissionspolicy.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if _, ok := ppc.mutation.IsGlobal(); !ok {
		v := permissionspolicy.DefaultIsGlobal
		ppc.mutation.SetIsGlobal(v)
	}
	var (
		err  error
		node *PermissionsPolicy
	)
	if len(ppc.hooks) == 0 {
		node, err = ppc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PermissionsPolicyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ppc.mutation = mutation
			node, err = ppc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ppc.hooks) - 1; i >= 0; i-- {
			mut = ppc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ppc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ppc *PermissionsPolicyCreate) SaveX(ctx context.Context) *PermissionsPolicy {
	v, err := ppc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ppc *PermissionsPolicyCreate) sqlSave(ctx context.Context) (*PermissionsPolicy, error) {
	var (
		pp    = &PermissionsPolicy{config: ppc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: permissionspolicy.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: permissionspolicy.FieldID,
			},
		}
	)
	if value, ok := ppc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: permissionspolicy.FieldCreateTime,
		})
		pp.CreateTime = value
	}
	if value, ok := ppc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: permissionspolicy.FieldUpdateTime,
		})
		pp.UpdateTime = value
	}
	if value, ok := ppc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: permissionspolicy.FieldName,
		})
		pp.Name = value
	}
	if value, ok := ppc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: permissionspolicy.FieldDescription,
		})
		pp.Description = value
	}
	if value, ok := ppc.mutation.IsGlobal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: permissionspolicy.FieldIsGlobal,
		})
		pp.IsGlobal = value
	}
	if value, ok := ppc.mutation.InventoryPolicy(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: permissionspolicy.FieldInventoryPolicy,
		})
		pp.InventoryPolicy = value
	}
	if value, ok := ppc.mutation.WorkforcePolicy(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: permissionspolicy.FieldWorkforcePolicy,
		})
		pp.WorkforcePolicy = value
	}
	if nodes := ppc.mutation.GroupsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ppc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pp.ID = int(id)
	return pp, nil
}
