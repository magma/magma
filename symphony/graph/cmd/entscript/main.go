// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"

	"go.uber.org/zap"
)

type cliFlags struct {
	Dsn    string `env:"MYSQL_DSN" long:"dsn" description:"data source name"`
	Tenant string `long:"tenant" required:"true" description:"target specific tenant"`
	User   string `long:"user" required:"true" description:"target specific user"`
}

func main() {
	logger, _ := log.Config{Format: "console"}.Build()
	ctx := context.Background()

	var cf cliFlags
	if _, err := flags.Parse(&cf); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	logger.For(ctx).Info("params", zap.String("dsn", cf.Dsn), zap.String("tenant", cf.Tenant), zap.String("user", cf.User))

	tenancy, err := viewer.NewMySQLTenancy(cf.Dsn)
	if err != nil {
		logger.For(ctx).Fatal("cannot connect to graph database", zap.String("dsn", cf.Dsn), zap.Error(err))
		return
	}

	mysql.SetLogger(logger)

	v := &viewer.Viewer{Tenant: cf.Tenant, User: cf.User}

	ctx = log.NewFieldsContext(ctx, zap.Object("viewer", v))
	ctx = viewer.NewContext(ctx, v)

	client, err := tenancy.ClientFor(ctx, cf.Tenant)
	if err != nil {
		logger.For(ctx).Fatal("cannot get ent client for tenant", zap.String("tenant", cf.Tenant), zap.Error(err))
		return
	}
	ctx = ent.NewContext(ctx, client)
	utilityFunc(ctx, client, logger)
}

func utilityFunc(ctx context.Context, client *ent.Client, logger log.Logger) {
	/**
	Add your Go code in this function
	You need to run this code from the same version production is at to avoid schema mismatches
	DO NOT LAND THE CODE AFTER THIS COMMENT
	*/
	/*
		Example code:
		count, err := client.EquipmentPosition.Delete().Where(equipmentposition.ID("30064771558")).Exec(ctx)
		if err != nil {
			logger.For(ctx).Fatal("failed to delete equipment position", zap.String("ID", "30064771558"))
		}
		logger.For(ctx).Info("equipment position deleted", zap.Int("count", count))
	*/
}
