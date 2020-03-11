/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package dispatcher_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/dispatcher"
	"magma/orc8r/cloud/go/services/dispatcher/test_init"

	"github.com/stretchr/testify/assert"
)

func TestGetHostnameForHwid(t *testing.T) {
	test_init.StartTestService(t)

	// Values seeded during dispatcher test service init
	hostname, err := dispatcher.GetHostnameForHwid("some_hwid_0")
	assert.NoError(t, err)
	assert.Equal(t, "some_hostname_0", hostname)
}
