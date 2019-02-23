/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/lte/cloud/go/protos"

	"github.com/golang/glog"
)

// MemorySubscriberStore is an in memory implementation of SubscriberStore.
type MemorySubscriberStore struct {
	accounts map[string]*protos.SubscriberData
}

// NewMemorySubscriberStore initializes a MemorySubscriberStore with an empty accounts map.
// Output: a new MemorySubscriberStore
func NewMemorySubscriberStore() *MemorySubscriberStore {
	return &MemorySubscriberStore{
		accounts: make(map[string]*protos.SubscriberData),
	}
}

// AddSubscriber tries to add this subscriber to the server.
// This function returns an AlreadyExists error if the subscriber has already
// been added.
// Input: The subscriber data which will be added.
func (store *MemorySubscriberStore) AddSubscriber(data *protos.SubscriberData) error {
	if err := validateSubscriberData(data); err != nil {
		return err
	}

	// Check that we are not adding a duplicate subscriber.
	id := data.GetSid().Id
	_, exists := store.accounts[id]
	if exists {
		glog.Errorf("Subscriber '%s' already added", id)
		return NewAlreadyExistsError(id)
	}

	store.accounts[data.GetSid().Id] = data
	return nil
}

// UpdateSubscriber changes the data stored for an existing subscriber.
// If the subscriber cannot be found, an error is returned instead.
// Input: The new subscriber data to store
func (store *MemorySubscriberStore) UpdateSubscriber(data *protos.SubscriberData) error {
	if err := validateSubscriberData(data); err != nil {
		return err
	}

	id := data.Sid.Id
	_, exists := store.accounts[id]
	if !exists {
		glog.Errorf("Subscriber '%s' not found", id)
		return NewUnknownSubscriberError(id)
	}

	store.accounts[id] = data
	return nil
}

// GetSubscriberData looks up a subscriber by their id.
// If the subscriber cannot be found, an error is returned instead.
// Input: The id of the subscriber to be looked up.
// Output: The data of the corresponding subscriber or an error.
func (store *MemorySubscriberStore) GetSubscriberData(id string) (*protos.SubscriberData, error) {
	data, exists := store.accounts[id]
	if exists {
		return data, nil
	}
	glog.Errorf("Subscriber '%s' not found", id)
	return nil, NewUnknownSubscriberError(id)
}

// DeleteSubscriber deletes a subscriber by their id.
// If the subscriber is not found, then this call is ignored.
// Input: The id of the subscriber to be deleted.
func (store *MemorySubscriberStore) DeleteSubscriber(id string) error {
	delete(store.accounts, id)
	return nil
}

// validateSubscriberData ensures that a subscriber data proto is not nil and
// that itcontains a subscriber id.
func validateSubscriberData(subscriber *protos.SubscriberData) error {
	if subscriber == nil {
		return NewInvalidArgumentError("Subscriber data cannot be nil")
	}

	if subscriber.Sid == nil {
		return NewInvalidArgumentError("Subscriber data must contain a subscriber id")
	}
	return nil
}
