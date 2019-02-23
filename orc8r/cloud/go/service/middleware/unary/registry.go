/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package middleware implements cloud service middleware layer which
// facilitates injection of cloudwide request & context decorators or filters
// (interceptors)
package unary

import (
	"log"
	"reflect"
	"runtime"

	"magma/orc8r/cloud/go/service/middleware/unary/interceptors"
)

// unaryRegistry is a list of all registered unary Interceptors, they can be
// either objects satisfying Interceptor interface OR functions of type
// InterceptorHandler
// For example:
// 	a handler in magma/service/middleware/unary/interceptors/identity_decorator.go
//	can be included as either
//		var unaryRegistry = []interface{} {
//			InterceptorHandler(interceptors.SetIdentityFromContext),
// 		...
// 	or
//		var unaryRegistry = []interface{} {
//			interceptors.NewIdentityDecorator(),
// 		...
var unaryRegistry = []interface{}{
	interceptors.NewIdentityDecorator(),
	InterceptorHandler(interceptors.BlockUnregisteredGateways),
}

// registry lists all Unary Interceptors invoked on every service RPC call
var registry []Interceptor

func init() {
	registry = make([]Interceptor, 0, len(unaryRegistry))
	for idx, x := range unaryRegistry {
		if x == nil {
			log.Printf("Nil Unary Interceptor in %d position", idx)
			continue
		}
		if interceptor, ok := x.(Interceptor); ok {
			registry = append(registry, interceptor)
		} else if handlerFunc, ok := x.(InterceptorHandler); ok {
			adaptor := &interceptorFunc{
				handlerFunc,
				runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name(),
			}
			registry = append(registry, adaptor)
		} else {
			log.Printf("Invalid Unary Interceptor Type %T at %d position", x, idx)
		}
	}
}
