// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/facebookincubator/symphony/frontier/ent/migrate"

	"github.com/facebookincubator/symphony/frontier/ent/auditlog"
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
	"github.com/facebookincubator/symphony/frontier/ent/token"
	"github.com/facebookincubator/symphony/frontier/ent/user"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// AuditLog is the client for interacting with the AuditLog builders.
	AuditLog *AuditLogClient
	// Tenant is the client for interacting with the Tenant builders.
	Tenant *TenantClient
	// Token is the client for interacting with the Token builders.
	Token *TokenClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.AuditLog = NewAuditLogClient(c.config)
	c.Tenant = NewTenantClient(c.config)
	c.Token = NewTokenClient(c.config)
	c.User = NewUserClient(c.config)
}

// Open opens a connection to the database specified by the driver name and a
// driver-specific data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: tx, log: c.log, debug: c.debug, hooks: c.hooks}
	return &Tx{
		config:   cfg,
		AuditLog: NewAuditLogClient(cfg),
		Tenant:   NewTenantClient(cfg),
		Token:    NewTokenClient(cfg),
		User:     NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		AuditLog.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.AuditLog.Use(hooks...)
	c.Tenant.Use(hooks...)
	c.Token.Use(hooks...)
	c.User.Use(hooks...)
}

// AuditLogClient is a client for the AuditLog schema.
type AuditLogClient struct {
	config
}

// NewAuditLogClient returns a client for the AuditLog from the given config.
func NewAuditLogClient(c config) *AuditLogClient {
	return &AuditLogClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `auditlog.Hooks(f(g(h())))`.
func (c *AuditLogClient) Use(hooks ...Hook) {
	c.hooks.AuditLog = append(c.hooks.AuditLog, hooks...)
}

// Create returns a create builder for AuditLog.
func (c *AuditLogClient) Create() *AuditLogCreate {
	mutation := newAuditLogMutation(c.config, OpCreate)
	return &AuditLogCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for AuditLog.
func (c *AuditLogClient) Update() *AuditLogUpdate {
	mutation := newAuditLogMutation(c.config, OpUpdate)
	return &AuditLogUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *AuditLogClient) UpdateOne(al *AuditLog) *AuditLogUpdateOne {
	return c.UpdateOneID(al.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *AuditLogClient) UpdateOneID(id int) *AuditLogUpdateOne {
	mutation := newAuditLogMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &AuditLogUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for AuditLog.
func (c *AuditLogClient) Delete() *AuditLogDelete {
	mutation := newAuditLogMutation(c.config, OpDelete)
	return &AuditLogDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *AuditLogClient) DeleteOne(al *AuditLog) *AuditLogDeleteOne {
	return c.DeleteOneID(al.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *AuditLogClient) DeleteOneID(id int) *AuditLogDeleteOne {
	builder := c.Delete().Where(auditlog.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &AuditLogDeleteOne{builder}
}

// Create returns a query builder for AuditLog.
func (c *AuditLogClient) Query() *AuditLogQuery {
	return &AuditLogQuery{config: c.config}
}

// Get returns a AuditLog entity by its id.
func (c *AuditLogClient) Get(ctx context.Context, id int) (*AuditLog, error) {
	return c.Query().Where(auditlog.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *AuditLogClient) GetX(ctx context.Context, id int) *AuditLog {
	al, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return al
}

// Hooks returns the client hooks.
func (c *AuditLogClient) Hooks() []Hook {
	return c.hooks.AuditLog
}

// TenantClient is a client for the Tenant schema.
type TenantClient struct {
	config
}

// NewTenantClient returns a client for the Tenant from the given config.
func NewTenantClient(c config) *TenantClient {
	return &TenantClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `tenant.Hooks(f(g(h())))`.
func (c *TenantClient) Use(hooks ...Hook) {
	c.hooks.Tenant = append(c.hooks.Tenant, hooks...)
}

// Create returns a create builder for Tenant.
func (c *TenantClient) Create() *TenantCreate {
	mutation := newTenantMutation(c.config, OpCreate)
	return &TenantCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Tenant.
func (c *TenantClient) Update() *TenantUpdate {
	mutation := newTenantMutation(c.config, OpUpdate)
	return &TenantUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TenantClient) UpdateOne(t *Tenant) *TenantUpdateOne {
	return c.UpdateOneID(t.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *TenantClient) UpdateOneID(id int) *TenantUpdateOne {
	mutation := newTenantMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &TenantUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Tenant.
func (c *TenantClient) Delete() *TenantDelete {
	mutation := newTenantMutation(c.config, OpDelete)
	return &TenantDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *TenantClient) DeleteOne(t *Tenant) *TenantDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *TenantClient) DeleteOneID(id int) *TenantDeleteOne {
	builder := c.Delete().Where(tenant.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TenantDeleteOne{builder}
}

// Create returns a query builder for Tenant.
func (c *TenantClient) Query() *TenantQuery {
	return &TenantQuery{config: c.config}
}

// Get returns a Tenant entity by its id.
func (c *TenantClient) Get(ctx context.Context, id int) (*Tenant, error) {
	return c.Query().Where(tenant.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TenantClient) GetX(ctx context.Context, id int) *Tenant {
	t, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return t
}

// Hooks returns the client hooks.
func (c *TenantClient) Hooks() []Hook {
	return c.hooks.Tenant
}

// TokenClient is a client for the Token schema.
type TokenClient struct {
	config
}

// NewTokenClient returns a client for the Token from the given config.
func NewTokenClient(c config) *TokenClient {
	return &TokenClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `token.Hooks(f(g(h())))`.
func (c *TokenClient) Use(hooks ...Hook) {
	c.hooks.Token = append(c.hooks.Token, hooks...)
}

// Create returns a create builder for Token.
func (c *TokenClient) Create() *TokenCreate {
	mutation := newTokenMutation(c.config, OpCreate)
	return &TokenCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for Token.
func (c *TokenClient) Update() *TokenUpdate {
	mutation := newTokenMutation(c.config, OpUpdate)
	return &TokenUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TokenClient) UpdateOne(t *Token) *TokenUpdateOne {
	return c.UpdateOneID(t.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *TokenClient) UpdateOneID(id int) *TokenUpdateOne {
	mutation := newTokenMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &TokenUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Token.
func (c *TokenClient) Delete() *TokenDelete {
	mutation := newTokenMutation(c.config, OpDelete)
	return &TokenDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *TokenClient) DeleteOne(t *Token) *TokenDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *TokenClient) DeleteOneID(id int) *TokenDeleteOne {
	builder := c.Delete().Where(token.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TokenDeleteOne{builder}
}

// Create returns a query builder for Token.
func (c *TokenClient) Query() *TokenQuery {
	return &TokenQuery{config: c.config}
}

// Get returns a Token entity by its id.
func (c *TokenClient) Get(ctx context.Context, id int) (*Token, error) {
	return c.Query().Where(token.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TokenClient) GetX(ctx context.Context, id int) *Token {
	t, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return t
}

// QueryUser queries the user edge of a Token.
func (c *TokenClient) QueryUser(t *Token) *UserQuery {
	query := &UserQuery{config: c.config}
	id := t.ID
	step := sqlgraph.NewStep(
		sqlgraph.From(token.Table, token.FieldID, id),
		sqlgraph.To(user.Table, user.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, token.UserTable, token.UserColumn),
	)
	query.sql = sqlgraph.Neighbors(t.driver.Dialect(), step)

	return query
}

// Hooks returns the client hooks.
func (c *TokenClient) Hooks() []Hook {
	return c.hooks.Token
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Create returns a create builder for User.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	return c.UpdateOneID(u.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id int) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne)
	mutation.id = &id
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UserClient) DeleteOneID(id int) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Create returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{config: c.config}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id int) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id int) *User {
	u, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return u
}

// QueryTokens queries the tokens edge of a User.
func (c *UserClient) QueryTokens(u *User) *TokenQuery {
	query := &TokenQuery{config: c.config}
	id := u.ID
	step := sqlgraph.NewStep(
		sqlgraph.From(user.Table, user.FieldID, id),
		sqlgraph.To(token.Table, token.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, user.TokensTable, user.TokensColumn),
	)
	query.sql = sqlgraph.Neighbors(u.driver.Dialect(), step)

	return query
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}
