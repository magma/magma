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

package analytics

const (
	// MaxExecBatchSize the maximum batch execution size for which there's a bucket in stats
	MaxExecBatchSize = 10
)

type (
	// Queue A queue for Analytics Requests to run on
	Queue struct {
		items      chan Request
		terminate  chan bool // signal termination to worker goroutine, with 'drain' param
		done       chan bool // signal back to main context that drain completed
		terminated bool
	}

	// Request the interface to be implemented by a task that is pushed to the
	// analytics.Queue, so it can be executed
	Request interface {
		Run(m ModuleCtx)
	}
)

// NewAnalyticsQueue creates a new queue
func NewAnalyticsQueue(mCtx ModuleCtx) *Queue {
	result := &Queue{
		items:      make(chan Request, 9999),
		terminate:  make(chan bool, 1),
		done:       make(chan bool, 1),
		terminated: false,
	}

	// Start a go routine to truncate the queue
	go func(q *Queue) {
		for {
			select {
			case request := <-q.items:
				request.Run(mCtx)
				break
			case drain := <-q.terminate:
				go func(queue *Queue, shouldDrain bool) {
					if shouldDrain {
						for len(queue.items) > 0 {
							(<-queue.items).Run(mCtx)
						}
					}
					q.done <- true
				}(q, drain)
				return
			}
		}
	}(result)

	return result
}

// Push push a task into the queue
func (q *Queue) Push(t Request) {
	if q.terminated {
		return
	}
	q.items <- t
}

// Close closes the queue in prep for its memory reclamation. the existing task being executed will complete its execution.
// when calling this, caller must ensure that no calls tasks are being pushed into the instance & no such calls will be done in the future.
// indicate whether the existing tasks should be executed or not
// when this function returns, the queue is closed
func (q *Queue) Close(drain bool) {
	q.terminated = true
	q.terminate <- drain
	<-q.done
}

// Count returns number of items left in queue for handling
func (q *Queue) Count() int {
	return len(q.items)
}
