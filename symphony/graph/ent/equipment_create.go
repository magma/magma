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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/property"
)

// EquipmentCreate is the builder for creating a Equipment entity.
type EquipmentCreate struct {
	config
	create_time     *time.Time
	update_time     *time.Time
	name            *string
	future_state    *string
	device_id       *string
	_type           map[string]struct{}
	location        map[string]struct{}
	parent_position map[string]struct{}
	positions       map[string]struct{}
	ports           map[string]struct{}
	work_order      map[string]struct{}
	properties      map[string]struct{}
	service         map[string]struct{}
	files           map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (ec *EquipmentCreate) SetCreateTime(t time.Time) *EquipmentCreate {
	ec.create_time = &t
	return ec
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableCreateTime(t *time.Time) *EquipmentCreate {
	if t != nil {
		ec.SetCreateTime(*t)
	}
	return ec
}

// SetUpdateTime sets the update_time field.
func (ec *EquipmentCreate) SetUpdateTime(t time.Time) *EquipmentCreate {
	ec.update_time = &t
	return ec
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableUpdateTime(t *time.Time) *EquipmentCreate {
	if t != nil {
		ec.SetUpdateTime(*t)
	}
	return ec
}

// SetName sets the name field.
func (ec *EquipmentCreate) SetName(s string) *EquipmentCreate {
	ec.name = &s
	return ec
}

// SetFutureState sets the future_state field.
func (ec *EquipmentCreate) SetFutureState(s string) *EquipmentCreate {
	ec.future_state = &s
	return ec
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableFutureState(s *string) *EquipmentCreate {
	if s != nil {
		ec.SetFutureState(*s)
	}
	return ec
}

// SetDeviceID sets the device_id field.
func (ec *EquipmentCreate) SetDeviceID(s string) *EquipmentCreate {
	ec.device_id = &s
	return ec
}

// SetNillableDeviceID sets the device_id field if the given value is not nil.
func (ec *EquipmentCreate) SetNillableDeviceID(s *string) *EquipmentCreate {
	if s != nil {
		ec.SetDeviceID(*s)
	}
	return ec
}

// SetTypeID sets the type edge to EquipmentType by id.
func (ec *EquipmentCreate) SetTypeID(id string) *EquipmentCreate {
	if ec._type == nil {
		ec._type = make(map[string]struct{})
	}
	ec._type[id] = struct{}{}
	return ec
}

// SetType sets the type edge to EquipmentType.
func (ec *EquipmentCreate) SetType(e *EquipmentType) *EquipmentCreate {
	return ec.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (ec *EquipmentCreate) SetLocationID(id string) *EquipmentCreate {
	if ec.location == nil {
		ec.location = make(map[string]struct{})
	}
	ec.location[id] = struct{}{}
	return ec
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableLocationID(id *string) *EquipmentCreate {
	if id != nil {
		ec = ec.SetLocationID(*id)
	}
	return ec
}

// SetLocation sets the location edge to Location.
func (ec *EquipmentCreate) SetLocation(l *Location) *EquipmentCreate {
	return ec.SetLocationID(l.ID)
}

// SetParentPositionID sets the parent_position edge to EquipmentPosition by id.
func (ec *EquipmentCreate) SetParentPositionID(id string) *EquipmentCreate {
	if ec.parent_position == nil {
		ec.parent_position = make(map[string]struct{})
	}
	ec.parent_position[id] = struct{}{}
	return ec
}

// SetNillableParentPositionID sets the parent_position edge to EquipmentPosition by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableParentPositionID(id *string) *EquipmentCreate {
	if id != nil {
		ec = ec.SetParentPositionID(*id)
	}
	return ec
}

// SetParentPosition sets the parent_position edge to EquipmentPosition.
func (ec *EquipmentCreate) SetParentPosition(e *EquipmentPosition) *EquipmentCreate {
	return ec.SetParentPositionID(e.ID)
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (ec *EquipmentCreate) AddPositionIDs(ids ...string) *EquipmentCreate {
	if ec.positions == nil {
		ec.positions = make(map[string]struct{})
	}
	for i := range ids {
		ec.positions[ids[i]] = struct{}{}
	}
	return ec
}

// AddPositions adds the positions edges to EquipmentPosition.
func (ec *EquipmentCreate) AddPositions(e ...*EquipmentPosition) *EquipmentCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddPositionIDs(ids...)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (ec *EquipmentCreate) AddPortIDs(ids ...string) *EquipmentCreate {
	if ec.ports == nil {
		ec.ports = make(map[string]struct{})
	}
	for i := range ids {
		ec.ports[ids[i]] = struct{}{}
	}
	return ec
}

// AddPorts adds the ports edges to EquipmentPort.
func (ec *EquipmentCreate) AddPorts(e ...*EquipmentPort) *EquipmentCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (ec *EquipmentCreate) SetWorkOrderID(id string) *EquipmentCreate {
	if ec.work_order == nil {
		ec.work_order = make(map[string]struct{})
	}
	ec.work_order[id] = struct{}{}
	return ec
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (ec *EquipmentCreate) SetNillableWorkOrderID(id *string) *EquipmentCreate {
	if id != nil {
		ec = ec.SetWorkOrderID(*id)
	}
	return ec
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (ec *EquipmentCreate) SetWorkOrder(w *WorkOrder) *EquipmentCreate {
	return ec.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (ec *EquipmentCreate) AddPropertyIDs(ids ...string) *EquipmentCreate {
	if ec.properties == nil {
		ec.properties = make(map[string]struct{})
	}
	for i := range ids {
		ec.properties[ids[i]] = struct{}{}
	}
	return ec
}

// AddProperties adds the properties edges to Property.
func (ec *EquipmentCreate) AddProperties(p ...*Property) *EquipmentCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ec.AddPropertyIDs(ids...)
}

// AddServiceIDs adds the service edge to Service by ids.
func (ec *EquipmentCreate) AddServiceIDs(ids ...string) *EquipmentCreate {
	if ec.service == nil {
		ec.service = make(map[string]struct{})
	}
	for i := range ids {
		ec.service[ids[i]] = struct{}{}
	}
	return ec
}

// AddService adds the service edges to Service.
func (ec *EquipmentCreate) AddService(s ...*Service) *EquipmentCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ec.AddServiceIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (ec *EquipmentCreate) AddFileIDs(ids ...string) *EquipmentCreate {
	if ec.files == nil {
		ec.files = make(map[string]struct{})
	}
	for i := range ids {
		ec.files[ids[i]] = struct{}{}
	}
	return ec
}

// AddFiles adds the files edges to File.
func (ec *EquipmentCreate) AddFiles(f ...*File) *EquipmentCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddFileIDs(ids...)
}

// Save creates the Equipment in the database.
func (ec *EquipmentCreate) Save(ctx context.Context) (*Equipment, error) {
	if ec.create_time == nil {
		v := equipment.DefaultCreateTime()
		ec.create_time = &v
	}
	if ec.update_time == nil {
		v := equipment.DefaultUpdateTime()
		ec.update_time = &v
	}
	if ec.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := equipment.NameValidator(*ec.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if len(ec._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if ec._type == nil {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	if len(ec.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(ec.parent_position) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent_position\"")
	}
	if len(ec.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return ec.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EquipmentCreate) SaveX(ctx context.Context) *Equipment {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ec *EquipmentCreate) sqlSave(ctx context.Context) (*Equipment, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ec.driver.Dialect())
		e       = &Equipment{config: ec.config}
	)
	tx, err := ec.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(equipment.Table).Default()
	if value := ec.create_time; value != nil {
		insert.Set(equipment.FieldCreateTime, *value)
		e.CreateTime = *value
	}
	if value := ec.update_time; value != nil {
		insert.Set(equipment.FieldUpdateTime, *value)
		e.UpdateTime = *value
	}
	if value := ec.name; value != nil {
		insert.Set(equipment.FieldName, *value)
		e.Name = *value
	}
	if value := ec.future_state; value != nil {
		insert.Set(equipment.FieldFutureState, *value)
		e.FutureState = *value
	}
	if value := ec.device_id; value != nil {
		insert.Set(equipment.FieldDeviceID, *value)
		e.DeviceID = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(equipment.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	e.ID = strconv.FormatInt(id, 10)
	if len(ec._type) > 0 {
		for eid := range ec._type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipment.TypeTable).
				Set(equipment.TypeColumn, eid).
				Where(sql.EQ(equipment.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(ec.location) > 0 {
		for eid := range ec.location {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipment.LocationTable).
				Set(equipment.LocationColumn, eid).
				Where(sql.EQ(equipment.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(ec.parent_position) > 0 {
		eid, err := strconv.Atoi(keys(ec.parent_position)[0])
		if err != nil {
			return nil, err
		}
		query, args := builder.Update(equipment.ParentPositionTable).
			Set(equipment.ParentPositionColumn, eid).
			Where(sql.EQ(equipment.FieldID, id).And().IsNull(equipment.ParentPositionColumn)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ec.parent_position) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"parent_position\" %v already connected to a different \"Equipment\"", keys(ec.parent_position))})
		}
	}
	if len(ec.positions) > 0 {
		p := sql.P()
		for eid := range ec.positions {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(equipmentposition.FieldID, eid)
		}
		query, args := builder.Update(equipment.PositionsTable).
			Set(equipment.PositionsColumn, id).
			Where(sql.And(p, sql.IsNull(equipment.PositionsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ec.positions) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"positions\" %v already connected to a different \"Equipment\"", keys(ec.positions))})
		}
	}
	if len(ec.ports) > 0 {
		p := sql.P()
		for eid := range ec.ports {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(equipmentport.FieldID, eid)
		}
		query, args := builder.Update(equipment.PortsTable).
			Set(equipment.PortsColumn, id).
			Where(sql.And(p, sql.IsNull(equipment.PortsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ec.ports) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"Equipment\"", keys(ec.ports))})
		}
	}
	if len(ec.work_order) > 0 {
		for eid := range ec.work_order {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipment.WorkOrderTable).
				Set(equipment.WorkOrderColumn, eid).
				Where(sql.EQ(equipment.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(ec.properties) > 0 {
		p := sql.P()
		for eid := range ec.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(property.FieldID, eid)
		}
		query, args := builder.Update(equipment.PropertiesTable).
			Set(equipment.PropertiesColumn, id).
			Where(sql.And(p, sql.IsNull(equipment.PropertiesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ec.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Equipment\"", keys(ec.properties))})
		}
	}
	if len(ec.service) > 0 {
		for eid := range ec.service {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(equipment.ServiceTable).
				Columns(equipment.ServicePrimaryKey[1], equipment.ServicePrimaryKey[0]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(ec.files) > 0 {
		p := sql.P()
		for eid := range ec.files {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(file.FieldID, eid)
		}
		query, args := builder.Update(equipment.FilesTable).
			Set(equipment.FilesColumn, id).
			Where(sql.And(p, sql.IsNull(equipment.FilesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ec.files) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"files\" %v already connected to a different \"Equipment\"", keys(ec.files))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return e, nil
}
