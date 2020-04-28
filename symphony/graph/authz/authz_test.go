package authz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/stretchr/testify/require"
)

func TestAuthHandler(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	u := viewer.MustGetOrCreateUser(ctx, viewertest.DefaultUser, viewer.SuperUserRole)
	v := viewer.New(viewertest.DefaultTenant, u)
	ctx = viewer.NewContext(ctx, v)
	h := authz.AuthHandler{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions := authz.FromContext(r.Context())
			require.NotNil(t, permissions)
			require.EqualValues(t, authz.FullPermissions(), permissions)
			w.WriteHeader(http.StatusAccepted)
		}),
		Logger: log.NewNopLogger(),
	}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req.WithContext(ctx))
	require.Equal(t, http.StatusAccepted, rec.Code)
}
