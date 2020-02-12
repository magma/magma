/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"

	lteprotos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	orc8rprotos "magma/orc8r/lib/go/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EPSAuthServer struct {
	Store *storage.SubscriberDBStorage
}

// NewEPSAuthServer returns a Server with the provided store.
func NewEPSAuthServer(store *storage.SubscriberDBStorage) (*EPSAuthServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Cannot initialize eps authentication server with nil store")
	}
	return &EPSAuthServer{Store: store}, nil
}

// lookupSubscriber returns a subscriber's data or an error.
func (srv *EPSAuthServer) lookupSubscriber(userName, networkID string) (*lteprotos.SubscriberData, lteprotos.ErrorCode, error) {
	lookup := &lteprotos.SubscriberLookup{
		Sid:       &lteprotos.SubscriberID{Id: userName},
		NetworkId: &orc8rprotos.NetworkID{Id: networkID},
	}
	subscriber, err := srv.Store.GetSubscriberData(lookup)
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return nil, lteprotos.ErrorCode_USER_UNKNOWN, err
		}
		return nil, lteprotos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE, err
	}
	return subscriber, lteprotos.ErrorCode_SUCCESS, nil
}
