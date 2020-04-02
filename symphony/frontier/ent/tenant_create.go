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
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
)

// TenantCreate is the builder for creating a Tenant entity.
type TenantCreate struct {
	config
	mutation *TenantMutation
	hooks    []Hook
}

// SetCreatedAt sets the created_at field.
func (tc *TenantCreate) SetCreatedAt(t time.Time) *TenantCreate {
	tc.mutation.SetCreatedAt(t)
	return tc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (tc *TenantCreate) SetNillableCreatedAt(t *time.Time) *TenantCreate {
	if t != nil {
		tc.SetCreatedAt(*t)
	}
	return tc
}

// SetUpdatedAt sets the updated_at field.
func (tc *TenantCreate) SetUpdatedAt(t time.Time) *TenantCreate {
	tc.mutation.SetUpdatedAt(t)
	return tc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (tc *TenantCreate) SetNillableUpdatedAt(t *time.Time) *TenantCreate {
	if t != nil {
		tc.SetUpdatedAt(*t)
	}
	return tc
}

// SetName sets the name field.
func (tc *TenantCreate) SetName(s string) *TenantCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetDomains sets the domains field.
func (tc *TenantCreate) SetDomains(s []string) *TenantCreate {
	tc.mutation.SetDomains(s)
	return tc
}

// SetNetworks sets the networks field.
func (tc *TenantCreate) SetNetworks(s []string) *TenantCreate {
	tc.mutation.SetNetworks(s)
	return tc
}

// SetTabs sets the tabs field.
func (tc *TenantCreate) SetTabs(s []string) *TenantCreate {
	tc.mutation.SetTabs(s)
	return tc
}

// SetSSOCert sets the SSOCert field.
func (tc *TenantCreate) SetSSOCert(s string) *TenantCreate {
	tc.mutation.SetSSOCert(s)
	return tc
}

// SetNillableSSOCert sets the SSOCert field if the given value is not nil.
func (tc *TenantCreate) SetNillableSSOCert(s *string) *TenantCreate {
	if s != nil {
		tc.SetSSOCert(*s)
	}
	return tc
}

// SetSSOEntryPoint sets the SSOEntryPoint field.
func (tc *TenantCreate) SetSSOEntryPoint(s string) *TenantCreate {
	tc.mutation.SetSSOEntryPoint(s)
	return tc
}

// SetNillableSSOEntryPoint sets the SSOEntryPoint field if the given value is not nil.
func (tc *TenantCreate) SetNillableSSOEntryPoint(s *string) *TenantCreate {
	if s != nil {
		tc.SetSSOEntryPoint(*s)
	}
	return tc
}

// SetSSOIssuer sets the SSOIssuer field.
func (tc *TenantCreate) SetSSOIssuer(s string) *TenantCreate {
	tc.mutation.SetSSOIssuer(s)
	return tc
}

// SetNillableSSOIssuer sets the SSOIssuer field if the given value is not nil.
func (tc *TenantCreate) SetNillableSSOIssuer(s *string) *TenantCreate {
	if s != nil {
		tc.SetSSOIssuer(*s)
	}
	return tc
}

// Save creates the Tenant in the database.
func (tc *TenantCreate) Save(ctx context.Context) (*Tenant, error) {
	if _, ok := tc.mutation.CreatedAt(); !ok {
		v := tenant.DefaultCreatedAt()
		tc.mutation.SetCreatedAt(v)
	}
	if _, ok := tc.mutation.UpdatedAt(); !ok {
		v := tenant.DefaultUpdatedAt()
		tc.mutation.SetUpdatedAt(v)
	}
	if _, ok := tc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := tc.mutation.Name(); ok {
		if err := tenant.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if _, ok := tc.mutation.Domains(); !ok {
		return nil, errors.New("ent: missing required field \"domains\"")
	}
	if _, ok := tc.mutation.Networks(); !ok {
		return nil, errors.New("ent: missing required field \"networks\"")
	}
	if _, ok := tc.mutation.SSOCert(); !ok {
		v := tenant.DefaultSSOCert
		tc.mutation.SetSSOCert(v)
	}
	if _, ok := tc.mutation.SSOEntryPoint(); !ok {
		v := tenant.DefaultSSOEntryPoint
		tc.mutation.SetSSOEntryPoint(v)
	}
	if _, ok := tc.mutation.SSOIssuer(); !ok {
		v := tenant.DefaultSSOIssuer
		tc.mutation.SetSSOIssuer(v)
	}
	var (
		err  error
		node *Tenant
	)
	if len(tc.hooks) == 0 {
		node, err = tc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TenantMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tc.mutation = mutation
			node, err = tc.sqlSave(ctx)
			return node, err
		})
		for i := len(tc.hooks) - 1; i >= 0; i-- {
			mut = tc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TenantCreate) SaveX(ctx context.Context) *Tenant {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tc *TenantCreate) sqlSave(ctx context.Context) (*Tenant, error) {
	var (
		t     = &Tenant{config: tc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: tenant.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tenant.FieldID,
			},
		}
	)
	if value, ok := tc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: tenant.FieldCreatedAt,
		})
		t.CreatedAt = value
	}
	if value, ok := tc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: tenant.FieldUpdatedAt,
		})
		t.UpdatedAt = value
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldName,
		})
		t.Name = value
	}
	if value, ok := tc.mutation.Domains(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldDomains,
		})
		t.Domains = value
	}
	if value, ok := tc.mutation.Networks(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldNetworks,
		})
		t.Networks = value
	}
	if value, ok := tc.mutation.Tabs(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: tenant.FieldTabs,
		})
		t.Tabs = value
	}
	if value, ok := tc.mutation.SSOCert(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOCert,
		})
		t.SSOCert = value
	}
	if value, ok := tc.mutation.SSOEntryPoint(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOEntryPoint,
		})
		t.SSOEntryPoint = value
	}
	if value, ok := tc.mutation.SSOIssuer(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tenant.FieldSSOIssuer,
		})
		t.SSOIssuer = value
	}
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	t.ID = int(id)
	return t, nil
}
