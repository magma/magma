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

package zap

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

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

type Logger struct {
	names []string
	*zap.Logger
	zap.Config

	printers map[log.Level]*printer
}

var _ log.Logger = (*Logger)(nil)

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

func NewLogger(c zap.Config, opts ...zap.Option) log.Logger {
	l, err := c.Build(append(opts, zap.AddCallerSkip(1))...)
	if err != nil {
		panic(errors.Wrapf(err, "zap.Config=%+v", c))
	}
	return &Logger{
		Config:   c,
		Logger:   l,
		printers: newPrinters(l),
	}
}

func (l *Logger) Error() log.Printer {
	return l.printers[log.ErrorLevel]
}

func (l *Logger) Warning() log.Printer {
	return l.printers[log.WarnLevel]
}

func (l *Logger) Info() log.Printer {
	return l.printers[log.InfoLevel]
}

func (l *Logger) Debug() log.Printer {
	return l.printers[log.DebugLevel]
}

func (l *Logger) Level() log.Level {
	lvl := l.Config.Level.Level()
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

func (l *Logger) SetLevel(level log.Level) {
	switch level {
	case log.DebugLevel:
		l.Config.Level.SetLevel(zap.DebugLevel)
	case log.InfoLevel:
		l.Config.Level.SetLevel(zap.InfoLevel)
	case log.WarnLevel:
		l.Config.Level.SetLevel(zap.WarnLevel)
	case log.ErrorLevel:
		l.Config.Level.SetLevel(zap.ErrorLevel)
	default:
		panic("unsupported log level " + level.String())
	}
}

func (l *Logger) Named(name string) log.Logger {
	names := append(l.names, name)

	// create a new config so we can override log levels without affecting
	// existing logger.
	newConfig := l.Config
	newConfig.Level = zap.NewAtomicLevelAt(l.Config.Level.Level())

	newLogger, err := newConfig.Build()
	if err != nil {
		panic(errors.Wrapf(err, "zap.Config=%+v", newConfig))
	}

	for _, name := range names {
		newLogger = newLogger.Named(name)
	}

	return &Logger{
		names:    names,
		Config:   newConfig,
		Logger:   newLogger,
		printers: newPrinters(newLogger),
	}
}

func (l *Logger) With(field string, value interface{}) log.Logger {
	newLogger := l.Logger.With(zap.Any(field, value))
	return &Logger{
		Config:   l.Config,
		Logger:   newLogger,
		printers: newPrinters(newLogger),
	}
}
