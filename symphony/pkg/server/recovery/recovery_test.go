// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testRecoverer struct {
	mock.Mock
}

func (tr *testRecoverer) HandlePanic(ctx context.Context, p interface{}) error {
	args := tr.Called(ctx, p)
	return args.Error(0)
}

func TestRecoveryHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := errors.New("bad handler")

	var recoverer testRecoverer
	recoverer.On("HandlePanic", ctx, err).
		Return(err).
		Once()
	defer recoverer.AssertExpectations(t)

	handler := &Handler{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(err)
		}),
		HandlerFunc: recoverer.HandlePanic,
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.EqualError(t, err, strings.TrimSuffix(rec.Body.String(), "\n"))
}
