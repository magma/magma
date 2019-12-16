// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct{}

// NewFieldsContext returns a new context with the given fields attached.
func NewFieldsContext(parent context.Context, fields ...zap.Field) context.Context {
	f := FieldsFromContext(parent)
	return context.WithValue(parent, contextKey{}, append(f[:len(f):len(f)], fields...))
}

// FieldsFromContext returns the fields stored in a context, or nil if there isn't one.
func FieldsFromContext(ctx context.Context) []zap.Field {
	f, _ := ctx.Value(contextKey{}).([]zap.Field)
	return f
}
