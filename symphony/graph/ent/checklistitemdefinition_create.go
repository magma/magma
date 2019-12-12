// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
)

// CheckListItemDefinitionCreate is the builder for creating a CheckListItemDefinition entity.
type CheckListItemDefinitionCreate struct {
	config
	title           *string
	_type           *string
	index           *int
	enum_values     *string
	help_text       *string
	work_order_type map[string]struct{}
}

// SetTitle sets the title field.
func (clidc *CheckListItemDefinitionCreate) SetTitle(s string) *CheckListItemDefinitionCreate {
	clidc.title = &s
	return clidc
}

// SetType sets the type field.
func (clidc *CheckListItemDefinitionCreate) SetType(s string) *CheckListItemDefinitionCreate {
	clidc._type = &s
	return clidc
}

// SetIndex sets the index field.
func (clidc *CheckListItemDefinitionCreate) SetIndex(i int) *CheckListItemDefinitionCreate {
	clidc.index = &i
	return clidc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (clidc *CheckListItemDefinitionCreate) SetNillableIndex(i *int) *CheckListItemDefinitionCreate {
	if i != nil {
		clidc.SetIndex(*i)
	}
	return clidc
}

// SetEnumValues sets the enum_values field.
func (clidc *CheckListItemDefinitionCreate) SetEnumValues(s string) *CheckListItemDefinitionCreate {
	clidc.enum_values = &s
	return clidc
}

// SetNillableEnumValues sets the enum_values field if the given value is not nil.
func (clidc *CheckListItemDefinitionCreate) SetNillableEnumValues(s *string) *CheckListItemDefinitionCreate {
	if s != nil {
		clidc.SetEnumValues(*s)
	}
	return clidc
}

// SetHelpText sets the help_text field.
func (clidc *CheckListItemDefinitionCreate) SetHelpText(s string) *CheckListItemDefinitionCreate {
	clidc.help_text = &s
	return clidc
}

// SetNillableHelpText sets the help_text field if the given value is not nil.
func (clidc *CheckListItemDefinitionCreate) SetNillableHelpText(s *string) *CheckListItemDefinitionCreate {
	if s != nil {
		clidc.SetHelpText(*s)
	}
	return clidc
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (clidc *CheckListItemDefinitionCreate) SetWorkOrderTypeID(id string) *CheckListItemDefinitionCreate {
	if clidc.work_order_type == nil {
		clidc.work_order_type = make(map[string]struct{})
	}
	clidc.work_order_type[id] = struct{}{}
	return clidc
}

// SetNillableWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id if the given value is not nil.
func (clidc *CheckListItemDefinitionCreate) SetNillableWorkOrderTypeID(id *string) *CheckListItemDefinitionCreate {
	if id != nil {
		clidc = clidc.SetWorkOrderTypeID(*id)
	}
	return clidc
}

// SetWorkOrderType sets the work_order_type edge to WorkOrderType.
func (clidc *CheckListItemDefinitionCreate) SetWorkOrderType(w *WorkOrderType) *CheckListItemDefinitionCreate {
	return clidc.SetWorkOrderTypeID(w.ID)
}

// Save creates the CheckListItemDefinition in the database.
func (clidc *CheckListItemDefinitionCreate) Save(ctx context.Context) (*CheckListItemDefinition, error) {
	if clidc.title == nil {
		return nil, errors.New("ent: missing required field \"title\"")
	}
	if clidc._type == nil {
		return nil, errors.New("ent: missing required field \"type\"")
	}
	if len(clidc.work_order_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order_type\"")
	}
	return clidc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (clidc *CheckListItemDefinitionCreate) SaveX(ctx context.Context) *CheckListItemDefinition {
	v, err := clidc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (clidc *CheckListItemDefinitionCreate) sqlSave(ctx context.Context) (*CheckListItemDefinition, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(clidc.driver.Dialect())
		clid    = &CheckListItemDefinition{config: clidc.config}
	)
	tx, err := clidc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(checklistitemdefinition.Table).Default()
	if value := clidc.title; value != nil {
		insert.Set(checklistitemdefinition.FieldTitle, *value)
		clid.Title = *value
	}
	if value := clidc._type; value != nil {
		insert.Set(checklistitemdefinition.FieldType, *value)
		clid.Type = *value
	}
	if value := clidc.index; value != nil {
		insert.Set(checklistitemdefinition.FieldIndex, *value)
		clid.Index = *value
	}
	if value := clidc.enum_values; value != nil {
		insert.Set(checklistitemdefinition.FieldEnumValues, *value)
		clid.EnumValues = value
	}
	if value := clidc.help_text; value != nil {
		insert.Set(checklistitemdefinition.FieldHelpText, *value)
		clid.HelpText = value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(checklistitemdefinition.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	clid.ID = strconv.FormatInt(id, 10)
	if len(clidc.work_order_type) > 0 {
		for eid := range clidc.work_order_type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(checklistitemdefinition.WorkOrderTypeTable).
				Set(checklistitemdefinition.WorkOrderTypeColumn, eid).
				Where(sql.EQ(checklistitemdefinition.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return clid, nil
}
