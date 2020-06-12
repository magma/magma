// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"net/http"

	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/log"
	"go.uber.org/zap"
)

// Handler returns a Handler that runs h with permissions.
func Handler(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		permissions, err := Permissions(NewContext(ctx, AdminPermissions()))
		if err != nil {
			const msg = "cannot get permissions"
			logger.For(ctx).Warn(msg, zap.Error(err))
			http.Error(w, msg, http.StatusServiceUnavailable)
			return
		}
		ctx = NewContext(ctx, permissions)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextKey struct{}

// FromContext returns the PermissionSettings stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) *models.PermissionSettings {
	v, _ := ctx.Value(contextKey{}).(*models.PermissionSettings)
	return v
}

// NewContext returns a new context with the given PermissionSettings attached.
func NewContext(parent context.Context, p *models.PermissionSettings) context.Context {
	return context.WithValue(parent, contextKey{}, p)
}
