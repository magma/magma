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
	"github.com/facebookincubator/symphony/cloud/log"
	entmigrate "github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/migrate"

	"github.com/facebookincubator/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
)

func main() {
	drv := flag.String("db_driver", "", "driver name")
	dsn := flag.String("db_dsn", "", "data source name")
	dc := flag.Bool("drop_column", false, "enable column drop")
	di := flag.Bool("drop_index", false, "enable index drop")
	plan := flag.Bool("plan", false, "print the execution plan")
	flag.Parse()

	logger, _ := log.Config{Format: "console"}.Build()
	driver, err := sql.Open(*drv, *dsn)
	if err != nil {
		logger.Background().Fatal("opening database", zap.Error(err))
	}

	tenants, err := graphgrpc.NewTenantService(driver.DB()).
		List(context.Background(), &empty.Empty{})
	if err != nil {
		logger.Background().Fatal("listing tenants", zap.Error(err))
	}

	names := make([]string, len(tenants.Tenants))
	for i, tenant := range tenants.Tenants {
		names[i] = tenant.Name
	}

	cfg := migrate.MigratorConfig{
		Driver: driver,
		Logger: logger,
		Options: []schema.MigrateOption{
			schema.WithDropColumn(*dc),
			schema.WithDropIndex(*di),
		},
	}
	if *plan {
		cfg.Creator = func(driver dialect.Driver) migrate.Creator {
			return planner{entmigrate.NewSchema(driver)}
		}
	}

	if err := migrate.NewMigrator(cfg).Migrate(context.Background(), names...); err != nil {
		os.Exit(1)
	}
}

type planner struct {
	*entmigrate.Schema
}

func (c planner) Create(ctx context.Context, opts ...schema.MigrateOption) error {
	return c.WriteTo(ctx, os.Stdout, opts...)
}
