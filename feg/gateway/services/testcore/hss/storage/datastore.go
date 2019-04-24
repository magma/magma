/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/datastore"

	"github.com/golang/protobuf/proto"
)

const tableName = "subscribers_hss"

// SubscriberDataStore is an implementation of SubscriberStore using a datastore.Api as the backend.
type SubscriberDataStore struct {
	api datastore.Api
}

// NewSubscriberDataStore initializes a SubscriberDataStore & returns it.
func NewSubscriberDataStore(api datastore.Api) *SubscriberDataStore {
	return &SubscriberDataStore{api: api}
}

// AddSubscriber tries to add this subscriber to the server.
// This function returns an AlreadyExists error if the subscriber has already
// been added.
// Input: The subscriber data which will be added.
func (store *SubscriberDataStore) AddSubscriber(data *protos.SubscriberData) error {
	if err := validateSubscriberData(data); err != nil {
		return err
	}

	// Check that we are not adding a duplicate subscriber.
	id := data.GetSid().Id
	exists, err := store.api.DoesKeyExist(tableName, id)
	if err != nil {
		return err
	}
	if exists {
		return NewAlreadyExistsError(id)
	}

	marshaledData, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return store.api.Put(tableName, id, marshaledData)
}

// UpdateSubscriber changes the data stored for an existing subscriber.
// If the subscriber cannot be found, an error is returned instead.
// Input: The new subscriber data to store
func (store *SubscriberDataStore) UpdateSubscriber(data *protos.SubscriberData) error {
	if err := validateSubscriberData(data); err != nil {
		return err
	}

	// Check that the subscriber exists.
	id := data.GetSid().Id
	exists, err := store.api.DoesKeyExist(tableName, id)
	if err != nil {
		return err
	}
	if !exists {
		return NewUnknownSubscriberError(id)
	}

	marshaledData, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return store.api.Put(tableName, id, marshaledData)
}

// GetSubscriberData looks up a subscriber by their id.
// If the subscriber cannot be found, an error is returned instead.
// Input: The id of the subscriber to be looked up.
// Output: The data of the corresponding subscriber or an error.
func (store *SubscriberDataStore) GetSubscriberData(id string) (*protos.SubscriberData, error) {
	if err := validateSubscriberID(id); err != nil {
		return nil, err
	}

	bytes, _, err := store.api.Get(tableName, id)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, NewUnknownSubscriberError(id)
		}
		return nil, err
	}

	var subscriber = new(protos.SubscriberData)
	err = proto.Unmarshal(bytes, subscriber)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}

// DeleteSubscriber deletes a subscriber by their id.
// If the subscriber is not found, then this call is ignored.
// Input: The id of the subscriber to be deleted.
func (store *SubscriberDataStore) DeleteSubscriber(id string) error {
	if err := validateSubscriberID(id); err != nil {
		return err
	}
	return store.api.Delete(tableName, id)
}

func (store *SubscriberDataStore) DeleteAllSubscribers() error {
	subs, err := store.api.ListKeys(tableName)
	if err != nil {
		return err
	}
	_, err = store.api.DeleteMany(tableName, subs)
	return err
}
