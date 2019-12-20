// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tracer

import (
	"context"

	"github.com/99designs/gqlgen-contrib/gqlopencensus"
	"github.com/99designs/gqlgen/graphql"
)

type tracer struct{ graphql.Tracer }

// New creates a graphql traces
func New() graphql.Tracer {
	return tracer{Tracer: gqlopencensus.New()}
}

// StartFieldExecution blocks span creation on fields
func (tracer) StartFieldExecution(ctx context.Context, _ graphql.CollectedField) context.Context {
	return ctx
}
