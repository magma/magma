// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/customer"
)

// CustomerCreate is the builder for creating a Customer entity.
type CustomerCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	external_id *string
	services    map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (cc *CustomerCreate) SetCreateTime(t time.Time) *CustomerCreate {
	cc.create_time = &t
	return cc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (cc *CustomerCreate) SetNillableCreateTime(t *time.Time) *CustomerCreate {
	if t != nil {
		cc.SetCreateTime(*t)
	}
	return cc
}

// SetUpdateTime sets the update_time field.
func (cc *CustomerCreate) SetUpdateTime(t time.Time) *CustomerCreate {
	cc.update_time = &t
	return cc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (cc *CustomerCreate) SetNillableUpdateTime(t *time.Time) *CustomerCreate {
	if t != nil {
		cc.SetUpdateTime(*t)
	}
	return cc
}

// SetName sets the name field.
func (cc *CustomerCreate) SetName(s string) *CustomerCreate {
	cc.name = &s
	return cc
}

// SetExternalID sets the external_id field.
func (cc *CustomerCreate) SetExternalID(s string) *CustomerCreate {
	cc.external_id = &s
	return cc
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (cc *CustomerCreate) SetNillableExternalID(s *string) *CustomerCreate {
	if s != nil {
		cc.SetExternalID(*s)
	}
	return cc
}

// AddServiceIDs adds the services edge to Service by ids.
func (cc *CustomerCreate) AddServiceIDs(ids ...string) *CustomerCreate {
	if cc.services == nil {
		cc.services = make(map[string]struct{})
	}
	for i := range ids {
		cc.services[ids[i]] = struct{}{}
	}
	return cc
}

// AddServices adds the services edges to Service.
func (cc *CustomerCreate) AddServices(s ...*Service) *CustomerCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cc.AddServiceIDs(ids...)
}

// Save creates the Customer in the database.
func (cc *CustomerCreate) Save(ctx context.Context) (*Customer, error) {
	if cc.create_time == nil {
		v := customer.DefaultCreateTime()
		cc.create_time = &v
	}
	if cc.update_time == nil {
		v := customer.DefaultUpdateTime()
		cc.update_time = &v
	}
	if cc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := customer.NameValidator(*cc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if cc.external_id != nil {
		if err := customer.ExternalIDValidator(*cc.external_id); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	return cc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CustomerCreate) SaveX(ctx context.Context) *Customer {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (cc *CustomerCreate) sqlSave(ctx context.Context) (*Customer, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(cc.driver.Dialect())
		c       = &Customer{config: cc.config}
	)
	tx, err := cc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(customer.Table).Default()
	if value := cc.create_time; value != nil {
		insert.Set(customer.FieldCreateTime, *value)
		c.CreateTime = *value
	}
	if value := cc.update_time; value != nil {
		insert.Set(customer.FieldUpdateTime, *value)
		c.UpdateTime = *value
	}
	if value := cc.name; value != nil {
		insert.Set(customer.FieldName, *value)
		c.Name = *value
	}
	if value := cc.external_id; value != nil {
		insert.Set(customer.FieldExternalID, *value)
		c.ExternalID = value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(customer.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	c.ID = strconv.FormatInt(id, 10)
	if len(cc.services) > 0 {
		for eid := range cc.services {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(customer.ServicesTable).
				Columns(customer.ServicesPrimaryKey[1], customer.ServicesPrimaryKey[0]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return c, nil
}
