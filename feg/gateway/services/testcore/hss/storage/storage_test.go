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
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"google.golang.org/genproto/protobuf/field_mask"

	"magma/lte/cloud/go/protos"
	orc8rprotos "magma/orc8r/lib/go/protos"
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
	suite.Exactly(NewInvalidArgumentError("Update request cannot be nil"), err)

	subUpdt := &protos.SubscriberUpdate{Data: &protos.SubscriberData{}}
	err = store.UpdateSubscriber(subUpdt)
	suite.Exactly(NewInvalidArgumentError("Subscriber data must contain a subscriber id"), err)

	sub := &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	subUpdt = &protos.SubscriberUpdate{Data: sub}

	err = store.UpdateSubscriber(subUpdt)
	suite.Exactly(NewUnknownSubscriberError("1"), err)

	err = store.AddSubscriber(sub)
	suite.NoError(err)

	updatedSub := &protos.SubscriberUpdate{
		Data: &protos.SubscriberData{
			Sid:        &protos.SubscriberID{Id: "1"},
			SubProfile: "test",
		},
	}
	err = store.UpdateSubscriber(updatedSub)
	suite.NoError(err)

	retrievedSub, err := store.GetSubscriberData("1")
	suite.NoError(err)
	suite.True(proto.Equal(updatedSub.Data, retrievedSub))
}

func (suite *SubscriberStoreTestSuite) TestPartialUpdateSubscriberData() {
	store := suite.store

	// initial subscriber
	sub := &protos.SubscriberData{
		Sid: &protos.SubscriberID{Id: "1"},
		Gsm: &protos.GSMSubscription{
			State:      1,
			AuthAlgo:   1,
			AuthKey:    []byte{1, 1, 1},
			AuthTuples: nil,
		},
		Lte: &protos.LTESubscription{
			State:             protos.LTESubscription_ACTIVE,
			AuthAlgo:          1,
			AuthKey:           []byte{1, 1, 1},
			AuthOpc:           []byte{9, 9, 9},
			AssignedBaseNames: nil,
			AssignedPolicies:  nil,
		},
	}

	// update request with mask (only update fields in mask)
	updatedSub := &protos.SubscriberUpdate{
		// Mask indicates the fields that should be updated
		Mask: &field_mask.FieldMask{
			Paths: []string{"SubProfile", "Lte.AuthAlgo", "Non_3Gpp"},
		},
		Data: &protos.SubscriberData{
			Sid:        &protos.SubscriberID{Id: "1"},
			SubProfile: "test",
			Lte: &protos.LTESubscription{
				AuthAlgo: 0,
				AuthKey:  []byte{2, 2, 2},
			},
			Non_3Gpp: &protos.Non3GPPUserProfile{
				Msisdn:              "12345",
				Non_3GppIpAccess:    10,
				Non_3GppIpAccessApn: 11,
			},
			// NetworkId should not be updated since it is not in Mask
			NetworkId: &orc8rprotos.NetworkID{Id: "test"},
		},
	}

	// expected result after update with mask
	expectedRes := &protos.SubscriberData{
		Sid:        &protos.SubscriberID{Id: "1"},
		SubProfile: "test", // Modified (new)
		Gsm: &protos.GSMSubscription{
			State:      1,
			AuthAlgo:   1,
			AuthKey:    []byte{1, 1, 1},
			AuthTuples: nil,
		},
		Lte: &protos.LTESubscription{
			State:             protos.LTESubscription_ACTIVE,
			AuthAlgo:          0,               // Modified
			AuthKey:           []byte{1, 1, 1}, // Not Modified
			AuthOpc:           []byte{9, 9, 9},
			AssignedBaseNames: nil,
			AssignedPolicies:  nil,
		},
		Non_3Gpp: &protos.Non3GPPUserProfile{ // Modified (new)
			Msisdn:              "12345",
			Non_3GppIpAccess:    10,
			Non_3GppIpAccessApn: 11,
		},
	}

	err := store.AddSubscriber(sub)
	suite.NoError(err)

	err = store.UpdateSubscriber(updatedSub)
	suite.NoError(err)

	retrievedSub, err := store.GetSubscriberData("1")
	suite.NoError(err)

	suite.True(proto.Equal(expectedRes, retrievedSub))
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
			err := store.UpdateSubscriber(&protos.SubscriberUpdate{Data: localSub})
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
