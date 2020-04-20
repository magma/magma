// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/schema"
	entmigrate "github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/migrate"
	"github.com/facebookincubator/symphony/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	_ "github.com/facebookincubator/symphony/graph/ent/runtime"
)

func main() {
	kingpin.HelpFlag.Short('h')
	drv := kingpin.Flag("db-driver", "database driver name").Default("mysql").String()
	dsn := kingpin.Flag("db-dsn", "data source name").Required().String()
	dropColumn := kingpin.Flag("drop-column", "enable column drop").Bool()
	dropIndex := kingpin.Flag("drop-index", "enable index drop").Bool()
	fixture := kingpin.Flag("fixture", "run ent@v0.1.0 migrate fixture").Bool()
	dryRun := kingpin.Flag("dry-run", "run in dry run mode").Bool()
	tenantName := kingpin.Flag("tenant", "target specific tenant").String()
	logcfg := log.AddFlags(kingpin.CommandLine)
	kingpin.Parse()

	logger, _, _ := log.Provider(*logcfg)
	driver, err := sql.Open(*drv, *dsn)
	if err != nil {
		logger.Background().Fatal("opening database", zap.Error(err))
	}

	tenants, err := graphgrpc.NewTenantService(
		func(context.Context) graphgrpc.ExecQueryer {
			return driver.DB()
		},
	).List(context.Background(), &empty.Empty{})
	if err != nil {
		logger.Background().Fatal("listing tenants", zap.Error(err))
	}

	names := make([]string, 0, len(tenants.Tenants))
	for _, tenant := range tenants.Tenants {
		if *tenantName == "" || *tenantName == tenant.Name {
			names = append(names, tenant.Name)
		}
	}

	cfg := migrate.MigratorConfig{
		Logger: logger,
		Driver: dialect.Debug(driver),
		Options: []schema.MigrateOption{
			schema.WithFixture(*fixture),
			schema.WithDropColumn(*dropColumn),
			schema.WithDropIndex(*dropIndex),
		},
	}
	if *dryRun {
		cfg.Creator = func(driver dialect.Driver) migrate.Creator {
			entSchema := entmigrate.NewSchema(driver)
			return migrate.CreatorFunc(func(ctx context.Context, opts ...schema.MigrateOption) error {
				return entSchema.WriteTo(ctx, os.Stdout, opts...)
			})
		}
	}

	if err := migrate.NewMigrator(cfg).Migrate(context.Background(), names...); err != nil {
		os.Exit(1)
	}
}
