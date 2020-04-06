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

// MySQLTenancy provides logical database per tenant.
type MySQLTenancy struct {
	health.Checker
	clients sync.Map
	mu      sync.Mutex
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

// WithLogger sets tenancy logger.
func (m *MySQLTenancy) WithLogger(logger log.Logger) *MySQLTenancy {
	m.logger = logger
	return m
}

// ClientFor implements Tenancy interface.
func (m *MySQLTenancy) ClientFor(ctx context.Context, name string) (*ent.Client, error) {
	if client, ok := m.clients.Load(name); ok {
		return client.(*ent.Client), nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if client, ok := m.clients.Load(name); ok {
		return client.(*ent.Client), nil
	}
	client, err := m.clientFor(ctx, name)
	if err != nil {
		return nil, err
	}
	m.clients.Store(name, client)
	return client, nil
}

func (m *MySQLTenancy) clientFor(ctx context.Context, name string) (*ent.Client, error) {
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
