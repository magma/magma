/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/test_utils"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSubscriberdb(t *testing.T) {
	ds := test_utils.NewMockDatastore()
	subscriberDBStore, err := storage.NewSubscriberDBStorage(ds)
	assert.NoError(t, err)
	ctx := context.Background()

	networkId := orcprotos.NetworkID{Id: "test"}
	sid := protos.SubscriberID{Id: "12345"}
	lookup := protos.SubscriberLookup{NetworkId: &networkId, Sid: &sid}
	subs := protos.SubscriberData{Sid: &sid, NetworkId: &networkId}

	srv, err := servicers.NewSubscriberDBServer(subscriberDBStore)
	assert.NoError(t, err)

	_, err = srv.AddSubscriber(ctx, &subs)
	assert.NoError(t, err)
	subs.Sid.Id = "67890"
	_, err = srv.AddSubscriber(ctx, &subs)
	assert.NoError(t, err)

	sids, err := srv.ListSubscribers(ctx, &networkId)
	assert.NoError(t, err)

	if len(sids.Sids) != 2 {
		t.Fatalf("Got %d subs, 2 expected...", len(sids.Sids))
	}

	_, err = srv.DeleteSubscriber(ctx, &lookup)
	assert.NoError(t, err)
	subs.Sid.Id = "12345"
	_, err = srv.DeleteSubscriber(ctx, &lookup)
	assert.NoError(t, err)

	// Add a subscriber
	_, err = srv.AddSubscriber(ctx, &subs)
	assert.NoError(t, err)
	_, err = srv.AddSubscriber(ctx, &subs)
	assert.Error(t, err) // duplicate addition

	res, err := srv.GetSubscriberData(ctx, &lookup)
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(&subs), orcprotos.TestMarshal(res))

	// Update the subscriber
	subs.Lte = &protos.LTESubscription{}
	_, err = srv.UpdateSubscriber(ctx, &subs)
	assert.NoError(t, err)

	res, err = srv.GetSubscriberData(ctx, &lookup)
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(&subs), orcprotos.TestMarshal(res))

	// List the subscribers
	sids, err = srv.ListSubscribers(ctx, &networkId)
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(sids), orcprotos.TestMarshal(&protos.SubscriberIDSet{
		Sids: []*protos.SubscriberID{&sid}}))
	allSubs, err := srv.GetAllSubscriberData(ctx, &networkId)
	assert.NoError(t, err)
	assert.Equal(
		t,
		orcprotos.TestMarshal(&protos.GetAllSubscriberDataResponse{Subscribers: []*protos.SubscriberData{&subs}}),
		orcprotos.TestMarshal(allSubs),
	)

	// Delete the subscriber
	_, err = srv.DeleteSubscriber(ctx, &lookup)
	assert.NoError(t, err)
	_, err = srv.GetSubscriberData(ctx, &lookup)
	assert.Error(t, err) // subscriber already removed

	sids, err = srv.ListSubscribers(ctx, &networkId)
	assert.NoError(t, err)
	assert.Equal(
		t, orcprotos.TestMarshal(sids), orcprotos.TestMarshal(&protos.SubscriberIDSet{Sids: []*protos.SubscriberID{}}))

}
