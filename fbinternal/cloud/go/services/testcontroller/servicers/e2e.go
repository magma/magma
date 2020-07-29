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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type e2eServicer struct {
	store storage.TestControllerStorage
}

func NewTestControllerServicer(store storage.TestControllerStorage) tcprotos.TestControllerServer {
	return &e2eServicer{store: store}
}

func (e *e2eServicer) GetTestCases(_ context.Context, req *tcprotos.GetTestCasesRequest) (*tcprotos.GetTestCasesResponse, error) {
	tcs, err := e.store.GetTestCases(req.Pks)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tcprotos.GetTestCasesResponse{Tests: tcs}, nil
}

func (e *e2eServicer) CreateOrUpdateTestCase(_ context.Context, req *tcprotos.CreateTestCaseRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if req.Test == nil {
		return ret, status.Error(codes.InvalidArgument, "test case in request must be non-nil")
	}
	err := e.store.CreateOrUpdateTestCase(req.Test)
	if err != nil {
		return ret, status.Error(codes.Internal, err.Error())
	}
	return ret, nil
}

func (e *e2eServicer) DeleteTestCase(_ context.Context, req *tcprotos.DeleteTestCaseRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	err := e.store.DeleteTestCase(req.Pk)
	if err != nil {
		return ret, status.Error(codes.Internal, err.Error())
	}
	return ret, nil
}
