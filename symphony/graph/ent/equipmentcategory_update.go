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
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentCategoryUpdate is the builder for updating EquipmentCategory entities.
type EquipmentCategoryUpdate struct {
	config

	update_time  *time.Time
	name         *string
	types        map[string]struct{}
	removedTypes map[string]struct{}
	predicates   []predicate.EquipmentCategory
}

// Where adds a new predicate for the builder.
func (ecu *EquipmentCategoryUpdate) Where(ps ...predicate.EquipmentCategory) *EquipmentCategoryUpdate {
	ecu.predicates = append(ecu.predicates, ps...)
	return ecu
}

// SetName sets the name field.
func (ecu *EquipmentCategoryUpdate) SetName(s string) *EquipmentCategoryUpdate {
	ecu.name = &s
	return ecu
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecu *EquipmentCategoryUpdate) AddTypeIDs(ids ...string) *EquipmentCategoryUpdate {
	if ecu.types == nil {
		ecu.types = make(map[string]struct{})
	}
	for i := range ids {
		ecu.types[ids[i]] = struct{}{}
	}
	return ecu
}

// AddTypes adds the types edges to EquipmentType.
func (ecu *EquipmentCategoryUpdate) AddTypes(e ...*EquipmentType) *EquipmentCategoryUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.AddTypeIDs(ids...)
}

// RemoveTypeIDs removes the types edge to EquipmentType by ids.
func (ecu *EquipmentCategoryUpdate) RemoveTypeIDs(ids ...string) *EquipmentCategoryUpdate {
	if ecu.removedTypes == nil {
		ecu.removedTypes = make(map[string]struct{})
	}
	for i := range ids {
		ecu.removedTypes[ids[i]] = struct{}{}
	}
	return ecu
}

// RemoveTypes removes types edges to EquipmentType.
func (ecu *EquipmentCategoryUpdate) RemoveTypes(e ...*EquipmentType) *EquipmentCategoryUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.RemoveTypeIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ecu *EquipmentCategoryUpdate) Save(ctx context.Context) (int, error) {
	if ecu.update_time == nil {
		v := equipmentcategory.UpdateDefaultUpdateTime()
		ecu.update_time = &v
	}
	return ecu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ecu *EquipmentCategoryUpdate) SaveX(ctx context.Context) int {
	affected, err := ecu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ecu *EquipmentCategoryUpdate) Exec(ctx context.Context) error {
	_, err := ecu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecu *EquipmentCategoryUpdate) ExecX(ctx context.Context) {
	if err := ecu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ecu *EquipmentCategoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(ecu.driver.Dialect())
		selector = builder.Select(equipmentcategory.FieldID).From(builder.Table(equipmentcategory.Table))
	)
	for _, p := range ecu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ecu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := ecu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentcategory.Table).Where(sql.InInts(equipmentcategory.FieldID, ids...))
	)
	if value := ecu.update_time; value != nil {
		updater.Set(equipmentcategory.FieldUpdateTime, *value)
	}
	if value := ecu.name; value != nil {
		updater.Set(equipmentcategory.FieldName, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ecu.removedTypes) > 0 {
		eids := make([]int, len(ecu.removedTypes))
		for eid := range ecu.removedTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentcategory.TypesTable).
			SetNull(equipmentcategory.TypesColumn).
			Where(sql.InInts(equipmentcategory.TypesColumn, ids...)).
			Where(sql.InInts(equipmenttype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(ecu.types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ecu.types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmenttype.FieldID, eid)
			}
			query, args := builder.Update(equipmentcategory.TypesTable).
				Set(equipmentcategory.TypesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentcategory.TypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(ecu.types) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"types\" %v already connected to a different \"EquipmentCategory\"", keys(ecu.types))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// EquipmentCategoryUpdateOne is the builder for updating a single EquipmentCategory entity.
type EquipmentCategoryUpdateOne struct {
	config
	id string

	update_time  *time.Time
	name         *string
	types        map[string]struct{}
	removedTypes map[string]struct{}
}

// SetName sets the name field.
func (ecuo *EquipmentCategoryUpdateOne) SetName(s string) *EquipmentCategoryUpdateOne {
	ecuo.name = &s
	return ecuo
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecuo *EquipmentCategoryUpdateOne) AddTypeIDs(ids ...string) *EquipmentCategoryUpdateOne {
	if ecuo.types == nil {
		ecuo.types = make(map[string]struct{})
	}
	for i := range ids {
		ecuo.types[ids[i]] = struct{}{}
	}
	return ecuo
}

// AddTypes adds the types edges to EquipmentType.
func (ecuo *EquipmentCategoryUpdateOne) AddTypes(e ...*EquipmentType) *EquipmentCategoryUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.AddTypeIDs(ids...)
}

// RemoveTypeIDs removes the types edge to EquipmentType by ids.
func (ecuo *EquipmentCategoryUpdateOne) RemoveTypeIDs(ids ...string) *EquipmentCategoryUpdateOne {
	if ecuo.removedTypes == nil {
		ecuo.removedTypes = make(map[string]struct{})
	}
	for i := range ids {
		ecuo.removedTypes[ids[i]] = struct{}{}
	}
	return ecuo
}

// RemoveTypes removes types edges to EquipmentType.
func (ecuo *EquipmentCategoryUpdateOne) RemoveTypes(e ...*EquipmentType) *EquipmentCategoryUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.RemoveTypeIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (ecuo *EquipmentCategoryUpdateOne) Save(ctx context.Context) (*EquipmentCategory, error) {
	if ecuo.update_time == nil {
		v := equipmentcategory.UpdateDefaultUpdateTime()
		ecuo.update_time = &v
	}
	return ecuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (ecuo *EquipmentCategoryUpdateOne) SaveX(ctx context.Context) *EquipmentCategory {
	ec, err := ecuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ec
}

// Exec executes the query on the entity.
func (ecuo *EquipmentCategoryUpdateOne) Exec(ctx context.Context) error {
	_, err := ecuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecuo *EquipmentCategoryUpdateOne) ExecX(ctx context.Context) {
	if err := ecuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ecuo *EquipmentCategoryUpdateOne) sqlSave(ctx context.Context) (ec *EquipmentCategory, err error) {
	var (
		builder  = sql.Dialect(ecuo.driver.Dialect())
		selector = builder.Select(equipmentcategory.Columns...).From(builder.Table(equipmentcategory.Table))
	)
	equipmentcategory.ID(ecuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = ecuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		ec = &EquipmentCategory{config: ecuo.config}
		if err := ec.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into EquipmentCategory: %v", err)
		}
		id = ec.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("EquipmentCategory with id: %v", ecuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one EquipmentCategory with the same id: %v", ecuo.id)
	}

	tx, err := ecuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentcategory.Table).Where(sql.InInts(equipmentcategory.FieldID, ids...))
	)
	if value := ecuo.update_time; value != nil {
		updater.Set(equipmentcategory.FieldUpdateTime, *value)
		ec.UpdateTime = *value
	}
	if value := ecuo.name; value != nil {
		updater.Set(equipmentcategory.FieldName, *value)
		ec.Name = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ecuo.removedTypes) > 0 {
		eids := make([]int, len(ecuo.removedTypes))
		for eid := range ecuo.removedTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentcategory.TypesTable).
			SetNull(equipmentcategory.TypesColumn).
			Where(sql.InInts(equipmentcategory.TypesColumn, ids...)).
			Where(sql.InInts(equipmenttype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(ecuo.types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range ecuo.types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmenttype.FieldID, eid)
			}
			query, args := builder.Update(equipmentcategory.TypesTable).
				Set(equipmentcategory.TypesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentcategory.TypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(ecuo.types) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"types\" %v already connected to a different \"EquipmentCategory\"", keys(ecuo.types))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ec, nil
}
