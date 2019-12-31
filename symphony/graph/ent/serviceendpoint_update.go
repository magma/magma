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
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// ServiceEndpointUpdate is the builder for updating ServiceEndpoint entities.
type ServiceEndpointUpdate struct {
	config

	update_time    *time.Time
	role           *string
	port           map[string]struct{}
	service        map[string]struct{}
	clearedPort    bool
	clearedService bool
	predicates     []predicate.ServiceEndpoint
}

// Where adds a new predicate for the builder.
func (seu *ServiceEndpointUpdate) Where(ps ...predicate.ServiceEndpoint) *ServiceEndpointUpdate {
	seu.predicates = append(seu.predicates, ps...)
	return seu
}

// SetRole sets the role field.
func (seu *ServiceEndpointUpdate) SetRole(s string) *ServiceEndpointUpdate {
	seu.role = &s
	return seu
}

// SetPortID sets the port edge to EquipmentPort by id.
func (seu *ServiceEndpointUpdate) SetPortID(id string) *ServiceEndpointUpdate {
	if seu.port == nil {
		seu.port = make(map[string]struct{})
	}
	seu.port[id] = struct{}{}
	return seu
}

// SetNillablePortID sets the port edge to EquipmentPort by id if the given value is not nil.
func (seu *ServiceEndpointUpdate) SetNillablePortID(id *string) *ServiceEndpointUpdate {
	if id != nil {
		seu = seu.SetPortID(*id)
	}
	return seu
}

// SetPort sets the port edge to EquipmentPort.
func (seu *ServiceEndpointUpdate) SetPort(e *EquipmentPort) *ServiceEndpointUpdate {
	return seu.SetPortID(e.ID)
}

// SetServiceID sets the service edge to Service by id.
func (seu *ServiceEndpointUpdate) SetServiceID(id string) *ServiceEndpointUpdate {
	if seu.service == nil {
		seu.service = make(map[string]struct{})
	}
	seu.service[id] = struct{}{}
	return seu
}

// SetNillableServiceID sets the service edge to Service by id if the given value is not nil.
func (seu *ServiceEndpointUpdate) SetNillableServiceID(id *string) *ServiceEndpointUpdate {
	if id != nil {
		seu = seu.SetServiceID(*id)
	}
	return seu
}

// SetService sets the service edge to Service.
func (seu *ServiceEndpointUpdate) SetService(s *Service) *ServiceEndpointUpdate {
	return seu.SetServiceID(s.ID)
}

// ClearPort clears the port edge to EquipmentPort.
func (seu *ServiceEndpointUpdate) ClearPort() *ServiceEndpointUpdate {
	seu.clearedPort = true
	return seu
}

// ClearService clears the service edge to Service.
func (seu *ServiceEndpointUpdate) ClearService() *ServiceEndpointUpdate {
	seu.clearedService = true
	return seu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (seu *ServiceEndpointUpdate) Save(ctx context.Context) (int, error) {
	if seu.update_time == nil {
		v := serviceendpoint.UpdateDefaultUpdateTime()
		seu.update_time = &v
	}
	if len(seu.port) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"port\"")
	}
	if len(seu.service) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"service\"")
	}
	return seu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (seu *ServiceEndpointUpdate) SaveX(ctx context.Context) int {
	affected, err := seu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (seu *ServiceEndpointUpdate) Exec(ctx context.Context) error {
	_, err := seu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (seu *ServiceEndpointUpdate) ExecX(ctx context.Context) {
	if err := seu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (seu *ServiceEndpointUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(seu.driver.Dialect())
		selector = builder.Select(serviceendpoint.FieldID).From(builder.Table(serviceendpoint.Table))
	)
	for _, p := range seu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = seu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := seu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(serviceendpoint.Table)
	)
	updater = updater.Where(sql.InInts(serviceendpoint.FieldID, ids...))
	if value := seu.update_time; value != nil {
		updater.Set(serviceendpoint.FieldUpdateTime, *value)
	}
	if value := seu.role; value != nil {
		updater.Set(serviceendpoint.FieldRole, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if seu.clearedPort {
		query, args := builder.Update(serviceendpoint.PortTable).
			SetNull(serviceendpoint.PortColumn).
			Where(sql.InInts(equipmentport.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(seu.port) > 0 {
		for eid := range seu.port {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(serviceendpoint.PortTable).
				Set(serviceendpoint.PortColumn, eid).
				Where(sql.InInts(serviceendpoint.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if seu.clearedService {
		query, args := builder.Update(serviceendpoint.ServiceTable).
			SetNull(serviceendpoint.ServiceColumn).
			Where(sql.InInts(service.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(seu.service) > 0 {
		for eid := range seu.service {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(serviceendpoint.ServiceTable).
				Set(serviceendpoint.ServiceColumn, eid).
				Where(sql.InInts(serviceendpoint.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// ServiceEndpointUpdateOne is the builder for updating a single ServiceEndpoint entity.
type ServiceEndpointUpdateOne struct {
	config
	id string

	update_time    *time.Time
	role           *string
	port           map[string]struct{}
	service        map[string]struct{}
	clearedPort    bool
	clearedService bool
}

// SetRole sets the role field.
func (seuo *ServiceEndpointUpdateOne) SetRole(s string) *ServiceEndpointUpdateOne {
	seuo.role = &s
	return seuo
}

// SetPortID sets the port edge to EquipmentPort by id.
func (seuo *ServiceEndpointUpdateOne) SetPortID(id string) *ServiceEndpointUpdateOne {
	if seuo.port == nil {
		seuo.port = make(map[string]struct{})
	}
	seuo.port[id] = struct{}{}
	return seuo
}

// SetNillablePortID sets the port edge to EquipmentPort by id if the given value is not nil.
func (seuo *ServiceEndpointUpdateOne) SetNillablePortID(id *string) *ServiceEndpointUpdateOne {
	if id != nil {
		seuo = seuo.SetPortID(*id)
	}
	return seuo
}

// SetPort sets the port edge to EquipmentPort.
func (seuo *ServiceEndpointUpdateOne) SetPort(e *EquipmentPort) *ServiceEndpointUpdateOne {
	return seuo.SetPortID(e.ID)
}

// SetServiceID sets the service edge to Service by id.
func (seuo *ServiceEndpointUpdateOne) SetServiceID(id string) *ServiceEndpointUpdateOne {
	if seuo.service == nil {
		seuo.service = make(map[string]struct{})
	}
	seuo.service[id] = struct{}{}
	return seuo
}

// SetNillableServiceID sets the service edge to Service by id if the given value is not nil.
func (seuo *ServiceEndpointUpdateOne) SetNillableServiceID(id *string) *ServiceEndpointUpdateOne {
	if id != nil {
		seuo = seuo.SetServiceID(*id)
	}
	return seuo
}

// SetService sets the service edge to Service.
func (seuo *ServiceEndpointUpdateOne) SetService(s *Service) *ServiceEndpointUpdateOne {
	return seuo.SetServiceID(s.ID)
}

// ClearPort clears the port edge to EquipmentPort.
func (seuo *ServiceEndpointUpdateOne) ClearPort() *ServiceEndpointUpdateOne {
	seuo.clearedPort = true
	return seuo
}

// ClearService clears the service edge to Service.
func (seuo *ServiceEndpointUpdateOne) ClearService() *ServiceEndpointUpdateOne {
	seuo.clearedService = true
	return seuo
}

// Save executes the query and returns the updated entity.
func (seuo *ServiceEndpointUpdateOne) Save(ctx context.Context) (*ServiceEndpoint, error) {
	if seuo.update_time == nil {
		v := serviceendpoint.UpdateDefaultUpdateTime()
		seuo.update_time = &v
	}
	if len(seuo.port) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"port\"")
	}
	if len(seuo.service) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"service\"")
	}
	return seuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (seuo *ServiceEndpointUpdateOne) SaveX(ctx context.Context) *ServiceEndpoint {
	se, err := seuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return se
}

// Exec executes the query on the entity.
func (seuo *ServiceEndpointUpdateOne) Exec(ctx context.Context) error {
	_, err := seuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (seuo *ServiceEndpointUpdateOne) ExecX(ctx context.Context) {
	if err := seuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (seuo *ServiceEndpointUpdateOne) sqlSave(ctx context.Context) (se *ServiceEndpoint, err error) {
	var (
		builder  = sql.Dialect(seuo.driver.Dialect())
		selector = builder.Select(serviceendpoint.Columns...).From(builder.Table(serviceendpoint.Table))
	)
	serviceendpoint.ID(seuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = seuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		se = &ServiceEndpoint{config: seuo.config}
		if err := se.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into ServiceEndpoint: %v", err)
		}
		id = se.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("ServiceEndpoint with id: %v", seuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one ServiceEndpoint with the same id: %v", seuo.id)
	}

	tx, err := seuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(serviceendpoint.Table)
	)
	updater = updater.Where(sql.InInts(serviceendpoint.FieldID, ids...))
	if value := seuo.update_time; value != nil {
		updater.Set(serviceendpoint.FieldUpdateTime, *value)
		se.UpdateTime = *value
	}
	if value := seuo.role; value != nil {
		updater.Set(serviceendpoint.FieldRole, *value)
		se.Role = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if seuo.clearedPort {
		query, args := builder.Update(serviceendpoint.PortTable).
			SetNull(serviceendpoint.PortColumn).
			Where(sql.InInts(equipmentport.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(seuo.port) > 0 {
		for eid := range seuo.port {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(serviceendpoint.PortTable).
				Set(serviceendpoint.PortColumn, eid).
				Where(sql.InInts(serviceendpoint.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if seuo.clearedService {
		query, args := builder.Update(serviceendpoint.ServiceTable).
			SetNull(serviceendpoint.ServiceColumn).
			Where(sql.InInts(service.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(seuo.service) > 0 {
		for eid := range seuo.service {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(serviceendpoint.ServiceTable).
				Set(serviceendpoint.ServiceColumn, eid).
				Where(sql.InInts(serviceendpoint.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return se, nil
}
