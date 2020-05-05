// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	_ "github.com/facebookincubator/symphony/graph/ent/runtime"
)

func main() {
	kingpin.HelpFlag.Short('h')
	dsn := kingpin.Flag("db-dsn", "data source name").Envar("MYSQL_DSN").Required().String()
	tenantName := kingpin.Flag("tenant", "tenant name to target. \"ALL\" for running on all tenants").Required().String()
	username := kingpin.Flag("user", "who is running the script (for logging purposes)").Required().String()
	logcfg := log.AddFlags(kingpin.CommandLine)
	kingpin.Parse()

	logger, _, _ := log.Provider(*logcfg)
	ctx := context.Background()

	logger.For(ctx).Info("params",
		zap.Stringp("dsn", dsn),
		zap.Stringp("tenant", tenantName),
		zap.Stringp("user", username),
	)
	tenancy, err := viewer.NewMySQLTenancy(*dsn)
	if err != nil {
		logger.For(ctx).Fatal("cannot connect to graph database",
			zap.Stringp("dsn", dsn),
			zap.Error(err),
		)
	}
	mysql.SetLogger(logger)

	driver, err := sql.Open("mysql", *dsn)
	if err != nil {
		logger.For(ctx).Fatal("cannot connect sql database",
			zap.Stringp("dsn", dsn),
			zap.Error(err),
		)
	}

	tenants, err := getTenantList(ctx, driver, tenantName)
	if err != nil {
		logger.For(ctx).Fatal("cannot get tenants to run on",
			zap.Stringp("dsn", dsn),
			zap.Stringp("tenant", tenantName),
			zap.Error(err),
		)
	}

	for _, tenant := range tenants {
		client, err := tenancy.ClientFor(ctx, tenant)
		if err != nil {
			logger.For(ctx).Fatal("cannot get ent client for tenant",
				zap.String("tenant", tenant),
				zap.Error(err),
			)
		}
		ctx := ent.NewContext(ctx, client)
		v := viewer.NewAutomation(tenant, *username, user.RoleOWNER)
		ctx = log.NewFieldsContext(ctx, zap.Object("viewer", v))
		ctx = viewer.NewContext(ctx, v)
		permissions, err := authz.Permissions(ctx)
		if err != nil {
			logger.For(ctx).Fatal("cannot get permissions",
				zap.String("tenant", tenant),
				zap.Stringp("user", username),
				zap.Error(err),
			)
		}
		ctx = authz.NewContext(ctx, permissions)

		tx, err := client.Tx(ctx)
		if err != nil {
			logger.For(ctx).Fatal("cannot begin transaction", zap.Error(err))
		}
		defer func() {
			if r := recover(); r != nil {
				if err := tx.Rollback(); err != nil {
					logger.For(ctx).Error("cannot rollback transaction", zap.Error(err))
				}
				logger.For(ctx).Panic("application panic", zap.Reflect("error", r))
			}
		}()

		ctx = ent.NewContext(ctx, tx.Client())
		// Since the client is already uses transaction we can't have transactions on graphql also
		r := resolver.New(
			resolver.Config{
				Logger:     logger,
				Subscriber: event.NewNopSubscriber(),
			},
			resolver.WithTransaction(false),
		)

		if err := utilityFunc(ctx, r, logger, tenant); err != nil {
			logger.For(ctx).Error("failed to run function", zap.Error(err))
			if err := tx.Rollback(); err != nil {
				logger.For(ctx).Error("cannot rollback transaction", zap.Error(err))
			}
			return
		}

		if err := tx.Commit(); err != nil {
			logger.For(ctx).Error("cannot commit transaction", zap.Error(err))
			return
		}
	}
}

func getTenantList(ctx context.Context, driver *sql.Driver, tenant *string) ([]string, error) {
	if *tenant != "ALL" {
		return []string{*tenant}, nil
	}
	rows, err := driver.DB().QueryContext(ctx,
		"SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME LIKE ?", viewer.DBName("%"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tenants []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		name = viewer.FromDBName(name)
		tenants = append(tenants, name)
	}
	return tenants, nil
}

func utilityFunc(ctx context.Context, _ generated.ResolverRoot, logger log.Logger, tenant string) error {
	/**
	Add your Go code in this function
	You need to run this code from the same version production is at to avoid schema mismatches
	DO NOT LAND THE CODE AFTER THIS COMMENT
	*/
	/*
		Example code:
		client := ent.FromContext(ctx)
		eqt, err := r.Mutation().AddEquipmentType(ctx, models.AddEquipmentTypeInput{Name: "My new type"})
		if err != nil {
			return fmt.Errorf("cannot create equipment type: %w", err)
		}
		logger.For(ctx).Info("equipment created", zap.Int("ID", eqt.ID))
		client.EquipmentType.UpdateOneID(eqt.ID).SetName("My new type 2").ExecX(ctx)
		if err != nil {
			return fmt.Errorf("cannot update equipment type: id=%q, %w", eqt.ID, err)
		}
	*/
	return nil
}
