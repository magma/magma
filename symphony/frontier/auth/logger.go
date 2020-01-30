// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
)

// logger wraps zap.Logger to implement authboss.Logger.
type logger struct{ *zap.Logger }

// logs an info message.
func (l logger) Info(msg string) {
	l.Logger.Info(msg)
}

// logs an error message.
func (l logger) Error(msg string) {
	l.Logger.Error(msg)
}

// check if logger implements necessary interface.
var _ authboss.Logger = (*logger)(nil)

// ctxlogger wraps log.logger to implement authboss.ContextLogger.
type ctxlogger struct{ log.Logger }

// logs an info message.
func (l ctxlogger) Info(msg string) {
	l.Background().Info(msg)
}

// logs an error message.
func (l ctxlogger) Error(msg string) {
	l.Background().Error(msg)
}

// creates zap based logger from context.
func (l ctxlogger) FromContext(ctx context.Context) authboss.Logger {
	return logger{l.For(ctx)}
}

// check if ctxlogger implements necessary interfaces.
var (
	_ authboss.Logger        = (*ctxlogger)(nil)
	_ authboss.ContextLogger = (*ctxlogger)(nil)
)
