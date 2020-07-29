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

package memstore_test

import (
	"strconv"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/dispatcher/broker/memstore"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

// process response when no reqId is set up, should return err, not panic
func TestResponseTableImpl_ProcessGatewayResponseWithInvalidReqId(t *testing.T) {
	table := memstore.NewResponseTable(time.Second * 3)
	gwResp := &protos.GatewayResponse{Status: "200", Payload: []byte("test payload")}
	assert.NotPanics(t, func() {
		err := table.SendResponse(&protos.SyncRPCResponse{ReqId: 128, RespBody: gwResp})
		assert.EqualError(t, err, "No response channel found for reqId 128\n")
	})
}

func TestResponseTableImpl_SimpleFlow(t *testing.T) {
	expectedGwResp := &protos.GatewayResponse{
		Status:  "400",
		Payload: []byte("test bytes"),
		Headers: map[string]string{"te": "tailers"}}
	table := memstore.NewResponseTable(time.Second * 3)
	mockClient(t, table, expectedGwResp)
}

func TestResponseTableImpl_ConcurrentFlow(t *testing.T) {
	table := memstore.NewResponseTable(time.Second * 3)
	for i := 0; i < 10; i++ {
		expectedGwResp := &protos.GatewayResponse{
			Status:  "200",
			Payload: []byte("test byte number: " + strconv.Itoa(i)),
			Headers: map[string]string{"test header key": "test header val"},
		}
		go mockClient(t, table, expectedGwResp)
	}
}

func TestResponseTableImpl_SendResponseTimeout(t *testing.T) {
	table := memstore.NewResponseTable(time.Second * 3)
	_, reqId := table.InitializeResponse()
	gwResp := &protos.GatewayResponse{Status: "200", Payload: []byte("test payload")}
	syncResp := &protos.SyncRPCResponse{ReqId: reqId, RespBody: gwResp}
	err := table.SendResponse(syncResp)
	assert.EqualError(t, err, "sendResponse timed out as respChan is not being actively waited on")
}

// initialize response with respChan, waits for a response and expect it to equal expectedGwResp
// create a response channel and waits on it for response. Meanwhile in another goroutine, initializeResponse on the
// table, and send response back. Assert the response is sent back to the correct channel.
func mockClient(
	t *testing.T,
	table memstore.ResponseTable,
	expectedGwResp *protos.GatewayResponse) {
	respChanChan := make(chan chan *protos.GatewayResponse)
	go func(chan chan *protos.GatewayResponse) {
		// initialize response
		respChan, reqId := table.InitializeResponse()
		respChanChan <- respChan
		close(respChanChan)
		// send response
		table.SendResponse(&protos.SyncRPCResponse{ReqId: reqId, RespBody: expectedGwResp})
	}(respChanChan)
	respChan := <-respChanChan
	// wait for response
	gwResp := <-respChan
	assertGatewayRespEqual(t, gwResp, expectedGwResp)
}

func assertGatewayRespEqual(t *testing.T, resp1 *protos.GatewayResponse, resp2 *protos.GatewayResponse) {
	if resp1 == nil && resp2 == nil {
		return
	}
	if resp1 == nil || resp2 == nil {
		assert.Fail(t, "not equal: resp1: %v, resp2: %v\n", resp1, resp2)
	}
	assert.Equal(t, resp1.Status, resp2.Status)
	assert.Equal(t, resp1.Payload, resp2.Payload)
	assert.Equal(t, resp1.Headers, resp2.Headers)
}
