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
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// AuditLogUpdate is the builder for updating AuditLog entities.
type AuditLogUpdate struct {
	config
	hooks      []Hook
	mutation   *AuditLogMutation
	predicates []predicate.AuditLog
}

// Where adds a new predicate for the builder.
func (alu *AuditLogUpdate) Where(ps ...predicate.AuditLog) *AuditLogUpdate {
	alu.predicates = append(alu.predicates, ps...)
	return alu
}

// SetActingUserID sets the acting_user_id field.
func (alu *AuditLogUpdate) SetActingUserID(i int) *AuditLogUpdate {
	alu.mutation.ResetActingUserID()
	alu.mutation.SetActingUserID(i)
	return alu
}

// AddActingUserID adds i to acting_user_id.
func (alu *AuditLogUpdate) AddActingUserID(i int) *AuditLogUpdate {
	alu.mutation.AddActingUserID(i)
	return alu
}

// SetOrganization sets the organization field.
func (alu *AuditLogUpdate) SetOrganization(s string) *AuditLogUpdate {
	alu.mutation.SetOrganization(s)
	return alu
}

// SetMutationType sets the mutation_type field.
func (alu *AuditLogUpdate) SetMutationType(s string) *AuditLogUpdate {
	alu.mutation.SetMutationType(s)
	return alu
}

// SetObjectID sets the object_id field.
func (alu *AuditLogUpdate) SetObjectID(s string) *AuditLogUpdate {
	alu.mutation.SetObjectID(s)
	return alu
}

// SetObjectType sets the object_type field.
func (alu *AuditLogUpdate) SetObjectType(s string) *AuditLogUpdate {
	alu.mutation.SetObjectType(s)
	return alu
}

// SetObjectDisplayName sets the object_display_name field.
func (alu *AuditLogUpdate) SetObjectDisplayName(s string) *AuditLogUpdate {
	alu.mutation.SetObjectDisplayName(s)
	return alu
}

// SetMutationData sets the mutation_data field.
func (alu *AuditLogUpdate) SetMutationData(m map[string]string) *AuditLogUpdate {
	alu.mutation.SetMutationData(m)
	return alu
}

// SetURL sets the url field.
func (alu *AuditLogUpdate) SetURL(s string) *AuditLogUpdate {
	alu.mutation.SetURL(s)
	return alu
}

// SetIPAddress sets the ip_address field.
func (alu *AuditLogUpdate) SetIPAddress(s string) *AuditLogUpdate {
	alu.mutation.SetIPAddress(s)
	return alu
}

// SetStatus sets the status field.
func (alu *AuditLogUpdate) SetStatus(s string) *AuditLogUpdate {
	alu.mutation.SetStatus(s)
	return alu
}

// SetStatusCode sets the status_code field.
func (alu *AuditLogUpdate) SetStatusCode(s string) *AuditLogUpdate {
	alu.mutation.SetStatusCode(s)
	return alu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (alu *AuditLogUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := alu.mutation.UpdatedAt(); !ok {
		v := auditlog.UpdateDefaultUpdatedAt()
		alu.mutation.SetUpdatedAt(v)
	}
	var (
		err      error
		affected int
	)
	if len(alu.hooks) == 0 {
		affected, err = alu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AuditLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			alu.mutation = mutation
			affected, err = alu.sqlSave(ctx)
			return affected, err
		})
		for i := len(alu.hooks); i > 0; i-- {
			mut = alu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, alu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (alu *AuditLogUpdate) SaveX(ctx context.Context) int {
	affected, err := alu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (alu *AuditLogUpdate) Exec(ctx context.Context) error {
	_, err := alu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (alu *AuditLogUpdate) ExecX(ctx context.Context) {
	if err := alu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (alu *AuditLogUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   auditlog.Table,
			Columns: auditlog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: auditlog.FieldID,
			},
		},
	}
	if ps := alu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := alu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: auditlog.FieldUpdatedAt,
		})
	}
	if value, ok := alu.mutation.ActingUserID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value, ok := alu.mutation.AddedActingUserID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value, ok := alu.mutation.Organization(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldOrganization,
		})
	}
	if value, ok := alu.mutation.MutationType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldMutationType,
		})
	}
	if value, ok := alu.mutation.ObjectID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectID,
		})
	}
	if value, ok := alu.mutation.ObjectType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectType,
		})
	}
	if value, ok := alu.mutation.ObjectDisplayName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectDisplayName,
		})
	}
	if value, ok := alu.mutation.MutationData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: auditlog.FieldMutationData,
		})
	}
	if value, ok := alu.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldURL,
		})
	}
	if value, ok := alu.mutation.IPAddress(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldIPAddress,
		})
	}
	if value, ok := alu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldStatus,
		})
	}
	if value, ok := alu.mutation.StatusCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldStatusCode,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, alu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{auditlog.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// AuditLogUpdateOne is the builder for updating a single AuditLog entity.
type AuditLogUpdateOne struct {
	config
	hooks    []Hook
	mutation *AuditLogMutation
}

// SetActingUserID sets the acting_user_id field.
func (aluo *AuditLogUpdateOne) SetActingUserID(i int) *AuditLogUpdateOne {
	aluo.mutation.ResetActingUserID()
	aluo.mutation.SetActingUserID(i)
	return aluo
}

// AddActingUserID adds i to acting_user_id.
func (aluo *AuditLogUpdateOne) AddActingUserID(i int) *AuditLogUpdateOne {
	aluo.mutation.AddActingUserID(i)
	return aluo
}

// SetOrganization sets the organization field.
func (aluo *AuditLogUpdateOne) SetOrganization(s string) *AuditLogUpdateOne {
	aluo.mutation.SetOrganization(s)
	return aluo
}

// SetMutationType sets the mutation_type field.
func (aluo *AuditLogUpdateOne) SetMutationType(s string) *AuditLogUpdateOne {
	aluo.mutation.SetMutationType(s)
	return aluo
}

// SetObjectID sets the object_id field.
func (aluo *AuditLogUpdateOne) SetObjectID(s string) *AuditLogUpdateOne {
	aluo.mutation.SetObjectID(s)
	return aluo
}

// SetObjectType sets the object_type field.
func (aluo *AuditLogUpdateOne) SetObjectType(s string) *AuditLogUpdateOne {
	aluo.mutation.SetObjectType(s)
	return aluo
}

// SetObjectDisplayName sets the object_display_name field.
func (aluo *AuditLogUpdateOne) SetObjectDisplayName(s string) *AuditLogUpdateOne {
	aluo.mutation.SetObjectDisplayName(s)
	return aluo
}

// SetMutationData sets the mutation_data field.
func (aluo *AuditLogUpdateOne) SetMutationData(m map[string]string) *AuditLogUpdateOne {
	aluo.mutation.SetMutationData(m)
	return aluo
}

// SetURL sets the url field.
func (aluo *AuditLogUpdateOne) SetURL(s string) *AuditLogUpdateOne {
	aluo.mutation.SetURL(s)
	return aluo
}

// SetIPAddress sets the ip_address field.
func (aluo *AuditLogUpdateOne) SetIPAddress(s string) *AuditLogUpdateOne {
	aluo.mutation.SetIPAddress(s)
	return aluo
}

// SetStatus sets the status field.
func (aluo *AuditLogUpdateOne) SetStatus(s string) *AuditLogUpdateOne {
	aluo.mutation.SetStatus(s)
	return aluo
}

// SetStatusCode sets the status_code field.
func (aluo *AuditLogUpdateOne) SetStatusCode(s string) *AuditLogUpdateOne {
	aluo.mutation.SetStatusCode(s)
	return aluo
}

// Save executes the query and returns the updated entity.
func (aluo *AuditLogUpdateOne) Save(ctx context.Context) (*AuditLog, error) {
	if _, ok := aluo.mutation.UpdatedAt(); !ok {
		v := auditlog.UpdateDefaultUpdatedAt()
		aluo.mutation.SetUpdatedAt(v)
	}
	var (
		err  error
		node *AuditLog
	)
	if len(aluo.hooks) == 0 {
		node, err = aluo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AuditLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			aluo.mutation = mutation
			node, err = aluo.sqlSave(ctx)
			return node, err
		})
		for i := len(aluo.hooks); i > 0; i-- {
			mut = aluo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, aluo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (aluo *AuditLogUpdateOne) SaveX(ctx context.Context) *AuditLog {
	al, err := aluo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return al
}

// Exec executes the query on the entity.
func (aluo *AuditLogUpdateOne) Exec(ctx context.Context) error {
	_, err := aluo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aluo *AuditLogUpdateOne) ExecX(ctx context.Context) {
	if err := aluo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (aluo *AuditLogUpdateOne) sqlSave(ctx context.Context) (al *AuditLog, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   auditlog.Table,
			Columns: auditlog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: auditlog.FieldID,
			},
		},
	}
	id, ok := aluo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing AuditLog.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := aluo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: auditlog.FieldUpdatedAt,
		})
	}
	if value, ok := aluo.mutation.ActingUserID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value, ok := aluo.mutation.AddedActingUserID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value, ok := aluo.mutation.Organization(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldOrganization,
		})
	}
	if value, ok := aluo.mutation.MutationType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldMutationType,
		})
	}
	if value, ok := aluo.mutation.ObjectID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectID,
		})
	}
	if value, ok := aluo.mutation.ObjectType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectType,
		})
	}
	if value, ok := aluo.mutation.ObjectDisplayName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectDisplayName,
		})
	}
	if value, ok := aluo.mutation.MutationData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: auditlog.FieldMutationData,
		})
	}
	if value, ok := aluo.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldURL,
		})
	}
	if value, ok := aluo.mutation.IPAddress(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldIPAddress,
		})
	}
	if value, ok := aluo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldStatus,
		})
	}
	if value, ok := aluo.mutation.StatusCode(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldStatusCode,
		})
	}
	al = &AuditLog{config: aluo.config}
	_spec.Assign = al.assignValues
	_spec.ScanValues = al.scanValues()
	if err = sqlgraph.UpdateNode(ctx, aluo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{auditlog.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return al, nil
}
