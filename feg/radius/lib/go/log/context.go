/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

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
