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
	"google.golang.org/grpc/metadata"
)

const (
	// Client Certificate Serial Number Header
	CLIENT_CERT_SN_KEY = "x-magma-client-cert-serial"

	// Magic Certificate Value for orc8r service clients
	ORC8R_CLIENT_CERT_VALUE = "7ZZXAF7CAETF241KL22B8YRR7B5UF401"
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

	return TimeoutInterceptor(outgoingCloudClientCtx(ctx), method, req, resp, cc, invoker, opts...)
}

// outgoingCloudClientCtx amends outgoing cloud client context with magic client CSN metadata
// if the original context doesn't already have client certificate serial number metadata
func outgoingCloudClientCtx(ctx context.Context) context.Context {
	md, exists := metadata.FromOutgoingContext(ctx)
	if exists {
		if sns := md.Get(CLIENT_CERT_SN_KEY); len(sns) == 1 && len(sns[0]) > 0 {
			// Do not alter outgoing CTX if it already has client certificate serial header
			// this should allow dual use services (serving external as well as internal clients) to be called
			// internally passing external client cert SN (see directoryd update servicer)
			return ctx
		}
		md = md.Copy()
		md.Set(CLIENT_CERT_SN_KEY, ORC8R_CLIENT_CERT_VALUE)
	} else {
		md = metadata.Pairs(CLIENT_CERT_SN_KEY, ORC8R_CLIENT_CERT_VALUE)
	}
	return metadata.NewOutgoingContext(ctx, md)
}
