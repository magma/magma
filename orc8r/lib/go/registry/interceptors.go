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

// Package registry for Magma microservices
package registry

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"magma/orc8r/lib/go/service/middleware/unary"
)

var defaultTimeoutDuration = GrpcMaxTimeoutSec * time.Second

// TimeoutInterceptor is a generic client connection interceptor which sets default timeout option for RPC if the
// currently used CTX does not already specify it's own deadline option
func TimeoutInterceptor(ctx context.Context, method string, req, resp interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	// check if given CTX already has a deadline & only add default deadline if not
	if _, deadlineIsSet := ctx.Deadline(); !deadlineIsSet {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimeoutDuration)
		// cleanup timer after invoke call chain completion
		defer cancel()
	}
	return invoker(ctx, method, req, resp, cc, opts...)
}

// CloudClientTimeoutInterceptor is the same as TimeoutInterceptor,
// but - it also amends outgoing CTX with magic cloud client CSN header
func CloudClientTimeoutInterceptor(ctx context.Context, method string, req, resp interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	return TimeoutInterceptor(unary.OutgoingCloudClientCtx(ctx), method, req, resp, cc, invoker, opts...)
}
