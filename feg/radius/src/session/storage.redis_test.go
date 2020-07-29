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

	"github.com/alicebob/miniredis"
)

func TestBasicInsertGetRedis(t *testing.T) {

	// Arrange
	sessionID := "test"

	mr, _ := miniredis.Run()
	defer mr.Close()

	storage := NewMultiSessionRedisStorage(mr.Addr(), "", 0)

	// Act and Assert
	performSignleReadWriteDeleteReadTest(t, storage, sessionID)
}

func TestMultipleConcurrentInsertDeleteGetRedis(t *testing.T) {
	// Arrange
	degOfParallelism := 100
	reqPerConcurrentContext := 100
	onComplete := sync.WaitGroup{}

	mr, _ := miniredis.Run()
	defer mr.Close()

	storage := NewMultiSessionRedisStorage(mr.Addr(), "", 0)

	// Act
	for i := 0; i < degOfParallelism; i++ {
		go func(called string, calling string) {
			sessionID := fmt.Sprintf("session_%s_%s", calling, called)
			loopReadWriteDelete(t, storage, sessionID, reqPerConcurrentContext, &onComplete)
		}(fmt.Sprintf("called%d", i), fmt.Sprintf("calling%d", i))
	}
	onComplete.Wait()
}
