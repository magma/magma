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

package servicers

import (
	"context"

	tcprotos "magma/fbinternal/cloud/go/services/testcontroller/protos"
	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type nodeLeasorServicer struct {
	store storage.NodeLeasorStorage
}

func NewNodeLeasorServicer(store storage.NodeLeasorStorage) tcprotos.NodeLeasorServer {
	return &nodeLeasorServicer{store: store}
}

func (n *nodeLeasorServicer) GetNodes(_ context.Context, req *tcprotos.GetNodesRequest) (*tcprotos.GetNodesResponse, error) {
	nodes, err := n.store.GetNodes(req.Ids, stringWrapperAsStrPtr(req.Tag))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tcprotos.GetNodesResponse{Nodes: nodes}, nil
}

func (n *nodeLeasorServicer) CreateOrUpdateNode(_ context.Context, req *tcprotos.CreateOrUpdateNodeRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if req.Node == nil {
		return ret, status.Error(codes.InvalidArgument, "node in request must be non-nil")
	}

	err := n.store.CreateOrUpdateNode(req.Node)
	if err != nil {
		return ret, status.Error(codes.Internal, err.Error())
	}
	return ret, nil
}

func (n *nodeLeasorServicer) DeleteNode(_ context.Context, req *tcprotos.DeleteNodeRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	err := n.store.DeleteNode(req.Id)
	if err != nil {
		return ret, status.Error(codes.Internal, err.Error())
	}
	return ret, nil
}

func (n *nodeLeasorServicer) ReserveNode(_ context.Context, req *tcprotos.ReserveNodeRequest) (*tcprotos.LeaseNodeResponse, error) {
	lease, err := n.store.ReserveNode(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tcprotos.LeaseNodeResponse{Lease: lease}, nil
}

func (n *nodeLeasorServicer) LeaseNode(_ context.Context, req *tcprotos.LeaseNodeRequest) (*tcprotos.LeaseNodeResponse, error) {
	lease, err := n.store.LeaseNode(strPtr(req.Tag))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tcprotos.LeaseNodeResponse{Lease: lease}, nil
}

func (n *nodeLeasorServicer) ReleaseNode(_ context.Context, req *tcprotos.ReleaseNodeRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	err := n.store.ReleaseNode(req.NodeID, req.LeaseID)
	switch {
	case err == nil:
		return ret, nil
	case err == storage.ErrBadRelease:
		return ret, status.Error(codes.InvalidArgument, err.Error())
	default:
		return ret, status.Error(codes.Internal, err.Error())
	}
}

func stringWrapperAsStrPtr(s *wrappers.StringValue) *string {
	if s == nil {
		return nil
	}
	return strPtr(s.GetValue())
}

func strPtr(s string) *string {
	return &s
}
