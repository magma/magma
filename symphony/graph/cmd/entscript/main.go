// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.HelpFlag.Short('h')
	dsn := kingpin.Flag("db-dsn", "data source name").Envar("MYSQL_DSN").Required().String()
	tenant := kingpin.Flag("tenant", "tenant name to target").Required().String()
	user := kingpin.Flag("user", "user name to target").Required().String()
	logcfg := log.AddFlags(kingpin.CommandLine)
	kingpin.Parse()

	logger, _, _ := log.Provide(*logcfg)
	ctx := context.Background()

	logger.For(ctx).Info("params",
		zap.Stringp("dsn", dsn),
		zap.Stringp("tenant", tenant),
		zap.Stringp("user", user),
	)
	tenancy, err := viewer.NewMySQLTenancy(*dsn)
	if err != nil {
		logger.For(ctx).Fatal("cannot connect to graph database",
			zap.Stringp("dsn", dsn),
			zap.Error(err),
		)
	}
	mysql.SetLogger(logger)

	v := &viewer.Viewer{Tenant: *tenant, User: *user}
	ctx = log.NewFieldsContext(ctx, zap.Object("viewer", v))
	ctx = viewer.NewContext(ctx, v)
	client, err := tenancy.ClientFor(ctx, *tenant)
	if err != nil {
		logger.For(ctx).Fatal("cannot get ent client for tenant",
			zap.Stringp("tenant", tenant),
			zap.Error(err),
		)
	}

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
			Emitter:    event.NewNopEmitter(),
			Subscriber: event.NewNopSubscriber(),
		},
		resolver.WithTransaction(false),
	)

	if err := utilityFunc(ctx, r, logger); err != nil {
		logger.For(ctx).Error("failed to run function", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			logger.For(ctx).Error("cannot rollback transaction", zap.Error(err))
		}
		return
	}

	if err := tx.Commit(); err != nil {
		logger.For(ctx).Error("cannot commit transaction", zap.Error(err))
	}
}

func utilityFunc(_ context.Context, _ generated.ResolverRoot, _ log.Logger) error {
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
		logger.For(ctx).Info("equipment created", zap.String("ID", eqt.ID))
		client.EquipmentType.UpdateOneID(eqt.ID).SetName("My new type 2").ExecX(ctx)
		if err != nil {
			return fmt.Errorf("cannot update equipment type: id=%q, %w", eqt.ID, err)
		}
	*/
	return nil
}
