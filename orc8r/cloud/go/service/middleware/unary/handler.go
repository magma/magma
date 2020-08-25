/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package unary implements cloud service middleware layer which
// facilitates injection of cloud-wide request & context decorators or filters
// (interceptors) for unary RPC methods
package unary

import (
	"runtime/debug"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var uncaughtCounterVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_handler_panic",
		Help: "There was a panic in the gateway",
	},
	[]string{"fullMethod"},
)

func init() {
	prometheus.MustRegister(uncaughtCounterVec)
}

var registry = []Interceptor{
	{
		Handler:     SetIdentityFromContext,
		Name:        "Unary Identity Decorator",
		Description: "Identity Decorator injects protos.Identity instance into RPC context",
	},
	{
		Handler:     BlockUnregisteredGateways,
		Name:        "BlockUnregisteredGateways",
		Description: "interceptor which blocks unregistered gateways from making RPC calls",
	},
}

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
type InterceptorHandler func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) (newCtx context.Context, newReq interface{}, resp interface{}, err error)

// Interceptor defines an interface to be implemented by all Unary Interceptors
// In addition to a receiver form of InterceptorHandler it provides Name &
// Description methods to aid diagnostic & logging of Interceptor related issues
type Interceptor struct {
	// Interceptor's Handler, has the same signature as
	// the non-receiver InterceptorHandler
	Handler func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) (newCtx context.Context, newReq interface{}, resp interface{}, err error)
	// Name returns name of the Interceptor implementation
	Name string
	// Description returns a string describing Interceptor
	Description string
}

// MiddlewareHandler iterates through and calls all registered unary
// middleware interceptors and 'decorates' RPC parameters before invoking
// the original server RPC method
func MiddlewareHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	for _, unaryInterceptor := range registry {
		newCtx, newReq, resp, err := unaryInterceptor.Handler(ctx, req, info)
		if err != nil {
			glog.Errorf("Error %s from unary interceptor %s ", err, unaryInterceptor.Name)
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

	resp, err = callHandler(ctx, req, info, handler)
	return
}

// callHandler wraps the handler with error recovery, logging, and metrics.
func callHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = status.Errorf(codes.Unknown, "handler panic: %s; stack trace: %s", r, debug.Stack())
			uncaughtCounterVec.WithLabelValues(info.FullMethod).Inc()
		}
	}()

	resp, err = handler(ctx, req)
	if err != nil {
		glog.Errorf("[ERROR %s]: %+v", info.FullMethod, err)
	}

	return
}
