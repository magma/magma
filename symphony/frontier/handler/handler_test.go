// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/justinas/nosurf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyHandler(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, err := io.Copy(w, r.Body)
		assert.NoError(t, err)
	}))
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	h := NewHandler(Config{
		ProxyTarget: u,
		Logger:      logtest.NewTestLogger(t),
	})

	req := httptest.NewRequest(
		http.MethodPost, "/",
		strings.NewReader("proxy body"),
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "proxy body", rec.Body.String())
}

func TestNoSurfing(t *testing.T) {
	h := NewHandler(Config{
		ProxyTarget: &url.URL{},
		Logger:      logtest.NewTestLogger(t),
	})
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, nosurf.Token(r))
		assert.NoError(t, err)
	}).
		Methods(http.MethodGet)
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}).
		Methods(http.MethodPost)

	srv := httptest.NewTLSServer(h)
	defer srv.Close()
	client := srv.Client()
	client.Jar, _ = cookiejar.New(nil)

	rsp, err := client.Get(srv.URL)
	require.NoError(t, err)
	token, err := ioutil.ReadAll(rsp.Body)
	require.NoError(t, err)
	rsp.Body.Close()

	t.Run("WithToken", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, srv.URL, nil)
		require.NoError(t, err)
		req.Header.Set("X-CSRF-Token", string(token))
		req.Header.Set("Content-Type", "text/plain")
		rsp, err := client.Do(req)
		require.NoError(t, err)
		defer rsp.Body.Close()
		assert.Equal(t, http.StatusAccepted, rsp.StatusCode)
	})
	t.Run("WithoutToken", func(t *testing.T) {
		rsp, err := client.Post(srv.URL, "text/plain", nil)
		require.NoError(t, err)
		defer rsp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})
}
