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
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
)

// EquipmentCategoryCreate is the builder for creating a EquipmentCategory entity.
type EquipmentCategoryCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	types       map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (ecc *EquipmentCategoryCreate) SetCreateTime(t time.Time) *EquipmentCategoryCreate {
	ecc.create_time = &t
	return ecc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ecc *EquipmentCategoryCreate) SetNillableCreateTime(t *time.Time) *EquipmentCategoryCreate {
	if t != nil {
		ecc.SetCreateTime(*t)
	}
	return ecc
}

// SetUpdateTime sets the update_time field.
func (ecc *EquipmentCategoryCreate) SetUpdateTime(t time.Time) *EquipmentCategoryCreate {
	ecc.update_time = &t
	return ecc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ecc *EquipmentCategoryCreate) SetNillableUpdateTime(t *time.Time) *EquipmentCategoryCreate {
	if t != nil {
		ecc.SetUpdateTime(*t)
	}
	return ecc
}

// SetName sets the name field.
func (ecc *EquipmentCategoryCreate) SetName(s string) *EquipmentCategoryCreate {
	ecc.name = &s
	return ecc
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (ecc *EquipmentCategoryCreate) AddTypeIDs(ids ...string) *EquipmentCategoryCreate {
	if ecc.types == nil {
		ecc.types = make(map[string]struct{})
	}
	for i := range ids {
		ecc.types[ids[i]] = struct{}{}
	}
	return ecc
}

// AddTypes adds the types edges to EquipmentType.
func (ecc *EquipmentCategoryCreate) AddTypes(e ...*EquipmentType) *EquipmentCategoryCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecc.AddTypeIDs(ids...)
}

// Save creates the EquipmentCategory in the database.
func (ecc *EquipmentCategoryCreate) Save(ctx context.Context) (*EquipmentCategory, error) {
	if ecc.create_time == nil {
		v := equipmentcategory.DefaultCreateTime()
		ecc.create_time = &v
	}
	if ecc.update_time == nil {
		v := equipmentcategory.DefaultUpdateTime()
		ecc.update_time = &v
	}
	if ecc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	return ecc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (ecc *EquipmentCategoryCreate) SaveX(ctx context.Context) *EquipmentCategory {
	v, err := ecc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ecc *EquipmentCategoryCreate) sqlSave(ctx context.Context) (*EquipmentCategory, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ecc.driver.Dialect())
		ec      = &EquipmentCategory{config: ecc.config}
	)
	tx, err := ecc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(equipmentcategory.Table).Default()
	if value := ecc.create_time; value != nil {
		insert.Set(equipmentcategory.FieldCreateTime, *value)
		ec.CreateTime = *value
	}
	if value := ecc.update_time; value != nil {
		insert.Set(equipmentcategory.FieldUpdateTime, *value)
		ec.UpdateTime = *value
	}
	if value := ecc.name; value != nil {
		insert.Set(equipmentcategory.FieldName, *value)
		ec.Name = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(equipmentcategory.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	ec.ID = strconv.FormatInt(id, 10)
	if len(ecc.types) > 0 {
		p := sql.P()
		for eid := range ecc.types {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(ecc.types) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"types\" %v already connected to a different \"EquipmentCategory\"", keys(ecc.types))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ec, nil
}
