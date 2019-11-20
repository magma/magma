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
	"go.opencensus.io/trace"
)

type testTracer struct {
	graphql.NopTracer
	mock.Mock
}

func (t *testTracer) StartFieldExecution(ctx context.Context, field graphql.CollectedField) context.Context {
	return t.Called(ctx, field).Get(0).(context.Context)
}

func TestTracerStartFieldExecution(t *testing.T) {
	t.Run("user", func(t *testing.T) {
		var m testTracer
		m.On("StartFieldExecution", mock.Anything, mock.Anything).
			Return(context.Background()).
			Once()
		defer m.AssertExpectations(t)
		field := graphql.CollectedField{Field: &ast.Field{Name: "name"}}
		_ = tracer{&m}.StartFieldExecution(context.Background(), field)
	})
	t.Run("internal", func(t *testing.T) {
		var m testTracer
		defer m.AssertExpectations(t)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		ctx := tracer{&m}.StartFieldExecution(
			context.Background(),
			graphql.CollectedField{
				Field: &ast.Field{
					Name:             "__schema",
					ObjectDefinition: &ast.Definition{Name: "Query"},
				},
			},
		)
		assert.False(t, trace.FromContext(ctx).IsRecordingEvents())
	})
}
