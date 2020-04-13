// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/facebookincubator/ent/dialect"
	entsql "github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/pkg/log"
	pkgmysql "github.com/facebookincubator/symphony/pkg/mysql"
	"go.opencensus.io/trace"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gocloud.dev/server/health"
	"gocloud.dev/server/health/sqlhealth"
)

// Tenancy provides tenant client for key.
type Tenancy interface {
	ClientFor(context.Context, string) (*ent.Client, error)
}

// FixedTenancy returns a fixed client.
type FixedTenancy struct {
	client *ent.Client
}

// NewFixedTenancy creates fixed tenancy from client.
func NewFixedTenancy(client *ent.Client) FixedTenancy {
	return FixedTenancy{client}
}

// ClientFor implements Tenancy interface.
func (f FixedTenancy) ClientFor(context.Context, string) (*ent.Client, error) {
	return f.Client(), nil
}

// Client returns the client stored in fixed tenancy.
func (f FixedTenancy) Client() *ent.Client {
	return f.client
}

// CacheTenancy is a tenancy wrapper cashing underlying clients.
type CacheTenancy struct {
	tenancy  Tenancy
	initFunc func(*ent.Client)
	clients  map[string]*ent.Client
	mu       sync.RWMutex
}

// NewCacheTenancy creates a tenancy cache.
func NewCacheTenancy(tenancy Tenancy, initFunc func(*ent.Client)) *CacheTenancy {
	return &CacheTenancy{
		tenancy:  tenancy,
		initFunc: initFunc,
		clients:  map[string]*ent.Client{},
	}
}

// ClientFor implements Tenancy interface.
func (c *CacheTenancy) ClientFor(ctx context.Context, name string) (*ent.Client, error) {
	c.mu.RLock()
	client, ok := c.clients[name]
	c.mu.RUnlock()
	if ok {
		return client, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if client, ok := c.clients[name]; ok {
		return client, nil
	}
	client, err := c.tenancy.ClientFor(ctx, name)
	if err != nil {
		return client, err
	}
	if c.initFunc != nil {
		c.initFunc(client)
	}
	c.clients[name] = client
	return client, nil
}

// CheckHealth implements health.Checker interface.
func (c *CacheTenancy) CheckHealth() error {
	if checker, ok := c.tenancy.(health.Checker); ok {
		return checker.CheckHealth()
	}
	return nil
}

// MySQLTenancy provides logical database per tenant.
type MySQLTenancy struct {
	health.Checker
	logger  log.Logger
	config  *mysql.Config
	closers []func()
}

// NewMySQLTenancy creates mysql tenancy for data source.
func NewMySQLTenancy(dsn string) (*MySQLTenancy, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing dsn: %w", err)
	}
	db := pkgmysql.Open(dsn)
	checker := sqlhealth.New(db)
	tenancy := &MySQLTenancy{
		Checker: checker,
		config:  config,
		logger:  log.NewNopLogger(),
		closers: []func(){checker.Stop},
	}
	runtime.SetFinalizer(tenancy, func(tenancy *MySQLTenancy) {
		for _, closer := range tenancy.closers {
			closer()
		}
	})
	return tenancy, nil
}

// SetLogger sets tenancy logger.
func (m *MySQLTenancy) SetLogger(logger log.Logger) {
	m.logger = logger
}

// ClientFor implements Tenancy interface.
func (m *MySQLTenancy) ClientFor(ctx context.Context, name string) (*ent.Client, error) {
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.MySQL, m.dbFor(name))))
	if err := m.migrate(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (m *MySQLTenancy) migrate(ctx context.Context, client *ent.Client) error {
	ctx, span := trace.StartSpan(ctx, "tenancy.Migrate")
	defer span.End()
	if err := client.Schema.Create(ctx,
		migrate.WithFixture(false),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		m.logger.For(ctx).Error("tenancy migrate", zap.Error(err))
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		return fmt.Errorf("running tenancy migration: %w", err)
	}
	return nil
}

func (m *MySQLTenancy) dbFor(name string) *sql.DB {
	m.config.DBName = DBName(name)
	db := pkgmysql.Open(m.config.FormatDSN())
	db.SetMaxOpenConns(10)
	m.closers = append(m.closers, pkgmysql.RecordStats(db))
	return db
}

// DBName returns the prefixed database name in order to avoid collision with MySQL internal databases.
func DBName(name string) string {
	return "tenant_" + name
}

// FromDBName returns the source name of the tenant.
func FromDBName(name string) string {
	return strings.TrimPrefix(name, "tenant_")
}
