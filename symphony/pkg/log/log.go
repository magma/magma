// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A Logger can create *zap.Logger instances either for a given context or context-less.
type Logger interface {
	Background() *zap.Logger
	For(context.Context) *zap.Logger
}

// DefaultLogger implements Logger interface.
type DefaultLogger struct {
	bg *zap.Logger
}

// NewDefaultLogger creates a new default logger.
func NewDefaultLogger(logger *zap.Logger) *DefaultLogger {
	return &DefaultLogger{logger}
}

// Background returns a context-unaware logger.
func (l DefaultLogger) Background() *zap.Logger {
	return l.bg
}

// For returns a context-aware logger.
func (l DefaultLogger) For(ctx context.Context) *zap.Logger {
	logger := l.Background()
	if span := trace.FromContext(ctx); span != nil {
		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, spanCore{LevelEnabler: core, span: span})
		}))
	}
	return logger.With(FieldsFromContext(ctx)...)
}

// NopLogger is a no-op Logger implementation.
type NopLogger struct {
	logger *zap.Logger
}

// NewNopLogger creates a new nop logger.
func NewNopLogger() *NopLogger {
	return &NopLogger{zap.NewNop()}
}

// Background belongs to Logger interface.
func (l NopLogger) Background() *zap.Logger {
	return l.logger
}

// For belongs to Logger interface.
func (l NopLogger) For(context.Context) *zap.Logger {
	return l.logger
}
