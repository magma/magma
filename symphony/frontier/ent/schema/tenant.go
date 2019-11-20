// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Tenant defines tenant schema.
type Tenant struct {
	schema
}

// Config configures tenant schema.
func (Tenant) Config() ent.Config {
	return ent.Config{
		Table: "Organizations",
	}
}

// Fields of tenant entity.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.Strings("domains").
			StorageKey("customDomains"),
		field.Strings("networks").
			StorageKey("networkIDs"),
		field.Strings("tabs").
			Optional(),
		field.Text("SSOCert").
			StorageKey("ssoCert").
			Default(""),
		field.String("SSOEntryPoint").
			StorageKey("ssoEntrypoint").
			Default(""),
		field.String("SSOIssuer").
			StorageKey("ssoIssuer").
			Default(""),
	}
}
