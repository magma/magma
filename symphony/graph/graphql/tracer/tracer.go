// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tracer

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen-contrib/gqlopencensus"
	"github.com/99designs/gqlgen/graphql"
	"go.opencensus.io/trace"
)

type tracer struct{ graphql.Tracer }

// New creates a graphql traces
func New() graphql.Tracer {
	return tracer{Tracer: gqlopencensus.New()}
}

// StartFieldExecution blocks span creation on introspection fields
func (t tracer) StartFieldExecution(ctx context.Context, field graphql.CollectedField) context.Context {
	if strings.HasPrefix(field.Name, "__") {
		ctx, _ = trace.StartSpan(ctx,
			field.ObjectDefinition.Name+"/"+field.Name,
			trace.WithSampler(trace.NeverSample()),
		)
		return ctx
	}
	return t.Tracer.StartFieldExecution(ctx, field)
}
