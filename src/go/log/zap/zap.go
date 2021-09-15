// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package zap implements github.com/magma/magma/log via github.com/uber-go/zap
//
// Basic usage:
//
//	import (
//		uber_zap "go.uber.org/zap"
//
//		"github.com/magma/magma/log"
//		"github.com/magma/magma/log/zap"
//	)
//
//	lm := log.NewManager(zap.NewLogger())
//	lm.LoggerFor("thing").Info().Print("hello")
//	// Output: [thing] hello
package zap

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/magma/magma/log"
)

type printer struct {
	print  func(args ...interface{})
	printf func(format string, args ...interface{})
}

func (p *printer) Print(args ...interface{}) {
	p.print(args...)
}

func (p *printer) Printf(format string, args ...interface{}) {
	p.printf(format, args...)
}

func newPrinters(l *zap.Logger) map[log.Level]*printer {
	s := l.Sugar()
	return map[log.Level]*printer{
		log.DebugLevel: {
			print:  s.Debug,
			printf: s.Debugf,
		},
		log.InfoLevel: {
			print:  s.Info,
			printf: s.Infof,
		},
		log.WarnLevel: {
			print:  s.Warn,
			printf: s.Warnf,
		},
		log.ErrorLevel: {
			print:  s.Error,
			printf: s.Errorf,
		},
	}
}

// Logger adapts *zap.Logger to github.com/magma/magma/log.Logger.
type Logger struct {
	*zap.Logger

	names        []string
	level        zap.AtomicLevel
	newZapLogger func(zap.AtomicLevel) *zap.Logger
	printers     map[log.Level]*printer
}

var _ log.Logger = (*Logger)(nil)

// Error returns a log.ErrorLevel log.Printer backed by a zap.SugaredLogger.
func (l *Logger) Error() log.Printer {
	return l.printers[log.ErrorLevel]
}

// Warning returns a log.WarnLevel log.Printer backed by a zap.SugaredLogger.
func (l *Logger) Warning() log.Printer {
	return l.printers[log.WarnLevel]
}

// Info returns a log.InfoLevel log.Printer backed by a zap.SugaredLogger.
func (l *Logger) Info() log.Printer {
	return l.printers[log.InfoLevel]
}

// Debug returns a log.DebugLevel log.Printer backed by a zap.SugaredLogger.
func (l *Logger) Debug() log.Printer {
	return l.printers[log.DebugLevel]
}

// Level returns the currently configured log.Level.
func (l *Logger) Level() log.Level {
	lvl := l.level.Level()
	switch lvl {
	case zap.DebugLevel:
		return log.DebugLevel
	case zap.InfoLevel:
		return log.InfoLevel
	case zap.WarnLevel:
		return log.WarnLevel
	case zap.ErrorLevel:
		return log.ErrorLevel
	}
	panic("unsupported log level " + lvl.String())
}

func zapLevel(level log.Level) zapcore.Level {
	switch level {
	case log.DebugLevel:
		return zap.DebugLevel
	case log.InfoLevel:
		return zap.InfoLevel
	case log.WarnLevel:
		return zap.WarnLevel
	case log.ErrorLevel:
		return zap.ErrorLevel
	}
	panic("unsupported log level " + level.String())
}

// SetLevel updates the logger's log.Level. This takes effect immediately.
func (l *Logger) SetLevel(level log.Level) {
	l.level.SetLevel(zapLevel(level))
}

// Named returns a new logger scoped to the provided name.
func (l *Logger) Named(name string) log.Logger {
	names := append(l.names, name)

	// create a new level so that we can override log levels without affecting
	// existing logger.
	level := zap.NewAtomicLevelAt(l.level.Level())
	newLogger := l.newZapLogger(level)
	for _, name := range names {
		newLogger = newLogger.Named(name)
	}

	return &Logger{
		Logger:       newLogger,
		names:        names,
		level:        level,
		newZapLogger: l.newZapLogger,
		printers:     newPrinters(newLogger),
	}
}

// With returns a new logger with the field=value annotated on each message.
func (l *Logger) With(field string, value interface{}) log.Logger {
	newLogger := l.Logger.With(zap.Any(field, value))
	return &Logger{
		Logger:       newLogger,
		names:        l.names,
		level:        l.level,
		newZapLogger: l.newZapLogger,
		printers:     newPrinters(newLogger),
	}
}

// New returns a new *Logger, which satisfies magma/log.Logger. This function
// is available for fine-tuning logging config. For most usages, see NewLogger
// and NewLoggerAtLevel.
func New(
	enc zapcore.Encoder,
	ws zapcore.WriteSyncer,
	lvl log.Level,
	opts ...zap.Option,
) *Logger {
	if !lvl.Valid() {
		panic("invalid log.Level, lvl=" + lvl.String())
	}

	newZapLogger := func(al zap.AtomicLevel) *zap.Logger {
		nc := zapcore.NewCore(enc, ws, al)
		return zap.New(nc, append(opts, zap.AddCallerSkip(1))...)
	}
	level := zap.NewAtomicLevelAt(zapLevel(lvl))
	l := newZapLogger(level)
	return &Logger{
		Logger:       l,
		level:        level,
		newZapLogger: newZapLogger,
		printers:     newPrinters(l),
	}
}

// NewLoggerAtLevel creates a new *Logger at a specified log.Level with default
// console encoding outputting to the specified paths. See uber_go/zap.Open
// docs for more info on supported path formats.
func NewLoggerAtLevel(lvl log.Level, paths ...string) *Logger {
	if len(paths) == 0 {
		paths = []string{"stdout"}
	}
	ws, _, err := zap.Open(paths...)
	if err != nil {
		panic(errors.Wrapf(err, "paths=%s", paths))
	}

	return New(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		ws,
		lvl,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel))
}

// NewLogger returns a new root *Logger with default config, outputting to the
// specified paths. See uber_go/zap.Open docs for more info on supported path
// formats.
func NewLogger(paths ...string) *Logger {
	return NewLoggerAtLevel(log.InfoLevel, paths...)
}
