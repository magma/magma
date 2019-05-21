/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package checkind

import (
	"context"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ServiceName = "CHECKIND"

func getCheckindClient() (protos.CheckindClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewCheckindClient(conn), err
}

// GetStatus returns the gateway status for the gateway with logicalID in the network specified by networkID
func GetStatus(networkID string, logicalID string) (*protos.GatewayStatus, error) {
	client, err := getCheckindClient()
	if err != nil {
		return nil, err
	}

	ret, err := client.GetStatus(context.Background(), &protos.GatewayStatusRequest{
		NetworkId: networkID,
		LogicalId: logicalID,
	})
	// Special handling for 404. If err == nil, note that status.Code returns
	// codes.OK
	switch status.Code(err) {
	case codes.NotFound:
		return nil, errors.ErrNotFound
	default:
		return ret, err
	}
}

// DeleteGatewayStatus removes the gateway status record from the gateway's network table
// NOTE: the record will be created again after next successful gateway checkin
func DeleteGatewayStatus(networkID string, logicalID string) error {
	client, err := getCheckindClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteGatewayStatus(context.Background(), &protos.GatewayStatusRequest{
		NetworkId: networkID,
		LogicalId: logicalID,
	})
	return err
}

// DeleteNetwork deletes the network's status table. All gateway statuses must
// be deleted prior to removal.
func DeleteNetwork(networkID string) error {
	client, err := getCheckindClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNetwork(context.Background(), &protos.NetworkID{Id: networkID})
	return err
}

// List returns a list of all logical gateway IDs for the given network which
// have been stored in the service DB
func List(networkID string) (*protos.IDList, error) {
	client, err := getCheckindClient()
	if err != nil {
		return nil, err
	}
	return client.List(context.Background(), &protos.NetworkID{Id: networkID})
}
