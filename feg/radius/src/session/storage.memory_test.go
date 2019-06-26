/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package session

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicInsertGet(t *testing.T) {
	// Arrange
	storage := NewMultiSessionMemoryStorage()

	// Act and Assert
	performSignleReadWriteDeleteReadTest(t, storage, "test")
}

func performSignleReadWriteDeleteReadTest(t *testing.T, storage GlobalStorage, sessionID string) {
	// Arrange
	msisdn := fmt.Sprintf("+%d", rand.Intn(999999))

	// Act
	stateBeforeWrite, errBeforeWrite := storage.Get(sessionID)
	writeErr := storage.Set(sessionID, State{
		MSISDN: msisdn,
	})
	stateAfterWrite, errAfterWrite := storage.Get(sessionID)
	storage.Reset(sessionID)
	stateAfterReset, errAfterReset := storage.Get(sessionID)

	// Assert
	require.Equal(t, nil, writeErr)
	require.True(t, stateBeforeWrite == nil)
	require.True(t, errBeforeWrite != nil)
	require.Equal(t, errBeforeWrite.Error(), fmt.Sprintf("session %s no found in storage", sessionID))

	require.True(t, stateAfterWrite != nil)
	require.True(t, errAfterWrite == nil)
	require.Equal(t, stateAfterWrite.MACAddress, "")
	require.Equal(t, stateAfterWrite.MSISDN, msisdn)

	require.True(t, stateAfterReset == nil)
	require.True(t, errAfterReset != nil)
	require.Equal(t, errAfterReset.Error(), fmt.Sprintf("session %s no found in storage", sessionID))
}

func TestMultipleConcurrentInsertDeleteGet(t *testing.T) {
	// Arrange
	degOfParallelism := 100
	reqPerConcurrentContext := 100
	onComplete := sync.WaitGroup{}
	storage := NewMultiSessionMemoryStorage()

	// Act
	for i := 0; i < degOfParallelism; i++ {
		go func(called string, calling string) {
			sessionID := fmt.Sprintf("session_%s_%s", calling, called)
			loopReadWriteDelete(t, storage, sessionID, reqPerConcurrentContext, &onComplete)
		}(fmt.Sprintf("called%d", i), fmt.Sprintf("calling%d", i))
	}
	onComplete.Wait()

	// Assert
	// nothing to do (assert will happen in the go routines spawned above)
}

func loopReadWriteDelete(
	t *testing.T,
	storage GlobalStorage,
	sessionID string,
	count int,
	onComplete *sync.WaitGroup,
) {
	for i := 1; i < count; i++ {
		performSignleReadWriteDeleteReadTest(t, storage, fmt.Sprintf("%s_%d", sessionID, i))
	}
	onComplete.Done()
}
