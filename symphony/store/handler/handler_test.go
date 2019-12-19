// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/store/sign/signtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	signer := &signtest.MockSigner{}
	signer.On("Sign", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", nil).
		Once()
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusSeeOther, rec.Code)
}

func TestPut(t *testing.T) {
	signer := &signtest.MockSigner{}
	signer.On("Sign", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("example.com", nil).
		Once()
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	req := httptest.NewRequest(http.MethodGet, "/put", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var result struct {
		URL string
		Key string
	}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, result.URL, "example.com")
	assert.NotPanics(t, func() { uuid.Must(uuid.Parse(result.Key)) })
}

func TestSignError(t *testing.T) {
	signer := &signtest.MockSigner{}
	signer.On("Sign", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("sign error")).
		Twice()
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	for _, target := range []string{"/get", "/put"} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestWithoutKey(t *testing.T) {
	signer := &signtest.MockSigner{}
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	for _, target := range []string{"/get", "/download"} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestGetWithTenant(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get?"+url.Values{"key": []string{"test"}}.Encode(), nil)
	req.Header.Add("x-auth-organization", "root")
	rec := httptest.NewRecorder()

	signer := &signtest.MockSigner{}
	signer.On("Sign", mock.Anything, mock.Anything, "root/test", mock.Anything, mock.Anything).
		Return("", nil).
		Once()
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)
	h.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusSeeOther)
}

func TestDelete(t *testing.T) {
	signer := &signtest.MockSigner{}
	signer.On("Sign", mock.Anything, mock.Anything, "test", mock.Anything, mock.Anything).
		Return("example.com", nil).
		Once()
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestDeleteBadMethod(t *testing.T) {
	signer := &signtest.MockSigner{}
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	req := httptest.NewRequest(http.MethodGet, "/delete", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestDownload(t *testing.T) {
	signer := &signtest.MockSigner{}
	signer.On("Sign", mock.Anything, mock.Anything, "test", "file", mock.Anything).
		Return("example.com", nil).
		Once()
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}, "fileName": []string{"file"}}.Encode()
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusSeeOther, rec.Code)
}

func TestDownloadNoFilename(t *testing.T) {
	signer := &signtest.MockSigner{}
	defer signer.AssertExpectations(t)

	h := New(Config{Logger: logtest.NewTestLogger(t), Signer: signer})
	require.NotNil(t, h)

	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
