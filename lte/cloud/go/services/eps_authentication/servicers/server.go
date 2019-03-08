/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"

	"magma/lte/cloud/go/services/subscriberdb/storage"
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
