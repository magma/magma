// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AllowedLevel is a settable identifier for the
// minimum level a log entry must have.
type AllowedLevel zapcore.Level

// String returns a ASCII representation of the allowed log level.
func (l AllowedLevel) String() string {
	return (zapcore.Level)(l).String()
}

// Set updates the value of the allowed level.
func (l *AllowedLevel) Set(s string) error {
	var level zapcore.Level
	if err := level.Set(s); err != nil {
		return err
	}
	if level > zapcore.ErrorLevel {
		return fmt.Errorf("unrecognized level: %q", s)
	}
	*l = AllowedLevel(level)
	return nil
}

// UnmarshalFlag updates the value of the allowed level.
func (l *AllowedLevel) UnmarshalFlag(s string) error {
	return l.Set(s)
}

// AllowedFormat is a settable identifier for the output
// format that the logger can have.
type AllowedFormat string

// String returns a ASCII representation of the allowed log format.
func (f AllowedFormat) String() string {
	return string(f)
}

// Set updates the value of the allowed format.
func (f *AllowedFormat) Set(s string) error {
	switch s {
	case "console", "json":
		*f = AllowedFormat(s)
	default:
		return fmt.Errorf("unrecognized format: %q", s)
	}
	return nil
}

// UnmarshalFlag updates the value of the allowed format.
func (f *AllowedFormat) UnmarshalFlag(s string) error {
	return f.Set(s)
}

// Config is a struct containing configurable settings for the logger.
type Config struct {
	Level  AllowedLevel  `env:"LEVEL" long:"level" default:"info" description:"Only log messages with the given severity or above."`
	Format AllowedFormat `env:"FORMAT" long:"format" default:"console" choice:"console" choice:"json" description:"Output format of log messages."`
}

// empty returns true when the config is equal to its zero value.
func (c Config) empty() bool {
	return c == Config{}
}

// New returns a new leveled contextual logger.
func New(config Config) (Logger, error) {
	if config.empty() {
		return NewNopLogger(), nil
	}
	var cfg zap.Config
	switch config.Format {
	case "json":
		cfg = zap.NewProductionConfig()
	case "console":
		cfg = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("unrecognized format: %q", config.Format)
	}
	cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(config.Level))
	logger, err := cfg.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return nil, fmt.Errorf("building logger: %w", err)
	}
	return NewDefaultLogger(logger), nil
}

// MustNew returns a new leveled contextual logger, and panic on error.
func MustNew(config Config) Logger {
	logger, err := New(config)
	if err != nil {
		panic(err)
	}
	return logger
}
