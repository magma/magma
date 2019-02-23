/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package broker

import "magma/orc8r/cloud/go/protos"

// ========================
// GatewayRPCBroker is the bridge between httpServer and SyncRPC grpc servicer,
// where httpServer handles requests from the cloud service instances, and grpc servicer talks directly
// to the gateways using grpc bidirectional stream.
// ==========================
type GatewayRPCBroker interface {
	// httpServer sends request to a certain gateway, and waits on the response channel for response
	SendRequestToGateway(gwReq *protos.GatewayRequest) (chan *protos.GatewayResponse, error)
	// receive a SyncRPCResponse from grpc servicer, and send the corresponding GatewayResponse to httpServer
	ProcessGatewayResponse(response *protos.SyncRPCResponse) error
	// Initialize the necessary datastructure for a gwId when the gw connects to SyncRPC grpc servicer so the dispatcher
	// is ready to take any requests for this gateway, and returns
	// a request queue for grpc servicer to listen on for incoming requests from httpserver.
	InitializeGateway(gwId string) chan *protos.SyncRPCRequest
	// Cleanup the data and resources for a gwId when the gw loses SyncRPC connection to the cloud
	CleanupGateway(gwId string) error
}
