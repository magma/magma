/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package middleware/unary implements cloud service middleware layer which
// facilitates injection of cloudwide request & context decorators or filters
// (interceptors) for unary RPC methods
package unary

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// InterceptorHandler is a function type to intercept the execution of a unary
// RPC on the server.
// ctx, req & info contains all the information of this RPC the interceptor can
// operate on,
// If Handler returns an error, the chain of Interceptor calls will be
// interrupted and the error will be returned to the RPC client
// If returned CTX is not nil, it'll be used for the remaining interceptors and
// original RPC
// If resp return value is not nil - , the chain of Interceptor calls will be
// interrupted and the resp will be returned to the RPC client
type InterceptorHandler func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo) (
	newCtx context.Context, newReq interface{}, resp interface{}, err error)

// Inerceptor defines an interface to be implemented by all Unary Interceptors
// In addition to a receiver form of InterceptorHandler it provides Name &
// Description methods to aid diagnostic & logging of Interceptor related issues
type Interceptor interface {
	// Interceptor's Handler, has the same signature as
	// the non-receiver InterceptorHandler
	Handle(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo) (
		newCtx context.Context, newReq interface{}, resp interface{}, err error)
	// Name returns name of the Interceptor implementation
	Name() string
	// Description returns a string describing Interceptor
	Description() string
}

// interceptorFunc is an Interceptor interface adaptor for 'naked' functional
// interceptors
type interceptorFunc struct {
	handler InterceptorHandler
	name    string
}

func (adaptor *interceptorFunc) Handle(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) (
	newCtx context.Context, newReq interface{}, resp interface{}, err error) {
	return adaptor.handler(ctx, req, info)
}

func (adaptor *interceptorFunc) Name() string {
	return adaptor.name
}

func (adaptor *interceptorFunc) Description() string {
	return "Interceptor Adaptor For " + adaptor.name
}

func ListAllInterceptors() []string {
	var res = make([]string, len(registry))
	for i, unaryInterceptor := range registry {
		res[i] = unaryInterceptor.Name()
	}
	return res
}

// unary.MiddlewareHandler iterates through and calls all registered unary
// middleware interceptors and 'decorates' RPC parameters before invoking
// the original server RPC method
func MiddlewareHandler(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	for _, unaryInterceptor := range registry {
		newCtx, newReq, resp, err := unaryInterceptor.Handle(ctx, req, info)
		if err != nil {
			log.Printf(
				"Error %s from unary interceptor %s ",
				err, unaryInterceptor.Name())
			return resp, err
		}
		if resp != nil {
			return resp, err
		}
		if newCtx != nil {
			ctx = newCtx
		}
		if newReq != nil {
			req = newReq
		}
	}
	return handler(ctx, req)
}
