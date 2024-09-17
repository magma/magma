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

package test_utils

import (
	"context"
	"time"

	"google.golang.org/grpc"

	lib_protos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

// GetConnectionWithAuthority provides a gRPC connection to a service in the registry with Authority header.
func GetConnectionWithAuthority(service string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
	defer cancel()
	return registry.GetConnectionImpl(ctx, service, lib_protos.ServiceType_SOUTHBOUND, grpc.WithBackoffMaxDelay(registry.GrpcMaxDelaySec*time.Second), grpc.WithBlock(), grpc.WithAuthority(service))
}
