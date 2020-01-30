// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestWithoutNameSampler(t *testing.T) {
	sampler := WithoutNameSampler("foo", "bar")
	decision := sampler(trace.SamplingParameters{Name: "foo"})
	assert.False(t, decision.Sample)
	decision = sampler(trace.SamplingParameters{Name: "bar"})
	assert.False(t, decision.Sample)
	decision = sampler(trace.SamplingParameters{Name: "baz"})
	assert.True(t, decision.Sample)
}
