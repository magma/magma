// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tracer

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vektah/gqlparser/ast"
)

type testTracer struct {
	graphql.NopTracer
	mock.Mock
}

func (t *testTracer) StartFieldExecution(ctx context.Context, field graphql.CollectedField) context.Context {
	return t.Called(ctx, field).Get(0).(context.Context)
}

func TestTracerStartFieldExecution(t *testing.T) {
	var m testTracer
	defer m.AssertExpectations(t)
	ctx := context.Background()
	field := graphql.CollectedField{Field: &ast.Field{Name: "name"}}
	assert.Equal(t, ctx, tracer{&m}.StartFieldExecution(ctx, field))
}
