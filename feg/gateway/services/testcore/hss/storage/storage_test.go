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
	"magma/lte/cloud/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"
)

// SubscriberStoreTestSuite is a test suite which can be run against any implementation
// of the SubscriberStore interface.
type SubscriberStoreTestSuite struct {
	suite.Suite
	store SubscriberStore

	// createStore is run before every test to recreate the subscriber store
	createStore func() SubscriberStore
}

func (suite *SubscriberStoreTestSuite) SetupTest() {
	suite.store = suite.createStore()
}

func (suite *SubscriberStoreTestSuite) TestAddSubscriber() {
	store := suite.store

	sub1 := &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	sub2 := &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "2"}}

	err := store.AddSubscriber(sub1)
	suite.NoError(err)

	err = store.AddSubscriber(sub1)
	suite.Exactly(NewAlreadyExistsError("1"), err)

	err = store.AddSubscriber(sub2)
	suite.NoError(err)

	err = store.AddSubscriber(sub1)
	suite.Exactly(NewAlreadyExistsError("1"), err)

	err = store.AddSubscriber(nil)
	suite.Exactly(NewInvalidArgumentError("Subscriber data cannot be nil"), err)

	sub := &protos.SubscriberData{}
	err = store.AddSubscriber(sub)
	suite.Exactly(NewInvalidArgumentError("Subscriber data must contain a subscriber id"), err)
}

func (suite *SubscriberStoreTestSuite) TestGetSubscriberData() {
	store := suite.store
	sub := protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}

	_, err := store.GetSubscriberData("1")
	suite.Exactly(NewUnknownSubscriberError("1"), err)

	err = store.AddSubscriber(&sub)
	suite.NoError(err)

	result, err := store.GetSubscriberData("1")
	suite.NoError(err)
	suite.True(proto.Equal(&sub, result))
}

func (suite *SubscriberStoreTestSuite) TestUpdateSubscriberData() {
	store := suite.store

	err := store.UpdateSubscriber(nil)
	suite.Exactly(NewInvalidArgumentError("Subscriber data cannot be nil"), err)

	sub := &protos.SubscriberData{}
	err = store.UpdateSubscriber(sub)
	suite.Exactly(NewInvalidArgumentError("Subscriber data must contain a subscriber id"), err)

	sub = &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	err = store.UpdateSubscriber(sub)
	suite.Exactly(NewUnknownSubscriberError("1"), err)

	err = store.AddSubscriber(sub)
	suite.NoError(err)

	updatedSub := &protos.SubscriberData{
		Sid:        &protos.SubscriberID{Id: "1"},
		SubProfile: "test",
	}
	err = store.UpdateSubscriber(updatedSub)
	suite.NoError(err)

	retrievedSub, err := store.GetSubscriberData("1")
	suite.NoError(err)
	suite.True(proto.Equal(updatedSub, retrievedSub))
}

func (suite *SubscriberStoreTestSuite) TestDeleteSubscriber() {
	store := suite.store
	sub := protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}

	err := store.AddSubscriber(&sub)
	suite.NoError(err)

	result, err := store.GetSubscriberData("1")
	suite.NoError(err)
	suite.True(proto.Equal(&sub, result))

	err = store.DeleteSubscriber("1")
	suite.NoError(err)

	_, err = store.GetSubscriberData("1")
	suite.Exactly(NewUnknownSubscriberError("1"), err)
}

func (suite *SubscriberStoreTestSuite) TestDeleteAllSubscribers() {
	store := suite.store

	err := store.AddSubscriber(&protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}})
	suite.NoError(err)

	err = store.AddSubscriber(&protos.SubscriberData{Sid: &protos.SubscriberID{Id: "2"}})
	suite.NoError(err)

	err = store.DeleteAllSubscribers()
	suite.NoError(err)

	_, err = store.GetSubscriberData("1")
	suite.Exactly(NewUnknownSubscriberError("1"), err)

	_, err = store.GetSubscriberData("2")
	suite.Exactly(NewUnknownSubscriberError("2"), err)
}

func (suite *SubscriberStoreTestSuite) TestRaceCondition() {
	store := suite.store
	sub := &protos.SubscriberData{
		Sid:   &protos.SubscriberID{Id: "1"},
		State: &protos.SubscriberState{LteAuthNextSeq: 0},
	}
	err := store.AddSubscriber(sub)
	suite.NoError(err)

	writers := uint64(3)
	doneSignal := make(chan struct{})

	for i := uint64(1); i <= writers; i++ {
		localSub := proto.Clone(sub).(*protos.SubscriberData)
		localSub.State.LteAuthNextSeq = i

		go func() {
			err := store.UpdateSubscriber(localSub)
			suite.NoError(err)
			doneSignal <- struct{}{}
		}()
	}

	for i := uint64(0); i < writers; i++ {
		<-doneSignal
	}

	data, err := store.GetSubscriberData("1")
	suite.NoError(err)
	seq := data.GetState().GetLteAuthNextSeq()
	if seq <= 0 || seq > writers {
		suite.Fail("invalid seq: %d", seq)
	}
}
