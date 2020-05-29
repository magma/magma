/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package reindex

import (
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"

	"github.com/pkg/errors"
)

const (
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

	indexer0  = mocks.NewMockIndexer(id0, version0, nil, nil, nil, nil)
	indexer0a = mocks.NewMockIndexer(id0, version0a, nil, nil, nil, nil)
	indexer1  = mocks.NewMockIndexer(id1, version1, nil, nil, nil, nil)
	indexer1a = mocks.NewMockIndexer(id1, version1a, nil, nil, nil, nil)
	indexer2  = mocks.NewMockIndexer(id2, version2, nil, nil, nil, nil)
)
