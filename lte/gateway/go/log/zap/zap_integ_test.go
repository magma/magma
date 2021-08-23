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
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	uber_zap "go.uber.org/zap"

	"magma/lte/gateway/log"
	"magma/lte/gateway/log/zap"
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

func mustTempDir() string {
	td, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return td
}

func checkedRemoveAll(t *testing.T, path string) {
	if t.Failed() {
		return
	}
	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	td := mustTempDir()
	defer checkedRemoveAll(t, td)
	logPath := td + "/logs"
	t.Log(logPath)

	c := uber_zap.NewDevelopmentConfig()
	c.DisableStacktrace = true
	c.OutputPaths = []string{logPath}
	now := time.Now()
	l := zap.NewLogger(c, uber_zap.WithClock(&frozenZapClock{now: now}))

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

	ts := now.Format("2006-01-02T15:04:05.000Z0700")
	expectedLogLines := []string{
		fmt.Sprintf("%s\tDEBUG\tzap/zap_integ_test.go:%d\tDEBUG world", ts, line+1),
		fmt.Sprintf("%s\tINFO\tzap/zap_integ_test.go:%d\tINFO world", ts, line+2),
		fmt.Sprintf("%s\tWARN\tzap/zap_integ_test.go:%d\tWARN world", ts, line+3),
		fmt.Sprintf("%s\tERROR\tzap/zap_integ_test.go:%d\tERROR world", ts, line+4),
	}
	assert.Equal(
		t,
		strings.Join(expectedLogLines, "\n")+"\n",
		logs)
}

func TestLogger_Named(t *testing.T) {
	t.Parallel()

	td := mustTempDir()
	defer checkedRemoveAll(t, td)
	logPath := td + "/logs"
	t.Log(logPath)

	c := uber_zap.NewDevelopmentConfig()
	c.DisableStacktrace = true
	c.OutputPaths = []string{logPath}
	c.EncoderConfig.TimeKey = ""
	c.EncoderConfig.CallerKey = ""
	l := zap.NewLogger(c)

	l.SetLevel(log.ErrorLevel)
	fooLog := l.Named("foo")
	fooLog.SetLevel(log.InfoLevel)
	foobarLog := fooLog.Named("bar")
	foobarLog.SetLevel(log.WarnLevel)

	assert.Equal(t, l.Level(), log.ErrorLevel)
	assert.Equal(t, fooLog.Level(), log.InfoLevel)
	assert.Equal(t, foobarLog.Level(), log.WarnLevel)

	l.Warning().Print("should not print")
	l.Error().Print("should print")
	fooLog.Debug().Print("should not print")
	fooLog.Info().Print("should print")
	foobarLog.Info().Print("should not print")
	foobarLog.Warning().Print("should print")

	fooLog.SetLevel(log.DebugLevel)
	fooLog.Debug().Print("now should print")
	assert.Equal(t, fooLog.Level(), log.DebugLevel)

	bs, err := ioutil.ReadFile(logPath)
	if err != nil {
		panic(err)
	}
	logs := string(bs)

	expectedLogLines := []string{
		"ERROR\tshould print",
		"INFO\tfoo\tshould print",
		"WARN\tfoo.bar\tshould print",
		"DEBUG\tfoo\tnow should print",
	}
	assert.Equal(
		t,
		strings.Join(expectedLogLines, "\n")+"\n",
		logs)
}

func TestLogger_With(t *testing.T) {
	t.Parallel()

	td := mustTempDir()
	defer checkedRemoveAll(t, td)
	logPath := td + "/logs"
	t.Log(logPath)

	c := uber_zap.NewDevelopmentConfig()
	c.DisableStacktrace = true
	c.OutputPaths = []string{logPath}
	c.EncoderConfig.TimeKey = ""
	c.EncoderConfig.CallerKey = ""
	l := zap.NewLogger(c).Named("chat")
	john := l.With("name", "John Smith")
	jane := l.With("name", "Jane Doe")

	john.Info().Print("how are you?")
	jane.Info().Print("doing great")

	bs, err := ioutil.ReadFile(logPath)
	if err != nil {
		panic(err)
	}
	logs := string(bs)

	expectedLogLines := []string{
		"INFO\tchat\thow are you?\t{\"name\": \"John Smith\"}",
		"INFO\tchat\tdoing great\t{\"name\": \"Jane Doe\"}",
	}
	assert.Equal(
		t,
		strings.Join(expectedLogLines, "\n")+"\n",
		logs)
}
