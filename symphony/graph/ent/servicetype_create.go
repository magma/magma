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
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceTypeCreate is the builder for creating a ServiceType entity.
type ServiceTypeCreate struct {
	config
	create_time    *time.Time
	update_time    *time.Time
	name           *string
	has_customer   *bool
	services       map[string]struct{}
	property_types map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (stc *ServiceTypeCreate) SetCreateTime(t time.Time) *ServiceTypeCreate {
	stc.create_time = &t
	return stc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableCreateTime(t *time.Time) *ServiceTypeCreate {
	if t != nil {
		stc.SetCreateTime(*t)
	}
	return stc
}

// SetUpdateTime sets the update_time field.
func (stc *ServiceTypeCreate) SetUpdateTime(t time.Time) *ServiceTypeCreate {
	stc.update_time = &t
	return stc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableUpdateTime(t *time.Time) *ServiceTypeCreate {
	if t != nil {
		stc.SetUpdateTime(*t)
	}
	return stc
}

// SetName sets the name field.
func (stc *ServiceTypeCreate) SetName(s string) *ServiceTypeCreate {
	stc.name = &s
	return stc
}

// SetHasCustomer sets the has_customer field.
func (stc *ServiceTypeCreate) SetHasCustomer(b bool) *ServiceTypeCreate {
	stc.has_customer = &b
	return stc
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stc *ServiceTypeCreate) SetNillableHasCustomer(b *bool) *ServiceTypeCreate {
	if b != nil {
		stc.SetHasCustomer(*b)
	}
	return stc
}

// AddServiceIDs adds the services edge to Service by ids.
func (stc *ServiceTypeCreate) AddServiceIDs(ids ...string) *ServiceTypeCreate {
	if stc.services == nil {
		stc.services = make(map[string]struct{})
	}
	for i := range ids {
		stc.services[ids[i]] = struct{}{}
	}
	return stc
}

// AddServices adds the services edges to Service.
func (stc *ServiceTypeCreate) AddServices(s ...*Service) *ServiceTypeCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stc.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stc *ServiceTypeCreate) AddPropertyTypeIDs(ids ...string) *ServiceTypeCreate {
	if stc.property_types == nil {
		stc.property_types = make(map[string]struct{})
	}
	for i := range ids {
		stc.property_types[ids[i]] = struct{}{}
	}
	return stc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stc *ServiceTypeCreate) AddPropertyTypes(p ...*PropertyType) *ServiceTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stc.AddPropertyTypeIDs(ids...)
}

// Save creates the ServiceType in the database.
func (stc *ServiceTypeCreate) Save(ctx context.Context) (*ServiceType, error) {
	if stc.create_time == nil {
		v := servicetype.DefaultCreateTime()
		stc.create_time = &v
	}
	if stc.update_time == nil {
		v := servicetype.DefaultUpdateTime()
		stc.update_time = &v
	}
	if stc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if stc.has_customer == nil {
		v := servicetype.DefaultHasCustomer
		stc.has_customer = &v
	}
	return stc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (stc *ServiceTypeCreate) SaveX(ctx context.Context) *ServiceType {
	v, err := stc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stc *ServiceTypeCreate) sqlSave(ctx context.Context) (*ServiceType, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(stc.driver.Dialect())
		st      = &ServiceType{config: stc.config}
	)
	tx, err := stc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(servicetype.Table).Default()
	if value := stc.create_time; value != nil {
		insert.Set(servicetype.FieldCreateTime, *value)
		st.CreateTime = *value
	}
	if value := stc.update_time; value != nil {
		insert.Set(servicetype.FieldUpdateTime, *value)
		st.UpdateTime = *value
	}
	if value := stc.name; value != nil {
		insert.Set(servicetype.FieldName, *value)
		st.Name = *value
	}
	if value := stc.has_customer; value != nil {
		insert.Set(servicetype.FieldHasCustomer, *value)
		st.HasCustomer = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(servicetype.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	st.ID = strconv.FormatInt(id, 10)
	if len(stc.services) > 0 {
		p := sql.P()
		for eid := range stc.services {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(stc.services) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"services\" %v already connected to a different \"ServiceType\"", keys(stc.services))})
		}
	}
	if len(stc.property_types) > 0 {
		p := sql.P()
		for eid := range stc.property_types {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(stc.property_types) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"ServiceType\"", keys(stc.property_types))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return st, nil
}
