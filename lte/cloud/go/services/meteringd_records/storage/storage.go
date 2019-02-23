/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/lte/cloud/go/protos"
)

/*
	Storage interface for metering records. All datastore writes from
	meteringd_records service must go through this interface.
*/
type MeteringRecordsStorage interface {

	// Do some initialization work on startup (e.g. create tables)
	InitTables() error

	// Update existing flow records with new usage counts or persist new
	// flow records
	UpdateOrCreateRecords(networkId string, flows []*protos.FlowRecord) error

	GetRecord(networkId string, recordId string) (*protos.FlowRecord, error)

	// Get all the flow records for a subscriber in a network.
	GetRecordsForSubscriber(networkId string, sid string) ([]*protos.FlowRecord, error)

	// Delete all flow records for a subscriber in a network
	DeleteRecordsForSubscriber(networkId string, sid string) error
}
