/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test

import (
	"testing"

	"magma/feg/gateway/services/testcore/hss/servicers"

	"github.com/stretchr/testify/assert"
)

func TestAllZero(t *testing.T) {
	assert.Equal(t, true, servicers.AllZero(nil))
	assert.Equal(t, true, servicers.AllZero(make([]byte, 50)))

	bytes := make([]byte, 30)
	bytes[25] = 1
	assert.Equal(t, false, servicers.AllZero(bytes))
}
