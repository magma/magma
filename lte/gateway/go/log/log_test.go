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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level Level
		want  string
	}{
		{
			level: DebugLevel,
			want:  "DEBUG",
		},
		{
			level: InfoLevel,
			want:  "INFO",
		},
		{
			level: WarnLevel,
			want:  "WARN",
		},
		{
			level: ErrorLevel,
			want:  "ERROR",
		},
		{
			level: 100,
			want:  "INVALID (100)",
		},
	}

	for _, test := range tests {
		have := test.level.String()
		assert.Equal(t, test.want, have)
	}
}

func TestFullName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		names []string
		want  string
	}{
		"empty": {},
		"one name": {
			names: []string{"a"},
			want:  "a",
		},
		"multiple names": {
			names: []string{"a", "b", "c"},
			want:  "a.b.c",
		},
		"duplicate names": {
			names: []string{"a", "a", "a"},
			want:  "a.a.a",
		},
	}

	for desc, test := range tests {
		have := FullName(test.names)
		assert.Equal(t, test.want, have, desc)
	}
}

//go:generate go run github.com/golang/mock/mockgen -package log -destination mock_logger_test.go . Logger

func TestNewManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)
	m := NewManager(ml)

	assert.Len(t, m.loggers, 1)
	managed, ok := m.loggers[""]
	assert.True(t, ok, "root logger must be in log.Manager")
	assert.Same(t, m, managed.Manager)
	assert.Same(t, ml, managed.Logger)
}

func TestManagedLogger_Named(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)
	fooml := NewMockLogger(ctrl)
	ml.EXPECT().
		Named("foo").
		Return(fooml)
	fooml.EXPECT().
		Named("bar").
		Return(nil)

	m := NewManager(ml)

	root := m.LoggerFor("")
	foo := root.Named("foo")
	managedFoo, ok := foo.(*ManagedLogger)
	assert.True(t, ok, "logger in log.Manager must be *ManagedLogger")
	assert.Same(t, fooml, managedFoo.Logger)

	assert.Len(t, m.loggers, 2)
	assert.Nil(t, m.loggers[""].names)
	assert.Equal(t, []string{"foo"}, m.loggers["foo"].names)

	// ml.Named("foo") should not be called again
	fooAgain := root.Named("foo")
	assert.Same(t, foo, fooAgain)
	assert.Len(t, m.loggers, 2)

	foobar := foo.Named("bar")
	managedFooBar := m.LoggerFor("foo.bar")
	assert.Same(t, foobar, managedFooBar)
	assert.Len(t, m.loggers, 3)
	assert.Nil(t, m.loggers[""].names)
	assert.Equal(t, []string{"foo"}, m.loggers["foo"].names)
	assert.Equal(t, []string{"foo", "bar"}, m.loggers["foo.bar"].names)
}

func TestManagedLogger_Named_MissingRoot(t *testing.T) {
	m := &Manager{}
	ml := &ManagedLogger{Manager: m}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		assert.Equal(t, "missing root logger, see log.NewManager", r)
	}()
	ml.Named("foo")
}

func TestManagedLogger_Named_MissingIntermediate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)
	ml.EXPECT().
		Named("foo").
		Return(nil)

	m := &Manager{}
	rootLogger := &ManagedLogger{Manager: m, Logger: ml}
	foobarLogger := &ManagedLogger{Manager: m, names: []string{"foo", "bar"}}
	m.loggers = map[string]*ManagedLogger{
		"": rootLogger,
		"foo.bar": foobarLogger,
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		assert.Equal(t, "sublogger foo.bar exists, but foo didn't", r)
	}()
	m.LoggerFor("foo.bar.baz")
}
