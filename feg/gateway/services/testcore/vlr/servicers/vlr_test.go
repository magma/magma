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

package servicers

import (
	"context"
	"errors"
	"net"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/encode/message"
	"magma/feg/gateway/services/csfb/servicers/mocks"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestVLRServer_GetNextRequestReplyPair(t *testing.T) {
	// initialize the server and the queue
	conn, _ := servicers.NewSCTPServerConnection()
	srv, _ := NewVLRServer(conn)
	testPair := RequestReply{
		request: Request{
			encodedRequest:    []byte{0x04, 0x03, 0x02, 0x01},
			requestType:       decode.SGsAPPagingReject,
			marshalledRequest: nil,
		},
		reply: Reply{
			serverBehavior: protos.Reply_REPLY_INSTANTLY,
			delayingTime:   0,
			encodedReply:   []byte{0x01, 0x02, 0x03, 0x04},
		},
	}
	srv.requestReplyQueue = append(srv.requestReplyQueue, &testPair)
	srv.queueIndex = 0

	// successfully get the next pair
	nextPair, err := srv.GetNextRequestReplyPair()
	assert.NoError(t, err)
	assert.Equal(t, srv.queueIndex, 1)
	assert.Equal(t, *nextPair, testPair)

	// fail to get the next pair since the queue is empty
	_, err = srv.GetNextRequestReplyPair()
	assert.Error(t, err, errors.New("reply queue is used up"))
}

func TestVLRServer_Reset(t *testing.T) {
	connSCTP, _ := servicers.NewSCTPServerConnection()
	conn, srv := getConnAndTestVLRGRPCServer(t, connSCTP)
	defer conn.Close()

	srv.requestReplyQueue = append(srv.requestReplyQueue, &RequestReply{})
	srv.queueIndex = 1

	client := protos.NewMockCoreConfiguratorClient(conn)
	reply, err := client.Reset(context.Background(), &orcprotos.Void{})
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	assert.Equal(t, []*RequestReply{}, srv.requestReplyQueue)
	assert.Equal(t, 0, srv.queueIndex)
}

func TestVLRServer_ReplyClient(t *testing.T) {
	mockInterface := &mocks.ServerConnectionInterface{}
	srv, err := NewVLRServer(mockInterface)
	assert.NoError(t, err)

	mockInterface.On(
		"SendFromServer",
		[]byte{0x01, 0x02, 0x03, 0x04}).Return(nil)

	// populate the reply queue
	testPair := RequestReply{
		request: Request{
			encodedRequest:    []byte{0x04, 0x03, 0x02, 0x01},
			requestType:       decode.SGsAPPagingReject,
			marshalledRequest: nil,
		},
		reply: Reply{
			serverBehavior: protos.Reply_REPLY_INSTANTLY,
			delayingTime:   0,
			encodedReply:   []byte{0x01, 0x02, 0x03, 0x04},
		},
	}
	srv.requestReplyQueue = append(srv.requestReplyQueue, &testPair)
	srv.queueIndex = 0

	// successful case
	err = srv.ReplyClient([]byte{0x04, 0x03, 0x02, 0x01})
	assert.NoError(t, err)
	mockInterface.AssertNumberOfCalls(
		t,
		"SendFromServer",
		1,
	)
	mockInterface.AssertExpectations(t)

	// failing case of unmatched request
	srv.queueIndex = 0
	err = srv.ReplyClient([]byte{0x01, 0x01, 0x01, 0x01})
	assert.EqualError(
		t, err, "received unexpected request, expected 04 03 02 01 but got 01 01 01 01")
}

func TestVLRServer_ConfigServer(t *testing.T) {
	connSCTP, _ := servicers.NewSCTPServerConnection()
	conn, srv := getConnAndTestVLRGRPCServer(t, connSCTP)
	defer conn.Close()

	// request
	epsDetachIndication := protos.EPSDetachIndication{
		Imsi:                         "001010000000001",
		MmeName:                      ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
		ImsiDetachFromEpsServiceType: []byte{byte(0x11)},
	}
	expectedRequest := protos.ExpectedRequest{
		SgsMessage: &protos.ExpectedRequest_EpsDetachIndication{
			EpsDetachIndication: &epsDetachIndication,
		},
	}

	// reply
	epsDetachAck := protos.EPSDetachAck{
		Imsi: "001010000000001",
	}
	reply := protos.Reply{
		ServerBehavior: protos.Reply_REPLY_INSTANTLY,
		ReplyDelay:     0,
		SgsMessage: &protos.Reply_EpsDetachAck{
			EpsDetachAck: &epsDetachAck,
		},
	}

	requestReply := protos.RequestReply{
		Request: &expectedRequest,
		Reply:   &reply,
	}
	config := protos.ServerConfiguration{
		RequestReply: []*protos.RequestReply{&requestReply},
	}

	client := protos.NewMockCoreConfiguratorClient(conn)
	srvReply, err := client.ConfigServer(context.Background(), &config)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, srvReply)

	// encoded request
	encodedEPSDetachIndication, err := message.EncodeSGsAPEPSDetachIndication(
		&epsDetachIndication,
	)
	assert.NoError(t, err)
	marshalledEPSDetachIndication, err := ptypes.MarshalAny(
		&epsDetachIndication,
	)
	assert.NoError(t, err)
	encodedRequest := Request{
		encodedRequest:    encodedEPSDetachIndication,
		requestType:       decode.SGsAPEPSDetachIndication,
		marshalledRequest: marshalledEPSDetachIndication,
	}

	// encoded reply
	encodedEPSDetachAck, err := message.EncodeSGsAPEPSDetachAck(
		&epsDetachAck,
	)
	assert.NoError(t, err)
	encodedReply := Reply{
		serverBehavior: protos.Reply_REPLY_INSTANTLY,
		delayingTime:   0,
		encodedReply:   encodedEPSDetachAck,
		replyType:      decode.SGsAPEPSDetachAck,
	}

	encodedRequestReply := RequestReply{
		request: encodedRequest,
		reply:   encodedReply,
	}

	assert.Equal(t, &encodedRequestReply, srv.requestReplyQueue[0])
}

func getConnAndTestVLRGRPCServer(
	t *testing.T,
	connectionInterface servicers.ServerConnectionInterface,
) (*grpc.ClientConn, *VLRServer) {
	srv, err := NewVLRServer(connectionInterface)
	assert.NoError(t, err)

	s := grpc.NewServer()
	protos.RegisterMockCoreConfiguratorServer(s, srv)

	lis, err := net.Listen("tcp", "")
	assert.NoError(t, err)

	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	addr := lis.Addr()
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	assert.NoError(t, err)
	return conn, srv
}
