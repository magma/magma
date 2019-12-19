// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oc

import (
	"encoding/json"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/google/wire"
	"go.opencensus.io/trace"

	"github.com/jessevdk/go-flags"
)

type (
	// Options offers a declarative way to construct exporter options.
	Options struct {
		SamplingProbability float64 `env:"SAMPLING_PROBABILITY" long:"sampling-probability" default:"1.0" description:"Trace sampling probability" json:"sampling_probability"`
		Jaeger              *Jaeger `env:"JAEGER" long:"jaeger" description:"Jaeger exporter options as json" json:"jaeger"`
	}

	// Jaeger attaches flags methods to jaeger.Options.
	Jaeger jaeger.Options
)

// Set is a Wire provider set that produces exporter specific options given Options.
var Set = wire.NewSet(
	TraceSampler,
	JaegerOptions,
)

// TraceSampler returns trace sampler from options.
func TraceSampler(o Options) trace.Sampler {
	return trace.ProbabilitySampler(o.SamplingProbability)
}

// JaegerOptions returns jaeger options stored in Options.
func JaegerOptions(o Options) (opts jaeger.Options) {
	if o.Jaeger != nil {
		opts = jaeger.Options(*o.Jaeger)
	}
	return
}

// UnmarshalFlag implements flags.Unmarshaler interface.
func (j *Jaeger) UnmarshalFlag(value string) error {
	var opts jaeger.Options
	if err := json.Unmarshal([]byte(value), &opts); err != nil {
		return &flags.Error{
			Type:    flags.ErrMarshal,
			Message: err.Error(),
		}
	}
	*j = Jaeger(opts)
	return nil
}

var _ flags.Unmarshaler = (*Jaeger)(nil)
