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
	"context"
	"runtime/debug"

	"github.com/golang/glog"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
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

var interceptor = grpc.ChainUnaryInterceptor(
	errlogInterceptor,
	recoveryInterceptor,
	grpc_prometheus.UnaryServerInterceptor,
	makeInterceptorFromTemplate(SetIdentityFromContext),
	makeInterceptorFromTemplate(BlockUnregisteredGateways),
)

func GetInterceptorOpt() grpc.ServerOption {
	return interceptor
}

// recoveryInterceptor converts handler panics to gRPC errors.
// Ref: https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/recovery/interceptors.go.
func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	panicked := true
	defer func() {
		if r := recover(); r != nil || panicked {
			err = status.Errorf(codes.Unknown, "handler panic: %s; stack trace: %s", r, debug.Stack())
			uncaughtCounterVec.WithLabelValues(info.FullMethod).Inc()
		}
	}()
	res, err = handler(ctx, req)
	panicked = false
	return res, err
}

// errlogInterceptor logs errors when gRPC handlers return errors.
func errlogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		glog.Errorf("[ERROR %s]: %+v", info.FullMethod, err)
	}
	return res, err
}

type interceptorTemplate func(context.Context, interface{}, *grpc.UnaryServerInfo) (context.Context, interface{}, interface{}, error)

func makeInterceptorFromTemplate(fn interceptorTemplate) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx, _, res, err := fn(ctx, req, info)
		if err != nil {
			return res, err
		}
		if newCtx != nil {
			ctx = newCtx
		}
		return handler(ctx, req)
	}
}
