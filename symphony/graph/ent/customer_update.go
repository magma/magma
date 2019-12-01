// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CustomerUpdate is the builder for updating Customer entities.
type CustomerUpdate struct {
	config

	update_time      *time.Time
	name             *string
	external_id      *string
	clearexternal_id bool
	services         map[string]struct{}
	removedServices  map[string]struct{}
	predicates       []predicate.Customer
}

// Where adds a new predicate for the builder.
func (cu *CustomerUpdate) Where(ps ...predicate.Customer) *CustomerUpdate {
	cu.predicates = append(cu.predicates, ps...)
	return cu
}

// SetName sets the name field.
func (cu *CustomerUpdate) SetName(s string) *CustomerUpdate {
	cu.name = &s
	return cu
}

// SetExternalID sets the external_id field.
func (cu *CustomerUpdate) SetExternalID(s string) *CustomerUpdate {
	cu.external_id = &s
	return cu
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (cu *CustomerUpdate) SetNillableExternalID(s *string) *CustomerUpdate {
	if s != nil {
		cu.SetExternalID(*s)
	}
	return cu
}

// ClearExternalID clears the value of external_id.
func (cu *CustomerUpdate) ClearExternalID() *CustomerUpdate {
	cu.external_id = nil
	cu.clearexternal_id = true
	return cu
}

// AddServiceIDs adds the services edge to Service by ids.
func (cu *CustomerUpdate) AddServiceIDs(ids ...string) *CustomerUpdate {
	if cu.services == nil {
		cu.services = make(map[string]struct{})
	}
	for i := range ids {
		cu.services[ids[i]] = struct{}{}
	}
	return cu
}

// AddServices adds the services edges to Service.
func (cu *CustomerUpdate) AddServices(s ...*Service) *CustomerUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cu.AddServiceIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (cu *CustomerUpdate) RemoveServiceIDs(ids ...string) *CustomerUpdate {
	if cu.removedServices == nil {
		cu.removedServices = make(map[string]struct{})
	}
	for i := range ids {
		cu.removedServices[ids[i]] = struct{}{}
	}
	return cu
}

// RemoveServices removes services edges to Service.
func (cu *CustomerUpdate) RemoveServices(s ...*Service) *CustomerUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cu.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (cu *CustomerUpdate) Save(ctx context.Context) (int, error) {
	if cu.update_time == nil {
		v := customer.UpdateDefaultUpdateTime()
		cu.update_time = &v
	}
	if cu.name != nil {
		if err := customer.NameValidator(*cu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if cu.external_id != nil {
		if err := customer.ExternalIDValidator(*cu.external_id); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	return cu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CustomerUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CustomerUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CustomerUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *CustomerUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(cu.driver.Dialect())
		selector = builder.Select(customer.FieldID).From(builder.Table(customer.Table))
	)
	for _, p := range cu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := cu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(customer.Table)
	)
	updater = updater.Where(sql.InInts(customer.FieldID, ids...))
	if value := cu.update_time; value != nil {
		updater.Set(customer.FieldUpdateTime, *value)
	}
	if value := cu.name; value != nil {
		updater.Set(customer.FieldName, *value)
	}
	if value := cu.external_id; value != nil {
		updater.Set(customer.FieldExternalID, *value)
	}
	if cu.clearexternal_id {
		updater.SetNull(customer.FieldExternalID)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(cu.removedServices) > 0 {
		eids := make([]int, len(cu.removedServices))
		for eid := range cu.removedServices {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(customer.ServicesTable).
			Where(sql.InInts(customer.ServicesPrimaryKey[1], ids...)).
			Where(sql.InInts(customer.ServicesPrimaryKey[0], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(cu.services) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range cu.services {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(customer.ServicesTable).
			Columns(customer.ServicesPrimaryKey[1], customer.ServicesPrimaryKey[0])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// CustomerUpdateOne is the builder for updating a single Customer entity.
type CustomerUpdateOne struct {
	config
	id string

	update_time      *time.Time
	name             *string
	external_id      *string
	clearexternal_id bool
	services         map[string]struct{}
	removedServices  map[string]struct{}
}

// SetName sets the name field.
func (cuo *CustomerUpdateOne) SetName(s string) *CustomerUpdateOne {
	cuo.name = &s
	return cuo
}

// SetExternalID sets the external_id field.
func (cuo *CustomerUpdateOne) SetExternalID(s string) *CustomerUpdateOne {
	cuo.external_id = &s
	return cuo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (cuo *CustomerUpdateOne) SetNillableExternalID(s *string) *CustomerUpdateOne {
	if s != nil {
		cuo.SetExternalID(*s)
	}
	return cuo
}

// ClearExternalID clears the value of external_id.
func (cuo *CustomerUpdateOne) ClearExternalID() *CustomerUpdateOne {
	cuo.external_id = nil
	cuo.clearexternal_id = true
	return cuo
}

// AddServiceIDs adds the services edge to Service by ids.
func (cuo *CustomerUpdateOne) AddServiceIDs(ids ...string) *CustomerUpdateOne {
	if cuo.services == nil {
		cuo.services = make(map[string]struct{})
	}
	for i := range ids {
		cuo.services[ids[i]] = struct{}{}
	}
	return cuo
}

// AddServices adds the services edges to Service.
func (cuo *CustomerUpdateOne) AddServices(s ...*Service) *CustomerUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cuo.AddServiceIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (cuo *CustomerUpdateOne) RemoveServiceIDs(ids ...string) *CustomerUpdateOne {
	if cuo.removedServices == nil {
		cuo.removedServices = make(map[string]struct{})
	}
	for i := range ids {
		cuo.removedServices[ids[i]] = struct{}{}
	}
	return cuo
}

// RemoveServices removes services edges to Service.
func (cuo *CustomerUpdateOne) RemoveServices(s ...*Service) *CustomerUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cuo.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (cuo *CustomerUpdateOne) Save(ctx context.Context) (*Customer, error) {
	if cuo.update_time == nil {
		v := customer.UpdateDefaultUpdateTime()
		cuo.update_time = &v
	}
	if cuo.name != nil {
		if err := customer.NameValidator(*cuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if cuo.external_id != nil {
		if err := customer.ExternalIDValidator(*cuo.external_id); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	return cuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CustomerUpdateOne) SaveX(ctx context.Context) *Customer {
	c, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

// Exec executes the query on the entity.
func (cuo *CustomerUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CustomerUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *CustomerUpdateOne) sqlSave(ctx context.Context) (c *Customer, err error) {
	var (
		builder  = sql.Dialect(cuo.driver.Dialect())
		selector = builder.Select(customer.Columns...).From(builder.Table(customer.Table))
	)
	customer.ID(cuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		c = &Customer{config: cuo.config}
		if err := c.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Customer: %v", err)
		}
		id = c.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Customer with id: %v", cuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Customer with the same id: %v", cuo.id)
	}

	tx, err := cuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(customer.Table)
	)
	updater = updater.Where(sql.InInts(customer.FieldID, ids...))
	if value := cuo.update_time; value != nil {
		updater.Set(customer.FieldUpdateTime, *value)
		c.UpdateTime = *value
	}
	if value := cuo.name; value != nil {
		updater.Set(customer.FieldName, *value)
		c.Name = *value
	}
	if value := cuo.external_id; value != nil {
		updater.Set(customer.FieldExternalID, *value)
		c.ExternalID = value
	}
	if cuo.clearexternal_id {
		c.ExternalID = nil
		updater.SetNull(customer.FieldExternalID)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(cuo.removedServices) > 0 {
		eids := make([]int, len(cuo.removedServices))
		for eid := range cuo.removedServices {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(customer.ServicesTable).
			Where(sql.InInts(customer.ServicesPrimaryKey[1], ids...)).
			Where(sql.InInts(customer.ServicesPrimaryKey[0], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(cuo.services) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range cuo.services {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(customer.ServicesTable).
			Columns(customer.ServicesPrimaryKey[1], customer.ServicesPrimaryKey[0])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return c, nil
}
