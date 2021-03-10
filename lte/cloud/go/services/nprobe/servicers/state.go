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

package servicers

import (
	"context"

	"magma/lte/cloud/go/services/nprobe/protos"
	"magma/lte/cloud/go/services/nprobe/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// nProbeStateServicer manages nprobe states in a blobstore
// state is stored as a SQL table with two string columns
//	- LastExported
//	- SequenceNumber
type nProbeStateServicer struct {
	store storage.NProbeStateService
}

// NewNProbeStateServicer returns a new nprobe state servicer.
// Stores should be initialized by the caller.
func NewNProbeStateServicer(store storage.NProbeStateService) protos.NProbeStateServiceServer {
	return &nProbeStateServicer{store: store}
}

func (s *nProbeStateServicer) GetNProbeState(ctx context.Context, req *protos.GetNProbeStateRequest) (*protos.GetNProbeStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	state, err := s.store.GetNProbeState(req.NetworkId, req.TaskId)
	if err != nil {
		return nil, makeErr(err, "get state from store")
	}

	res := &protos.GetNProbeStateResponse{LastExported: state.LastExported, SequenceNumber: state.SequenceNumber}
	return res, nil
}

func (s *nProbeStateServicer) SetNProbeState(ctx context.Context, req *protos.SetNProbeStateRequest) (*protos.SetNProbeStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := s.store.SetNProbeState(req.NetworkId, req.TaskId, req.TargetId, req.NprobeState)
	if err != nil {
		return nil, makeErr(err, "set state to store")
	}

	res := &protos.SetNProbeStateResponse{}
	return res, nil
}

func (s *nProbeStateServicer) DeleteNProbeState(ctx context.Context, req *protos.DeleteNProbeStateRequest) (*protos.DeleteNProbeStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := s.store.DeleteNProbeState(req.NetworkId, req.TaskId)
	if err != nil {
		return nil, makeErr(err, "delete state from store")
	}
	return &protos.DeleteNProbeStateResponse{}, nil
}

func makeErr(err error, wrap string) error {
	e := errors.Wrap(err, wrap)
	code := codes.Internal
	if err == merrors.ErrNotFound {
		code = codes.NotFound
	}
	return status.Error(code, e.Error())
}
