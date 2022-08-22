/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

func TestConfigParse(t *testing.T) {
	tests := []struct {
		args    []string
		env     map[string]string
		prepare func()
		expect  func(*testing.T, *Config, error)
	}{
		{
			args: []string{
				"--no-stats", "true",
				"--view", "proc",
				"--view", "http",
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert.NoError(t, err)
				assert.True(t, cfg.DisableStats)
				assert.Equal(t, StatViews{"proc", "http"}, cfg.StatViews)
				assert.Nil(t, cfg.Jaeger)
			},
		},
		{
			args: []string{
				"--jaeger", `{"Process":{"ServiceName":"test"}}`,
			},
			env: map[string]string{
				"VIEWS": "http,proc",
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert.NoError(t, err)
				assert.Equal(t, StatViews{"http", "proc"}, cfg.StatViews)
				require.NotNil(t, cfg.Jaeger)
				assert.Equal(t, "test", cfg.Jaeger.Process.ServiceName)
			},
		},
		{
			args: []string{
				"--view", "http",
				"--view", "none",
			},
			expect: func(t *testing.T, _ *Config, err error) {
				assert.Error(t, err)
			},
		},
		{
			args: []string{
				"--view", "custom",
				"--view", "proc",
			},
			prepare: func() {
				MustRegisterViewer("custom", Views{})
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert.NoError(t, err)
				assert.Equal(t, StatViews{"custom", "proc"}, cfg.StatViews)
			},
		},
		{
			env: map[string]string{
				"JAEGER": `{"AgentEndpoint":"localhost:6831"}`,
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert.NoError(t, err)
				require.NotNil(t, cfg.Jaeger)
				assert.Equal(t, "localhost:6831", cfg.Jaeger.AgentEndpoint)
			},
		},
		{
			args: []string{
				"--jaeger", `{"AgentEndpoint:"localhost:6831"}`,
			},
			expect: func(t *testing.T, _ *Config, err error) {
				assert.Error(t, err)
			},
		},
	}

	var opts flags.Options = flags.HelpFlag | flags.PassDoubleDash
	for _, tc := range tests {
		for k, v := range tc.env {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}
		if tc.prepare != nil {
			tc.prepare()
		}
		var cfg Config
		_, err := flags.NewParser(&cfg, opts).ParseArgs(tc.args)
		tc.expect(t, &cfg, err)
		for k := range tc.env {
			err := os.Unsetenv(k)
			require.NoError(t, err)
		}
	}
}

func TestConfigBuildStats(t *testing.T) {
	census, err := Config{
		StatViews: StatViews{"proc", "http", "http"},
	}.Build(
		WithNamespace("test"),
		WithLogger(nil),
	)
	require.NoError(t, err)
	require.NotNil(t, census)
	assert.NoError(t, census.Close())
}

func TestConfigBuildDisableStats(t *testing.T) {
	census, err := Config{DisableStats: true}.Build()
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	census.StatsHandler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestConfigBuildBadView(t *testing.T) {
	_, err := Config{StatViews: StatViews{"http", "none"}}.Build()
	assert.Error(t, err)
}

func TestConfigBuildJaeger(t *testing.T) {
	census, err := Config{
		DisableStats:        true,
		SamplingProbability: 1,
		Jaeger: &Jaeger{
			Options: jaeger.Options{
				CollectorEndpoint: "http://jaeger:14268/api/traces",
			},
		},
	}.Build(WithService("tester"))
	assert.NoError(t, err)
	_, span := trace.StartSpan(context.Background(), "test")
	assert.True(t, span.SpanContext().IsSampled())
	span.End()
	census.Close()
	_, span = trace.StartSpan(context.Background(), "test")
	assert.False(t, span.SpanContext().IsSampled())
	span.End()
}
