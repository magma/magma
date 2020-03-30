// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directive

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// The following server GraphQL measures are supported for use in custom views.
var (
	ServerDeprecatedInputs = stats.Int64(
		"graphql/server/deprecated_inputs_count",
		"Number of GraphQL deprecated input fields",
		stats.UnitDimensionless,
	)
)

// The following tags are applied to stats recorded by this package.
var (
	// Field is the GraphQL object field being resolved.
	Field = tag.MustNewKey("graphql.field")
)

var (
	ServerDeprecatedCountByObjectInputField = &view.View{
		Name:        "graphql/server/deprecated_count_by_object_input_field",
		Description: "Count of GraphQL deprecated input fields by object and field",
		TagKeys:     []tag.Key{Field},
		Measure:     ServerDeprecatedInputs,
		Aggregation: view.Count(),
	}
)
