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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceTypeUpdate is the builder for updating ServiceType entities.
type ServiceTypeUpdate struct {
	config

	update_time          *time.Time
	name                 *string
	has_customer         *bool
	services             map[string]struct{}
	property_types       map[string]struct{}
	removedServices      map[string]struct{}
	removedPropertyTypes map[string]struct{}
	predicates           []predicate.ServiceType
}

// Where adds a new predicate for the builder.
func (stu *ServiceTypeUpdate) Where(ps ...predicate.ServiceType) *ServiceTypeUpdate {
	stu.predicates = append(stu.predicates, ps...)
	return stu
}

// SetName sets the name field.
func (stu *ServiceTypeUpdate) SetName(s string) *ServiceTypeUpdate {
	stu.name = &s
	return stu
}

// SetHasCustomer sets the has_customer field.
func (stu *ServiceTypeUpdate) SetHasCustomer(b bool) *ServiceTypeUpdate {
	stu.has_customer = &b
	return stu
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stu *ServiceTypeUpdate) SetNillableHasCustomer(b *bool) *ServiceTypeUpdate {
	if b != nil {
		stu.SetHasCustomer(*b)
	}
	return stu
}

// AddServiceIDs adds the services edge to Service by ids.
func (stu *ServiceTypeUpdate) AddServiceIDs(ids ...string) *ServiceTypeUpdate {
	if stu.services == nil {
		stu.services = make(map[string]struct{})
	}
	for i := range ids {
		stu.services[ids[i]] = struct{}{}
	}
	return stu
}

// AddServices adds the services edges to Service.
func (stu *ServiceTypeUpdate) AddServices(s ...*Service) *ServiceTypeUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stu *ServiceTypeUpdate) AddPropertyTypeIDs(ids ...string) *ServiceTypeUpdate {
	if stu.property_types == nil {
		stu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		stu.property_types[ids[i]] = struct{}{}
	}
	return stu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stu *ServiceTypeUpdate) AddPropertyTypes(p ...*PropertyType) *ServiceTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stu.AddPropertyTypeIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (stu *ServiceTypeUpdate) RemoveServiceIDs(ids ...string) *ServiceTypeUpdate {
	if stu.removedServices == nil {
		stu.removedServices = make(map[string]struct{})
	}
	for i := range ids {
		stu.removedServices[ids[i]] = struct{}{}
	}
	return stu
}

// RemoveServices removes services edges to Service.
func (stu *ServiceTypeUpdate) RemoveServices(s ...*Service) *ServiceTypeUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.RemoveServiceIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (stu *ServiceTypeUpdate) RemovePropertyTypeIDs(ids ...string) *ServiceTypeUpdate {
	if stu.removedPropertyTypes == nil {
		stu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		stu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return stu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (stu *ServiceTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *ServiceTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stu.RemovePropertyTypeIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stu *ServiceTypeUpdate) Save(ctx context.Context) (int, error) {
	if stu.update_time == nil {
		v := servicetype.UpdateDefaultUpdateTime()
		stu.update_time = &v
	}
	return stu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stu *ServiceTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := stu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (stu *ServiceTypeUpdate) Exec(ctx context.Context) error {
	_, err := stu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stu *ServiceTypeUpdate) ExecX(ctx context.Context) {
	if err := stu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stu *ServiceTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(stu.driver.Dialect())
		selector = builder.Select(servicetype.FieldID).From(builder.Table(servicetype.Table))
	)
	for _, p := range stu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = stu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := stu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(servicetype.Table)
	)
	updater = updater.Where(sql.InInts(servicetype.FieldID, ids...))
	if value := stu.update_time; value != nil {
		updater.Set(servicetype.FieldUpdateTime, *value)
	}
	if value := stu.name; value != nil {
		updater.Set(servicetype.FieldName, *value)
	}
	if value := stu.has_customer; value != nil {
		updater.Set(servicetype.FieldHasCustomer, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(stu.removedServices) > 0 {
		eids := make([]int, len(stu.removedServices))
		for eid := range stu.removedServices {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(servicetype.ServicesTable).
			SetNull(servicetype.ServicesColumn).
			Where(sql.InInts(servicetype.ServicesColumn, ids...)).
			Where(sql.InInts(service.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(stu.services) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range stu.services {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(service.FieldID, eid)
			}
			query, args := builder.Update(servicetype.ServicesTable).
				Set(servicetype.ServicesColumn, id).
				Where(sql.And(p, sql.IsNull(servicetype.ServicesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(stu.services) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"services\" %v already connected to a different \"ServiceType\"", keys(stu.services))})
			}
		}
	}
	if len(stu.removedPropertyTypes) > 0 {
		eids := make([]int, len(stu.removedPropertyTypes))
		for eid := range stu.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(servicetype.PropertyTypesTable).
			SetNull(servicetype.PropertyTypesColumn).
			Where(sql.InInts(servicetype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(stu.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range stu.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(servicetype.PropertyTypesTable).
				Set(servicetype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(servicetype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(stu.property_types) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"ServiceType\"", keys(stu.property_types))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// ServiceTypeUpdateOne is the builder for updating a single ServiceType entity.
type ServiceTypeUpdateOne struct {
	config
	id string

	update_time          *time.Time
	name                 *string
	has_customer         *bool
	services             map[string]struct{}
	property_types       map[string]struct{}
	removedServices      map[string]struct{}
	removedPropertyTypes map[string]struct{}
}

// SetName sets the name field.
func (stuo *ServiceTypeUpdateOne) SetName(s string) *ServiceTypeUpdateOne {
	stuo.name = &s
	return stuo
}

// SetHasCustomer sets the has_customer field.
func (stuo *ServiceTypeUpdateOne) SetHasCustomer(b bool) *ServiceTypeUpdateOne {
	stuo.has_customer = &b
	return stuo
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stuo *ServiceTypeUpdateOne) SetNillableHasCustomer(b *bool) *ServiceTypeUpdateOne {
	if b != nil {
		stuo.SetHasCustomer(*b)
	}
	return stuo
}

// AddServiceIDs adds the services edge to Service by ids.
func (stuo *ServiceTypeUpdateOne) AddServiceIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.services == nil {
		stuo.services = make(map[string]struct{})
	}
	for i := range ids {
		stuo.services[ids[i]] = struct{}{}
	}
	return stuo
}

// AddServices adds the services edges to Service.
func (stuo *ServiceTypeUpdateOne) AddServices(s ...*Service) *ServiceTypeUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stuo *ServiceTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.property_types == nil {
		stuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		stuo.property_types[ids[i]] = struct{}{}
	}
	return stuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stuo *ServiceTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *ServiceTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stuo.AddPropertyTypeIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (stuo *ServiceTypeUpdateOne) RemoveServiceIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.removedServices == nil {
		stuo.removedServices = make(map[string]struct{})
	}
	for i := range ids {
		stuo.removedServices[ids[i]] = struct{}{}
	}
	return stuo
}

// RemoveServices removes services edges to Service.
func (stuo *ServiceTypeUpdateOne) RemoveServices(s ...*Service) *ServiceTypeUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.RemoveServiceIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (stuo *ServiceTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.removedPropertyTypes == nil {
		stuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		stuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return stuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (stuo *ServiceTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *ServiceTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stuo.RemovePropertyTypeIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (stuo *ServiceTypeUpdateOne) Save(ctx context.Context) (*ServiceType, error) {
	if stuo.update_time == nil {
		v := servicetype.UpdateDefaultUpdateTime()
		stuo.update_time = &v
	}
	return stuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stuo *ServiceTypeUpdateOne) SaveX(ctx context.Context) *ServiceType {
	st, err := stuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return st
}

// Exec executes the query on the entity.
func (stuo *ServiceTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := stuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stuo *ServiceTypeUpdateOne) ExecX(ctx context.Context) {
	if err := stuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stuo *ServiceTypeUpdateOne) sqlSave(ctx context.Context) (st *ServiceType, err error) {
	var (
		builder  = sql.Dialect(stuo.driver.Dialect())
		selector = builder.Select(servicetype.Columns...).From(builder.Table(servicetype.Table))
	)
	servicetype.ID(stuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = stuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		st = &ServiceType{config: stuo.config}
		if err := st.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into ServiceType: %v", err)
		}
		id = st.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("ServiceType with id: %v", stuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one ServiceType with the same id: %v", stuo.id)
	}

	tx, err := stuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(servicetype.Table)
	)
	updater = updater.Where(sql.InInts(servicetype.FieldID, ids...))
	if value := stuo.update_time; value != nil {
		updater.Set(servicetype.FieldUpdateTime, *value)
		st.UpdateTime = *value
	}
	if value := stuo.name; value != nil {
		updater.Set(servicetype.FieldName, *value)
		st.Name = *value
	}
	if value := stuo.has_customer; value != nil {
		updater.Set(servicetype.FieldHasCustomer, *value)
		st.HasCustomer = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(stuo.removedServices) > 0 {
		eids := make([]int, len(stuo.removedServices))
		for eid := range stuo.removedServices {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(servicetype.ServicesTable).
			SetNull(servicetype.ServicesColumn).
			Where(sql.InInts(servicetype.ServicesColumn, ids...)).
			Where(sql.InInts(service.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(stuo.services) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range stuo.services {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(service.FieldID, eid)
			}
			query, args := builder.Update(servicetype.ServicesTable).
				Set(servicetype.ServicesColumn, id).
				Where(sql.And(p, sql.IsNull(servicetype.ServicesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(stuo.services) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"services\" %v already connected to a different \"ServiceType\"", keys(stuo.services))})
			}
		}
	}
	if len(stuo.removedPropertyTypes) > 0 {
		eids := make([]int, len(stuo.removedPropertyTypes))
		for eid := range stuo.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(servicetype.PropertyTypesTable).
			SetNull(servicetype.PropertyTypesColumn).
			Where(sql.InInts(servicetype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(stuo.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range stuo.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(servicetype.PropertyTypesTable).
				Set(servicetype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(servicetype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(stuo.property_types) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"ServiceType\"", keys(stuo.property_types))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return st, nil
}
