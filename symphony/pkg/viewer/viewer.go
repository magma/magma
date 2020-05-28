// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/user"
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
	// UserHeader is the http user header.
	AutomationHeader = "x-auth-automation-name"
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

// Option enables viewer customization.
type Option func(Viewer)

// WithFeatures overrides default feature set.
func WithFeatures(features ...string) Option {
	return func(v Viewer) {
		switch x := v.(type) {
		case *UserViewer:
			x.features = NewFeatureSet(features...)
		case *AutomationViewer:
			x.features = NewFeatureSet(features...)
		}
	}
}

type viewer struct {
	tenant   string
	features FeatureSet
}

// Tenant is the tenant of the viewer.
func (v *viewer) Tenant() string {
	return v.tenant
}

// Features is the features applied for the viewer.
func (v *viewer) Features() FeatureSet {
	return v.features
}

// Viewer is the interface to hold additional per request information.
type Viewer interface {
	zapcore.ObjectMarshaler
	Tenant() string
	Features() FeatureSet
	Name() string
	Role() user.Role
}

// UserViewer is a viewer that holds a user ent.
type UserViewer struct {
	viewer
	user atomic.Value
}

// Name implements Viewer.Name by getting user's Auth ID.
func (v *UserViewer) Name() string {
	return v.User().AuthID
}

// Name implements Viewer.Name by getting user's Role.
func (v *UserViewer) Role() user.Role {
	return v.User().Role
}

// AutomationViewer is a viewer that holds information for the automation process.
type AutomationViewer struct {
	viewer
	name string
	role user.Role
}

// Name implements Viewer.Name by returning stored name.
func (v *AutomationViewer) Name() string {
	return v.name
}

// Role implements Viewer.Role by returning stored role.
func (v *AutomationViewer) Role() user.Role {
	return v.role
}

// NewFeatureSet create FeatureSet from a list of features.
func NewFeatureSet(features ...string) FeatureSet {
	set := make(FeatureSet, len(features))
	for _, feature := range features {
		set[feature] = struct{}{}
	}
	return set
}

// NewUser initializes and return UserViewer.
func NewUser(tenant string, user *ent.User, options ...Option) *UserViewer {
	v := &UserViewer{viewer: viewer{
		tenant: tenant,
	}}
	v.user.Store(user)
	for _, option := range options {
		option(v)
	}
	return v
}

// NewAutomation initializes and return AutomationViewer.
func NewAutomation(tenant, name string, role user.Role, options ...Option) *AutomationViewer {
	v := &AutomationViewer{
		viewer: viewer{
			tenant: tenant,
		},
		name: name,
		role: role}
	for _, option := range options {
		option(v)
	}
	return v
}

// String returns the textual representation of a feature set.
func (f FeatureSet) String() string {
	features := make([]string, 0, len(f))
	for feature := range f {
		features = append(features, feature)
	}
	sort.Strings(features)
	return strings.Join(features, ",")
}

// Enabled check if feature is in FeatureSet.
func (f FeatureSet) Enabled(feature string) bool {
	_, ok := f[feature]
	return ok
}

func (v *UserViewer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("tenant", v.Tenant())
	enc.AddString("user", v.Name())
	enc.AddString("role", v.Role().String())
	return nil
}

func (v *AutomationViewer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("tenant", v.Tenant())
	enc.AddString("user", v.Name())
	enc.AddString("role", v.Role().String())
	return nil
}

func traceAttrs(v Viewer) []trace.Attribute {
	return []trace.Attribute{
		trace.StringAttribute(TenantAttribute, v.Tenant()),
		trace.StringAttribute(UserAttribute, v.Name()),
		trace.StringAttribute(RoleAttribute, v.Role().String()),
	}
}

func tags(r *http.Request, v Viewer) []tag.Mutator {
	var userAgent string
	if parts := strings.SplitN(r.UserAgent(), " ", 2); len(parts) > 0 {
		userAgent = parts[0]
	}
	return []tag.Mutator{
		tag.Upsert(KeyTenant, v.Tenant()),
		tag.Upsert(KeyUser, v.Name()),
		tag.Upsert(KeyRole, v.Role().String()),
		tag.Upsert(KeyUserAgent, userAgent),
	}
}

// User returns the ent user of the viewer.
func (v *UserViewer) User() *ent.User {
	u, _ := v.user.Load().(*ent.User)
	return u
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

		var current struct {
			Tenant string `json:"organization"`
			User   string `json:"email"`
			Role   string `json:"role"`
		}
		if err := json.NewDecoder(rsp.Body).Decode(&current); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Header.Set(TenantHeader, current.Tenant)
		r.Header.Set(UserHeader, current.User)
		r.Header.Set(RoleHeader, current.Role)
		h.ServeHTTP(w, r)
	})
}

// TenancyHandler adds viewer / tenancy into incoming requests.
func TenancyHandler(h http.Handler, tenancy Tenancy, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.For(r.Context())
		tenant := r.Header.Get(TenantHeader)
		if tenant == "" {
			logger.Warn("request missing tenant header")
			http.Error(w, "missing tenant header", http.StatusBadRequest)
			return
		}
		logger = logger.With(zap.String("tenant", tenant))

		role := user.Role(r.Header.Get(RoleHeader))
		if err := user.RoleValidator(role); err != nil {
			logger.Warn("request contains invalid role",
				zap.Stringer("role", role),
				zap.Error(err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client, err := tenancy.ClientFor(r.Context(), tenant)
		if err != nil {
			const msg = "cannot get tenancy client"
			logger.Warn(msg, zap.Error(err))
			http.Error(w, msg, http.StatusServiceUnavailable)
			return
		}

		var (
			ctx  = ent.NewContext(r.Context(), client)
			opts = make([]Option, 0, 1)
			v    Viewer
		)
		if features := r.Header.Get(FeaturesHeader); features != "" {
			opts = append(opts, WithFeatures(strings.Split(features, ",")...))
		}
		if username := r.Header.Get(UserHeader); username != "" {
			u, err := GetOrCreateUser(
				privacy.DecisionContext(ctx, privacy.Allow),
				username,
				role)
			if err != nil {
				logger.Warn("cannot get user ent", zap.Error(err))
				http.Error(w, "getting user entity", http.StatusServiceUnavailable)
				return
			}
			v = NewUser(tenant, u, opts...)
		} else if automation := r.Header.Get(AutomationHeader); automation != "" {
			v = NewAutomation(tenant, automation, role, opts...)
		} else {
			logger.Warn("request missing identity header")
			http.Error(w, "missing identity header", http.StatusBadRequest)
			return
		}

		ctx = log.NewFieldsContext(ctx, zap.Object("viewer", v))
		trace.FromContext(ctx).AddAttributes(traceAttrs(v)...)
		ctx, _ = tag.New(ctx, tags(r, v)...)
		ctx = NewContext(ctx, v)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetOrCreateUser creates or returns existing user with given authID and role.
func GetOrCreateUser(ctx context.Context, authID string, role user.Role) (*ent.User, error) {
	client := ent.FromContext(ctx).User
	if u, err := client.Query().
		Where(user.AuthID(authID)).
		Only(ctx); err == nil || !ent.IsNotFound(err) {
		return u, err
	}
	u, err := client.Create().
		SetAuthID(authID).
		SetEmail(authID).
		SetRole(role).
		Save(ctx)
	if ent.IsConstraintError(err) {
		u, err = client.Query().
			Where(user.AuthID(authID)).
			Only(ctx)
	}
	return u, err
}

// MustGetOrCreateUser creates or returns existing user ent with given authID and role
func MustGetOrCreateUser(ctx context.Context, authID string, role user.Role) *ent.User {
	u, err := GetOrCreateUser(ctx, authID, role)
	if err != nil {
		panic(err)
	}
	return u
}

// UserHandler blocks deactivated users.
func UserHandler(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v, ok := FromContext(r.Context()).(*UserViewer); ok {
			if u := v.User(); u.Status == user.StatusDEACTIVATED {
				logger.For(r.Context()).Debug("blocking deactivated user",
					zap.String("authid", u.AuthID),
					zap.String("email", u.Email),
				)
				http.Error(w, "user is deactivated", http.StatusForbidden)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

type contextKey struct{}

// FromContext returns the Viewer stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) Viewer {
	v, _ := ctx.Value(contextKey{}).(Viewer)
	return v
}

// NewContext returns a new context with the given Viewer attached.
func NewContext(parent context.Context, v Viewer) context.Context {
	return context.WithValue(parent, contextKey{}, v)
}
