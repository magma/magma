// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
)

// TenantUpdate is the builder for updating Tenant entities.
type TenantUpdate struct {
	config

	updated_at    *time.Time
	name          *string
	domains       *[]string
	networks      *[]string
	tabs          *[]string
	cleartabs     bool
	SSOCert       *string
	SSOEntryPoint *string
	SSOIssuer     *string
	predicates    []predicate.Tenant
}

// Where adds a new predicate for the builder.
func (tu *TenantUpdate) Where(ps ...predicate.Tenant) *TenantUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetName sets the name field.
func (tu *TenantUpdate) SetName(s string) *TenantUpdate {
	tu.name = &s
	return tu
}

// SetDomains sets the domains field.
func (tu *TenantUpdate) SetDomains(s []string) *TenantUpdate {
	tu.domains = &s
	return tu
}

// SetNetworks sets the networks field.
func (tu *TenantUpdate) SetNetworks(s []string) *TenantUpdate {
	tu.networks = &s
	return tu
}

// SetTabs sets the tabs field.
func (tu *TenantUpdate) SetTabs(s []string) *TenantUpdate {
	tu.tabs = &s
	return tu
}

// ClearTabs clears the value of tabs.
func (tu *TenantUpdate) ClearTabs() *TenantUpdate {
	tu.tabs = nil
	tu.cleartabs = true
	return tu
}

// SetSSOCert sets the SSOCert field.
func (tu *TenantUpdate) SetSSOCert(s string) *TenantUpdate {
	tu.SSOCert = &s
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
	tu.SSOEntryPoint = &s
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
	tu.SSOIssuer = &s
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
	if tu.updated_at == nil {
		v := tenant.UpdateDefaultUpdatedAt()
		tu.updated_at = &v
	}
	if tu.name != nil {
		if err := tenant.NameValidator(*tu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	return tu.sqlSave(ctx)
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
	if value := tu.updated_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: tenant.FieldUpdatedAt,
		})
	}
	if value := tu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: tenant.FieldName,
		})
	}
	if value := tu.domains; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: tenant.FieldDomains,
		})
	}
	if value := tu.networks; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: tenant.FieldNetworks,
		})
	}
	if value := tu.tabs; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: tenant.FieldTabs,
		})
	}
	if tu.cleartabs {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: tenant.FieldTabs,
		})
	}
	if value := tu.SSOCert; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: tenant.FieldSSOCert,
		})
	}
	if value := tu.SSOEntryPoint; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: tenant.FieldSSOEntryPoint,
		})
	}
	if value := tu.SSOIssuer; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
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
	id int

	updated_at    *time.Time
	name          *string
	domains       *[]string
	networks      *[]string
	tabs          *[]string
	cleartabs     bool
	SSOCert       *string
	SSOEntryPoint *string
	SSOIssuer     *string
}

// SetName sets the name field.
func (tuo *TenantUpdateOne) SetName(s string) *TenantUpdateOne {
	tuo.name = &s
	return tuo
}

// SetDomains sets the domains field.
func (tuo *TenantUpdateOne) SetDomains(s []string) *TenantUpdateOne {
	tuo.domains = &s
	return tuo
}

// SetNetworks sets the networks field.
func (tuo *TenantUpdateOne) SetNetworks(s []string) *TenantUpdateOne {
	tuo.networks = &s
	return tuo
}

// SetTabs sets the tabs field.
func (tuo *TenantUpdateOne) SetTabs(s []string) *TenantUpdateOne {
	tuo.tabs = &s
	return tuo
}

// ClearTabs clears the value of tabs.
func (tuo *TenantUpdateOne) ClearTabs() *TenantUpdateOne {
	tuo.tabs = nil
	tuo.cleartabs = true
	return tuo
}

// SetSSOCert sets the SSOCert field.
func (tuo *TenantUpdateOne) SetSSOCert(s string) *TenantUpdateOne {
	tuo.SSOCert = &s
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
	tuo.SSOEntryPoint = &s
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
	tuo.SSOIssuer = &s
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
	if tuo.updated_at == nil {
		v := tenant.UpdateDefaultUpdatedAt()
		tuo.updated_at = &v
	}
	if tuo.name != nil {
		if err := tenant.NameValidator(*tuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	return tuo.sqlSave(ctx)
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
				Value:  tuo.id,
				Type:   field.TypeInt,
				Column: tenant.FieldID,
			},
		},
	}
	if value := tuo.updated_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: tenant.FieldUpdatedAt,
		})
	}
	if value := tuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: tenant.FieldName,
		})
	}
	if value := tuo.domains; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: tenant.FieldDomains,
		})
	}
	if value := tuo.networks; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: tenant.FieldNetworks,
		})
	}
	if value := tuo.tabs; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: tenant.FieldTabs,
		})
	}
	if tuo.cleartabs {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: tenant.FieldTabs,
		})
	}
	if value := tuo.SSOCert; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: tenant.FieldSSOCert,
		})
	}
	if value := tuo.SSOEntryPoint; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: tenant.FieldSSOEntryPoint,
		})
	}
	if value := tuo.SSOIssuer; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
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
