/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package unary

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// Client Certificate Serial Number Header
	CLIENT_CERT_SN_KEY = "x-magma-client-cert-serial"

	// Magic Certificate Value for orc8r service clients
	ORC8R_CLIENT_CERT_VALUE = "7ZZXAF7CAETF241KL22B8YRR7B5UF401"
)

// CloudClientInterceptor sets Magic Certificate Value for orc8r service clients in the outgoing CTX
// if the CTX metadata already has a client certificate SN key, CloudClientInterceptor will overwrite it with
// the Magic Certificate Value
func CloudClientInterceptor(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	return invoker(OutgoingCloudClientCtx(ctx), method, req, reply, cc, opts...)
}

// OutgoingCloudClientCtx amends outgoing cloud client context with magic client CSN metadata
// if the original context doesn't already have client certificate serial number metadata
func OutgoingCloudClientCtx(ctx context.Context) context.Context {
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
