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
	"bytes"
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/magma/magma/log"
)

func TestPrinter_Print(t *testing.T) {
	t.Parallel()

	var calls int
	p := &printer{
		print: func(args ...interface{}) {
			calls++
			assert.Equal(t, args, []interface{}{"a", "b", "c"})
		},
	}
	p.Print("a", "b", "c")
	assert.Equal(t, 1, calls)
}

func TestPrinter_Printf(t *testing.T) {
	t.Parallel()

	var calls int
	p := &printer{
		printf: func(format string, args ...interface{}) {
			calls++
			assert.Equal(t, format, "hello %s")
			assert.Equal(t, args, []interface{}{"world"})
		},
	}
	p.Printf("hello %s", "world")
	assert.Equal(t, 1, calls)
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	c := zap.NewProductionConfig()
	l := NewLogger(c)
	logger, ok := l.(*Logger)
	assert.True(
		t, ok, "NewLogger must return *zap.Logger, have=%+v", l)
	assert.NotNil(t, logger.Config)
	assert.NotNil(t, logger.Logger)
	assert.NotNil(t, logger.printers)
	for level := log.DebugLevel; level <= log.ErrorLevel; level++ {
		assert.NotNil(t, logger.printers[level])
	}

	assert.Same(t, l.Error(), logger.printers[log.ErrorLevel])
	assert.Same(t, l.Warning(), logger.printers[log.WarnLevel])
	assert.Same(t, l.Info(), logger.printers[log.InfoLevel])
	assert.Same(t, l.Debug(), logger.printers[log.DebugLevel])
}

func TestNewLogger_Error(t *testing.T) {
	t.Parallel()

	c := zap.Config{Level: zap.NewAtomicLevelAt(0)}

	// we expect error from zap.Config.Build() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("panic value not err, r=%+v", r)
			return
		}
		assert.Regexp(
			t, "^zap.Config=.*no encoder name specified$", err.Error())
	}()
	_ = NewLogger(c)
}

func TestLogger_Level(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level zapcore.Level
		want  log.Level
	}{
		{
			level: zapcore.DebugLevel,
			want:  log.DebugLevel,
		},
		{
			level: zapcore.InfoLevel,
			want:  log.InfoLevel,
		},
		{
			level: zapcore.WarnLevel,
			want:  log.WarnLevel,
		},
		{
			level: zapcore.ErrorLevel,
			want:  log.ErrorLevel,
		},
	}

	for _, test := range tests {
		c := zap.NewProductionConfig()
		c.Level.SetLevel(test.level)
		l := NewLogger(c)
		assert.Equal(t, test.want, l.Level())
	}
}

func TestLogger_Level_Error(t *testing.T) {
	t.Parallel()

	// we expect error from zap.Config.Build() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		assert.Equal(t, r, "unsupported log level fatal")
	}()
	c := zap.NewProductionConfig()
	c.Level.SetLevel(zapcore.FatalLevel)
	l := NewLogger(c)
	l.Level()
}

func TestLogger_SetLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level log.Level
		want  zapcore.Level
	}{
		{
			level: log.DebugLevel,
			want:  zapcore.DebugLevel,
		},
		{
			level: log.InfoLevel,
			want:  zapcore.InfoLevel,
		},
		{
			level: log.WarnLevel,
			want:  zapcore.WarnLevel,
		},
		{
			level: log.ErrorLevel,
			want:  zapcore.ErrorLevel,
		},
	}

	for _, test := range tests {
		l := NewLogger(zap.NewProductionConfig())
		l.SetLevel(test.level)
		logger, ok := l.(*Logger)
		assert.True(
			t, ok, "NewLogger must return *zap.Logger, have=%+v", l)
		assert.Equal(t, test.want, logger.Config.Level.Level())
	}
}

func TestLogger_SetLevel_Error(t *testing.T) {
	t.Parallel()

	// we expect error from zap.Config.Build() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		assert.Equal(t, r, "unsupported log level INVALID (5)")
	}()
	l := NewLogger(zap.NewProductionConfig())
	l.SetLevel(5)
}

// memorySink implements zap.Sink by writing all messages to a buffer.
type memorySink struct {
	*bytes.Buffer
}

func (s *memorySink) Close() error { return nil }
func (s *memorySink) Sync() error  { return nil }

var memorySinks map[string]*memorySink
var memorySinksM sync.Mutex

func getMemorySink(name string) *memorySink {
	memorySinksM.Lock()
	defer memorySinksM.Unlock()
	return memorySinks[name]
}

func newMemorySink(name string) *memorySink {
	memorySinksM.Lock()
	defer memorySinksM.Unlock()

	if s, ok := memorySinks[name]; ok {
		return s
	}
	s := &memorySink{&bytes.Buffer{}}
	memorySinks[name] = s
	return s
}

func init() {
	memorySinks = make(map[string]*memorySink)
	zap.RegisterSink("memory", func(u *url.URL) (zap.Sink, error) {
		return newMemorySink(u.Host), nil
	})
}

func TestLoggerNamed(t *testing.T) {
	t.Parallel()

	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"memory://TestLoggerNamed"}

	l := NewLogger(c)
	logger, ok := l.(*Logger)
	assert.True(
		t, ok, "NewLogger must return *zap.Logger, have=%+v", l)
	foo := l.Named("foo")
	fooLogger, ok := foo.(*Logger)
	assert.True(
		t, ok, "NewLogger must return *zap.Logger, have=%+v", l)

	assert.NotSame(t, logger, fooLogger)
	assert.Empty(t, logger.names)
	assert.Equal(t, fooLogger.names, []string{"foo"})

	l.Info().Print("a")
	foo.Info().Print("b")

	out := getMemorySink("TestLoggerNamed")
	assert.NotNil(t, out)

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	assert.Equal(t, 2, len(lines))

	var first, second map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &first))
	assert.NoError(t, json.Unmarshal([]byte(lines[1]), &second))

	assert.Nil(t, first["logger"])
	assert.Equal(t, "a", first["msg"])
	assert.Equal(t, "foo", second["logger"])
	assert.Equal(t, "b", second["msg"])
}

func TestLoggerNamed_Error(t *testing.T) {
	t.Parallel()

	l := NewLogger(zap.NewProductionConfig())
	logger, ok := l.(*Logger)
	assert.True(
		t, ok, "NewLogger must return *zap.Logger, have=%+v", l)
	logger.Config.Encoding = ""

	// we expect error from zap.Config.Build() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("panic value not err, r=%+v", r)
			return
		}
		assert.Regexp(
			t, "(?s)^zap.Config=.*no encoder name specified", err.Error())
	}()
	_ = l.Named("bar")
}

func TestLogger_With(t *testing.T) {
	t.Parallel()

	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"memory://TestLoggerWith"}
	l := NewLogger(c)
	test := l.With("env", "test")
	nested := test.With("foo", "bar")
	overwrite := nested.With("env", "test2")

	l.Info().Print("hi")
	test.Info().Print("hi")
	nested.Info().Print("hi")
	overwrite.Info().Print("hi")

	out := getMemorySink("TestLoggerWith")
	assert.NotNil(t, out)

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	assert.Equal(t, 4, len(lines))

	var first, second, third, fourth map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &first))
	assert.NoError(t, json.Unmarshal([]byte(lines[1]), &second))
	assert.NoError(t, json.Unmarshal([]byte(lines[2]), &third))
	assert.NoError(t, json.Unmarshal([]byte(lines[3]), &fourth))

	assert.Nil(t, first["env"])
	assert.Nil(t, first["foo"])
	assert.Equal(t, "hi", first["msg"])

	assert.Equal(t, "test", second["env"])
	assert.Nil(t, second["foo"])
	assert.Equal(t, "hi", second["msg"])

	assert.Equal(t, "test", third["env"])
	assert.Equal(t, "bar", third["foo"])
	assert.Equal(t, "hi", third["msg"])

	assert.Equal(t, "test2", fourth["env"])
	assert.Equal(t, "bar", fourth["foo"])
	assert.Equal(t, "hi", fourth["msg"])
}
