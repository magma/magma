/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"context"
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
func StartNewTestIndexer(t *testing.T, serviceName string, idx indexer.Indexer) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, serviceName)
	servicer := &indexerServicer{idx}
	protos.RegisterIndexerServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}

func (i *indexerServicer) GetIndexerInfo(ctx context.Context, req *protos.GetIndexerInfoRequest) (*protos.GetIndexerInfoResponse, error) {
	res := &protos.GetIndexerInfoResponse{Version: uint32(i.idx.GetVersion()), StateTypes: i.idx.GetTypes()}
	return res, nil
}

func (i *indexerServicer) Index(ctx context.Context, req *protos.IndexRequest) (*protos.IndexResponse, error) {
	states, err := types.MakeStatesByID(req.States)
	if err != nil {
		return nil, err
	}
	stErrs, err := i.idx.Index(req.NetworkId, states)
	res := &protos.IndexResponse{StateErrors: types.MakeProtoStateErrors(stErrs)}
	return res, nil
}

func (i *indexerServicer) PrepareReindex(ctx context.Context, req *protos.PrepareReindexRequest) (*protos.PrepareReindexResponse, error) {
	err := i.idx.PrepareReindex(indexer.Version(req.FromVersion), indexer.Version(req.ToVersion), req.IsFirst)
	return &protos.PrepareReindexResponse{}, err
}

func (i *indexerServicer) CompleteReindex(ctx context.Context, req *protos.CompleteReindexRequest) (*protos.CompleteReindexResponse, error) {
	err := i.idx.CompleteReindex(indexer.Version(req.FromVersion), indexer.Version(req.ToVersion))
	return &protos.CompleteReindexResponse{}, err
}
