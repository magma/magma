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
	"github.com/facebookincubator/symphony/graph/ent/property"
)

// EquipmentPortCreate is the builder for creating a EquipmentPort entity.
type EquipmentPortCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	definition  map[string]struct{}
	parent      map[string]struct{}
	link        map[string]struct{}
	properties  map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (epc *EquipmentPortCreate) SetCreateTime(t time.Time) *EquipmentPortCreate {
	epc.create_time = &t
	return epc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableCreateTime(t *time.Time) *EquipmentPortCreate {
	if t != nil {
		epc.SetCreateTime(*t)
	}
	return epc
}

// SetUpdateTime sets the update_time field.
func (epc *EquipmentPortCreate) SetUpdateTime(t time.Time) *EquipmentPortCreate {
	epc.update_time = &t
	return epc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableUpdateTime(t *time.Time) *EquipmentPortCreate {
	if t != nil {
		epc.SetUpdateTime(*t)
	}
	return epc
}

// SetDefinitionID sets the definition edge to EquipmentPortDefinition by id.
func (epc *EquipmentPortCreate) SetDefinitionID(id string) *EquipmentPortCreate {
	if epc.definition == nil {
		epc.definition = make(map[string]struct{})
	}
	epc.definition[id] = struct{}{}
	return epc
}

// SetDefinition sets the definition edge to EquipmentPortDefinition.
func (epc *EquipmentPortCreate) SetDefinition(e *EquipmentPortDefinition) *EquipmentPortCreate {
	return epc.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epc *EquipmentPortCreate) SetParentID(id string) *EquipmentPortCreate {
	if epc.parent == nil {
		epc.parent = make(map[string]struct{})
	}
	epc.parent[id] = struct{}{}
	return epc
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableParentID(id *string) *EquipmentPortCreate {
	if id != nil {
		epc = epc.SetParentID(*id)
	}
	return epc
}

// SetParent sets the parent edge to Equipment.
func (epc *EquipmentPortCreate) SetParent(e *Equipment) *EquipmentPortCreate {
	return epc.SetParentID(e.ID)
}

// SetLinkID sets the link edge to Link by id.
func (epc *EquipmentPortCreate) SetLinkID(id string) *EquipmentPortCreate {
	if epc.link == nil {
		epc.link = make(map[string]struct{})
	}
	epc.link[id] = struct{}{}
	return epc
}

// SetNillableLinkID sets the link edge to Link by id if the given value is not nil.
func (epc *EquipmentPortCreate) SetNillableLinkID(id *string) *EquipmentPortCreate {
	if id != nil {
		epc = epc.SetLinkID(*id)
	}
	return epc
}

// SetLink sets the link edge to Link.
func (epc *EquipmentPortCreate) SetLink(l *Link) *EquipmentPortCreate {
	return epc.SetLinkID(l.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (epc *EquipmentPortCreate) AddPropertyIDs(ids ...string) *EquipmentPortCreate {
	if epc.properties == nil {
		epc.properties = make(map[string]struct{})
	}
	for i := range ids {
		epc.properties[ids[i]] = struct{}{}
	}
	return epc
}

// AddProperties adds the properties edges to Property.
func (epc *EquipmentPortCreate) AddProperties(p ...*Property) *EquipmentPortCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return epc.AddPropertyIDs(ids...)
}

// Save creates the EquipmentPort in the database.
func (epc *EquipmentPortCreate) Save(ctx context.Context) (*EquipmentPort, error) {
	if epc.create_time == nil {
		v := equipmentport.DefaultCreateTime()
		epc.create_time = &v
	}
	if epc.update_time == nil {
		v := equipmentport.DefaultUpdateTime()
		epc.update_time = &v
	}
	if len(epc.definition) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"definition\"")
	}
	if epc.definition == nil {
		return nil, errors.New("ent: missing required edge \"definition\"")
	}
	if len(epc.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	if len(epc.link) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"link\"")
	}
	return epc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (epc *EquipmentPortCreate) SaveX(ctx context.Context) *EquipmentPort {
	v, err := epc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epc *EquipmentPortCreate) sqlSave(ctx context.Context) (*EquipmentPort, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(epc.driver.Dialect())
		ep      = &EquipmentPort{config: epc.config}
	)
	tx, err := epc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(equipmentport.Table).Default()
	if value := epc.create_time; value != nil {
		insert.Set(equipmentport.FieldCreateTime, *value)
		ep.CreateTime = *value
	}
	if value := epc.update_time; value != nil {
		insert.Set(equipmentport.FieldUpdateTime, *value)
		ep.UpdateTime = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(equipmentport.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	ep.ID = strconv.FormatInt(id, 10)
	if len(epc.definition) > 0 {
		for eid := range epc.definition {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipmentport.DefinitionTable).
				Set(equipmentport.DefinitionColumn, eid).
				Where(sql.EQ(equipmentport.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(epc.parent) > 0 {
		for eid := range epc.parent {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipmentport.ParentTable).
				Set(equipmentport.ParentColumn, eid).
				Where(sql.EQ(equipmentport.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(epc.link) > 0 {
		for eid := range epc.link {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipmentport.LinkTable).
				Set(equipmentport.LinkColumn, eid).
				Where(sql.EQ(equipmentport.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(epc.properties) > 0 {
		p := sql.P()
		for eid := range epc.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(epc.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"EquipmentPort\"", keys(epc.properties))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ep, nil
}
