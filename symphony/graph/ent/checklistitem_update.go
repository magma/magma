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

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// CheckListItemUpdate is the builder for updating CheckListItem entities.
type CheckListItemUpdate struct {
	config
	title            *string
	_type            *string
	index            *int
	addindex         *int
	clearindex       bool
	checked          *bool
	clearchecked     bool
	string_val       *string
	clearstring_val  bool
	enum_values      *string
	clearenum_values bool
	help_text        *string
	clearhelp_text   bool
	work_order       map[string]struct{}
	clearedWorkOrder bool
	predicates       []predicate.CheckListItem
}

// Where adds a new predicate for the builder.
func (cliu *CheckListItemUpdate) Where(ps ...predicate.CheckListItem) *CheckListItemUpdate {
	cliu.predicates = append(cliu.predicates, ps...)
	return cliu
}

// SetTitle sets the title field.
func (cliu *CheckListItemUpdate) SetTitle(s string) *CheckListItemUpdate {
	cliu.title = &s
	return cliu
}

// SetType sets the type field.
func (cliu *CheckListItemUpdate) SetType(s string) *CheckListItemUpdate {
	cliu._type = &s
	return cliu
}

// SetIndex sets the index field.
func (cliu *CheckListItemUpdate) SetIndex(i int) *CheckListItemUpdate {
	cliu.index = &i
	cliu.addindex = nil
	return cliu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableIndex(i *int) *CheckListItemUpdate {
	if i != nil {
		cliu.SetIndex(*i)
	}
	return cliu
}

// AddIndex adds i to index.
func (cliu *CheckListItemUpdate) AddIndex(i int) *CheckListItemUpdate {
	if cliu.addindex == nil {
		cliu.addindex = &i
	} else {
		*cliu.addindex += i
	}
	return cliu
}

// ClearIndex clears the value of index.
func (cliu *CheckListItemUpdate) ClearIndex() *CheckListItemUpdate {
	cliu.index = nil
	cliu.clearindex = true
	return cliu
}

// SetChecked sets the checked field.
func (cliu *CheckListItemUpdate) SetChecked(b bool) *CheckListItemUpdate {
	cliu.checked = &b
	return cliu
}

// SetNillableChecked sets the checked field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableChecked(b *bool) *CheckListItemUpdate {
	if b != nil {
		cliu.SetChecked(*b)
	}
	return cliu
}

// ClearChecked clears the value of checked.
func (cliu *CheckListItemUpdate) ClearChecked() *CheckListItemUpdate {
	cliu.checked = nil
	cliu.clearchecked = true
	return cliu
}

// SetStringVal sets the string_val field.
func (cliu *CheckListItemUpdate) SetStringVal(s string) *CheckListItemUpdate {
	cliu.string_val = &s
	return cliu
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableStringVal(s *string) *CheckListItemUpdate {
	if s != nil {
		cliu.SetStringVal(*s)
	}
	return cliu
}

// ClearStringVal clears the value of string_val.
func (cliu *CheckListItemUpdate) ClearStringVal() *CheckListItemUpdate {
	cliu.string_val = nil
	cliu.clearstring_val = true
	return cliu
}

// SetEnumValues sets the enum_values field.
func (cliu *CheckListItemUpdate) SetEnumValues(s string) *CheckListItemUpdate {
	cliu.enum_values = &s
	return cliu
}

// SetNillableEnumValues sets the enum_values field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableEnumValues(s *string) *CheckListItemUpdate {
	if s != nil {
		cliu.SetEnumValues(*s)
	}
	return cliu
}

// ClearEnumValues clears the value of enum_values.
func (cliu *CheckListItemUpdate) ClearEnumValues() *CheckListItemUpdate {
	cliu.enum_values = nil
	cliu.clearenum_values = true
	return cliu
}

// SetHelpText sets the help_text field.
func (cliu *CheckListItemUpdate) SetHelpText(s string) *CheckListItemUpdate {
	cliu.help_text = &s
	return cliu
}

// SetNillableHelpText sets the help_text field if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableHelpText(s *string) *CheckListItemUpdate {
	if s != nil {
		cliu.SetHelpText(*s)
	}
	return cliu
}

// ClearHelpText clears the value of help_text.
func (cliu *CheckListItemUpdate) ClearHelpText() *CheckListItemUpdate {
	cliu.help_text = nil
	cliu.clearhelp_text = true
	return cliu
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (cliu *CheckListItemUpdate) SetWorkOrderID(id string) *CheckListItemUpdate {
	if cliu.work_order == nil {
		cliu.work_order = make(map[string]struct{})
	}
	cliu.work_order[id] = struct{}{}
	return cliu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (cliu *CheckListItemUpdate) SetNillableWorkOrderID(id *string) *CheckListItemUpdate {
	if id != nil {
		cliu = cliu.SetWorkOrderID(*id)
	}
	return cliu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (cliu *CheckListItemUpdate) SetWorkOrder(w *WorkOrder) *CheckListItemUpdate {
	return cliu.SetWorkOrderID(w.ID)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (cliu *CheckListItemUpdate) ClearWorkOrder() *CheckListItemUpdate {
	cliu.clearedWorkOrder = true
	return cliu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (cliu *CheckListItemUpdate) Save(ctx context.Context) (int, error) {
	if len(cliu.work_order) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return cliu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cliu *CheckListItemUpdate) SaveX(ctx context.Context) int {
	affected, err := cliu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cliu *CheckListItemUpdate) Exec(ctx context.Context) error {
	_, err := cliu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cliu *CheckListItemUpdate) ExecX(ctx context.Context) {
	if err := cliu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cliu *CheckListItemUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(cliu.driver.Dialect())
		selector = builder.Select(checklistitem.FieldID).From(builder.Table(checklistitem.Table))
	)
	for _, p := range cliu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cliu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := cliu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(checklistitem.Table).Where(sql.InInts(checklistitem.FieldID, ids...))
	)
	if value := cliu.title; value != nil {
		updater.Set(checklistitem.FieldTitle, *value)
	}
	if value := cliu._type; value != nil {
		updater.Set(checklistitem.FieldType, *value)
	}
	if value := cliu.index; value != nil {
		updater.Set(checklistitem.FieldIndex, *value)
	}
	if value := cliu.addindex; value != nil {
		updater.Add(checklistitem.FieldIndex, *value)
	}
	if cliu.clearindex {
		updater.SetNull(checklistitem.FieldIndex)
	}
	if value := cliu.checked; value != nil {
		updater.Set(checklistitem.FieldChecked, *value)
	}
	if cliu.clearchecked {
		updater.SetNull(checklistitem.FieldChecked)
	}
	if value := cliu.string_val; value != nil {
		updater.Set(checklistitem.FieldStringVal, *value)
	}
	if cliu.clearstring_val {
		updater.SetNull(checklistitem.FieldStringVal)
	}
	if value := cliu.enum_values; value != nil {
		updater.Set(checklistitem.FieldEnumValues, *value)
	}
	if cliu.clearenum_values {
		updater.SetNull(checklistitem.FieldEnumValues)
	}
	if value := cliu.help_text; value != nil {
		updater.Set(checklistitem.FieldHelpText, *value)
	}
	if cliu.clearhelp_text {
		updater.SetNull(checklistitem.FieldHelpText)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if cliu.clearedWorkOrder {
		query, args := builder.Update(checklistitem.WorkOrderTable).
			SetNull(checklistitem.WorkOrderColumn).
			Where(sql.InInts(workorder.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(cliu.work_order) > 0 {
		for eid := range cliu.work_order {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(checklistitem.WorkOrderTable).
				Set(checklistitem.WorkOrderColumn, eid).
				Where(sql.InInts(checklistitem.FieldID, ids...)).
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

// CheckListItemUpdateOne is the builder for updating a single CheckListItem entity.
type CheckListItemUpdateOne struct {
	config
	id               string
	title            *string
	_type            *string
	index            *int
	addindex         *int
	clearindex       bool
	checked          *bool
	clearchecked     bool
	string_val       *string
	clearstring_val  bool
	enum_values      *string
	clearenum_values bool
	help_text        *string
	clearhelp_text   bool
	work_order       map[string]struct{}
	clearedWorkOrder bool
}

// SetTitle sets the title field.
func (cliuo *CheckListItemUpdateOne) SetTitle(s string) *CheckListItemUpdateOne {
	cliuo.title = &s
	return cliuo
}

// SetType sets the type field.
func (cliuo *CheckListItemUpdateOne) SetType(s string) *CheckListItemUpdateOne {
	cliuo._type = &s
	return cliuo
}

// SetIndex sets the index field.
func (cliuo *CheckListItemUpdateOne) SetIndex(i int) *CheckListItemUpdateOne {
	cliuo.index = &i
	cliuo.addindex = nil
	return cliuo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableIndex(i *int) *CheckListItemUpdateOne {
	if i != nil {
		cliuo.SetIndex(*i)
	}
	return cliuo
}

// AddIndex adds i to index.
func (cliuo *CheckListItemUpdateOne) AddIndex(i int) *CheckListItemUpdateOne {
	if cliuo.addindex == nil {
		cliuo.addindex = &i
	} else {
		*cliuo.addindex += i
	}
	return cliuo
}

// ClearIndex clears the value of index.
func (cliuo *CheckListItemUpdateOne) ClearIndex() *CheckListItemUpdateOne {
	cliuo.index = nil
	cliuo.clearindex = true
	return cliuo
}

// SetChecked sets the checked field.
func (cliuo *CheckListItemUpdateOne) SetChecked(b bool) *CheckListItemUpdateOne {
	cliuo.checked = &b
	return cliuo
}

// SetNillableChecked sets the checked field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableChecked(b *bool) *CheckListItemUpdateOne {
	if b != nil {
		cliuo.SetChecked(*b)
	}
	return cliuo
}

// ClearChecked clears the value of checked.
func (cliuo *CheckListItemUpdateOne) ClearChecked() *CheckListItemUpdateOne {
	cliuo.checked = nil
	cliuo.clearchecked = true
	return cliuo
}

// SetStringVal sets the string_val field.
func (cliuo *CheckListItemUpdateOne) SetStringVal(s string) *CheckListItemUpdateOne {
	cliuo.string_val = &s
	return cliuo
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableStringVal(s *string) *CheckListItemUpdateOne {
	if s != nil {
		cliuo.SetStringVal(*s)
	}
	return cliuo
}

// ClearStringVal clears the value of string_val.
func (cliuo *CheckListItemUpdateOne) ClearStringVal() *CheckListItemUpdateOne {
	cliuo.string_val = nil
	cliuo.clearstring_val = true
	return cliuo
}

// SetEnumValues sets the enum_values field.
func (cliuo *CheckListItemUpdateOne) SetEnumValues(s string) *CheckListItemUpdateOne {
	cliuo.enum_values = &s
	return cliuo
}

// SetNillableEnumValues sets the enum_values field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableEnumValues(s *string) *CheckListItemUpdateOne {
	if s != nil {
		cliuo.SetEnumValues(*s)
	}
	return cliuo
}

// ClearEnumValues clears the value of enum_values.
func (cliuo *CheckListItemUpdateOne) ClearEnumValues() *CheckListItemUpdateOne {
	cliuo.enum_values = nil
	cliuo.clearenum_values = true
	return cliuo
}

// SetHelpText sets the help_text field.
func (cliuo *CheckListItemUpdateOne) SetHelpText(s string) *CheckListItemUpdateOne {
	cliuo.help_text = &s
	return cliuo
}

// SetNillableHelpText sets the help_text field if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableHelpText(s *string) *CheckListItemUpdateOne {
	if s != nil {
		cliuo.SetHelpText(*s)
	}
	return cliuo
}

// ClearHelpText clears the value of help_text.
func (cliuo *CheckListItemUpdateOne) ClearHelpText() *CheckListItemUpdateOne {
	cliuo.help_text = nil
	cliuo.clearhelp_text = true
	return cliuo
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (cliuo *CheckListItemUpdateOne) SetWorkOrderID(id string) *CheckListItemUpdateOne {
	if cliuo.work_order == nil {
		cliuo.work_order = make(map[string]struct{})
	}
	cliuo.work_order[id] = struct{}{}
	return cliuo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (cliuo *CheckListItemUpdateOne) SetNillableWorkOrderID(id *string) *CheckListItemUpdateOne {
	if id != nil {
		cliuo = cliuo.SetWorkOrderID(*id)
	}
	return cliuo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (cliuo *CheckListItemUpdateOne) SetWorkOrder(w *WorkOrder) *CheckListItemUpdateOne {
	return cliuo.SetWorkOrderID(w.ID)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (cliuo *CheckListItemUpdateOne) ClearWorkOrder() *CheckListItemUpdateOne {
	cliuo.clearedWorkOrder = true
	return cliuo
}

// Save executes the query and returns the updated entity.
func (cliuo *CheckListItemUpdateOne) Save(ctx context.Context) (*CheckListItem, error) {
	if len(cliuo.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return cliuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cliuo *CheckListItemUpdateOne) SaveX(ctx context.Context) *CheckListItem {
	cli, err := cliuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return cli
}

// Exec executes the query on the entity.
func (cliuo *CheckListItemUpdateOne) Exec(ctx context.Context) error {
	_, err := cliuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cliuo *CheckListItemUpdateOne) ExecX(ctx context.Context) {
	if err := cliuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cliuo *CheckListItemUpdateOne) sqlSave(ctx context.Context) (cli *CheckListItem, err error) {
	var (
		builder  = sql.Dialect(cliuo.driver.Dialect())
		selector = builder.Select(checklistitem.Columns...).From(builder.Table(checklistitem.Table))
	)
	checklistitem.ID(cliuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = cliuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		cli = &CheckListItem{config: cliuo.config}
		if err := cli.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into CheckListItem: %v", err)
		}
		id = cli.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("CheckListItem with id: %v", cliuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one CheckListItem with the same id: %v", cliuo.id)
	}

	tx, err := cliuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(checklistitem.Table).Where(sql.InInts(checklistitem.FieldID, ids...))
	)
	if value := cliuo.title; value != nil {
		updater.Set(checklistitem.FieldTitle, *value)
		cli.Title = *value
	}
	if value := cliuo._type; value != nil {
		updater.Set(checklistitem.FieldType, *value)
		cli.Type = *value
	}
	if value := cliuo.index; value != nil {
		updater.Set(checklistitem.FieldIndex, *value)
		cli.Index = *value
	}
	if value := cliuo.addindex; value != nil {
		updater.Add(checklistitem.FieldIndex, *value)
		cli.Index += *value
	}
	if cliuo.clearindex {
		var value int
		cli.Index = value
		updater.SetNull(checklistitem.FieldIndex)
	}
	if value := cliuo.checked; value != nil {
		updater.Set(checklistitem.FieldChecked, *value)
		cli.Checked = *value
	}
	if cliuo.clearchecked {
		var value bool
		cli.Checked = value
		updater.SetNull(checklistitem.FieldChecked)
	}
	if value := cliuo.string_val; value != nil {
		updater.Set(checklistitem.FieldStringVal, *value)
		cli.StringVal = *value
	}
	if cliuo.clearstring_val {
		var value string
		cli.StringVal = value
		updater.SetNull(checklistitem.FieldStringVal)
	}
	if value := cliuo.enum_values; value != nil {
		updater.Set(checklistitem.FieldEnumValues, *value)
		cli.EnumValues = *value
	}
	if cliuo.clearenum_values {
		var value string
		cli.EnumValues = value
		updater.SetNull(checklistitem.FieldEnumValues)
	}
	if value := cliuo.help_text; value != nil {
		updater.Set(checklistitem.FieldHelpText, *value)
		cli.HelpText = value
	}
	if cliuo.clearhelp_text {
		cli.HelpText = nil
		updater.SetNull(checklistitem.FieldHelpText)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if cliuo.clearedWorkOrder {
		query, args := builder.Update(checklistitem.WorkOrderTable).
			SetNull(checklistitem.WorkOrderColumn).
			Where(sql.InInts(workorder.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(cliuo.work_order) > 0 {
		for eid := range cliuo.work_order {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(checklistitem.WorkOrderTable).
				Set(checklistitem.WorkOrderColumn, eid).
				Where(sql.InInts(checklistitem.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return cli, nil
}
