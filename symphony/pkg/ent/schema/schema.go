// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/mixin"
	"github.com/facebookincubator/symphony/pkg/authz"
)

// schema adds time mixin to underlying ents.
type schema struct {
	ent.Schema
}

// Mixin returns schema mixins.
func (schema) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Policy returns schema policy.
func (schema) Policy() ent.Policy {
	return authz.NewPolicy()
}
