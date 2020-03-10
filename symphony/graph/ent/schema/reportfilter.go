// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// ReportFilter defines the schema
type ReportFilter struct {
	schema
}

// Fields returns ReportFilter fields.
func (ReportFilter) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty(),
		field.Enum("entity").
			Values("WORK_ORDER", "PORT", "EQUIPMENT", "LINK", "LOCATION", "SERVICE"),
		field.Text("filters").Default("[]"),
	}
}
