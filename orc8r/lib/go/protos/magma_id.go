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

package protos

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// X_MAGMA_ID is a request/response tracing header which must be carried between RPC calls
const X_MAGMA_ID = "x-magma-id"

// AddRequestId creates and returns new outgoing GRPC CTX with specified Id's header
// Intended for GRPC client use
func AddRequestId(ctx context.Context, id string) context.Context {
	md := metadata.New(map[string]string{X_MAGMA_ID: id})
	return metadata.NewOutgoingContext(ctx, md)
}

// ForwardRequestId fetches the ID from incoming GRPC CTX and appends it to the server's response
// Should be called immediately before return from RPC handler
func ForwardRequestId(ctx context.Context) error {
	if ctx == nil {
		return status.Errorf(codes.InvalidArgument, "nil '%s' ctx", X_MAGMA_ID)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || md == nil {
		return status.Errorf(codes.DataLoss, "missing metadata for %s", X_MAGMA_ID)
	}
	hdrs := md[X_MAGMA_ID]
	if len(hdrs) == 0 {
		return status.Errorf(codes.InvalidArgument, "missing '%s' header", X_MAGMA_ID)
	}
	id := strings.TrimSpace(hdrs[0])
	if len(id) == 0 {
		return status.Errorf(codes.InvalidArgument, "empty '%s' header", X_MAGMA_ID)
	}
	return SendRequestId(ctx, id)
}

// SendRequestId sends given ID with GRPC CTX to the client
// Should be called immediately before return from RPC handler
func SendRequestId(ctx context.Context, id string) error {
	if ctx == nil {
		return status.Errorf(codes.InvalidArgument, "nil '%s' ctx", X_MAGMA_ID)
	}
	header := metadata.New(map[string]string{X_MAGMA_ID: id})
	err := grpc.SendHeader(ctx, header)
	if err != nil {
		err = status.Errorf(codes.Internal, "error sending '%s' header: %v", X_MAGMA_ID, err)
	}
	return err
}
