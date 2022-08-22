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

package session

import (
	"fmt"
	"sync"
	"testing"
)

func TestBasicInsertGet(t *testing.T) {
	// Arrange
	storage := NewMultiSessionMemoryStorage()

	// Act and Assert
	performSignleReadWriteDeleteReadTest(t, storage, "test")
}

func TestMultipleConcurrentInsertDeleteGet(t *testing.T) {
	// Arrange
	degOfParallelism := 100
	reqPerConcurrentContext := 100
	onComplete := sync.WaitGroup{}
	storage := NewMultiSessionMemoryStorage()

	// Act
	for i := 0; i < degOfParallelism; i++ {
		go func(called string, calling string) {
			sessionID := fmt.Sprintf("session_%s_%s", calling, called)
			loopReadWriteDelete(t, storage, sessionID, reqPerConcurrentContext, &onComplete)
		}(fmt.Sprintf("called%d", i), fmt.Sprintf("calling%d", i))
	}
	onComplete.Wait()

	// Assert
	// nothing to do (assert will happen in the go routines spawned above)
}
