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

package memstore

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

type ResponseTable interface {
	InitializeResponse() (chan *protos.GatewayResponse, uint32)
	SendResponse(*protos.SyncRPCResponse) error
}

type ResponseTableImpl struct {
	// respChanByReqId is a map: <uint32, chan *protos.GatewayResponse>
	respChanByReqId *sync.Map
	// reqIdCounter strictly increases, making sure all reqIds will be unique.
	reqIdCounter uint32
	timeout      time.Duration
}

func NewResponseTable(timeout time.Duration) *ResponseTableImpl {
	return &ResponseTableImpl{respChanByReqId: &sync.Map{}, timeout: timeout}
}

// InitializeResponse creates a request ID to bind to the GatewayResponse
// channel, so when SyncRPCResponse comes back, it can be written to the
// corresponding GatewayResponse channel identified by the request id.
func (table *ResponseTableImpl) InitializeResponse() (chan *protos.GatewayResponse, uint32) {
	reqId := generateReqId(&table.reqIdCounter)
	respChan := make(chan *protos.GatewayResponse)
	table.respChanByReqId.Store(reqId, respChan)
	return respChan, reqId
}

// SendResponse sends the response to the corresponding response channel
func (table *ResponseTableImpl) SendResponse(resp *protos.SyncRPCResponse) error {
	if resp == nil {
		return errors.New("cannot send nil SyncRPCResponse")
	}

	respChanVal, ok := table.respChanByReqId.Load(resp.ReqId)
	if !ok {
		return fmt.Errorf("No response channel found for reqId %v\n", resp.ReqId)
	}
	respChan := respChanVal.(chan *protos.GatewayResponse)
	if resp.RespBody == nil {
		glog.Errorf("Nil response body received, forward to httpServer anyways\n")
	}
	select {
	case respChan <- resp.RespBody:
		return nil
	case <-time.After(table.timeout):
		// give up sending, close the channel and delete the table entry
		close(respChan)
		table.respChanByReqId.Delete(resp.ReqId)
		return errors.New("sendResponse timed out as respChan is not being actively waited on")
	}
}

func generateReqId(counter *uint32) uint32 {
	return atomic.AddUint32(counter, 1)
}
