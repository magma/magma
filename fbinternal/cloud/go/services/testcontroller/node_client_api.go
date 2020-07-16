/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package testcontroller

import (
	"context"

	merrors "magma/orc8r/lib/go/errors"
	protos2 "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
	"orc8r/fbinternal/cloud/go/services/testcontroller/protos"
	"orc8r/fbinternal/cloud/go/services/testcontroller/storage"

	"github.com/golang/glog"
)

func getNodeClient() (protos.NodeLeasorClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewNodeLeasorClient(conn), nil
}

func GetNodes(ids []string) (map[string]*storage.CINode, error) {
	client, err := getNodeClient()
	if err != nil {
		return nil, err
	}
	res, err := client.GetNodes(context.Background(), &protos.GetNodesRequest{Ids: ids})
	if err != nil {
		return nil, err
	}
	return res.Nodes, nil
}

func CreateOrUpdateNode(node *storage.MutableCINode) error {
	client, err := getNodeClient()
	if err != nil {
		return err
	}
	_, err = client.CreateOrUpdateNode(context.Background(), &protos.CreateOrUpdateNodeRequest{Node: node})
	return err
}

func DeleteNode(id string) error {
	client, err := getNodeClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNode(context.Background(), &protos.DeleteNodeRequest{Id: id})
	return err
}

func ReserveNode(id string) (*storage.NodeLease, error) {
	client, err := getNodeClient()
	if err != nil {
		return nil, err
	}
	res, err := client.ReserveNode(context.Background(), &protos.ReserveNodeRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return res.Lease, nil
}

func LeaseNode() (*storage.NodeLease, error) {
	client, err := getNodeClient()
	if err != nil {
		return nil, err
	}
	res, err := client.LeaseNode(context.Background(), &protos2.Void{})
	if err != nil {
		return nil, err
	}
	return res.Lease, nil
}

func ReleaseNode(id string, leaseID string) error {
	client, err := getNodeClient()
	if err != nil {
		return err
	}
	_, err = client.ReleaseNode(context.Background(), &protos.ReleaseNodeRequest{NodeID: id, LeaseID: leaseID})
	return err
}
