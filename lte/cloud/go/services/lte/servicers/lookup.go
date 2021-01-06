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

	"magma/lte/cloud/go/services/lte/protos"
	lte_storage "magma/lte/cloud/go/services/lte/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// lookupServicer stores reported enodeb state with additional gatewayID as
// a primary key.
type lookupServicer struct {
	store lte_storage.EnodebStateLookup
}

// NewLookupServicer returns a new enodeb lookup servicer.
// Stores should be initialized by the caller.
func NewLookupServicer(enbStore lte_storage.EnodebStateLookup) protos.EnodebStateLookupServer {
	return &lookupServicer{store: enbStore}
}

func (l *lookupServicer) GetEnodebState(ctx context.Context, req *protos.GetEnodebStateRequest) (*protos.GetEnodebStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	enbState, err := l.store.GetEnodebState(req.NetworkId, req.GatewayId, req.EnodebSn)
	if err != nil {
		return nil, makeErr(err, "get eNB state from store")
	}

	res := &protos.GetEnodebStateResponse{SerializedState: enbState}
	return res, nil
}

func (l *lookupServicer) SetEnodebState(ctx context.Context, req *protos.SetEnodebStateRequest) (*protos.SetEnodebStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := l.store.SetEnodebState(req.NetworkId, req.GatewayId, req.EnodebSn, req.SerializedState)
	if err != nil {
		return nil, makeErr(err, "write eNB state to store")
	}
	return &protos.SetEnodebStateResponse{}, nil
}

func makeErr(err error, wrap string) error {
	e := errors.Wrap(err, wrap)
	code := codes.Internal
	if err == merrors.ErrNotFound {
		code = codes.NotFound
	}
	return status.Error(code, e.Error())
}
