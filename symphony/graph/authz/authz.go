// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/facebookincubator/symphony/pkg/log"
)

type AuthHandler struct {
	Handler http.Handler
	Logger  log.Logger
}

// AuthHandler calculates permissions of viewer and put in context.
func (h AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	permissions, err := Permissions(NewContext(ctx, AdminPermissions()))
	if err != nil {
		http.Error(w, fmt.Sprintf("get permissions: %s", err.Error()), http.StatusServiceUnavailable)
		return
	}
	ctx = NewContext(ctx, permissions)
	h.Handler.ServeHTTP(w, r.WithContext(ctx))
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
