// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler_test

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/store/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gocloud.dev/blob/driver"
	"gocloud.dev/blob/fileblob"
)

type mockSigner struct {
	mock.Mock
}

func (m *mockSigner) URLFromKey(ctx context.Context, key string, opts *driver.SignedURLOptions) (*url.URL, error) {
	args := m.Called(ctx, key, opts)
	u, _ := args.Get(0).(*url.URL)
	return u, args.Error(1)
}

func (m *mockSigner) KeyFromURL(ctx context.Context, u *url.URL) (string, error) {
	args := m.Called(ctx, u)
	return args.String(0), args.Error(1)
}

type handlerSuite struct {
	suite.Suite
	handler http.Handler
	signer  *mockSigner
}

func (s *handlerSuite) SetupTest() {
	s.signer = &mockSigner{}
	bucket, err := fileblob.OpenBucket(
		os.TempDir(),
		&fileblob.Options{
			URLSigner: s.signer,
		},
	)
	s.Require().NoError(err)
	s.handler = handler.New(handler.Config{
		Logger: logtest.NewTestLogger(s.T()),
		Bucket: bucket,
	})
	s.Require().NotNil(s.handler)
}

func (s *handlerSuite) TearDownTest() {
	s.signer.AssertExpectations(s.T())
}

func TestHandler(t *testing.T) {
	suite.Run(t, &handlerSuite{})
}

func (s *handlerSuite) TestGet() {
	const key = "test"
	s.signer.On("URLFromKey", mock.Anything, key, mock.AnythingOfType("*driver.SignedURLOptions")).
		Run(func(args mock.Arguments) {
			opts := args.Get(2).(*driver.SignedURLOptions)
			s.Assert().Equal(http.MethodGet, opts.Method)
		}).
		Return(&url.URL{}, nil).
		Once()
	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.URL.RawQuery = url.Values{"key": []string{key}}.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(http.StatusSeeOther, rec.Code)
}

func (s *handlerSuite) TestPut() {
	const host = "example.com"
	s.signer.On("URLFromKey", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*driver.SignedURLOptions")).
		Run(func(args mock.Arguments) {
			key := args.String(1)
			s.Assert().NotPanics(func() { uuid.Must(uuid.Parse(key)) })
			opts := args.Get(2).(*driver.SignedURLOptions)
			s.Assert().Equal(http.MethodPut, opts.Method)
		}).
		Return(&url.URL{Host: host}, nil).
		Once()
	req := httptest.NewRequest(http.MethodGet, "/put", nil)
	query := req.URL.Query()
	query.Set("contentType", "image/png")
	req.URL.RawQuery = query.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(http.StatusOK, rec.Code)
	s.Assert().Contains(rec.Header().Get("Content-Type"), "application/json")

	var result struct {
		URL string
		Key string
	}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	s.Require().NoError(err)

	u, err := url.Parse(result.URL)
	s.Require().NoError(err)
	s.Assert().Equal(host, u.Host)
	s.Assert().NotPanics(func() { uuid.Must(uuid.Parse(result.Key)) })
}

func (s *handlerSuite) TestSignError() {
	s.signer.On("URLFromKey", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("signing error")).
		Once()
	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (s *handlerSuite) TestWithoutKey() {
	for _, target := range []string{"/get", "/download"} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rec := httptest.NewRecorder()
		s.handler.ServeHTTP(rec, req)
		s.Assert().Equal(http.StatusNotFound, rec.Code)
	}
}

func (s *handlerSuite) TestGetWithTenant() {
	s.signer.On("URLFromKey", mock.Anything, "root/test", mock.Anything).
		Return(&url.URL{Host: "example.com"}, nil).
		Once()
	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.Header.Add("x-auth-organization", "root")
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(rec.Code, http.StatusSeeOther)
}

func (s *handlerSuite) TestDelete() {
	const key = "test"
	s.signer.On("URLFromKey", mock.Anything, key, mock.AnythingOfType("*driver.SignedURLOptions")).
		Run(func(args mock.Arguments) {
			opts := args.Get(2).(*driver.SignedURLOptions)
			s.Assert().Equal(http.MethodDelete, opts.Method)
		}).
		Return(&url.URL{Host: "example.com"}, nil).
		Once()
	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	req.URL.RawQuery = url.Values{"key": []string{key}}.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(http.StatusTemporaryRedirect, rec.Code)
}

func (s *handlerSuite) TestDeleteBadMethod() {
	req := httptest.NewRequest(http.MethodGet, "/delete", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(http.StatusMethodNotAllowed, rec.Code)
}

func (s *handlerSuite) TestDownload() {
	s.T().SkipNow()
	const key = "test"
	s.signer.On("URLFromKey", mock.Anything, key, mock.AnythingOfType("*driver.SignedURLOptions")).
		Run(func(args mock.Arguments) {
			opts := args.Get(2).(*driver.SignedURLOptions)
			s.Assert().Equal(http.MethodGet, opts.Method)
		}).
		Return(&url.URL{Host: "example.com"}, nil).
		Once()
	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	const filename = "file"
	req.URL.RawQuery = url.Values{"key": []string{"test"}, "fileName": []string{filename}}.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(http.StatusSeeOther, rec.Code)
	var ref struct {
		URL string `xml:"href,attr"`
	}
	err := xml.NewDecoder(rec.Body).Decode(&ref)
	s.Require().NoError(err)
	u, err := url.Parse(ref.URL)
	s.Require().NoError(err)
	s.Assert().Equal("attachment; filename="+filename, u.Query().Get("response-content-disposition"))
}

func (s *handlerSuite) TestDownloadNoFilename() {
	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	req.URL.RawQuery = url.Values{"key": []string{"test"}}.Encode()
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	s.Assert().Equal(http.StatusNotFound, rec.Code)
}
