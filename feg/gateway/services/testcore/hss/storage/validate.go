/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import "magma/lte/cloud/go/protos"

// validateSubscriberData ensures that a subscriber data proto is not nil and
// that it contains a valid subscriber id.
func validateSubscriberData(subscriber *protos.SubscriberData) error {
	if subscriber == nil {
		return NewInvalidArgumentError("Subscriber data cannot be nil")
	}
	if subscriber.Sid == nil {
		return NewInvalidArgumentError("Subscriber data must contain a subscriber id")
	}
	return validateSubscriberID(subscriber.Sid.Id)
}

// validateSubscriberID ensures that a subscriber ID can be stored
func validateSubscriberID(id string) error {
	if len(id) == 0 {
		return NewInvalidArgumentError("Subscriber id cannot be the empty string")
	}
	return nil
}
