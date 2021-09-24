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

package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/proto"

	"github.com/magma/magma/src/go/log"
	"github.com/magma/magma/src/go/protos/magma/mconfig"
)

func TestLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		mconfigLevel mconfig.AgwD_LogLevel
		want         log.Level
	}{
		{
			mconfigLevel: mconfig.AgwD_DEBUG,
			want:         log.DebugLevel,
		},
		{
			mconfigLevel: mconfig.AgwD_INFO,
			want:         log.InfoLevel,
		},
		{
			mconfigLevel: mconfig.AgwD_WARN,
			want:         log.WarnLevel,
		},
		{
			mconfigLevel: mconfig.AgwD_ERROR,
			want:         log.ErrorLevel,
		},
		{
			mconfigLevel: mconfig.AgwD_UNSET,
			want:         log.InfoLevel,
		},
		{
			mconfigLevel: mconfig.AgwD_LogLevel(5),
			want:         log.InfoLevel,
		},
		{
			mconfigLevel: mconfig.AgwD_LogLevel(-1),
			want:         log.InfoLevel,
		},
	}

	for _, test := range tests {
		got := LogLevel(test.mconfigLevel)
		assert.Equal(
			t,
			test.want,
			got,
			"LogLevel(%s) = %s, want %s",
			test.mconfigLevel,
			got,
			test.want)
	}
}

func TestParseTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target string
		want   resolver.Target
	}{
		{},
		{
			target: "unix:///tmp/sctpd_downstream.sock",
			want: resolver.Target{
				Scheme:    "unix",
				Authority: "",
				Endpoint:  "/tmp/sctpd_downstream.sock",
			},
		},
		{
			target: "ipv4:localhost:1234",
			want: resolver.Target{
				Scheme:    "ipv4",
				Authority: "",
				Endpoint:  "localhost:1234",
			},
		},
		{
			target: "ipv6:[2607:f8b0:400e:c00::ef]:443",
			want: resolver.Target{
				Scheme:    "ipv6",
				Authority: "",
				Endpoint:  "[2607:f8b0:400e:c00::ef]:443",
			},
		},
		{
			target: "ipv6:[::]:1234",
			want: resolver.Target{
				Scheme:    "ipv6",
				Authority: "",
				Endpoint:  "[::]:1234",
			},
		},
	}

	for _, test := range tests {
		got := ParseTarget(test.target)
		assert.Equal(
			t,
			test.want,
			got,
			"ParseTarget(%s) = %v, want %v",
			test.target,
			got,
			test.want)
	}
}

func TestFilterCStyleComments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{},
		{
			in:   "/*",
			want: "",
		},
		{
			in:   "/**/",
			want: "",
		},
		{
			in:   "/*\nasdf\n*/",
			want: "",
		},
		{
			in:   "a/*b*/c",
			want: "ac",
		},
		{
			in:   "a /* b /* c /* d */ e",
			want: "a  e",
		},
		{
			in:   "a/*b*//*c*/d",
			want: "ad",
		},
	}
	for _, test := range tests {
		got := filterCStyleComments(test.in)
		assert.Equal(
			t,
			test.want,
			got,
			"filterCStyleComments(%s), got=%s want=%s",
			test.in,
			got,
			test.want)
	}
}

func TestNewConfigManager(t *testing.T) {
	t.Parallel()

	cm := NewConfigManager()
	assert.Equal(t, mconfig.AgwD_INFO, cm.Config().LogLevel)
	assert.Equal(t, "unix:///tmp/sctpd_downstream.sock", cm.Config().SctpdDownstreamServiceTarget)
	assert.Equal(t, "unix:///tmp/sctpd_upstream.sock", cm.Config().SctpdUpstreamServiceTarget)
	assert.Equal(t, "unix:///tmp/mme_sctpd_downstream.sock", cm.Config().MmeSctpdDownstreamServiceTarget)
	assert.Equal(t, "unix:///tmp/mme_sctpd_upstream.sock", cm.Config().MmeSctpdUpstreamServiceTarget)
}

func TestLoadConfigFile(t *testing.T) {
	t.Parallel()

	cm := NewConfigManager()
	err := LoadConfigFile(cm, filepath.Join("testdata", "accessd_mconfig.json"))
	assert.NoError(t, err)
	assert.Equal(t, mconfig.AgwD_WARN, cm.Config().LogLevel)
	assert.Equal(t, "a", cm.Config().SctpdDownstreamServiceTarget)
	assert.Equal(t, "b", cm.Config().SctpdUpstreamServiceTarget)
	assert.Equal(t, "c", cm.Config().MmeSctpdDownstreamServiceTarget)
	assert.Equal(t, "d", cm.Config().MmeSctpdUpstreamServiceTarget)
}

func TestNewConfigManager_DefaultNotFound(t *testing.T) {
	t.Parallel()

	cm := NewConfigManager()
	err := LoadConfigFile(cm, filepath.Join("testdata", "doesnotexist.json"))
	assert.True(t, os.IsNotExist(errors.Cause(err)))
	assert.True(t, proto.Equal(newDefaultConfig(), cm.Config()))
}

func TestNewConfigManager_OverrideSome(t *testing.T) {
	t.Parallel()

	cm := NewConfigManager()
	err := LoadConfigFile(cm, filepath.Join("testdata", "override_some.json"))
	assert.Nil(t, err)
	assert.Equal(t, mconfig.AgwD_DEBUG, cm.Config().LogLevel)
	assert.Equal(t, "unix:///tmp/sctpd_downstream.sock", cm.Config().SctpdDownstreamServiceTarget)
	assert.Equal(t, "foo", cm.Config().SctpdUpstreamServiceTarget)
	assert.Equal(t, "bar", cm.Config().MmeSctpdDownstreamServiceTarget)
	assert.Equal(t, "unix:///tmp/mme_sctpd_upstream.sock", cm.Config().MmeSctpdUpstreamServiceTarget)
}

func TestLoadConfigFile_StatErr(t *testing.T) {
	t.Parallel()

	path := "foo"
	statErr := errors.New("stat error")
	config, err := loadConfigFile(
		func(statPath string) (os.FileInfo, error) {
			assert.Equal(t, path, statPath)
			return nil, statErr
		},
		func(readPath string) ([]byte, error) {
			assert.Fail(
				t, "readFile shouldn't be called, path=%s", readPath)
			return nil, nil
		},
		func([]byte, proto.Message) error {
			return nil
		},
		path)
	assert.Nil(t, config)
	assert.EqualError(t, err, "path=foo: stat error")
}

func TestLoadConfigFile_ReadErr(t *testing.T) {
	t.Parallel()

	path := "foo"
	readError := errors.New("read error")
	config, err := loadConfigFile(
		func(statPath string) (os.FileInfo, error) {
			return nil, nil
		},
		func(readPath string) ([]byte, error) {
			assert.Equal(t, path, readPath)
			return nil, readError
		},
		func([]byte, proto.Message) error {
			return nil
		},
		path)
	assert.Nil(t, config)
	assert.EqualError(t, err, "path=foo: read error")
}

func TestNewConfigManager_BadSyntax(t *testing.T) {
	t.Parallel()

	path := filepath.Join("testdata", "bad_syntax.json")
	err := LoadConfigFile(nil, path)

	assert.True(
		t,
		strings.HasPrefix(
			err.Error(),
			"path="+path+" filtered={\"foo\": \"bar\"}"),
		"err=%s", err.Error())
}

func TestConfigManager_Merge(t *testing.T) {
	t.Parallel()

	cm := NewConfigManager()
	update := &mconfig.AgwD{MmeSctpdDownstreamServiceTarget: "a"}
	cm.Merge(update)
	assert.Equal(
		t,
		update.GetMmeSctpdDownstreamServiceTarget(),
		cm.Config().GetMmeSctpdDownstreamServiceTarget())
	update2 := &mconfig.AgwD{SctpdDownstreamServiceTarget: "b"}
	cm.Merge(update2)
	assert.Equal(
		t,
		update2.GetSctpdDownstreamServiceTarget(),
		cm.Config().GetSctpdDownstreamServiceTarget())

	assert.Equal(t, "b", cm.Config().SctpdDownstreamServiceTarget)
	assert.Equal(t, "unix:///tmp/sctpd_upstream.sock", cm.Config().SctpdUpstreamServiceTarget)
	assert.Equal(t, "a", cm.Config().MmeSctpdDownstreamServiceTarget)
	assert.Equal(t, "unix:///tmp/mme_sctpd_upstream.sock", cm.Config().MmeSctpdUpstreamServiceTarget)
}

func TestConfigManager_Race(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	wg.Add(2)

	cm := NewConfigManager()
	go func() {
		c := cm.Config()
		_ = c.GetLogLevel()
		wg.Done()
	}()
	go func() {
		cm.Merge(newDefaultConfig())
		wg.Done()
	}()

	wg.Wait()
}
