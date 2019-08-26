/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"testing"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSubscriberDataStore(t *testing.T) {
	testSuite := new(SubscriberStoreTestSuite)
	testSuite.createStore = func() SubscriberStore {
		db, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
		assert.NoError(t, err)
		return NewSubscriberDataStore(db)
	}
	suite.Run(t, testSuite)
}
