// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/google/wire"
	"go.uber.org/zap"
)

// Provider is a wire provider of this package.
var Provider = wire.NewSet(
	ProvideLogger,
	ProvideZapLogger,
)

// Provider is a wire provider that produces a logger from config.
func ProvideLogger(config Config) (Logger, func(), error) {
	logger, err := New(config)
	if err != nil {
		return nil, nil, err
	}
	restoreGlobal := zap.ReplaceGlobals(logger.Background())
	restoreStdLog := zap.RedirectStdLog(logger.Background())
	return logger, func() { restoreStdLog(); restoreGlobal() }, nil
}

// ProvideZapLogger is a wire provider that produces zap logger from logger.
func ProvideZapLogger(logger Logger) *zap.Logger {
	return logger.Background()
}
