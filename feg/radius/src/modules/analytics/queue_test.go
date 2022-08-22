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

import (
	"math/rand"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/stretchr/testify/require"
)

// maxTimeBetweenTasks the time of random sleep between firing tasks & for task execution.
const maxTimeBetweenTasks = 10 // msec

func TestAnalyticsQueue(t *testing.T) {
	var mCtx ModuleCtx
	numTasks := 100
	for withDrain := 0; withDrain < 2; withDrain++ {
		isDrain := false
		if withDrain == 1 {
			isDrain = true
		}
		lastExecTaskID := atomic.NewInt32(0)
		queue := NewAnalyticsQueue(mCtx)
		for i := 1; i <= numTasks; i++ {
			t := LazyExecSerializerOrderedTask{
				Testing:        t,
				Assertions:     require.New(t),
				LastExecTaskID: lastExecTaskID,
				TaskID:         i,
			}
			queue.Push(&t)
			sl := rand.Intn(maxTimeBetweenTasks)
			time.Sleep(time.Duration(sl) * time.Millisecond)
		}
		queue.Close(isDrain)

		// Assert
		queueSize := queue.Count()
		if isDrain {
			require.Equal(t, 0, queueSize)
		} else {
			require.True(t, numTasks >= queueSize)
		}
	}
}

// LazyExecSerializerOrderedTask a task that will test the ordered execution.
type LazyExecSerializerOrderedTask struct {
	Testing        *testing.T
	Assertions     *require.Assertions
	LastExecTaskID *atomic.Int32 // this needs to be atomic bcz go-routines that execute & update this might run on different cores
	TaskID         int
}

func (t *LazyExecSerializerOrderedTask) Run(_ ModuleCtx) {
	expectedID := 1 + t.LastExecTaskID.Load()
	// check order of tasks is maintained
	require.True(t.Testing, int32(t.TaskID) == expectedID, "expecting task ID %d but executing ID %d", expectedID, t.TaskID)
	// update the last executed task as self.
	t.LastExecTaskID.Store(int32(t.TaskID))
	sl := rand.Intn(maxTimeBetweenTasks)
	time.Sleep(time.Duration(sl) * time.Millisecond)
}
