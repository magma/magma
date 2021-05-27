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

package test_init

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/protos"
	"magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/test_utils"
)

type indexerServicer struct {
	idx indexer.Indexer
}

// StartNewTestIndexer starts a new indexing service which forwards calls to the passed indexer.
func StartNewTestIndexer(t *testing.T, idx indexer.Indexer) {
	labels := map[string]string{
		orc8r.StateIndexerLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StateIndexerVersionAnnotation: strconv.Itoa(int(idx.GetVersion())),
		orc8r.StateIndexerTypesAnnotation:   strings.Join(idx.GetTypes(), orc8r.AnnotationFieldSeparator),
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, idx.GetID(), labels, annotations)
	servicer := &indexerServicer{idx: idx}
	protos.RegisterIndexerServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}

func (i *indexerServicer) Index(ctx context.Context, req *protos.IndexRequest) (*protos.IndexResponse, error) {
	states, err := types.MakeSerializedStatesByID(req.States)
	if err != nil {
		return nil, err
	}
	stErrs, err := i.idx.Index(req.NetworkId, states)
	res := &protos.IndexResponse{StateErrors: types.MakeProtoStateErrors(stErrs)}
	return res, err
}

func (i *indexerServicer) PrepareReindex(ctx context.Context, req *protos.PrepareReindexRequest) (*protos.PrepareReindexResponse, error) {
	err := i.idx.PrepareReindex(indexer.Version(req.FromVersion), indexer.Version(req.ToVersion), req.IsFirst)
	return &protos.PrepareReindexResponse{}, err
}

func (i *indexerServicer) CompleteReindex(ctx context.Context, req *protos.CompleteReindexRequest) (*protos.CompleteReindexResponse, error) {
	err := i.idx.CompleteReindex(indexer.Version(req.FromVersion), indexer.Version(req.ToVersion))
	return &protos.CompleteReindexResponse{}, err
}
