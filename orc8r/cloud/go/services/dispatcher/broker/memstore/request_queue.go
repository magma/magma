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
	"time"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

// InitializedQueue contains an initialized NewQueue, and an OldQueue to cleanup if any
type InitializedQueue struct {
	NewQueue chan *protos.SyncRPCRequest
	OldQueue chan *protos.SyncRPCRequest
}

type RequestQueue interface {
	InitializeQueue(gwId string) InitializedQueue
	CleanupQueue(gwId string) chan *protos.SyncRPCRequest
	Enqueue(req *protos.SyncRPCRequest) error
}

type requestQueueImpl struct {
	reqQueueByGwId map[string]chan *protos.SyncRPCRequest
	*sync.RWMutex
	maxQueueLen int
}

func NewRequestQueue(queueLen int) *requestQueueImpl {
	return &requestQueueImpl{
		make(map[string]chan *protos.SyncRPCRequest),
		&sync.RWMutex{},
		queueLen,
	}
}

// InitializeQueue returns a new queue to dequeue from, and the old queue to clean up
func (queues *requestQueueImpl) InitializeQueue(gwId string) InitializedQueue {
	queues.Lock()
	defer queues.Unlock()
	newQueue := make(chan *protos.SyncRPCRequest, queues.maxQueueLen)
	if oldQueue, ok := queues.reqQueueByGwId[gwId]; ok {
		// sends on a closed queue will panic, but no one will send onto this queue,
		// because all sends are through enqueue.
		close(oldQueue)
		queues.reqQueueByGwId[gwId] = newQueue
		return InitializedQueue{newQueue, oldQueue}
	}
	queues.reqQueueByGwId[gwId] = newQueue
	return InitializedQueue{newQueue, nil}
}

// CleanupQueue return the old queue to clean up
func (queues *requestQueueImpl) CleanupQueue(gwId string) chan *protos.SyncRPCRequest {
	queues.Lock()
	defer queues.Unlock()
	if queue, ok := queues.reqQueueByGwId[gwId]; ok {
		// sends on a closed queue will panic, but no one will send onto this queue,
		// because all sends are through enqueue.
		close(queue)
		delete(queues.reqQueueByGwId, gwId)
		// the broker will cleanup requests in the queue
		return queue
	} else {
		glog.Warningf("HWID %v: no request queue found to clean up", gwId)
	}
	return nil
}

// Enqueue adds a SyncRPCRequest to the queue of gatewayId gwId.
// gwId cannot be empty string, gwReq or ReqId of gwReq cannot be nil.
// gwId: key of the syncRPCReqQueue map
// gwReq: to append to []protos.SyncRPCRequest of the syncRPCReqQueue
func (queues *requestQueueImpl) Enqueue(gwReq *protos.SyncRPCRequest) error {
	const (
		totalQueuingTimeout = time.Second
		minQueuingTimeout   = time.Millisecond * 700
		waitIncrement       = time.Millisecond * 10
	)
	if gwReq == nil || gwReq.ReqId <= 0 || gwReq.ReqBody == nil || len(gwReq.ReqBody.GwId) == 0 {
		return errors.New("SyncRPCRequest cannot be nil and gwId of ReqBody has to be valid")
	}
	queueingTimeout := totalQueuingTimeout
	gwId := gwReq.ReqBody.GwId
	queues.RLock()
	reqQueue, ok := queues.reqQueueByGwId[gwId]
	for !ok {
		queues.RUnlock()
		if queueingTimeout < minQueuingTimeout {
			return fmt.Errorf("Queue does not exist for gwId %v\n", gwId)
		}
		time.Sleep(waitIncrement)
		queueingTimeout -= waitIncrement
		queues.RLock()
		reqQueue, ok = queues.reqQueueByGwId[gwId]
	}
	defer queues.RUnlock()
	select {
	case reqQueue <- gwReq:
		return nil
	case <-time.After(queueingTimeout):
		return fmt.Errorf("Failed to enqueue %v because queue for gwId %v is full\n", gwReq, gwId)
	}
}
