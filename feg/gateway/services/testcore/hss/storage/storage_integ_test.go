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

package storage_test

import (
	"testing"

	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/lte/cloud/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestTestcoreStorageMemory_Integration(t *testing.T) {
	store := storage.NewMemorySubscriberStore()
	testTestcoreStorageImpl(t, store)
}

func testTestcoreStorageImpl(t *testing.T, store storage.SubscriberStore) {
	sub0 := "subscriber_id_0"
	sub1 := "subscriber_id_1"

	data0 := &protos.SubscriberData{
		Sid:        &protos.SubscriberID{Id: sub0},
		SubProfile: "some_sub_profile",
	}
	data0Updated := &protos.SubscriberData{
		Sid:        &protos.SubscriberID{Id: sub0},
		SubProfile: "some_new_sub_profile",
	}

	data1 := &protos.SubscriberData{
		Sid:        &protos.SubscriberID{Id: sub1},
		SubProfile: "some_alternate_sub_profile",
	}

	// Initially no subscribers
	_, err := store.GetSubscriberData(sub0)
	assert.Error(t, err)

	// Idempotent deletions
	err = store.DeleteAllSubscribers()
	assert.NoError(t, err)
	err = store.DeleteSubscriber(sub0)
	assert.NoError(t, err)

	// Add, get data
	err = store.AddSubscriber(data0)
	assert.NoError(t, err)
	dataRecvd, err := store.GetSubscriberData(sub0)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(dataRecvd, data0))

	// Update, get data
	err = store.UpdateSubscriber(data0Updated)
	assert.NoError(t, err)
	dataRecvd, err = store.GetSubscriberData(sub0)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(dataRecvd, data0Updated))

	// Add, delete, get data
	err = store.AddSubscriber(data1)
	assert.NoError(t, err)
	err = store.DeleteSubscriber(sub1)
	assert.NoError(t, err)
	_, err = store.GetSubscriberData(sub1)
	assert.Error(t, err)
	dataRecvd, err = store.GetSubscriberData(sub0)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(dataRecvd, data0Updated))

	// Add back data1, delete all, get
	err = store.AddSubscriber(data1)
	assert.NoError(t, err)
	err = store.DeleteAllSubscribers()
	assert.NoError(t, err)
	_, err = store.GetSubscriberData(sub0)
	assert.Error(t, err)
	_, err = store.GetSubscriberData(sub1)
	assert.Error(t, err)
}
