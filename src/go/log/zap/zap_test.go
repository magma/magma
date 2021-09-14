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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

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

	l := NewLogger()

	assert.Equal(t, zapcore.InfoLevel, l.level.Level())
	assert.NotNil(t, l.Logger)
	assert.Nil(t, l.names)
	assert.NotNil(t, l.newZapLogger)
	assert.NotNil(t, l.printers)
	for level := log.DebugLevel; level <= log.ErrorLevel; level++ {
		assert.NotNil(t, l.printers[level])
	}

	// compiler ensures github.com/magma/magma/log/zap.Logger is a log.Logger
	var logger log.Logger = l

	assert.Same(t, logger.Error(), l.printers[log.ErrorLevel])
	assert.Same(t, logger.Warning(), l.printers[log.WarnLevel])
	assert.Same(t, logger.Info(), l.printers[log.InfoLevel])
	assert.Same(t, logger.Debug(), l.printers[log.DebugLevel])
}

func TestNewLogger_Error(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_ = NewLogger("///badpath")
	})
}

func TestNewLoggerAtLevel_Error(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(
		t,
		"invalid log.Level, lvl=INVALID (100)",
		func() {
			_ = NewLoggerAtLevel(log.Level(100))
		})
}

func TestLogger_Level(t *testing.T) {
	t.Parallel()

	tests := []struct {
		zapLevel zapcore.Level
		logLevel log.Level
	}{
		{
			zapLevel: zapcore.DebugLevel,
			logLevel: log.DebugLevel,
		},
		{
			zapLevel: zapcore.InfoLevel,
			logLevel: log.InfoLevel,
		},
		{
			zapLevel: zapcore.WarnLevel,
			logLevel: log.WarnLevel,
		},
		{
			zapLevel: zapcore.ErrorLevel,
			logLevel: log.ErrorLevel,
		},
	}

	for _, test := range tests {
		l := NewLoggerAtLevel(test.logLevel)
		assert.Equal(t, test.zapLevel, l.level.Level())
		var logger log.Logger = l
		assert.Equal(t, test.logLevel, logger.Level())
	}
}

func TestLogger_Level_Error(t *testing.T) {
	t.Parallel()

	l := NewLogger()
	l.level.SetLevel(zapcore.FatalLevel)
	assert.PanicsWithValue(t, "unsupported log level fatal", func() {
		l.Level()
	})
}

func TestLogger_SetLevel_Error(t *testing.T) {
	t.Parallel()

	l := NewLogger()
	assert.PanicsWithValue(t, "unsupported log level INVALID (5)", func() {
		l.SetLevel(5)
	})
}

func TestLoggerNamed(t *testing.T) {
	t.Parallel()

	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	ws := &zaptest.Buffer{}
	logger := New(enc, ws, log.InfoLevel)
	var l log.Logger = logger

	foo := l.Named("foo")
	fooLogger, ok := foo.(*Logger)
	assert.True(
		t, ok, "NewLogger must return *zap.Logger, got=%+v", l)

	assert.NotSame(t, logger, fooLogger)
	assert.Empty(t, logger.names)
	assert.Equal(t, fooLogger.names, []string{"foo"})

	l.Info().Print("a")
	foo.Info().Print("b")

	lines := ws.Lines()
	assert.Equal(t, 2, len(lines))

	var first, second map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &first))
	assert.NoError(t, json.Unmarshal([]byte(lines[1]), &second))

	assert.Nil(t, first["logger"])
	assert.Equal(t, "a", first["msg"])
	assert.Equal(t, "foo", second["logger"])
	assert.Equal(t, "b", second["msg"])
}

func TestLogger_With(t *testing.T) {
	t.Parallel()

	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	ws := &zaptest.Buffer{}
	logger := New(enc, ws, log.InfoLevel)
	var l log.Logger = logger

	test := l.With("env", "test")
	nested := test.With("foo", "bar")
	overwrite := nested.With("env", "test2")

	l.Info().Print("hi")
	test.Info().Print("hi")
	nested.Info().Print("hi")
	overwrite.Info().Print("hi")

	lines := ws.Lines()
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
