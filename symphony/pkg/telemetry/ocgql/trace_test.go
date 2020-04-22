// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql_test

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/facebookincubator/symphony/pkg/telemetry/ocgql"
	"github.com/stretchr/testify/suite"
	"go.opencensus.io/trace"
)

type tracerTestSuite struct {
	suite.Suite
	sample bool
	client *client.Client
	spans  map[string]*trace.SpanData
}

func (s *tracerTestSuite) SetupSuite() {
	srv := testserver.New()
	srv.AddTransport(transport.POST{})
	srv.Use(extension.FixedComplexityLimit(1000))
	srv.SetCalculatedComplexity(100)
	srv.Use(ocgql.Tracer{
		AllowRoot: true,
		Field:     true,
		DefaultAttributes: []trace.Attribute{
			trace.BoolAttribute("graphql.test.value", true),
		},
		Sampler: func(trace.SamplingParameters) trace.SamplingDecision {
			return trace.SamplingDecision{Sample: s.sample}
		},
	})
	srv.AroundOperations(
		func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			ctx, span := trace.StartSpan(ctx, "test",
				trace.WithSampler(trace.AlwaysSample()),
			)
			defer span.End()
			return next(ctx)
		},
	)
	s.client = client.New(srv)
	trace.RegisterExporter(s)
}

func (s *tracerTestSuite) TearDownSuite() {
	trace.UnregisterExporter(s)
}

func (s *tracerTestSuite) SetupTest() {
	s.sample = true
	s.spans = map[string]*trace.SpanData{}
}

func (s *tracerTestSuite) ExportSpan(span *trace.SpanData) {
	s.spans[span.Name] = span
}

func TestTracer(t *testing.T) {
	suite.Run(t, &tracerTestSuite{})
}

func (s *tracerTestSuite) TestWithSampling() {
	const (
		query = "query($id: Int!) { name: find(id: $id) }"
		id    = "42"
	)
	err := s.post(query, client.Var("id", id))
	s.Require().NoError(err)
	s.Assert().Len(s.spans, 3)

	span, ok := s.spans["query"]
	s.Require().True(ok)
	s.Assert().EqualValues(trace.SpanKindServer, span.SpanKind)
	s.Assert().Equal(query, span.Attributes["graphql.query"])
	s.Assert().Equal(id, span.Attributes["graphql.vars.id"])
	s.Assert().EqualValues(100, span.Attributes["graphql.complexity.value"])
	s.Assert().EqualValues(1000, span.Attributes["graphql.complexity.limit"])
	s.Assert().EqualValues(trace.StatusCodeOK, span.Code)
	s.Assert().Empty(span.Message)

	span, ok = s.spans["name"]
	s.Require().True(ok)
	for _, attr := range []string{"path", "name", "alias"} {
		s.Assert().Equal("name", span.Attributes["graphql.field."+attr])
	}
	s.Assert().EqualValues(trace.StatusCodeOK, span.Code)
	s.Assert().Empty(span.Message)
}

func (s *tracerTestSuite) TestWithoutSampling() {
	s.sample = false
	err := s.post("query { name }")
	s.Require().NoError(err)
	_, ok := s.spans["query"]
	s.Assert().False(ok)
}

func (s *tracerTestSuite) TestUnsupportedOperation() {
	const query = "mutation { name }"
	err := s.post(query)
	s.Require().Error(err)

	span, ok := s.spans["mutation"]
	s.Require().True(ok)
	s.Assert().EqualValues(trace.SpanKindServer, span.SpanKind)
	s.Assert().Equal(query, span.Attributes["graphql.query"])
	s.Assert().EqualValues(trace.StatusCodeUnknown, span.Code)
	const message = "mutations are not supported"
	s.Assert().Contains(span.Message, message)
	s.Assert().Contains(err.Error(), message)
}

func (s *tracerTestSuite) TestDefaultAttributes() {
	err := s.post("query { name }")
	s.Require().NoError(err)

	span, ok := s.spans["query"]
	s.Require().True(ok)
	s.Assert().Equal(true, span.Attributes["graphql.test.value"])
}

func (s *tracerTestSuite) post(query string, opts ...client.Option) error {
	return s.client.Post(query, &struct{ Name string }{}, opts...)
}
