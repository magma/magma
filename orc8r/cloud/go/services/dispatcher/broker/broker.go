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

package broker

import (
	"errors"
	"time"

	"magma/orc8r/cloud/go/services/dispatcher/broker/memstore"
	"magma/orc8r/lib/go/protos"
)

const (
	processResponseTimeout = time.Second * 3
	queueLen               = 50
)

type GatewayResponseChannel struct {
	RespChan chan *protos.GatewayResponse
	ReqId    uint32
}

// GatewayRPCBroker is the bridge between httpServer and SyncRPC servicer,
// where httpServer handles requests from the cloud service instances, and SyncRPC servicer talks directly
// to the gateways using a gRPC bidirectional stream.
type GatewayRPCBroker interface {
	// SendRequestToGateway is called by the HTTP server to send a request
	// to a certain gateway, and waits on the response channel for response.
	// The caller should time out on the response channel.
	SendRequestToGateway(gwReq *protos.GatewayRequest) (*GatewayResponseChannel, error)
	// ProcessGatewayResponse is called by the SyncRPC servicer. It receives
	// a SyncRPCResponse from the SyncRPC servicer, and send the corresponding GatewayResponse to the HTTP server
	ProcessGatewayResponse(response *protos.SyncRPCResponse) error
	// InitializeGateway initializes the necessary data structures for a gwId
	// when the gateway connects to the SyncRPC servicer so the dispatcher
	// is ready to take any requests for this gateway, and returns
	// a request queue for gRPC servicer to listen on for incoming requests
	// from HTTP servers.
	InitializeGateway(gwId string) chan *protos.SyncRPCRequest
	// CleanupGateway cleans up the data and resources for a gwId when the gw loses SyncRPC connection to the cloud.
	CleanupGateway(gwId string) error
	// CancelGatewayRequest notifies the gateway to stop handling the request with ID reqId.
	CancelGatewayRequest(gwId string, reqId uint32) error
}

// GatewayRPCBrokerImpl implements a GatewayRPCBroker, managing a response table and request queue.
type GatewayRPCBrokerImpl struct {
	responseTable memstore.ResponseTable
	requests      memstore.RequestQueue
}

func NewGatewayReqRespBroker() *GatewayRPCBrokerImpl {
	respTable := memstore.NewResponseTable(processResponseTimeout)
	requests := memstore.NewRequestQueue(queueLen)
	return &GatewayRPCBrokerImpl{responseTable: respTable, requests: requests}
}

func (broker *GatewayRPCBrokerImpl) SendRequestToGateway(
	gwReq *protos.GatewayRequest,
) (*GatewayResponseChannel, error) {
	if gwReq == nil || len(gwReq.GwId) == 0 {
		return nil, errors.New("gwReq cannot be nil and gwId cannot be empty string")
	}
	respChan, reqId := broker.responseTable.InitializeResponse()
	// Add request to queue.
	syncRPCReq := &protos.SyncRPCRequest{ReqId: reqId, ReqBody: gwReq}
	if err := broker.requests.Enqueue(syncRPCReq); err != nil {
		return nil, err
	}
	return &GatewayResponseChannel{RespChan: respChan, ReqId: reqId}, nil
}

func (broker *GatewayRPCBrokerImpl) ProcessGatewayResponse(response *protos.SyncRPCResponse) error {
	return broker.responseTable.SendResponse(response)
}

func (broker *GatewayRPCBrokerImpl) InitializeGateway(gwId string) chan *protos.SyncRPCRequest {
	// Also returns the old queue that requests in which can be canceled.
	// As we don't do anything now, the requests will just time out.
	initializedQueue := broker.requests.InitializeQueue(gwId)
	return initializedQueue.NewQueue
}

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
