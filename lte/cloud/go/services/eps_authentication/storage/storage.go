/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package storage

import (
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/protos"
)

type SubscriberDBStorage interface {
	ListSubscribers(networkID *protos.NetworkID) (*lteprotos.SubscriberIDSet, error)
	GetSubscriberData(lookup *lteprotos.SubscriberLookup) (*lteprotos.SubscriberData, error)
	GetAllSubscriberData(networkID *protos.NetworkID) (*lteprotos.GetAllSubscriberDataResponse, error)
	AddSubscriber(subs *lteprotos.SubscriberData) (*protos.Void, error)
	UpdateSubscriber(subs *lteprotos.SubscriberData) (*protos.Void, error)
	DeleteSubscriber(lookup *lteprotos.SubscriberLookup) (*protos.Void, error)
}

// DEPRECATED -- temporarily deprecated, needs to be implemented
type subscriberDBStorageImpl struct{}

func NewSubscriberDBStorage() SubscriberDBStorage {
	return &subscriberDBStorageImpl{}
}

func (s *subscriberDBStorageImpl) ListSubscribers(networkID *protos.NetworkID) (*lteprotos.SubscriberIDSet, error) {
	panic("implement me")
}

func (s *subscriberDBStorageImpl) GetSubscriberData(lookup *lteprotos.SubscriberLookup) (*lteprotos.SubscriberData, error) {
	panic("implement me")
}

func (s *subscriberDBStorageImpl) GetAllSubscriberData(networkID *protos.NetworkID) (*lteprotos.GetAllSubscriberDataResponse, error) {
	panic("implement me")
}

func (s *subscriberDBStorageImpl) AddSubscriber(subs *lteprotos.SubscriberData) (*protos.Void, error) {
	panic("implement me")
}

func (s *subscriberDBStorageImpl) UpdateSubscriber(subs *lteprotos.SubscriberData) (*protos.Void, error) {
	panic("implement me")
}

func (s *subscriberDBStorageImpl) DeleteSubscriber(lookup *lteprotos.SubscriberLookup) (*protos.Void, error) {
	panic("implement me")
}
