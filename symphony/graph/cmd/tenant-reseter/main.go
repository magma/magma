// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	endpoint := flag.String("target", "", "graph target address")
	tenant := flag.String("tenant", "", "tenant to reset")
	force := flag.Bool("force", false, "force tenant reset")
	flag.Parse()

	logger, _ := zap.NewProduction()
	if !*force && !strings.Contains(strings.ToLower(*tenant), "test") {
		logger.Fatal("reset of non test tenant requires -force")
	}

	conn, err := grpc.Dial(*endpoint, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect to target", zap.Error(err))
	}
	client := graphgrpc.NewTenantServiceClient(conn)

	logger = logger.With(zap.String("name", *tenant))
	value := &wrappers.StringValue{Value: *tenant}
	if _, err := client.Delete(context.Background(), value); err != nil {
		log.Fatal("cannot delete tenant")
	}
	if _, err := client.Create(context.Background(), value); err != nil {
		log.Fatal("cannot create tenant")
	}
	logger.Info("tenant reset finished successfully")
}
