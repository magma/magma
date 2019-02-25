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

func TestValidateSubscriberData(t *testing.T) {
	err := validateSubscriberData(nil)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data cannot be nil"), err)

	sub := &protos.SubscriberData{}
	err = validateSubscriberData(sub)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber data must contain a subscriber id"), err)

	sub = &protos.SubscriberData{Sid: &protos.SubscriberID{Id: ""}}
	err = validateSubscriberData(sub)
	assert.Exactly(t, NewInvalidArgumentError("Subscriber id cannot be the empty string"), err)

	sub = &protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	err = validateSubscriberData(sub)
	assert.NoError(t, err)
}

func TestValidateSubscriberID(t *testing.T) {
	err := validateSubscriberID("")
	assert.Exactly(t, NewInvalidArgumentError("Subscriber id cannot be the empty string"), err)

	err = validateSubscriberID("53425542542332")
	assert.NoError(t, err)
}
