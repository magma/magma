// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/google/wire"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Config offers a declarative way to construct a logger.
	Config struct {
		Level  Level  `env:"LEVEL" long:"level" default:"info" description:"Only log messages with the given severity or above."`
		Format string `env:"FORMAT" long:"format" default:"console" choice:"console" choice:"json" description:"Output format of log messages."`
	}

	// Level attaches flags methods to zapcore.Level.
	Level zapcore.Level
)

// Set is a Wire provider set that produces a Logger.
var Set = wire.NewSet(New)

// New creates and sets a global logger from config.
func New(cfg Config) (Logger, func(), error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}
	bg := logger.Background()
	restoreGlobal := zap.ReplaceGlobals(bg)
	restoreStdLog := zap.RedirectStdLog(bg)
	return logger, func() {
		restoreStdLog()
		restoreGlobal()
	}, nil
}

// Build constructs a logger from Config.
func (cfg Config) Build() (Logger, error) {
	if cfg == (Config{}) {
		return NewNopLogger(), nil
	}

	var c zap.Config
	switch cfg.Format {
	case "console":
		c = zap.NewDevelopmentConfig()
	case "json":
		c = zap.NewProductionConfig()
	default:
		return nil, errors.Errorf("unsupported logging format: %q", cfg.Format)
	}
	c.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))

	logger, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return nil, errors.Wrap(err, "creating logger")
	}
	return NewDefaultLogger(logger), nil
}

// UnmarshalFlag implements flags.Unmarshaler.
func (l *Level) UnmarshalFlag(value string) error {
	var level zapcore.Level
	if err := level.Set(value); err != nil {
		return &flags.Error{
			Type:    flags.ErrMarshal,
			Message: err.Error(),
		}
	}
	*l = Level(level)
	return nil
}

// Ensure Level correctly implements flags.Unmarshaler.
var _ flags.Unmarshaler = (*Level)(nil)
