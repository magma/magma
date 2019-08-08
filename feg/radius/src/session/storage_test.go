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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicSessionStorageInsertGet(t *testing.T) {
	// Arrange
	sessionID := "sessionid"
	storage := NewSessionStorage(NewMultiSessionMemoryStorage(), sessionID)

	// Act and Assert
	performSignleSessionStorageTest(t, storage, sessionID)
}

func performSignleSessionStorageTest(t *testing.T, storage Storage, sessionID string) {
	// Arrange
	msisdn := fmt.Sprintf("+%d", rand.Intn(999999))

	// Act
	stateBeforeWrite, errBeforeWrite := storage.Get()
	writeErr := storage.Set(State{
		MSISDN: msisdn,
	})
	stateAfterWrite, errAfterWrite := storage.Get()
	storage.Reset()
	stateAfterReset, errAfterReset := storage.Get()

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
