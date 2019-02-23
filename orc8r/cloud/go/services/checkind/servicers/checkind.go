/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
Checkind service provides the gRPC interface for the Gateways, REST server to
collect/report and retreive Gateway runtime statuses.

The service relies on it's own DBs & tables for storage Gateway statuses and
magmad API for mapping Gateway IDs & networks.

It's not intended for storage of Gateway configs or any long term persistent
data/configuration .
*/
package servicers

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind/store"
)

type checkindServer struct {
	Store *store.CheckinStore
}

func NewCheckindServer(store *store.CheckinStore) (*checkindServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Cannot initialize Checkin Server with Nil store")
	}
	return &checkindServer{store}, nil
}

// Gateway periodic checkin - records given GW status into the GW's network table
// Checkin RPC is used by registered Gateways to report their real time status
// A Gateway identifies itself by its hardware ID (in CheckinRequest), it's the
// responsibility of the service to associate the provided HW ID with appropriate
// network table & corresponding GW Logical ID.
// Checkin requests for unregistered Gateways will result in error
func (srv *checkindServer) Checkin(ctx context.Context, req *protos.CheckinRequest) (*protos.CheckinResponse, error) {
	respPtr := &protos.CheckinResponse{
		Action: protos.CheckinResponse_NONE,
		Time:   uint64(time.Now().UnixNano()) / uint64(time.Millisecond)}
	if req == nil {
		return respPtr, fmt.Errorf("Nil CheckinRequest")
	}
	// Get gateway id from context
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		return respPtr, status.Errorf(codes.PermissionDenied, "Missing Gateway Identity")
	}
	if !gw.Registered() {
		return respPtr, status.Errorf(codes.PermissionDenied, "Gateway is not registered")
	}
	// Overwrite GW ID with verified Id From ctx
	req.GatewayId = gw.HardwareId
	err := srv.Store.UpdateRegisteredGatewayStatus(
		gw.GetNetworkId(),
		gw.GetLogicalId(),
		&protos.GatewayStatus{
			Time:               respPtr.Time,
			Checkin:            req,
			CertExpirationTime: protos.GetClientCertExpiration(ctx),
		},
	)
	if err != nil {
		err = fmt.Errorf("Update Gateway Status Error: '%s' for Gateway: %s", err, req.GatewayId)
	}
	return respPtr, err
}

// Gateway real time status retrieval from the GW's network table
// Gatway status can only be queried using Network ID & Gateway Logical ID of
// the network the Gateway is registered on.
// Status requests for unregistered Gateways will result in error
func (srv *checkindServer) GetStatus(ctx context.Context, req *protos.GatewayStatusRequest) (*protos.GatewayStatus, error) {
	if req == nil {
		return new(protos.GatewayStatus), fmt.Errorf("Nil GatewayStatusRequest")
	}
	ret, err := srv.Store.GetGatewayStatus(req)
	if err == store.ErrNotFound {
		return nil, status.Errorf(codes.NotFound, "No status found")
	}
	return ret, err
}

// Removes Gateway status record from the Gateway's network table
// NOTE: the record will be created again after next successfull Gateway checkin
func (srv *checkindServer) DeleteGatewayStatus(ctx context.Context, req *protos.GatewayStatusRequest) (*protos.Void, error) {
	if req == nil {
		return &protos.Void{}, fmt.Errorf("Nil GatewayStatusRequest")
	}
	return &protos.Void{}, srv.Store.DeleteGatewayStatus(req)
}

// Deletes the network's status table, the table must be emptied prior to
// removal or the operation will fail
// A caller can use magmad to list all available network gateways
// and iterate to remove every gateway status (see: DeleteGatewayStatus above)
func (srv *checkindServer) DeleteNetwork(ctx context.Context, networkId *protos.NetworkID) (*protos.Void, error) {
	if networkId == nil {
		return &protos.Void{}, fmt.Errorf("Nil Network ID")
	}
	return &protos.Void{}, srv.Store.DeleteNetworkTable(networkId.Id)
}

// Returns a list of all logical gateway IDs for the given network which have
// status stored in the service DB
func (srv *checkindServer) List(ctx context.Context, networkId *protos.NetworkID) (*protos.IDList, error) {
	list := new(protos.IDList)
	if networkId == nil {
		return list, fmt.Errorf("Nil Network ID")
	}
	ids, err := srv.Store.List(networkId.Id)

	if err == nil {
		list.Ids = ids
	}
	return list, err
}
