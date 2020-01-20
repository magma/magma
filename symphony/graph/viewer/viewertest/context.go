// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewertest

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"
)

// Option enables viewer customization.
type Option func(*viewer.Viewer)

// WithTenant overrides viewer tenant name.
func WithTenant(name string) Option {
	return func(v *viewer.Viewer) {
		v.Tenant = name
	}
}

// NewContext returns viewer context for tests.
func NewContext(c *ent.Client, opts ...Option) context.Context {
	v := &viewer.Viewer{Tenant: "test"}
	for _, opt := range opts {
		opt(v)
	}
	ctx := viewer.NewContext(context.Background(), v)
	return ent.NewContext(ctx, c)
}
