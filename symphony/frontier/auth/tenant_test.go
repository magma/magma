// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/facebookincubator/symphony/frontier/ent/enttest"
	"github.com/facebookincubator/symphony/frontier/ent/migrate"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testHandler struct {
	mock.Mock
}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Called(w, r)
}

func (t *testHandler) Load(ctx context.Context, name string) (*ent.Tenant, error) {
	args := t.Called(ctx, name)
	tenant, _ := args.Get(0).(*ent.Tenant)
	return tenant, args.Error(1)
}

func TestTenantHandler(t *testing.T) {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	client := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(sql.OpenDB(name, db))),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	defer client.Close()

	want, err := client.Tenant.
		Create().
		SetName("test").
		SetDomains([]string{}).
		SetNetworks([]string{}).
		Save(context.Background())
	require.NoError(t, err)

	t.Run("Simple", func(t *testing.T) {
		var m testHandler
		m.On("ServeHTTP", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				r := args.Get(1).(*http.Request)
				got := CurrentTenant(r.Context())
				require.NotNil(t, got)
				assert.Equal(t, want.ID, got.ID)
				assert.Equal(t, want.Name, got.Name)
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(http.StatusMultiStatus)
			}).
			Once()
		defer m.AssertExpectations(t)

		h := TenantHandler(&m, TenantClientLoader(client.Tenant, nil))
		req := httptest.NewRequest(http.MethodGet, "http://app.test.example.com", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusMultiStatus, rec.Code)
	})

	t.Run("NoSubdomain", func(t *testing.T) {
		var m testHandler
		defer m.AssertExpectations(t)
		h := TenantHandler(&m, TenantClientLoader(
			client.Tenant, logtest.NewTestLogger(t),
		))

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("LoadError", func(t *testing.T) {
		var m testHandler
		m.On("Load", mock.Anything, "test").
			Return(nil, &ent.NotFoundError{}).
			Once()
		m.On("Load", mock.Anything, "test").
			Return(nil, &ent.NotSingularError{}).
			Once()
		defer m.AssertExpectations(t)
		h := TenantHandler(&m, &m)

		for _, code := range []int{http.StatusNotFound, http.StatusInternalServerError} {
			req := httptest.NewRequest(http.MethodPost, "http://test.example.com", nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			assert.Equal(t, code, rec.Code)
		}
	})
}
