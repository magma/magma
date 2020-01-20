// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orc8r

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixedHostTransport(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	client := srv.Client()
	client.Transport = Transport{
		Base: client.Transport,
		Host: u.Host,
	}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	rsp, err := client.Do(req)
	require.NoError(t, err)
	defer rsp.Body.Close()
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
}
