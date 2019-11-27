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
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPortDefinitionUpdate is the builder for updating EquipmentPortDefinition entities.
type EquipmentPortDefinitionUpdate struct {
	config

	update_time              *time.Time
	name                     *string
	_type                    *string
	index                    *int
	addindex                 *int
	clearindex               bool
	bandwidth                *string
	clearbandwidth           bool
	visibility_label         *string
	clearvisibility_label    bool
	equipment_port_type      map[string]struct{}
	ports                    map[string]struct{}
	equipment_type           map[string]struct{}
	clearedEquipmentPortType bool
	removedPorts             map[string]struct{}
	clearedEquipmentType     bool
	predicates               []predicate.EquipmentPortDefinition
}

// Where adds a new predicate for the builder.
func (epdu *EquipmentPortDefinitionUpdate) Where(ps ...predicate.EquipmentPortDefinition) *EquipmentPortDefinitionUpdate {
	epdu.predicates = append(epdu.predicates, ps...)
	return epdu
}

// SetName sets the name field.
func (epdu *EquipmentPortDefinitionUpdate) SetName(s string) *EquipmentPortDefinitionUpdate {
	epdu.name = &s
	return epdu
}

// SetType sets the type field.
func (epdu *EquipmentPortDefinitionUpdate) SetType(s string) *EquipmentPortDefinitionUpdate {
	epdu._type = &s
	return epdu
}

// SetIndex sets the index field.
func (epdu *EquipmentPortDefinitionUpdate) SetIndex(i int) *EquipmentPortDefinitionUpdate {
	epdu.index = &i
	epdu.addindex = nil
	return epdu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableIndex(i *int) *EquipmentPortDefinitionUpdate {
	if i != nil {
		epdu.SetIndex(*i)
	}
	return epdu
}

// AddIndex adds i to index.
func (epdu *EquipmentPortDefinitionUpdate) AddIndex(i int) *EquipmentPortDefinitionUpdate {
	if epdu.addindex == nil {
		epdu.addindex = &i
	} else {
		*epdu.addindex += i
	}
	return epdu
}

// ClearIndex clears the value of index.
func (epdu *EquipmentPortDefinitionUpdate) ClearIndex() *EquipmentPortDefinitionUpdate {
	epdu.index = nil
	epdu.clearindex = true
	return epdu
}

// SetBandwidth sets the bandwidth field.
func (epdu *EquipmentPortDefinitionUpdate) SetBandwidth(s string) *EquipmentPortDefinitionUpdate {
	epdu.bandwidth = &s
	return epdu
}

// SetNillableBandwidth sets the bandwidth field if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableBandwidth(s *string) *EquipmentPortDefinitionUpdate {
	if s != nil {
		epdu.SetBandwidth(*s)
	}
	return epdu
}

// ClearBandwidth clears the value of bandwidth.
func (epdu *EquipmentPortDefinitionUpdate) ClearBandwidth() *EquipmentPortDefinitionUpdate {
	epdu.bandwidth = nil
	epdu.clearbandwidth = true
	return epdu
}

// SetVisibilityLabel sets the visibility_label field.
func (epdu *EquipmentPortDefinitionUpdate) SetVisibilityLabel(s string) *EquipmentPortDefinitionUpdate {
	epdu.visibility_label = &s
	return epdu
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableVisibilityLabel(s *string) *EquipmentPortDefinitionUpdate {
	if s != nil {
		epdu.SetVisibilityLabel(*s)
	}
	return epdu
}

// ClearVisibilityLabel clears the value of visibility_label.
func (epdu *EquipmentPortDefinitionUpdate) ClearVisibilityLabel() *EquipmentPortDefinitionUpdate {
	epdu.visibility_label = nil
	epdu.clearvisibility_label = true
	return epdu
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentPortTypeID(id string) *EquipmentPortDefinitionUpdate {
	if epdu.equipment_port_type == nil {
		epdu.equipment_port_type = make(map[string]struct{})
	}
	epdu.equipment_port_type[id] = struct{}{}
	return epdu
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableEquipmentPortTypeID(id *string) *EquipmentPortDefinitionUpdate {
	if id != nil {
		epdu = epdu.SetEquipmentPortTypeID(*id)
	}
	return epdu
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentPortType(e *EquipmentPortType) *EquipmentPortDefinitionUpdate {
	return epdu.SetEquipmentPortTypeID(e.ID)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (epdu *EquipmentPortDefinitionUpdate) AddPortIDs(ids ...string) *EquipmentPortDefinitionUpdate {
	if epdu.ports == nil {
		epdu.ports = make(map[string]struct{})
	}
	for i := range ids {
		epdu.ports[ids[i]] = struct{}{}
	}
	return epdu
}

// AddPorts adds the ports edges to EquipmentPort.
func (epdu *EquipmentPortDefinitionUpdate) AddPorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdu.AddPortIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentTypeID(id string) *EquipmentPortDefinitionUpdate {
	if epdu.equipment_type == nil {
		epdu.equipment_type = make(map[string]struct{})
	}
	epdu.equipment_type[id] = struct{}{}
	return epdu
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableEquipmentTypeID(id *string) *EquipmentPortDefinitionUpdate {
	if id != nil {
		epdu = epdu.SetEquipmentTypeID(*id)
	}
	return epdu
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentType(e *EquipmentType) *EquipmentPortDefinitionUpdate {
	return epdu.SetEquipmentTypeID(e.ID)
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (epdu *EquipmentPortDefinitionUpdate) ClearEquipmentPortType() *EquipmentPortDefinitionUpdate {
	epdu.clearedEquipmentPortType = true
	return epdu
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (epdu *EquipmentPortDefinitionUpdate) RemovePortIDs(ids ...string) *EquipmentPortDefinitionUpdate {
	if epdu.removedPorts == nil {
		epdu.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		epdu.removedPorts[ids[i]] = struct{}{}
	}
	return epdu
}

// RemovePorts removes ports edges to EquipmentPort.
func (epdu *EquipmentPortDefinitionUpdate) RemovePorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdu.RemovePortIDs(ids...)
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (epdu *EquipmentPortDefinitionUpdate) ClearEquipmentType() *EquipmentPortDefinitionUpdate {
	epdu.clearedEquipmentType = true
	return epdu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (epdu *EquipmentPortDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if epdu.update_time == nil {
		v := equipmentportdefinition.UpdateDefaultUpdateTime()
		epdu.update_time = &v
	}
	if len(epdu.equipment_port_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"equipment_port_type\"")
	}
	if len(epdu.equipment_type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"equipment_type\"")
	}
	return epdu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (epdu *EquipmentPortDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := epdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (epdu *EquipmentPortDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := epdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epdu *EquipmentPortDefinitionUpdate) ExecX(ctx context.Context) {
	if err := epdu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epdu *EquipmentPortDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(epdu.driver.Dialect())
		selector = builder.Select(equipmentportdefinition.FieldID).From(builder.Table(equipmentportdefinition.Table))
	)
	for _, p := range epdu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = epdu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := epdu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentportdefinition.Table)
	)
	updater = updater.Where(sql.InInts(equipmentportdefinition.FieldID, ids...))
	if value := epdu.update_time; value != nil {
		updater.Set(equipmentportdefinition.FieldUpdateTime, *value)
	}
	if value := epdu.name; value != nil {
		updater.Set(equipmentportdefinition.FieldName, *value)
	}
	if value := epdu._type; value != nil {
		updater.Set(equipmentportdefinition.FieldType, *value)
	}
	if value := epdu.index; value != nil {
		updater.Set(equipmentportdefinition.FieldIndex, *value)
	}
	if value := epdu.addindex; value != nil {
		updater.Add(equipmentportdefinition.FieldIndex, *value)
	}
	if epdu.clearindex {
		updater.SetNull(equipmentportdefinition.FieldIndex)
	}
	if value := epdu.bandwidth; value != nil {
		updater.Set(equipmentportdefinition.FieldBandwidth, *value)
	}
	if epdu.clearbandwidth {
		updater.SetNull(equipmentportdefinition.FieldBandwidth)
	}
	if value := epdu.visibility_label; value != nil {
		updater.Set(equipmentportdefinition.FieldVisibilityLabel, *value)
	}
	if epdu.clearvisibility_label {
		updater.SetNull(equipmentportdefinition.FieldVisibilityLabel)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if epdu.clearedEquipmentPortType {
		query, args := builder.Update(equipmentportdefinition.EquipmentPortTypeTable).
			SetNull(equipmentportdefinition.EquipmentPortTypeColumn).
			Where(sql.InInts(equipmentporttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epdu.equipment_port_type) > 0 {
		for eid := range epdu.equipment_port_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentportdefinition.EquipmentPortTypeTable).
				Set(equipmentportdefinition.EquipmentPortTypeColumn, eid).
				Where(sql.InInts(equipmentportdefinition.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if len(epdu.removedPorts) > 0 {
		eids := make([]int, len(epdu.removedPorts))
		for eid := range epdu.removedPorts {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentportdefinition.PortsTable).
			SetNull(equipmentportdefinition.PortsColumn).
			Where(sql.InInts(equipmentportdefinition.PortsColumn, ids...)).
			Where(sql.InInts(equipmentport.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epdu.ports) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range epdu.ports {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentport.FieldID, eid)
			}
			query, args := builder.Update(equipmentportdefinition.PortsTable).
				Set(equipmentportdefinition.PortsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentportdefinition.PortsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(epdu.ports) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"EquipmentPortDefinition\"", keys(epdu.ports))})
			}
		}
	}
	if epdu.clearedEquipmentType {
		query, args := builder.Update(equipmentportdefinition.EquipmentTypeTable).
			SetNull(equipmentportdefinition.EquipmentTypeColumn).
			Where(sql.InInts(equipmenttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epdu.equipment_type) > 0 {
		for eid := range epdu.equipment_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentportdefinition.EquipmentTypeTable).
				Set(equipmentportdefinition.EquipmentTypeColumn, eid).
				Where(sql.InInts(equipmentportdefinition.FieldID, ids...)).
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

// EquipmentPortDefinitionUpdateOne is the builder for updating a single EquipmentPortDefinition entity.
type EquipmentPortDefinitionUpdateOne struct {
	config
	id string

	update_time              *time.Time
	name                     *string
	_type                    *string
	index                    *int
	addindex                 *int
	clearindex               bool
	bandwidth                *string
	clearbandwidth           bool
	visibility_label         *string
	clearvisibility_label    bool
	equipment_port_type      map[string]struct{}
	ports                    map[string]struct{}
	equipment_type           map[string]struct{}
	clearedEquipmentPortType bool
	removedPorts             map[string]struct{}
	clearedEquipmentType     bool
}

// SetName sets the name field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetName(s string) *EquipmentPortDefinitionUpdateOne {
	epduo.name = &s
	return epduo
}

// SetType sets the type field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetType(s string) *EquipmentPortDefinitionUpdateOne {
	epduo._type = &s
	return epduo
}

// SetIndex sets the index field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetIndex(i int) *EquipmentPortDefinitionUpdateOne {
	epduo.index = &i
	epduo.addindex = nil
	return epduo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableIndex(i *int) *EquipmentPortDefinitionUpdateOne {
	if i != nil {
		epduo.SetIndex(*i)
	}
	return epduo
}

// AddIndex adds i to index.
func (epduo *EquipmentPortDefinitionUpdateOne) AddIndex(i int) *EquipmentPortDefinitionUpdateOne {
	if epduo.addindex == nil {
		epduo.addindex = &i
	} else {
		*epduo.addindex += i
	}
	return epduo
}

// ClearIndex clears the value of index.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearIndex() *EquipmentPortDefinitionUpdateOne {
	epduo.index = nil
	epduo.clearindex = true
	return epduo
}

// SetBandwidth sets the bandwidth field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetBandwidth(s string) *EquipmentPortDefinitionUpdateOne {
	epduo.bandwidth = &s
	return epduo
}

// SetNillableBandwidth sets the bandwidth field if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableBandwidth(s *string) *EquipmentPortDefinitionUpdateOne {
	if s != nil {
		epduo.SetBandwidth(*s)
	}
	return epduo
}

// ClearBandwidth clears the value of bandwidth.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearBandwidth() *EquipmentPortDefinitionUpdateOne {
	epduo.bandwidth = nil
	epduo.clearbandwidth = true
	return epduo
}

// SetVisibilityLabel sets the visibility_label field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetVisibilityLabel(s string) *EquipmentPortDefinitionUpdateOne {
	epduo.visibility_label = &s
	return epduo
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableVisibilityLabel(s *string) *EquipmentPortDefinitionUpdateOne {
	if s != nil {
		epduo.SetVisibilityLabel(*s)
	}
	return epduo
}

// ClearVisibilityLabel clears the value of visibility_label.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearVisibilityLabel() *EquipmentPortDefinitionUpdateOne {
	epduo.visibility_label = nil
	epduo.clearvisibility_label = true
	return epduo
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentPortTypeID(id string) *EquipmentPortDefinitionUpdateOne {
	if epduo.equipment_port_type == nil {
		epduo.equipment_port_type = make(map[string]struct{})
	}
	epduo.equipment_port_type[id] = struct{}{}
	return epduo
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableEquipmentPortTypeID(id *string) *EquipmentPortDefinitionUpdateOne {
	if id != nil {
		epduo = epduo.SetEquipmentPortTypeID(*id)
	}
	return epduo
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentPortType(e *EquipmentPortType) *EquipmentPortDefinitionUpdateOne {
	return epduo.SetEquipmentPortTypeID(e.ID)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (epduo *EquipmentPortDefinitionUpdateOne) AddPortIDs(ids ...string) *EquipmentPortDefinitionUpdateOne {
	if epduo.ports == nil {
		epduo.ports = make(map[string]struct{})
	}
	for i := range ids {
		epduo.ports[ids[i]] = struct{}{}
	}
	return epduo
}

// AddPorts adds the ports edges to EquipmentPort.
func (epduo *EquipmentPortDefinitionUpdateOne) AddPorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epduo.AddPortIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentTypeID(id string) *EquipmentPortDefinitionUpdateOne {
	if epduo.equipment_type == nil {
		epduo.equipment_type = make(map[string]struct{})
	}
	epduo.equipment_type[id] = struct{}{}
	return epduo
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableEquipmentTypeID(id *string) *EquipmentPortDefinitionUpdateOne {
	if id != nil {
		epduo = epduo.SetEquipmentTypeID(*id)
	}
	return epduo
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentType(e *EquipmentType) *EquipmentPortDefinitionUpdateOne {
	return epduo.SetEquipmentTypeID(e.ID)
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearEquipmentPortType() *EquipmentPortDefinitionUpdateOne {
	epduo.clearedEquipmentPortType = true
	return epduo
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (epduo *EquipmentPortDefinitionUpdateOne) RemovePortIDs(ids ...string) *EquipmentPortDefinitionUpdateOne {
	if epduo.removedPorts == nil {
		epduo.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		epduo.removedPorts[ids[i]] = struct{}{}
	}
	return epduo
}

// RemovePorts removes ports edges to EquipmentPort.
func (epduo *EquipmentPortDefinitionUpdateOne) RemovePorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epduo.RemovePortIDs(ids...)
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearEquipmentType() *EquipmentPortDefinitionUpdateOne {
	epduo.clearedEquipmentType = true
	return epduo
}

// Save executes the query and returns the updated entity.
func (epduo *EquipmentPortDefinitionUpdateOne) Save(ctx context.Context) (*EquipmentPortDefinition, error) {
	if epduo.update_time == nil {
		v := equipmentportdefinition.UpdateDefaultUpdateTime()
		epduo.update_time = &v
	}
	if len(epduo.equipment_port_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_port_type\"")
	}
	if len(epduo.equipment_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_type\"")
	}
	return epduo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (epduo *EquipmentPortDefinitionUpdateOne) SaveX(ctx context.Context) *EquipmentPortDefinition {
	epd, err := epduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return epd
}

// Exec executes the query on the entity.
func (epduo *EquipmentPortDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := epduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epduo *EquipmentPortDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := epduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epduo *EquipmentPortDefinitionUpdateOne) sqlSave(ctx context.Context) (epd *EquipmentPortDefinition, err error) {
	var (
		builder  = sql.Dialect(epduo.driver.Dialect())
		selector = builder.Select(equipmentportdefinition.Columns...).From(builder.Table(equipmentportdefinition.Table))
	)
	equipmentportdefinition.ID(epduo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = epduo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		epd = &EquipmentPortDefinition{config: epduo.config}
		if err := epd.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into EquipmentPortDefinition: %v", err)
		}
		id = epd.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("EquipmentPortDefinition with id: %v", epduo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one EquipmentPortDefinition with the same id: %v", epduo.id)
	}

	tx, err := epduo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentportdefinition.Table)
	)
	updater = updater.Where(sql.InInts(equipmentportdefinition.FieldID, ids...))
	if value := epduo.update_time; value != nil {
		updater.Set(equipmentportdefinition.FieldUpdateTime, *value)
		epd.UpdateTime = *value
	}
	if value := epduo.name; value != nil {
		updater.Set(equipmentportdefinition.FieldName, *value)
		epd.Name = *value
	}
	if value := epduo._type; value != nil {
		updater.Set(equipmentportdefinition.FieldType, *value)
		epd.Type = *value
	}
	if value := epduo.index; value != nil {
		updater.Set(equipmentportdefinition.FieldIndex, *value)
		epd.Index = *value
	}
	if value := epduo.addindex; value != nil {
		updater.Add(equipmentportdefinition.FieldIndex, *value)
		epd.Index += *value
	}
	if epduo.clearindex {
		var value int
		epd.Index = value
		updater.SetNull(equipmentportdefinition.FieldIndex)
	}
	if value := epduo.bandwidth; value != nil {
		updater.Set(equipmentportdefinition.FieldBandwidth, *value)
		epd.Bandwidth = *value
	}
	if epduo.clearbandwidth {
		var value string
		epd.Bandwidth = value
		updater.SetNull(equipmentportdefinition.FieldBandwidth)
	}
	if value := epduo.visibility_label; value != nil {
		updater.Set(equipmentportdefinition.FieldVisibilityLabel, *value)
		epd.VisibilityLabel = *value
	}
	if epduo.clearvisibility_label {
		var value string
		epd.VisibilityLabel = value
		updater.SetNull(equipmentportdefinition.FieldVisibilityLabel)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if epduo.clearedEquipmentPortType {
		query, args := builder.Update(equipmentportdefinition.EquipmentPortTypeTable).
			SetNull(equipmentportdefinition.EquipmentPortTypeColumn).
			Where(sql.InInts(equipmentporttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epduo.equipment_port_type) > 0 {
		for eid := range epduo.equipment_port_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentportdefinition.EquipmentPortTypeTable).
				Set(equipmentportdefinition.EquipmentPortTypeColumn, eid).
				Where(sql.InInts(equipmentportdefinition.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(epduo.removedPorts) > 0 {
		eids := make([]int, len(epduo.removedPorts))
		for eid := range epduo.removedPorts {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentportdefinition.PortsTable).
			SetNull(equipmentportdefinition.PortsColumn).
			Where(sql.InInts(equipmentportdefinition.PortsColumn, ids...)).
			Where(sql.InInts(equipmentport.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epduo.ports) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range epduo.ports {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentport.FieldID, eid)
			}
			query, args := builder.Update(equipmentportdefinition.PortsTable).
				Set(equipmentportdefinition.PortsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentportdefinition.PortsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(epduo.ports) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"EquipmentPortDefinition\"", keys(epduo.ports))})
			}
		}
	}
	if epduo.clearedEquipmentType {
		query, args := builder.Update(equipmentportdefinition.EquipmentTypeTable).
			SetNull(equipmentportdefinition.EquipmentTypeColumn).
			Where(sql.InInts(equipmenttype.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epduo.equipment_type) > 0 {
		for eid := range epduo.equipment_type {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentportdefinition.EquipmentTypeTable).
				Set(equipmentportdefinition.EquipmentTypeColumn, eid).
				Where(sql.InInts(equipmentportdefinition.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return epd, nil
}
