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
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
)

// AuditLogCreate is the builder for creating a AuditLog entity.
type AuditLogCreate struct {
	config
	mutation *AuditLogMutation
	hooks    []Hook
}

// SetCreatedAt sets the created_at field.
func (alc *AuditLogCreate) SetCreatedAt(t time.Time) *AuditLogCreate {
	alc.mutation.SetCreatedAt(t)
	return alc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (alc *AuditLogCreate) SetNillableCreatedAt(t *time.Time) *AuditLogCreate {
	if t != nil {
		alc.SetCreatedAt(*t)
	}
	return alc
}

// SetUpdatedAt sets the updated_at field.
func (alc *AuditLogCreate) SetUpdatedAt(t time.Time) *AuditLogCreate {
	alc.mutation.SetUpdatedAt(t)
	return alc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (alc *AuditLogCreate) SetNillableUpdatedAt(t *time.Time) *AuditLogCreate {
	if t != nil {
		alc.SetUpdatedAt(*t)
	}
	return alc
}

// SetActingUserID sets the acting_user_id field.
func (alc *AuditLogCreate) SetActingUserID(i int) *AuditLogCreate {
	alc.mutation.SetActingUserID(i)
	return alc
}

// SetOrganization sets the organization field.
func (alc *AuditLogCreate) SetOrganization(s string) *AuditLogCreate {
	alc.mutation.SetOrganization(s)
	return alc
}

// SetMutationType sets the mutation_type field.
func (alc *AuditLogCreate) SetMutationType(s string) *AuditLogCreate {
	alc.mutation.SetMutationType(s)
	return alc
}

// SetObjectID sets the object_id field.
func (alc *AuditLogCreate) SetObjectID(s string) *AuditLogCreate {
	alc.mutation.SetObjectID(s)
	return alc
}

// SetObjectType sets the object_type field.
func (alc *AuditLogCreate) SetObjectType(s string) *AuditLogCreate {
	alc.mutation.SetObjectType(s)
	return alc
}

// SetObjectDisplayName sets the object_display_name field.
func (alc *AuditLogCreate) SetObjectDisplayName(s string) *AuditLogCreate {
	alc.mutation.SetObjectDisplayName(s)
	return alc
}

// SetMutationData sets the mutation_data field.
func (alc *AuditLogCreate) SetMutationData(m map[string]string) *AuditLogCreate {
	alc.mutation.SetMutationData(m)
	return alc
}

// SetURL sets the url field.
func (alc *AuditLogCreate) SetURL(s string) *AuditLogCreate {
	alc.mutation.SetURL(s)
	return alc
}

// SetIPAddress sets the ip_address field.
func (alc *AuditLogCreate) SetIPAddress(s string) *AuditLogCreate {
	alc.mutation.SetIPAddress(s)
	return alc
}

// SetStatus sets the status field.
func (alc *AuditLogCreate) SetStatus(s string) *AuditLogCreate {
	alc.mutation.SetStatus(s)
	return alc
}

// SetStatusCode sets the status_code field.
func (alc *AuditLogCreate) SetStatusCode(s string) *AuditLogCreate {
	alc.mutation.SetStatusCode(s)
	return alc
}

// Save creates the AuditLog in the database.
func (alc *AuditLogCreate) Save(ctx context.Context) (*AuditLog, error) {
	if _, ok := alc.mutation.CreatedAt(); !ok {
		v := auditlog.DefaultCreatedAt()
		alc.mutation.SetCreatedAt(v)
	}
	if _, ok := alc.mutation.UpdatedAt(); !ok {
		v := auditlog.DefaultUpdatedAt()
		alc.mutation.SetUpdatedAt(v)
	}
	if _, ok := alc.mutation.ActingUserID(); !ok {
		return nil, errors.New("ent: missing required field \"acting_user_id\"")
	}
	if _, ok := alc.mutation.Organization(); !ok {
		return nil, errors.New("ent: missing required field \"organization\"")
	}
	if _, ok := alc.mutation.MutationType(); !ok {
		return nil, errors.New("ent: missing required field \"mutation_type\"")
	}
	if _, ok := alc.mutation.ObjectID(); !ok {
		return nil, errors.New("ent: missing required field \"object_id\"")
	}
	if _, ok := alc.mutation.ObjectType(); !ok {
		return nil, errors.New("ent: missing required field \"object_type\"")
	}
	if _, ok := alc.mutation.ObjectDisplayName(); !ok {
		return nil, errors.New("ent: missing required field \"object_display_name\"")
	}
	if _, ok := alc.mutation.MutationData(); !ok {
		return nil, errors.New("ent: missing required field \"mutation_data\"")
	}
	if _, ok := alc.mutation.URL(); !ok {
		return nil, errors.New("ent: missing required field \"url\"")
	}
	if _, ok := alc.mutation.IPAddress(); !ok {
		return nil, errors.New("ent: missing required field \"ip_address\"")
	}
	if _, ok := alc.mutation.Status(); !ok {
		return nil, errors.New("ent: missing required field \"status\"")
	}
	if _, ok := alc.mutation.StatusCode(); !ok {
		return nil, errors.New("ent: missing required field \"status_code\"")
	}
	var (
		err  error
		node *AuditLog
	)
	if len(alc.hooks) == 0 {
		node, err = alc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AuditLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			alc.mutation = mutation
			node, err = alc.sqlSave(ctx)
			return node, err
		})
		for i := len(alc.hooks) - 1; i >= 0; i-- {
			mut = alc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, alc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (alc *AuditLogCreate) SaveX(ctx context.Context) *AuditLog {
	v, err := alc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (alc *AuditLogCreate) sqlSave(ctx context.Context) (*AuditLog, error) {
	var (
		al    = &AuditLog{config: alc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: auditlog.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: auditlog.FieldID,
			},
		}
	)
	if value, ok := alc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: auditlog.FieldCreatedAt,
		})
		al.CreatedAt = value
	}
	if value, ok := alc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: auditlog.FieldUpdatedAt,
		})
		al.UpdatedAt = value
	}
	if value, ok := alc.mutation.ActingUserID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: auditlog.FieldActingUserID,
		})
		al.ActingUserID = value
	}
	if value, ok := alc.mutation.Organization(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldOrganization,
		})
		al.Organization = value
	}
	if value, ok := alc.mutation.MutationType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldMutationType,
		})
		al.MutationType = value
	}
	if value, ok := alc.mutation.ObjectID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectID,
		})
		al.ObjectID = value
	}
	if value, ok := alc.mutation.ObjectType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectType,
		})
		al.ObjectType = value
	}
	if value, ok := alc.mutation.ObjectDisplayName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldObjectDisplayName,
		})
		al.ObjectDisplayName = value
	}
	if value, ok := alc.mutation.MutationData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: auditlog.FieldMutationData,
		})
		al.MutationData = value
	}
	if value, ok := alc.mutation.URL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldURL,
		})
		al.URL = value
	}
	if value, ok := alc.mutation.IPAddress(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldIPAddress,
		})
		al.IPAddress = value
	}
	if value, ok := alc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldStatus,
		})
		al.Status = value
	}
	if value, ok := alc.mutation.StatusCode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: auditlog.FieldStatusCode,
		})
		al.StatusCode = value
	}
	if err := sqlgraph.CreateNode(ctx, alc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	al.ID = int(id)
	return al, nil
}
