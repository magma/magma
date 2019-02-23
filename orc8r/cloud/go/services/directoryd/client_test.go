/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package directoryd_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"

	"github.com/stretchr/testify/assert"
)

const (
	testGwId1  = "gw1"
	testGwId2  = "gw2"
	testSubId1 = "sub1"
	testSubId2 = "sub2"
	testSubId3 = "sub3"
)

func TestDirectorydControllerClientMethods(t *testing.T) {
	directoryd_test_init.StartTestService(t)

	// Get empty DB
	_, err := directoryd.GetHardwareIdByIMSI(testSubId1)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Repeat using other table

	// Get empty DB
	_, err = directoryd.GetHostNameByIMSI(testSubId1)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Add two locations
	err = directoryd.UpdateHostNameByHwId(testSubId1, testGwId1)
	assert.NoError(t, err)

	err = directoryd.UpdateHostNameByHwId(testSubId2, testGwId2)
	assert.NoError(t, err)

	// Read back
	record, err := directoryd.GetHostNameByIMSI(testSubId1)
	assert.NoError(t, err)
	assert.Equal(t, testGwId1, record)

	record, err = directoryd.GetHostNameByIMSI(testSubId2)
	assert.NoError(t, err)
	assert.Equal(t, testGwId2, record)

	record, err = directoryd.GetHostNameByIMSI(testSubId3)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Delete
	err = directoryd.DeleteHostNameByIMSI(testSubId1)
	assert.NoError(t, err)

	record, err = directoryd.GetHostNameByIMSI(testSubId1)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Delete unknown
	err = directoryd.DeleteHostNameByIMSI(testSubId3)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error finding location record: No record for query")
}
