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

// package interceptors implements all cloud service framework unary interceptors
package unary

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/lib/go/protos"
)

// BlockUnregisteredGateways is an Interceptor blocking calls from Gateways
// which were not registered on the cloud.
// BlockUnregisteredGateways must be invoked after Identity Decorator since
// it relies on the Identity Decorator's results
func BlockUnregisteredGateways(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo) (
	newCtx context.Context, newReq interface{}, resp interface{}, err error,
) {
	gw := protos.GetClientGateway(ctx)
	if gw != nil && !gw.Registered() {
		var rpc string
		if info != nil {
			rpc = info.FullMethod
		} else {
			rpc = "Undefined"
		}
		log.Printf("Blocking %s call from unregisterd Gateway %+v", rpc, gw)
		err = status.Errorf(codes.PermissionDenied, "Unregistered Gateway %s", gw.GetHardwareId())
	}
	return
}
