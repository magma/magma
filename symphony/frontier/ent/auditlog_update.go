// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
)

// AuditLogUpdate is the builder for updating AuditLog entities.
type AuditLogUpdate struct {
	config

	updated_at          *time.Time
	acting_user_id      *int
	addacting_user_id   *int
	organization        *string
	mutation_type       *string
	object_id           *string
	object_type         *string
	object_display_name *string
	mutation_data       *map[string]string
	url                 *string
	ip_address          *string
	status              *string
	status_code         *string
	predicates          []predicate.AuditLog
}

// Where adds a new predicate for the builder.
func (alu *AuditLogUpdate) Where(ps ...predicate.AuditLog) *AuditLogUpdate {
	alu.predicates = append(alu.predicates, ps...)
	return alu
}

// SetActingUserID sets the acting_user_id field.
func (alu *AuditLogUpdate) SetActingUserID(i int) *AuditLogUpdate {
	alu.acting_user_id = &i
	alu.addacting_user_id = nil
	return alu
}

// AddActingUserID adds i to acting_user_id.
func (alu *AuditLogUpdate) AddActingUserID(i int) *AuditLogUpdate {
	if alu.addacting_user_id == nil {
		alu.addacting_user_id = &i
	} else {
		*alu.addacting_user_id += i
	}
	return alu
}

// SetOrganization sets the organization field.
func (alu *AuditLogUpdate) SetOrganization(s string) *AuditLogUpdate {
	alu.organization = &s
	return alu
}

// SetMutationType sets the mutation_type field.
func (alu *AuditLogUpdate) SetMutationType(s string) *AuditLogUpdate {
	alu.mutation_type = &s
	return alu
}

// SetObjectID sets the object_id field.
func (alu *AuditLogUpdate) SetObjectID(s string) *AuditLogUpdate {
	alu.object_id = &s
	return alu
}

// SetObjectType sets the object_type field.
func (alu *AuditLogUpdate) SetObjectType(s string) *AuditLogUpdate {
	alu.object_type = &s
	return alu
}

// SetObjectDisplayName sets the object_display_name field.
func (alu *AuditLogUpdate) SetObjectDisplayName(s string) *AuditLogUpdate {
	alu.object_display_name = &s
	return alu
}

// SetMutationData sets the mutation_data field.
func (alu *AuditLogUpdate) SetMutationData(m map[string]string) *AuditLogUpdate {
	alu.mutation_data = &m
	return alu
}

// SetURL sets the url field.
func (alu *AuditLogUpdate) SetURL(s string) *AuditLogUpdate {
	alu.url = &s
	return alu
}

// SetIPAddress sets the ip_address field.
func (alu *AuditLogUpdate) SetIPAddress(s string) *AuditLogUpdate {
	alu.ip_address = &s
	return alu
}

// SetStatus sets the status field.
func (alu *AuditLogUpdate) SetStatus(s string) *AuditLogUpdate {
	alu.status = &s
	return alu
}

// SetStatusCode sets the status_code field.
func (alu *AuditLogUpdate) SetStatusCode(s string) *AuditLogUpdate {
	alu.status_code = &s
	return alu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (alu *AuditLogUpdate) Save(ctx context.Context) (int, error) {
	if alu.updated_at == nil {
		v := auditlog.UpdateDefaultUpdatedAt()
		alu.updated_at = &v
	}
	return alu.sqlSave(ctx)
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
	if value := alu.updated_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: auditlog.FieldUpdatedAt,
		})
	}
	if value := alu.acting_user_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value := alu.addacting_user_id; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value := alu.organization; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldOrganization,
		})
	}
	if value := alu.mutation_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldMutationType,
		})
	}
	if value := alu.object_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldObjectID,
		})
	}
	if value := alu.object_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldObjectType,
		})
	}
	if value := alu.object_display_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldObjectDisplayName,
		})
	}
	if value := alu.mutation_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: auditlog.FieldMutationData,
		})
	}
	if value := alu.url; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldURL,
		})
	}
	if value := alu.ip_address; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldIPAddress,
		})
	}
	if value := alu.status; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldStatus,
		})
	}
	if value := alu.status_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldStatusCode,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, alu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// AuditLogUpdateOne is the builder for updating a single AuditLog entity.
type AuditLogUpdateOne struct {
	config
	id int

	updated_at          *time.Time
	acting_user_id      *int
	addacting_user_id   *int
	organization        *string
	mutation_type       *string
	object_id           *string
	object_type         *string
	object_display_name *string
	mutation_data       *map[string]string
	url                 *string
	ip_address          *string
	status              *string
	status_code         *string
}

// SetActingUserID sets the acting_user_id field.
func (aluo *AuditLogUpdateOne) SetActingUserID(i int) *AuditLogUpdateOne {
	aluo.acting_user_id = &i
	aluo.addacting_user_id = nil
	return aluo
}

// AddActingUserID adds i to acting_user_id.
func (aluo *AuditLogUpdateOne) AddActingUserID(i int) *AuditLogUpdateOne {
	if aluo.addacting_user_id == nil {
		aluo.addacting_user_id = &i
	} else {
		*aluo.addacting_user_id += i
	}
	return aluo
}

// SetOrganization sets the organization field.
func (aluo *AuditLogUpdateOne) SetOrganization(s string) *AuditLogUpdateOne {
	aluo.organization = &s
	return aluo
}

// SetMutationType sets the mutation_type field.
func (aluo *AuditLogUpdateOne) SetMutationType(s string) *AuditLogUpdateOne {
	aluo.mutation_type = &s
	return aluo
}

// SetObjectID sets the object_id field.
func (aluo *AuditLogUpdateOne) SetObjectID(s string) *AuditLogUpdateOne {
	aluo.object_id = &s
	return aluo
}

// SetObjectType sets the object_type field.
func (aluo *AuditLogUpdateOne) SetObjectType(s string) *AuditLogUpdateOne {
	aluo.object_type = &s
	return aluo
}

// SetObjectDisplayName sets the object_display_name field.
func (aluo *AuditLogUpdateOne) SetObjectDisplayName(s string) *AuditLogUpdateOne {
	aluo.object_display_name = &s
	return aluo
}

// SetMutationData sets the mutation_data field.
func (aluo *AuditLogUpdateOne) SetMutationData(m map[string]string) *AuditLogUpdateOne {
	aluo.mutation_data = &m
	return aluo
}

// SetURL sets the url field.
func (aluo *AuditLogUpdateOne) SetURL(s string) *AuditLogUpdateOne {
	aluo.url = &s
	return aluo
}

// SetIPAddress sets the ip_address field.
func (aluo *AuditLogUpdateOne) SetIPAddress(s string) *AuditLogUpdateOne {
	aluo.ip_address = &s
	return aluo
}

// SetStatus sets the status field.
func (aluo *AuditLogUpdateOne) SetStatus(s string) *AuditLogUpdateOne {
	aluo.status = &s
	return aluo
}

// SetStatusCode sets the status_code field.
func (aluo *AuditLogUpdateOne) SetStatusCode(s string) *AuditLogUpdateOne {
	aluo.status_code = &s
	return aluo
}

// Save executes the query and returns the updated entity.
func (aluo *AuditLogUpdateOne) Save(ctx context.Context) (*AuditLog, error) {
	if aluo.updated_at == nil {
		v := auditlog.UpdateDefaultUpdatedAt()
		aluo.updated_at = &v
	}
	return aluo.sqlSave(ctx)
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
				Value:  aluo.id,
				Type:   field.TypeInt,
				Column: auditlog.FieldID,
			},
		},
	}
	if value := aluo.updated_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: auditlog.FieldUpdatedAt,
		})
	}
	if value := aluo.acting_user_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value := aluo.addacting_user_id; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: auditlog.FieldActingUserID,
		})
	}
	if value := aluo.organization; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldOrganization,
		})
	}
	if value := aluo.mutation_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldMutationType,
		})
	}
	if value := aluo.object_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldObjectID,
		})
	}
	if value := aluo.object_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldObjectType,
		})
	}
	if value := aluo.object_display_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldObjectDisplayName,
		})
	}
	if value := aluo.mutation_data; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: auditlog.FieldMutationData,
		})
	}
	if value := aluo.url; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldURL,
		})
	}
	if value := aluo.ip_address; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldIPAddress,
		})
	}
	if value := aluo.status; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldStatus,
		})
	}
	if value := aluo.status_code; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: auditlog.FieldStatusCode,
		})
	}
	al = &AuditLog{config: aluo.config}
	_spec.Assign = al.assignValues
	_spec.ScanValues = al.scanValues()
	if err = sqlgraph.UpdateNode(ctx, aluo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return al, nil
}
