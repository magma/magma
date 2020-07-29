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

package csfb

import (
	"context"
	"errors"
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/csfb/servicers/decode"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

type grpcClient func(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error)

var clientFuncMap = map[decode.SGsMessageType]grpcClient{
	decode.SGsAPAlertRequest:         alertRequestClient,
	decode.SGsAPDownlinkUnitdata:     downlinkUnitdataClient,
	decode.SGsAPEPSDetachAck:         epsDetachAckClient,
	decode.SGsAPIMSIDetachAck:        imsiDetachAckClient,
	decode.SGsAPLocationUpdateAccept: locationUpdateAcceptClient,
	decode.SGsAPLocationUpdateReject: locationUpdateRejectClient,
	decode.SGsAPMMInformationRequest: mmInformationRequestClient,
	decode.SGsAPPagingRequest:        pagingRequestClient,
	decode.SGsAPReleaseRequest:       releaseRequestClient,
	decode.SGsAPServiceAbortRequest:  serviceAbortClient,
	decode.SGsAPResetAck:             resetAckClient,
	decode.SGsAPResetIndication:      resetIndicationClient,
	decode.SGsAPStatus:               statusClient,
}

func SendSGsMessageToGateway(messageType decode.SGsMessageType, msg *any.Any) (*orcprotos.Void, error) {
	conn, err := registry.Get().GetCloudConnection(feg_relay.ServiceName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to establish connection to cloud FegToGwRelayClient: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	defer conn.Close()

	client := protos.NewCSFBGatewayServiceClient(conn)

	if clientFunc, ok := clientFuncMap[messageType]; ok {
		return clientFunc(msg, client)
	}
	return &orcprotos.Void{}, errors.New("unknown message type")
}

func alertRequestClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.AlertRequest{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.AlertReq(context.Background(), unmarshalledMsg)
}

func downlinkUnitdataClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.DownlinkUnitdata{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.Downlink(context.Background(), unmarshalledMsg)
}

func epsDetachAckClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.EPSDetachAck{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.EPSDetachAc(context.Background(), unmarshalledMsg)
}

func imsiDetachAckClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.IMSIDetachAck{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.IMSIDetachAc(context.Background(), unmarshalledMsg)
}

func locationUpdateAcceptClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.LocationUpdateAccept{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.LocationUpdateAcc(context.Background(), unmarshalledMsg)
}

func locationUpdateRejectClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.LocationUpdateReject{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.LocationUpdateRej(context.Background(), unmarshalledMsg)
}

func mmInformationRequestClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.MMInformationRequest{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.MMInformationReq(context.Background(), unmarshalledMsg)
}

func pagingRequestClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.PagingRequest{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.PagingReq(context.Background(), unmarshalledMsg)
}

func releaseRequestClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.ReleaseRequest{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.ReleaseReq(context.Background(), unmarshalledMsg)
}

func serviceAbortClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.ServiceAbortRequest{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.ServiceAbort(context.Background(), unmarshalledMsg)
}

func resetAckClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.ResetAck{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.VLRResetAck(context.Background(), unmarshalledMsg)
}

func resetIndicationClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.ResetIndication{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.VLRResetIndication(context.Background(), unmarshalledMsg)
}

func statusClient(msg *any.Any, client protos.CSFBGatewayServiceClient) (*orcprotos.Void, error) {
	unmarshalledMsg := &protos.Status{}
	ptypes.UnmarshalAny(msg, unmarshalledMsg)
	return client.VLRStatus(context.Background(), unmarshalledMsg)
}
