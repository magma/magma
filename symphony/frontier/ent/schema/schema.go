// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

type timeMixin struct{}

func (timeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			StorageKey("createdAt").
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StorageKey("updatedAt").
			Immutable(),
	}
}

type schema struct {
	ent.Schema
}

func (schema) Mixin() []ent.Mixin {
	return []ent.Mixin{
		timeMixin{},
	}
}
