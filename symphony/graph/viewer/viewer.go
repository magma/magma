// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

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
	// FeaturesHeader is the http feature header.
	FeaturesHeader = "x-auth-features"
)

// Attributes recorded on the span of the requests.
const (
	TenantAttribute    = "viewer.tenant"
	UserAttribute      = "viewer.user"
	RoleAttribute      = "viewer.role"
	UserAgentAttribute = "viewer.user_agent"
)

// The following tags are applied to context recorded by this package.
var (
	KeyTenant    = tag.MustNewKey(TenantAttribute)
	KeyUser      = tag.MustNewKey(UserAttribute)
	KeyRole      = tag.MustNewKey(RoleAttribute)
	KeyUserAgent = tag.MustNewKey(UserAgentAttribute)
)

// FeatureSet holds the list of features of the viewer
type FeatureSet map[string]struct{}

// Viewer holds additional per request information.
type Viewer struct {
	Tenant   string     `json:"organization"`
	User     string     `json:"email"`
	Role     string     `json:"role"`
	Features FeatureSet `json:"-"`
}

type UserHandler struct {
	Handler http.Handler
	Logger  log.Logger
}

// NewFeatureSet create FeatureSet from a list of features.
func NewFeatureSet(features ...string) FeatureSet {
	set := make(FeatureSet, len(features))
	for _, feature := range features {
		set[feature] = struct{}{}
	}
	return set
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
	var userAgent string
	if parts := strings.SplitN(r.UserAgent(), " ", 2); len(parts) > 0 {
		userAgent = parts[0]
	}
	return []tag.Mutator{
		tag.Upsert(KeyTenant, v.Tenant),
		tag.Upsert(KeyUser, v.User),
		tag.Upsert(KeyRole, v.Role),
		tag.Upsert(KeyUserAgent, userAgent),
	}
}

// Enabled check if feature is in FeatureSet.
func (f FeatureSet) Enabled(feature string) bool {
	_, ok := f[feature]
	return ok
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
			Tenant:   tenant,
			User:     r.Header.Get(UserHeader),
			Role:     r.Header.Get(RoleHeader),
			Features: NewFeatureSet(strings.Split(r.Header.Get(FeaturesHeader), ",")...),
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
			ctx = ent.NewContext(ctx, client)
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h UserHandler) getOrCreateUser(ctx context.Context) (*ent.User, error) {
	ctx, span := trace.StartSpan(ctx, "viewer.getOrCreateUser")
	defer span.End()
	v := FromContext(ctx)
	client := ent.FromContext(ctx)

	u, err := client.User.Query().Where(user.AuthID(v.User)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}
		role := user.RoleUSER
		if v.Role == "superuser" {
			role = user.RoleOWNER
		}
		u, err = client.User.Create().SetAuthID(v.User).SetEmail(v.User).SetRole(role).Save(ctx)
		if err != nil {
			if !ent.IsConstraintError(err) {
				return nil, err
			}
			return client.User.Query().Where(user.AuthID(v.User)).Only(ctx)
		}
		h.Logger.For(ctx).Info("New user created", zap.String("AuthID", v.User))
	}
	return u, err
}

// UserHandler adds users if request is from user that is not found.
func (h UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, err := h.getOrCreateUser(ctx)
	if err != nil {
		http.Error(w, "get user ent", http.StatusServiceUnavailable)
		return
	}
	if u.Status == user.StatusDEACTIVATED {
		http.Error(w, "user is deactivated", http.StatusForbidden)
		return
	}
	// TODO(T64743627): Stop checking read only
	readOnly, err := IsUserReadOnly(ctx, u)
	if err != nil {
		http.Error(w, "check is read only", http.StatusServiceUnavailable)
		return
	}
	if readOnly {
		client := ent.FromContext(ctx).ReadOnly()
		ctx = ent.NewContext(ctx, client)
	}
	h.Handler.ServeHTTP(w, r.WithContext(ctx))
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
