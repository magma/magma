/*
 * Copyright 2021 The Magma Authors.
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

// package unary implements all cloud service framework unary interceptors
package unary

import (
	"context"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	unarylib "magma/orc8r/lib/go/service/middleware/unary"
)

// EnforceIntraCloudRPC verifies that the RPC originated from a service on the cloud (and not a remote Gateway),
// it'll will also bypass the checks for allowed RPCs (methods in identityDecoratorBypassList)
func EnforceIntraCloudRPC(ctx context.Context, info *grpc.UnaryServerInfo) error {
	if err := ensureLocalPeer(ctx); err == nil {
		// a call from local peer - allow
		return nil
	}
	if info != nil {
		// Check if the call is for a globally allowed/bypassed method (bootstrapper, etc.)
		if _, ok := identityDecoratorBypassList[info.FullMethod]; ok {
			// Bypass method (Bootstrapper & Co.) - allow
			return nil
		}
	}
	ctxMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok || ctxMetadata == nil {
		// no metadata - fail
		glog.Error(ERROR_MSG_NO_METADATA)
		return status.Error(codes.Unauthenticated, ERROR_MSG_NO_METADATA)
	}
	// First, try to find the caller's identity
	snlist, ok := ctxMetadata[CLIENT_CERT_SN_KEY]
	if !ok || len(snlist) != 1 {
		// there can be only one CSN, error out if not
		glog.Errorf("Single CSN is expected in metadata: %+v", ctxMetadata)
		return status.Error(codes.InvalidArgument, "Single CSN is expected")
	}
	if snlist[0] != unarylib.ORC8R_CLIENT_CERT_VALUE {
		// Not an inter-service CSN - fail
		glog.Errorf("Non inter-service CSN in request to internal service: %s", snlist[0])
		return status.Error(codes.PermissionDenied, "Invalid CSN in internal RPC")
	}
	return nil
}
