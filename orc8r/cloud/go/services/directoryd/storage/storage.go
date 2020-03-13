/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

// DirectorydStorage is the persistence service interface for location records.
// All Directoryd data accesses from directoryd service must go through this interface.
type DirectorydStorage interface {
	// GetHostnameForHWID returns the hostname mapped to by hardware ID.
	GetHostnameForHWID(hwid string) (string, error)

	// MapHWIDsToHostnames maps {hwid -> hostname}.
	MapHWIDsToHostnames(hwidToHostname map[string]string) error

	// GetIMSIForSessionID returns the IMSI mapped to by session ID.
	GetIMSIForSessionID(networkID, sessionID string) (string, error)

	// MapSessionIDsToIMSIs maps {session ID -> IMSI}.
	MapSessionIDsToIMSIs(networkID string, sessionIDToIMSI map[string]string) error
}
