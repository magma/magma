// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := logger{zap.New(core)}
	logger.Info("info message")
	logger.Error("error message")
	assert.NotEmpty(t, o.FilterMessage("info message").Len())
	assert.NotEmpty(t, o.FilterMessage("error message").Len())
}

type testLogger struct {
	mock.Mock
}

func (m *testLogger) Background() *zap.Logger {
	return m.Called().Get(0).(*zap.Logger)
}

func (m *testLogger) For(ctx context.Context) *zap.Logger {
	return m.Called(ctx).Get(0).(*zap.Logger)
}

func TestContextLogger(t *testing.T) {
	var (
		ctx    = context.Background()
		logger = zap.NewNop()
	)
	var m testLogger
	m.On("For", ctx).
		Return(logger).
		Once()
	m.On("Background").
		Return(logger).
		Twice()
	defer m.AssertExpectations(t)

	ctxlogger := ctxlogger{&m}
	ctxlogger.Info("info message")
	ctxlogger.Error("error message")
	_ = ctxlogger.FromContext(ctx)
}
