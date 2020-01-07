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
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// EquipmentUpdate is the builder for updating Equipment entities.
type EquipmentUpdate struct {
	config

	update_time           *time.Time
	name                  *string
	future_state          *string
	clearfuture_state     bool
	device_id             *string
	cleardevice_id        bool
	external_id           *string
	clearexternal_id      bool
	_type                 map[string]struct{}
	location              map[string]struct{}
	parent_position       map[string]struct{}
	positions             map[string]struct{}
	ports                 map[string]struct{}
	work_order            map[string]struct{}
	properties            map[string]struct{}
	files                 map[string]struct{}
	clearedType           bool
	clearedLocation       bool
	clearedParentPosition bool
	removedPositions      map[string]struct{}
	removedPorts          map[string]struct{}
	clearedWorkOrder      bool
	removedProperties     map[string]struct{}
	removedFiles          map[string]struct{}
	predicates            []predicate.Equipment
}

// Where adds a new predicate for the builder.
func (eu *EquipmentUpdate) Where(ps ...predicate.Equipment) *EquipmentUpdate {
	eu.predicates = append(eu.predicates, ps...)
	return eu
}

// SetName sets the name field.
func (eu *EquipmentUpdate) SetName(s string) *EquipmentUpdate {
	eu.name = &s
	return eu
}

// SetFutureState sets the future_state field.
func (eu *EquipmentUpdate) SetFutureState(s string) *EquipmentUpdate {
	eu.future_state = &s
	return eu
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableFutureState(s *string) *EquipmentUpdate {
	if s != nil {
		eu.SetFutureState(*s)
	}
	return eu
}

// ClearFutureState clears the value of future_state.
func (eu *EquipmentUpdate) ClearFutureState() *EquipmentUpdate {
	eu.future_state = nil
	eu.clearfuture_state = true
	return eu
}

// SetDeviceID sets the device_id field.
func (eu *EquipmentUpdate) SetDeviceID(s string) *EquipmentUpdate {
	eu.device_id = &s
	return eu
}

// SetNillableDeviceID sets the device_id field if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableDeviceID(s *string) *EquipmentUpdate {
	if s != nil {
		eu.SetDeviceID(*s)
	}
	return eu
}

// ClearDeviceID clears the value of device_id.
func (eu *EquipmentUpdate) ClearDeviceID() *EquipmentUpdate {
	eu.device_id = nil
	eu.cleardevice_id = true
	return eu
}

// SetExternalID sets the external_id field.
func (eu *EquipmentUpdate) SetExternalID(s string) *EquipmentUpdate {
	eu.external_id = &s
	return eu
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableExternalID(s *string) *EquipmentUpdate {
	if s != nil {
		eu.SetExternalID(*s)
	}
	return eu
}

// ClearExternalID clears the value of external_id.
func (eu *EquipmentUpdate) ClearExternalID() *EquipmentUpdate {
	eu.external_id = nil
	eu.clearexternal_id = true
	return eu
}

// SetTypeID sets the type edge to EquipmentType by id.
func (eu *EquipmentUpdate) SetTypeID(id string) *EquipmentUpdate {
	if eu._type == nil {
		eu._type = make(map[string]struct{})
	}
	eu._type[id] = struct{}{}
	return eu
}

// SetType sets the type edge to EquipmentType.
func (eu *EquipmentUpdate) SetType(e *EquipmentType) *EquipmentUpdate {
	return eu.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (eu *EquipmentUpdate) SetLocationID(id string) *EquipmentUpdate {
	if eu.location == nil {
		eu.location = make(map[string]struct{})
	}
	eu.location[id] = struct{}{}
	return eu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableLocationID(id *string) *EquipmentUpdate {
	if id != nil {
		eu = eu.SetLocationID(*id)
	}
	return eu
}

// SetLocation sets the location edge to Location.
func (eu *EquipmentUpdate) SetLocation(l *Location) *EquipmentUpdate {
	return eu.SetLocationID(l.ID)
}

// SetParentPositionID sets the parent_position edge to EquipmentPosition by id.
func (eu *EquipmentUpdate) SetParentPositionID(id string) *EquipmentUpdate {
	if eu.parent_position == nil {
		eu.parent_position = make(map[string]struct{})
	}
	eu.parent_position[id] = struct{}{}
	return eu
}

// SetNillableParentPositionID sets the parent_position edge to EquipmentPosition by id if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableParentPositionID(id *string) *EquipmentUpdate {
	if id != nil {
		eu = eu.SetParentPositionID(*id)
	}
	return eu
}

// SetParentPosition sets the parent_position edge to EquipmentPosition.
func (eu *EquipmentUpdate) SetParentPosition(e *EquipmentPosition) *EquipmentUpdate {
	return eu.SetParentPositionID(e.ID)
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (eu *EquipmentUpdate) AddPositionIDs(ids ...string) *EquipmentUpdate {
	if eu.positions == nil {
		eu.positions = make(map[string]struct{})
	}
	for i := range ids {
		eu.positions[ids[i]] = struct{}{}
	}
	return eu
}

// AddPositions adds the positions edges to EquipmentPosition.
func (eu *EquipmentUpdate) AddPositions(e ...*EquipmentPosition) *EquipmentUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.AddPositionIDs(ids...)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (eu *EquipmentUpdate) AddPortIDs(ids ...string) *EquipmentUpdate {
	if eu.ports == nil {
		eu.ports = make(map[string]struct{})
	}
	for i := range ids {
		eu.ports[ids[i]] = struct{}{}
	}
	return eu
}

// AddPorts adds the ports edges to EquipmentPort.
func (eu *EquipmentUpdate) AddPorts(e ...*EquipmentPort) *EquipmentUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (eu *EquipmentUpdate) SetWorkOrderID(id string) *EquipmentUpdate {
	if eu.work_order == nil {
		eu.work_order = make(map[string]struct{})
	}
	eu.work_order[id] = struct{}{}
	return eu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (eu *EquipmentUpdate) SetNillableWorkOrderID(id *string) *EquipmentUpdate {
	if id != nil {
		eu = eu.SetWorkOrderID(*id)
	}
	return eu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (eu *EquipmentUpdate) SetWorkOrder(w *WorkOrder) *EquipmentUpdate {
	return eu.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (eu *EquipmentUpdate) AddPropertyIDs(ids ...string) *EquipmentUpdate {
	if eu.properties == nil {
		eu.properties = make(map[string]struct{})
	}
	for i := range ids {
		eu.properties[ids[i]] = struct{}{}
	}
	return eu
}

// AddProperties adds the properties edges to Property.
func (eu *EquipmentUpdate) AddProperties(p ...*Property) *EquipmentUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.AddPropertyIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (eu *EquipmentUpdate) AddFileIDs(ids ...string) *EquipmentUpdate {
	if eu.files == nil {
		eu.files = make(map[string]struct{})
	}
	for i := range ids {
		eu.files[ids[i]] = struct{}{}
	}
	return eu
}

// AddFiles adds the files edges to File.
func (eu *EquipmentUpdate) AddFiles(f ...*File) *EquipmentUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return eu.AddFileIDs(ids...)
}

// ClearType clears the type edge to EquipmentType.
func (eu *EquipmentUpdate) ClearType() *EquipmentUpdate {
	eu.clearedType = true
	return eu
}

// ClearLocation clears the location edge to Location.
func (eu *EquipmentUpdate) ClearLocation() *EquipmentUpdate {
	eu.clearedLocation = true
	return eu
}

// ClearParentPosition clears the parent_position edge to EquipmentPosition.
func (eu *EquipmentUpdate) ClearParentPosition() *EquipmentUpdate {
	eu.clearedParentPosition = true
	return eu
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (eu *EquipmentUpdate) RemovePositionIDs(ids ...string) *EquipmentUpdate {
	if eu.removedPositions == nil {
		eu.removedPositions = make(map[string]struct{})
	}
	for i := range ids {
		eu.removedPositions[ids[i]] = struct{}{}
	}
	return eu
}

// RemovePositions removes positions edges to EquipmentPosition.
func (eu *EquipmentUpdate) RemovePositions(e ...*EquipmentPosition) *EquipmentUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.RemovePositionIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (eu *EquipmentUpdate) RemovePortIDs(ids ...string) *EquipmentUpdate {
	if eu.removedPorts == nil {
		eu.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		eu.removedPorts[ids[i]] = struct{}{}
	}
	return eu
}

// RemovePorts removes ports edges to EquipmentPort.
func (eu *EquipmentUpdate) RemovePorts(e ...*EquipmentPort) *EquipmentUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eu.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (eu *EquipmentUpdate) ClearWorkOrder() *EquipmentUpdate {
	eu.clearedWorkOrder = true
	return eu
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (eu *EquipmentUpdate) RemovePropertyIDs(ids ...string) *EquipmentUpdate {
	if eu.removedProperties == nil {
		eu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		eu.removedProperties[ids[i]] = struct{}{}
	}
	return eu
}

// RemoveProperties removes properties edges to Property.
func (eu *EquipmentUpdate) RemoveProperties(p ...*Property) *EquipmentUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.RemovePropertyIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (eu *EquipmentUpdate) RemoveFileIDs(ids ...string) *EquipmentUpdate {
	if eu.removedFiles == nil {
		eu.removedFiles = make(map[string]struct{})
	}
	for i := range ids {
		eu.removedFiles[ids[i]] = struct{}{}
	}
	return eu
}

// RemoveFiles removes files edges to File.
func (eu *EquipmentUpdate) RemoveFiles(f ...*File) *EquipmentUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return eu.RemoveFileIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (eu *EquipmentUpdate) Save(ctx context.Context) (int, error) {
	if eu.update_time == nil {
		v := equipment.UpdateDefaultUpdateTime()
		eu.update_time = &v
	}
	if eu.name != nil {
		if err := equipment.NameValidator(*eu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(eu._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if eu.clearedType && eu._type == nil {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(eu.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(eu.parent_position) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"parent_position\"")
	}
	if len(eu.work_order) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return eu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EquipmentUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EquipmentUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EquipmentUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eu *EquipmentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(eu.driver.Dialect())
		selector = builder.Select(equipment.FieldID).From(builder.Table(equipment.Table))
	)
	for _, p := range eu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = eu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := eu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipment.Table)
	)
	updater = updater.Where(sql.InInts(equipment.FieldID, ids...))
	if value := eu.update_time; value != nil {
		updater.Set(equipment.FieldUpdateTime, *value)
	}
	if value := eu.name; value != nil {
		updater.Set(equipment.FieldName, *value)
	}
	if value := eu.future_state; value != nil {
		updater.Set(equipment.FieldFutureState, *value)
	}
	if eu.clearfuture_state {
		updater.SetNull(equipment.FieldFutureState)
	}
	if value := eu.device_id; value != nil {
		updater.Set(equipment.FieldDeviceID, *value)
	}
	if eu.cleardevice_id {
		updater.SetNull(equipment.FieldDeviceID)
	}
	if value := eu.external_id; value != nil {
		updater.Set(equipment.FieldExternalID, *value)
	}
	if eu.clearexternal_id {
		updater.SetNull(equipment.FieldExternalID)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if eu.clearedType {
		query, args := builder.Update(equipment.TypeTable).
			SetNull(equipment.TypeColumn).
			Where(sql.InInts(equipmenttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu._type) > 0 {
		for eid := range eu._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipment.TypeTable).
				Set(equipment.TypeColumn, eid).
				Where(sql.InInts(equipment.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if eu.clearedLocation {
		query, args := builder.Update(equipment.LocationTable).
			SetNull(equipment.LocationColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu.location) > 0 {
		for eid := range eu.location {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipment.LocationTable).
				Set(equipment.LocationColumn, eid).
				Where(sql.InInts(equipment.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if eu.clearedParentPosition {
		query, args := builder.Update(equipment.ParentPositionTable).
			SetNull(equipment.ParentPositionColumn).
			Where(sql.InInts(equipmentposition.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu.parent_position) > 0 {
		for _, id := range ids {
			eid, serr := strconv.Atoi(keys(eu.parent_position)[0])
			if serr != nil {
				return 0, rollback(tx, err)
			}
			query, args := builder.Update(equipment.ParentPositionTable).
				Set(equipment.ParentPositionColumn, eid).
				Where(sql.EQ(equipment.FieldID, id).And().IsNull(equipment.ParentPositionColumn)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eu.parent_position) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"parent_position\" %v already connected to a different \"Equipment\"", keys(eu.parent_position))})
			}
		}
	}
	if len(eu.removedPositions) > 0 {
		eids := make([]int, len(eu.removedPositions))
		for eid := range eu.removedPositions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.PositionsTable).
			SetNull(equipment.PositionsColumn).
			Where(sql.InInts(equipment.PositionsColumn, ids...)).
			Where(sql.InInts(equipmentposition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu.positions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eu.positions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentposition.FieldID, eid)
			}
			query, args := builder.Update(equipment.PositionsTable).
				Set(equipment.PositionsColumn, id).
				Where(sql.And(p, sql.IsNull(equipment.PositionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eu.positions) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"positions\" %v already connected to a different \"Equipment\"", keys(eu.positions))})
			}
		}
	}
	if len(eu.removedPorts) > 0 {
		eids := make([]int, len(eu.removedPorts))
		for eid := range eu.removedPorts {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.PortsTable).
			SetNull(equipment.PortsColumn).
			Where(sql.InInts(equipment.PortsColumn, ids...)).
			Where(sql.InInts(equipmentport.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu.ports) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eu.ports {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentport.FieldID, eid)
			}
			query, args := builder.Update(equipment.PortsTable).
				Set(equipment.PortsColumn, id).
				Where(sql.And(p, sql.IsNull(equipment.PortsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eu.ports) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"Equipment\"", keys(eu.ports))})
			}
		}
	}
	if eu.clearedWorkOrder {
		query, args := builder.Update(equipment.WorkOrderTable).
			SetNull(equipment.WorkOrderColumn).
			Where(sql.InInts(workorder.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu.work_order) > 0 {
		for eid := range eu.work_order {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipment.WorkOrderTable).
				Set(equipment.WorkOrderColumn, eid).
				Where(sql.InInts(equipment.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if len(eu.removedProperties) > 0 {
		eids := make([]int, len(eu.removedProperties))
		for eid := range eu.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.PropertiesTable).
			SetNull(equipment.PropertiesColumn).
			Where(sql.InInts(equipment.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eu.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(equipment.PropertiesTable).
				Set(equipment.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(equipment.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eu.properties) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Equipment\"", keys(eu.properties))})
			}
		}
	}
	if len(eu.removedFiles) > 0 {
		eids := make([]int, len(eu.removedFiles))
		for eid := range eu.removedFiles {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.FilesTable).
			SetNull(equipment.FilesColumn).
			Where(sql.InInts(equipment.FilesColumn, ids...)).
			Where(sql.InInts(file.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eu.files) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eu.files {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(file.FieldID, eid)
			}
			query, args := builder.Update(equipment.FilesTable).
				Set(equipment.FilesColumn, id).
				Where(sql.And(p, sql.IsNull(equipment.FilesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eu.files) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"files\" %v already connected to a different \"Equipment\"", keys(eu.files))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// EquipmentUpdateOne is the builder for updating a single Equipment entity.
type EquipmentUpdateOne struct {
	config
	id string

	update_time           *time.Time
	name                  *string
	future_state          *string
	clearfuture_state     bool
	device_id             *string
	cleardevice_id        bool
	external_id           *string
	clearexternal_id      bool
	_type                 map[string]struct{}
	location              map[string]struct{}
	parent_position       map[string]struct{}
	positions             map[string]struct{}
	ports                 map[string]struct{}
	work_order            map[string]struct{}
	properties            map[string]struct{}
	files                 map[string]struct{}
	clearedType           bool
	clearedLocation       bool
	clearedParentPosition bool
	removedPositions      map[string]struct{}
	removedPorts          map[string]struct{}
	clearedWorkOrder      bool
	removedProperties     map[string]struct{}
	removedFiles          map[string]struct{}
}

// SetName sets the name field.
func (euo *EquipmentUpdateOne) SetName(s string) *EquipmentUpdateOne {
	euo.name = &s
	return euo
}

// SetFutureState sets the future_state field.
func (euo *EquipmentUpdateOne) SetFutureState(s string) *EquipmentUpdateOne {
	euo.future_state = &s
	return euo
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableFutureState(s *string) *EquipmentUpdateOne {
	if s != nil {
		euo.SetFutureState(*s)
	}
	return euo
}

// ClearFutureState clears the value of future_state.
func (euo *EquipmentUpdateOne) ClearFutureState() *EquipmentUpdateOne {
	euo.future_state = nil
	euo.clearfuture_state = true
	return euo
}

// SetDeviceID sets the device_id field.
func (euo *EquipmentUpdateOne) SetDeviceID(s string) *EquipmentUpdateOne {
	euo.device_id = &s
	return euo
}

// SetNillableDeviceID sets the device_id field if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableDeviceID(s *string) *EquipmentUpdateOne {
	if s != nil {
		euo.SetDeviceID(*s)
	}
	return euo
}

// ClearDeviceID clears the value of device_id.
func (euo *EquipmentUpdateOne) ClearDeviceID() *EquipmentUpdateOne {
	euo.device_id = nil
	euo.cleardevice_id = true
	return euo
}

// SetExternalID sets the external_id field.
func (euo *EquipmentUpdateOne) SetExternalID(s string) *EquipmentUpdateOne {
	euo.external_id = &s
	return euo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableExternalID(s *string) *EquipmentUpdateOne {
	if s != nil {
		euo.SetExternalID(*s)
	}
	return euo
}

// ClearExternalID clears the value of external_id.
func (euo *EquipmentUpdateOne) ClearExternalID() *EquipmentUpdateOne {
	euo.external_id = nil
	euo.clearexternal_id = true
	return euo
}

// SetTypeID sets the type edge to EquipmentType by id.
func (euo *EquipmentUpdateOne) SetTypeID(id string) *EquipmentUpdateOne {
	if euo._type == nil {
		euo._type = make(map[string]struct{})
	}
	euo._type[id] = struct{}{}
	return euo
}

// SetType sets the type edge to EquipmentType.
func (euo *EquipmentUpdateOne) SetType(e *EquipmentType) *EquipmentUpdateOne {
	return euo.SetTypeID(e.ID)
}

// SetLocationID sets the location edge to Location by id.
func (euo *EquipmentUpdateOne) SetLocationID(id string) *EquipmentUpdateOne {
	if euo.location == nil {
		euo.location = make(map[string]struct{})
	}
	euo.location[id] = struct{}{}
	return euo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableLocationID(id *string) *EquipmentUpdateOne {
	if id != nil {
		euo = euo.SetLocationID(*id)
	}
	return euo
}

// SetLocation sets the location edge to Location.
func (euo *EquipmentUpdateOne) SetLocation(l *Location) *EquipmentUpdateOne {
	return euo.SetLocationID(l.ID)
}

// SetParentPositionID sets the parent_position edge to EquipmentPosition by id.
func (euo *EquipmentUpdateOne) SetParentPositionID(id string) *EquipmentUpdateOne {
	if euo.parent_position == nil {
		euo.parent_position = make(map[string]struct{})
	}
	euo.parent_position[id] = struct{}{}
	return euo
}

// SetNillableParentPositionID sets the parent_position edge to EquipmentPosition by id if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableParentPositionID(id *string) *EquipmentUpdateOne {
	if id != nil {
		euo = euo.SetParentPositionID(*id)
	}
	return euo
}

// SetParentPosition sets the parent_position edge to EquipmentPosition.
func (euo *EquipmentUpdateOne) SetParentPosition(e *EquipmentPosition) *EquipmentUpdateOne {
	return euo.SetParentPositionID(e.ID)
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (euo *EquipmentUpdateOne) AddPositionIDs(ids ...string) *EquipmentUpdateOne {
	if euo.positions == nil {
		euo.positions = make(map[string]struct{})
	}
	for i := range ids {
		euo.positions[ids[i]] = struct{}{}
	}
	return euo
}

// AddPositions adds the positions edges to EquipmentPosition.
func (euo *EquipmentUpdateOne) AddPositions(e ...*EquipmentPosition) *EquipmentUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.AddPositionIDs(ids...)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (euo *EquipmentUpdateOne) AddPortIDs(ids ...string) *EquipmentUpdateOne {
	if euo.ports == nil {
		euo.ports = make(map[string]struct{})
	}
	for i := range ids {
		euo.ports[ids[i]] = struct{}{}
	}
	return euo
}

// AddPorts adds the ports edges to EquipmentPort.
func (euo *EquipmentUpdateOne) AddPorts(e ...*EquipmentPort) *EquipmentUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (euo *EquipmentUpdateOne) SetWorkOrderID(id string) *EquipmentUpdateOne {
	if euo.work_order == nil {
		euo.work_order = make(map[string]struct{})
	}
	euo.work_order[id] = struct{}{}
	return euo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (euo *EquipmentUpdateOne) SetNillableWorkOrderID(id *string) *EquipmentUpdateOne {
	if id != nil {
		euo = euo.SetWorkOrderID(*id)
	}
	return euo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (euo *EquipmentUpdateOne) SetWorkOrder(w *WorkOrder) *EquipmentUpdateOne {
	return euo.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (euo *EquipmentUpdateOne) AddPropertyIDs(ids ...string) *EquipmentUpdateOne {
	if euo.properties == nil {
		euo.properties = make(map[string]struct{})
	}
	for i := range ids {
		euo.properties[ids[i]] = struct{}{}
	}
	return euo
}

// AddProperties adds the properties edges to Property.
func (euo *EquipmentUpdateOne) AddProperties(p ...*Property) *EquipmentUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.AddPropertyIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (euo *EquipmentUpdateOne) AddFileIDs(ids ...string) *EquipmentUpdateOne {
	if euo.files == nil {
		euo.files = make(map[string]struct{})
	}
	for i := range ids {
		euo.files[ids[i]] = struct{}{}
	}
	return euo
}

// AddFiles adds the files edges to File.
func (euo *EquipmentUpdateOne) AddFiles(f ...*File) *EquipmentUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return euo.AddFileIDs(ids...)
}

// ClearType clears the type edge to EquipmentType.
func (euo *EquipmentUpdateOne) ClearType() *EquipmentUpdateOne {
	euo.clearedType = true
	return euo
}

// ClearLocation clears the location edge to Location.
func (euo *EquipmentUpdateOne) ClearLocation() *EquipmentUpdateOne {
	euo.clearedLocation = true
	return euo
}

// ClearParentPosition clears the parent_position edge to EquipmentPosition.
func (euo *EquipmentUpdateOne) ClearParentPosition() *EquipmentUpdateOne {
	euo.clearedParentPosition = true
	return euo
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (euo *EquipmentUpdateOne) RemovePositionIDs(ids ...string) *EquipmentUpdateOne {
	if euo.removedPositions == nil {
		euo.removedPositions = make(map[string]struct{})
	}
	for i := range ids {
		euo.removedPositions[ids[i]] = struct{}{}
	}
	return euo
}

// RemovePositions removes positions edges to EquipmentPosition.
func (euo *EquipmentUpdateOne) RemovePositions(e ...*EquipmentPosition) *EquipmentUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.RemovePositionIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (euo *EquipmentUpdateOne) RemovePortIDs(ids ...string) *EquipmentUpdateOne {
	if euo.removedPorts == nil {
		euo.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		euo.removedPorts[ids[i]] = struct{}{}
	}
	return euo
}

// RemovePorts removes ports edges to EquipmentPort.
func (euo *EquipmentUpdateOne) RemovePorts(e ...*EquipmentPort) *EquipmentUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return euo.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (euo *EquipmentUpdateOne) ClearWorkOrder() *EquipmentUpdateOne {
	euo.clearedWorkOrder = true
	return euo
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (euo *EquipmentUpdateOne) RemovePropertyIDs(ids ...string) *EquipmentUpdateOne {
	if euo.removedProperties == nil {
		euo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		euo.removedProperties[ids[i]] = struct{}{}
	}
	return euo
}

// RemoveProperties removes properties edges to Property.
func (euo *EquipmentUpdateOne) RemoveProperties(p ...*Property) *EquipmentUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.RemovePropertyIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (euo *EquipmentUpdateOne) RemoveFileIDs(ids ...string) *EquipmentUpdateOne {
	if euo.removedFiles == nil {
		euo.removedFiles = make(map[string]struct{})
	}
	for i := range ids {
		euo.removedFiles[ids[i]] = struct{}{}
	}
	return euo
}

// RemoveFiles removes files edges to File.
func (euo *EquipmentUpdateOne) RemoveFiles(f ...*File) *EquipmentUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return euo.RemoveFileIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (euo *EquipmentUpdateOne) Save(ctx context.Context) (*Equipment, error) {
	if euo.update_time == nil {
		v := equipment.UpdateDefaultUpdateTime()
		euo.update_time = &v
	}
	if euo.name != nil {
		if err := equipment.NameValidator(*euo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if len(euo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if euo.clearedType && euo._type == nil {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}
	if len(euo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(euo.parent_position) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent_position\"")
	}
	if len(euo.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return euo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EquipmentUpdateOne) SaveX(ctx context.Context) *Equipment {
	e, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return e
}

// Exec executes the query on the entity.
func (euo *EquipmentUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EquipmentUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (euo *EquipmentUpdateOne) sqlSave(ctx context.Context) (e *Equipment, err error) {
	var (
		builder  = sql.Dialect(euo.driver.Dialect())
		selector = builder.Select(equipment.Columns...).From(builder.Table(equipment.Table))
	)
	equipment.ID(euo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = euo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		e = &Equipment{config: euo.config}
		if err := e.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Equipment: %v", err)
		}
		id = e.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Equipment with id: %v", euo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Equipment with the same id: %v", euo.id)
	}

	tx, err := euo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipment.Table)
	)
	updater = updater.Where(sql.InInts(equipment.FieldID, ids...))
	if value := euo.update_time; value != nil {
		updater.Set(equipment.FieldUpdateTime, *value)
		e.UpdateTime = *value
	}
	if value := euo.name; value != nil {
		updater.Set(equipment.FieldName, *value)
		e.Name = *value
	}
	if value := euo.future_state; value != nil {
		updater.Set(equipment.FieldFutureState, *value)
		e.FutureState = *value
	}
	if euo.clearfuture_state {
		var value string
		e.FutureState = value
		updater.SetNull(equipment.FieldFutureState)
	}
	if value := euo.device_id; value != nil {
		updater.Set(equipment.FieldDeviceID, *value)
		e.DeviceID = *value
	}
	if euo.cleardevice_id {
		var value string
		e.DeviceID = value
		updater.SetNull(equipment.FieldDeviceID)
	}
	if value := euo.external_id; value != nil {
		updater.Set(equipment.FieldExternalID, *value)
		e.ExternalID = *value
	}
	if euo.clearexternal_id {
		var value string
		e.ExternalID = value
		updater.SetNull(equipment.FieldExternalID)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if euo.clearedType {
		query, args := builder.Update(equipment.TypeTable).
			SetNull(equipment.TypeColumn).
			Where(sql.InInts(equipmenttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo._type) > 0 {
		for eid := range euo._type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipment.TypeTable).
				Set(equipment.TypeColumn, eid).
				Where(sql.InInts(equipment.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if euo.clearedLocation {
		query, args := builder.Update(equipment.LocationTable).
			SetNull(equipment.LocationColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo.location) > 0 {
		for eid := range euo.location {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipment.LocationTable).
				Set(equipment.LocationColumn, eid).
				Where(sql.InInts(equipment.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if euo.clearedParentPosition {
		query, args := builder.Update(equipment.ParentPositionTable).
			SetNull(equipment.ParentPositionColumn).
			Where(sql.InInts(equipmentposition.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo.parent_position) > 0 {
		for _, id := range ids {
			eid, serr := strconv.Atoi(keys(euo.parent_position)[0])
			if serr != nil {
				return nil, rollback(tx, err)
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
			if int(affected) < len(euo.parent_position) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"parent_position\" %v already connected to a different \"Equipment\"", keys(euo.parent_position))})
			}
		}
	}
	if len(euo.removedPositions) > 0 {
		eids := make([]int, len(euo.removedPositions))
		for eid := range euo.removedPositions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.PositionsTable).
			SetNull(equipment.PositionsColumn).
			Where(sql.InInts(equipment.PositionsColumn, ids...)).
			Where(sql.InInts(equipmentposition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo.positions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range euo.positions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(euo.positions) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"positions\" %v already connected to a different \"Equipment\"", keys(euo.positions))})
			}
		}
	}
	if len(euo.removedPorts) > 0 {
		eids := make([]int, len(euo.removedPorts))
		for eid := range euo.removedPorts {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.PortsTable).
			SetNull(equipment.PortsColumn).
			Where(sql.InInts(equipment.PortsColumn, ids...)).
			Where(sql.InInts(equipmentport.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo.ports) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range euo.ports {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(euo.ports) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"Equipment\"", keys(euo.ports))})
			}
		}
	}
	if euo.clearedWorkOrder {
		query, args := builder.Update(equipment.WorkOrderTable).
			SetNull(equipment.WorkOrderColumn).
			Where(sql.InInts(workorder.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo.work_order) > 0 {
		for eid := range euo.work_order {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipment.WorkOrderTable).
				Set(equipment.WorkOrderColumn, eid).
				Where(sql.InInts(equipment.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(euo.removedProperties) > 0 {
		eids := make([]int, len(euo.removedProperties))
		for eid := range euo.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.PropertiesTable).
			SetNull(equipment.PropertiesColumn).
			Where(sql.InInts(equipment.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range euo.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(euo.properties) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Equipment\"", keys(euo.properties))})
			}
		}
	}
	if len(euo.removedFiles) > 0 {
		eids := make([]int, len(euo.removedFiles))
		for eid := range euo.removedFiles {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipment.FilesTable).
			SetNull(equipment.FilesColumn).
			Where(sql.InInts(equipment.FilesColumn, ids...)).
			Where(sql.InInts(file.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(euo.files) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range euo.files {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(euo.files) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"files\" %v already connected to a different \"Equipment\"", keys(euo.files))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return e, nil
}
