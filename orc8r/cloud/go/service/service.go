/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package service outlines the Magma microservices framework in the cloud.
// The framework helps to create a microservice easily, and provides
// the common service logic like service303, config, etc.
package service

import (
	"flag"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/service/middleware/unary"
	platform_service "magma/orc8r/lib/go/service"

	"google.golang.org/grpc"
)

// NewOrchestratorService returns a new GRPC orchestrator service
// implementing service303. This service will implement a middleware
// interceptor to perform identity check. If your service does not or can not
// perform identity checks, (e.g. federation), use NewServiceWithOptions.
func NewOrchestratorService(moduleName string, serviceName string) (*platform_service.Service, error) {
	flag.Parse()
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})
	return platform_service.NewServiceWithOptionsImpl(moduleName, serviceName, grpc.UnaryInterceptor(unary.MiddlewareHandler))
}

// NewOrchestratorServiceWithOptions returns a new GRPC orchestrator service
// implementing service303 with the specified grpc server options. This service
// will implement a middleware interceptor to perform identity check.
func NewOrchestratorServiceWithOptions(moduleName string, serviceName string, serverOptions ...grpc.ServerOption) (*platform_service.Service, error) {
	flag.Parse()
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})
	serverOptions = append(serverOptions, grpc.UnaryInterceptor(unary.MiddlewareHandler))
	return platform_service.NewServiceWithOptionsImpl(moduleName, serviceName, serverOptions...)
}
