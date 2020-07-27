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

package storage

import (
	lteprotos "magma/lte/cloud/go/protos"
	orc8rprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

type subscriberListener struct {
	store SubscriberStore
}

func NewSubscriberListener(store SubscriberStore) *subscriberListener {
	return &subscriberListener{store: store}
}

func (listener *subscriberListener) GetName() string {
	return "subscriberdb"
}

func (listener *subscriberListener) ReportError(err error) error {
	glog.Errorf("hss subscriber stream error: %s", err.Error())
	return nil
}

func (listener *subscriberListener) Update(batch *orc8rprotos.DataUpdateBatch) bool {
	glog.V(2).Infof("streaming %d subscriber update(s)", len(batch.GetUpdates()))
	store := listener.store

	if batch.GetResync() {
		err := store.DeleteAllSubscribers()
		if err != nil {
			glog.Errorf("failed to clear subscriber database: %s", err.Error())
		}
	}

	for _, update := range batch.GetUpdates() {
		subscriber := &lteprotos.SubscriberData{}
		if err := proto.Unmarshal(update.GetValue(), subscriber); err != nil {
			glog.Errorf("failed to unmarshal subscriber update for %s: %s", update.GetKey(), err.Error())
			continue
		}

		id := subscriber.GetSid().GetId()
		oldSub, err := store.GetSubscriberData(id)
		if err == nil {
			if oldSub.State != nil {
				subscriber.State = oldSub.State
			}
			err = store.UpdateSubscriber(subscriber)
			glog.Errorf("failed to update subscriber(%s): %s", id, err.Error())
		} else {
			err = store.AddSubscriber(subscriber)
			glog.Errorf("failed to add subscriber(%s): %s", id, err.Error())
		}
	}

	return true
}

func (listener *subscriberListener) GetExtraArgs() *any.Any {
	return nil
}
