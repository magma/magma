// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/schema"
	entmigrate "github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/migrate"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/facebookincubator/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
)

func main() {
	drv := flag.String("db-driver", "mysql", "driver name")
	dsn := flag.String("db-dsn", "", "data source name")
	dropColumn := flag.Bool("drop-column", false, "enable column drop")
	dropIndex := flag.Bool("drop-index", false, "enable index drop")
	dryRun := flag.Bool("dry-run", false, "run in dry run mode")
	tenantName := flag.String("tenant", "", "target specific tenant")
	flag.Parse()

	logger, _, _ := log.New(log.Config{Format: "console"})
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
