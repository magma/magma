// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"github.com/facebookincubator/symphony/cloud/log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newServer(logger log.Logger, srv TenantServiceServer) (*grpc.Server, func(), error) {
	grpc_zap.ReplaceGrpcLogger(logger.Background())

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger.Background()),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger.Background()),
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)

	RegisterTenantServiceServer(s, srv)
	reflection.Register(s)

	err := view.Register(ocgrpc.DefaultServerViews...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "registering grpc views")
	}
	return s, func() { view.Unregister(ocgrpc.DefaultServerViews...) }, nil
}
