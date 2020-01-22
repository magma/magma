// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Hyperlink defines the hyperlink schema.
type Hyperlink struct {
	schema
}

// Fields returns hyperlink fields.
func (Hyperlink) Fields() []ent.Field {
	return []ent.Field{
		field.String("url"),
		field.String("name").
			StructTag(`gqlgen:"displayName"`).
			Optional(),
		field.String("category").
			Optional(),
	}
}
