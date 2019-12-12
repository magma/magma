// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
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
	var (
		builder  = sql.Dialect(alu.driver.Dialect())
		selector = builder.Select(auditlog.FieldID).From(builder.Table(auditlog.Table))
	)
	for _, p := range alu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = alu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := alu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(auditlog.Table)
	)
	updater = updater.Where(sql.InInts(auditlog.FieldID, ids...))
	if value := alu.updated_at; value != nil {
		updater.Set(auditlog.FieldUpdatedAt, *value)
	}
	if value := alu.acting_user_id; value != nil {
		updater.Set(auditlog.FieldActingUserID, *value)
	}
	if value := alu.addacting_user_id; value != nil {
		updater.Add(auditlog.FieldActingUserID, *value)
	}
	if value := alu.organization; value != nil {
		updater.Set(auditlog.FieldOrganization, *value)
	}
	if value := alu.mutation_type; value != nil {
		updater.Set(auditlog.FieldMutationType, *value)
	}
	if value := alu.object_id; value != nil {
		updater.Set(auditlog.FieldObjectID, *value)
	}
	if value := alu.object_type; value != nil {
		updater.Set(auditlog.FieldObjectType, *value)
	}
	if value := alu.object_display_name; value != nil {
		updater.Set(auditlog.FieldObjectDisplayName, *value)
	}
	if value := alu.mutation_data; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(auditlog.FieldMutationData, buf)
	}
	if value := alu.url; value != nil {
		updater.Set(auditlog.FieldURL, *value)
	}
	if value := alu.ip_address; value != nil {
		updater.Set(auditlog.FieldIPAddress, *value)
	}
	if value := alu.status; value != nil {
		updater.Set(auditlog.FieldStatus, *value)
	}
	if value := alu.status_code; value != nil {
		updater.Set(auditlog.FieldStatusCode, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
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
	var (
		builder  = sql.Dialect(aluo.driver.Dialect())
		selector = builder.Select(auditlog.Columns...).From(builder.Table(auditlog.Table))
	)
	auditlog.ID(aluo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = aluo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		al = &AuditLog{config: aluo.config}
		if err := al.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into AuditLog: %v", err)
		}
		id = al.ID
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("AuditLog with id: %v", aluo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one AuditLog with the same id: %v", aluo.id)
	}

	tx, err := aluo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(auditlog.Table)
	)
	updater = updater.Where(sql.InInts(auditlog.FieldID, ids...))
	if value := aluo.updated_at; value != nil {
		updater.Set(auditlog.FieldUpdatedAt, *value)
		al.UpdatedAt = *value
	}
	if value := aluo.acting_user_id; value != nil {
		updater.Set(auditlog.FieldActingUserID, *value)
		al.ActingUserID = *value
	}
	if value := aluo.addacting_user_id; value != nil {
		updater.Add(auditlog.FieldActingUserID, *value)
		al.ActingUserID += *value
	}
	if value := aluo.organization; value != nil {
		updater.Set(auditlog.FieldOrganization, *value)
		al.Organization = *value
	}
	if value := aluo.mutation_type; value != nil {
		updater.Set(auditlog.FieldMutationType, *value)
		al.MutationType = *value
	}
	if value := aluo.object_id; value != nil {
		updater.Set(auditlog.FieldObjectID, *value)
		al.ObjectID = *value
	}
	if value := aluo.object_type; value != nil {
		updater.Set(auditlog.FieldObjectType, *value)
		al.ObjectType = *value
	}
	if value := aluo.object_display_name; value != nil {
		updater.Set(auditlog.FieldObjectDisplayName, *value)
		al.ObjectDisplayName = *value
	}
	if value := aluo.mutation_data; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(auditlog.FieldMutationData, buf)
		al.MutationData = *value
	}
	if value := aluo.url; value != nil {
		updater.Set(auditlog.FieldURL, *value)
		al.URL = *value
	}
	if value := aluo.ip_address; value != nil {
		updater.Set(auditlog.FieldIPAddress, *value)
		al.IPAddress = *value
	}
	if value := aluo.status; value != nil {
		updater.Set(auditlog.FieldStatus, *value)
		al.Status = *value
	}
	if value := aluo.status_code; value != nil {
		updater.Set(auditlog.FieldStatusCode, *value)
		al.StatusCode = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return al, nil
}
