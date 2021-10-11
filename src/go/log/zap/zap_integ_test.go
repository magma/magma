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

package zap_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	uber_zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/magma/magma/src/go/internal/testutil"
	"github.com/magma/magma/src/go/log"
	"github.com/magma/magma/src/go/log/zap"
)

type frozenZapClock struct {
	now time.Time
}

func (f *frozenZapClock) Now() time.Time {
	return f.now
}

func (f *frozenZapClock) NewTicker(duration time.Duration) *time.Ticker {
	panic("implement me")
}

func newLogFileURL(path string) string {
	if runtime.GOOS == "windows" {
		return "winfile:///" + path
	}
	return path
}

func newTestLogger(now time.Time, path string) (log.Logger, func()) {
	ec := uber_zap.NewDevelopmentEncoderConfig()
	enc := zapcore.NewConsoleEncoder(ec)
	ws, done, err := uber_zap.Open(path)
	if err != nil {
		panic(err)
	}
	l := zap.New(
		enc,
		ws,
		log.DebugLevel,
		uber_zap.AddCaller(),
		uber_zap.WithClock(&frozenZapClock{now: now}))
	return l, done
}

type logLine struct {
	level   string
	name    string
	line    int
	message string
	fields  string
}

func expectedLogOutput(t time.Time, lls []logLine) string {
	var lines []string
	ts := t.Format("2006-01-02T15:04:05.000Z0700")
	for _, ll := range lls {
		lineParts := []string{ts, ll.level}
		if ll.name != "" {
			lineParts = append(lineParts, ll.name)
		}
		lineParts = append(
			lineParts,
			fmt.Sprintf("zap/zap_integ_test.go:%d", ll.line),
			ll.message)
		if ll.fields != "" {
			lineParts = append(lineParts, ll.fields)
		}
		lines = append(lines, strings.Join(lineParts, "\t"))
	}
	return strings.Join(lines, "\n") + "\n"
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	td, tdDone := testutil.MustTempDir()
	defer tdDone()

	now := time.Now()
	t.Log(now)

	logPath := filepath.Join(td, "/logs")
	logURL := newLogFileURL(logPath)
	t.Log(logURL)

	l, logDone := newTestLogger(now, logURL)
	defer logDone()

	_, _, line, _ := runtime.Caller(0)
	l.Debug().Printf("%s world", log.DebugLevel)
	l.Info().Printf("%s world", log.InfoLevel)
	l.Warning().Printf("%s world", log.WarnLevel)
	l.Error().Printf("%s world", log.ErrorLevel)

	bs, err := ioutil.ReadFile(logPath)
	if err != nil {
		panic(err)
	}
	logs := string(bs)

	assert.Equal(
		t,
		expectedLogOutput(now,
			[]logLine{
				{level: "DEBUG", line: line + 1, message: "DEBUG world"},
				{level: "INFO", line: line + 2, message: "INFO world"},
				{level: "WARN", line: line + 3, message: "WARN world"},
				{level: "ERROR", line: line + 4, message: "ERROR world"},
			}),
		logs)
}

func TestLogger_Named(t *testing.T) {
	t.Parallel()

	td, tdDone := testutil.MustTempDir()
	defer tdDone()

	now := time.Now()
	t.Log(now)

	logPath := filepath.Join(td, "/logs")
	logURL := newLogFileURL(logPath)
	t.Log(logURL)

	l, logDone := newTestLogger(now, logURL)
	defer logDone()

	l.SetLevel(log.ErrorLevel)
	fooLog := l.Named("foo")
	fooLog.SetLevel(log.InfoLevel)
	foobarLog := fooLog.Named("bar")
	foobarLog.SetLevel(log.WarnLevel)

	assert.Equal(t, l.Level(), log.ErrorLevel)
	assert.Equal(t, fooLog.Level(), log.InfoLevel)
	assert.Equal(t, foobarLog.Level(), log.WarnLevel)

	_, _, line, _ := runtime.Caller(0)
	l.Warning().Print("should not print")
	l.Error().Print("should print")
	fooLog.Debug().Print("should not print")
	fooLog.Info().Print("should print")
	foobarLog.Info().Print("should not print")
	foobarLog.Warning().Print("should print")

	fooLog.SetLevel(log.DebugLevel)
	_, _, line2, _ := runtime.Caller(0)
	fooLog.Debug().Print("now should print")
	assert.Equal(t, fooLog.Level(), log.DebugLevel)

	bs, err := ioutil.ReadFile(logPath)
	if err != nil {
		panic(err)
	}
	logs := string(bs)

	assert.Equal(
		t,
		expectedLogOutput(now,
			[]logLine{
				{level: "ERROR", line: line + 2, message: "should print"},
				{level: "INFO", name: "foo", line: line + 4, message: "should print"},
				{level: "WARN", name: "foo.bar", line: line + 6, message: "should print"},
				{level: "DEBUG", name: "foo", line: line2 + 1, message: "now should print"},
			}),
		logs)
}

func TestLogger_With(t *testing.T) {
	t.Parallel()

	td, tdDone := testutil.MustTempDir()
	defer tdDone()

	now := time.Now()
	t.Log(now)

	logPath := filepath.Join(td, "/logs")
	logURL := newLogFileURL(logPath)
	t.Log(logURL)

	l, logDone := newTestLogger(now, logURL)
	defer logDone()

	l = l.Named("chat")
	john := l.With("name", "John Smith")
	jane := l.With("name", "Jane Doe")

	_, _, line, _ := runtime.Caller(0)
	john.Info().Print("how are you?")
	jane.Info().Print("doing great")

	bs, err := ioutil.ReadFile(logPath)
	if err != nil {
		panic(err)
	}
	logs := string(bs)

	assert.Equal(
		t,
		expectedLogOutput(now,
			[]logLine{
				{level: "INFO", name: "chat", line: line + 1, message: "how are you?", fields: "{\"name\": \"John Smith\"}"},
				{level: "INFO", name: "chat", line: line + 2, message: "doing great", fields: "{\"name\": \"Jane Doe\"}"},
			}),
		logs)
}
