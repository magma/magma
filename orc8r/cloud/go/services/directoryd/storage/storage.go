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
	// GetHostname gets the hostname mapped to by hwid.
	GetHostname(hwid string) (string, error)

	// PutHostname maps hwid to hostname.
	PutHostname(hwid, hostname string) error
}
