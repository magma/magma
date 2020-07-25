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
	"sync"
	"testing"

	"magma/orc8r/cloud/go/services/dispatcher/broker/memstore"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestSyncRPCReqQueuesMap_EnqueueRequest(t *testing.T) {
	queuesMap := memstore.NewRequestQueue(10)
	// initialize queues
	queuesMap.InitializeQueue("10")
	queuesMap.InitializeQueue("20")
	queuesMap.InitializeQueue("30")
	// concurrent enqueue
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			enqueueGwWithReqId(t, queuesMap, "10", uint32(10+i))
			enqueueGwWithReqId(t, queuesMap, "20", uint32(20+i))
			enqueueGwWithReqId(t, queuesMap, "30", uint32(30+i))
		}(i)
	}

	// wait for enqueueing to finish
	wg.Wait()
}

func TestRequestQueueImpl_EnqueueFullQueue(t *testing.T) {
	queue := memstore.NewRequestQueue(1)
	req := protos.SyncRPCRequest{ReqId: 3, ReqBody: &protos.GatewayRequest{GwId: "gwId1"}}
	queue.InitializeQueue("gwId1")
	queue.InitializeQueue("gwId2")
	err := queue.Enqueue(&req)
	assert.NoError(t, err)
	err = queue.Enqueue(&req)
	assert.EqualError(t, err, "Failed to enqueue reqId:3 reqBody:<gwId:\"gwId1\" >  because queue for gwId gwId1 is full\n")
	// should still be able to enqueue for other gateways
	req = protos.SyncRPCRequest{ReqId: 4, ReqBody: &protos.GatewayRequest{GwId: "gwId2"}}
	err = queue.Enqueue(&req)
	assert.NoError(t, err)
}

func TestEnqueueNonExistingQueue(t *testing.T) {
	queue := memstore.NewRequestQueue(1)
	req := protos.SyncRPCRequest{ReqId: 3, ReqBody: &protos.GatewayRequest{GwId: "gwId1"}}
	err := queue.Enqueue(&req)
	assert.EqualError(t, err, "Queue does not exist for gwId gwId1\n")
}

func enqueueGwWithReqId(t *testing.T, queue memstore.RequestQueue, gwId string, reqId uint32) {
	req := &protos.SyncRPCRequest{ReqId: reqId, ReqBody: &protos.GatewayRequest{GwId: gwId}}
	err := queue.Enqueue(req)
	assert.NoError(t, err)
}
