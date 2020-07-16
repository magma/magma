/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"context"

	"magma/orc8r/lib/go/protos"
	tcprotos "orc8r/fbinternal/cloud/go/services/testcontroller/protos"
	"orc8r/fbinternal/cloud/go/services/testcontroller/storage"

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
	nodes, err := n.store.GetNodes(req.Ids)
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

func (n *nodeLeasorServicer) LeaseNode(context.Context, *protos.Void) (*tcprotos.LeaseNodeResponse, error) {
	lease, err := n.store.LeaseNode()
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
