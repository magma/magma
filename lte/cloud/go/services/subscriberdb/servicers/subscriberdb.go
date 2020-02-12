/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
Subscriberdb servicer provides the gRPC interface for the REST and
services to interact with the Subscriber data.

The servicer require a backing Datastore (which is typically Postgres)
for storing and retrieving the data.
*/
package servicers

import (
	"fmt"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	orcprotos "magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
)

type SubscriberDBServer struct {
	store *storage.SubscriberDBStorage
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func NewSubscriberDBServer(store *storage.SubscriberDBStorage) (*SubscriberDBServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Cannot initialize SubscriberDBServer with Nil store")
	}
	return &SubscriberDBServer{store}, nil
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func (srv *SubscriberDBServer) AddSubscriber(
	ctx context.Context,
	subs *protos.SubscriberData,
) (*orcprotos.Void, error) {
	if err := validateSubscriberData(subs); err != nil {
		return nil, err
	}
	return srv.store.AddSubscriber(subs)
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func (srv *SubscriberDBServer) DeleteSubscriber(
	ctx context.Context,
	lookup *protos.SubscriberLookup,
) (*orcprotos.Void, error) {
	if err := validateSubscriberLookup(lookup); err != nil {
		return nil, err
	}
	return srv.store.DeleteSubscriber(lookup)
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func (srv *SubscriberDBServer) UpdateSubscriber(
	ctx context.Context,
	subs *protos.SubscriberData,
) (*orcprotos.Void, error) {
	if err := validateSubscriberData(subs); err != nil {
		return nil, err
	}
	return srv.store.UpdateSubscriber(subs)
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func (srv *SubscriberDBServer) GetSubscriberData(
	ctx context.Context,
	lookup *protos.SubscriberLookup,
) (*protos.SubscriberData, error) {
	if err := validateSubscriberLookup(lookup); err != nil {
		return nil, err
	}
	return srv.store.GetSubscriberData(lookup)
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func (srv *SubscriberDBServer) ListSubscribers(
	ctx context.Context,
	networkID *orcprotos.NetworkID,
) (*protos.SubscriberIDSet, error) {
	if networkID == nil {
		return nil, fmt.Errorf("No network ID provided")
	}
	return srv.store.ListSubscribers(networkID)
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func (srv *SubscriberDBServer) GetAllSubscriberData(
	ctx context.Context,
	networkID *orcprotos.NetworkID,
) (*protos.GetAllSubscriberDataResponse, error) {
	if networkID == nil {
		return nil, fmt.Errorf("No network ID provided")
	}
	return srv.store.GetAllSubscriberData(networkID)
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func validateSubscriberLookup(lookup *protos.SubscriberLookup) error {
	if lookup == nil {
		return fmt.Errorf("No subscriber data provided")
	}
	if lookup.GetSid() == nil {
		return fmt.Errorf("No subscriber ID provided")
	}
	if lookup.GetNetworkId() == nil {
		return fmt.Errorf("No network ID provided")
	}
	return nil
}

// DEPRECATED
// THIS ENTIRE SERVICER IS DEPRECATED
func validateSubscriberData(subs *protos.SubscriberData) error {
	if subs == nil {
		return fmt.Errorf("No subscriber data provided")
	}
	if subs.GetSid() == nil {
		return fmt.Errorf("No subscriber ID provided")
	}
	if subs.GetNetworkId() == nil {
		return fmt.Errorf("No network ID provided")
	}
	return nil
}
