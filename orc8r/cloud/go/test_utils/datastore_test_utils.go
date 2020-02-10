/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

/*
	Assert the datastore has some expected rows
*/
func AssertDatastoreHasRows(
	t *testing.T,
	store datastore.Api,
	tableKey string,
	expectedRows map[string]interface{},
	deserializer func([]byte) (interface{}, error),
) {
	marshaledValueWrappers, err := store.GetMany(tableKey, getMapKeys(expectedRows))
	assert.NoError(t, err)
	assert.Equal(
		t, len(expectedRows), len(marshaledValueWrappers),
		"Expected %d rows in datastore, actual %d", len(expectedRows), len(marshaledValueWrappers),
	)

	for k, v := range marshaledValueWrappers {
		unmarshaledVal, err := deserializer(v.Value)
		assert.NoError(t, err)

		expectedVal, ok := expectedRows[k]
		assert.True(t, ok)
		valMsg, ok := unmarshaledVal.(proto.Message)
		if ok {
			expectedMsg, ok := expectedVal.(proto.Message)
			assert.True(t, ok)
			assert.Equal(t, protos.TestMarshal(expectedMsg), protos.TestMarshal(valMsg))
		} else {
			assert.Equal(t, expectedVal, unmarshaledVal)
		}

	}
}

func getMapKeys(in map[string]interface{}) []string {
	ret := make([]string, 0, len(in))
	for k := range in {
		ret = append(ret, k)
	}
	return ret
}

/*
  Assert that the datastore does not have an entry for a specific key.
*/
func AssertDatastoreDoesNotHaveRow(
	t *testing.T,
	store datastore.Api,
	tableKey string,
	rowKey string,
) {
	allKeys, err := store.ListKeys(tableKey)
	assert.NoError(t, err)
	for _, k := range allKeys {
		if k == rowKey {
			assert.Fail(
				t,
				fmt.Sprintf("Found table key %s which is not supposed to exist", rowKey))
		}
	}
}

/*
  Insert some test fixture data into the datastore
*/
func SetupTestFixtures(
	t *testing.T,
	store datastore.Api,
	tableKey string,
	fixtures map[string]interface{},
	serializer func(interface{}) ([]byte, error),
) {
	for k, val := range fixtures {
		marshaledVal, err := serializer(val)
		assert.NoError(t, err)
		err = store.Put(tableKey, k, marshaledVal)
		assert.NoError(t, err)
	}
}
