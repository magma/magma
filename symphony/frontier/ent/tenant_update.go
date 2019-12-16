// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
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
	var (
		builder  = sql.Dialect(tu.driver.Dialect())
		selector = builder.Select(tenant.FieldID).From(builder.Table(tenant.Table))
	)
	for _, p := range tu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = tu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := tu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(tenant.Table)
	)
	updater = updater.Where(sql.InInts(tenant.FieldID, ids...))
	if value := tu.updated_at; value != nil {
		updater.Set(tenant.FieldUpdatedAt, *value)
	}
	if value := tu.name; value != nil {
		updater.Set(tenant.FieldName, *value)
	}
	if value := tu.domains; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(tenant.FieldDomains, buf)
	}
	if value := tu.networks; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(tenant.FieldNetworks, buf)
	}
	if value := tu.tabs; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(tenant.FieldTabs, buf)
	}
	if tu.cleartabs {
		updater.SetNull(tenant.FieldTabs)
	}
	if value := tu.SSOCert; value != nil {
		updater.Set(tenant.FieldSSOCert, *value)
	}
	if value := tu.SSOEntryPoint; value != nil {
		updater.Set(tenant.FieldSSOEntryPoint, *value)
	}
	if value := tu.SSOIssuer; value != nil {
		updater.Set(tenant.FieldSSOIssuer, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
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
	var (
		builder  = sql.Dialect(tuo.driver.Dialect())
		selector = builder.Select(tenant.Columns...).From(builder.Table(tenant.Table))
	)
	tenant.ID(tuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = tuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		t = &Tenant{config: tuo.config}
		if err := t.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Tenant: %v", err)
		}
		id = t.ID
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Tenant with id: %v", tuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Tenant with the same id: %v", tuo.id)
	}

	tx, err := tuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(tenant.Table)
	)
	updater = updater.Where(sql.InInts(tenant.FieldID, ids...))
	if value := tuo.updated_at; value != nil {
		updater.Set(tenant.FieldUpdatedAt, *value)
		t.UpdatedAt = *value
	}
	if value := tuo.name; value != nil {
		updater.Set(tenant.FieldName, *value)
		t.Name = *value
	}
	if value := tuo.domains; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(tenant.FieldDomains, buf)
		t.Domains = *value
	}
	if value := tuo.networks; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(tenant.FieldNetworks, buf)
		t.Networks = *value
	}
	if value := tuo.tabs; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(tenant.FieldTabs, buf)
		t.Tabs = *value
	}
	if tuo.cleartabs {
		var value []string
		t.Tabs = value
		updater.SetNull(tenant.FieldTabs)
	}
	if value := tuo.SSOCert; value != nil {
		updater.Set(tenant.FieldSSOCert, *value)
		t.SSOCert = *value
	}
	if value := tuo.SSOEntryPoint; value != nil {
		updater.Set(tenant.FieldSSOEntryPoint, *value)
		t.SSOEntryPoint = *value
	}
	if value := tuo.SSOIssuer; value != nil {
		updater.Set(tenant.FieldSSOIssuer, *value)
		t.SSOIssuer = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}
