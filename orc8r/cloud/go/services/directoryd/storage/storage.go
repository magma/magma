/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/orc8r/lib/go/protos"
)

/*
	Persistence service interface for location records. All Directoryd data accesses from
	directoryd service must go through this interface.
*/
type DirectorydPersistenceService interface {

	// Get location record by ID
	GetRecord(tableId protos.TableID, recordId string) (*protos.LocationRecord, error)

	// Update existing location record or persist new location record
	UpdateOrCreateRecord(tableId protos.TableID, recordId string, record *protos.LocationRecord) error

	// Delete location record. MAY return non-nil error if recordId does not exist.
	DeleteRecord(tableId protos.TableID, recordId string) error
}
