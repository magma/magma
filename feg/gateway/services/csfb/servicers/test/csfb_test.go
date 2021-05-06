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

package test

import (
	"context"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/encode/message"
	"magma/feg/gateway/services/csfb/servicers/mocks"
	"magma/feg/gateway/services/csfb/test_init"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

const mandatoryFieldLength = decode.LengthIEI + decode.LengthLengthIndicator

func TestCsfbServer_AlertAc(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.AlertAck{
		Imsi: "111111",
	}
	encodedMsg, _ := message.EncodeSGsAPAlertAck(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.AlertAc(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_AlertRej(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.AlertReject{
		Imsi:     "111111",
		SgsCause: make([]byte, decode.IELengthSGsCause-mandatoryFieldLength),
	}
	encodedMsg, _ := message.EncodeSGsAPAlertReject(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.AlertRej(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_EPSDetachInd(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.EPSDetachIndication{
		Imsi:                         "111111",
		MmeName:                      "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcde",
		ImsiDetachFromEpsServiceType: []byte{byte(0x11)},
	}
	encodedMsg, _ := message.EncodeSGsAPEPSDetachIndication(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.EPSDetachInd(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_IMSIDetachInd(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.IMSIDetachIndication{
		Imsi:                            "111111",
		MmeName:                         "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcde",
		ImsiDetachFromNonEpsServiceType: []byte{byte(0x11)},
	}
	encodedMsg, _ := message.EncodeSGsAPIMSIDetachIndication(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.IMSIDetachInd(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_LocationUpdateReq(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.LocationUpdateRequest{
		Imsi:                      "111111",
		MmeName:                   "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcde",
		EpsLocationUpdateType:     make([]byte, decode.IELengthEPSLocationUpdateType-mandatoryFieldLength),
		NewLocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	}
	encodedMsg, _ := message.EncodeSGsAPLocationUpdateRequest(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.LocationUpdateReq(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_PagingRej(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.PagingReject{
		Imsi:     "111111",
		SgsCause: make([]byte, decode.IELengthSGsCause-mandatoryFieldLength),
	}
	encodedMsg, _ := message.EncodeSGsAPPagingReject(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.PagingRej(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_ServiceReq(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.ServiceRequest{
		Imsi:             "111111",
		ServiceIndicator: make([]byte, decode.IELengthServiceIndicator-mandatoryFieldLength),
	}
	encodedMsg, _ := message.EncodeSGsAPServiceRequest(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.ServiceReq(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_TMSIReallocationComp(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.TMSIReallocationComplete{
		Imsi: "111111",
	}
	encodedMsg, _ := message.EncodeSGsAPTMSIReallocationComplete(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.TMSIReallocationComp(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_UEActivityInd(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.UEActivityIndication{
		Imsi: "111111",
	}
	encodedMsg, _ := message.EncodeSGsAPUEActivityIndication(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.UEActivityInd(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_UEUnreach(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.UEUnreachable{
		Imsi:     "111111",
		SgsCause: make([]byte, decode.IELengthSGsCause-mandatoryFieldLength),
	}
	encodedMsg, _ := message.EncodeSGsAPUEUnreachable(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.UEUnreach(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_Uplink(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.UplinkUnitdata{
		Imsi:                "111111",
		NasMessageContainer: make([]byte, decode.IELengthNASMessageContainerMax-mandatoryFieldLength),
	}
	encodedMsg, _ := message.EncodeSGsAPUplinkUnitdata(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.Uplink(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_MMEResetAck(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.ResetAck{
		MmeName: "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcde",
	}
	encodedMsg, _ := message.EncodeSGsAPResetAck(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.MMEResetAck(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_MMEResetIndication(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.ResetIndication{
		MmeName: "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcde",
	}
	encodedMsg, _ := message.EncodeSGsAPResetIndication(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.MMEResetIndication(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}

func TestCsfbServer_MMEStatus(t *testing.T) {
	mockInterface := &mocks.ClientConnectionInterface{}
	req := &protos.Status{
		Imsi:             "111111",
		SgsCause:         make([]byte, decode.IELengthSGsCause-mandatoryFieldLength),
		ErroneousMessage: make([]byte, decode.IELengthErroneousMessageMin-mandatoryFieldLength),
	}
	encodedMsg, _ := message.EncodeSGsAPStatus(req)
	mockInterface.On("Send", encodedMsg).Return(nil)

	conn := test_init.GetConnToTestFedGWServiceServer(t, mockInterface)
	defer conn.Close()

	client := protos.NewCSFBFedGWServiceClient(conn)
	reply, err := client.MMEStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, reply)

	mockInterface.AssertNumberOfCalls(t, "Send", 1)
	mockInterface.AssertExpectations(t)
}
