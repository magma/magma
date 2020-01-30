// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"go.opencensus.io/trace"
)

// WithoutNameSampler returns a trace sampler filtering out a set of span names.
func WithoutNameSampler(name string, names ...string) trace.Sampler {
	return func(params trace.SamplingParameters) trace.SamplingDecision {
		if params.Name == name {
			return trace.SamplingDecision{Sample: false}
		}
		for _, name := range names {
			if params.Name == name {
				return trace.SamplingDecision{Sample: false}
			}
		}
		return trace.SamplingDecision{Sample: true}
	}
}
