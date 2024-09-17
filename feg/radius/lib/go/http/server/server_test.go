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

package server

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"sort"
	"sync"
	"testing"

	"fbc/lib/go/http/header"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type mockCloser struct {
	mock.Mock
}

func (m *mockCloser) Close() error {
	args := m.Called()
	err := args.Error(0)
	return err
}

func TestServerRun(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:")
	require.NoError(t, err)

	core, o := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	closer := &mockCloser{}
	errBadCloser := errors.New("bad closer")
	closer.On("Close").Return(errBadCloser).Once()

	srv, err := New(
		Config{},
		Listener(listener),
		Logger(logger),
		Closer(closer),
	)
	require.NotNil(t, srv)
	require.NoError(t, err)
	srv.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}))

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan os.Signal)
	defer close(done)
	go func() {
		_ = srv.run(done)
		wg.Done()
	}()

	resp, err := http.Get("http://" + srv.Addr + "/")
	require.NoError(t, err)
	defer resp.Body.Close()

	done <- os.Interrupt
	wg.Wait()

	assert.Equal(t, 1, o.FilterMessageSnippet("starting http server").
		FilterField(zap.String("address", srv.Addr)).Len())
	assert.Equal(t, 1, o.FilterMessage("terminating http server").
		FilterField(zap.Error(http.ErrServerClosed)).Len())
	assert.Equal(t, 1, o.FilterField(zap.Error(errBadCloser)).Len())
	closer.AssertExpectations(t)
}

func TestServerLoggerConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  LoggerConfig
		wantErr bool
	}{
		{
			name: "discard",
			config: LoggerConfig{
				Type: DiscardLogger,
			},
		},
		{
			name: "development",
			config: LoggerConfig{
				Type: DevelopmentLogger,
			},
		},
		{
			name: "production",
			config: LoggerConfig{
				Type: ProductionLogger,
			},
		},
		{
			name: "production/verbose",
			config: LoggerConfig{
				Type:    ProductionLogger,
				Verbose: true,
			},
		},
		{
			name: "bad",
			config: LoggerConfig{
				Type: "badtype",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			srv, err := New(Config{Logger: &tt.config})
			assert.True(t, (err != nil) == tt.wantErr)
			assert.True(t, (srv == nil) == tt.wantErr)
		})
	}

	_, err := New(Config{Logger: &LoggerConfig{Type: "badtype"}})
	assert.Error(t, err)
}

func TestServerStaticFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	f, err := os.Create(path.Join(dir, "data.txt"))
	require.NoError(t, err)
	defer f.Close()

	text := "data text"
	_, err = f.WriteString(text)
	require.NoError(t, err)

	err = f.Sync()
	require.NoError(t, err)

	srv, err := New(Config{})
	require.NoError(t, err)
	srv.ServeFiles("/static/", http.Dir(dir))

	req := httptest.NewRequest(http.MethodGet, "/static/data.txt", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assert.Equal(t, text, rec.Body.String())
}

func TestServerRequestLogging(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	srv, err := New(Config{}, Logger(logger))
	require.NoError(t, err)
	srv.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	o = o.FilterMessage("HTTP request").
		FilterField(zap.String("method", req.Method)).
		FilterField(zap.Int("status", rec.Code))

	require.Equal(t, 1, o.Len())
	assert.Condition(t, func() bool {
		for _, field := range o.TakeAll()[0].Context {
			if field.Key == "url" {
				return true
			}
		}
		return false
	})
}

func TestServerOptionFunc(t *testing.T) {
	f := func(*Server) error { return nil }
	opt := OptionFunc(f)
	assert.Equal(t, reflect.ValueOf(f).Pointer(), reflect.ValueOf(opt.apply).Pointer())
	assert.Zero(t, opt.weight)
}

func TestServerOptionsOrder(t *testing.T) {
	option := func(weights []int, weight int) Option {
		return Option{
			apply: func(*Server) error {
				weights = append(weights, weight)
				return nil
			},
			weight: weight,
		}
	}

	options, weights := []Option{}, []int{}
	for i := 0; i < 10; i++ {
		options = append(options, option(weights, rand.Int()))
	}

	_, err := New(Config{}, options...)
	require.NoError(t, err)
	assert.True(t, sort.IntsAreSorted(weights))
}

func TestServerFailingOption(t *testing.T) {
	opt := Option{
		apply:  func(*Server) error { return errors.New("bad option") },
		weight: -1,
	}
	srv, err := New(Config{}, opt)
	assert.Nil(t, srv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "early")
	assert.Contains(t, err.Error(), "bad option")

	opt.weight = 1
	srv, err = New(Config{}, opt)
	assert.Nil(t, srv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "late")
	assert.Contains(t, err.Error(), "bad option")
}

func TestServerPanicRecovery(t *testing.T) {
	core, o := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	srv, err := New(Config{}, Logger(logger))
	require.NoError(t, err)
	errBadHandler := errors.New("bad handler func")
	srv.HandleFunc("/panic", func(http.ResponseWriter, *http.Request) {
		panic(errBadHandler)
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	o = o.FilterMessage("panic recovery").FilterField(zap.Error(errBadHandler))
	require.Equal(t, 1, o.Len())
	assert.Condition(t, func() bool {
		for _, field := range o.TakeAll()[0].Context {
			if field.Key == "stacktrace" {
				return true
			}
		}
		return false
	})
}

func TestServerRequestIDGeneration(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*http.Request)
		expect  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "without-id",
			prepare: func(req *http.Request) {},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.NotEmpty(t, rec.Header().Get(header.XRequestID))
			},
		},
		{
			name: "with-id",
			prepare: func(req *http.Request) {
				req.Header.Set(header.XRequestID, "f2314c55814a")
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, "f2314c55814a", rec.Header().Get(header.XRequestID))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			srv, err := New(Config{})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			tt.prepare(req)
			srv.ServeHTTP(rec, req)
			tt.expect(t, rec)
		})
	}
}

func TestServerProfilingHandler(t *testing.T) {
	srv, err := New(Config{})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/index", nil)
	handler, pattern := srv.Mux.Handler(req)
	assert.Equal(t, http.DefaultServeMux, handler)
	assert.Equal(t, "/debug/pprof/", pattern)
	_, pattern = http.DefaultServeMux.Handler(req)
	assert.NotEmpty(t, pattern)
}

func TestServerHealthHandler(t *testing.T) {
	srv, err := New(Config{})
	require.NoError(t, err)

	for _, target := range []string{"/health", "/healthz"} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}
}
