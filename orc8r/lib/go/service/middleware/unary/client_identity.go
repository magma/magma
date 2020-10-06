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

func UnaryCloudClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	appendedCtx := metadata.AppendToOutgoingContext(ctx, CLIENT_CERT_SN_KEY, ORC8R_CLIENT_CERT_VALUE)
	err := invoker(appendedCtx, method, req, reply, cc, opts...)
	return err
}
