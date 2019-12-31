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
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
)

// EquipmentPortUpdate is the builder for updating EquipmentPort entities.
type EquipmentPortUpdate struct {
	config

	update_time       *time.Time
	definition        map[string]struct{}
	parent            map[string]struct{}
	link              map[string]struct{}
	properties        map[string]struct{}
	endpoints         map[string]struct{}
	clearedDefinition bool
	clearedParent     bool
	clearedLink       bool
	removedProperties map[string]struct{}
	removedEndpoints  map[string]struct{}
	predicates        []predicate.EquipmentPort
}

// Where adds a new predicate for the builder.
func (epu *EquipmentPortUpdate) Where(ps ...predicate.EquipmentPort) *EquipmentPortUpdate {
	epu.predicates = append(epu.predicates, ps...)
	return epu
}

// SetDefinitionID sets the definition edge to EquipmentPortDefinition by id.
func (epu *EquipmentPortUpdate) SetDefinitionID(id string) *EquipmentPortUpdate {
	if epu.definition == nil {
		epu.definition = make(map[string]struct{})
	}
	epu.definition[id] = struct{}{}
	return epu
}

// SetDefinition sets the definition edge to EquipmentPortDefinition.
func (epu *EquipmentPortUpdate) SetDefinition(e *EquipmentPortDefinition) *EquipmentPortUpdate {
	return epu.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epu *EquipmentPortUpdate) SetParentID(id string) *EquipmentPortUpdate {
	if epu.parent == nil {
		epu.parent = make(map[string]struct{})
	}
	epu.parent[id] = struct{}{}
	return epu
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epu *EquipmentPortUpdate) SetNillableParentID(id *string) *EquipmentPortUpdate {
	if id != nil {
		epu = epu.SetParentID(*id)
	}
	return epu
}

// SetParent sets the parent edge to Equipment.
func (epu *EquipmentPortUpdate) SetParent(e *Equipment) *EquipmentPortUpdate {
	return epu.SetParentID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (epu *EquipmentPortUpdate) SetLinkID(id string) *EquipmentPortUpdate {
	if epu.link == nil {
		epu.link = make(map[string]struct{})
	}
	epu.link[id] = struct{}{}
	return epu
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (epu *EquipmentPortUpdate) SetNillableLinkID(id *string) *EquipmentPortUpdate {
	if id != nil {
		epu = epu.SetLinkID(*id)
	}
	return epu
}

// SetLink sets the link edge to Link.
func (epu *EquipmentPortUpdate) SetLink(l *Link) *EquipmentPortUpdate {
	return epu.SetLinkID(l.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (epu *EquipmentPortUpdate) AddPropertyIDs(ids ...string) *EquipmentPortUpdate {
	if epu.properties == nil {
		epu.properties = make(map[string]struct{})
	}
	for i := range ids {
		epu.properties[ids[i]] = struct{}{}
	}
	return epu
}

// AddProperties adds the properties edges to Property.
func (epu *EquipmentPortUpdate) AddProperties(p ...*Property) *EquipmentPortUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epu.AddPropertyIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (epu *EquipmentPortUpdate) AddEndpointIDs(ids ...string) *EquipmentPortUpdate {
	if epu.endpoints == nil {
		epu.endpoints = make(map[string]struct{})
	}
	for i := range ids {
		epu.endpoints[ids[i]] = struct{}{}
	}
	return epu
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (epu *EquipmentPortUpdate) AddEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epu.AddEndpointIDs(ids...)
}

// ClearDefinition clears the definition edge to EquipmentPortDefinition.
func (epu *EquipmentPortUpdate) ClearDefinition() *EquipmentPortUpdate {
	epu.clearedDefinition = true
	return epu
}

// ClearParent clears the parent edge to Equipment.
func (epu *EquipmentPortUpdate) ClearParent() *EquipmentPortUpdate {
	epu.clearedParent = true
	return epu
}

// ClearLink clears the link edge to Link.
func (epu *EquipmentPortUpdate) ClearLink() *EquipmentPortUpdate {
	epu.clearedLink = true
	return epu
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (epu *EquipmentPortUpdate) RemovePropertyIDs(ids ...string) *EquipmentPortUpdate {
	if epu.removedProperties == nil {
		epu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		epu.removedProperties[ids[i]] = struct{}{}
	}
	return epu
}

// RemoveProperties removes properties edges to Property.
func (epu *EquipmentPortUpdate) RemoveProperties(p ...*Property) *EquipmentPortUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epu.RemovePropertyIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (epu *EquipmentPortUpdate) RemoveEndpointIDs(ids ...string) *EquipmentPortUpdate {
	if epu.removedEndpoints == nil {
		epu.removedEndpoints = make(map[string]struct{})
	}
	for i := range ids {
		epu.removedEndpoints[ids[i]] = struct{}{}
	}
	return epu
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (epu *EquipmentPortUpdate) RemoveEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epu.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (epu *EquipmentPortUpdate) Save(ctx context.Context) (int, error) {
	if epu.update_time == nil {
		v := equipmentport.UpdateDefaultUpdateTime()
		epu.update_time = &v
	}
	if len(epu.definition) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"definition\"")
	}
	if epu.clearedDefinition && epu.definition == nil {
		return 0, errors.New("ent: clearing a unique edge \"definition\"")
	}
	if len(epu.parent) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	if len(epu.link) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"link\"")
	}
	return epu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (epu *EquipmentPortUpdate) SaveX(ctx context.Context) int {
	affected, err := epu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (epu *EquipmentPortUpdate) Exec(ctx context.Context) error {
	_, err := epu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epu *EquipmentPortUpdate) ExecX(ctx context.Context) {
	if err := epu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epu *EquipmentPortUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(epu.driver.Dialect())
		selector = builder.Select(equipmentport.FieldID).From(builder.Table(equipmentport.Table))
	)
	for _, p := range epu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = epu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := epu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentport.Table)
	)
	updater = updater.Where(sql.InInts(equipmentport.FieldID, ids...))
	if value := epu.update_time; value != nil {
		updater.Set(equipmentport.FieldUpdateTime, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if epu.clearedDefinition {
		query, args := builder.Update(equipmentport.DefinitionTable).
			SetNull(equipmentport.DefinitionColumn).
			Where(sql.InInts(equipmentportdefinition.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epu.definition) > 0 {
		for eid := range epu.definition {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentport.DefinitionTable).
				Set(equipmentport.DefinitionColumn, eid).
				Where(sql.InInts(equipmentport.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if epu.clearedParent {
		query, args := builder.Update(equipmentport.ParentTable).
			SetNull(equipmentport.ParentColumn).
			Where(sql.InInts(equipment.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epu.parent) > 0 {
		for eid := range epu.parent {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentport.ParentTable).
				Set(equipmentport.ParentColumn, eid).
				Where(sql.InInts(equipmentport.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if epu.clearedLink {
		query, args := builder.Update(equipmentport.LinkTable).
			SetNull(equipmentport.LinkColumn).
			Where(sql.InInts(link.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epu.link) > 0 {
		for eid := range epu.link {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentport.LinkTable).
				Set(equipmentport.LinkColumn, eid).
				Where(sql.InInts(equipmentport.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if len(epu.removedProperties) > 0 {
		eids := make([]int, len(epu.removedProperties))
		for eid := range epu.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentport.PropertiesTable).
			SetNull(equipmentport.PropertiesColumn).
			Where(sql.InInts(equipmentport.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epu.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range epu.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(equipmentport.PropertiesTable).
				Set(equipmentport.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentport.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(epu.properties) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"EquipmentPort\"", keys(epu.properties))})
			}
		}
	}
	if len(epu.removedEndpoints) > 0 {
		eids := make([]int, len(epu.removedEndpoints))
		for eid := range epu.removedEndpoints {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentport.EndpointsTable).
			SetNull(equipmentport.EndpointsColumn).
			Where(sql.InInts(equipmentport.EndpointsColumn, ids...)).
			Where(sql.InInts(serviceendpoint.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(epu.endpoints) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range epu.endpoints {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(serviceendpoint.FieldID, eid)
			}
			query, args := builder.Update(equipmentport.EndpointsTable).
				Set(equipmentport.EndpointsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentport.EndpointsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(epu.endpoints) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"endpoints\" %v already connected to a different \"EquipmentPort\"", keys(epu.endpoints))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// EquipmentPortUpdateOne is the builder for updating a single EquipmentPort entity.
type EquipmentPortUpdateOne struct {
	config
	id string

	update_time       *time.Time
	definition        map[string]struct{}
	parent            map[string]struct{}
	link              map[string]struct{}
	properties        map[string]struct{}
	endpoints         map[string]struct{}
	clearedDefinition bool
	clearedParent     bool
	clearedLink       bool
	removedProperties map[string]struct{}
	removedEndpoints  map[string]struct{}
}

// SetDefinitionID sets the definition edge to EquipmentPortDefinition by id.
func (epuo *EquipmentPortUpdateOne) SetDefinitionID(id string) *EquipmentPortUpdateOne {
	if epuo.definition == nil {
		epuo.definition = make(map[string]struct{})
	}
	epuo.definition[id] = struct{}{}
	return epuo
}

// SetDefinition sets the definition edge to EquipmentPortDefinition.
func (epuo *EquipmentPortUpdateOne) SetDefinition(e *EquipmentPortDefinition) *EquipmentPortUpdateOne {
	return epuo.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epuo *EquipmentPortUpdateOne) SetParentID(id string) *EquipmentPortUpdateOne {
	if epuo.parent == nil {
		epuo.parent = make(map[string]struct{})
	}
	epuo.parent[id] = struct{}{}
	return epuo
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epuo *EquipmentPortUpdateOne) SetNillableParentID(id *string) *EquipmentPortUpdateOne {
	if id != nil {
		epuo = epuo.SetParentID(*id)
	}
	return epuo
}

// SetParent sets the parent edge to Equipment.
func (epuo *EquipmentPortUpdateOne) SetParent(e *Equipment) *EquipmentPortUpdateOne {
	return epuo.SetParentID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (epuo *EquipmentPortUpdateOne) SetLinkID(id string) *EquipmentPortUpdateOne {
	if epuo.link == nil {
		epuo.link = make(map[string]struct{})
	}
	epuo.link[id] = struct{}{}
	return epuo
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (epuo *EquipmentPortUpdateOne) SetNillableLinkID(id *string) *EquipmentPortUpdateOne {
	if id != nil {
		epuo = epuo.SetLinkID(*id)
	}
	return epuo
}

// SetLink sets the link edge to Link.
func (epuo *EquipmentPortUpdateOne) SetLink(l *Link) *EquipmentPortUpdateOne {
	return epuo.SetLinkID(l.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (epuo *EquipmentPortUpdateOne) AddPropertyIDs(ids ...string) *EquipmentPortUpdateOne {
	if epuo.properties == nil {
		epuo.properties = make(map[string]struct{})
	}
	for i := range ids {
		epuo.properties[ids[i]] = struct{}{}
	}
	return epuo
}

// AddProperties adds the properties edges to Property.
func (epuo *EquipmentPortUpdateOne) AddProperties(p ...*Property) *EquipmentPortUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epuo.AddPropertyIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (epuo *EquipmentPortUpdateOne) AddEndpointIDs(ids ...string) *EquipmentPortUpdateOne {
	if epuo.endpoints == nil {
		epuo.endpoints = make(map[string]struct{})
	}
	for i := range ids {
		epuo.endpoints[ids[i]] = struct{}{}
	}
	return epuo
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (epuo *EquipmentPortUpdateOne) AddEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epuo.AddEndpointIDs(ids...)
}

// ClearDefinition clears the definition edge to EquipmentPortDefinition.
func (epuo *EquipmentPortUpdateOne) ClearDefinition() *EquipmentPortUpdateOne {
	epuo.clearedDefinition = true
	return epuo
}

// ClearParent clears the parent edge to Equipment.
func (epuo *EquipmentPortUpdateOne) ClearParent() *EquipmentPortUpdateOne {
	epuo.clearedParent = true
	return epuo
}

// ClearLink clears the link edge to Link.
func (epuo *EquipmentPortUpdateOne) ClearLink() *EquipmentPortUpdateOne {
	epuo.clearedLink = true
	return epuo
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (epuo *EquipmentPortUpdateOne) RemovePropertyIDs(ids ...string) *EquipmentPortUpdateOne {
	if epuo.removedProperties == nil {
		epuo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		epuo.removedProperties[ids[i]] = struct{}{}
	}
	return epuo
}

// RemoveProperties removes properties edges to Property.
func (epuo *EquipmentPortUpdateOne) RemoveProperties(p ...*Property) *EquipmentPortUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epuo.RemovePropertyIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (epuo *EquipmentPortUpdateOne) RemoveEndpointIDs(ids ...string) *EquipmentPortUpdateOne {
	if epuo.removedEndpoints == nil {
		epuo.removedEndpoints = make(map[string]struct{})
	}
	for i := range ids {
		epuo.removedEndpoints[ids[i]] = struct{}{}
	}
	return epuo
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (epuo *EquipmentPortUpdateOne) RemoveEndpoints(s ...*ServiceEndpoint) *EquipmentPortUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return epuo.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (epuo *EquipmentPortUpdateOne) Save(ctx context.Context) (*EquipmentPort, error) {
	if epuo.update_time == nil {
		v := equipmentport.UpdateDefaultUpdateTime()
		epuo.update_time = &v
	}
	if len(epuo.definition) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"definition\"")
	}
	if epuo.clearedDefinition && epuo.definition == nil {
		return nil, errors.New("ent: clearing a unique edge \"definition\"")
	}
	if len(epuo.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	if len(epuo.link) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"link\"")
	}
	return epuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (epuo *EquipmentPortUpdateOne) SaveX(ctx context.Context) *EquipmentPort {
	ep, err := epuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ep
}

// Exec executes the query on the entity.
func (epuo *EquipmentPortUpdateOne) Exec(ctx context.Context) error {
	_, err := epuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epuo *EquipmentPortUpdateOne) ExecX(ctx context.Context) {
	if err := epuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epuo *EquipmentPortUpdateOne) sqlSave(ctx context.Context) (ep *EquipmentPort, err error) {
	var (
		builder  = sql.Dialect(epuo.driver.Dialect())
		selector = builder.Select(equipmentport.Columns...).From(builder.Table(equipmentport.Table))
	)
	equipmentport.ID(epuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = epuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		ep = &EquipmentPort{config: epuo.config}
		if err := ep.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into EquipmentPort: %v", err)
		}
		id = ep.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("EquipmentPort with id: %v", epuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one EquipmentPort with the same id: %v", epuo.id)
	}

	tx, err := epuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentport.Table)
	)
	updater = updater.Where(sql.InInts(equipmentport.FieldID, ids...))
	if value := epuo.update_time; value != nil {
		updater.Set(equipmentport.FieldUpdateTime, *value)
		ep.UpdateTime = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if epuo.clearedDefinition {
		query, args := builder.Update(equipmentport.DefinitionTable).
			SetNull(equipmentport.DefinitionColumn).
			Where(sql.InInts(equipmentportdefinition.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epuo.definition) > 0 {
		for eid := range epuo.definition {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentport.DefinitionTable).
				Set(equipmentport.DefinitionColumn, eid).
				Where(sql.InInts(equipmentport.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if epuo.clearedParent {
		query, args := builder.Update(equipmentport.ParentTable).
			SetNull(equipmentport.ParentColumn).
			Where(sql.InInts(equipment.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epuo.parent) > 0 {
		for eid := range epuo.parent {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentport.ParentTable).
				Set(equipmentport.ParentColumn, eid).
				Where(sql.InInts(equipmentport.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if epuo.clearedLink {
		query, args := builder.Update(equipmentport.LinkTable).
			SetNull(equipmentport.LinkColumn).
			Where(sql.InInts(link.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epuo.link) > 0 {
		for eid := range epuo.link {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmentport.LinkTable).
				Set(equipmentport.LinkColumn, eid).
				Where(sql.InInts(equipmentport.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(epuo.removedProperties) > 0 {
		eids := make([]int, len(epuo.removedProperties))
		for eid := range epuo.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentport.PropertiesTable).
			SetNull(equipmentport.PropertiesColumn).
			Where(sql.InInts(equipmentport.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epuo.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range epuo.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(equipmentport.PropertiesTable).
				Set(equipmentport.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentport.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(epuo.properties) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"EquipmentPort\"", keys(epuo.properties))})
			}
		}
	}
	if len(epuo.removedEndpoints) > 0 {
		eids := make([]int, len(epuo.removedEndpoints))
		for eid := range epuo.removedEndpoints {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentport.EndpointsTable).
			SetNull(equipmentport.EndpointsColumn).
			Where(sql.InInts(equipmentport.EndpointsColumn, ids...)).
			Where(sql.InInts(serviceendpoint.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(epuo.endpoints) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range epuo.endpoints {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(serviceendpoint.FieldID, eid)
			}
			query, args := builder.Update(equipmentport.EndpointsTable).
				Set(equipmentport.EndpointsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentport.EndpointsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(epuo.endpoints) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"endpoints\" %v already connected to a different \"EquipmentPort\"", keys(epuo.endpoints))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ep, nil
}
