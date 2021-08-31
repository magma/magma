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

package log

import (
	"fmt"
	"strings"
	"sync"
)

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
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return fmt.Sprintf("INVALID (%d)", l)
}

type Printer interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
}

type Logger interface {
	Error() Printer
	Warning() Printer
	Info() Printer
	Debug() Printer

	Level() Level
	SetLevel(level Level)

	Named(name string) Logger
	With(field string, value interface{}) Logger
}

func FullName(names []string) string {
	return strings.Join(names, ".")
}

type Manager struct {
	loggers map[string]*ManagedLogger
	sync.Mutex
}

func NewManager(l Logger) *Manager {
	m := &Manager{
		loggers: make(map[string]*ManagedLogger),
	}

	newLogger := &ManagedLogger{Manager: m, Logger: l}
	m.loggers[""] = newLogger
	return m
}

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
				panic(fmt.Sprintf("sublogger %s exists, but %s didn't", newFullName, FullName(names[0:j])))
			}

			newNames := append(prev.names, names[j])
			newLogger := &ManagedLogger{
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

type ManagedLogger struct {
	names []string

	*Manager
	Logger
}

func (ml *ManagedLogger) Named(name string) Logger {
	return ml.LoggerFor(FullName(append(ml.names, name)))
}
