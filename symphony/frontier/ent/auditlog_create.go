// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
)

// AuditLogCreate is the builder for creating a AuditLog entity.
type AuditLogCreate struct {
	config
	created_at          *time.Time
	updated_at          *time.Time
	acting_user_id      *int
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

// SetCreatedAt sets the created_at field.
func (alc *AuditLogCreate) SetCreatedAt(t time.Time) *AuditLogCreate {
	alc.created_at = &t
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
	alc.updated_at = &t
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
	alc.acting_user_id = &i
	return alc
}

// SetOrganization sets the organization field.
func (alc *AuditLogCreate) SetOrganization(s string) *AuditLogCreate {
	alc.organization = &s
	return alc
}

// SetMutationType sets the mutation_type field.
func (alc *AuditLogCreate) SetMutationType(s string) *AuditLogCreate {
	alc.mutation_type = &s
	return alc
}

// SetObjectID sets the object_id field.
func (alc *AuditLogCreate) SetObjectID(s string) *AuditLogCreate {
	alc.object_id = &s
	return alc
}

// SetObjectType sets the object_type field.
func (alc *AuditLogCreate) SetObjectType(s string) *AuditLogCreate {
	alc.object_type = &s
	return alc
}

// SetObjectDisplayName sets the object_display_name field.
func (alc *AuditLogCreate) SetObjectDisplayName(s string) *AuditLogCreate {
	alc.object_display_name = &s
	return alc
}

// SetMutationData sets the mutation_data field.
func (alc *AuditLogCreate) SetMutationData(m map[string]string) *AuditLogCreate {
	alc.mutation_data = &m
	return alc
}

// SetURL sets the url field.
func (alc *AuditLogCreate) SetURL(s string) *AuditLogCreate {
	alc.url = &s
	return alc
}

// SetIPAddress sets the ip_address field.
func (alc *AuditLogCreate) SetIPAddress(s string) *AuditLogCreate {
	alc.ip_address = &s
	return alc
}

// SetStatus sets the status field.
func (alc *AuditLogCreate) SetStatus(s string) *AuditLogCreate {
	alc.status = &s
	return alc
}

// SetStatusCode sets the status_code field.
func (alc *AuditLogCreate) SetStatusCode(s string) *AuditLogCreate {
	alc.status_code = &s
	return alc
}

// Save creates the AuditLog in the database.
func (alc *AuditLogCreate) Save(ctx context.Context) (*AuditLog, error) {
	if alc.created_at == nil {
		v := auditlog.DefaultCreatedAt()
		alc.created_at = &v
	}
	if alc.updated_at == nil {
		v := auditlog.DefaultUpdatedAt()
		alc.updated_at = &v
	}
	if alc.acting_user_id == nil {
		return nil, errors.New("ent: missing required field \"acting_user_id\"")
	}
	if alc.organization == nil {
		return nil, errors.New("ent: missing required field \"organization\"")
	}
	if alc.mutation_type == nil {
		return nil, errors.New("ent: missing required field \"mutation_type\"")
	}
	if alc.object_id == nil {
		return nil, errors.New("ent: missing required field \"object_id\"")
	}
	if alc.object_type == nil {
		return nil, errors.New("ent: missing required field \"object_type\"")
	}
	if alc.object_display_name == nil {
		return nil, errors.New("ent: missing required field \"object_display_name\"")
	}
	if alc.mutation_data == nil {
		return nil, errors.New("ent: missing required field \"mutation_data\"")
	}
	if alc.url == nil {
		return nil, errors.New("ent: missing required field \"url\"")
	}
	if alc.ip_address == nil {
		return nil, errors.New("ent: missing required field \"ip_address\"")
	}
	if alc.status == nil {
		return nil, errors.New("ent: missing required field \"status\"")
	}
	if alc.status_code == nil {
		return nil, errors.New("ent: missing required field \"status_code\"")
	}
	return alc.sqlSave(ctx)
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
		builder = sql.Dialect(alc.driver.Dialect())
		al      = &AuditLog{config: alc.config}
	)
	tx, err := alc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(auditlog.Table).Default()
	if value := alc.created_at; value != nil {
		insert.Set(auditlog.FieldCreatedAt, *value)
		al.CreatedAt = *value
	}
	if value := alc.updated_at; value != nil {
		insert.Set(auditlog.FieldUpdatedAt, *value)
		al.UpdatedAt = *value
	}
	if value := alc.acting_user_id; value != nil {
		insert.Set(auditlog.FieldActingUserID, *value)
		al.ActingUserID = *value
	}
	if value := alc.organization; value != nil {
		insert.Set(auditlog.FieldOrganization, *value)
		al.Organization = *value
	}
	if value := alc.mutation_type; value != nil {
		insert.Set(auditlog.FieldMutationType, *value)
		al.MutationType = *value
	}
	if value := alc.object_id; value != nil {
		insert.Set(auditlog.FieldObjectID, *value)
		al.ObjectID = *value
	}
	if value := alc.object_type; value != nil {
		insert.Set(auditlog.FieldObjectType, *value)
		al.ObjectType = *value
	}
	if value := alc.object_display_name; value != nil {
		insert.Set(auditlog.FieldObjectDisplayName, *value)
		al.ObjectDisplayName = *value
	}
	if value := alc.mutation_data; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(auditlog.FieldMutationData, buf)
		al.MutationData = *value
	}
	if value := alc.url; value != nil {
		insert.Set(auditlog.FieldURL, *value)
		al.URL = *value
	}
	if value := alc.ip_address; value != nil {
		insert.Set(auditlog.FieldIPAddress, *value)
		al.IPAddress = *value
	}
	if value := alc.status; value != nil {
		insert.Set(auditlog.FieldStatus, *value)
		al.Status = *value
	}
	if value := alc.status_code; value != nil {
		insert.Set(auditlog.FieldStatusCode, *value)
		al.StatusCode = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(auditlog.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	al.ID = int(id)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return al, nil
}
