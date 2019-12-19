// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
	"github.com/facebookincubator/symphony/pkg/log"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// TenantLoader loads tenants by name.
type TenantLoader interface {
	Load(context.Context, string) (*ent.Tenant, error)
}

// The TenantLoaderFunc type is an adapter to allow the use of
// ordinary functions as tenant loaders.
type TenantLoaderFunc func(context.Context, string) (*ent.Tenant, error)

// Load calls f(ctx, name)
func (f TenantLoaderFunc) Load(ctx context.Context, name string) (*ent.Tenant, error) {
	return f(ctx, name)
}

// TenantClientLoader creates tenant loader from TenantClient.
func TenantClientLoader(client *ent.TenantClient, logger log.Logger) TenantLoader {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return TenantLoaderFunc(func(ctx context.Context, name string) (*ent.Tenant, error) {
		t, err := client.Query().Where(tenant.Name(name)).Only(ctx)
		if err != nil {
			logger.For(ctx).
				Error("loading tenant",
					zap.String("name", name),
					zap.Error(err),
				)
		}
		return t, err
	})
}

type tenantCtxKey struct{}

// TenantHandler returns a Handler that populates request tenant from subdomain.
func TenantHandler(handler http.Handler, loader TenantLoader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := domainutil.Subdomain(r.Host)
		if name == "" {
			http.Error(w, "missing tenant name", http.StatusBadRequest)
			return
		}
		if idx := strings.LastIndex(name, "."); idx != -1 {
			name = name[idx+1:]
		}
		ctx := r.Context()
		switch t, err := loader.Load(ctx, name); err.(type) {
		case nil:
			ctx = context.WithValue(ctx, tenantCtxKey{}, t)
			ctx = log.NewFieldsContext(ctx, zap.String("tenant", name))
			trace.FromContext(ctx).AddAttributes(
				trace.StringAttribute("tenant", name),
			)
			handler.ServeHTTP(w, r.WithContext(ctx))
		case *ent.ErrNotFound:
			http.Error(w, "tenant not found", http.StatusNotFound)
		default:
			http.Error(w, "cannot load tenant", http.StatusInternalServerError)
		}
	})
}

// CurrentTenant returns the Tenant stored in a context, or nil if there isn't one.
func CurrentTenant(ctx context.Context) *ent.Tenant {
	t, _ := ctx.Value(tenantCtxKey{}).(*ent.Tenant)
	return t
}
