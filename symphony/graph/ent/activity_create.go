// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/activity"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// ActivityCreate is the builder for creating a Activity entity.
type ActivityCreate struct {
	config
	mutation *ActivityMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (ac *ActivityCreate) SetCreateTime(t time.Time) *ActivityCreate {
	ac.mutation.SetCreateTime(t)
	return ac
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ac *ActivityCreate) SetNillableCreateTime(t *time.Time) *ActivityCreate {
	if t != nil {
		ac.SetCreateTime(*t)
	}
	return ac
}

// SetUpdateTime sets the update_time field.
func (ac *ActivityCreate) SetUpdateTime(t time.Time) *ActivityCreate {
	ac.mutation.SetUpdateTime(t)
	return ac
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ac *ActivityCreate) SetNillableUpdateTime(t *time.Time) *ActivityCreate {
	if t != nil {
		ac.SetUpdateTime(*t)
	}
	return ac
}

// SetChangedField sets the changed_field field.
func (ac *ActivityCreate) SetChangedField(af activity.ChangedField) *ActivityCreate {
	ac.mutation.SetChangedField(af)
	return ac
}

// SetIsCreate sets the is_create field.
func (ac *ActivityCreate) SetIsCreate(b bool) *ActivityCreate {
	ac.mutation.SetIsCreate(b)
	return ac
}

// SetNillableIsCreate sets the is_create field if the given value is not nil.
func (ac *ActivityCreate) SetNillableIsCreate(b *bool) *ActivityCreate {
	if b != nil {
		ac.SetIsCreate(*b)
	}
	return ac
}

// SetOldValue sets the old_value field.
func (ac *ActivityCreate) SetOldValue(s string) *ActivityCreate {
	ac.mutation.SetOldValue(s)
	return ac
}

// SetNillableOldValue sets the old_value field if the given value is not nil.
func (ac *ActivityCreate) SetNillableOldValue(s *string) *ActivityCreate {
	if s != nil {
		ac.SetOldValue(*s)
	}
	return ac
}

// SetNewValue sets the new_value field.
func (ac *ActivityCreate) SetNewValue(s string) *ActivityCreate {
	ac.mutation.SetNewValue(s)
	return ac
}

// SetNillableNewValue sets the new_value field if the given value is not nil.
func (ac *ActivityCreate) SetNillableNewValue(s *string) *ActivityCreate {
	if s != nil {
		ac.SetNewValue(*s)
	}
	return ac
}

// SetAuthorID sets the author edge to User by id.
func (ac *ActivityCreate) SetAuthorID(id int) *ActivityCreate {
	ac.mutation.SetAuthorID(id)
	return ac
}

// SetNillableAuthorID sets the author edge to User by id if the given value is not nil.
func (ac *ActivityCreate) SetNillableAuthorID(id *int) *ActivityCreate {
	if id != nil {
		ac = ac.SetAuthorID(*id)
	}
	return ac
}

// SetAuthor sets the author edge to User.
func (ac *ActivityCreate) SetAuthor(u *User) *ActivityCreate {
	return ac.SetAuthorID(u.ID)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (ac *ActivityCreate) SetWorkOrderID(id int) *ActivityCreate {
	ac.mutation.SetWorkOrderID(id)
	return ac
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (ac *ActivityCreate) SetNillableWorkOrderID(id *int) *ActivityCreate {
	if id != nil {
		ac = ac.SetWorkOrderID(*id)
	}
	return ac
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (ac *ActivityCreate) SetWorkOrder(w *WorkOrder) *ActivityCreate {
	return ac.SetWorkOrderID(w.ID)
}

// Save creates the Activity in the database.
func (ac *ActivityCreate) Save(ctx context.Context) (*Activity, error) {
	if _, ok := ac.mutation.CreateTime(); !ok {
		v := activity.DefaultCreateTime()
		ac.mutation.SetCreateTime(v)
	}
	if _, ok := ac.mutation.UpdateTime(); !ok {
		v := activity.DefaultUpdateTime()
		ac.mutation.SetUpdateTime(v)
	}
	if _, ok := ac.mutation.ChangedField(); !ok {
		return nil, errors.New("ent: missing required field \"changed_field\"")
	}
	if v, ok := ac.mutation.ChangedField(); ok {
		if err := activity.ChangedFieldValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"changed_field\": %v", err)
		}
	}
	if _, ok := ac.mutation.IsCreate(); !ok {
		v := activity.DefaultIsCreate
		ac.mutation.SetIsCreate(v)
	}
	var (
		err  error
		node *Activity
	)
	if len(ac.hooks) == 0 {
		node, err = ac.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ActivityMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ac.mutation = mutation
			node, err = ac.sqlSave(ctx)
			return node, err
		})
		for i := len(ac.hooks) - 1; i >= 0; i-- {
			mut = ac.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ac.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ac *ActivityCreate) SaveX(ctx context.Context) *Activity {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ac *ActivityCreate) sqlSave(ctx context.Context) (*Activity, error) {
	var (
		a     = &Activity{config: ac.config}
		_spec = &sqlgraph.CreateSpec{
			Table: activity.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: activity.FieldID,
			},
		}
	)
	if value, ok := ac.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: activity.FieldCreateTime,
		})
		a.CreateTime = value
	}
	if value, ok := ac.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: activity.FieldUpdateTime,
		})
		a.UpdateTime = value
	}
	if value, ok := ac.mutation.ChangedField(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: activity.FieldChangedField,
		})
		a.ChangedField = value
	}
	if value, ok := ac.mutation.IsCreate(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: activity.FieldIsCreate,
		})
		a.IsCreate = value
	}
	if value, ok := ac.mutation.OldValue(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: activity.FieldOldValue,
		})
		a.OldValue = value
	}
	if value, ok := ac.mutation.NewValue(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: activity.FieldNewValue,
		})
		a.NewValue = value
	}
	if nodes := ac.mutation.AuthorIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.WorkOrderIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	a.ID = int(id)
	return a, nil
}
