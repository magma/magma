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

// Package log provides a generic logging abstraction for Magma.
//
// Goals:
//
//   1. Allow packages to log against an abstraction so that they can be
//      deployed in different environments and logging implementations can be
//      injected. e.g. `magma/imsi` could be used in AGW as well as cloud
//      components. On an AGW machine, code could log out to console via a
//      logging library like `uber/zap` or it could use `log/syslog` to stream
//      directly to a syslog-ng server. When `magma/imsi` is used in a cloud
//      service, an adapter to a cloud native logging solution, e.g. Google
//      Cloud Logging or Amazon CLoudWatch Logs, could be injected instead.
//
//   2. Provide a scoped/named log tree, e.g. "module.operation.function" tags
//      on log output. This allows shared code to be re-used in various
//      contexts with logging output that tracks with its calling scope.
//
//   3. Explicit key/value parameterized logging, e.g. "IMSI=1234". This is
//      important so parameters can be machine-parsed easily and reliably. A
//      common use case is to elide sensitive information before logs are
//      shared.
//
// Examples:
//
//	logger := zap.NewLogger()
//	logger.Debug().Printf("hello %s", "world")
//
// With management:
//
//	type Component struct {
//		log.Logger
//	}
//
//	func (c *Component) DoFoo(val int) {
//		c.Named("DoFoo").
//			With("val", val).
//			Info().
//			Print("doing stuff")
//	}
//
//	rootLogger := zap.NewLogger()
//	logManager := log.NewManager(rootLogger)
//	c := &Component{Logger: logManager.LoggerFor("component")}
//	c.DoFoo(42)
//
//	// Output: [component.DoFoo] [val=42] doing stuff
//
// FAQs
//
// Q. Why isn't `Fatal` a supported level?
// A. Implementations of `Fatal` often call `os.Exit(1)`. This immediately ends
//    the process and does not allow any panic recovery / crash handling to
//    occur. We strongly discourage `os.Exit` anywhere in the codebase; handle
//    errors when possible and `panic` if not possible.
//
// Q. Why `Debug` instead of verbosity levels (e.g. `V(level Level)` in glog)?
// A. Using a scoped/named log tree is a more precise way to control logging
//    and combined with a `Debug` level should satisfy most use cases. Also,
//    those unfamiliar with glog may be confused by `V` level logging usage.
//
// Q. What about the built-in `log` package or using `glog` directly?
// A. Many log packages lack features (see Goals above) and do not have
//    mockable interfaces. While log output can be redirected to achieve
//    similar results, our code then becomes tightly coupled/dependent on the
//    specific logging library.
package log

import (
	"fmt"
	"strings"
	"sync"
)

// A Level is a logging priority. Higher levels are more important.
type Level int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running
	// smoothly, it shouldn't generate any error-level logs.
	ErrorLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

// String returns a upper-case ASCII representation of the log level.
func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return fmt.Sprintf("INVALID (%d)", l)
}

// A Printer abstracts outputting a log message of any Level.
type Printer interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
}

// A Logger provides access to output logs at priority Levels, enables
// accessing Named subloggers, and attaching key=value parameters to log
// output scope.
type Logger interface {
	// Error provides a Printer for Error Level logging.
	Error() Printer
	// Warning provides a Printer for Warning Level logging.
	Warning() Printer
	// Info provides a Printer for Info Level logging.
	Info() Printer
	// Debug provides a Printer for Debug Level logging.
	Debug() Printer

	// Level returns the currently active output Level for this Logger.
	Level() Level
	// SetLevel sets the active output Level for this Logger. It is meant to
	// be changeable at runtime.
	SetLevel(level Level)

	// Named returns a new Logger, scoped with the provided name. Nested calls
	// to Named create a log tree.
	Named(name string) Logger
	// With returns a new Logger with a field=value parameter added to the
	// scope.
	With(field string, value interface{}) Logger
}

// FullName is a convenience function to create a fully qualified Logger name
// from a slice of names.
func FullName(names []string) string {
	return strings.Join(names, ".")
}

// A Manager manages a tree of Loggers. The root Logger will be available
// via loggers[""].
type Manager struct {
	loggers map[string]*managedLogger
	sync.Mutex
}

// NewManager creates a new Manager with the provided Logger as the root.
func NewManager(l Logger) *Manager {
	m := &Manager{
		loggers: make(map[string]*managedLogger),
	}

	newLogger := &managedLogger{Manager: m, Logger: l}
	m.loggers[""] = newLogger
	return m
}

// LoggerFor returns a logger for the fully specified named scope, creating
// one if necessary. This function is threadsafe.
func (m *Manager) LoggerFor(fullName string) Logger {
	m.Lock()
	defer m.Unlock()
	if l, ok := m.loggers[fullName]; ok {
		return l
	}

	prev, ok := m.loggers[""]
	if !ok {
		panic("missing root logger, see log.NewManager")
	}
	names := strings.Split(fullName, ".")
	for i := 0; i < len(names); i++ {
		fullName := FullName(names[0 : i+1])
		l, ok := m.loggers[fullName]
		if ok {
			prev = l
			continue
		}

		for j := i; j < len(names); j++ {
			newFullName := FullName(names[0 : j+1])
			if _, ok := m.loggers[newFullName]; ok {
				panic(
					fmt.Sprintf(
						"sublogger %s exists, but %s didn't",
						newFullName,
						FullName(names[0:j])))
			}

			newNames := append(prev.names, names[j])
			newLogger := &managedLogger{
				names:   newNames,
				Manager: m,
				Logger:  prev.Logger.Named(names[j]),
			}
			m.loggers[newFullName] = newLogger
			prev = newLogger
		}
		break
	}
	return prev
}

type managedLogger struct {
	names []string

	*Manager
	Logger
}

// Named overrides the Logger's Named function to integrate with log.Manager.
// This ensures only one Logger per fully qualified name is created.
func (ml *managedLogger) Named(name string) Logger {
	return ml.LoggerFor(FullName(append(ml.names, name)))
}
