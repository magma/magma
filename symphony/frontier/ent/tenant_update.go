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
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
)

// TenantUpdate is the builder for updating Tenant entities.
type TenantUpdate struct {
	config
	hooks      []Hook
	mutation   *TenantMutation
	predicates []predicate.Tenant
}

// Where adds a new predicate for the builder.
func (tu *TenantUpdate) Where(ps ...predicate.Tenant) *TenantUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetName sets the name field.
func (tu *TenantUpdate) SetName(s string) *TenantUpdate {
	tu.mutation.SetName(s)
	return tu
}

// SetDomains sets the domains field.
func (tu *TenantUpdate) SetDomains(s []string) *TenantUpdate {
	tu.mutation.SetDomains(s)
	return tu
}

// SetNetworks sets the networks field.
func (tu *TenantUpdate) SetNetworks(s []string) *TenantUpdate {
	tu.mutation.SetNetworks(s)
	return tu
}

// SetTabs sets the tabs field.
func (tu *TenantUpdate) SetTabs(s []string) *TenantUpdate {
	tu.mutation.SetTabs(s)
	return tu
}

// ClearTabs clears the value of tabs.
func (tu *TenantUpdate) ClearTabs() *TenantUpdate {
	tu.mutation.ClearTabs()
	return tu
}

// SetSSOCert sets the SSOCert field.
func (tu *TenantUpdate) SetSSOCert(s string) *TenantUpdate {
	tu.mutation.SetSSOCert(s)
	return tu
}

// SetNillableSSOCert sets the SSOCert field if the given value is not nil.
func (tu *TenantUpdate) SetNillableSSOCert(s *string) *TenantUpdate {
	if s != nil {
		tu.SetSSOCert(*s)
	}
	return tu
}

// SetSSOEntryPoint sets the SSOEntryPoint field.
func (tu *TenantUpdate) SetSSOEntryPoint(s string) *TenantUpdate {
	tu.mutation.SetSSOEntryPoint(s)
	return tu
}

// SetNillableSSOEntryPoint sets the SSOEntryPoint field if the given value is not nil.
func (tu *TenantUpdate) SetNillableSSOEntryPoint(s *string) *TenantUpdate {
	if s != nil {
		tu.SetSSOEntryPoint(*s)
	}
	return tu
}

// SetSSOIssuer sets the SSOIssuer field.
func (tu *TenantUpdate) SetSSOIssuer(s string) *TenantUpdate {
	tu.mutation.SetSSOIssuer(s)
	return tu
}

// SetNillableSSOIssuer sets the SSOIssuer field if the given value is not nil.
func (tu *TenantUpdate) SetNillableSSOIssuer(s *string) *TenantUpdate {
	if s != nil {
		tu.SetSSOIssuer(*s)
	}
	return tu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (tu *TenantUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := tu.mutation.UpdatedAt(); !ok {
		v := tenant.UpdateDefaultUpdatedAt()
		tu.mutation.SetUpdatedAt(v)
	}
	if v, ok := tu.mutation.Name(); ok {
		if err := tenant.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	var (
		err      error
		affected int
	)
	if len(tu.hooks) == 0 {
		affected, err = tu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TenantMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tu.mutation = mutation
			affected, err = tu.sqlSave(ctx)
			return affected, err
		})
		for i := len(tu.hooks); i > 0; i-- {
			mut = tu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, tu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TenantUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TenantUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TenantUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TenantUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tenant.Table,
			Columns: tenant.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tenant.FieldID,
			},
		},
	}
	if ps := tu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: tenant.FieldUpdatedAt,
		})
	}
	if value, ok := tu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldName,
		})
	}
	if value, ok := tu.mutation.Domains(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldDomains,
		})
	}
	if value, ok := tu.mutation.Networks(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldNetworks,
		})
	}
	if value, ok := tu.mutation.Tabs(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldTabs,
		})
	}
	if tu.mutation.TabsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: tenant.FieldTabs,
		})
	}
	if value, ok := tu.mutation.SSOCert(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOCert,
		})
	}
	if value, ok := tu.mutation.SSOEntryPoint(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOEntryPoint,
		})
	}
	if value, ok := tu.mutation.SSOIssuer(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOIssuer,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tenant.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// TenantUpdateOne is the builder for updating a single Tenant entity.
type TenantUpdateOne struct {
	config
	hooks    []Hook
	mutation *TenantMutation
}

// SetName sets the name field.
func (tuo *TenantUpdateOne) SetName(s string) *TenantUpdateOne {
	tuo.mutation.SetName(s)
	return tuo
}

// SetDomains sets the domains field.
func (tuo *TenantUpdateOne) SetDomains(s []string) *TenantUpdateOne {
	tuo.mutation.SetDomains(s)
	return tuo
}

// SetNetworks sets the networks field.
func (tuo *TenantUpdateOne) SetNetworks(s []string) *TenantUpdateOne {
	tuo.mutation.SetNetworks(s)
	return tuo
}

// SetTabs sets the tabs field.
func (tuo *TenantUpdateOne) SetTabs(s []string) *TenantUpdateOne {
	tuo.mutation.SetTabs(s)
	return tuo
}

// ClearTabs clears the value of tabs.
func (tuo *TenantUpdateOne) ClearTabs() *TenantUpdateOne {
	tuo.mutation.ClearTabs()
	return tuo
}

// SetSSOCert sets the SSOCert field.
func (tuo *TenantUpdateOne) SetSSOCert(s string) *TenantUpdateOne {
	tuo.mutation.SetSSOCert(s)
	return tuo
}

// SetNillableSSOCert sets the SSOCert field if the given value is not nil.
func (tuo *TenantUpdateOne) SetNillableSSOCert(s *string) *TenantUpdateOne {
	if s != nil {
		tuo.SetSSOCert(*s)
	}
	return tuo
}

// SetSSOEntryPoint sets the SSOEntryPoint field.
func (tuo *TenantUpdateOne) SetSSOEntryPoint(s string) *TenantUpdateOne {
	tuo.mutation.SetSSOEntryPoint(s)
	return tuo
}

// SetNillableSSOEntryPoint sets the SSOEntryPoint field if the given value is not nil.
func (tuo *TenantUpdateOne) SetNillableSSOEntryPoint(s *string) *TenantUpdateOne {
	if s != nil {
		tuo.SetSSOEntryPoint(*s)
	}
	return tuo
}

// SetSSOIssuer sets the SSOIssuer field.
func (tuo *TenantUpdateOne) SetSSOIssuer(s string) *TenantUpdateOne {
	tuo.mutation.SetSSOIssuer(s)
	return tuo
}

// SetNillableSSOIssuer sets the SSOIssuer field if the given value is not nil.
func (tuo *TenantUpdateOne) SetNillableSSOIssuer(s *string) *TenantUpdateOne {
	if s != nil {
		tuo.SetSSOIssuer(*s)
	}
	return tuo
}

// Save executes the query and returns the updated entity.
func (tuo *TenantUpdateOne) Save(ctx context.Context) (*Tenant, error) {
	if _, ok := tuo.mutation.UpdatedAt(); !ok {
		v := tenant.UpdateDefaultUpdatedAt()
		tuo.mutation.SetUpdatedAt(v)
	}
	if v, ok := tuo.mutation.Name(); ok {
		if err := tenant.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	var (
		err  error
		node *Tenant
	)
	if len(tuo.hooks) == 0 {
		node, err = tuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TenantMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tuo.mutation = mutation
			node, err = tuo.sqlSave(ctx)
			return node, err
		})
		for i := len(tuo.hooks); i > 0; i-- {
			mut = tuo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, tuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TenantUpdateOne) SaveX(ctx context.Context) *Tenant {
	t, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return t
}

// Exec executes the query on the entity.
func (tuo *TenantUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TenantUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TenantUpdateOne) sqlSave(ctx context.Context) (t *Tenant, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tenant.Table,
			Columns: tenant.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tenant.FieldID,
			},
		},
	}
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Tenant.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := tuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: tenant.FieldUpdatedAt,
		})
	}
	if value, ok := tuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldName,
		})
	}
	if value, ok := tuo.mutation.Domains(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldDomains,
		})
	}
	if value, ok := tuo.mutation.Networks(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldNetworks,
		})
	}
	if value, ok := tuo.mutation.Tabs(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldTabs,
		})
	}
	if tuo.mutation.TabsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: tenant.FieldTabs,
		})
	}
	if value, ok := tuo.mutation.SSOCert(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOCert,
		})
	}
	if value, ok := tuo.mutation.SSOEntryPoint(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOEntryPoint,
		})
	}
	if value, ok := tuo.mutation.SSOIssuer(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOIssuer,
		})
	}
	t = &Tenant{config: tuo.config}
	_spec.Assign = t.assignValues
	_spec.ScanValues = t.scanValues()
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tenant.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return t, nil
}
