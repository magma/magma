// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

// ReportFilter defines the schema
type ReportFilter struct {
	schema
}

// Fields returns ReportFilter fields.
func (ReportFilter) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.Enum("entity").
			Values("WORK_ORDER", "PORT", "EQUIPMENT", "LINK", "LOCATION", "SERVICE"),
		field.Text("filters").Default("[]"),
	}
}

// Indexes of the ReportFilter.
func (ReportFilter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "entity").
			Unique(),
	}
}

// Policy returns ReportFilter policy.
func (ReportFilter) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			privacy.AlwaysAllowRule(),
		),
	)
}
