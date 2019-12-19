// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spa

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/log/logtest"

	"github.com/justinas/nosurf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/constantvar"
)

func waitForVar(t *testing.T, v *runtimevar.Variable) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err := v.Watch(ctx)
	require.NotEqual(t, context.DeadlineExceeded, err)
}

func TestSinglePageHandler(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, distPath+"/manifest.json", r.URL.Path)
		_, err := io.WriteString(w, `{
			"vendor.js":     "/inventory/static/dist/vendor.ec3cd6c53dea9177230f.js",
			"login.js":      "/inventory/static/dist/login.33ee93e09d062b3f7794.js",
			"main.js":       "/inventory/static/dist/main.b96c37ebdd5682a41d25.js",
			"master.js":     "/inventory/static/dist/master.b784f9c100062efe30c3.js",
			"onboarding.js": "/inventory/static/dist/onboarding.7d0af91b2281cddc7f25.js"
		}`)
		assert.NoError(t, err)
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	sp, err := SinglePageHandler(
		"main", OriginManifester(u),
	)
	require.NoError(t, err)
	defer sp.Close()
	waitForVar(t, sp.manifest)

	h := nosurf.New(sp)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "<html")
	assert.Contains(t, body, "/main.b96c37ebdd5682a41d25.js")
	assert.Contains(t, body, "/vendor.ec3cd6c53dea9177230f.js")
	assert.Contains(t, body, nosurf.Token(req))
}

func TestSinglePageOriginNoManifest(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	sp, err := SinglePageHandler(
		"main", OriginManifester(u),
		WithLogger(logtest.NewTestLogger(t)),
	)
	require.NoError(t, err)
	defer sp.Close()
	waitForVar(t, sp.manifest)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	sp.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "<html")
	assert.Contains(t, body, "/main.js")
	assert.Contains(t, body, "/vendor.js")
}

func TestSinglePageOriginError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	sp, err := SinglePageHandler(
		"main", OriginManifester(u),
		WithLogger(logtest.NewTestLogger(t)),
	)
	require.NoError(t, err)
	defer sp.Close()
	waitForVar(t, sp.manifest)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	sp.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestSinglePageBadManifest(t *testing.T) {
	tests := []struct {
		name     string
		manifest interface{}
		errmsg   string
	}{
		{
			name:     "InvalidType",
			manifest: 42,
			errmsg:   "incompatible manifest type int",
		},
		{
			name:     "NoMainJS",
			manifest: map[string]string{},
			errmsg:   `missing manifest key "main.js"`,
		},
		{
			name:     "NoMainJS",
			manifest: map[string]string{"main.js": "/main.js"},
			errmsg:   `missing manifest key "vendor.js"`,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			core, o := observer.New(zap.InfoLevel)
			logger := log.NewDefaultLogger(zap.New(core))
			sp, err := SinglePageHandler(
				"main.js", func() (*runtimevar.Variable, error) {
					cv := constantvar.New(tc.manifest)
					waitForVar(t, cv)
					return cv, nil
				}, WithLogger(logger),
			)
			require.NoError(t, err)
			defer sp.Close()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			sp.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
			entries := o.FilterMessage("getting index data").TakeAll()
			require.Len(t, entries, 1)
			assert.Equal(t, tc.errmsg, entries[0].ContextMap()["error"])
		})
	}
}

func TestSinglePageBadManifester(t *testing.T) {
	_, err := SinglePageHandler("root", func() (*runtimevar.Variable, error) {
		return nil, errors.New("bad manifester")
	})
	assert.EqualError(t, err, "resolving manifest var: bad manifester")
}

func TestSinglePageCheckHealth(t *testing.T) {
	tests := []struct {
		manifest  *runtimevar.Variable
		assertion func(assert.TestingT, error, ...interface{}) bool
	}{
		{
			manifest: constantvar.New(`{
				"main.js": "/main.js",
				"vendor.js": "/vendor.js"
			}`),
			assertion: assert.NoError,
		},
		{
			manifest:  constantvar.NewError(errors.New("bad manifest")),
			assertion: assert.Error,
		},
	}
	for _, tc := range tests {
		sp, err := SinglePageHandler("health", func() (*runtimevar.Variable, error) {
			waitForVar(t, tc.manifest)
			return tc.manifest, nil
		})
		require.NoError(t, err)
		err = sp.CheckHealth()
		tc.assertion(t, err)
	}
}
