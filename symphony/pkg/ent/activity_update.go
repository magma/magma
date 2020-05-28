// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/activity"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

// ActivityUpdate is the builder for updating Activity entities.
type ActivityUpdate struct {
	config
	hooks      []Hook
	mutation   *ActivityMutation
	predicates []predicate.Activity
}

// Where adds a new predicate for the builder.
func (au *ActivityUpdate) Where(ps ...predicate.Activity) *ActivityUpdate {
	au.predicates = append(au.predicates, ps...)
	return au
}

// SetChangedField sets the changed_field field.
func (au *ActivityUpdate) SetChangedField(af activity.ChangedField) *ActivityUpdate {
	au.mutation.SetChangedField(af)
	return au
}

// SetIsCreate sets the is_create field.
func (au *ActivityUpdate) SetIsCreate(b bool) *ActivityUpdate {
	au.mutation.SetIsCreate(b)
	return au
}

// SetNillableIsCreate sets the is_create field if the given value is not nil.
func (au *ActivityUpdate) SetNillableIsCreate(b *bool) *ActivityUpdate {
	if b != nil {
		au.SetIsCreate(*b)
	}
	return au
}

// SetOldValue sets the old_value field.
func (au *ActivityUpdate) SetOldValue(s string) *ActivityUpdate {
	au.mutation.SetOldValue(s)
	return au
}

// SetNillableOldValue sets the old_value field if the given value is not nil.
func (au *ActivityUpdate) SetNillableOldValue(s *string) *ActivityUpdate {
	if s != nil {
		au.SetOldValue(*s)
	}
	return au
}

// ClearOldValue clears the value of old_value.
func (au *ActivityUpdate) ClearOldValue() *ActivityUpdate {
	au.mutation.ClearOldValue()
	return au
}

// SetNewValue sets the new_value field.
func (au *ActivityUpdate) SetNewValue(s string) *ActivityUpdate {
	au.mutation.SetNewValue(s)
	return au
}

// SetNillableNewValue sets the new_value field if the given value is not nil.
func (au *ActivityUpdate) SetNillableNewValue(s *string) *ActivityUpdate {
	if s != nil {
		au.SetNewValue(*s)
	}
	return au
}

// ClearNewValue clears the value of new_value.
func (au *ActivityUpdate) ClearNewValue() *ActivityUpdate {
	au.mutation.ClearNewValue()
	return au
}

// SetAuthorID sets the author edge to User by id.
func (au *ActivityUpdate) SetAuthorID(id int) *ActivityUpdate {
	au.mutation.SetAuthorID(id)
	return au
}

// SetNillableAuthorID sets the author edge to User by id if the given value is not nil.
func (au *ActivityUpdate) SetNillableAuthorID(id *int) *ActivityUpdate {
	if id != nil {
		au = au.SetAuthorID(*id)
	}
	return au
}

// SetAuthor sets the author edge to User.
func (au *ActivityUpdate) SetAuthor(u *User) *ActivityUpdate {
	return au.SetAuthorID(u.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (au *ActivityUpdate) SetWorkOrderID(id int) *ActivityUpdate {
	au.mutation.SetWorkOrderID(id)
	return au
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (au *ActivityUpdate) SetNillableWorkOrderID(id *int) *ActivityUpdate {
	if id != nil {
		au = au.SetWorkOrderID(*id)
	}
	return au
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (au *ActivityUpdate) SetWorkOrder(w *WorkOrder) *ActivityUpdate {
	return au.SetWorkOrderID(w.ID)
}

// ClearAuthor clears the author edge to User.
func (au *ActivityUpdate) ClearAuthor() *ActivityUpdate {
	au.mutation.ClearAuthor()
	return au
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (au *ActivityUpdate) ClearWorkOrder() *ActivityUpdate {
	au.mutation.ClearWorkOrder()
	return au
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (au *ActivityUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := au.mutation.UpdateTime(); !ok {
		v := activity.UpdateDefaultUpdateTime()
		au.mutation.SetUpdateTime(v)
	}
	if v, ok := au.mutation.ChangedField(); ok {
		if err := activity.ChangedFieldValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"changed_field\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(au.hooks) == 0 {
		affected, err = au.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ActivityMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			au.mutation = mutation
			affected, err = au.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(au.hooks) - 1; i >= 0; i-- {
			mut = au.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, au.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (au *ActivityUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *ActivityUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *ActivityUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

func (au *ActivityUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   activity.Table,
			Columns: activity.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: activity.FieldID,
			},
		},
	}
	if ps := au.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: activity.FieldUpdateTime,
		})
	}
	if value, ok := au.mutation.ChangedField(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: activity.FieldChangedField,
		})
	}
	if value, ok := au.mutation.IsCreate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: activity.FieldIsCreate,
		})
	}
	if value, ok := au.mutation.OldValue(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: activity.FieldOldValue,
		})
	}
	if au.mutation.OldValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: activity.FieldOldValue,
		})
	}
	if value, ok := au.mutation.NewValue(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: activity.FieldNewValue,
		})
	}
	if au.mutation.NewValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: activity.FieldNewValue,
		})
	}
	if au.mutation.AuthorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   activity.AuthorTable,
			Columns: []string{activity.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.AuthorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   activity.AuthorTable,
			Columns: []string{activity.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if au.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   activity.WorkOrderTable,
			Columns: []string{activity.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   activity.WorkOrderTable,
			Columns: []string{activity.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{activity.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ActivityUpdateOne is the builder for updating a single Activity entity.
type ActivityUpdateOne struct {
	config
	hooks    []Hook
	mutation *ActivityMutation
}

// SetChangedField sets the changed_field field.
func (auo *ActivityUpdateOne) SetChangedField(af activity.ChangedField) *ActivityUpdateOne {
	auo.mutation.SetChangedField(af)
	return auo
}

// SetIsCreate sets the is_create field.
func (auo *ActivityUpdateOne) SetIsCreate(b bool) *ActivityUpdateOne {
	auo.mutation.SetIsCreate(b)
	return auo
}

// SetNillableIsCreate sets the is_create field if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableIsCreate(b *bool) *ActivityUpdateOne {
	if b != nil {
		auo.SetIsCreate(*b)
	}
	return auo
}

// SetOldValue sets the old_value field.
func (auo *ActivityUpdateOne) SetOldValue(s string) *ActivityUpdateOne {
	auo.mutation.SetOldValue(s)
	return auo
}

// SetNillableOldValue sets the old_value field if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableOldValue(s *string) *ActivityUpdateOne {
	if s != nil {
		auo.SetOldValue(*s)
	}
	return auo
}

// ClearOldValue clears the value of old_value.
func (auo *ActivityUpdateOne) ClearOldValue() *ActivityUpdateOne {
	auo.mutation.ClearOldValue()
	return auo
}

// SetNewValue sets the new_value field.
func (auo *ActivityUpdateOne) SetNewValue(s string) *ActivityUpdateOne {
	auo.mutation.SetNewValue(s)
	return auo
}

// SetNillableNewValue sets the new_value field if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableNewValue(s *string) *ActivityUpdateOne {
	if s != nil {
		auo.SetNewValue(*s)
	}
	return auo
}

// ClearNewValue clears the value of new_value.
func (auo *ActivityUpdateOne) ClearNewValue() *ActivityUpdateOne {
	auo.mutation.ClearNewValue()
	return auo
}

// SetAuthorID sets the author edge to User by id.
func (auo *ActivityUpdateOne) SetAuthorID(id int) *ActivityUpdateOne {
	auo.mutation.SetAuthorID(id)
	return auo
}

// SetNillableAuthorID sets the author edge to User by id if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableAuthorID(id *int) *ActivityUpdateOne {
	if id != nil {
		auo = auo.SetAuthorID(*id)
	}
	return auo
}

// SetAuthor sets the author edge to User.
func (auo *ActivityUpdateOne) SetAuthor(u *User) *ActivityUpdateOne {
	return auo.SetAuthorID(u.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (auo *ActivityUpdateOne) SetWorkOrderID(id int) *ActivityUpdateOne {
	auo.mutation.SetWorkOrderID(id)
	return auo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableWorkOrderID(id *int) *ActivityUpdateOne {
	if id != nil {
		auo = auo.SetWorkOrderID(*id)
	}
	return auo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (auo *ActivityUpdateOne) SetWorkOrder(w *WorkOrder) *ActivityUpdateOne {
	return auo.SetWorkOrderID(w.ID)
}

// ClearAuthor clears the author edge to User.
func (auo *ActivityUpdateOne) ClearAuthor() *ActivityUpdateOne {
	auo.mutation.ClearAuthor()
	return auo
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (auo *ActivityUpdateOne) ClearWorkOrder() *ActivityUpdateOne {
	auo.mutation.ClearWorkOrder()
	return auo
}

// Save executes the query and returns the updated entity.
func (auo *ActivityUpdateOne) Save(ctx context.Context) (*Activity, error) {
	if _, ok := auo.mutation.UpdateTime(); !ok {
		v := activity.UpdateDefaultUpdateTime()
		auo.mutation.SetUpdateTime(v)
	}
	if v, ok := auo.mutation.ChangedField(); ok {
		if err := activity.ChangedFieldValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"changed_field\": %v", err)
		}
	}

	var (
		err  error
		node *Activity
	)
	if len(auo.hooks) == 0 {
		node, err = auo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ActivityMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			auo.mutation = mutation
			node, err = auo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(auo.hooks) - 1; i >= 0; i-- {
			mut = auo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, auo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (auo *ActivityUpdateOne) SaveX(ctx context.Context) *Activity {
	a, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return a
}

// Exec executes the query on the entity.
func (auo *ActivityUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *ActivityUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (auo *ActivityUpdateOne) sqlSave(ctx context.Context) (a *Activity, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   activity.Table,
			Columns: activity.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: activity.FieldID,
			},
		},
	}
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Activity.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := auo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: activity.FieldUpdateTime,
		})
	}
	if value, ok := auo.mutation.ChangedField(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: activity.FieldChangedField,
		})
	}
	if value, ok := auo.mutation.IsCreate(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: activity.FieldIsCreate,
		})
	}
	if value, ok := auo.mutation.OldValue(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: activity.FieldOldValue,
		})
	}
	if auo.mutation.OldValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: activity.FieldOldValue,
		})
	}
	if value, ok := auo.mutation.NewValue(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: activity.FieldNewValue,
		})
	}
	if auo.mutation.NewValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: activity.FieldNewValue,
		})
	}
	if auo.mutation.AuthorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   activity.AuthorTable,
			Columns: []string{activity.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.AuthorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   activity.AuthorTable,
			Columns: []string{activity.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if auo.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   activity.WorkOrderTable,
			Columns: []string{activity.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   activity.WorkOrderTable,
			Columns: []string{activity.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	a = &Activity{config: auo.config}
	_spec.Assign = a.assignValues
	_spec.ScanValues = a.scanValues()
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{activity.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return a, nil
}
