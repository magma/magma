// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enttest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	assert.NotNil(t, client)
	assert.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)
}
