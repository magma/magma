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

// package protos includes generated GRPC sources as well as corresponding helper functions
package protos

import (
	"strconv"
	"sync"
)

const sidLocksSize = 64

var subscriberLocks [sidLocksSize]sync.Mutex

func (sid *SubscriberID) hash() uint64 {
	h, _ := strconv.ParseUint(sid.GetId(), 10, 64)
	h += uint64(sid.GetType()) << 46
	return h
}

// Lock - lockable interface implementation
func (subscr *SubscriberData) Lock() {
	subscriberLocks[subscr.GetSid().hash()%sidLocksSize].Lock()
}

// Unlock - lockable interface implementation
func (subscr *SubscriberData) Unlock() {
	subscriberLocks[subscr.GetSid().hash()%sidLocksSize].Unlock()
}
