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

package servicers_test

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/services/state/indexer"
	reindex_mocks "magma/orc8r/cloud/go/services/state/indexer/reindex/mocks"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	state_proto_mocks "magma/orc8r/cloud/go/services/state/protos/mocks"
	"magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	id0 = "SOME_INDEXERID_0"
	id1 = "SOME_INDEXERID_1"
	id2 = "SOME_INDEXERID_2"

	zero      = 0
	version0  = 10
	version1  = 20
	version1a = 200
	version2  = 30
)

var (
	someErr = errors.New("some_error")

	ctxBlank        = context.Background()
	ctxWithIdentity = protos.
			NewGatewayIdentity("some_hwid", "some_nwid", "some_logicalid").
			NewContextWithIdentity(ctxBlank)
)

func TestIndexerServicer_GetIndexers(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		composed := []*indexer.Versions{
			{IndexerID: id0, Actual: zero, Desired: version0},
			{IndexerID: id1, Actual: version1, Desired: version1a},
			{IndexerID: id2, Actual: version2, Desired: version2},
		}
		asProtos := map[string]*indexer_protos.IndexerInfo{
			id0: {IndexerId: id0, ActualVersion: zero, DesiredVersion: version0},
			id1: {IndexerId: id1, ActualVersion: version1, DesiredVersion: version1a},
			id2: {IndexerId: id2, ActualVersion: version2, DesiredVersion: version2},
		}

		r := &reindex_mocks.Reindexer{}
		r.On("GetIndexerVersions").Return(composed, nil)
		srv := servicers.NewIndexerManagerServicer(r, false)

		got, err := srv.GetIndexers(ctxBlank, &indexer_protos.GetIndexersRequest{})
		assert.NoError(t, err)
		want := &indexer_protos.GetIndexersResponse{IndexersById: asProtos}
		assert.Equal(t, want, got)
	})

	t.Run("fail when blankCtx identity is present", func(t *testing.T) {
		srv := servicers.NewIndexerManagerServicer(&reindex_mocks.Reindexer{}, false)

		_, err := srv.GetIndexers(ctxWithIdentity, &indexer_protos.GetIndexersRequest{})
		assert.Error(t, err)
		e, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, e.Code())
	})

	t.Run("reindexer err", func(t *testing.T) {
		r := &reindex_mocks.Reindexer{}
		r.On("GetIndexerVersions").Return(nil, someErr)
		srv := servicers.NewIndexerManagerServicer(r, false)

		_, err := srv.GetIndexers(ctxBlank, &indexer_protos.GetIndexersRequest{})
		assert.Error(t, err)
	})
}

func TestIndexerServicer_StartReindex(t *testing.T) {
	stream := &state_proto_mocks.IndexerManager_StartReindexServer{}
	stream.On("Send", mock.Anything).Return(nil)
	stream.On("Context").Return(ctxBlank)

	// reindex idx0
	t.Run("reindex one", func(t *testing.T) {
		r := &reindex_mocks.Reindexer{}
		r.On("RunUnsafe", ctxBlank, id0, mock.Anything).Return(nil)

		srv := servicers.NewIndexerManagerServicer(r, false)
		err := srv.StartReindex(&indexer_protos.StartReindexRequest{IndexerId: id0}, stream)
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	// reindex idx1 and idx2
	t.Run("reindex multiple with override", func(t *testing.T) {
		r := &reindex_mocks.Reindexer{}
		r.On("RunUnsafe", ctxBlank, "", mock.Anything).Return(nil)
		srv := servicers.NewIndexerManagerServicer(r, true)

		err := srv.StartReindex(&indexer_protos.StartReindexRequest{IndexerId: "", Force: true}, stream)
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("fail when auto reindex enabled", func(t *testing.T) {
		srv := servicers.NewIndexerManagerServicer(&reindex_mocks.Reindexer{}, true)

		err := srv.StartReindex(&indexer_protos.StartReindexRequest{Force: false}, stream)
		assert.Error(t, err)
		e, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.FailedPrecondition, e.Code())
	})

	t.Run("fail when blankCtx identity is present", func(t *testing.T) {
		stream := &state_proto_mocks.IndexerManager_StartReindexServer{}
		stream.On("Context").Return(ctxWithIdentity)
		srv := servicers.NewIndexerManagerServicer(&reindex_mocks.Reindexer{}, true)

		err := srv.StartReindex(&indexer_protos.StartReindexRequest{}, stream)
		assert.Error(t, err)
		e, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, e.Code())
	})

	t.Run("reindexer err", func(t *testing.T) {
		r := &reindex_mocks.Reindexer{}
		r.On("RunUnsafe", ctxBlank, "", mock.Anything).Return(someErr)
		srv := servicers.NewIndexerManagerServicer(r, false)

		err := srv.StartReindex(&indexer_protos.StartReindexRequest{}, stream)
		assert.Error(t, err)
		r.AssertExpectations(t)
	})
}
