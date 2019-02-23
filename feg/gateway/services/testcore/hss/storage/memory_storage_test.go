/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"testing"

	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestAddSubscriber(t *testing.T) {
	store := NewMemorySubscriberStore()

	sub1 := &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	sub2 := &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "2"}}

	err := store.AddSubscriber(sub1)
	assert.NoError(t, err)

	err = store.AddSubscriber(sub1)
	assert.Exactly(t, NewAlreadyExistsError("1"), err)

	err = store.AddSubscriber(sub2)
	assert.NoError(t, err)

	err = store.AddSubscriber(sub1)
	assert.Exactly(t, NewAlreadyExistsError("1"), err)

	err = store.AddSubscriber(nil)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data cannot be nil"), err)

	sub := &protos.SubscriberData{}
	err = store.AddSubscriber(sub)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data must contain a subscriber id"), err)
}

func TestGetSubscriberData(t *testing.T) {
	store := NewMemorySubscriberStore()
	sub := protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}

	_, err := store.GetSubscriberData("1")
	assert.Exactly(t, NewUnknownSubscriberError("1"), err)

	err = store.AddSubscriber(&sub)
	assert.NoError(t, err)

	result, err := store.GetSubscriberData("1")
	assert.NoError(t, err)
	assert.Equal(t, sub, *result)
}

func TestUpdateSubscriberData(t *testing.T) {
	store := NewMemorySubscriberStore()

	err := store.UpdateSubscriber(nil)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data cannot be nil"), err)

	sub := &protos.SubscriberData{}
	err = store.UpdateSubscriber(sub)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data must contain a subscriber id"), err)

	sub = &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	err = store.UpdateSubscriber(sub)
	assert.Exactly(t, NewUnknownSubscriberError("1"), err)

	err = store.AddSubscriber(sub)
	assert.NoError(t, err)

	updatedSub := &protos.SubscriberData{
		Sid:        &protos.SubscriberID{Id: "1"},
		SubProfile: "test",
	}
	err = store.UpdateSubscriber(updatedSub)
	assert.NoError(t, err)

	retreivedSub, err := store.GetSubscriberData("1")
	assert.NoError(t, err)
	assert.Equal(t, updatedSub, retreivedSub)
}

func TestDeleteSubscriber(t *testing.T) {
	store := NewMemorySubscriberStore()
	sub := protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}

	err := store.AddSubscriber(&sub)
	assert.NoError(t, err)

	result, err := store.GetSubscriberData("1")
	assert.NoError(t, err)
	assert.Equal(t, sub, *result)

	err = store.DeleteSubscriber("1")
	assert.NoError(t, err)

	_, err = store.GetSubscriberData("1")
	assert.Exactly(t, NewUnknownSubscriberError("1"), err)
}

func TestValidateSubscriberData(t *testing.T) {
	err := validateSubscriberData(nil)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data cannot be nil"), err)

	sub := &protos.SubscriberData{}
	err = validateSubscriberData(sub)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data must contain a subscriber id"), err)

	sub = &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	err = validateSubscriberData(sub)
	assert.NoError(t, err)
}
