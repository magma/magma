/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package alert_test

import (
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/alert/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileLocker_Lock(t *testing.T) {
	locks, err := alert.NewFileLocker(newHealthyDirClient("test"))
	assert.NoError(t, err)
	fname := "file1"
	var events []int

	locks.Lock(fname)
	events = append(events, 1)
	go func() {
		locks.Lock(fname)
		events = append(events, 2)
	}()
	events = append(events, 3)

	// Assert event 2 has not happened yet since it is locked
	assert.Equal(t, []int{1, 3}, events)
}

func TestFileLocker_Unlock(t *testing.T) {
	locks, err := alert.NewFileLocker(newHealthyDirClient("test"))
	assert.NoError(t, err)
	fname := "file1"
	var events []int

	locks.Lock(fname)
	events = append(events, 1)
	go func() {
		locks.Lock(fname)
		events = append(events, 2)
	}()
	events = append(events, 3)
	locks.Unlock(fname)

	// sleep 1ms to let go func finish
	time.Sleep(1 * time.Millisecond)
	// Assert event 2 happened after event 3
	assert.Equal(t, []int{1, 3, 2}, events)
}

func TestFileLocker_RLock(t *testing.T) {
	locks, err := alert.NewFileLocker(newHealthyDirClient("test"))
	assert.NoError(t, err)
	fname := "file1"
	var events []int

	locks.RLock(fname)
	events = append(events, 1)
	go func() {
		locks.RLock(fname)
		events = append(events, 2)
	}()
	time.Sleep(1 * time.Millisecond)
	events = append(events, 3)
	go func() {
		locks.Lock(fname)
		events = append(events, 4)
	}()
	// Fname can be RLocked multiple times so events are in order, but cannot
	// be Locked so event 4 did not happen
	assert.Equal(t, []int{1, 2, 3}, events)
}

func TestFileLocker_RUnlock(t *testing.T) {
	locks, err := alert.NewFileLocker(newHealthyDirClient("test"))
	assert.NoError(t, err)
	fname := "file1"
	var events []int

	locks.RLock(fname)
	events = append(events, 1)
	go func() {
		locks.RLock(fname)
		events = append(events, 2)
		locks.RUnlock(fname)
	}()
	time.Sleep(1 * time.Millisecond)

	go func() {
		locks.Lock(fname)
		events = append(events, 3)
	}()

	events = append(events, 4)
	locks.RUnlock(fname)
	time.Sleep(1 * time.Millisecond)

	// Assert event 3 happened after 4
	assert.Equal(t, []int{1, 2, 4, 3}, events)
}

// creates mock directory client that doesn't return errors
func newHealthyDirClient(rulesDir string) *mocks.DirectoryClient {
	client := &mocks.DirectoryClient{}
	client.On("Stat", mock.AnythingOfType("string")).Return(nil, nil)
	client.On("Mkdir", mock.AnythingOfType("os.FileMode")).Return(nil)
	client.On("ReadDir").Return(nil, nil)
	client.On("Dir").Return(rulesDir)

	return client
}
