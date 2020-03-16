// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
	"github.com/facebookincubator/symphony/frontier/ent/token"
	"github.com/facebookincubator/symphony/frontier/ent/user"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeAuditLog = "AuditLog"
	TypeTenant   = "Tenant"
	TypeToken    = "Token"
	TypeUser     = "User"
)

// AuditLogMutation represents an operation that mutate the AuditLogs
// nodes in the graph.
type AuditLogMutation struct {
	config
	op                  Op
	typ                 string
	id                  *int
	created_at          *time.Time
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
	clearedFields       map[string]bool
}

var _ ent.Mutation = (*AuditLogMutation)(nil)

// newAuditLogMutation creates new mutation for $n.Name.
func newAuditLogMutation(c config, op Op) *AuditLogMutation {
	return &AuditLogMutation{
		config:        c,
		op:            op,
		typ:           TypeAuditLog,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m AuditLogMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m AuditLogMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *AuditLogMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the created_at field.
func (m *AuditLogMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the created_at value in the mutation.
func (m *AuditLogMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreatedAt reset all changes of the created_at field.
func (m *AuditLogMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the updated_at field.
func (m *AuditLogMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the updated_at value in the mutation.
func (m *AuditLogMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdatedAt reset all changes of the updated_at field.
func (m *AuditLogMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetActingUserID sets the acting_user_id field.
func (m *AuditLogMutation) SetActingUserID(i int) {
	m.acting_user_id = &i
	m.addacting_user_id = nil
}

// ActingUserID returns the acting_user_id value in the mutation.
func (m *AuditLogMutation) ActingUserID() (r int, exists bool) {
	v := m.acting_user_id
	if v == nil {
		return
	}
	return *v, true
}

// AddActingUserID adds i to acting_user_id.
func (m *AuditLogMutation) AddActingUserID(i int) {
	if m.addacting_user_id != nil {
		*m.addacting_user_id += i
	} else {
		m.addacting_user_id = &i
	}
}

// AddedActingUserID returns the value that was added to the acting_user_id field in this mutation.
func (m *AuditLogMutation) AddedActingUserID() (r int, exists bool) {
	v := m.addacting_user_id
	if v == nil {
		return
	}
	return *v, true
}

// ResetActingUserID reset all changes of the acting_user_id field.
func (m *AuditLogMutation) ResetActingUserID() {
	m.acting_user_id = nil
	m.addacting_user_id = nil
}

// SetOrganization sets the organization field.
func (m *AuditLogMutation) SetOrganization(s string) {
	m.organization = &s
}

// Organization returns the organization value in the mutation.
func (m *AuditLogMutation) Organization() (r string, exists bool) {
	v := m.organization
	if v == nil {
		return
	}
	return *v, true
}

// ResetOrganization reset all changes of the organization field.
func (m *AuditLogMutation) ResetOrganization() {
	m.organization = nil
}

// SetMutationType sets the mutation_type field.
func (m *AuditLogMutation) SetMutationType(s string) {
	m.mutation_type = &s
}

// MutationType returns the mutation_type value in the mutation.
func (m *AuditLogMutation) MutationType() (r string, exists bool) {
	v := m.mutation_type
	if v == nil {
		return
	}
	return *v, true
}

// ResetMutationType reset all changes of the mutation_type field.
func (m *AuditLogMutation) ResetMutationType() {
	m.mutation_type = nil
}

// SetObjectID sets the object_id field.
func (m *AuditLogMutation) SetObjectID(s string) {
	m.object_id = &s
}

// ObjectID returns the object_id value in the mutation.
func (m *AuditLogMutation) ObjectID() (r string, exists bool) {
	v := m.object_id
	if v == nil {
		return
	}
	return *v, true
}

// ResetObjectID reset all changes of the object_id field.
func (m *AuditLogMutation) ResetObjectID() {
	m.object_id = nil
}

// SetObjectType sets the object_type field.
func (m *AuditLogMutation) SetObjectType(s string) {
	m.object_type = &s
}

// ObjectType returns the object_type value in the mutation.
func (m *AuditLogMutation) ObjectType() (r string, exists bool) {
	v := m.object_type
	if v == nil {
		return
	}
	return *v, true
}

// ResetObjectType reset all changes of the object_type field.
func (m *AuditLogMutation) ResetObjectType() {
	m.object_type = nil
}

// SetObjectDisplayName sets the object_display_name field.
func (m *AuditLogMutation) SetObjectDisplayName(s string) {
	m.object_display_name = &s
}

// ObjectDisplayName returns the object_display_name value in the mutation.
func (m *AuditLogMutation) ObjectDisplayName() (r string, exists bool) {
	v := m.object_display_name
	if v == nil {
		return
	}
	return *v, true
}

// ResetObjectDisplayName reset all changes of the object_display_name field.
func (m *AuditLogMutation) ResetObjectDisplayName() {
	m.object_display_name = nil
}

// SetMutationData sets the mutation_data field.
func (m *AuditLogMutation) SetMutationData(value map[string]string) {
	m.mutation_data = &value
}

// MutationData returns the mutation_data value in the mutation.
func (m *AuditLogMutation) MutationData() (r map[string]string, exists bool) {
	v := m.mutation_data
	if v == nil {
		return
	}
	return *v, true
}

// ResetMutationData reset all changes of the mutation_data field.
func (m *AuditLogMutation) ResetMutationData() {
	m.mutation_data = nil
}

// SetURL sets the url field.
func (m *AuditLogMutation) SetURL(s string) {
	m.url = &s
}

// URL returns the url value in the mutation.
func (m *AuditLogMutation) URL() (r string, exists bool) {
	v := m.url
	if v == nil {
		return
	}
	return *v, true
}

// ResetURL reset all changes of the url field.
func (m *AuditLogMutation) ResetURL() {
	m.url = nil
}

// SetIPAddress sets the ip_address field.
func (m *AuditLogMutation) SetIPAddress(s string) {
	m.ip_address = &s
}

// IPAddress returns the ip_address value in the mutation.
func (m *AuditLogMutation) IPAddress() (r string, exists bool) {
	v := m.ip_address
	if v == nil {
		return
	}
	return *v, true
}

// ResetIPAddress reset all changes of the ip_address field.
func (m *AuditLogMutation) ResetIPAddress() {
	m.ip_address = nil
}

// SetStatus sets the status field.
func (m *AuditLogMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the status value in the mutation.
func (m *AuditLogMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatus reset all changes of the status field.
func (m *AuditLogMutation) ResetStatus() {
	m.status = nil
}

// SetStatusCode sets the status_code field.
func (m *AuditLogMutation) SetStatusCode(s string) {
	m.status_code = &s
}

// StatusCode returns the status_code value in the mutation.
func (m *AuditLogMutation) StatusCode() (r string, exists bool) {
	v := m.status_code
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatusCode reset all changes of the status_code field.
func (m *AuditLogMutation) ResetStatusCode() {
	m.status_code = nil
}

// Op returns the operation name.
func (m *AuditLogMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (AuditLog).
func (m *AuditLogMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *AuditLogMutation) Fields() []string {
	fields := make([]string, 0, 13)
	if m.created_at != nil {
		fields = append(fields, auditlog.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, auditlog.FieldUpdatedAt)
	}
	if m.acting_user_id != nil {
		fields = append(fields, auditlog.FieldActingUserID)
	}
	if m.organization != nil {
		fields = append(fields, auditlog.FieldOrganization)
	}
	if m.mutation_type != nil {
		fields = append(fields, auditlog.FieldMutationType)
	}
	if m.object_id != nil {
		fields = append(fields, auditlog.FieldObjectID)
	}
	if m.object_type != nil {
		fields = append(fields, auditlog.FieldObjectType)
	}
	if m.object_display_name != nil {
		fields = append(fields, auditlog.FieldObjectDisplayName)
	}
	if m.mutation_data != nil {
		fields = append(fields, auditlog.FieldMutationData)
	}
	if m.url != nil {
		fields = append(fields, auditlog.FieldURL)
	}
	if m.ip_address != nil {
		fields = append(fields, auditlog.FieldIPAddress)
	}
	if m.status != nil {
		fields = append(fields, auditlog.FieldStatus)
	}
	if m.status_code != nil {
		fields = append(fields, auditlog.FieldStatusCode)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *AuditLogMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case auditlog.FieldCreatedAt:
		return m.CreatedAt()
	case auditlog.FieldUpdatedAt:
		return m.UpdatedAt()
	case auditlog.FieldActingUserID:
		return m.ActingUserID()
	case auditlog.FieldOrganization:
		return m.Organization()
	case auditlog.FieldMutationType:
		return m.MutationType()
	case auditlog.FieldObjectID:
		return m.ObjectID()
	case auditlog.FieldObjectType:
		return m.ObjectType()
	case auditlog.FieldObjectDisplayName:
		return m.ObjectDisplayName()
	case auditlog.FieldMutationData:
		return m.MutationData()
	case auditlog.FieldURL:
		return m.URL()
	case auditlog.FieldIPAddress:
		return m.IPAddress()
	case auditlog.FieldStatus:
		return m.Status()
	case auditlog.FieldStatusCode:
		return m.StatusCode()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *AuditLogMutation) SetField(name string, value ent.Value) error {
	switch name {
	case auditlog.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case auditlog.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case auditlog.FieldActingUserID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetActingUserID(v)
		return nil
	case auditlog.FieldOrganization:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOrganization(v)
		return nil
	case auditlog.FieldMutationType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMutationType(v)
		return nil
	case auditlog.FieldObjectID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetObjectID(v)
		return nil
	case auditlog.FieldObjectType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetObjectType(v)
		return nil
	case auditlog.FieldObjectDisplayName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetObjectDisplayName(v)
		return nil
	case auditlog.FieldMutationData:
		v, ok := value.(map[string]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMutationData(v)
		return nil
	case auditlog.FieldURL:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetURL(v)
		return nil
	case auditlog.FieldIPAddress:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIPAddress(v)
		return nil
	case auditlog.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case auditlog.FieldStatusCode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatusCode(v)
		return nil
	}
	return fmt.Errorf("unknown AuditLog field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *AuditLogMutation) AddedFields() []string {
	var fields []string
	if m.addacting_user_id != nil {
		fields = append(fields, auditlog.FieldActingUserID)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *AuditLogMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case auditlog.FieldActingUserID:
		return m.AddedActingUserID()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *AuditLogMutation) AddField(name string, value ent.Value) error {
	switch name {
	case auditlog.FieldActingUserID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddActingUserID(v)
		return nil
	}
	return fmt.Errorf("unknown AuditLog numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *AuditLogMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *AuditLogMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *AuditLogMutation) ClearField(name string) error {
	return fmt.Errorf("unknown AuditLog nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *AuditLogMutation) ResetField(name string) error {
	switch name {
	case auditlog.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case auditlog.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case auditlog.FieldActingUserID:
		m.ResetActingUserID()
		return nil
	case auditlog.FieldOrganization:
		m.ResetOrganization()
		return nil
	case auditlog.FieldMutationType:
		m.ResetMutationType()
		return nil
	case auditlog.FieldObjectID:
		m.ResetObjectID()
		return nil
	case auditlog.FieldObjectType:
		m.ResetObjectType()
		return nil
	case auditlog.FieldObjectDisplayName:
		m.ResetObjectDisplayName()
		return nil
	case auditlog.FieldMutationData:
		m.ResetMutationData()
		return nil
	case auditlog.FieldURL:
		m.ResetURL()
		return nil
	case auditlog.FieldIPAddress:
		m.ResetIPAddress()
		return nil
	case auditlog.FieldStatus:
		m.ResetStatus()
		return nil
	case auditlog.FieldStatusCode:
		m.ResetStatusCode()
		return nil
	}
	return fmt.Errorf("unknown AuditLog field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *AuditLogMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *AuditLogMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *AuditLogMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *AuditLogMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *AuditLogMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *AuditLogMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *AuditLogMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown AuditLog unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *AuditLogMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown AuditLog edge %s", name)
}

// TenantMutation represents an operation that mutate the Tenants
// nodes in the graph.
type TenantMutation struct {
	config
	op             Op
	typ            string
	id             *int
	created_at     *time.Time
	updated_at     *time.Time
	name           *string
	domains        *[]string
	networks       *[]string
	tabs           *[]string
	_SSOCert       *string
	_SSOEntryPoint *string
	_SSOIssuer     *string
	clearedFields  map[string]bool
}

var _ ent.Mutation = (*TenantMutation)(nil)

// newTenantMutation creates new mutation for $n.Name.
func newTenantMutation(c config, op Op) *TenantMutation {
	return &TenantMutation{
		config:        c,
		op:            op,
		typ:           TypeTenant,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m TenantMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m TenantMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *TenantMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the created_at field.
func (m *TenantMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the created_at value in the mutation.
func (m *TenantMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreatedAt reset all changes of the created_at field.
func (m *TenantMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the updated_at field.
func (m *TenantMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the updated_at value in the mutation.
func (m *TenantMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdatedAt reset all changes of the updated_at field.
func (m *TenantMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetName sets the name field.
func (m *TenantMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *TenantMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *TenantMutation) ResetName() {
	m.name = nil
}

// SetDomains sets the domains field.
func (m *TenantMutation) SetDomains(s []string) {
	m.domains = &s
}

// Domains returns the domains value in the mutation.
func (m *TenantMutation) Domains() (r []string, exists bool) {
	v := m.domains
	if v == nil {
		return
	}
	return *v, true
}

// ResetDomains reset all changes of the domains field.
func (m *TenantMutation) ResetDomains() {
	m.domains = nil
}

// SetNetworks sets the networks field.
func (m *TenantMutation) SetNetworks(s []string) {
	m.networks = &s
}

// Networks returns the networks value in the mutation.
func (m *TenantMutation) Networks() (r []string, exists bool) {
	v := m.networks
	if v == nil {
		return
	}
	return *v, true
}

// ResetNetworks reset all changes of the networks field.
func (m *TenantMutation) ResetNetworks() {
	m.networks = nil
}

// SetTabs sets the tabs field.
func (m *TenantMutation) SetTabs(s []string) {
	m.tabs = &s
}

// Tabs returns the tabs value in the mutation.
func (m *TenantMutation) Tabs() (r []string, exists bool) {
	v := m.tabs
	if v == nil {
		return
	}
	return *v, true
}

// ClearTabs clears the value of tabs.
func (m *TenantMutation) ClearTabs() {
	m.tabs = nil
	m.clearedFields[tenant.FieldTabs] = true
}

// TabsCleared returns if the field tabs was cleared in this mutation.
func (m *TenantMutation) TabsCleared() bool {
	return m.clearedFields[tenant.FieldTabs]
}

// ResetTabs reset all changes of the tabs field.
func (m *TenantMutation) ResetTabs() {
	m.tabs = nil
	delete(m.clearedFields, tenant.FieldTabs)
}

// SetSSOCert sets the SSOCert field.
func (m *TenantMutation) SetSSOCert(s string) {
	m._SSOCert = &s
}

// SSOCert returns the SSOCert value in the mutation.
func (m *TenantMutation) SSOCert() (r string, exists bool) {
	v := m._SSOCert
	if v == nil {
		return
	}
	return *v, true
}

// ResetSSOCert reset all changes of the SSOCert field.
func (m *TenantMutation) ResetSSOCert() {
	m._SSOCert = nil
}

// SetSSOEntryPoint sets the SSOEntryPoint field.
func (m *TenantMutation) SetSSOEntryPoint(s string) {
	m._SSOEntryPoint = &s
}

// SSOEntryPoint returns the SSOEntryPoint value in the mutation.
func (m *TenantMutation) SSOEntryPoint() (r string, exists bool) {
	v := m._SSOEntryPoint
	if v == nil {
		return
	}
	return *v, true
}

// ResetSSOEntryPoint reset all changes of the SSOEntryPoint field.
func (m *TenantMutation) ResetSSOEntryPoint() {
	m._SSOEntryPoint = nil
}

// SetSSOIssuer sets the SSOIssuer field.
func (m *TenantMutation) SetSSOIssuer(s string) {
	m._SSOIssuer = &s
}

// SSOIssuer returns the SSOIssuer value in the mutation.
func (m *TenantMutation) SSOIssuer() (r string, exists bool) {
	v := m._SSOIssuer
	if v == nil {
		return
	}
	return *v, true
}

// ResetSSOIssuer reset all changes of the SSOIssuer field.
func (m *TenantMutation) ResetSSOIssuer() {
	m._SSOIssuer = nil
}

// Op returns the operation name.
func (m *TenantMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Tenant).
func (m *TenantMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *TenantMutation) Fields() []string {
	fields := make([]string, 0, 9)
	if m.created_at != nil {
		fields = append(fields, tenant.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, tenant.FieldUpdatedAt)
	}
	if m.name != nil {
		fields = append(fields, tenant.FieldName)
	}
	if m.domains != nil {
		fields = append(fields, tenant.FieldDomains)
	}
	if m.networks != nil {
		fields = append(fields, tenant.FieldNetworks)
	}
	if m.tabs != nil {
		fields = append(fields, tenant.FieldTabs)
	}
	if m._SSOCert != nil {
		fields = append(fields, tenant.FieldSSOCert)
	}
	if m._SSOEntryPoint != nil {
		fields = append(fields, tenant.FieldSSOEntryPoint)
	}
	if m._SSOIssuer != nil {
		fields = append(fields, tenant.FieldSSOIssuer)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *TenantMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case tenant.FieldCreatedAt:
		return m.CreatedAt()
	case tenant.FieldUpdatedAt:
		return m.UpdatedAt()
	case tenant.FieldName:
		return m.Name()
	case tenant.FieldDomains:
		return m.Domains()
	case tenant.FieldNetworks:
		return m.Networks()
	case tenant.FieldTabs:
		return m.Tabs()
	case tenant.FieldSSOCert:
		return m.SSOCert()
	case tenant.FieldSSOEntryPoint:
		return m.SSOEntryPoint()
	case tenant.FieldSSOIssuer:
		return m.SSOIssuer()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TenantMutation) SetField(name string, value ent.Value) error {
	switch name {
	case tenant.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case tenant.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case tenant.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case tenant.FieldDomains:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDomains(v)
		return nil
	case tenant.FieldNetworks:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetNetworks(v)
		return nil
	case tenant.FieldTabs:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTabs(v)
		return nil
	case tenant.FieldSSOCert:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSSOCert(v)
		return nil
	case tenant.FieldSSOEntryPoint:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSSOEntryPoint(v)
		return nil
	case tenant.FieldSSOIssuer:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSSOIssuer(v)
		return nil
	}
	return fmt.Errorf("unknown Tenant field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *TenantMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *TenantMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TenantMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Tenant numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *TenantMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[tenant.FieldTabs] {
		fields = append(fields, tenant.FieldTabs)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *TenantMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *TenantMutation) ClearField(name string) error {
	switch name {
	case tenant.FieldTabs:
		m.ClearTabs()
		return nil
	}
	return fmt.Errorf("unknown Tenant nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *TenantMutation) ResetField(name string) error {
	switch name {
	case tenant.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case tenant.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case tenant.FieldName:
		m.ResetName()
		return nil
	case tenant.FieldDomains:
		m.ResetDomains()
		return nil
	case tenant.FieldNetworks:
		m.ResetNetworks()
		return nil
	case tenant.FieldTabs:
		m.ResetTabs()
		return nil
	case tenant.FieldSSOCert:
		m.ResetSSOCert()
		return nil
	case tenant.FieldSSOEntryPoint:
		m.ResetSSOEntryPoint()
		return nil
	case tenant.FieldSSOIssuer:
		m.ResetSSOIssuer()
		return nil
	}
	return fmt.Errorf("unknown Tenant field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *TenantMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *TenantMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *TenantMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *TenantMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *TenantMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *TenantMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *TenantMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Tenant unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *TenantMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Tenant edge %s", name)
}

// TokenMutation represents an operation that mutate the Tokens
// nodes in the graph.
type TokenMutation struct {
	config
	op            Op
	typ           string
	id            *int
	created_at    *time.Time
	updated_at    *time.Time
	value         *string
	clearedFields map[string]bool
	user          *int
	cleareduser   bool
}

var _ ent.Mutation = (*TokenMutation)(nil)

// newTokenMutation creates new mutation for $n.Name.
func newTokenMutation(c config, op Op) *TokenMutation {
	return &TokenMutation{
		config:        c,
		op:            op,
		typ:           TypeToken,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m TokenMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m TokenMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *TokenMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the created_at field.
func (m *TokenMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the created_at value in the mutation.
func (m *TokenMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreatedAt reset all changes of the created_at field.
func (m *TokenMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the updated_at field.
func (m *TokenMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the updated_at value in the mutation.
func (m *TokenMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdatedAt reset all changes of the updated_at field.
func (m *TokenMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetValue sets the value field.
func (m *TokenMutation) SetValue(s string) {
	m.value = &s
}

// Value returns the value value in the mutation.
func (m *TokenMutation) Value() (r string, exists bool) {
	v := m.value
	if v == nil {
		return
	}
	return *v, true
}

// ResetValue reset all changes of the value field.
func (m *TokenMutation) ResetValue() {
	m.value = nil
}

// SetUserID sets the user edge to User by id.
func (m *TokenMutation) SetUserID(id int) {
	m.user = &id
}

// ClearUser clears the user edge to User.
func (m *TokenMutation) ClearUser() {
	m.cleareduser = true
}

// UserCleared returns if the edge user was cleared.
func (m *TokenMutation) UserCleared() bool {
	return m.cleareduser
}

// UserID returns the user id in the mutation.
func (m *TokenMutation) UserID() (id int, exists bool) {
	if m.user != nil {
		return *m.user, true
	}
	return
}

// UserIDs returns the user ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// UserID instead. It exists only for internal usage by the builders.
func (m *TokenMutation) UserIDs() (ids []int) {
	if id := m.user; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetUser reset all changes of the user edge.
func (m *TokenMutation) ResetUser() {
	m.user = nil
	m.cleareduser = false
}

// Op returns the operation name.
func (m *TokenMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Token).
func (m *TokenMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *TokenMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.created_at != nil {
		fields = append(fields, token.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, token.FieldUpdatedAt)
	}
	if m.value != nil {
		fields = append(fields, token.FieldValue)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *TokenMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case token.FieldCreatedAt:
		return m.CreatedAt()
	case token.FieldUpdatedAt:
		return m.UpdatedAt()
	case token.FieldValue:
		return m.Value()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TokenMutation) SetField(name string, value ent.Value) error {
	switch name {
	case token.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case token.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case token.FieldValue:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetValue(v)
		return nil
	}
	return fmt.Errorf("unknown Token field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *TokenMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *TokenMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TokenMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Token numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *TokenMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *TokenMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *TokenMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Token nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *TokenMutation) ResetField(name string) error {
	switch name {
	case token.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case token.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case token.FieldValue:
		m.ResetValue()
		return nil
	}
	return fmt.Errorf("unknown Token field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *TokenMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.user != nil {
		edges = append(edges, token.EdgeUser)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *TokenMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case token.EdgeUser:
		if id := m.user; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *TokenMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *TokenMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *TokenMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.cleareduser {
		edges = append(edges, token.EdgeUser)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *TokenMutation) EdgeCleared(name string) bool {
	switch name {
	case token.EdgeUser:
		return m.cleareduser
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *TokenMutation) ClearEdge(name string) error {
	switch name {
	case token.EdgeUser:
		m.ClearUser()
		return nil
	}
	return fmt.Errorf("unknown Token unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *TokenMutation) ResetEdge(name string) error {
	switch name {
	case token.EdgeUser:
		m.ResetUser()
		return nil
	}
	return fmt.Errorf("unknown Token edge %s", name)
}

// UserMutation represents an operation that mutate the Users
// nodes in the graph.
type UserMutation struct {
	config
	op            Op
	typ           string
	id            *int
	created_at    *time.Time
	updated_at    *time.Time
	email         *string
	password      *string
	role          *int
	addrole       *int
	tenant        *string
	networks      *[]string
	tabs          *[]string
	clearedFields map[string]bool
	tokens        map[int]struct{}
	removedtokens map[int]struct{}
}

var _ ent.Mutation = (*UserMutation)(nil)

// newUserMutation creates new mutation for $n.Name.
func newUserMutation(c config, op Op) *UserMutation {
	return &UserMutation{
		config:        c,
		op:            op,
		typ:           TypeUser,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m UserMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m UserMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *UserMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the created_at field.
func (m *UserMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the created_at value in the mutation.
func (m *UserMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreatedAt reset all changes of the created_at field.
func (m *UserMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the updated_at field.
func (m *UserMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the updated_at value in the mutation.
func (m *UserMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdatedAt reset all changes of the updated_at field.
func (m *UserMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetEmail sets the email field.
func (m *UserMutation) SetEmail(s string) {
	m.email = &s
}

// Email returns the email value in the mutation.
func (m *UserMutation) Email() (r string, exists bool) {
	v := m.email
	if v == nil {
		return
	}
	return *v, true
}

// ResetEmail reset all changes of the email field.
func (m *UserMutation) ResetEmail() {
	m.email = nil
}

// SetPassword sets the password field.
func (m *UserMutation) SetPassword(s string) {
	m.password = &s
}

// Password returns the password value in the mutation.
func (m *UserMutation) Password() (r string, exists bool) {
	v := m.password
	if v == nil {
		return
	}
	return *v, true
}

// ResetPassword reset all changes of the password field.
func (m *UserMutation) ResetPassword() {
	m.password = nil
}

// SetRole sets the role field.
func (m *UserMutation) SetRole(i int) {
	m.role = &i
	m.addrole = nil
}

// Role returns the role value in the mutation.
func (m *UserMutation) Role() (r int, exists bool) {
	v := m.role
	if v == nil {
		return
	}
	return *v, true
}

// AddRole adds i to role.
func (m *UserMutation) AddRole(i int) {
	if m.addrole != nil {
		*m.addrole += i
	} else {
		m.addrole = &i
	}
}

// AddedRole returns the value that was added to the role field in this mutation.
func (m *UserMutation) AddedRole() (r int, exists bool) {
	v := m.addrole
	if v == nil {
		return
	}
	return *v, true
}

// ResetRole reset all changes of the role field.
func (m *UserMutation) ResetRole() {
	m.role = nil
	m.addrole = nil
}

// SetTenant sets the tenant field.
func (m *UserMutation) SetTenant(s string) {
	m.tenant = &s
}

// Tenant returns the tenant value in the mutation.
func (m *UserMutation) Tenant() (r string, exists bool) {
	v := m.tenant
	if v == nil {
		return
	}
	return *v, true
}

// ResetTenant reset all changes of the tenant field.
func (m *UserMutation) ResetTenant() {
	m.tenant = nil
}

// SetNetworks sets the networks field.
func (m *UserMutation) SetNetworks(s []string) {
	m.networks = &s
}

// Networks returns the networks value in the mutation.
func (m *UserMutation) Networks() (r []string, exists bool) {
	v := m.networks
	if v == nil {
		return
	}
	return *v, true
}

// ResetNetworks reset all changes of the networks field.
func (m *UserMutation) ResetNetworks() {
	m.networks = nil
}

// SetTabs sets the tabs field.
func (m *UserMutation) SetTabs(s []string) {
	m.tabs = &s
}

// Tabs returns the tabs value in the mutation.
func (m *UserMutation) Tabs() (r []string, exists bool) {
	v := m.tabs
	if v == nil {
		return
	}
	return *v, true
}

// ClearTabs clears the value of tabs.
func (m *UserMutation) ClearTabs() {
	m.tabs = nil
	m.clearedFields[user.FieldTabs] = true
}

// TabsCleared returns if the field tabs was cleared in this mutation.
func (m *UserMutation) TabsCleared() bool {
	return m.clearedFields[user.FieldTabs]
}

// ResetTabs reset all changes of the tabs field.
func (m *UserMutation) ResetTabs() {
	m.tabs = nil
	delete(m.clearedFields, user.FieldTabs)
}

// AddTokenIDs adds the tokens edge to Token by ids.
func (m *UserMutation) AddTokenIDs(ids ...int) {
	if m.tokens == nil {
		m.tokens = make(map[int]struct{})
	}
	for i := range ids {
		m.tokens[ids[i]] = struct{}{}
	}
}

// RemoveTokenIDs removes the tokens edge to Token by ids.
func (m *UserMutation) RemoveTokenIDs(ids ...int) {
	if m.removedtokens == nil {
		m.removedtokens = make(map[int]struct{})
	}
	for i := range ids {
		m.removedtokens[ids[i]] = struct{}{}
	}
}

// RemovedTokens returns the removed ids of tokens.
func (m *UserMutation) RemovedTokensIDs() (ids []int) {
	for id := range m.removedtokens {
		ids = append(ids, id)
	}
	return
}

// TokensIDs returns the tokens ids in the mutation.
func (m *UserMutation) TokensIDs() (ids []int) {
	for id := range m.tokens {
		ids = append(ids, id)
	}
	return
}

// ResetTokens reset all changes of the tokens edge.
func (m *UserMutation) ResetTokens() {
	m.tokens = nil
	m.removedtokens = nil
}

// Op returns the operation name.
func (m *UserMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (User).
func (m *UserMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *UserMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.created_at != nil {
		fields = append(fields, user.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, user.FieldUpdatedAt)
	}
	if m.email != nil {
		fields = append(fields, user.FieldEmail)
	}
	if m.password != nil {
		fields = append(fields, user.FieldPassword)
	}
	if m.role != nil {
		fields = append(fields, user.FieldRole)
	}
	if m.tenant != nil {
		fields = append(fields, user.FieldTenant)
	}
	if m.networks != nil {
		fields = append(fields, user.FieldNetworks)
	}
	if m.tabs != nil {
		fields = append(fields, user.FieldTabs)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *UserMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case user.FieldCreatedAt:
		return m.CreatedAt()
	case user.FieldUpdatedAt:
		return m.UpdatedAt()
	case user.FieldEmail:
		return m.Email()
	case user.FieldPassword:
		return m.Password()
	case user.FieldRole:
		return m.Role()
	case user.FieldTenant:
		return m.Tenant()
	case user.FieldNetworks:
		return m.Networks()
	case user.FieldTabs:
		return m.Tabs()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *UserMutation) SetField(name string, value ent.Value) error {
	switch name {
	case user.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case user.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case user.FieldEmail:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEmail(v)
		return nil
	case user.FieldPassword:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPassword(v)
		return nil
	case user.FieldRole:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRole(v)
		return nil
	case user.FieldTenant:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTenant(v)
		return nil
	case user.FieldNetworks:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetNetworks(v)
		return nil
	case user.FieldTabs:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTabs(v)
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *UserMutation) AddedFields() []string {
	var fields []string
	if m.addrole != nil {
		fields = append(fields, user.FieldRole)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *UserMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case user.FieldRole:
		return m.AddedRole()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *UserMutation) AddField(name string, value ent.Value) error {
	switch name {
	case user.FieldRole:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRole(v)
		return nil
	}
	return fmt.Errorf("unknown User numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *UserMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[user.FieldTabs] {
		fields = append(fields, user.FieldTabs)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *UserMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *UserMutation) ClearField(name string) error {
	switch name {
	case user.FieldTabs:
		m.ClearTabs()
		return nil
	}
	return fmt.Errorf("unknown User nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *UserMutation) ResetField(name string) error {
	switch name {
	case user.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case user.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case user.FieldEmail:
		m.ResetEmail()
		return nil
	case user.FieldPassword:
		m.ResetPassword()
		return nil
	case user.FieldRole:
		m.ResetRole()
		return nil
	case user.FieldTenant:
		m.ResetTenant()
		return nil
	case user.FieldNetworks:
		m.ResetNetworks()
		return nil
	case user.FieldTabs:
		m.ResetTabs()
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *UserMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.tokens != nil {
		edges = append(edges, user.EdgeTokens)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *UserMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case user.EdgeTokens:
		ids := make([]ent.Value, 0, len(m.tokens))
		for id := range m.tokens {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *UserMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedtokens != nil {
		edges = append(edges, user.EdgeTokens)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *UserMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case user.EdgeTokens:
		ids := make([]ent.Value, 0, len(m.removedtokens))
		for id := range m.removedtokens {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *UserMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *UserMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *UserMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown User unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *UserMutation) ResetEdge(name string) error {
	switch name {
	case user.EdgeTokens:
		m.ResetTokens()
		return nil
	}
	return fmt.Errorf("unknown User edge %s", name)
}
