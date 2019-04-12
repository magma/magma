/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package broker

import (
	"errors"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/dispatcher/broker/memstore"
)

const QUEUE_LEN = 50
const PROCESS_RESPONSE_TIMEOUT = time.Second * 3

type GatewayRPCBrokerImpl struct {
	responseTable memstore.ResponseTable
	requests      memstore.RequestQueue
}

func NewGatewayReqRespBroker() *GatewayRPCBrokerImpl {
	respTable := memstore.NewResponseTable(PROCESS_RESPONSE_TIMEOUT)
	requests := memstore.NewRequestQueue(QUEUE_LEN)
	return &GatewayRPCBrokerImpl{
		responseTable: respTable,
		requests:      requests,
	}
}

// =========
//  APIs
// =========

// called by httpServer
// caller should time out on respChan.
func (broker *GatewayRPCBrokerImpl) SendRequestToGateway(
	gwReq *protos.GatewayRequest,
) (*GatewayResponseChannel, error) {
	if gwReq == nil || len(gwReq.GwId) == 0 {
		return nil, errors.New("gwReq cannot be nil and gwId cannot be empty string")
	}
	respChan, reqId := broker.responseTable.InitializeResponse()
	// add request to queue
	syncRPCReq := &protos.SyncRPCRequest{ReqId: reqId, ReqBody: gwReq}
	if err := broker.requests.Enqueue(syncRPCReq); err != nil {
		return nil, err
	}
	return &GatewayResponseChannel{respChan, reqId}, nil
}

// called by grpc servicer
func (broker *GatewayRPCBrokerImpl) ProcessGatewayResponse(response *protos.SyncRPCResponse) error {
	return broker.responseTable.SendResponse(response)
}

// called by grpc servicer
func (broker *GatewayRPCBrokerImpl) InitializeGateway(gwId string) chan *protos.SyncRPCRequest {
	// also returns old queue that requests in which can be cancelled. As we don't do anything now,
	// the requests will just time out.
	initializedQueue := broker.requests.InitializeQueue(gwId)
	return initializedQueue.NewQueue
}

// called by grpc servicer
func (broker *GatewayRPCBrokerImpl) CleanupGateway(gwId string) error {
	broker.requests.CleanupQueue(gwId)
	// todo: the above returns the old queue to cleanup. if receive from the queue, it would race with the receiver
	// that listens on the queue to send down to the stream, for now forget
	return nil
}

func (broker *GatewayRPCBrokerImpl) CancelGatewayRequest(gwId string, reqId uint32) error {
	syncRPCRequest := &protos.SyncRPCRequest{ReqId: reqId, ReqBody: &protos.GatewayRequest{GwId: gwId}, ConnClosed: true}
	if err := broker.requests.Enqueue(syncRPCRequest); err != nil {
		return err
	}
	return nil
}
