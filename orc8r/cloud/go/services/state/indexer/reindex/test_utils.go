/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package reindex

import (
	"errors"

	"magma/orc8r/cloud/go/services/state/indexer"
)

const (
	connectionStringPostgres = "dbname=magma_test user=magma_test password=magma_test host=postgres_test sslmode=disable"

	maxAttempts = 2

	nid0 = "some_networkid_0"
	nid1 = "some_networkid_1"
	nid2 = "some_networkid_2"

	hwid0 = "some_hwid_0"
	hwid1 = "some_hwid_1"
	hwid2 = "some_hwid_2"

	id0 = "some_indexerid_0"
	id1 = "some_indexerid_1"
	id2 = "some_indexerid_2"
	id3 = "some_indexerid_3"

	zero      indexer.Version = 0
	version0  indexer.Version = 100
	version0a indexer.Version = 1000
	version0b indexer.Version = 10000
	version1  indexer.Version = 200
	version1a indexer.Version = 2000
	version2  indexer.Version = 300
	version3  indexer.Version = 400
)

var (
	someErr  = errors.New("some_error")
	someErr1 = errors.New("some_error_1")
	someErr2 = errors.New("some_error_2")
	someErr3 = errors.New("some_error_3")

	indexer0  = indexer.NewTestIndexer(id0, version0)
	indexer0a = indexer.NewTestIndexer(id0, version0a)
	indexer1  = indexer.NewTestIndexer(id1, version1)
	indexer1a = indexer.NewTestIndexer(id1, version1a)
	indexer2  = indexer.NewTestIndexer(id2, version2)
)
