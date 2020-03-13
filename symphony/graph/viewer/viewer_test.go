// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/schema"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/testdb"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func newTestClient(t *testing.T) *ent.Client {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	drv := sql.OpenDB(name, db)
	client := ent.NewClient(ent.Driver(drv))
	require.NoError(t, client.Schema.Create(context.Background(), schema.WithGlobalUniqueID(true)))
	return client
}

func TestViewerHandler(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*http.Request)
		expect  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "TestTenant",
			prepare: func(req *http.Request) {
				req.Header.Set(TenantHeader, "test")
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "test", rec.Body.String())
			},
		},
		{
			name: "NoTenant",
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.NotZero(t, rec.Body.Len())
			},
		},
		{
			name: "ReadOnlyUser",
			prepare: func(req *http.Request) {
				req.Header.Set(TenantHeader, "test")
				req.Header.Set(RoleHeader, "readonly")
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.NotZero(t, rec.Body.Len())
			},
		},
	}

	h := TenancyHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			viewer := FromContext(ctx)
			require.NotNil(t, viewer)
			assert.NotNil(t, log.FieldsFromContext(ctx))
			_, _ = io.WriteString(w, viewer.Tenant)
			if r.Header.Get(RoleHeader) == "readonly" {
				_, err := ent.FromContext(ctx).Tx(ctx)
				require.EqualError(t, err, "ent: starting a transaction: permission denied: read-only user")
			}
		}),
		NewFixedTenancy(&ent.Client{}),
	)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.prepare != nil {
				tc.prepare(req)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			tc.expect(t, rec)
		})
	}
}

func TestWebSocketUpgradeHandler(t *testing.T) {
	var upgrader websocket.Upgrader
	srv := httptest.NewServer(
		WebSocketUpgradeHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if websocket.IsWebSocketUpgrade(r) {
					conn, err := upgrader.Upgrade(w, r, nil)
					require.NoError(t, err)
					defer conn.Close()
					for {
						if _, _, err := conn.ReadMessage(); err != nil {
							return
						}
					}
				}
				w.WriteHeader(http.StatusOK)
			}),
			"",
		),
	)
	defer srv.Close()

	t.Run("NoUpgrade", func(t *testing.T) {
		rsp, err := srv.Client().Get(srv.URL)
		require.NoError(t, err)
		defer rsp.Body.Close()
		assert.Equal(t, http.StatusOK, rsp.StatusCode)
	})
	t.Run("AuthenticatedUpgrade", func(t *testing.T) {
		header := http.Header{}
		header.Set(TenantHeader, "test")
		u, _ := url.Parse(srv.URL)
		u.Scheme = "ws"
		conn, rsp, err := websocket.DefaultDialer.DialContext(
			context.Background(), u.String(), header,
		)
		require.NoError(t, err)
		rsp.Body.Close()
		conn.Close()
	})

	const host = "test.example.com"
	authenticator := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, host, r.Header.Get("X-Forwarded-Host"))
			if username, _, ok := r.BasicAuth(); ok {
				err := json.NewEncoder(w).Encode(&Viewer{
					Tenant: "test",
					User:   username,
					Role:   "user",
				})
				require.NoError(t, err)
			} else if cookie, err := r.Cookie("Viewer"); err == nil {
				data, err := base64.StdEncoding.DecodeString(cookie.Value)
				require.NoError(t, err)
				_, err = w.Write(data)
				require.NoError(t, err)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		}),
	)
	defer authenticator.Close()

	t.Run("Auth", func(t *testing.T) {
		authenticate := func(t *testing.T, authReq func(*http.Request)) {
			var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(TenantHeader, r.Header.Get(TenantHeader))
				w.Header().Set(UserHeader, r.Header.Get(UserHeader))
				w.Header().Set(RoleHeader, r.Header.Get(RoleHeader))
				w.WriteHeader(http.StatusOK)
			})
			handler = WebSocketUpgradeHandler(handler, authenticator.URL)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = host
			authReq(req)
			req.Header.Set("Connection", "upgrade")
			req.Header.Set("Upgrade", "websocket")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "test", rec.Header().Get(TenantHeader))
			assert.Equal(t, "tester", rec.Header().Get(UserHeader))
			assert.Equal(t, "user", rec.Header().Get(RoleHeader))
		}

		t.Run("Basic", func(t *testing.T) {
			authenticate(t, func(req *http.Request) {
				req.SetBasicAuth("tester", "tester")
			})
		})
		t.Run("Session", func(t *testing.T) {
			data, err := json.Marshal(&Viewer{
				Tenant: "test",
				User:   "tester",
				Role:   "user",
			})
			require.NoError(t, err)
			authenticate(t, func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "Viewer",
					Value: base64.StdEncoding.EncodeToString(data),
				})
			})
		})
	})
}

func TestUserHandler(t *testing.T) {
	client := newTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	u, err := client.User.Create().SetAuthID("user1").Save(ctx)
	require.NoError(t, err)

	t.Run("WithUserEntExists", func(t *testing.T) {
		h := UserHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count, err := client.User.Query().Count(ctx)
				require.NoError(t, err)
				require.Equal(t, 1, count)
				w.WriteHeader(http.StatusAccepted)
			}),
			log.NewNopLogger(),
		)
		v := &Viewer{
			Tenant: "test",
			User:   "user1",
			Role:   "",
		}
		ctx = NewContext(ctx, v)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithUserEntNotExist", func(t *testing.T) {
		h := UserHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				exist, err := client.User.Query().Where(user.AuthID("user2")).Exist(ctx)
				require.NoError(t, err)
				require.True(t, exist)
				count, err := client.User.Query().Count(ctx)
				require.NoError(t, err)
				require.Equal(t, 2, count)
				w.WriteHeader(http.StatusAccepted)
			}),
			log.NewNopLogger(),
		)
		v := &Viewer{
			Tenant: "test",
			User:   "user2",
			Role:   "",
		}
		ctx = NewContext(ctx, v)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithNoUserInViewer", func(t *testing.T) {
		h := UserHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			log.NewNopLogger(),
		)
		v := &Viewer{
			Tenant: "test",
			User:   "",
			Role:   "",
		}
		ctx = NewContext(ctx, v)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})
	t.Run("WithDeactivatedUserEntExists", func(t *testing.T) {
		h := UserHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			log.NewNopLogger(),
		)
		v := &Viewer{
			Tenant: "test",
			User:   "user1",
			Role:   "",
		}
		ctx = NewContext(ctx, v)
		err = client.User.UpdateOne(u).SetStatus(user.StatusDEACTIVATED).Exec(ctx)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestViewerMarshalLog(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	v := &Viewer{Tenant: "test", User: "tester"}
	logger.Info("viewer log test", zap.Object("viewer", v))

	logs := o.TakeAll()
	require.Len(t, logs, 1)
	field, ok := logs[0].ContextMap()["viewer"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, v.Tenant, field["tenant"])
	assert.Equal(t, v.User, field["user"])
}

type testExporter struct {
	mock.Mock
}

func (te *testExporter) ExportSpan(s *trace.SpanData) {
	te.Called(s)
}

func TestViewerSpanAttributes(t *testing.T) {
	h := TenancyHandler(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}),
		nil,
	)
	t.Run("WithSpan", func(t *testing.T) {
		var te testExporter
		trace.RegisterExporter(&te)
		defer trace.UnregisterExporter(&te)

		te.On("ExportSpan", mock.AnythingOfType("*trace.SpanData")).
			Run(func(args mock.Arguments) {
				s := args.Get(0).(*trace.SpanData)
				assert.Equal(t, "test", s.Attributes["viewer.tenant"])
				assert.Equal(t, "test", s.Attributes["viewer.user"])
			}).
			Once()
		defer te.AssertExpectations(t)

		ctx, span := trace.StartSpan(context.Background(), "test",
			trace.WithSampler(trace.AlwaysSample()),
		)
		defer span.End()

		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		req.Header.Set(TenantHeader, "test")
		req.Header.Set(UserHeader, "test")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithoutSpan", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(TenantHeader, "test")
		rec := httptest.NewRecorder()
		assert.NotPanics(t, func() { h.ServeHTTP(rec, req) })
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
}

func TestViewerTags(t *testing.T) {
	measure := stats.Int64("", "", stats.UnitDimensionless)
	v := &view.View{
		Name:        "viewer/test_tags",
		TagKeys:     []tag.Key{KeyTenant, KeyUser, KeyRole},
		Measure:     measure,
		Aggregation: view.Count(),
	}
	err := view.Register(v)
	require.NoError(t, err)
	defer view.Unregister(v)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(TenantHeader, "test-tenant")
	req.Header.Set(UserHeader, "test-user")
	req.Header.Set(RoleHeader, "readonly")
	rec := httptest.NewRecorder()
	TenancyHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := stats.RecordWithTags(r.Context(), nil, measure.M(1))
			require.NoError(t, err)
			w.WriteHeader(http.StatusNoContent)
		}),
		nil,
	).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	rows, err := view.RetrieveData(v.Name)
	require.NoError(t, err)
	require.Len(t, rows, 1)

	hasTag := func(key tag.Key, value string) assert.Comparison {
		return func() bool {
			for _, t := range rows[0].Tags {
				if t.Key.Name() == key.Name() {
					return t.Value == value
				}
			}
			return false
		}
	}
	assert.Condition(t, hasTag(KeyTenant, "test-tenant"))
	assert.Condition(t, hasTag(KeyUser, "test-user"))
	assert.Condition(t, hasTag(KeyRole, "readonly"))
}

func TestViewerTenancy(t *testing.T) {
	t.Run("WithTenancy", func(t *testing.T) {
		var client ent.Client
		h := TenancyHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.True(t, &client == ent.FromContext(r.Context()))
				w.WriteHeader(http.StatusAccepted)
			}),
			NewFixedTenancy(&client),
		)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(TenantHeader, "test")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithoutTenancy", func(t *testing.T) {
		h := TenancyHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Nil(t, ent.FromContext(r.Context()))
				w.WriteHeader(http.StatusAccepted)
			}),
			nil,
		)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(TenantHeader, "test")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
}
