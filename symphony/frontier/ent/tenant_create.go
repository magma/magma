// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
)

// TenantCreate is the builder for creating a Tenant entity.
type TenantCreate struct {
	config
	created_at    *time.Time
	updated_at    *time.Time
	name          *string
	domains       *[]string
	networks      *[]string
	tabs          *[]string
	SSOCert       *string
	SSOEntryPoint *string
	SSOIssuer     *string
}

// SetCreatedAt sets the created_at field.
func (tc *TenantCreate) SetCreatedAt(t time.Time) *TenantCreate {
	tc.created_at = &t
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
	tc.updated_at = &t
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
	tc.name = &s
	return tc
}

// SetDomains sets the domains field.
func (tc *TenantCreate) SetDomains(s []string) *TenantCreate {
	tc.domains = &s
	return tc
}

// SetNetworks sets the networks field.
func (tc *TenantCreate) SetNetworks(s []string) *TenantCreate {
	tc.networks = &s
	return tc
}

// SetTabs sets the tabs field.
func (tc *TenantCreate) SetTabs(s []string) *TenantCreate {
	tc.tabs = &s
	return tc
}

// SetSSOCert sets the SSOCert field.
func (tc *TenantCreate) SetSSOCert(s string) *TenantCreate {
	tc.SSOCert = &s
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
	tc.SSOEntryPoint = &s
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
	tc.SSOIssuer = &s
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
	if tc.created_at == nil {
		v := tenant.DefaultCreatedAt()
		tc.created_at = &v
	}
	if tc.updated_at == nil {
		v := tenant.DefaultUpdatedAt()
		tc.updated_at = &v
	}
	if tc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := tenant.NameValidator(*tc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if tc.domains == nil {
		return nil, errors.New("ent: missing required field \"domains\"")
	}
	if tc.networks == nil {
		return nil, errors.New("ent: missing required field \"networks\"")
	}
	if tc.SSOCert == nil {
		v := tenant.DefaultSSOCert
		tc.SSOCert = &v
	}
	if tc.SSOEntryPoint == nil {
		v := tenant.DefaultSSOEntryPoint
		tc.SSOEntryPoint = &v
	}
	if tc.SSOIssuer == nil {
		v := tenant.DefaultSSOIssuer
		tc.SSOIssuer = &v
	}
	return tc.sqlSave(ctx)
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
		builder = sql.Dialect(tc.driver.Dialect())
		t       = &Tenant{config: tc.config}
	)
	tx, err := tc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(tenant.Table).Default()
	if value := tc.created_at; value != nil {
		insert.Set(tenant.FieldCreatedAt, *value)
		t.CreatedAt = *value
	}
	if value := tc.updated_at; value != nil {
		insert.Set(tenant.FieldUpdatedAt, *value)
		t.UpdatedAt = *value
	}
	if value := tc.name; value != nil {
		insert.Set(tenant.FieldName, *value)
		t.Name = *value
	}
	if value := tc.domains; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(tenant.FieldDomains, buf)
		t.Domains = *value
	}
	if value := tc.networks; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(tenant.FieldNetworks, buf)
		t.Networks = *value
	}
	if value := tc.tabs; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(tenant.FieldTabs, buf)
		t.Tabs = *value
	}
	if value := tc.SSOCert; value != nil {
		insert.Set(tenant.FieldSSOCert, *value)
		t.SSOCert = *value
	}
	if value := tc.SSOEntryPoint; value != nil {
		insert.Set(tenant.FieldSSOEntryPoint, *value)
		t.SSOEntryPoint = *value
	}
	if value := tc.SSOIssuer; value != nil {
		insert.Set(tenant.FieldSSOIssuer, *value)
		t.SSOIssuer = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(tenant.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	t.ID = int(id)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}
