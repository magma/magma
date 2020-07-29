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

package oc

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/plugin/ochttp"
)

func TestClientCache(t *testing.T) {
	assert.Equal(t, DefaultTransport, DefaultClient.Transport)
	assert.IsType(t, (*ochttp.Transport)(nil), DefaultTransport)
	assert.Nil(t, DefaultTransport.(*ochttp.Transport).FormatSpanName)
	assert.Equal(t, ClientFor(""), ClientFor(""))

	cache := &ClientCache{}
	send, receive := cache.ClientFor("send"), cache.ClientFor("receive")
	assert.Equal(t, "send", send.Transport.(*ochttp.Transport).FormatSpanName(nil))
	assert.Equal(t, "receive", receive.Transport.(*ochttp.Transport).FormatSpanName(nil))
	assert.Equal(t, 2, func() int {
		var count int
		cache.clients.Range(func(_, _ interface{}) bool { count++; return true })
		return count
	}())
}

func TestClientCacheContention(t *testing.T) {
	numGoroutine := runtime.NumCPU() * 16
	var wg sync.WaitGroup
	wg.Add(numGoroutine)

	var (
		clients sync.Map
		nStores int32
	)
	ops := []string{"get", "set", "drop", "create", "upsert"}
	for i := 0; i < numGoroutine; i++ {
		go func() {
			for n := 0; n < 32; n++ {
				for _, op := range ops {
					client := ClientFor(op)
					if existing, ok := clients.LoadOrStore(op, client); ok {
						assert.Equal(t, client, existing)
					} else {
						assert.True(t, atomic.AddInt32(&nStores, 1) <= int32(len(ops)))
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
