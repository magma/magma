// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/gorilla/websocket"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// TenantHeader is the http tenant header.
	TenantHeader = "x-auth-organization"
	// UserHeader is the http user header.
	UserHeader = "x-auth-user-email"
	// RoleHeader is the http role header.
	RoleHeader = "x-auth-user-role"
)

// Attributes recorded on the span of the requests.
const (
	TenantAttribute    = "viewer.tenant"
	UserAttribute      = "viewer.user"
	RoleAttribute      = "viewer.role"
	UserAgentAttribute = "viewer.user_agent"
)

const (
	UserIsDeactivatedError = "USER_IS_DEACTIVATED"
)

// The following tags are applied to context recorded by this package.
var (
	KeyTenant    = tag.MustNewKey(TenantAttribute)
	KeyUser      = tag.MustNewKey(UserAttribute)
	KeyRole      = tag.MustNewKey(RoleAttribute)
	KeyUserAgent = tag.MustNewKey(UserAgentAttribute)
)

// Viewer holds additional per request information.
type Viewer struct {
	Tenant string `json:"organization"`
	User   string `json:"email"`
	Role   string `json:"role"`
}

// MarshalLogObject implements zapcore.ObjectMarshaler interface.
func (v *Viewer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("tenant", v.Tenant)
	enc.AddString("user", v.User)
	enc.AddString("role", v.Role)
	return nil
}

func (v *Viewer) traceAttrs() []trace.Attribute {
	return []trace.Attribute{
		trace.StringAttribute(TenantAttribute, v.Tenant),
		trace.StringAttribute(UserAttribute, v.User),
		trace.StringAttribute(RoleAttribute, v.Role),
	}
}

func (v *Viewer) tags(r *http.Request) []tag.Mutator {
	return []tag.Mutator{
		tag.Upsert(KeyTenant, v.Tenant),
		tag.Upsert(KeyUser, v.User),
		tag.Upsert(KeyRole, v.Role),
		tag.Upsert(KeyUserAgent, r.UserAgent()),
	}
}

// WebSocketUpgradeHandler authenticates websocket upgrade requests.
func WebSocketUpgradeHandler(h http.Handler, authurl string) http.Handler {
	client := &http.Client{
		Transport: &ochttp.Transport{
			FormatSpanName: func(*http.Request) string {
				return "viewer.authenticate"
			},
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r) || r.Header.Get(TenantHeader) != "" {
			h.ServeHTTP(w, r)
			return
		}
		req, err := http.NewRequestWithContext(
			r.Context(), http.MethodGet, authurl, nil,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("X-Forwarded-Host", r.Host)
		if username, password, ok := r.BasicAuth(); ok {
			req.SetBasicAuth(username, password)
		}
		for _, c := range r.Cookies() {
			req.AddCookie(c)
		}

		rsp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		var v Viewer
		if err := json.NewDecoder(rsp.Body).Decode(&v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Header.Set(TenantHeader, v.Tenant)
		r.Header.Set(UserHeader, v.User)
		r.Header.Set(RoleHeader, v.Role)
		h.ServeHTTP(w, r)
	})
}

// TenancyHandler adds viewer / tenancy into incoming requests.
func TenancyHandler(h http.Handler, tenancy Tenancy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := r.Header.Get(TenantHeader)
		if tenant == "" {
			http.Error(w, "missing tenant header", http.StatusBadRequest)
			return
		}

		v := &Viewer{
			Tenant: tenant,
			User:   r.Header.Get(UserHeader),
			Role:   r.Header.Get(RoleHeader),
		}

		ctx := log.NewFieldsContext(r.Context(), zap.Object("viewer", v))
		trace.FromContext(ctx).AddAttributes(v.traceAttrs()...)
		ctx, _ = tag.New(ctx, v.tags(r)...)

		ctx = NewContext(ctx, v)
		if tenancy != nil {
			client, err := tenancy.ClientFor(ctx, tenant)
			if err != nil {
				http.Error(w, "getting tenancy client", http.StatusServiceUnavailable)
				return
			}
			if v.Role == "readonly" {
				client = client.ReadOnly()
			}
			ctx = ent.NewContext(ctx, client)
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserHandler adds users if request is from user that is not found.
func UserHandler(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v := FromContext(ctx)
		client := ent.FromContext(ctx)
		u, err := client.User.Query().Where(user.AuthID(v.User)).Only(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				http.Error(w, "query user ent", http.StatusServiceUnavailable)
				return
			}
			_, err := client.User.Create().SetAuthID(v.User).Save(ctx)
			if err != nil {
				http.Error(w, "create user ent", http.StatusServiceUnavailable)
				return
			}
			logger.For(ctx).Info("New user created", zap.String("AuthID", v.User))
		} else if u.Status == user.StatusDEACTIVATED {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = fmt.Fprintln(
				w, "{\"errorCode\": \"USER_IS_DEACTIVATED\", \"description\"Error struct: \"User must be active to see this\"}")
			return
		}
		h.ServeHTTP(w, r)
	})
}

type contextKey struct{}

// FromContext returns the Viewer stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) *Viewer {
	v, _ := ctx.Value(contextKey{}).(*Viewer)
	return v
}

// NewContext returns a new context with the given Viewer attached.
func NewContext(parent context.Context, v *Viewer) context.Context {
	return context.WithValue(parent, contextKey{}, v)
}
