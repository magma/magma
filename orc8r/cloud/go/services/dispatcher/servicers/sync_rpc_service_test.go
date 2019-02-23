/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"context"
	"io"
	"testing"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"
	"magma/orc8r/cloud/go/services/dispatcher"
	dispatcher_test_init "magma/orc8r/cloud/go/services/dispatcher/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

const TestSyncRPCAgHwId = "Test-AGW-Hw-Id"

func TestSyncRPC(t *testing.T) {
	magmad_test_init.StartTestService(t)
	directoryd_test_init.StartTestService(t)
	mockBroker := dispatcher_test_init.StartTestService(t)
	gwReq := &protos.GatewayRequest{
		GwId:      TestSyncRPCAgHwId,
		Authority: "test_authority",
		Path:      "test path",
		Headers:   map[string]string{"te": "trailers", "content-type": "grpc"},
		Payload:   []byte("test payload"),
	}
	syncRPCReq := &protos.SyncRPCRequest{ReqId: 1, ReqBody: gwReq}
	mockBroker.On("CleanupGateway", TestSyncRPCAgHwId).Return(nil)
	queue := make(chan *protos.SyncRPCRequest, 10)
	queue <- syncRPCReq
	mockBroker.On("InitializeGateway", TestSyncRPCAgHwId).Return(queue)
	synResp1 := &protos.SyncRPCResponse{ReqId: 2}
	synResp2 := &protos.SyncRPCResponse{ReqId: 1, RespBody: &protos.GatewayResponse{Status: "200"}, HeartBeat: false}
	mockBroker.On("ProcessGatewayResponse", proto.Clone(synResp1).(*protos.SyncRPCResponse)).Return(nil)
	mockBroker.On("ProcessGatewayResponse", proto.Clone(synResp2).(*protos.SyncRPCResponse)).Return(nil)
	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"},
		"sync_rpc_test_network")
	assert.NoError(t, err)

	t.Logf("New Registered Network: %s", testNetworkId)
	hwId := protos.AccessGatewayID{Id: TestSyncRPCAgHwId}
	logicalId, err := magmad.RegisterGateway(testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &hwId, Name: "Test GW Name"})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalId, "")

	conn, err := registry.GetConnection(dispatcher.ServiceName)
	assert.NoError(t, err)
	syncRPCClient := protos.NewSyncRPCServiceClient(conn)
	stream, err := syncRPCClient.EstablishSyncRPCStream(context.Background())
	assert.NoError(t, err)
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			assert.NoError(t, err)
			if protos.TestMarshal(in) != protos.TestMarshal(syncRPCReq) {
				t.Fatalf("request received at gateway is different from request sent on the service: "+
					"received: %v, sent: %v\n", in, syncRPCReq)
			}

		}
	}()
	// ProcessGatewayResponse on broker should not be called as HeartBeat is true
	err = stream.Send(&protos.SyncRPCResponse{ReqId: 3, RespBody: &protos.GatewayResponse{Status: "200"},
		HeartBeat: true})
	assert.NoError(t, err)

	// ProcessGatewayResponse on broker should be called even when RespBody is nil
	err = stream.Send(synResp1)
	assert.NoError(t, err)

	// ProcessGatewayResponse on broker should be called
	err = stream.Send(synResp2)
	assert.NoError(t, err)
	stream.CloseSend()
	<-waitc
	// wait until server receives from the stream
	time.Sleep(time.Second * 3)
	mockBroker.AssertCalled(t, "InitializeGateway", TestSyncRPCAgHwId)
	// should only be called once
	mockBroker.AssertNumberOfCalls(t, "ProcessGatewayResponse", 2)
	mockBroker.AssertExpectations(t)
}
