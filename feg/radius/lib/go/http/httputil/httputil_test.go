/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package httputil

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloneRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	assert.NoError(t, err)

	clone := CloneRequest(req)
	assert.Equal(t, req, clone)
	assert.False(t, &req.Header == &clone.Header)
}

func TestCloneHeader(t *testing.T) {
	header := make(http.Header)
	header.Set("Content-Length", "123")
	header.Set("Content-Type", "text/plain")
	header.Set("Date", time.Now().Format(time.RFC3339))

	clone := CloneHeader(header)
	assert.False(t, &header == &clone)
	assert.Equal(t, header, clone)

	clone.Set("Content-Language", "en")
	assert.NotEqual(t, header, clone)
	assert.Len(t, clone, len(header)+1)
}
