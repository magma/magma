// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migrate

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/schema"

	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	// Migrator performs schema migration.
	Migrator struct {
		driver     dialect.Driver
		logger     log.Logger
		options    []schema.MigrateOption
		newCreator func(dialect.Driver) Creator
	}

	// MigratorConfig configures migrator.
	MigratorConfig struct {
		dialect.Driver
		log.Logger
		Options []schema.MigrateOption
		Creator func(dialect.Driver) Creator
	}

	// Creator defines the interface for schema creation.
	Creator interface {
		Create(context.Context, ...schema.MigrateOption) error
	}

	// CreatorFunc is a function adapter implementing Creator interface.
	CreatorFunc func(context.Context, ...schema.MigrateOption) error
)

// Create invokes f(ctx, opts...).
func (f CreatorFunc) Create(ctx context.Context, opts ...schema.MigrateOption) error {
	return f(ctx, opts...)
}

// NewMigrator create schema based migrator from config.
func NewMigrator(cfg MigratorConfig) *Migrator {
	if cfg.Logger == nil {
		cfg.Logger = log.NewNopLogger()
	}
	if cfg.Creator == nil {
		cfg.Creator = func(driver dialect.Driver) Creator {
			return migrate.NewSchema(driver)
		}
	}
	cfg.Options = append(cfg.Options,
		schema.WithGlobalUniqueID(true),
	)
	return &Migrator{
		driver:     cfg.Driver,
		logger:     cfg.Logger,
		options:    cfg.Options,
		newCreator: cfg.Creator,
	}
}

// Migrate perform schema migration.
func (m *Migrator) Migrate(ctx context.Context, tenants ...string) error {
	logger := m.logger.For(ctx)
	tx, err := m.driver.Tx(ctx)
	if err != nil {
		logger.Error("cannot begin transaction", zap.Error(err))
		return errors.Wrap(err, "beginning transaction")
	}
	migration := migration{m, tx}
	if err := migration.Do(ctx, tenants...); err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error("cannot rollback transaction", zap.Error(err))
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		logger.Error("cannot commit transaction", zap.Error(err))
		return errors.Wrap(err, "committing transaction")
	}
	logger.Info("finished migrations", zap.Strings("tenants", tenants))
	return nil
}

type migration struct {
	*Migrator
	tx dialect.Tx
}

func (m *migration) Do(ctx context.Context, tenants ...string) error {
	logger := m.logger.For(ctx)
	for _, tenant := range tenants {
		logger := logger.With(zap.String("tenant", tenant))
		logger.Info("running migration")
		if err := m.do(ctx, tenant); err != nil {
			logger.Error("cannot run migration", zap.Error(err))
			return err
		}
	}
	return nil
}

func (m *migration) do(ctx context.Context, tenant string) error {
	query := fmt.Sprintf("USE `%s`", viewer.DBName(tenant))
	if err := m.tx.Exec(ctx, query, []interface{}{}, new(sql.Result)); err != nil {
		return errors.Wrap(err, "switching database")
	}
	if err := m.newCreator(m).Create(ctx, m.options...); err != nil {
		return errors.Wrap(err, "migrating schema")
	}
	return nil
}

func (m *migration) Exec(ctx context.Context, query string, args, v interface{}) error {
	return m.tx.Exec(ctx, query, args, v)
}

func (m *migration) Query(ctx context.Context, query string, args, v interface{}) error {
	return m.tx.Query(ctx, query, args, v)
}

func (m *migration) Dialect() string {
	return m.driver.Dialect()
}

func (m *migration) Tx(context.Context) (dialect.Tx, error) {
	return dialect.NopTx(m), nil
}

func (migration) Close() error { return nil }
