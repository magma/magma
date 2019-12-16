// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logtest

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestLogger is a logger to be used for testing.
type TestLogger struct {
	logger *zap.Logger
}

// TestingT is a subset of the API provided by all *testing.T and *testing.B objects.
type TestingT = zaptest.TestingT

// NewTestLogger creates a new testing logger.
func NewTestLogger(t TestingT) *TestLogger {
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller()))
	return &TestLogger{logger}
}

// Background returns a context-unaware logger.
func (l TestLogger) Background() *zap.Logger {
	return l.logger
}

// For ignores context and returns background logger.
func (l TestLogger) For(context.Context) *zap.Logger {
	return l.Background()
}
