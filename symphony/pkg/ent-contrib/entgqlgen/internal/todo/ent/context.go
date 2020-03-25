// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
)

type contextKey struct{}

// FromContext returns the Client stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) *Client {
	c, _ := ctx.Value(contextKey{}).(*Client)
	return c
}

// NewContext returns a new context with the given Client attached.
func NewContext(parent context.Context, c *Client) context.Context {
	return context.WithValue(parent, contextKey{}, c)
}
