/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package subscriberdb_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/test_init"
	orcprotos "magma/orc8r/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestSubscriberdb(t *testing.T) {
	test_init.StartTestService(t)

	networkId := &orcprotos.NetworkID{Id: "test"}
	sid := &protos.SubscriberID{Id: "12345"}
	subs := &protos.SubscriberData{Sid: sid, NetworkId: networkId}

	ids, err := subscriberdb.ListSubscribers("test")
	assert.NoError(t, err)
	assert.Equal(t, ids, []string{})

	err = subscriberdb.AddSubscriber("test", subs)
	assert.NoError(t, err)
	ids, err = subscriberdb.ListSubscribers("test")
	assert.NoError(t, err)
	assert.Equal(t, ids, []string{"IMSI12345"})

	allSubs, err := subscriberdb.GetAllSubscriberData("test")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(allSubs))
	assert.Equal(t, orcprotos.TestMarshal(subs), orcprotos.TestMarshal(allSubs[0]))

	subs.Lte = &protos.LTESubscription{}
	err = subscriberdb.UpdateSubscriber("test", subs)
	assert.NoError(t, err)
	data, err := subscriberdb.GetSubscriber("test", "IMSI12345")
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(subs), orcprotos.TestMarshal(data))

	err = subscriberdb.DeleteSubscriber("test", "IMSI12345")
	assert.NoError(t, err)
	ids, err = subscriberdb.ListSubscribers("test")
	assert.NoError(t, err)
	assert.Equal(t, ids, []string{})
}
