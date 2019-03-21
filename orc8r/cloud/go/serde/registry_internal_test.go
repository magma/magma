/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package serde

import (
	"testing"

	"magma/orc8r/cloud/go/serde/mocks"

	"github.com/stretchr/testify/assert"
)

// https://www.youtube.com/watch?v=ndmB0bj7eyw, timestamp 32:00 for the
// concurrent testing technique used here

func TestRegisterSerdes(t *testing.T) {
	// Register 2 serdes on 2 new domains ["foo", "bar"] in the first call
	// Register 1 serde on ["foo"] in the second call
	// Pause before acquiring write lock in the first call, but proceed in
	// second call with no delay
	UnregisterAllSerdes(t)
	missingDomainsCalculatedCallback = func() {}
	defer func() {
		UnregisterAllSerdes(t)
		missingDomainsCalculatedCallback = func() {}
	}()

	fooSerde1 := &mocks.Serde{}
	fooSerde1.On("GetDomain").Return("foo")
	fooSerde1.On("GetType").Return("foo1")

	fooSerde2 := &mocks.Serde{}
	fooSerde2.On("GetDomain").Return("foo")
	fooSerde2.On("GetType").Return("foo2")

	barSerde := &mocks.Serde{}
	barSerde.On("GetDomain").Return("bar")
	barSerde.On("GetType").Return("bar1")

	waiter := make(chan error)
	missingDomainsCalculatedCallback = func() {
		// The first send signals that the callback has been reached and
		// begun execution so we don't clear the callback before the
		// goroutine for the first call gets to it.
		waiter <- nil

		// This second send will block until we receive from the channel again
		waiter <- nil
	}
	go func() {
		err := RegisterSerdes(fooSerde1, barSerde)
		// Signal to the test case that this call is finished
		waiter <- err
	}()

	// Clear the callback for the second call which should return immediately
	// Wait until the first call enters the blocking callback
	<-waiter
	missingDomainsCalculatedCallback = func() {}
	err := RegisterSerdes(fooSerde2)
	assert.NoError(t, err)

	// Only the second call should have gone through
	expected := &serdeRegistry{
		serdeRegistriesByDomain: map[string]*serdes{
			"foo": {
				serdesByKey: map[string]Serde{
					"foo2": fooSerde2,
				},
			},
		},
	}
	assert.Equal(t, expected, registry)

	// Unblock the first call
	<-waiter
	// Wait for the first call to finish
	err = <-waiter
	assert.NoError(t, err)
	expected = &serdeRegistry{
		serdeRegistriesByDomain: map[string]*serdes{
			"foo": {
				serdesByKey: map[string]Serde{
					"foo1": fooSerde1,
					"foo2": fooSerde2,
				},
			},
			"bar": {
				serdesByKey: map[string]Serde{
					"bar1": barSerde,
				},
			},
		},
	}
	assert.Equal(t, expected, registry)
}

func TestRegisterSerdesRollback(t *testing.T) {
	// Register 2 serdes on new domains ["foo", "bar"]
	// Register 1 new serde on foo, re-register serde on bar
	// 2nd call should result in a full rollback of the registry
	UnregisterAllSerdes(t)
	defer func() {
		UnregisterAllSerdes(t)
	}()

	fooSerde1 := &mocks.Serde{}
	fooSerde1.On("GetDomain").Return("foo")
	fooSerde1.On("GetType").Return("foo1")

	fooSerde2 := &mocks.Serde{}
	fooSerde2.On("GetDomain").Return("foo")
	fooSerde2.On("GetType").Return("foo2")

	barSerde := &mocks.Serde{}
	barSerde.On("GetDomain").Return("bar")
	barSerde.On("GetType").Return("bar1")

	err := RegisterSerdes(fooSerde1, barSerde)
	assert.NoError(t, err)
	err = RegisterSerdes(fooSerde2, barSerde)
	assert.EqualError(t, err, "Error registering serdes: Serde with key bar1 is already registered; registry has been rolled back")
	expected := &serdeRegistry{
		serdeRegistriesByDomain: map[string]*serdes{
			"foo": {
				serdesByKey: map[string]Serde{
					"foo1": fooSerde1,
				},
			},
			"bar": {
				serdesByKey: map[string]Serde{
					"bar1": barSerde,
				},
			},
		},
	}
	assert.Equal(t, expected, registry)
}

func TestRegisterSerdesRaceRollback(t *testing.T) {
	// Same as TestRegisterSerdesRollback except we pause the first call before
	// registering and allow the second call through. The conflict here is
	// expected to be on the first call when we allow it to continue.
	UnregisterAllSerdes(t)
	newDomainsCreatedCallback = func() {}
	defer func() {
		UnregisterAllSerdes(t)
		newDomainsCreatedCallback = func() {}
	}()

	fooSerde1 := &mocks.Serde{}
	fooSerde1.On("GetDomain").Return("foo")
	fooSerde1.On("GetType").Return("foo1")

	fooSerde2 := &mocks.Serde{}
	fooSerde2.On("GetDomain").Return("foo")
	fooSerde2.On("GetType").Return("foo2")

	barSerde := &mocks.Serde{}
	barSerde.On("GetDomain").Return("bar")
	barSerde.On("GetType").Return("bar1")

	// Check out TestRegisterSerdes for what we're doing here
	waiter := make(chan error)
	newDomainsCreatedCallback = func() {
		waiter <- nil
		waiter <- nil
	}
	go func() {
		err := RegisterSerdes(fooSerde1, barSerde)
		waiter <- err
	}()

	<-waiter
	newDomainsCreatedCallback = func() {}
	err := RegisterSerdes(fooSerde2, barSerde)
	assert.NoError(t, err)
	expected := &serdeRegistry{
		serdeRegistriesByDomain: map[string]*serdes{
			"foo": {
				serdesByKey: map[string]Serde{
					"foo2": fooSerde2,
				},
			},
			"bar": {
				serdesByKey: map[string]Serde{
					"bar1": barSerde,
				},
			},
		},
	}
	assert.Equal(t, expected, registry)

	<-waiter
	err = <-waiter
	assert.EqualError(t, err, "Error registering serdes: Serde with key bar1 is already registered; registry has been rolled back")
	assert.Equal(t, expected, registry)
}
