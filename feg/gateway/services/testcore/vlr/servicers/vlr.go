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
	"fmt"
	"reflect"
	"sync"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/encode/message"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

type Request struct {
	encodedRequest    []byte
	requestType       decode.SGsMessageType
	marshalledRequest *any.Any
}

type Reply struct {
	serverBehavior protos.Reply_ServerBehavior
	delayingTime   int
	encodedReply   []byte
	replyType      decode.SGsMessageType
}

type RequestReply struct {
	request Request
	reply   Reply
}

type VLRServer struct {
	Conn              servicers.ServerConnectionInterface
	queueMux          *sync.Mutex
	requestReplyQueue []*RequestReply
	queueIndex        int
}

func NewVLRServer(ConnectionInterface servicers.ServerConnectionInterface) (*VLRServer, error) {
	return &VLRServer{
		Conn:              ConnectionInterface,
		queueMux:          new(sync.Mutex),
		requestReplyQueue: []*RequestReply{},
		queueIndex:        0}, nil
}

func (srv *VLRServer) GetNextRequestReplyPair() (*RequestReply, error) {
	srv.queueMux.Lock()
	defer srv.queueMux.Unlock()
	if srv.queueIndex >= len(srv.requestReplyQueue) {
		return nil, fmt.Errorf("no mock response to return")
	}
	requestReplyPair := srv.requestReplyQueue[srv.queueIndex]
	srv.queueIndex++
	return requestReplyPair, nil
}

func (srv *VLRServer) ReplyClient(clientRequest []byte) error {
	requestReply, err := srv.GetNextRequestReplyPair()
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(requestReply.request.encodedRequest, clientRequest) {
		return fmt.Errorf(
			"received unexpected request, expected % x but got % x",
			requestReply.request.encodedRequest,
			clientRequest,
		)
	}

	glog.V(2).Infof(
		"Message for replying: %s",
		decode.MsgTypeNameByCode[requestReply.reply.replyType],
	)
	if requestReply.reply.serverBehavior == protos.Reply_REPLY_INSTANTLY {
		glog.V(2).Info("Reply to the client instantly")
		err = srv.Conn.SendFromServer(requestReply.reply.encodedReply)
	} else if requestReply.reply.serverBehavior == protos.Reply_REPLY_LATE {
		glog.V(2).Infof(
			"Reply to the client after waiting for %d seconds",
			requestReply.reply.delayingTime,
		)
		time.Sleep(
			time.Duration(requestReply.reply.delayingTime) * time.Second,
		)
		err = srv.Conn.SendFromServer(requestReply.reply.encodedReply)
	} else {
		glog.V(2).Info("Do not reply to the client")
	}
	return err
}

func (srv *VLRServer) Reset(
	ctx context.Context,
	req *orcprotos.Void,
) (*orcprotos.Void, error) {
	srv.queueMux.Lock()
	defer srv.queueMux.Unlock()

	if srv.queueIndex < len(srv.requestReplyQueue) {
		glog.Warning("Received reset request before all expected VLR requests were received")
	}
	srv.queueIndex = 0
	srv.requestReplyQueue = []*RequestReply{}
	return &orcprotos.Void{}, nil
}

func (srv *VLRServer) ConfigServer(
	ctx context.Context,
	config *protos.ServerConfiguration,
) (*orcprotos.Void, error) {
	srv.queueMux.Lock()
	defer srv.queueMux.Unlock()

	if srv.queueIndex < len(srv.requestReplyQueue) {
		glog.Warning("Received new server configuration before replies are used up")
	}
	srv.queueIndex = 0
	srv.requestReplyQueue = nil

	for _, requestReply := range config.RequestReply {

		request, err := constructRequest(requestReply.Request)
		if err != nil {
			glog.Errorf("Failed to construct expected requests: %s", err)
			return &orcprotos.Void{}, err
		}

		reply, err := constructReply(requestReply.Reply)
		if err != nil {
			glog.Errorf("Failed to construct replies: %s", err)
			return &orcprotos.Void{}, err
		}

		// append the pair of request and reply to the queue
		encodedRequestReply := RequestReply{
			request: *request,
			reply:   *reply,
		}
		srv.requestReplyQueue = append(
			srv.requestReplyQueue,
			&encodedRequestReply,
		)
	}
	return &orcprotos.Void{}, nil
}

func constructRequest(protoRequest *protos.ExpectedRequest) (*Request, error) {
	var encodedMsg []byte
	var err error
	var requestType decode.SGsMessageType
	var protoMsg proto.Message

	switch t := protoRequest.SgsMessage.(type) {
	case *protos.ExpectedRequest_AlertAck:
		requestType = decode.SGsAPAlertAck
		protoMsg = protoRequest.GetAlertAck()
		encodedMsg, err = message.EncodeSGsAPAlertAck(
			protoRequest.GetAlertAck(),
		)
	case *protos.ExpectedRequest_AlertReject:
		requestType = decode.SGsAPAlertReject
		protoMsg = protoRequest.GetAlertReject()
		encodedMsg, err = message.EncodeSGsAPAlertReject(
			protoRequest.GetAlertReject(),
		)
	case *protos.ExpectedRequest_EpsDetachIndication:
		requestType = decode.SGsAPEPSDetachIndication
		protoMsg = protoRequest.GetEpsDetachIndication()
		encodedMsg, err = message.EncodeSGsAPEPSDetachIndication(
			protoRequest.GetEpsDetachIndication(),
		)
	case *protos.ExpectedRequest_ImsiDetachIndication:
		requestType = decode.SGsAPIMSIDetachIndication
		protoMsg = protoRequest.GetImsiDetachIndication()
		encodedMsg, err = message.EncodeSGsAPIMSIDetachIndication(
			protoRequest.GetImsiDetachIndication(),
		)
	case *protos.ExpectedRequest_LocationUpdateRequest:
		requestType = decode.SGsAPLocationUpdateRequest
		protoMsg = protoRequest.GetLocationUpdateRequest()
		encodedMsg, err = message.EncodeSGsAPLocationUpdateRequest(
			protoRequest.GetLocationUpdateRequest(),
		)
	case *protos.ExpectedRequest_PagingReject:
		requestType = decode.SGsAPPagingReject
		protoMsg = protoRequest.GetPagingReject()
		encodedMsg, err = message.EncodeSGsAPPagingReject(
			protoRequest.GetPagingReject(),
		)
	case *protos.ExpectedRequest_ServiceRequest:
		requestType = decode.SGsAPServiceRequest
		protoMsg = protoRequest.GetServiceRequest()
		encodedMsg, err = message.EncodeSGsAPServiceRequest(
			protoRequest.GetServiceRequest(),
		)
	case *protos.ExpectedRequest_TmsiReallocationComplete:
		requestType = decode.SGsAPTMSIReallocationComplete
		protoMsg = protoRequest.GetTmsiReallocationComplete()
		encodedMsg, err = message.EncodeSGsAPTMSIReallocationComplete(
			protoRequest.GetTmsiReallocationComplete(),
		)
	case *protos.ExpectedRequest_UeActivityIndication:
		requestType = decode.SGsAPUEActivityIndication
		protoMsg = protoRequest.GetUeActivityIndication()
		encodedMsg, err = message.EncodeSGsAPUEActivityIndication(
			protoRequest.GetUeActivityIndication(),
		)
	case *protos.ExpectedRequest_UeUnreachable:
		requestType = decode.SGsAPUEUnreachable
		protoMsg = protoRequest.GetUeUnreachable()
		encodedMsg, err = message.EncodeSGsAPUEUnreachable(
			protoRequest.GetUeUnreachable(),
		)
	case *protos.ExpectedRequest_UplinkUnitdata:
		requestType = decode.SGsAPUplinkUnitdata
		protoMsg = protoRequest.GetUplinkUnitdata()
		encodedMsg, err = message.EncodeSGsAPUplinkUnitdata(
			protoRequest.GetUplinkUnitdata(),
		)
	case *protos.ExpectedRequest_ResetAck:
		requestType = decode.SGsAPResetAck
		protoMsg = protoRequest.GetResetAck()
		encodedMsg, err = message.EncodeSGsAPResetAck(
			protoRequest.GetResetAck(),
		)
	case *protos.ExpectedRequest_ResetIndication:
		requestType = decode.SGsAPResetIndication
		protoMsg = protoRequest.GetResetIndication()
		encodedMsg, err = message.EncodeSGsAPResetIndication(
			protoRequest.GetResetIndication(),
		)
	case *protos.ExpectedRequest_Status:
		requestType = decode.SGsAPStatus
		protoMsg = protoRequest.GetStatus()
		encodedMsg, err = message.EncodeSGsAPStatus(
			protoRequest.GetStatus(),
		)
	default:
		err = fmt.Errorf("unexpected type %T", t)
	}
	if err != nil {
		glog.Errorf("Fail to encoded request message: %s", err)
		return nil, err
	}
	marshalledMsg, err := ptypes.MarshalAny(protoMsg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling SGs message to Any: %s", err)
	}

	request := Request{
		encodedRequest:    encodedMsg,
		requestType:       requestType,
		marshalledRequest: marshalledMsg,
	}

	return &request, nil
}

func constructReply(protoReply *protos.Reply) (*Reply, error) {
	var replyType decode.SGsMessageType
	var encodedMsg []byte
	var err error

	switch t := protoReply.SgsMessage.(type) {
	case *protos.Reply_AlertRequest:
		replyType = decode.SGsAPAlertRequest
		encodedMsg, err = message.EncodeSGsAPAlertRequest(
			protoReply.GetAlertRequest(),
		)
	case *protos.Reply_DownlinkUnitdata:
		replyType = decode.SGsAPDownlinkUnitdata
		encodedMsg, err = message.EncodeSGsAPDownlinkUnitdata(
			protoReply.GetDownlinkUnitdata(),
		)
	case *protos.Reply_EpsDetachAck:
		replyType = decode.SGsAPEPSDetachAck
		encodedMsg, err = message.EncodeSGsAPEPSDetachAck(
			protoReply.GetEpsDetachAck(),
		)
	case *protos.Reply_ImsiDetachAck:
		replyType = decode.SGsAPIMSIDetachAck
		encodedMsg, err = message.EncodeSGsAPIMSIDetachAck(
			protoReply.GetImsiDetachAck(),
		)
	case *protos.Reply_LocationUpdateAccept:
		replyType = decode.SGsAPLocationUpdateAccept
		encodedMsg, err = message.EncodeSGsAPLocationUpdateAccept(
			protoReply.GetLocationUpdateAccept(),
		)
	case *protos.Reply_LocationUpdateReject:
		replyType = decode.SGsAPLocationUpdateReject
		encodedMsg, err = message.EncodeSGsAPLocationUpdateReject(
			protoReply.GetLocationUpdateReject(),
		)
	case *protos.Reply_MmInformationRequest:
		replyType = decode.SGsAPMMInformationRequest
		encodedMsg, err = message.EncodeSGsAPMMInformationRequest(
			protoReply.GetMmInformationRequest(),
		)
	case *protos.Reply_PagingRequest:
		replyType = decode.SGsAPPagingRequest
		encodedMsg, err = message.EncodeSGsAPPagingRequest(
			protoReply.GetPagingRequest(),
		)
	case *protos.Reply_ReleaseRequest:
		replyType = decode.SGsAPReleaseRequest
		encodedMsg, err = message.EncodeSGsAPReleaseRequest(
			protoReply.GetReleaseRequest(),
		)
	case *protos.Reply_ServiceAbortRequest:
		replyType = decode.SGsAPServiceAbortRequest
		encodedMsg, err = message.EncodeSGsAPServiceAbortRequest(
			protoReply.GetServiceAbortRequest(),
		)
	case *protos.Reply_ResetAck:
		replyType = decode.SGsAPResetAck
		encodedMsg, err = message.EncodeSGsAPResetAck(
			protoReply.GetResetAck(),
		)
	case *protos.Reply_ResetIndication:
		replyType = decode.SGsAPResetIndication
		encodedMsg, err = message.EncodeSGsAPResetIndication(
			protoReply.GetResetIndication(),
		)
	case *protos.Reply_Status:
		replyType = decode.SGsAPStatus
		encodedMsg, err = message.EncodeSGsAPStatus(
			protoReply.GetStatus(),
		)
	default:
		err = fmt.Errorf("unexpected type %T", t)
	}
	if err != nil {
		glog.Errorf("Fail to encoded reply message: %s", err)
		return nil, err
	}

	reply := Reply{
		serverBehavior: protoReply.ServerBehavior,
		delayingTime:   int(protoReply.ReplyDelay),
		encodedReply:   encodedMsg,
		replyType:      replyType,
	}

	return &reply, nil
}
