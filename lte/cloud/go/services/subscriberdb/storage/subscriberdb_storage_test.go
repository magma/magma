/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package storage_test tests subscriberdb storage
package storage_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/test_utils"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

const testNetworkID = "subscriberdb_storage_test_network"

// initTestDB initializes a mock datastore and adds a subscriber
func initTestDB(t *testing.T) (*storage.SubscriberDBStorage, *protos.SubscriberData) {
	store, err := storage.NewSubscriberDBStorage(test_utils.NewMockDatastore())
	assert.NoError(t, err)

	networkID := orcprotos.NetworkID{Id: testNetworkID}
	sid := protos.SubscriberID{Id: "12345"}
	subs := &protos.SubscriberData{Sid: &sid, NetworkId: &networkID}

	ret, err := store.AddSubscriber(subs)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	return store, subs
}

func TestSubscriberDBStorageAddSubscriber(t *testing.T) {
	store, subs := initTestDB(t)

	// add a second subscriber to the db
	subs.Sid.Id = "67890"
	ret, err := store.AddSubscriber(subs)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	// ensure there is an error on duplicate additions
	ret, err = store.AddSubscriber(subs)
	assert.Error(t, err)
	assert.Nil(t, ret)

	// ensure multiple adds worked properly
	sids, err := store.ListSubscribers(subs.GetNetworkId())
	assert.NoError(t, err)

	if len(sids.GetSids()) != 2 {
		t.Errorf("Got %d subs, 2 expected...", len(sids.GetSids()))
	}

	expectedSids := []*protos.SubscriberID{
		{Id: "12345"},
		{Id: "67890"},
	}
	assert.ElementsMatch(t, expectedSids, sids.GetSids())
}

func TestSubscriberDBStorageDeleteSubscriber(t *testing.T) {
	store, subs := initTestDB(t)

	lookup := protos.SubscriberLookup{NetworkId: subs.GetNetworkId(), Sid: subs.GetSid()}

	// test deletion
	ret, err := store.DeleteSubscriber(&lookup)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	// ensure there is no error on duplicate deletion
	ret, err = store.DeleteSubscriber(&lookup)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	sids, err := store.ListSubscribers(subs.GetNetworkId())
	assert.NoError(t, err)

	if len(sids.Sids) != 0 {
		t.Errorf("Got %d subs, 0 expected...", len(sids.Sids))
	}
}

func TestSubscriberDBStorageUpdateSubscriber(t *testing.T) {
	store, subs := initTestDB(t)

	networkID := subs.GetNetworkId()
	lookup := protos.SubscriberLookup{NetworkId: networkID, Sid: subs.GetSid()}
	sid2 := protos.SubscriberID{Id: "11111"}
	newSubs := protos.SubscriberData{Sid: &sid2, NetworkId: networkID}
	lookup2 := protos.SubscriberLookup{NetworkId: networkID, Sid: &sid2}

	// test update on existing subscriber
	subs.Lte = &protos.LTESubscription{}
	ret, err := store.UpdateSubscriber(subs)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	subData, err := store.GetSubscriberData(&lookup)
	assert.NoError(t, err)
	assert.Equal(t, &protos.LTESubscription{}, subData.Lte)

	// test insert on update if subscriber doesn't exist
	ret, err = store.UpdateSubscriber(&newSubs)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	newSubData, err := store.GetSubscriberData(&lookup2)
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(&newSubs), orcprotos.TestMarshal(newSubData))
}

func TestSubscriberDBStorageGetSubscriberData(t *testing.T) {
	store, subs := initTestDB(t)

	lookup := protos.SubscriberLookup{NetworkId: subs.GetNetworkId(), Sid: subs.GetSid()}

	// test get subscriber data
	res, err := store.GetSubscriberData(&lookup)
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(subs), orcprotos.TestMarshal(res))

	ret, err := store.DeleteSubscriber(&lookup)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	// ensure error on get non-existent subscriber
	ret2, err := store.GetSubscriberData(&lookup)
	assert.Error(t, err)
	assert.Nil(t, ret2)
}

func TestSubscriberDBStorageListSubscribers(t *testing.T) {
	store, subs := initTestDB(t)

	networkID := subs.GetNetworkId()
	sid2 := protos.SubscriberID{Id: "54321"}
	sid3 := protos.SubscriberID{Id: "55555"}
	subs2 := protos.SubscriberData{Sid: &sid2, NetworkId: networkID}
	subs3 := protos.SubscriberData{Sid: &sid3, NetworkId: networkID}

	ret, err := store.AddSubscriber(&subs2)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	ret, err = store.AddSubscriber(&subs3)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	expectedSids := []*protos.SubscriberID{
		{Id: "12345"},
		{Id: "54321"},
		{Id: "55555"},
	}

	// test retrieval of all subscribers that have been added
	sids, err := store.ListSubscribers(networkID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedSids, sids.GetSids())
}

func TestSubscriberDBStorageGetAllSubscriberData(t *testing.T) {
	store, subs := initTestDB(t)

	networkID := subs.GetNetworkId()
	sid2 := protos.SubscriberID{Id: "54321"}
	sid3 := protos.SubscriberID{Id: "55555"}
	subs2 := protos.SubscriberData{Sid: &sid2, NetworkId: networkID}
	subs3 := protos.SubscriberData{Sid: &sid3, NetworkId: networkID}

	ret, err := store.AddSubscriber(&subs2)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	ret, err = store.AddSubscriber(&subs3)
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, ret)

	// test retrieval of all subscriber data that has been added
	allSubs, err := store.GetAllSubscriberData(networkID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []*protos.SubscriberData{
		{Sid: &protos.SubscriberID{Id: "12345"}, NetworkId: &orcprotos.NetworkID{Id: testNetworkID}},
		{Sid: &protos.SubscriberID{Id: "54321"}, NetworkId: &orcprotos.NetworkID{Id: testNetworkID}},
		{Sid: &protos.SubscriberID{Id: "55555"}, NetworkId: &orcprotos.NetworkID{Id: testNetworkID}},
	},
		allSubs.GetSubscribers())
}
